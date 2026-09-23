package eventlog

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/overview"
)

type gatedWriter struct {
	mu               sync.Mutex
	body             bytes.Buffer
	entered, release chan struct{}
	once             sync.Once
}

func (w *gatedWriter) Write(p []byte) (int, error) {
	w.once.Do(func() { close(w.entered); <-w.release })
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.body.Write(p)
}
func (w *gatedWriter) text() string { w.mu.Lock(); defer w.mu.Unlock(); return w.body.String() }

// S02-T11: blocked OS writer, bounded drops, sanitization and sequence ordering.
func TestBlockedSinkAndSanitizedErrors(t *testing.T) {
	w := &gatedWriter{entered: make(chan struct{}), release: make(chan struct{})}
	l := New(w)
	defer l.Close()
	l.Lifecycle("starting", "", "")
	<-w.entered
	for i := 0; i < 1000; i++ {
		_, _ = l.Write([]byte("private-secret upstream http://credential@host body hostname"))
	}
	if len(l.queue) != 64 || l.dropped.Load() != 936 {
		t.Fatal(len(l.queue), l.dropped.Load())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	l.Flush(ctx)
	close(w.release)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	l.Flush(ctx)
	l.Emit(overview.Event{Sequence: 2, Outcome: "available"})
	l.Emit(overview.Event{Sequence: 1, Outcome: "timeout"})
	l.Flush(ctx)
	text := w.text()
	if strings.Contains(text, "private-secret") || strings.Contains(text, "timeout") || !strings.Contains(text, `"count":936`) || !strings.Contains(text, "http_transport_error") {
		t.Fatal(text)
	}
	l.Close()
	select {
	case <-l.done:
	case <-ctx.Done():
		t.Fatal("log worker did not exit")
	}
}
