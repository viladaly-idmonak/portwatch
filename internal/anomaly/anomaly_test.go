package anomaly

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func entry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewInvalidThresholdReturnsError(t *testing.T) {
	_, err := New(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestFirstEventsUnderThresholdNotFlagged(t *testing.T) {
	d, _ := New(time.Minute, 3)
	e := entry(80, "tcp")
	for i := 0; i < 3; i++ {
		if d.Record(e) {
			t.Fatalf("iteration %d: expected false, got true", i)
		}
	}
}

func TestExceedingThresholdFlagged(t *testing.T) {
	d, _ := New(time.Minute, 3)
	e := entry(443, "tcp")
	for i := 0; i < 3; i++ {
		d.Record(e)
	}
	if !d.Record(e) {
		t.Fatal("expected anomaly flag after exceeding threshold")
	}
}

func TestEventsOutsideWindowNotCounted(t *testing.T) {
	d, _ := New(50*time.Millisecond, 2)
	now := time.Now()
	calls := 0
	d.now = func() time.Time {
		calls++
		if calls <= 2 {
			return now.Add(-100 * time.Millisecond)
		}
		return now
	}
	e := entry(22, "tcp")
	d.Record(e)
	d.Record(e)
	// Both prior events are outside the window; this should be the first.
	if d.Record(e) {
		t.Fatal("expected no anomaly; old events should be evicted")
	}
}

func TestResetClearsHistory(t *testing.T) {
	d, _ := New(time.Minute, 2)
	e := entry(8080, "tcp")
	d.Record(e)
	d.Record(e)
	d.Reset(e)
	if d.Record(e) {
		t.Fatal("expected no anomaly after reset")
	}
}

func TestDistinctPortsAreIndependent(t *testing.T) {
	d, _ := New(time.Minute, 1)
	a := entry(80, "tcp")
	b := entry(443, "tcp")
	d.Record(a)
	d.Record(a) // flags a
	if d.Record(b) {
		t.Fatal("port 443 should not be flagged by port 80 events")
	}
}
