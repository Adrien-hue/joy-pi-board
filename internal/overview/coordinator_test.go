package overview

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/healthschema"
)

type fakeTimer struct {
	ch       chan time.Time
	clock    *fakeClock
	deadline time.Duration
	stopped  bool
}

func (t *fakeTimer) C() <-chan time.Time { return t.ch }
func (t *fakeTimer) Stop() bool {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	old := t.stopped
	t.stopped = true
	return !old
}

type fakeClock struct {
	mu      sync.Mutex
	reading Reading
	timers  []*fakeTimer
}

func newFake() *fakeClock {
	return &fakeClock{reading: Reading{time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC), 0, true}}
}
func (c *fakeClock) Now() Reading { c.mu.Lock(); defer c.mu.Unlock(); return c.reading }
func (c *fakeClock) Until(d time.Duration) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &fakeTimer{make(chan time.Time, 1), c, d, false}
	c.timers = append(c.timers, t)
	return t
}
func (c *fakeClock) advance(tick, wall time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reading.Tick += tick
	c.reading.Wall = c.reading.Wall.Add(wall)
	for _, t := range c.timers {
		if !t.stopped && c.reading.Tick >= t.deadline {
			t.stopped = true
			t.ch <- c.reading.Wall
		}
	}
}

type fetchFunc func(context.Context) health.Result

func (f fetchFunc) Fetch(ctx context.Context) health.Result { return f(ctx) }

type recordedEvents struct {
	mu     sync.Mutex
	events []Event
}

