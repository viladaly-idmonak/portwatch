package trend

import (
	"testing"
)

func TestTrendDefaultConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Window != 10 {
		t.Errorf("expected Window 10, got %d", cfg.Window)
	}
}

func TestTrendValidateDisabledAlwaysPasses(t *testing.T) {
	cfg := Config{Enabled: false, Window: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestTrendValidateEnabledRejectsZeroWindow(t *testing.T) {
	cfg := Config{Enabled: true, Window: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero window, got nil")
	}
}

func TestTrendValidateEnabledRejectsNegativeWindow(t *testing.T) {
	cfg := Config{Enabled: true, Window: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative window, got nil")
	}
}

func TestTrendValidateEnabledAcceptsPositiveWindow(t *testing.T) {
	cfg := Config{Enabled: true, Window: 5}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTrendNewFromConfigDisabledReturnsNil(t *testing.T) {
	cfg := Config{Enabled: false}
	tracker, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tracker != nil {
		t.Error("expected nil tracker for disabled config")
	}
}

func TestTrendNewFromConfigInvalidReturnsError(t *testing.T) {
	cfg := Config{Enabled: true, Window: -5}
	_, err := NewFromConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}

func TestTrendNewFromConfigValidCreatesTracker(t *testing.T) {
	cfg := Config{Enabled: true, Window: 3}
	tracker, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tracker == nil {
		t.Error("expected non-nil tracker")
	}
}
