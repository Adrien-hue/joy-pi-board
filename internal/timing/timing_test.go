package timing

import (
	"context"
	"testing"
	"time"
)

func TestBoundsCorrelationAndIncomplete(t *testing.T) {
	var records []Record
	r := New(func(v Record) bool { records = append(records, v); return true })
	ctx, q := Begin(context.Background())
	f := NewFlight(ctx, 17)
	Associate(ctx, f)
	start := time.Now()
	end := Network(WithFlight(ctx, f))
	end(true)
	end(false)
	r.Finish(q, start, 200, false, true)
	if len(records) != 1 || records[0].Flight != 17 || !records[0].Eligible {
		t.Fatal(records)
	}
	f.Share()
	r.Finish(q, start, 200, false, true)
	if records[1].Eligible || !records[1].Shared {
		t.Fatal(records[1])
	}
	_, incomplete := Begin(context.Background())
	r.Finish(incomplete, start, 200, false, true)
	if records[2].Eligible || records[2].Complete {
		t.Fatal(records[2])
	}
	r.Finish(q, start, 200, true, true)
	if records[3].Eligible || !records[3].Cancelled {
		t.Fatal(records[3])
	}
	for i := 4; i < MaximumRecords+9; i++ {
		r.Finish(q, start, 200, false, true)
	}
	s := r.Summary()
	if len(records) != MaximumRecords || s.Truncated != 9 || s.Queued != MaximumRecords {
		t.Fatal(s)
	}
	blocked := New(func(Record) bool { return false })
	blocked.Finish(q, start, 200, false, true)
	if blocked.Summary().Dropped != 1 {
		t.Fatal(blocked.Summary())
	}
}
