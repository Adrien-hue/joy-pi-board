// Package eventlog provides one bounded, nonblocking queue. A blocked OS writer
// can retain at most one writer goroutine; callers and shutdown never wait on it.
package eventlog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/Adrien-hue/joy-pi-board/internal/overview"
)

type entry struct {
	Event    string `json:"event"`
	Sequence uint64 `json:"sequence,omitempty"`
	Outcome  string `json:"outcome,omitempty"`
	Issues   string `json:"issues,omitempty"`
	Address  string `json:"address,omitempty"`
	Build    string `json:"build,omitempty"`
	flush    chan struct{}
}
type Log struct {
	queue   chan entry
	stop    chan struct{}
	done    chan struct{}
	closed  atomic.Bool
	dropped atomic.Uint64
}

func New(w io.Writer) *Log {
	l := &Log{queue: make(chan entry, 64), stop: make(chan struct{}), done: make(chan struct{})}
	go func() {
		defer close(l.done)
		var sequence uint64
		for {
			select {
			case <-l.stop:
				return
			case e := <-l.queue:
				if e.flush != nil {
					close(e.flush)
					continue
				}
				if e.Sequence != 0 {
					if e.Sequence <= sequence {
						continue
					}
					sequence = e.Sequence
				}
				if n := l.dropped.Swap(0); n != 0 {
					_, _ = fmt.Fprintf(w, "{\"event\":\"logs_dropped\",\"count\":%d}\n", n)
				}
				_ = json.NewEncoder(w).Encode(e)
			}
		}
	}()
	return l
}
func (l *Log) send(e entry) {
	if l.closed.Load() {
		return
	}
	select {
	case l.queue <- e:
	default:
		l.dropped.Add(1)
	}
}
func (l *Log) Emit(e overview.Event) {
	l.send(entry{Event: "provider_transition", Sequence: e.Sequence, Outcome: e.Outcome, Issues: e.Issues})
}

// Lifecycle arguments are validated/bound local address and build constants,
// never configuration URLs or raw errors.
func (l *Log) Lifecycle(event, address, build string) {
	l.send(entry{Event: event, Address: address, Build: build})
}
func (l *Log) Close() {
	if l.closed.CompareAndSwap(false, true) {
		close(l.stop)
	}
}
func (l *Log) Flush(ctx context.Context) {
	done := make(chan struct{})
	select {
	case l.queue <- entry{flush: done}:
	case <-ctx.Done():
		return
	}
	select {
	case <-done:
	case <-ctx.Done():
	}
}

// Write sanitizes net/http's arbitrary ErrorLog before entering the queue.
func (l *Log) Write(p []byte) (int, error) {
	l.send(entry{Event: "http_transport_error"})
	return len(p), nil
}
