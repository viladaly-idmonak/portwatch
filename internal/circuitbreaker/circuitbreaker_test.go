package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestInitialStateIsClosed(t *testing.T) {
	cb := circuitbreaker.New(3, time.Second)
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed, got %v", cb.State())
	}
}

func TestAllowPassesWhenClosed(t *testing.T) {
	cb := circuitbreaker.New(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpensAfterMaxFailures(t *testing.T) {
	cb := circuitbreaker.New(3, time.Second)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected Open, got %v", cb.State())
	}
	if err := cb.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestTransitionsToHalfOpenAfterTimeout(t *testing.T) {
	cb := circuitbreaker.New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil after reset timeout, got %v", err)
	}
	if cb.State() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", cb.State())
	}
}

func TestSuccessInHalfOpenCloses(t *testing.T) {
	cb := circuitbreaker.New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow()
	cb.RecordSuccess()
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after success, got %v", cb.State())
	}
}

func TestSuccessResetFailureCount(t *testing.T) {
	cb := circuitbreaker.New(3, time.Second)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	cb.RecordFailure()
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after reset, got %v", cb.State())
	}
}
