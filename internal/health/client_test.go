package health

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func full(t *testing.T) []byte {
	t.Helper()
	b, e := os.ReadFile("../../docs/examples/overview-current.json")
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
	return v.Health.Snapshot
}
func fetch(t *testing.T, c *Client) Result {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.Fetch(ctx)
}

// S02-T02/T03: actual HTTP status, header, framing and schema policy.
func TestTransportPolicy(t *testing.T) {
	body := string(full(t))
	for _, tc := range []struct {
		name    string
		status  int
		headers http.Header
		body    string
		want    Reason
	}{
		{"valid", 200, http.Header{"Content-Type": {"application/json"}}, body, ""},
		{"charset", 200, http.Header{"Content-Type": {"Application/JSON; charset=UTF-8"}, "Content-Encoding": {"identity"}}, body, ""},
		{"interim", 103, http.Header{"Content-Type": {"application/json"}}, body, ""},
		{"redirect", 302, http.Header{"Location": {"http://127.0.0.1:1/secret"}}, "", UpstreamError},
		{"204", 204, nil, "", UpstreamError}, {"500", 500, nil, "private-secret", UpstreamError}, {"503", 503, nil, "", UpstreamError},
		{"101", 101, nil, "", UpstreamError},
		{"missing-mime", 200, http.Header{"Content-Type": {}}, body, InvalidResponse},
		{"duplicate-mime", 200, http.Header{"Content-Type": {"application/json", "application/json"}}, body, InvalidResponse},
		{"charset-bad", 200, http.Header{"Content-Type": {"application/json; charset=ascii"}}, body, InvalidResponse},
		{"charset-repeat", 200, http.Header{"Content-Type": {"application/json; charset=utf-8; charset=utf-8"}}, body, InvalidResponse},
		{"parameter", 200, http.Header{"Content-Type": {"application/json; x=y"}}, body, InvalidResponse},
		{"suffix", 200, http.Header{"Content-Type": {"application/problem+json"}}, body, InvalidResponse},
		{"gzip", 200, http.Header{"Content-Type": {"application/json"}, "Content-Encoding": {"gzip"}}, body, InvalidResponse},
		{"enc-duplicate", 200, http.Header{"Content-Type": {"application/json"}, "Content-Encoding": {"identity", "identity"}}, body, InvalidResponse},
		{"enc-list", 200, http.Header{"Content-Type": {"application/json"}, "Content-Encoding": {"identity, identity"}}, body, InvalidResponse},
		{"trailer", 200, http.Header{"Content-Type": {"application/json"}, "Trailer": {"X-Late"}}, body, InvalidResponse},
		{"unknown", 200, http.Header{"Content-Type": {"application/json"}}, `{"schema_version":"2"}`, UnsupportedSchemaVersion},
		{"unknown-malformed", 200, http.Header{"Content-Type": {"application/json"}}, `{"schema_version":"2"`, InvalidResponse},
		{"unknown-duplicate", 200, http.Header{"Content-Type": {"application/json"}}, `{"schema_version":"2","x":1,"x":2}`, InvalidResponse},
		{"oversize-stream", 200, http.Header{"Content-Type": {"application/json"}}, strings.Repeat(" ", 65537), InvalidResponse},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method != "GET" || r.Header.Get("Accept") != "application/json" || r.Header.Get("Accept-Encoding") != "identity" || r.Header.Get("Cookie") != "" || !r.Close {
					t.Error("request policy", r.Header)
				}
				for k, v := range tc.headers {
					w.Header()[k] = v
				}
				if tc.status == 103 {
					w.WriteHeader(103)
					w.WriteHeader(200)
				} else {
					w.WriteHeader(tc.status)
				}
				_, _ = io.WriteString(w, tc.body)
			}))
			defer s.Close()
			c := New(s.URL + "/v1/snapshot")
			defer c.Close()
			r := fetch(t, c)
			if r.Reason != tc.want || r.Cancelled || r.Snapshot.Valid() != (tc.want == "") {
				t.Fatalf("%+v want %s", r, tc.want)
			}
			if requests.Load() != 1 {
				t.Fatal(requests.Load())
			}
		})
	}
}

func TestBodySizeBoundary(t *testing.T) {
	base := full(t)
	for _, size := range []int{65535, 65536, 65537} {
		for _, declared := range []bool{false, true} {
			t.Run(fmt.Sprint(size, declared), func(t *testing.T) {
				body := append(append([]byte{}, base...), []byte(strings.Repeat(" ", size-len(base)))...)
				s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if declared {
						w.Header().Set("Content-Length", fmt.Sprint(size))
					} else {
						w.(http.Flusher).Flush()
					}
					_, _ = w.Write(body)
				}))
				defer s.Close()
				c := New(s.URL)
				defer c.Close()
				result := fetch(t, c)
				if (size <= 65536) != result.Snapshot.Valid() {
					t.Fatal(size, result.Reason)
				}
			})
		}
	}
}

