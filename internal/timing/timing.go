// Package timing contains opt-in bounded request timing, never system collection.
package timing

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const MaximumRecords = 4096

type Record struct {
	Event     string `json:"event"`
	Request   uint64 `json:"request"`
	Flight    uint64 `json:"flight"`
	HandlerNS int64  `json:"handler_ns"`
	NetworkNS int64  `json:"network_ns"`
	Shared    bool   `json:"shared"`
	Cancelled bool   `json:"cancelled"`
	Complete  bool   `json:"complete"`
	Status    int    `json:"status"`
	Eligible  bool   `json:"eligible"`
}
type Summary struct {
	Event     string `json:"event"`
	Attempts  uint64 `json:"attempts"`
	Queued    uint64 `json:"queued"`
	Dropped   uint64 `json:"dropped"`
	Truncated uint64 `json:"truncated"`
}
type Recorder struct {
	next      atomic.Uint64
	queued    atomic.Uint64
	dropped   atomic.Uint64
	truncated atomic.Uint64
	emit      func(Record) bool
}

func New(emit func(Record) bool) *Recorder { return &Recorder{emit: emit} }
func (r *Recorder) Summary() Summary {
	return Summary{"timing_summary", r.next.Load(), r.queued.Load(), r.dropped.Load(), r.truncated.Load()}
}

type Flight struct {
	mu       sync.Mutex
	sequence uint64
	shared   bool
	ns       int64
	complete bool
}

func (f *Flight) Share() {
	if f != nil {
		f.mu.Lock()
		f.shared = true
		f.mu.Unlock()
	}
}

type Request struct{ flight *Flight }
type requestKey struct{}
type flightKey struct{}

func Begin(ctx context.Context) (context.Context, *Request) {
	r := &Request{}
	return context.WithValue(ctx, requestKey{}, r), r
}
func Associate(ctx context.Context, f *Flight) {
	if r, ok := ctx.Value(requestKey{}).(*Request); ok {
		r.flight = f
	}
}
func NewFlight(ctx context.Context, seq uint64) *Flight {
	if ctx.Value(requestKey{}) == nil {
		return nil
	}
	return &Flight{sequence: seq}
}
func WithFlight(ctx context.Context, f *Flight) context.Context {
	if f == nil {
		return ctx
	}
	return context.WithValue(ctx, flightKey{}, f)
}

// Complete denotes full body receipt, not successful schema validation.
func Network(ctx context.Context) func(bool) {
	f, _ := ctx.Value(flightKey{}).(*Flight)
	if f == nil {
		return func(bool) {}
	}
	start := time.Now()
	var once sync.Once
	return func(complete bool) {
		once.Do(func() { f.mu.Lock(); f.ns = time.Since(start).Nanoseconds(); f.complete = complete; f.mu.Unlock() })
	}
}
func (r *Recorder) Finish(q *Request, start time.Time, status int, cancelled bool, success bool) {
	n := r.next.Add(1)
	if n > MaximumRecords {
		r.truncated.Add(1)
		return
	}
	v := Record{Event: "request_timing", Request: n, HandlerNS: time.Since(start).Nanoseconds(), Status: status, Cancelled: cancelled}
	if f := q.flight; f != nil {
		f.mu.Lock()
		v.Flight = f.sequence
		v.Shared = f.shared
		v.NetworkNS = f.ns
		v.Complete = f.complete
		f.mu.Unlock()
	}
	v.Eligible = success && status == 200 && v.Flight != 0 && v.Complete && !v.Shared && !cancelled && v.NetworkNS >= 0 && v.HandlerNS >= v.NetworkNS
	if r.emit(v) {
		r.queued.Add(1)
	} else {
		r.dropped.Add(1)
	}
}
