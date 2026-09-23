package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/config"
	"github.com/Adrien-hue/joy-pi-board/internal/eventlog"
	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/httpapi"
	"github.com/Adrien-hue/joy-pi-board/internal/overview"
	"github.com/Adrien-hue/joy-pi-board/web"
)

var version = "unknown"
var revision = "unknown"
var buildTime = "unknown"

type exitError struct {
	code    int
	message string
	silent  bool
}

func (e *exitError) Error() string { return e.message }

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		var e *exitError
		if errors.As(err, &e) {
			if !e.silent {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(e.code)
		}
		_, _ = fmt.Fprintln(os.Stderr, "Board failed")
		os.Exit(1)
	}
}
func run(args []string, output io.Writer) error {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	return execute(args, output, os.LookupEnv, signals)
}
func execute(args []string, output io.Writer, env func(string) (string, bool), signals <-chan os.Signal) error {
	cfg, err := config.Parse(args, env)
	if err != nil {
		return &exitError{2, err.Error(), false}
	}
	identity := fmt.Sprintf("version=%s revision=%s built=%s", version, revision, buildTime)
	if cfg.Action == "--help" {
		_, err = io.WriteString(output, config.Help)
		return err
	}
	if cfg.Action == "--version" {
		_, err = fmt.Fprintln(output, "joy-pi-board", identity)
		return err
	}
	assets := web.Assets()
	if _, err = fs.ReadFile(assets, "index.html"); err != nil {
		return &exitError{1, "embedded interface unavailable", false}
	}
	provider := health.New(cfg.HealthURL)
	defer provider.Close()
	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return &exitError{1, "unable to bind Board listener", false}
	}
	defer listener.Close()
	logs := eventlog.New(output)
	defer logs.Close()
	logs.Lifecycle("starting", "", identity)
	coordinator := overview.New(provider, nil, logs)
	handler := httpapi.New(assets, coordinator)
	defer handler.Stop()
	server := httpapi.Server(handler)
	server.ErrorLog = log.New(logs, "", 0)
	served := make(chan error, 1)
	go func() { served <- server.Serve(httpapi.LimitConnections(listener)) }()
	logs.Lifecycle("listening", listener.Addr().String(), identity)
	select {
	case err = <-served:
		if !errors.Is(err, http.ErrServerClosed) {
			return &exitError{1, "Board server stopped unexpectedly", true}
		}
		return &exitError{1, "Board server closed unexpectedly", true}
	case <-signals:
	}
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	handler.Stop()
	logs.Lifecycle("stopping", "", "")
	drained := make(chan error, 1)
	go func() { drained <- server.Shutdown(ctx) }()
	forced := false
	select {
	case err = <-drained:
		forced = err != nil
	case <-signals:
		forced = true
	case <-ctx.Done():
		forced = true
	}
	if forced {
		cancel()
		_ = server.Close()
	}
	logs.Lifecycle("stopped", "", "")
	flushed := make(chan struct{})
	go func() { logs.Flush(ctx); close(flushed) }()
	select {
	case <-flushed:
	case <-ctx.Done():
	case <-signals:
		forced = true
		cancel()
		_ = server.Close()
	}
	if forced {
		return &exitError{1, "Board shutdown forced", true}
	}
	return nil
}
