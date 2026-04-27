package labeler_test

import (
	"testing"

	"github.com/user/portwatch/internal/labeler"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := labeler.DefaultConfig()
	if cfg.Enabled {
		t.Fatal("default config should be disabled")
	}
	if len(cfg.Rules) != 0 {
		t.Fatal("default config should have no rules")
	}
}

func TestValidateRejectsZeroPort(t *testing.T) {
	cfg := labeler.Config{
		Enabled: true,
		Rules:   []labeler.RuleConfig{{Port: 0, Protocol: "tcp", Label: "x"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero port")
	}
}

func TestValidateRejectsEmptyProtocol(t *testing.T) {
	cfg := labeler.Config{
		Enabled: true,
		Rules:   []labeler.RuleConfig{{Port: 80, Protocol: "", Label: "x"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty protocol")
	}
}

func TestValidateRejectsEmptyLabel(t *testing.T) {
	cfg := labeler.Config{
		Enabled: true,
		Rules:   []labeler.RuleConfig{{Port: 80, Protocol: "tcp", Label: ""}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	cfg := labeler.DefaultConfig()
	l, err := labeler.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l != nil {
		t.Fatal("expected nil labeler when disabled")
	}
}

func TestNewFromConfigEnabledCreatesLabeler(t *testing.T) {
	cfg := labeler.Config{
		Enabled: true,
		Rules: []labeler.RuleConfig{
			{Port: 22, Protocol: "tcp", Label: "ssh"},
		},
	}
	l, err := labeler.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil labeler")
	}
	if got := l.Label(22, "tcp"); got != "ssh" {
		t.Fatalf("want ssh, got %q", got)
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	cfg := labeler.Config{
		Enabled: true,
		Rules:   []labeler.RuleConfig{{Port: 0, Protocol: "tcp", Label: "bad"}},
	}
	_, err := labeler.NewFromConfig(cfg)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
