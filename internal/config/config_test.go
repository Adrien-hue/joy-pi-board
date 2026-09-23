package config

import (
	"strings"
	"testing"
)

// S02-T01: presence, precedence, syntax and no I/O by construction.
func TestConfiguration(t *testing.T) {
	empty := func(string) (string, bool) { return "", false }
	c, err := Parse(nil, empty)
	if err != nil || c.Listen != DefaultListen || c.HealthURL != DefaultHealth {
		t.Fatal(c, err)
	}
	for _, listen := range []string{"127.0.0.1:0", "0.0.0.0:65535", "[::1]:8081"} {
		if _, err := Parse([]string{"--listen", listen}, empty); err != nil {
			t.Fatal(err)
		}
	}
	for _, url := range []string{DefaultHealth, "http://health.local:8080/v1/snapshot", "http://[::1]:1/v1/snapshot"} {
		if _, err := Parse([]string{"--health-url", url}, empty); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"--help", "--version"}, {"--help", "--listen=127.0.0.1:0"}, {"--help=1"}, {"-h", "--help"}, {"--version", "--version"}, {"--listen"}, {"--listen=x", "--listen=y"}, {"private-secret"}, {"--unknown=private-secret"}} {
		if _, err := Parse(args, empty); err == nil || strings.Contains(err.Error(), "private-secret") {
			t.Fatal(args, err)
		}
	}
	for _, s := range []string{"", ":8081", "localhost:8081", "[fe80::1%eth0]:8081", "127.0.0.1:http", "127.0.0.1:-1", "127.0.0.1:65536", " 127.0.0.1:80"} {
		if _, err := Parse([]string{"--listen", s}, empty); err == nil {
			t.Fatal(s)
		}
	}
	for _, s := range []string{"", "https://127.0.0.1:8080/v1/snapshot", "http://127.0.0.1/v1/snapshot", "http://a:0/v1/snapshot", "http://a:65536/v1/snapshot", "http://user:private-secret@a:80/v1/snapshot", DefaultHealth + "?", DefaultHealth + "#", DefaultHealth + "/", "http://a:80/v1/%73napshot", "http://[fe80::1%25eth0]:80/v1/snapshot", "http://héalth:80/v1/snapshot", "http://-a:80/v1/snapshot"} {
		if _, err := Parse([]string{"--health-url", s}, empty); err == nil || strings.Contains(err.Error(), "private-secret") {
			t.Fatal(s, err)
		}
	}
	env := func(k string) (string, bool) { return "", true }
	if _, err := Parse(nil, env); err == nil {
		t.Fatal("present empty env")
	}
	if _, err := Parse([]string{"--listen", DefaultListen, "--health-url", DefaultHealth}, env); err != nil {
		t.Fatal(err)
	}
	for _, action := range []string{"--help", "-h", "--version"} {
		if _, err := Parse([]string{action}, env); err != nil {
			t.Fatal(err)
		}
	}
	env = func(k string) (string, bool) {
		if k == "JOY_PI_BOARD_LISTEN" {
			return "127.0.0.1:0", true
		}
		return DefaultHealth, true
	}
	c, err = Parse(nil, env)
	if err != nil || c.Listen != "127.0.0.1:0" {
		t.Fatal(c, err)
	}
}
