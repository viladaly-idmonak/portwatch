package decay_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/decay"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := decay.DefaultConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.HalfLife != 5*time.Minute {
		t.Errorf("expected default HalfLife 5m, got %v", cfg.HalfLife)
	}
}

func TestValidateDisabledAlwaysPasses(t *testing.T) {
	cfg := decay.Config{Enabled: false, HalfLife: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for disabled config, got %v", err)
	}
}

func TestValidateEnabledRejectsZeroHalfLife(t *testing.T) {
	cfg := decay.Config{Enabled: true, HalfLife: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero half-life when enabled")
	}
}

func TestValidateEnabledAcceptsPositiveHalfLife(t *testing.T) {
	cfg := decay.Config{Enabled: true, HalfLife: time.Minute}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	cfg := decay.Config{Enabled: false, HalfLife: time.Minute}
	d, err := decay.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != nil {
		t.Fatal("expected nil Decayer when disabled")
	}
}

func TestNewFromConfigEnabledCreatesDecayer(t *testing.T) {
	cfg := decay.Config{Enabled: true, HalfLife: time.Minute}
	d, err := decay.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Decayer when enabled")
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	cfg := decay.Config{Enabled: true, HalfLife: -time.Second}
	_, err := decay.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
