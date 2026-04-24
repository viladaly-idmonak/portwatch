package sampling_test

import (
	"testing"

	"github.com/user/portwatch/internal/sampling"
)

func TestDefaultConfigValues(t *testing.T) {
	c := sampling.DefaultConfig()
	if c.Enabled {
		t.Error("default config should be disabled")
	}
	if c.Rate != 1.0 {
		t.Errorf("default rate should be 1.0, got %v", c.Rate)
	}
}

func TestValidateDisabledAlwaysPasses(t *testing.T) {
	c := sampling.Config{Enabled: false, Rate: -5}
	if err := c.Validate(); err != nil {
		t.Fatalf("disabled config should always validate, got: %v", err)
	}
}

func TestValidateEnabledInvalidRate(t *testing.T) {
	c := sampling.Config{Enabled: true, Rate: 0}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestValidateEnabledValidRate(t *testing.T) {
	c := sampling.Config{Enabled: true, Rate: 0.25}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	s, err := sampling.NewFromConfig(sampling.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Fatal("expected nil sampler when disabled")
	}
}

func TestNewFromConfigEnabledCreatesSampler(t *testing.T) {
	c := sampling.Config{Enabled: true, Rate: 0.5}
	s, err := sampling.NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	c := sampling.Config{Enabled: true, Rate: 2.0}
	_, err := sampling.NewFromConfig(c)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
