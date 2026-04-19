package pipeline

import (
	"testing"
)

func TestDedupDefaultConfigValues(t *testing.T) {
	cfg := DefaultDedupConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if cfg.MaxEntries != 0 {
		t.Errorf("expected MaxEntries 0, got %d", cfg.MaxEntries)
	}
}

func TestDedupValidateAcceptsDefault(t *testing.T) {
	if err := DefaultDedupConfig().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDedupValidateRejectsNegativeMaxEntries(t *testing.T) {
	cfg := DedupConfig{Enabled: true, MaxEntries: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative MaxEntries")
	}
}

func TestDedupNewFromConfigValid(t *testing.T) {
	cfg, err := NewDedupFromConfig(true, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected Enabled true")
	}
	if cfg.MaxEntries != 100 {
		t.Errorf("expected MaxEntries 100, got %d", cfg.MaxEntries)
	}
}

func TestDedupNewFromConfigInvalid(t *testing.T) {
	_, err := NewDedupFromConfig(true, -5)
	if err == nil {
		t.Error("expected error for negative MaxEntries")
	}
}

func TestDedupNewFromConfigDisabled(t *testing.T) {
	cfg, err := NewDedupFromConfig(false, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected Enabled false")
	}
}