// Raw server avoids net/http repairing intentionally malformed responses.
func TestWireFailuresAndNoReplay(t *testing.T) {
	for _, tc := range []struct {
		wire string
		want Reason
	}{
		{"", ConnectionFailed}, {"not HTTP\r\n\r\n", ConnectionFailed},
		{"HTTP/1.1 200 OK\r\nX-Huge: " + strings.Repeat("x", 9000) + "\r\n\r\n", ConnectionFailed},
		{"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 100\r\n\r\n{}", InvalidResponse},
		{"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nTransfer-Encoding: chunked\r\n\r\nZ\r\n", InvalidResponse},
		{"HTTP/1.1 503 Nope\r\nContent-Length: 100\r\n\r\n", UpstreamError},
	} {
		t.Run(fmt.Sprint(tc.want, len(tc.wire)), func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			var accepts, requests atomic.Int32
			done := make(chan struct{})
			go func() {
				defer close(done)
				for i := 0; i < 2; i++ {
					conn, e := listener.Accept()
					if e != nil {
						return
					}
					accepts.Add(1)
					_, e = http.ReadRequest(bufio.NewReader(conn))
					if e == nil {
						requests.Add(1)
					}
					_, _ = io.WriteString(conn, tc.wire)
					_ = conn.Close()
				}
			}()
			c := New("http://" + listener.Addr().String())
			defer c.Close()
			for i := 0; i < 2; i++ {
				if r := fetch(t, c); r.Reason != tc.want {
					t.Fatal(r.Reason)
				}
			}
			<-done
			if accepts.Load() != 2 || requests.Load() != 2 {
				t.Fatal(accepts.Load(), requests.Load())
			}
		})
	}
}

func TestSingleAddressProxyAndRedirect(t *testing.T) {
	var trap atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { trap.Add(1) }))
	defer s.Close()
	t.Setenv("HTTP_PROXY", s.URL)
	t.Setenv("http_proxy", s.URL)
	t.Setenv("NO_PROXY", "")
	var dials int
	c := newClient("http://health.invalid:8080/v1/snapshot", func(context.Context, string) ([]net.IPAddr, error) {
		return []net.IPAddr{{IP: net.ParseIP("::1")}, {IP: net.ParseIP("127.0.0.2")}, {IP: net.ParseIP("127.0.0.1")}}, nil
	}, func(ctx context.Context, network, address string) (net.Conn, error) {
		dials++
		if address != "127.0.0.1:8080" {
			t.Error(address)
		}
		return nil, errors.New("private-secret")
	})
	defer c.Close()
	if r := fetch(t, c); r.Reason != ConnectionFailed || dials != 1 || trap.Load() != 0 {
		t.Fatal(r.Reason, dials, trap.Load())
	}
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Location", s.URL); w.WriteHeader(302) }))
	defer redirect.Close()
	c = New(redirect.URL)
	defer c.Close()
	if r := fetch(t, c); r.Reason != UpstreamError || trap.Load() != 0 {
		t.Fatal(r.Reason, trap.Load())
	}
}

func TestCancellationAndRefusal(t *testing.T) {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := l.Addr().String()
	_ = l.Close()
	c := New("http://" + address)
	defer c.Close()
	start := time.Now()
	if r := fetch(t, c); r.Reason != ConnectionFailed {
		t.Fatal(r.Reason)
	}
	t.Logf("loopback refusal duration=%s (workstation, not Pi gate)", time.Since(start))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if r := c.Fetch(ctx); !r.Cancelled || r.Reason != "" {
		t.Fatal(r)
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type countedBody struct {
	read   int
	closed bool
}

func (b *countedBody) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = ' '
	}
	b.read += len(p)
	return len(p), nil
}
func (b *countedBody) Close() error { b.closed = true; return nil }

type earlyTimeout struct{}

func (earlyTimeout) Error() string   { return "private-secret" }
func (earlyTimeout) Timeout() bool   { return true }
func (earlyTimeout) Temporary() bool { return false }

func TestCountedBodyLimitAndEarlyTransportTimeout(t *testing.T) {
	for _, declared := range []int64{-1, 65537} {
		body := new(countedBody)
		c := &Client{endpoint: "http://health.invalid:8080/v1/snapshot", http: &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: body, ContentLength: declared}, nil
		})}}
		if r := fetch(t, c); r.Reason != InvalidResponse {
			t.Fatal(r.Reason)
		}
		want := 65537
		if declared > 65536 {
			want = 0
		}
		if body.read != want || !body.closed {
			t.Fatal(body.read, body.closed)
		}
	}
	c := &Client{endpoint: "http://health.invalid:8080/v1/snapshot", http: &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) { return nil, earlyTimeout{} })}}
	if r := fetch(t, c); r.Reason != Timeout || r.Cancelled {
		t.Fatal(r)
	}
}