func (e *recordedEvents) Emit(v Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.events = append(e.events, v)
}
func snapshot(t *testing.T, example string) healthschema.Validated {
	t.Helper()
	b, e := os.ReadFile("../../docs/examples/overview-" + example + ".json")
	if e != nil {
		t.Fatal(e)
	}
	var v struct {
		Health struct {
			Snapshot json.RawMessage `json:"snapshot"`
		} `json:"health"`
	}
	if e = json.Unmarshal(b, &v); e != nil {
		t.Fatal(e)
	}
	s, e := healthschema.Decode(v.Health.Snapshot)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func demand(t *testing.T, c *Coordinator) View {
	t.Helper()
	v, e := c.Demand(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func wire(t *testing.T, v View, clock Clock) (map[string]any, Response) {
	t.Helper()
	p, e := v.Prepare()
	if e != nil {
		t.Fatal(e)
	}
	r, e := p.Finalize(clock.Now())
	if e != nil {
		t.Fatal(e)
	}
	var body map[string]any
	dec := json.NewDecoder(bytes.NewReader(r.Body))
	dec.UseNumber()
	if e = dec.Decode(&body); e != nil {
		t.Fatal(e)
	}
	return body["health"].(map[string]any), r
}
func until(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("barrier watchdog")
		}
		runtime.Gosched()
	}
}
func waiters(c *Coordinator, n int) bool { c.mu.Lock(); defer c.mu.Unlock(); return c.waiters == n }

// S02-T05/T08: every reason, partial replacement, monotonic boundaries, recovery.
func TestStateTimeAndReplacement(t *testing.T) {
	full, partial := snapshot(t, "current"), snapshot(t, "partial")
	clock := newFake()
	next := health.Result{Reason: health.ConnectionFailed}
	calls := 0
	c := New(fetchFunc(func(context.Context) health.Result { calls++; return next }), clock, nil)
	defer c.Stop()
	h, r := wire(t, demand(t, c), clock)
	if h["snapshot"] != nil || h["last_success_at"] != nil || r.Age != "" {
		t.Fatal(h, r.Age)
	}
	next = health.Result{Snapshot: full}
	v := demand(t, c)
	h, r = wire(t, v, clock)
	if h["snapshot_state"] != "current" || h["reason"] != nil || r.Age != "0" {
		t.Fatal(h, r.Age)
	}
	last := h["last_success_at"]
	for _, reason := range []health.Reason{health.Timeout, health.ConnectionFailed, health.UpstreamError, health.InvalidResponse, health.UnsupportedSchemaVersion} {
		next = health.Result{Reason: reason}
		h, _ = wire(t, demand(t, c), clock)
		if h["snapshot_state"] != "stale" || h["reason"] != string(reason) || h["last_success_at"] != last {
			t.Fatal(h)
		}
	}
	for _, step := range []struct {
		delta      time.Duration
		state, age string
	}{{29999 * time.Millisecond, "stale", "29999"}, {time.Millisecond, "stale", "30000"}, {time.Millisecond, "none", "30001"}} {
		clock.advance(step.delta, -24*time.Hour)
		h, r = wire(t, demand(t, c), clock)
		if h["snapshot_state"] != step.state || r.Age != step.age || h["last_success_at"] != last {
			t.Fatal(h, r.Age)
		}
		if (h["snapshot"] == nil) != (step.state == "none") {
			t.Fatal("expired metric exposed")
		}
	}
	next = health.Result{Snapshot: partial}
	h, _ = wire(t, demand(t, c), clock)
	if h["snapshot_state"] != "current" || h["snapshot"].(map[string]any)["memory"].(map[string]any)["total_bytes"] != nil {
		t.Fatal("partial did not replace complete", h)
	}
	before := calls
	clock.advance(24*time.Hour, 48*time.Hour)
	if calls != before {
		t.Fatal("autonomous call")
	}
	next = health.Result{Reason: health.ConnectionFailed}
	restart := New(c.fetch, clock, nil)
	defer restart.Stop()
	h, _ = wire(t, demand(t, restart), clock)
	if h["last_success_at"] != nil || h["snapshot"] != nil {
		t.Fatal("restart cache")
	}
}

// S02-T06/T07: deterministic ten waiters, cancellation and fresh successive calls.
func TestActiveSharingAndCancellation(t *testing.T) {
	clock := newFake()
	release := make(chan struct{})
	entered := make(chan struct{}, 4)
	var calls atomic.Int32
	full := snapshot(t, "current")
	c := New(fetchFunc(func(context.Context) health.Result {
		calls.Add(1)
		entered <- struct{}{}
		<-release
		return health.Result{Snapshot: full}
	}), clock, nil)
	defer c.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	first := make(chan error, 1)
	go func() { _, e := c.Demand(ctx); first <- e }()
	<-entered
	views := make(chan View, 9)
	for i := 0; i < 9; i++ {
		go func() {
			v, e := c.Demand(context.Background())
			if e != nil {
				t.Error(e)
			}
			views <- v
		}()
	}
	until(t, func() bool { return waiters(c, 10) })
	cancel()
	if !errors.Is(<-first, context.Canceled) {
		t.Fatal("cancellation")
	}
	if calls.Load() != 1 {
		t.Fatal(calls.Load())
	}
	close(release)
	for i := 0; i < 9; i++ {
		v := <-views
		if v.sequence != 1 || !v.cache.snapshot.Valid() {
			t.Fatal(v)
		}
	}
	v := demand(t, c)
	if calls.Load() != 2 || v.sequence != 2 {
		t.Fatal("reused completed flight")
	}
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	_, _ = c.Demand(ctx)
	if calls.Load() != 2 {
		t.Fatal("cancelled caller started work")
	}
}

func TestOrphanTimeoutLateWorkerAndSaturation(t *testing.T) {
	clock := newFake()
	release := make(chan struct{})
	entered := make(chan context.Context, 33)
	full := snapshot(t, "current")
	c := New(fetchFunc(func(ctx context.Context) health.Result {
		entered <- ctx
		<-release
		return health.Result{Snapshot: full}
	}), clock, nil)
	defer func() {
		close(release)
		c.Stop()
		until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.workers == 0 })
	}()
	for i := 0; i < 32; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { _, e := c.Demand(ctx); done <- e }()
		work := <-entered
		cancel()
		if !errors.Is(<-done, context.Canceled) {
			t.Fatal("detach")
		}
		if work.Err() != nil {
			t.Fatal("orphan cancelled before its deadline")
		}
		clock.advance(time.Second, 0)
		until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.active == nil })
		until(t, func() bool { return work.Err() != nil })
	}
	if _, e := c.Demand(context.Background()); !errors.Is(e, ErrUnavailable) {
		t.Fatal("worker saturation", e)
	}
	if c.cache.snapshot.Valid() {
		t.Fatal("late success published")
	}
}

