package metrics

import (
	"testing"
	"time"
)

func TestInitialSnapshotIsZero(t *testing.T) {
	c := New()
	s := c.Snapshot()
	if s.Opened != 0 || s.Closed != 0 || s.Scans != 0 {
		t.Fatalf("expected zero counters, got %+v", s)
	}
}

func TestRecordScanIncrements(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	s := c.Snapshot()
	if s.Scans != 2 {
		t.Fatalf("expected 2 scans, got %d", s.Scans)
	}
	if s.LastScan.IsZero() {
		t.Fatal("expected LastScan to be set")
	}
}

func TestRecordOpenedAndClosed(t *testing.T) {
	c := New()
	c.RecordOpened(3)
	c.RecordClosed(1)
	s := c.Snapshot()
	if s.Opened != 3 {
		t.Fatalf("expected 3 opened, got %d", s.Opened)
	}
	if s.Closed != 1 {
		t.Fatalf("expected 1 closed, got %d", s.Closed)
	}
}

func TestUptimeGrows(t *testing.T) {
	c := New()
	time.Sleep(10 * time.Millisecond)
	s := c.Snapshot()
	if s.Uptime < 10*time.Millisecond {
		t.Fatalf("expected uptime >= 10ms, got %v", s.Uptime)
	}
}

func TestConcurrentRecordIsSafe(t *testing.T) {
	c := New()
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.RecordScan()
			c.RecordOpened(1)
			c.RecordClosed(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	s := c.Snapshot()
	if s.Scans != 50 {
		t.Fatalf("expected 50 scans, got %d", s.Scans)
	}
}
