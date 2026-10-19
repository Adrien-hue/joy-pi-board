package config

import "testing"

func TestDiagnosticFlag(t *testing.T) {
	env := func(string) (string, bool) { return "", false }
	c, e := Parse(nil, env)
	if e != nil || c.DiagnosticTiming {
		t.Fatal(c, e)
	}
	c, e = Parse([]string{"--diagnostic-timing", "--listen", DefaultListen}, env)
	if e != nil || !c.DiagnosticTiming {
		t.Fatal(c, e)
	}
	for _, args := range [][]string{{"--diagnostic-timing=true"}, {"--diagnostic-timing=false"}, {"--diagnostic-timing", "false"}, {"--diagnostic-timing", "--diagnostic-timing"}, {"--diagnostic-timing", "--help"}, {"--version", "--diagnostic-timing"}} {
		if _, e := Parse(args, env); e == nil {
			t.Fatal(args)
		}
	}
	c, e = Parse(nil, func(k string) (string, bool) { return "true", k == "JOY_PI_BOARD_DIAGNOSTIC_TIMING" })
	if e != nil || c.DiagnosticTiming {
		t.Fatal("unexpected environment setting", c, e)
	}
}
