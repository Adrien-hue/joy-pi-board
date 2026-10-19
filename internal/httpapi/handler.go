// Package httpapi composes overview with the existing static router.
package httpapi

import (
	"errors"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/httpui"
	"github.com/Adrien-hue/joy-pi-board/internal/httpwire"
	"github.com/Adrien-hue/joy-pi-board/internal/overview"
	"github.com/Adrien-hue/joy-pi-board/internal/timing"
)

type Handler struct {
	Timing   *timing.Recorder
	service  *overview.Coordinator
	assets   http.Handler
	mu       sync.Mutex
	stopping bool
	requests int
	// Test seam models a known pause before finalization, never an HTTP hook.
	beforeFinalize func()
	afterAssemble  func()
}

func New(assets fs.FS, service *overview.Coordinator) *Handler {
	return &Handler{service: service, assets: httpui.NewHandler(assets)}
}

func (h *Handler) Stop() {
	h.mu.Lock()
	h.stopping = true
	h.mu.Unlock()
	h.service.Stop()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	success := false
	if h.Timing != nil && r.URL.Path == "/api/v1/overview" {
		ctx, q := timing.Begin(r.Context())
		r = r.WithContext(ctx)
		tracked := &statusWriter{ResponseWriter: w}
		w = tracked
		defer func() { h.Timing.Finish(q, start, tracked.status, r.Context().Err() != nil, success) }()
	}
	httpwire.Headers(w)
	h.mu.Lock()
	if h.stopping || h.requests >= 32 {
		h.mu.Unlock()
		routeError(w, r, 503)
		return
	}
	h.requests++
	h.mu.Unlock()
	defer func() { h.mu.Lock(); h.requests--; h.mu.Unlock() }()
	if r.URL.IsAbs() || r.URL.Host != "" || r.URL.Opaque != "" || !strings.HasPrefix(r.URL.Path, "/") {
		routeError(w, r, 400)
		return
	}
	if r.URL.Path != "/api/v1/overview" || r.URL.EscapedPath() != r.URL.Path {
		h.assets.ServeHTTP(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		httpwire.Error(w, r, 405)
		return
	}
	if httpwire.InvalidRequest(r) {
		w.Header().Set("Connection", "close")
		httpwire.Error(w, r, 400)
		return
	}
	view, err := h.service.Demand(r.Context())
	if r.Context().Err() != nil {
		return
	}
	if err != nil {
		h.fail(w, r, err)
		return
	}
	prepared, err := view.Prepare()
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if h.beforeFinalize != nil {
		h.beforeFinalize()
	}
	if r.Context().Err() != nil {
		return
	}
	// Snapshot compaction has already completed. Only bounded wrapper encoding
	// remains between this cut and commitment; socket/transit belongs to S03.
	var response overview.Response
	stable := false
	for attempt := 0; attempt < 3; attempt++ {
		cut := h.service.Now()
		response, err = prepared.Finalize(cut)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		if h.afterAssemble != nil {
			h.afterAssemble()
		}
		now := h.service.Now()
		if !now.Valid {
			h.fail(w, r, overview.ErrClock)
			return
		}
		if prepared.StableBetween(cut, now) {
			stable = true
			break
		}
	}
	if !stable {
		h.fail(w, r, overview.ErrUnavailable)
		return
	}
	if r.Context().Err() != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if response.Age != "" {
		w.Header().Set(overview.AgeHeader, response.Age)
	}
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(response.Body)
	success = writeErr == nil && view.Successful()
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(s int) { w.status = s; w.ResponseWriter.WriteHeader(s) }
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(b)
}
func routeError(w http.ResponseWriter, r *http.Request, status int) {
	if r.URL.Path == "/api/v1/overview" && r.URL.EscapedPath() == r.URL.Path {
		httpwire.Error(w, r, status)
	} else {
		httpui.Error(w, r, status)
	}
}

func (*Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status := 500
	if errors.Is(err, overview.ErrUnavailable) {
		status = 503
	}
	httpwire.Error(w, r, status)
}

func Server(handler http.Handler) *http.Server {
	return &http.Server{Handler: handler, DisableGeneralOptionsHandler: true, MaxHeaderBytes: 8192, ReadHeaderTimeout: 2 * time.Second, ReadTimeout: 2 * time.Second, WriteTimeout: 3 * time.Second, IdleTimeout: 15 * time.Second}
}
