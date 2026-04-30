package quota

import (
	"testing"
	"time"
)

func TestNewInvalidMaxReturnsError(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestFirstEventAlwaysAllowed(t *testing.T) {
	q, _ := New(3, time.Second)
	if err := q.Allow("k"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExceedingMaxIsBlocked(t *testing.T) {
	q, _ := New(2, time.Second)
	if err := q.Allow("k"); err != nil {
		t.Fatal(err)
	}
	if err := q.Allow("k"); err != nil {
		t.Fatal(err)
	}
	if err := q.Allow("k"); err != ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestWindowResetAllowsNewEvents(t *testing.T) {
	q, _ := New(1, 20*time.Millisecond)
	if err := q.Allow("k"); err != nil {
		t.Fatal(err)
	}
	if err := q.Allow("k"); err != ErrQuotaExceeded {
		t.Fatalf("expected exceeded, got %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if err := q.Allow("k"); err != nil {
		t.Fatalf("expected reset, got %v", err)
	}
}

func TestDistinctKeysAreIndependent(t *testing.T) {
	q, _ := New(1, time.Second)
	if err := q.Allow("a"); err != nil {
		t.Fatal(err)
	}
	if err := q.Allow("b"); err != nil {
		t.Fatal(err)
	}
}

func TestRemainingDecrements(t *testing.T) {
	q, _ := New(3, time.Second)
	if got := q.Remaining("k"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
	_ = q.Allow("k")
	if got := q.Remaining("k"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestPurgeRemovesExpiredBuckets(t *testing.T) {
	q, _ := New(5, 10*time.Millisecond)
	_ = q.Allow("x")
	time.Sleep(20 * time.Millisecond)
	q.Purge()
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.buckets) != 0 {
		t.Fatalf("expected empty buckets after purge, got %d", len(q.buckets))
	}
}
