// Package healthschema consumes the pinned Health 1.0 contract in memory.
// It performs no HTTP requests, observation, caching or polling.
package healthschema

import (
	"bytes"
	"encoding/json"
)

const MaximumBytes = 65536
const MaximumDepth = 32

type Classification string

const (
	InvalidResponse          Classification = "invalid_response"
	UnsupportedSchemaVersion Classification = "unsupported_schema_version"
)

// Error contains a bounded diagnostic, never the untrusted value or payload.
type Error struct {
	Kind       Classification
	Path, Rule string
}

func (e *Error) Error() string        { return string(e.Kind) + " at " + e.Path + ": " + e.Rule }
func invalid(path, rule string) error { return &Error{InvalidResponse, path, rule} }

type Host struct {
	Hostname string `json:"hostname"`
}
type CPU struct {
	UtilizationPercent *float64 `json:"utilization_percent"`
	LogicalCPUCount    *uint64  `json:"logical_cpu_count"`
}
type Load struct {
	OneMinute      *float64 `json:"one_minute"`
	FiveMinutes    *float64 `json:"five_minutes"`
	FifteenMinutes *float64 `json:"fifteen_minutes"`
}
type ByteGroup struct {
	TotalBytes     *uint64 `json:"total_bytes"`
	AvailableBytes *uint64 `json:"available_bytes"`
	UsedBytes      *uint64 `json:"used_bytes"`
}
type NetworkInterface struct {
	Name    string  `json:"name"`
	State   *string `json:"state"`
	RXBytes *uint64 `json:"rx_bytes"`
	TXBytes *uint64 `json:"tx_bytes"`
}
type Network struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}
type RaspberryPi struct {
	SoCTemperatureCelsius              *float64 `json:"soc_temperature_celsius"`
	ThermalThrottlingActive            *bool    `json:"thermal_throttling_active"`
	ThermalThrottlingOccurredSinceBoot *bool    `json:"thermal_throttling_occurred_since_boot"`
	UndervoltageActive                 *bool    `json:"undervoltage_active"`
	UndervoltageOccurredSinceBoot      *bool    `json:"undervoltage_occurred_since_boot"`
}
type Issue struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Snapshot is a projection of known fields AFTER presence and runtime validation.
// Pointers here mean explicit null, never missing. It is not the transport DTO:
// only Validated retains unknown members and original number tokens.
type Snapshot struct {
	SchemaVersion  string      `json:"schema_version"`
	ObservedAt     string      `json:"observed_at"`
	Host           Host        `json:"host"`
	CPU            CPU         `json:"cpu"`
	Load           Load        `json:"load"`
	Memory         ByteGroup   `json:"memory"`
	RootFilesystem ByteGroup   `json:"root_filesystem"`
	UptimeSeconds  *uint64     `json:"uptime_seconds"`
	Network        *Network    `json:"network"`
	RaspberryPi    RaspberryPi `json:"raspberry_pi"`
	Issues         []Issue     `json:"issues"`
}

// Validated owns only private, immutable bytes. Copies of this value may be
// shared concurrently. Every returned byte slice and projection is independent.
type Validated struct{ body []byte }

func (v Validated) Valid() bool   { return len(v.body) != 0 }
func (v Validated) Bytes() []byte { return bytes.Clone(v.body) }
func (v Validated) MarshalJSON() ([]byte, error) {
	if !v.Valid() {
		return nil, invalid("/", "uninitialized validated snapshot")
	}
	return v.Bytes(), nil
}

// Snapshot returns a fresh known-field projection; mutations cannot change v.
func (v Validated) Snapshot() (Snapshot, error) {
	if !v.Valid() {
		return Snapshot{}, invalid("/", "uninitialized validated snapshot")
	}
	// Re-read our own validated bytes using the same exact token decoder. This
	// avoids maintaining an alias-prone clone graph as the schema evolves.
	root, err := document(v.body)
	if err != nil {
		return Snapshot{}, err
	}
	return known(root)
}

// Decode validates the entire bounded document before interpreting its version.
// The caller must not mutate input concurrently with this call.
func Decode(input []byte) (Validated, error) {
	if len(input) > MaximumBytes {
		return Validated{}, invalid("/", "size exceeds 65536 bytes")
	}
	body := bytes.Clone(input)
	root, err := document(body)
	if err != nil {
		return Validated{}, err
	}
	if _, err = known(root); err != nil {
		return Validated{}, err
	}
	return Validated{body}, nil
}

var _ json.Marshaler = Validated{}
