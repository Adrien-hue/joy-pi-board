package httpapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
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
	"github.com/Adrien-hue/joy-pi-board/internal/overview"
	"github.com/Adrien-hue/joy-pi-board/web"
)

type fetchFunc func(context.Context) health.Result

func (f fetchFunc) Fetch(ctx context.Context) health.Result { return f(ctx) }
func full(t *testing.T) healthschema.Validated {
	t.Helper()
	b, e := os.ReadFile("../../docs/examples/overview-current.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Health struct {
			Snapshot json.RawMessage `json:"snapshot"`
		} `json:"health"`
	}
	if e = json.Unmarshal(b, &doc); e != nil {
		t.Fatal(e)
	}
	v, e := healthschema.Decode(doc.Health.Snapshot)
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func request(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
	return w
}
func spin(t *testing.T, p func() bool) {
	t.Helper()
	until := time.Now().Add(3 * time.Second)
	for !p() {
		if time.Now().After(until) {
			t.Fatal("barrier watchdog")
		}
		runtime.Gosched()
	}
}

// S02-T08/T10: actual composed routes, public errors and no work on bad requests.
func TestRoutesErrorsAndHeaders(t *testing.T) {
	var calls atomic.Int32
	c := overview.New(fetchFunc(func(context.Context) health.Result {
		calls.Add(1)
		return health.Result{Reason: health.ConnectionFailed}
	}), nil, nil)
	defer c.Stop()
	h := New(web.Assets(), c)
	for _, tc := range []struct {
		method, path, body string
		status             int
		allow              string
	}{
		{"GET", "/", "", 200, ""}, {"HEAD", "/", "", 200, ""},
		{"POST", "/", "", 405, "GET, HEAD"}, {"GET", "/?", "", 400, ""}, {"GET", "/?x=1", "", 400, ""}, {"GET", "/", "x", 400, ""},
		{"GET", "/api/v1/overview/", "", 404, ""}, {"GET", "/api/v1/%6fverview", "", 404, ""}, {"GET", "/index.html", "", 404, ""},
		{"POST", "/assets/missing.js", "", 404, ""}, {"GET", "/assets/../index.html", "", 404, ""},
		{"GET", "http://example.com/api/v1/overview", "", 400, ""},
		{"OPTIONS", "*", "", 400, ""},
		{"HEAD", "/api/v1/overview", "", 405, "GET"}, {"OPTIONS", "/api/v1/overview", "", 405, "GET"},
		{"GET", "/api/v1/overview?", "", 400, ""}, {"GET", "/api/v1/overview?x=y", "", 400, ""}, {"GET", "/api/v1/overview", "private-secret", 400, ""},
	} {
		t.Run(tc.method+tc.path+tc.body, func(t *testing.T) {
			w := request(h, tc.method, tc.path, tc.body)
			if w.Code != tc.status || w.Header().Get("Allow") != tc.allow {
				t.Fatal(w.Code, w.Header())
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal(w.Header())
			}
			if tc.method == "HEAD" && w.Body.Len() != 0 {
				t.Fatal("HEAD body")
			}
			if tc.status >= 400 && w.Header().Get("Content-Type") != "application/json" {
				t.Fatal(w.Header())
			}
			if strings.Contains(w.Body.String(), "private-secret") || w.Header().Get(overview.AgeHeader) != "" {
				t.Fatal(w.Body.String())
			}
		})
	}
	r := httptest.NewRequest("GET", "/api/v1/overview", nil)
	r.TransferEncoding = []string{"chunked"}
	r.ContentLength = -1
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 || w.Header().Get("Connection") != "close" {
		t.Fatal(w.Code, w.Header())
	}
	if calls.Load() != 0 {
		t.Fatal("invalid request contacted Health")
	}
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"reason":"connection_failed"`) || !strings.Contains(w.Body.String(), `"last_success_at":null`) {
		t.Fatal(w.Code, w.Body.String())
	}
	if calls.Load() != 1 {
		t.Fatal(calls.Load())
	}
	h.Stop()
	w = request(h, "GET", "/", "")
	if w.Code != 503 {
		t.Fatal(w.Code)
	}
}

type testTimer struct{ *time.Timer }

func (t testTimer) C() <-chan time.Time { return t.Timer.C }

type fixedClock struct {
	mu  sync.Mutex
	now overview.Reading
}

func (c *fixedClock) Now() overview.Reading              { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *fixedClock) Until(time.Duration) overview.Timer { return testTimer{time.NewTimer(time.Hour)} }
func (c *fixedClock) set(tick time.Duration, valid bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now.Tick = tick
	c.now.Valid = valid
}

func TestDelayedFinalization(t *testing.T) {
	clock := &fixedClock{now: overview.Reading{Wall: time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC), Valid: true}}
	next := health.Result{Snapshot: full(t)}
	c := overview.New(fetchFunc(func(context.Context) health.Result { return next }), clock, nil)
	defer c.Stop()
	h := New(web.Assets(), c)
	w := request(h, "GET", "/api/v1/overview", "")
	if w.Code != 200 || w.Header().Get(overview.AgeHeader) != "0" {
		t.Fatal(w.Code, w.Header())
	}
	next = health.Result{Reason: health.Timeout}
	h.beforeFinalize = func() { clock.set(30*time.Second+1, true) }
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 200 || w.Header().Get(overview.AgeHeader) != "30001" || !strings.Contains(w.Body.String(), `"snapshot":null`) {
		t.Fatal(w.Code, w.Header(), w.Body.String())
	}
	next = health.Result{Snapshot: full(t)}
	h.beforeFinalize = func() { clock.set(61*time.Second, true) }
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 503 || w.Header().Get(overview.AgeHeader) != "" {
		t.Fatal(w.Code, w.Body.String())
	}
	h.beforeFinalize = func() { clock.set(0, false) }
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 500 || !strings.Contains(w.Body.String(), "internal_error") {
		t.Fatal(w.Code, w.Body.String())
	}
}

type blockedWriter struct {
	header           http.Header
	entered, release chan struct{}
	once             sync.Once
}

func (w *blockedWriter) Header() http.Header         { return w.header }
func (w *blockedWriter) WriteHeader(int)             { w.once.Do(func() { close(w.entered) }); <-w.release }
func (w *blockedWriter) Write(b []byte) (int, error) { return len(b), nil }

func TestRequestCapacityAndBlockedWriters(t *testing.T) {
	var calls atomic.Int32
	c := overview.New(fetchFunc(func(context.Context) health.Result {
		calls.Add(1)
		return health.Result{Reason: health.ConnectionFailed}
	}), nil, nil)
	defer c.Stop()
	h := New(web.Assets(), c)
	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		w := &blockedWriter{make(http.Header), make(chan struct{}), release, sync.Once{}}
		wg.Add(1)
		go func() { defer wg.Done(); h.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/overview", nil)) }()
		<-w.entered
	}
	w := request(h, "GET", "/api/v1/overview", "")
	if w.Code != 503 || calls.Load() != 32 {
		t.Fatal(w.Code, calls.Load())
	}
	close(release)
	wg.Wait()
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 200 || calls.Load() != 33 {
		t.Fatal(w.Code, calls.Load())
	}
}

// S02-T10 uses real connections for overflow and parser/Expect behavior.
func TestConnectionCapacityAndHTTPParser(t *testing.T) {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	limited := LimitConnections(l)
	var calls atomic.Int32
	c := overview.New(fetchFunc(func(context.Context) health.Result {
		calls.Add(1)
		return health.Result{Reason: health.ConnectionFailed}
	}), nil, nil)
	defer c.Stop()
	s := Server(New(web.Assets(), c))
	defer s.Close()
	go func() { _ = s.Serve(limited) }()
	var sockets []net.Conn
	defer func() {
		for _, c := range sockets {
			_ = c.Close()
		}
	}()
	for i := 0; i < 64; i++ {
		conn, e := net.Dial("tcp", l.Addr().String())
		if e != nil {
			t.Fatal(e)
		}
		sockets = append(sockets, conn)
	}
	spin(t, func() bool { return len(limited.permits) == 64 })
	over, e := net.Dial("tcp", l.Addr().String())
	if e != nil {
		t.Fatal(e)
	}
	defer over.Close()
	_ = over.SetReadDeadline(time.Now().Add(time.Second))
	if _, e = over.Read(make([]byte, 1)); e == nil {
		t.Fatal("65th connection not closed")
	}
	for _, conn := range sockets {
		_ = conn.Close()
	}
	spin(t, func() bool { return len(limited.permits) == 0 })
	for _, wire := range []string{
		"GET /api/v1/overview HTTP/1.1\r\nHost: board\r\nContent-Length: 4\r\nExpect: 100-continue\r\n\r\n",
		"GET /api/v1/overview HTTP/1.1\r\nHost: board\r\nTransfer-Encoding: chunked\r\n\r\n",
		"GET http://board/api/v1/overview HTTP/1.1\r\nHost: board\r\n\r\n",
		"OPTIONS * HTTP/1.1\r\nHost: board\r\n\r\n",
	} {
		conn, e := net.Dial("tcp", l.Addr().String())
		if e != nil {
			t.Fatal(e)
		}
		_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
		_, _ = io.WriteString(conn, wire)
		r, e := http.ReadResponse(bufio.NewReader(conn), nil)
		if e != nil {
			t.Fatal(e)
		}
		if r.StatusCode != 400 {
			t.Fatal(r.StatusCode)
		}
		_ = conn.Close()
	}
	conn, e := net.Dial("tcp", l.Addr().String())
	if e != nil {
		t.Fatal(e)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	_, _ = fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: board\r\nX-Huge: %s\r\n\r\n", strings.Repeat("x", 20000))
	r, e := http.ReadResponse(bufio.NewReader(conn), nil)
	if e == nil && r.StatusCode != 431 {
		t.Fatal(r.StatusCode)
	}
	if calls.Load() != 0 {
		t.Fatal("bad requests contacted provider")
	}
	if s.ReadHeaderTimeout != 2*time.Second || s.ReadTimeout != 2*time.Second || s.WriteTimeout != 3*time.Second || s.IdleTimeout != 15*time.Second || s.MaxHeaderBytes != 8192 {
		t.Fatal("server limits")
	}
}

func TestMaximumMarkupEnvelope(t *testing.T) {
	base := full(t).Bytes()
	var compact bytes.Buffer
	if e := json.Compact(&compact, base); e != nil {
		t.Fatal(e)
	}
	base = compact.Bytes()
	suffix := []byte(`,"extra":"` + strings.Repeat("<>&", 21845) + `"}`)
	body := append(append([]byte{}, base[:len(base)-1]...), suffix...)
	// Trim only the derived ASCII markup member to hit exactly the Health limit.
	delta := len(body) - 65536
	body = append(body[:len(body)-2-delta], body[len(body)-2:]...)
	v, e := healthschema.Decode(body)
	if e != nil {
		t.Fatal(len(body), e)
	}
	c := overview.New(fetchFunc(func(context.Context) health.Result { return health.Result{Snapshot: v} }), nil, nil)
	defer c.Stop()
	w := request(New(web.Assets(), c), "GET", "/api/v1/overview", "")
	if w.Code != 200 || w.Body.Len() > overview.MaximumEnvelopeBytes || w.Body.Len()-len(body) > 1024 || bytes.Contains(w.Body.Bytes(), []byte(`\u003c`)) {
		t.Fatal(w.Code, w.Body.Len())
	}
	var out struct {
		Health struct {
			Snapshot json.RawMessage `json:"snapshot"`
		} `json:"health"`
	}
	if e = json.Unmarshal(w.Body.Bytes(), &out); e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(body, out.Health.Snapshot) {
		t.Fatal("raw snapshot changed")
	}
}

func TestExpirationDuringAssembly(t *testing.T) {
	clock := &fixedClock{now: overview.Reading{Wall: time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC), Valid: true}}
	next := health.Result{Snapshot: full(t)}
	c := overview.New(fetchFunc(func(context.Context) health.Result { return next }), clock, nil)
	defer c.Stop()
	h := New(web.Assets(), c)
	if w := request(h, "GET", "/api/v1/overview", ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
	next = health.Result{Reason: health.ConnectionFailed}
	clock.set(30*time.Second, true)
	var once sync.Once
	h.afterAssemble = func() { once.Do(func() { clock.set(30*time.Second+1, true) }) }
	w := request(h, "GET", "/api/v1/overview", "")
	if w.Code != 200 || w.Header().Get(overview.AgeHeader) != "30001" || !strings.Contains(w.Body.String(), `"snapshot":null`) {
		t.Fatal(w.Code, w.Body.String(), w.Header())
	}
	h.afterAssemble = func() { clock.set(clock.Now().Tick+time.Millisecond, true) }
	w = request(h, "GET", "/api/v1/overview", "")
	if w.Code != 503 {
		t.Fatal("persistent pause must be bounded", w.Code)
	}
}
