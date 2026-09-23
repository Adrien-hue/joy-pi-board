// Package config parses startup settings without network or filesystem access.
package config

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const DefaultListen = "0.0.0.0:8081"
const DefaultHealth = "http://127.0.0.1:8080/v1/snapshot"
const Help = "Joy Pi Board\n  --listen IP:port (default 0.0.0.0:8081; JOY_PI_BOARD_LISTEN)\n  --health-url URL (default http://127.0.0.1:8080/v1/snapshot; JOY_PI_BOARD_HEALTH_URL)\n  --help, -h\n  --version\nExplicit flag > present environment > default.\n"

type Config struct{ Listen, HealthURL, Action string }

func Parse(args []string, env func(string) (string, bool)) (Config, error) {
	c := Config{Listen: DefaultListen, HealthURL: DefaultHealth}
	values := make(map[string]string)
	for i := 0; i < len(args); i++ {
		name, value, equal := strings.Cut(args[i], "=")
		if name == "-h" {
			name = "--help"
		}
		if _, duplicate := values[name]; duplicate {
			return c, errors.New("duplicate option")
		}
		switch name {
		case "--help", "--version":
			if equal {
				return c, errors.New("informational option takes no value")
			}
		case "--listen", "--health-url":
			if !equal {
				i++
				if i >= len(args) || strings.HasPrefix(args[i], "--") {
					return c, errors.New("option requires a value")
				}
				value = args[i]
			}
		default:
			return c, errors.New("unsupported argument")
		}
		values[name] = value
	}
	for _, action := range []string{"--help", "--version"} {
		if _, ok := values[action]; ok {
			if len(values) != 1 {
				return c, errors.New("informational action must stand alone")
			}
			c.Action = action
			return c, nil
		}
	}
	for _, setting := range []struct {
		flag, variable string
		target         *string
	}{
		{"--listen", "JOY_PI_BOARD_LISTEN", &c.Listen}, {"--health-url", "JOY_PI_BOARD_HEALTH_URL", &c.HealthURL},
	} {
		if v, ok := values[setting.flag]; ok {
			*setting.target = v
		} else if v, ok := env(setting.variable); ok {
			*setting.target = v
		}
	}
	if !validListen(c.Listen) {
		return c, errors.New("listen: numeric IP and decimal port 0-65535 required")
	}
	if !validHealth(c.HealthURL) {
		return c, errors.New("health-url: explicit HTTP host/port and /v1/snapshot required; no credentials, query or fragment")
	}
	return c, nil
}

func port(s string, zero bool) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	n, err := strconv.ParseUint(s, 10, 16)
	return err == nil && (zero || n != 0)
}

func validListen(s string) bool {
	host, p, err := net.SplitHostPort(s)
	return err == nil && net.ParseIP(host) != nil && port(p, true)
}

func validHealth(s string) bool {
	if strings.ContainsAny(s, "?#%\\ \t\r\n") {
		return false
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "http" || u.Opaque != "" || u.User != nil || u.Path != "/v1/snapshot" || u.RawPath != "" {
		return false
	}
	host, p, err := net.SplitHostPort(u.Host)
	if err != nil || !port(p, false) || host == "" {
		return false
	}
	if net.ParseIP(host) != nil {
		return true
	}
	if strings.HasPrefix(u.Host, "[") {
		return false
	}
	if len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(strings.TrimSuffix(host, "."), ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, ch := range label {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-') {
				return false
			}
		}
	}
	return true
}
