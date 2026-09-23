package httpapi

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

type oneListener struct {
	conn      net.Conn
	once      sync.Once
	stopped   chan struct{}
	closeOnce sync.Once
}

func (l *oneListener) Accept() (net.Conn, error) {
	var c net.Conn
	l.once.Do(func() { c = l.conn })
	if c != nil {
		return c, nil
	}
	<-l.stopped
	return nil, net.ErrClosed
}
func (l *oneListener) Close() error   { l.closeOnce.Do(func() { close(l.stopped) }); return nil }
func (l *oneListener) Addr() net.Addr { return l.conn.LocalAddr() }

type deadlineConn struct {
	net.Conn
	reads, writes chan time.Time
}

func (c *deadlineConn) SetReadDeadline(d time.Time) error {
	c.reads <- d
	return c.Conn.SetReadDeadline(d)
}
func (c *deadlineConn) SetWriteDeadline(d time.Time) error {
	c.writes <- d
	return c.Conn.SetWriteDeadline(d)
}

// S02-T10: observe actual deadlines on a connection; no 15-second idle sleep.
func TestServerDeadlineWiring(t *testing.T) {
	serverSide, client := net.Pipe()
	defer client.Close()
	recorded := &deadlineConn{serverSide, make(chan time.Time, 32), make(chan time.Time, 32)}
	l := &oneListener{conn: recorded, stopped: make(chan struct{})}
	s := Server(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "okay") }))
	defer s.Close()
	go func() { _ = s.Serve(l) }()
	start := time.Now()
	_, _ = io.WriteString(client, "GET / HTTP/1.1\r\nHost: board\r\n\r\n")
	r, e := http.ReadResponse(bufio.NewReader(client), nil)
	if e != nil {
		t.Fatal(e)
	}
	_, _ = io.Copy(io.Discard, r.Body)
	_ = r.Body.Close()
	ctx := time.NewTimer(time.Second)
	defer ctx.Stop()
	foundIdle, foundRead := false, false
	for !foundIdle {
		select {
		case d := <-recorded.reads:
			delta := d.Sub(start)
			if delta > time.Second && delta < 3*time.Second {
				foundRead = true
			}
			if delta > 14*time.Second && delta < 16*time.Second {
				foundIdle = true
			}
		case <-ctx.C:
			t.Fatal("idle deadline missing")
		}
	}
	if !foundRead {
		t.Fatal("read/header deadline missing")
	}
	d := <-recorded.writes
	if delta := d.Sub(start); delta < 2*time.Second || delta > 4*time.Second {
		t.Fatal("write deadline", delta)
	}
}

func TestRealHeaderAndWriteTimeouts(t *testing.T) {
	for _, mode := range []string{"headers", "write"} {
		t.Run(mode, func(t *testing.T) {
			serverSide, client := net.Pipe()
			defer client.Close()
			l := &oneListener{conn: serverSide, stopped: make(chan struct{})}
			closed := make(chan struct{})
			s := Server(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write(make([]byte, 65536)) }))
			s.ConnState = func(_ net.Conn, state http.ConnState) {
				if state == http.StateClosed {
					close(closed)
				}
			}
			defer s.Close()
			go func() { _ = s.Serve(l) }()
			start := time.Now()
			if mode == "write" {
				_, _ = io.WriteString(client, "GET / HTTP/1.1\r\nHost: board\r\n\r\n")
			}
			select {
			case <-closed:
			case <-time.After(4 * time.Second):
				t.Fatal("network deadline did not release connection")
			}
			elapsed := time.Since(start)
			minimum := 1900 * time.Millisecond
			if mode == "write" {
				minimum = 2900 * time.Millisecond
			}
			if elapsed < minimum {
				t.Fatal("deadline fired early", elapsed)
			}
			t.Logf("real %s timeout=%s; net.Pipe workstation, not Pi", mode, elapsed)
		})
	}
}
