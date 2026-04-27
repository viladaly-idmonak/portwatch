package pipeline_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
)

func TestSuppressDefaultConfigValues(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Window != 30*time.Second {
		t.Errorf("expected Window=30s, got %s", c.Window)
	}
}

func TestSuppressValidateDisabledAlwaysPasses(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = false
	c.Window = 0 // would be invalid if enabled
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestSuppressValidateRejectsZeroWindow(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = true
	c.Window = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero window, got nil")
	}
}

func TestSuppressValidateRejectsNegativeWindow(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = true
	c.Window = -1 * time.Second
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative window, got nil")
	}
}

func TestNewSuppressFromConfigDisabledReturnsNil(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = false
	s, err := pipeline.NewSuppressFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Error("expected nil suppressor when disabled")
	}
}

func TestNewSuppressFromConfigEnabledCreatesSuppressor(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = true
	c.Window = 5 * time.Second
	s, err := pipeline.NewSuppressFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Error("expected non-nil suppressor when enabled")
	}
}

func TestNewSuppressFromConfigInvalidReturnsError(t *testing.T) {
	c := pipeline.DefaultSuppressConfig()
	c.Enabled = true
	c.Window = 0
	_, err := pipeline.NewSuppressFromConfig(c)
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}
