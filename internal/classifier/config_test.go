package classifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/classifier"
)

func TestDefaultConfigValues(t *testing.T) {
	c := classifier.DefaultConfig()
	if c.Enabled {
		t.Error("expected disabled by default")
	}
	if c.DefaultClass != string(classifier.ClassUnknown) {
		t.Errorf("expected default class %q, got %q", classifier.ClassUnknown, c.DefaultClass)
	}
}

func TestValidateDisabledAlwaysPasses(t *testing.T) {
	c := classifier.Config{Enabled: false, DefaultClass: ""}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateEnabledRequiresDefaultClass(t *testing.T) {
	c := classifier.Config{Enabled: true, DefaultClass: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty default_class")
	}
}

func TestValidateEnabledRejectsEmptyRuleProtocol(t *testing.T) {
	c := classifier.Config{
		Enabled:      true,
		DefaultClass: "unknown",
		Rules:        []classifier.RuleConfig{{Port: 80, Protocol: "", Class: "trusted"}},
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty protocol in rule")
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	c := classifier.DefaultConfig()
	clf, err := classifier.NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if clf != nil {
		t.Error("expected nil classifier when disabled")
	}
}

func TestNewFromConfigEnabledCreatesClassifier(t *testing.T) {
	c := classifier.Config{
		Enabled:      true,
		DefaultClass: "unknown",
		Rules: []classifier.RuleConfig{
			{Port: 22, Protocol: "tcp", Class: "trusted"},
		},
	}
	clf, err := classifier.NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if clf == nil {
		t.Fatal("expected non-nil classifier")
	}
}
