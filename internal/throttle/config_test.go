package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := throttle.DefaultConfig()
	if cfg.Interval != 500*time.Millisecond {
		t.Errorf("expected 500ms, got %s", cfg.Interval)
	}
}

func TestValidateRejectsZeroInterval(t *testing.T) {
	cfg := throttle.Config{Interval: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidateRejectsNegativeInterval(t *testing.T) {
	cfg := throttle.Config{Interval: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestNewFromConfigValid(t *testing.T) {
	cfg := throttle.Config{Interval: 100 * time.Millisecond}
	th, err := throttle.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !th.Allow() {
		t.Fatal("expected first allow from config-created throttle")
	}
}

func TestNewFromConfigInvalid(t *testing.T) {
	cfg := throttle.Config{Interval: 0}
	_, err := throttle.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
