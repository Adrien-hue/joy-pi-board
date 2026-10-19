package overview

import (
	"context"
	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/timing"
	"sync"
	"testing"
	"time"
)

// S04-T05: all waiters, including the initial caller, are marked shared.
func TestSharedDiagnosticFlight(t *testing.T) {
	release := make(chan struct{})
	entered := make(chan struct{})
	c := New(fetchFunc(func(ctx context.Context) health.Result {
		end := timing.Network(ctx)
		close(entered)
		<-release
		end(true)
		return health.Result{Reason: health.InvalidResponse}
	}), nil, nil)
	defer c.Stop()
	ctx1, q1 := timing.Begin(context.Background())
	ctx2, q2 := timing.Begin(context.Background())
	start := time.Now()
	var wg sync.WaitGroup
	wg.Go(func() { _, _ = c.Demand(ctx1) })
	<-entered
	wg.Go(func() { _, _ = c.Demand(ctx2) })
	deadline := time.Now().Add(time.Second)
	for !waiters(c, 2) {
		if time.Now().After(deadline) {
			t.Fatal("join watchdog")
		}
		time.Sleep(time.Millisecond)
	}
	close(release)
	wg.Wait()
	var rows []timing.Record
	r := timing.New(func(v timing.Record) bool { rows = append(rows, v); return true })
	r.Finish(q1, start, 200, false, false)
	r.Finish(q2, start, 200, false, false)
	if len(rows) != 2 || rows[0].Flight != rows[1].Flight || !rows[0].Shared || !rows[1].Shared || rows[0].Eligible || rows[1].Eligible {
		t.Fatal(rows)
	}
}