func TestOrphanCanBeRejoinedAndStopDoesNotInventReason(t *testing.T) {
	clock := newFake()
	entered := make(chan struct{})
	release := make(chan struct{})
	full := snapshot(t, "current")
	c := New(fetchFunc(func(context.Context) health.Result { close(entered); <-release; return health.Result{Snapshot: full} }), clock, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { _, e := c.Demand(ctx); done <- e }()
	<-entered
	cancel()
	<-done
	joined := make(chan error, 1)
	go func() { _, e := c.Demand(context.Background()); joined <- e }()
	until(t, func() bool { return waiters(c, 1) })
	c.Stop()
	if !errors.Is(<-joined, ErrUnavailable) {
		t.Fatal("shutdown")
	}
	close(release)
	until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.workers == 0 })
	if c.cache.snapshot.Valid() || c.sequence != 1 || c.transition != "" {
		t.Fatal("shutdown published provider result")
	}
}

func TestDeadlinePriorityAndImmutableViews(t *testing.T) {
	full, partial := snapshot(t, "current"), snapshot(t, "partial")
	for _, candidate := range []health.Result{{Snapshot: full}, {Reason: health.UpstreamError}, {Reason: health.InvalidResponse}, {Reason: health.UnsupportedSchemaVersion}} {
		clock := newFake()
		c := New(fetchFunc(func(context.Context) health.Result { clock.advance(time.Second, 0); return candidate }), clock, nil)
		v := demand(t, c)
		if v.reason != health.Timeout || v.cache.snapshot.Valid() {
			t.Fatal("equal deadline must win", v.reason)
		}
		c.Stop()
	}
	clock := newFake()
	next := health.Result{Snapshot: full}
	c := New(fetchFunc(func(context.Context) health.Result { return next }), clock, nil)
	defer c.Stop()
	demand(t, c)
	next = health.Result{Reason: health.UpstreamError}
	old := demand(t, c)
	clock.advance(time.Millisecond, time.Hour)
	next = health.Result{Snapshot: partial}
	newView := demand(t, c)
	oldBytes := old.cache.snapshot.Bytes()
	oldBytes[0] = '!'
	h, _ := wire(t, old, clock)
	if h["reason"] != "upstream_error" || h["snapshot"].(map[string]any)["memory"].(map[string]any)["total_bytes"] == nil {
		t.Fatal("mixed views")
	}
	h, _ = wire(t, newView, clock)
	if h["snapshot_state"] != "current" {
		t.Fatal(h)
	}
	clock.advance(30*time.Second, 0)
	p, e := newView.Prepare()
	if e != nil {
		t.Fatal(e)
	}
	// exactly 30s remains current; +1ns cannot be delivered as current.
	if _, e = p.Finalize(clock.Now()); e != nil {
		t.Fatal(e)
	}
	clock.advance(time.Nanosecond, 0)
	if _, e = p.Finalize(clock.Now()); !errors.Is(e, ErrUnavailable) {
		t.Fatal(e)
	}
	h, r := wire(t, old, clock)
	if h["snapshot"] != nil || r.Age != "30002" {
		t.Fatal(h, r.Age)
	}
}

