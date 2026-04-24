package enricher_test

import (
	"testing"

	"github.com/user/portwatch/internal/enricher"
)

func TestDefaultConfigValues(t *testing.T) {
	c := enricher.DefaultConfig()
	if c.Enabled {
		t.Error("default config should be disabled")
	}
	if len(c.Fields) != 0 {
		t.Error("default config should have no fields")
	}
}

func TestValidateRejectsEmptyKey(t *testing.T) {
	c := enricher.Config{
		Enabled: true,
		Fields:  map[string]string{"": "value"},
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty field key")
	}
}

func TestValidateAcceptsValidConfig(t *testing.T) {
	c := enricher.Config{
		Enabled: true,
		Fields:  map[string]string{"host": "box1"},
	}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	e, err := enricher.NewFromConfig(enricher.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Error("expected nil enricher when disabled")
	}
}

func TestNewFromConfigEnabledCreatesEnricher(t *testing.T) {
	c := enricher.Config{
		Enabled: true,
		Fields:  map[string]string{"dc": "eu-west"},
	}
	e, err := enricher.NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enricher")
	}
	if e.Fields()["dc"] != "eu-west" {
		t.Errorf("expected dc=eu-west, got %q", e.Fields()["dc"])
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	c := enricher.Config{
		Enabled: true,
		Fields:  map[string]string{"": "bad"},
	}
	_, err := enricher.NewFromConfig(c)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}
