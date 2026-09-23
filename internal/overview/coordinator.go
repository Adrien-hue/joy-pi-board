// Package overview coordinates demand-only Health attempts and one volatile,
// immutable snapshot. No HTTP response I/O is performed while holding its lock.
package overview

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/healthschema"
)

const AttemptBudget = time.Second
const Freshness = 30 * time.Second
const MaximumWorkers = 32

var ErrUnavailable = errors.New("Board is temporarily unavailable")
var ErrClock = errors.New("invalid elapsed clock")

// Reading keeps wall formatting separate from elapsed arithmetic. Tick is local
// to a Clock, never reconstructed from an RFC3339 timestamp. Valid fails closed.
type Reading struct {
	Wall  time.Time
	Tick  time.Duration
	Valid bool
}
type Timer interface {
	C() <-chan time.Time
	Stop() bool
}
type Clock interface {
	Now() Reading
	Until(time.Duration) Timer
}

// Serialize observations, detecting regressions across readings as well as
// relative to a particular snapshot anchor. Invalid clocks remain fail-closed.
type checkedClock struct {
	mu      sync.Mutex
	source  Clock
	last    time.Duration
	invalid bool
}

func (c *checkedClock) Now() Reading {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := c.source.Now()
	if !n.Valid || n.Tick < 0 || n.Tick < c.last {
		c.invalid = true
	}
	c.last = n.Tick
	n.Valid = n.Valid && !c.invalid
	return n
}
func (c *checkedClock) Until(d time.Duration) Timer { return c.source.Until(d) }

type realClock struct{ origin time.Time }
type realTimer struct{ *time.Timer }

func (t realTimer) C() <-chan time.Time { return t.Timer.C }
func NewClock() Clock                   { return &realClock{origin: time.Now()} }
func (c *realClock) Now() Reading {
	n := time.Now()
	return Reading{n.UTC(), n.Sub(c.origin), true}
}
func (c *realClock) Until(t time.Duration) Timer {
	return realTimer{time.NewTimer(t - time.Since(c.origin))}
}

type Fetcher interface {
	Fetch(context.Context) health.Result
}

// Event contains only a sequence, fixed outcome and approved issue paths/codes.
// Emit must return immediately; the production sink is bounded and nonblocking.
type Event struct {
	Sequence        uint64
	Outcome, Issues string
}
type Events interface{ Emit(Event) }
type cached struct {
	snapshot healthschema.Validated
	success  Reading
}
type View struct {
	sequence uint64
	reason   health.Reason
	cache    cached
}
type flight struct {
	sequence uint64
	start    Reading
	deadline time.Duration
	done     chan struct{}
	cancel   context.CancelFunc
	timer    Timer
	view     View
	err      error
}
type Coordinator struct {
	mu               sync.Mutex
	clock            Clock
	fetch            Fetcher
	events           Events
	active           *flight
	cache            cached
	sequence         uint64
	workers, waiters int
	stopping         bool
	transition       string
	beforePublish    func(uint64) // test barrier; always outside coordinator lock
}

func New(fetch Fetcher, clock Clock, events Events) *Coordinator {
	if clock == nil {
		clock = NewClock()
	}
	return &Coordinator{fetch: fetch, clock: &checkedClock{source: clock}, events: events}
}

func (c *Coordinator) Demand(ctx context.Context) (View, error) {
	c.mu.Lock()
	if err := ctx.Err(); err != nil {
		c.mu.Unlock()
		return View{}, err
	}
	if c.stopping || c.waiters >= 32 {
		c.mu.Unlock()
		return View{}, ErrUnavailable
	}
	f := c.active
	if f == nil {
		if c.workers >= MaximumWorkers {
			c.mu.Unlock()
			return View{}, ErrUnavailable
		}
		start := c.clock.Now()
		if !start.Valid || start.Tick < 0 || start.Tick > time.Duration(1<<63-1)-AttemptBudget {
			c.mu.Unlock()
			return View{}, ErrClock
		}
		work, cancel := context.WithCancel(context.Background())
		c.sequence++
		f = &flight{sequence: c.sequence, start: start, deadline: start.Tick + AttemptBudget, done: make(chan struct{}), cancel: cancel}
		f.timer = c.clock.Until(f.deadline)
		c.active = f
		c.workers++
		go c.supervise(f)
		go c.work(work, f)
	}
	c.waiters++
	c.mu.Unlock()
	defer func() { c.mu.Lock(); c.waiters--; c.mu.Unlock() }()
	select {
	case <-ctx.Done():
		return View{}, ctx.Err()
	case <-f.done:
		if err := ctx.Err(); err != nil {
			return View{}, err
		}
		return f.view, f.err
	}
}

func (c *Coordinator) supervise(f *flight) {
	select {
	case <-f.done:
	case <-f.timer.C():
		c.finish(f, health.Result{Reason: health.Timeout}, Reading{}, "")
	}
}

func (c *Coordinator) work(ctx context.Context, f *flight) {
	defer func() { c.mu.Lock(); c.workers--; c.mu.Unlock() }()
	result := c.fetch.Fetch(ctx)
	anchor := c.clock.Now() // validation-completion pair, before publication contention
	issues := ""
	if result.Snapshot.Valid() {
		projection, err := result.Snapshot.Snapshot()
		if err != nil {
			result = health.Result{Reason: health.InvalidResponse}
		} else {
			var b strings.Builder
			for _, issue := range projection.Issues {
				b.WriteString(issue.Path)
				b.WriteByte(':')
				b.WriteString(issue.Code)
				b.WriteByte(';')
			}
			issues = b.String()
		}
	}
	if c.beforePublish != nil {
		c.beforePublish(f.sequence)
	}
	c.finish(f, result, anchor, issues)
}

func (c *Coordinator) finish(f *flight, result health.Result, anchor Reading, issues string) {
	c.mu.Lock()
	if c.stopping || c.active != f {
		c.mu.Unlock()
		return
	}
	now := c.clock.Now()
	var err error
	if !now.Valid || now.Tick < f.start.Tick {
		err = ErrClock
	} else if now.Tick >= f.deadline {
		result = health.Result{Reason: health.Timeout}
		issues = ""
	} else if result.Cancelled {
		err = ErrUnavailable
	} else if result.Snapshot.Valid() {
		if !anchor.Valid || anchor.Tick < f.start.Tick || anchor.Tick > now.Tick {
			err = ErrClock
		} else {
			c.cache = cached{result.Snapshot, anchor}
		}
	} else if result.Reason == "" {
		err = ErrClock
	}
	f.view = View{f.sequence, result.Reason, c.cache}
	f.err = err
	c.active = nil
	var event *Event
	if err == nil {
		outcome := string(result.Reason)
		if outcome == "" {
			outcome = "available"
		}
		key := outcome + "|" + issues
		if c.transition != key {
			c.transition = key
			event = &Event{f.sequence, outcome, issues}
		}
	}
	close(f.done)
	c.mu.Unlock()
	f.timer.Stop()
	f.cancel()
	if event != nil && c.events != nil {
		c.events.Emit(*event)
	}
}

// Stop is idempotent and never waits for network, writers or abandoned workers.
func (c *Coordinator) Stop() {
	c.mu.Lock()
	c.stopping = true
	f := c.active
	c.active = nil
	if f != nil {
		f.err = ErrUnavailable
		close(f.done)
	}
	c.mu.Unlock()
	if f != nil {
		f.timer.Stop()
		f.cancel()
	}
}

func (c *Coordinator) Now() Reading { return c.clock.Now() }
