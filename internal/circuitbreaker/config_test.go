package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.MaxFailures != 5 {
		t.Errorf("expected MaxFailures=5, got %d", cfg.MaxFailures)
	}
	if cfg.ResetTimeout != 30*time.Second {
		t.Errorf("expected ResetTimeout=30s, got %v", cfg.ResetTimeout)
	}
}

func TestValidateRejectsZeroMaxFailures(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	cfg.MaxFailures = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero max_failures")
	}
}

func TestValidateRejectsZeroResetTimeout(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	cfg.ResetTimeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero reset_timeout")
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	cfg.Enabled = false
	cb, err := circuitbreaker.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb != nil {
		t.Error("expected nil circuit breaker when disabled")
	}
}

func TestNewFromConfigEnabledCreatesCB(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	cb, err := circuitbreaker.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb == nil {
		t.Fatal("expected non-nil circuit breaker")
	}
	if cb.State() != circuitbreaker.StateClosed {
		t.Errorf("expected Closed initial state, got %v", cb.State())
	}
}