func TestResponseGuardsAndAgeArithmetic(t *testing.T) {
	if _, err := timestamp(time.Date(9999, 12, 31, 23, 0, 0, 0, time.FixedZone("test", -14*60*60))); err == nil {
		t.Fatal("UTC year overflow accepted")
	}
	clock := newFake()
	full := snapshot(t, "current")
	c := New(fetchFunc(func(context.Context) health.Result { return health.Result{Snapshot: full} }), clock, nil)
	defer c.Stop()
	v := demand(t, c)
	p, e := v.Prepare()
	if e != nil {
		t.Fatal(e)
	}
	for _, now := range []Reading{{Wall: clock.Now().Wall, Tick: -1, Valid: true}, {Wall: clock.Now().Wall, Valid: false}, {Wall: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}} {
		if _, e = p.Finalize(now); !errors.Is(e, ErrClock) {
			t.Fatal(e)
		}
	}
	for _, n := range []time.Duration{0, 1, time.Millisecond, time.Millisecond + 1, 30 * time.Second} {
		now := clock.Now()
		now.Tick = n
		r, e := p.Finalize(now)
		if e != nil {
			t.Fatal(e)
		}
		want := (uint64(n) / 1000000)
		if n%time.Millisecond != 0 {
			want++
		}
		if r.Age != fmt.Sprint(want) {
			t.Fatal(r.Age, want)
		}
	}
	if capMilliseconds(^uint64(0)) != "9007199254740991" {
		t.Fatal("saturation")
	}
	if _, e = (View{}).Prepare(); e == nil {
		t.Fatal("zero view")
	}
}

func TestTransitionLogsAndLateResultDiscard(t *testing.T) {
	clock := newFake()
	events := new(recordedEvents)
	full, partial := snapshot(t, "current"), snapshot(t, "partial")
	next := health.Result{Reason: health.ConnectionFailed}
	c := New(fetchFunc(func(context.Context) health.Result { return next }), clock, events)
	defer c.Stop()
	for i := 0; i < 3; i++ {
		demand(t, c)
	}
	next = health.Result{Snapshot: full}
	for i := 0; i < 3; i++ {
		demand(t, c)
	}
	next = health.Result{Snapshot: partial}
	demand(t, c)
	until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.workers == 0 })
	events.mu.Lock()
	defer events.mu.Unlock()
	if len(events.events) != 3 {
		t.Fatal(events.events)
	}
	for _, e := range events.events {
		if strings.Contains(e.Issues, "metric") || strings.Contains(e.Issues, "hostname") {
			t.Fatal("untrusted log data")
		}
	}
}

// S02-T04: real sockets + the actual coordinator timer (not fake clock).
func TestRealDeadlineHeadersBodyAndProgress(t *testing.T) {
	for _, mode := range []string{"headers", "body", "progress"} {
		t.Run(mode, func(t *testing.T) {
			cancelled := make(chan struct{})
			entered := make(chan struct{})
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(cancelled)
				close(entered)
				if mode != "headers" {
					w.Header().Set("Content-Type", "application/json")
					w.(http.Flusher).Flush()
				}
				if mode == "progress" {
					ticker := time.NewTicker(100 * time.Millisecond)
					defer ticker.Stop()
					for {
						select {
						case <-r.Context().Done():
							return
						case <-ticker.C:
							_, _ = io.WriteString(w, " ")
							w.(http.Flusher).Flush()
						}
					}
				}
				<-r.Context().Done()
			}))
			defer s.Close()
			client := health.New(s.URL)
			defer client.Close()
			c := New(client, nil, nil)
			defer c.Stop()
			start := time.Now()
			v := demand(t, c)
			elapsed := time.Since(start)
			if v.reason != health.Timeout || elapsed < 900*time.Millisecond || elapsed > 2*time.Second {
				t.Fatal(v.reason, elapsed)
			}
			select {
			case <-cancelled:
			case <-time.After(time.Second):
				t.Fatal("socket not cancelled")
			}
			<-entered
			t.Logf("real %s deadline=%s; one loopback attempt, workstation only", mode, elapsed)
		})
	}
}

