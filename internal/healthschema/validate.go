package healthschema

import (
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var uintToken = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)
var utcToken = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?Z$`)

type reader struct{ err error }

func (r *reader) fail(path, rule string) {
	if r.err == nil {
		r.err = invalid(path, rule)
	}
}
func (r *reader) object(v any, path string) map[string]any {
	o, ok := v.(map[string]any)
	if !ok {
		r.fail(path, "required object")
	}
	return o
}
func (r *reader) field(o map[string]any, key, path string) any {
	v, present := o[key]
	if !present {
		r.fail(path, "required field absent")
	}
	return v
}
func (r *reader) text(v any, path string) string {
	s, ok := v.(string)
	if !ok || s == "" {
		r.fail(path, "required non-empty string")
	}
	return s
}
func (r *reader) integer(v any, path string) *uint64 {
	if v == nil {
		return nil
	}
	n, ok := v.(json.Number)
	if !ok || !uintToken.MatchString(string(n)) {
		r.fail(path, "uint64 decimal token required")
		return nil
	}
	i, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		r.fail(path, "uint64 out of range")
		return nil
	}
	return &i
}
func (r *reader) floating(v any, path string, min, max float64) *float64 {
	if v == nil {
		return nil
	}
	n, ok := v.(json.Number)
	if !ok {
		r.fail(path, "number required")
		return nil
	}
	f, err := strconv.ParseFloat(string(n), 64)
	if err != nil || math.IsInf(f, 0) || math.IsNaN(f) || f < min || f > max {
		r.fail(path, "finite number outside domain")
		return nil
	}
	return &f
}
func (r *reader) boolean(v any, path string) *bool {
	if v == nil {
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		r.fail(path, "boolean required")
		return nil
	}
	return &b
}
func (r *reader) group(root map[string]any, key string) ByteGroup {
	p := "/" + key
	o := r.object(r.field(root, key, p), p)
	return ByteGroup{r.integer(r.field(o, "total_bytes", p+"/total_bytes"), p+"/total_bytes"), r.integer(r.field(o, "available_bytes", p+"/available_bytes"), p+"/available_bytes"), r.integer(r.field(o, "used_bytes", p+"/used_bytes"), p+"/used_bytes")}
}

func known(root map[string]any) (Snapshot, error) {
	r := &reader{}
	version, ok := r.field(root, "schema_version", "/schema_version").(string)
	if r.err != nil {
		return Snapshot{}, r.err
	}
	if !ok {
		return Snapshot{}, invalid("/schema_version", "required string")
	}
	if version != "1.0" {
		return Snapshot{}, &Error{UnsupportedSchemaVersion, "/schema_version", "unsupported version"}
	}
	s := Snapshot{SchemaVersion: version}
	s.ObservedAt = r.text(r.field(root, "observed_at", "/observed_at"), "/observed_at")
	t, err := time.Parse(time.RFC3339Nano, s.ObservedAt)
	if !utcToken.MatchString(s.ObservedAt) || err != nil || t.IsZero() {
		r.fail("/observed_at", "nonzero RFC3339 UTC timestamp ending Z required")
	}
	host := r.object(r.field(root, "host", "/host"), "/host")
	s.Host.Hostname = r.text(r.field(host, "hostname", "/host/hostname"), "/host/hostname")
	cpu := r.object(r.field(root, "cpu", "/cpu"), "/cpu")
	s.CPU.UtilizationPercent = r.floating(r.field(cpu, "utilization_percent", "/cpu/utilization_percent"), "/cpu/utilization_percent", 0, 100)
	s.CPU.LogicalCPUCount = r.integer(r.field(cpu, "logical_cpu_count", "/cpu/logical_cpu_count"), "/cpu/logical_cpu_count")
	load := r.object(r.field(root, "load", "/load"), "/load")
	s.Load.OneMinute = r.floating(r.field(load, "one_minute", "/load/one_minute"), "/load/one_minute", 0, math.Inf(1))
	s.Load.FiveMinutes = r.floating(r.field(load, "five_minutes", "/load/five_minutes"), "/load/five_minutes", 0, math.Inf(1))
	s.Load.FifteenMinutes = r.floating(r.field(load, "fifteen_minutes", "/load/fifteen_minutes"), "/load/fifteen_minutes", 0, math.Inf(1))
	s.Memory = r.group(root, "memory")
	s.RootFilesystem = r.group(root, "root_filesystem")
	s.UptimeSeconds = r.integer(r.field(root, "uptime_seconds", "/uptime_seconds"), "/uptime_seconds")
	pi := r.object(r.field(root, "raspberry_pi", "/raspberry_pi"), "/raspberry_pi")
	s.RaspberryPi.SoCTemperatureCelsius = r.floating(r.field(pi, "soc_temperature_celsius", "/raspberry_pi/soc_temperature_celsius"), "/raspberry_pi/soc_temperature_celsius", math.Inf(-1), math.Inf(1))
	flags := []struct {
		key    string
		target **bool
	}{
		{"thermal_throttling_active", &s.RaspberryPi.ThermalThrottlingActive},
		{"thermal_throttling_occurred_since_boot", &s.RaspberryPi.ThermalThrottlingOccurredSinceBoot},
		{"undervoltage_active", &s.RaspberryPi.UndervoltageActive},
		{"undervoltage_occurred_since_boot", &s.RaspberryPi.UndervoltageOccurredSinceBoot},
	}
	for _, f := range flags {
		p := "/raspberry_pi/" + f.key
		*f.target = r.boolean(r.field(pi, f.key, p), p)
	}
	network := r.field(root, "network", "/network")
	if network != nil {
		n := r.object(network, "/network")
		entries, ok := r.field(n, "interfaces", "/network/interfaces").([]any)
		if !ok || len(entries) > 64 {
			r.fail("/network/interfaces", "array of at most 64 interfaces required")
		}
		s.Network = &Network{Interfaces: make([]NetworkInterface, 0, len(entries))}
		for i, entry := range entries {
			p := "/network/interfaces/" + strconv.Itoa(i)
			o := r.object(entry, p)
			ni := NetworkInterface{Name: r.text(r.field(o, "name", p+"/name"), p+"/name")}
			state := r.field(o, "state", p+"/state")
			if state != nil {
				v := r.text(state, p+"/state")
				ni.State = &v
			}
			ni.RXBytes = r.integer(r.field(o, "rx_bytes", p+"/rx_bytes"), p+"/rx_bytes")
			ni.TXBytes = r.integer(r.field(o, "tx_bytes", p+"/tx_bytes"), p+"/tx_bytes")
			s.Network.Interfaces = append(s.Network.Interfaces, ni)
		}
	}
	issues, ok := r.field(root, "issues", "/issues").([]any)
	if !ok {
		r.fail("/issues", "required array")
	}
	s.Issues = make([]Issue, 0, len(issues))
	for i, value := range issues {
		p := "/issues/" + strconv.Itoa(i)
		o := r.object(value, p)
		s.Issues = append(s.Issues, Issue{r.text(r.field(o, "path", p+"/path"), p+"/path"), r.text(r.field(o, "code", p+"/code"), p+"/code"), r.text(r.field(o, "message", p+"/message"), p+"/message")})
	}
	if r.err != nil {
		return Snapshot{}, r.err
	}
	if err := invariants(s); err != nil {
		return Snapshot{}, err
	}
	return s, nil
}

func allOrNone(present ...bool) bool {
	for _, p := range present {
		if p != present[0] {
			return false
		}
	}
	return true
}
func invariants(s Snapshot) error {
	if s.CPU.LogicalCPUCount != nil && *s.CPU.LogicalCPUCount == 0 {
		return invalid("/cpu/logical_cpu_count", "must be positive")
	}
	if !allOrNone(s.Load.OneMinute != nil, s.Load.FiveMinutes != nil, s.Load.FifteenMinutes != nil) {
		return invalid("/load", "all leaves available or null")
	}
	for _, g := range []struct {
		path  string
		value ByteGroup
	}{{"/memory", s.Memory}, {"/root_filesystem", s.RootFilesystem}} {
		v := g.value
		if !allOrNone(v.TotalBytes != nil, v.AvailableBytes != nil, v.UsedBytes != nil) {
			return invalid(g.path, "all leaves available or null")
		}
		if v.TotalBytes != nil && (*v.AvailableBytes > *v.TotalBytes || *v.UsedBytes != *v.TotalBytes-*v.AvailableBytes) {
			return invalid(g.path, "inconsistent byte group")
		}
	}
	pi := s.RaspberryPi
	if !allOrNone(pi.ThermalThrottlingActive != nil, pi.ThermalThrottlingOccurredSinceBoot != nil, pi.UndervoltageActive != nil, pi.UndervoltageOccurredSinceBoot != nil) {
		return invalid("/raspberry_pi", "firmware flags must be available together")
	}
	expected := []string{}
	useful := false
	add := func(available bool, path string) {
		if available {
			useful = true
		} else {
			expected = append(expected, path)
		}
	}
	add(s.UptimeSeconds != nil, "/uptime_seconds")
	add(s.CPU.LogicalCPUCount != nil, "/cpu/logical_cpu_count")
	add(s.Load.OneMinute != nil, "/load")
	add(s.Memory.TotalBytes != nil, "/memory")
	add(s.RootFilesystem.TotalBytes != nil, "/root_filesystem")
	add(pi.SoCTemperatureCelsius != nil, "/raspberry_pi/soc_temperature_celsius")
	add(pi.ThermalThrottlingActive != nil, "/raspberry_pi/thermal_throttling_active")
	add(pi.ThermalThrottlingOccurredSinceBoot != nil, "/raspberry_pi/thermal_throttling_occurred_since_boot")
	add(pi.UndervoltageActive != nil, "/raspberry_pi/undervoltage_active")
	add(pi.UndervoltageOccurredSinceBoot != nil, "/raspberry_pi/undervoltage_occurred_since_boot")
	if s.Network == nil {
		expected = append(expected, "/network")
	} else {
		previous := ""
		for i, n := range s.Network.Interfaces {
			if i > 0 && strings.Compare(previous, n.Name) >= 0 {
				return invalid("/network/interfaces", "names must be unique and sorted by UTF-8 bytes")
			}
			previous = n.Name
			useful = true
			p := "/network/interfaces/" + strconv.Itoa(i)
			add(n.State != nil, p+"/state")
			add(n.RXBytes != nil, p+"/rx_bytes")
			add(n.TXBytes != nil, p+"/tx_bytes")
		}
	}
	add(s.CPU.UtilizationPercent != nil, "/cpu/utilization_percent")
	if len(expected) != len(s.Issues) {
		return invalid("/issues", "issue count does not match unavailable observations")
	}
	var firmware *Issue
	for i, issue := range s.Issues {
		if issue.Path != expected[i] {
			return invalid("/issues/"+strconv.Itoa(i)+"/path", "unexpected issue path or order")
		}
		switch issue.Code {
		case "unsupported", "permission_denied", "temporarily_unavailable":
		default:
			return invalid("/issues/"+strconv.Itoa(i)+"/code", "unsupported issue code")
		}
		if strings.HasPrefix(issue.Path, "/raspberry_pi/") && issue.Path != "/raspberry_pi/soc_temperature_celsius" {
			if firmware == nil {
				copy := issue
				firmware = &copy
			} else if issue.Code != firmware.Code || issue.Message != firmware.Message {
				return invalid("/issues", "firmware issues must share code and message")
			}
		}
	}
	if !useful {
		return invalid("/", "no useful observation")
	}
	return nil
}
