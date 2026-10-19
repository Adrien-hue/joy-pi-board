package main

import (
	"encoding/json"
	"github.com/Adrien-hue/joy-pi-board/internal/timing"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestOptInDiagnosticEndToEnd(t *testing.T) {
	data, e := os.ReadFile("../../docs/examples/overview-current.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Health struct{ Snapshot json.RawMessage }
	}
	if e = json.Unmarshal(data, &doc); e != nil {
		t.Fatal(e)
	}
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") != "identity" {
			t.Error("upstream encoding")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(doc.Health.Snapshot)
	}))
	defer provider.Close()
	for _, enabled := range []bool{false, true} {
		w := &logCapture{address: make(chan string, 1)}
		signals := make(chan os.Signal, 2)
		done := make(chan error, 1)
		args := []string{"--listen=127.0.0.1:0", "--health-url=" + provider.URL + "/v1/snapshot"}
		if enabled {
			args = append(args, "--diagnostic-timing")
		}
		go func() { done <- execute(args, w, func(string) (string, bool) { return "", false }, signals) }()
		var address string
		select {
		case address = <-w.address:
		case <-time.After(2 * time.Second):
			t.Fatal("startup")
		}
		client := &http.Client{Timeout: 2 * time.Second}
		for i := 0; i < 3; i++ {
			res, err := client.Get("http://" + address + "/api/v1/overview")
			if err != nil {
				t.Fatal(err)
			}
			_, _ = io.Copy(io.Discard, res.Body)
			res.Body.Close()
			if res.StatusCode != 200 || res.Header.Get("Content-Encoding") != "" {
				t.Fatal(res.StatusCode)
			}
		}
		client.CloseIdleConnections()
		signals <- syscall.SIGTERM
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("shutdown")
		}
		logs := w.text()
		if strings.Contains(logs, provider.URL) || strings.Contains(logs, "snapshot") || strings.Contains(logs, "9007199254740993") {
			t.Fatal("diagnostic leak", logs)
		}
		var count, summaries int
		for _, line := range strings.Split(strings.TrimSpace(logs), "\n") {
			var record timing.Record
			_ = json.Unmarshal([]byte(line), &record)
			if record.Event == "request_timing" {
				count++
				if !record.Eligible || record.Request != uint64(count) || record.Flight != uint64(count) || record.NetworkNS <= 0 {
					t.Fatal(record)
				}
			}
			if record.Event == "timing_summary" {
				summaries++
				var summary timing.Summary
				_ = json.Unmarshal([]byte(line), &summary)
				if summary.Attempts != 3 || summary.Queued != 3 || summary.Dropped != 0 || summary.Truncated != 0 {
					t.Fatal(summary)
				}
			}
		}
		if enabled && (count != 3 || summaries != 1) || !enabled && (count != 0 || summaries != 0) {
			t.Fatal(count, summaries)
		}
	}
}