func TestValidationAnchorAndLatePublication(t *testing.T) {
	full, partial := snapshot(t, "current"), snapshot(t, "partial")
	for _, timeout := range []bool{false, true} {
		t.Run(fmt.Sprint(timeout), func(t *testing.T) {
			clock := newFake()
			entered, release := make(chan struct{}), make(chan struct{})
			var count atomic.Int32
			c := New(fetchFunc(func(context.Context) health.Result {
				if count.Add(1) == 1 {
					clock.advance(200*time.Millisecond, 0)
					return health.Result{Snapshot: full}
				}
				return health.Result{Snapshot: partial}
			}), clock, nil)
			defer c.Stop()
			c.beforePublish = func(seq uint64) {
				if seq == 1 {
					close(entered)
					<-release
				}
			}
			views := make(chan View, 1)
			go func() {
				v, e := c.Demand(context.Background())
				if e != nil {
					t.Error(e)
				}
				views <- v
			}()
			<-entered
			if !timeout {
				clock.advance(300*time.Millisecond, time.Hour)
				close(release)
				v := <-views
				if v.cache.success.Tick != 200*time.Millisecond {
					t.Fatal("publication refreshed validation anchor")
				}
				_, r := wire(t, v, clock)
				if r.Age != "300" {
					t.Fatal(r.Age)
				}
			} else {
				clock.advance(800*time.Millisecond, time.Hour)
				old := <-views
				if old.reason != health.Timeout {
					t.Fatal(old.reason)
				}
				newView := demand(t, c)
				close(release)
				until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.workers == 0 })
				if !bytes.Equal(c.cache.snapshot.Bytes(), partial.Bytes()) || c.cache.success != newView.cache.success {
					t.Fatal("late worker replaced newer partial")
				}
				if old.cache.snapshot.Valid() {
					t.Fatal("old view acquired newer cache")
				}
			}
		})
	}
}

func TestClockRegressionAndShutdownRace(t *testing.T) {
	clock := newFake()
	c := New(fetchFunc(func(context.Context) health.Result { return health.Result{Reason: health.UpstreamError} }), clock, nil)
	demand(t, c)
	clock.advance(time.Second, 0)
	_ = c.Now()
	clock.advance(-time.Millisecond, 0)
	if c.Now().Valid {
		t.Fatal("clock regression accepted")
	}
	if _, e := c.Demand(context.Background()); !errors.Is(e, ErrClock) {
		t.Fatal(e)
	}
	c.Stop()
	for i := 0; i < 50; i++ {
		clock := newFake()
		release := make(chan struct{})
		entered := make(chan struct{})
		c := New(fetchFunc(func(context.Context) health.Result {
			close(entered)
			<-release
			return health.Result{Reason: health.UpstreamError}
		}), clock, nil)
		done := make(chan struct{})
		go func() { _, _ = c.Demand(context.Background()); close(done) }()
		<-entered
		var wg sync.WaitGroup
		wg.Add(3)
		go func() { defer wg.Done(); c.Stop() }()
		go func() { defer wg.Done(); clock.advance(time.Second, 0) }()
		go func() { defer wg.Done(); close(release) }()
		wg.Wait()
		<-done
		until(t, func() bool { c.mu.Lock(); defer c.mu.Unlock(); return c.workers == 0 })
	}
}

func TestTenCallersShareOneRealHealthGET(t *testing.T) {
	full := snapshot(t, "current")
	var calls atomic.Int32
	gates := make(chan chan struct{}, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var gate chan struct{}
		select {
		case gate = <-gates:
		case <-r.Context().Done():
			return
		}
		select {
		case <-gate:
		case <-r.Context().Done():
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(full.Bytes())
	}))
	defer upstream.Close()
	client := health.New(upstream.URL)
	defer client.Close()
	clock := newFake()
	c := New(client, clock, nil)
	defer c.Stop()
	for wave := int32(1); wave <= 3; wave++ {
		release := make(chan struct{})
		gates <- release
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				v, e := c.Demand(context.Background())
				if e != nil {
					t.Error(e)
				} else if v.reason != "" {
					t.Error(v.reason)
				}
			}()
		}
		until(t, func() bool { return waiters(c, 10) && calls.Load() == wave })
		close(release)
		wg.Wait()
		clock.advance(5*time.Second, 5*time.Second)
		if calls.Load() != wave {
			t.Fatal("retry or autonomous call", calls.Load())
		}
	}
}
