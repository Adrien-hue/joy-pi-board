package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpDoesNotStartAListener(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"--help"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "0.0.0.0:8081") {
		t.Fatal("help must describe the approved default listener")
	}
}

func TestInvalidArgumentsFailBeforeServing(t *testing.T) {
	for _, args := range [][]string{{"--health-url", "http://127.0.0.1:8080/v1/snapshot"}, {"unexpected"}, {"--listen", "not-an-address"}} {
		var output bytes.Buffer
		if err := run(args, &output); err == nil {
			t.Fatalf("run(%v) unexpectedly succeeded", args)
		}
	}
}
