package pipeline_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := pipeline.DefaultConfig()
	if cfg.DebounceWait != 0 {
		t.Errorf("expected DebounceWait 0, got %v", cfg.DebounceWait)
	}
	if cfg.ThrottleEnabled {
		t.Error("expected ThrottleEnabled false")
	}
	if cfg.ThrottleInterval != 500*time.Millisecond {
		t.Errorf("expected ThrottleInterval 500ms, got %v", cfg.ThrottleInterval)
	}
}

func TestValidateAcceptsDefault(t *testing.T) {
	if err := pipeline.DefaultConfig().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsNegativeDebounce(t *testing.T) {
	cfg := pipeline.DefaultConfig()
	cfg.DebounceWait = -1 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative debounce_wait")
	}
}

func TestValidateRejectsThrottleWithZeroInterval(t *testing.T) {
	cfg := pipeline.DefaultConfig()
	cfg.ThrottleEnabled = true
	cfg.ThrottleInterval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for throttle enabled with zero interval")
	}
}
