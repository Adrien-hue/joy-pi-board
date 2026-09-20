package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/httpui"
	"github.com/Adrien-hue/joy-pi-board/web"
)

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	flags := flag.NewFlagSet("joy-pi-board", flag.ContinueOnError)
	flags.SetOutput(output)
	listen := flags.String("listen", "0.0.0.0:8081", "S00 HTTP listen address")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("positional arguments are not supported")
	}
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()
	logger := log.New(output, "joy-pi-board: ", 0)
	logger.Printf("S00 shell listening on %s", listener.Addr())
	server := &http.Server{
		Handler:           httpui.NewHandler(web.Assets()),
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       15 * time.Second,
		ErrorLog:          logger,
	}
	return server.Serve(listener)
}
