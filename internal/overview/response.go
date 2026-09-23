package overview

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Adrien-hue/joy-pi-board/internal/health"
)

const MaximumEnvelopeBytes = 81920
const MaximumEnvelopeDepth = 34
const MaximumAgeMilliseconds uint64 = 9007199254740991
const AgeHeader = "X-Joy-Pi-Snapshot-Age-Ms"

// Prepared owns compact private bytes. It has no aliases to input or cache.
type Prepared struct {
	view     View
	snapshot json.RawMessage
}

// StableBetween detects a known preparation pause that changes the conservative
// integer age (or validity). It never consults a newer flight/cache. HTTP makes
// at most three assembly attempts; persistent suspension is a delivery failure.
func (p Prepared) StableBetween(a, b Reading) bool {
	if !a.Valid || !b.Valid || b.Tick < a.Tick {
		return false
	}
	if p.snapshot == nil {
		return true
	}
	origin := p.view.cache.success.Tick
	if a.Tick < origin || b.Tick < origin {
		return false
	}
	age := func(d time.Duration) time.Duration {
		q := d / time.Millisecond
		if d%time.Millisecond != 0 {
			q++
		}
		return q
	}
	return age(a.Tick-origin) == age(b.Tick-origin) && (a.Tick-origin <= Freshness) == (b.Tick-origin <= Freshness)
}

type Response struct {
	Body []byte
	Age  string
}
type healthJSON struct {
	Availability  string          `json:"availability"`
	SnapshotState string          `json:"snapshot_state"`
	LastSuccessAt *string         `json:"last_success_at"`
	Reason        *health.Reason  `json:"reason"`
	Snapshot      json.RawMessage `json:"snapshot"`
}
type envelope struct {
	GeneratedAt string     `json:"generated_at"`
	Health      healthJSON `json:"health"`
}

func (v View) Prepare() (Prepared, error) {
	p := Prepared{view: v}
	if v.sequence == 0 {
		return p, ErrClock
	}
	if v.cache.snapshot.Valid() {
		body := v.cache.snapshot.Bytes()
		var buf bytes.Buffer
		if err := json.Compact(&buf, body); err != nil || buf.Len() > len(body) {
			return Prepared{}, ErrClock
		}
		p.snapshot = buf.Bytes()
	}
	return p, nil
}

func timestamp(t time.Time) (string, error) {
	t = t.UTC()
	if t.Year() < 0 || t.Year() > 9999 {
		return "", ErrClock
	}
	return t.Format(time.RFC3339Nano), nil
}

func capMilliseconds(ms uint64) string {
	if ms > MaximumAgeMilliseconds {
		ms = MaximumAgeMilliseconds
	}
	return strconv.FormatUint(ms, 10)
}

// Finalize samples no shared cache: all fields derive from this flight's view
// and one supplied paired reading. Duration arithmetic never uses wall time.
func (p Prepared) Finalize(now Reading) (Response, error) {
	if !now.Valid || now.Tick < 0 {
		return Response{}, ErrClock
	}
	generated, err := timestamp(now.Wall)
	if err != nil {
		return Response{}, err
	}
	h := healthJSON{Availability: "unavailable", SnapshotState: "none"}
	if p.view.reason != "" {
		reason := p.view.reason
		h.Reason = &reason
	}
	age := ""
	if p.snapshot != nil {
		origin := p.view.cache.success
		if !origin.Valid || origin.Tick < 0 || now.Tick < origin.Tick {
			return Response{}, ErrClock
		}
		elapsed := now.Tick - origin.Tick
		ms := uint64(elapsed / time.Millisecond)
		if elapsed%time.Millisecond != 0 {
			ms++
		}
		age = capMilliseconds(ms)
		last, err := timestamp(origin.Wall)
		if err != nil {
			return Response{}, err
		}
		h.LastSuccessAt = &last
		if elapsed <= Freshness {
			h.Snapshot = p.snapshot
			h.SnapshotState = "stale"
		}
		if p.view.reason == "" {
			if elapsed > Freshness {
				return Response{}, ErrUnavailable
			}
			h.Availability = "available"
			h.SnapshotState = "current"
		}
	} else if p.view.reason == "" {
		return Response{}, ErrClock
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	// Encode the bounded wrapper first, then insert already compact validated
	// JSON at its fixed final field. Re-encoding/scanning 64 KiB after the clock
	// cut would unnecessarily extend preparation, especially under race tooling.
	snapshot := h.Snapshot
	h.Snapshot = nil
	if err := encoder.Encode(envelope{generated, h}); err != nil {
		return Response{}, ErrClock
	}
	if buf.Len() > 1024 {
		return Response{}, ErrClock
	}
	body := buf.Bytes()
	if snapshot != nil {
		const suffix = "null}}\n"
		if !bytes.HasSuffix(body, []byte(suffix)) {
			return Response{}, ErrClock
		}
		result := make([]byte, 0, len(body)-4+len(snapshot))
		result = append(result, body[:len(body)-len(suffix)]...)
		result = append(result, snapshot...)
		body = append(result, "}}\n"...)
	}
	if len(body) > MaximumEnvelopeBytes || len(body)-len(snapshot) > 1024 {
		return Response{}, ErrClock
	}
	return Response{body, age}, nil
}
