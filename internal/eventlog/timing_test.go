package eventlog

import (
	"context"
	"github.com/Adrien-hue/joy-pi-board/internal/timing"
	"testing"
	"time"
)

func TestDiagnosticQueueNeverBlocks(t *testing.T) {
	w := &gatedWriter{entered: make(chan struct{}), release: make(chan struct{})}
	l := New(w)
	defer l.Close()
	l.Lifecycle("starting", "", "")
	<-w.entered
	r := timing.New(l.Timing)
	_, q := timing.Begin(context.Background())
	start := time.Now()
	for i := 0; i < timing.MaximumRecords+1; i++ {
		r.Finish(q, start, 503, false, false)
	}
	s := r.Summary()
	if s.Queued != 64 || s.Dropped != timing.MaximumRecords-64 || s.Truncated != 1 {
		t.Fatal(s)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	l.Flush(ctx)
	close(w.release)
}
