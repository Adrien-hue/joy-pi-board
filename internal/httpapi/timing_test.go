package httpapi

import (
	"context"
	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/overview"
	"github.com/Adrien-hue/joy-pi-board/internal/timing"
	"github.com/Adrien-hue/joy-pi-board/web"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestDiagnosticFlightCorrelation(t *testing.T) {
	snapshot := full(t)
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	c := overview.New(fetchFunc(func(ctx context.Context) health.Result {
		end := timing.Network(ctx)
		entered <- struct{}{}
		<-release
		end(true)
		return health.Result{Snapshot: snapshot}
	}), nil, nil)
	defer c.Stop()
	h := New(web.Assets(), c)
	var mu sync.Mutex
	var rows []timing.Record
	h.Timing = timing.New(func(v timing.Record) bool { mu.Lock(); defer mu.Unlock(); rows = append(rows, v); return true })
	var wg sync.WaitGroup
	wg.Go(func() { request(h, "GET", "/api/v1/overview", "") })
	<-entered
	// Race a second caller with completion; the deterministic sharing barrier
	// is covered separately by overview.TestSharedDiagnosticFlight.
	wg.Go(func() { request(h, "GET", "/api/v1/overview", "") })
	close(release)
	wg.Wait()
	if len(rows) != 2 {
		t.Fatal(rows)
	}
	for _, r := range rows {
		if r.Flight == 0 || !r.Complete || r.Cancelled {
			t.Fatal(r)
		}
		if r.Shared && r.Eligible {
			t.Fatal(r)
		}
	}
	// Distinct next demand receives a different flight, with validation retained in handler time.
	request(h, "GET", "/api/v1/overview", "")
	if rows[2].Flight <= rows[0].Flight || !rows[2].Eligible {
		t.Fatal(rows)
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "*;q=0")
	h.Stop()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 503 || w.Body.Len() != 0 || w.Header().Get("Content-Encoding") != "" {
		t.Fatal(w.Code, w.Header())
	}
}
