package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

type logCapture struct {
	mu      sync.Mutex
	body    bytes.Buffer
	address chan string
}

func (w *logCapture) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	var e struct{ Event, Address string }
	_ = json.Unmarshal(b, &e)
	if e.Event == "listening" {
		w.address <- e.Address
	}
	return w.body.Write(b)
}
func (w *logCapture) text() string { w.mu.Lock(); defer w.mu.Unlock(); return w.body.String() }

// S02-T01/T11: no startup Health work, shutdown releases active HTTP and logs.
func TestStartupAndShutdownWithHealthAbsent(t *testing.T) {
	var calls atomic.Int32
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		close(entered)
		<-r.Context().Done()
		close(cancelled)
	}))
	defer upstream.Close()
	w := &logCapture{address: make(chan string, 1)}
	signals := make(chan os.Signal, 2)
	done := make(chan error, 1)
	go func() {
		done <- execute([]string{"--listen=127.0.0.1:0", "--health-url=" + upstream.URL + "/v1/snapshot"}, w, func(string) (string, bool) { return "", false }, signals)
	}()
	var address string
	select {
	case address = <-w.address:
	case <-time.After(2 * time.Second):
		t.Fatal("startup")
	}
	client := &http.Client{Timeout: 2 * time.Second}
	defer client.CloseIdleConnections()
	res, e := client.Get("http://" + address + "/")
	if e != nil {
		t.Fatal(e)
	}
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()
	if res.StatusCode != 200 || calls.Load() != 0 {
		t.Fatal(res.StatusCode, calls.Load())
	}
	response := make(chan error, 1)
	go func() {
		res, e := client.Get("http://" + address + "/api/v1/overview")
		if e == nil {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()
			if res.StatusCode != 503 {
				e = errors.New("expected Board shutdown 503")
			}
		}
		response <- e
	}()
	<-entered
	start := time.Now()
	signals <- syscall.SIGTERM
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("shutdown budget")
	}
	if e := <-response; e != nil {
		t.Fatal(e)
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("Health not cancelled")
	}
	logs := w.text()
	for _, s := range []string{"starting", "listening", "stopping", "stopped"} {
		if !strings.Contains(logs, s) {
			t.Fatal(logs)
		}
	}
	if strings.Contains(logs, upstream.URL) || strings.Contains(logs, "connection_failed") || calls.Load() != 1 {
		t.Fatal(logs, calls.Load())
	}
	t.Logf("in-process signal shutdown with HTTP in flight=%s (workstation)", time.Since(start))
}

func TestSafeStartupExitCodes(t *testing.T) {
	var out bytes.Buffer
	for _, args := range [][]string{{"--bogus=private-secret"}, {"--listen=private-secret"}, {"--health-url=http://user:private-secret@a:80/v1/snapshot"}} {
		e := execute(args, &out, os.LookupEnv, nil)
		var exit *exitError
		if !errors.As(e, &exit) || exit.code != 2 || strings.Contains(e.Error(), "private-secret") {
			t.Fatal(e)
		}
	}
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	defer l.Close()
	e = execute([]string{"--listen", l.Addr().String()}, &out, os.LookupEnv, nil)
	var exit *exitError
	if !errors.As(e, &exit) || exit.code != 1 || strings.Contains(e.Error(), l.Addr().String()) {
		t.Fatal(e)
	}
	if e = execute([]string{"--version"}, &out, func(string) (string, bool) { return "", true }, nil); e != nil || !strings.Contains(out.String(), "version=unknown") {
		t.Fatal(e, out.String())
	}
}

func TestSignalProcessHelper(t *testing.T) {
	if os.Getenv("S02_SIGNAL_CHILD") != "1" {
		return
	}
	err := run([]string{"--listen=127.0.0.1:0", "--health-url=" + os.Getenv("S02_TEST_HEALTH")}, os.Stdout)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestRealSIGTERMProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("real SIGTERM runs in native Linux race/CI environment")
	}
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done(); close(cancelled) }))
	defer upstream.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestSignalProcessHelper$")
	cmd.Env = append(os.Environ(), "S02_SIGNAL_CHILD=1", "S02_TEST_HEALTH="+upstream.URL+"/v1/snapshot")
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	cmd.Stderr = io.Discard
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer cmd.Process.Kill()
	dec := json.NewDecoder(stdout)
	var address string
	for address == "" {
		var event struct{ Event, Address string }
		if e = dec.Decode(&event); e != nil {
			t.Fatal(e)
		}
		if event.Event == "listening" {
			address = event.Address
		}
	}
	go func() { _, _ = io.Copy(io.Discard, stdout) }()
	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		res, e := http.Get("http://" + address + "/api/v1/overview")
		if e == nil {
			_ = res.Body.Close()
		}
	}()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("request not started")
	}
	start := time.Now()
	if e = cmd.Process.Signal(syscall.SIGTERM); e != nil {
		t.Fatal(e)
	}
	if e = cmd.Wait(); e != nil {
		t.Fatal(e)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("signal budget")
	}
	<-requestDone
	select {
	case <-cancelled:
	case <-ctx.Done():
		t.Fatal("upstream not cancelled")
	}
	t.Logf("real Linux SIGTERM exit=%s; no Pi acceptance", time.Since(start))
}

type blockedLog struct {
	entered, release chan struct{}
	once             sync.Once
}

func (w *blockedLog) Write(b []byte) (int, error) {
	w.once.Do(func() { close(w.entered) })
	<-w.release
	return len(b), nil
}

func TestBlockedLogShutdownBudgetAndSecondSignal(t *testing.T) {
	for _, second := range []bool{false, true} {
		t.Run(strings.Join([]string{"second", fmt.Sprint(second)}, "-"), func(t *testing.T) {
			w := &blockedLog{entered: make(chan struct{}), release: make(chan struct{})}
			defer close(w.release)
			signals := make(chan os.Signal, 2)
			done := make(chan error, 1)
			go func() {
				done <- execute([]string{"--listen=127.0.0.1:0"}, w, func(string) (string, bool) { return "", false }, signals)
			}()
			<-w.entered
			start := time.Now()
			signals <- syscall.SIGTERM
			if second {
				signals <- os.Interrupt
			}
			select {
			case err := <-done:
				if second {
					var e *exitError
					if !errors.As(err, &e) || e.code != 1 || !e.silent {
						t.Fatal(err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
			case <-time.After(6 * time.Second):
				t.Fatal("blocked writer prevented bounded exit")
			}
			if second && time.Since(start) > time.Second {
				t.Fatal("second signal not immediate")
			}
			t.Logf("blocked logger shutdown second=%v duration=%s, one absolute 5s context", second, time.Since(start))
		})
	}
}
