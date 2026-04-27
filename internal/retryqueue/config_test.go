package retryqueue_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/retryqueue"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	if cfg.MaxSize <= 0 {
		t.Errorf("MaxSize should be positive, got %d", cfg.MaxSize)
	}
	if cfg.MaxAttempts <= 0 {
		t.Errorf("MaxAttempts should be positive, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay <= 0 {
		t.Errorf("BaseDelay should be positive, got %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay < cfg.BaseDelay {
		t.Errorf("MaxDelay should be >= BaseDelay")
	}
}

func TestValidateRejectsZeroMaxSize(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.MaxSize = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for MaxSize=0")
	}
}

func TestValidateRejectsZeroMaxAttempts(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.MaxAttempts = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for MaxAttempts=0")
	}
}

func TestValidateRejectsMaxDelayBelowBaseDelay(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.BaseDelay = time.Second
	cfg.MaxDelay = time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when MaxDelay < BaseDelay")
	}
}

func TestNewFromConfigValidReturnsConfig(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	out, err := retryqueue.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.MaxSize != cfg.MaxSize {
		t.Errorf("MaxSize mismatch: got %d, want %d", out.MaxSize, cfg.MaxSize)
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.BaseDelay = -1
	_, err := retryqueue.NewFromConfig(cfg)
	if err == nil {
		t.Error("expected error for negative BaseDelay")
	}
}
