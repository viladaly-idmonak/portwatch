package classifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/classifier"
	"github.com/user/portwatch/internal/scanner"
)

func entry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto}
}

func TestClassifyKnownPortReturnsConfiguredClass(t *testing.T) {
	c, err := classifier.New([]classifier.Rule{
		{Port: 22, Protocol: "tcp", Class: classifier.ClassTrusted},
		{Port: 4444, Protocol: "tcp", Class: classifier.ClassCritical},
	}, classifier.ClassUnknown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Classify(entry(22, "tcp")); got != classifier.ClassTrusted {
		t.Errorf("expected trusted, got %q", got)
	}
	if got := c.Classify(entry(4444, "tcp")); got != classifier.ClassCritical {
		t.Errorf("expected critical, got %q", got)
	}
}

func TestClassifyUnknownPortReturnsDefault(t *testing.T) {
	c, _ := classifier.New(nil, classifier.ClassSuspect)
	if got := c.Classify(entry(9999, "tcp")); got != classifier.ClassSuspect {
		t.Errorf("expected suspect, got %q", got)
	}
}

func TestClassifyProtocolDistinguished(t *testing.T) {
	c, _ := classifier.New([]classifier.Rule{
		{Port: 53, Protocol: "tcp", Class: classifier.ClassTrusted},
		{Port: 53, Protocol: "udp", Class: classifier.ClassSuspect},
	}, classifier.ClassUnknown)
	if got := c.Classify(entry(53, "tcp")); got != classifier.ClassTrusted {
		t.Errorf("tcp/53: expected trusted, got %q", got)
	}
	if got := c.Classify(entry(53, "udp")); got != classifier.ClassSuspect {
		t.Errorf("udp/53: expected suspect, got %q", got)
	}
}

func TestNewEmptyDefaultClassReturnsError(t *testing.T) {
	_, err := classifier.New(nil, "")
	if err == nil {
		t.Fatal("expected error for empty default class")
	}
}

func TestAddRuleUpdatesClassification(t *testing.T) {
	c, _ := classifier.New(nil, classifier.ClassUnknown)
	if got := c.Classify(entry(8080, "tcp")); got != classifier.ClassUnknown {
		t.Fatalf("pre-add: expected unknown, got %q", got)
	}
	c.AddRule(classifier.Rule{Port: 8080, Protocol: "tcp", Class: classifier.ClassTrusted})
	if got := c.Classify(entry(8080, "tcp")); got != classifier.ClassTrusted {
		t.Errorf("post-add: expected trusted, got %q", got)
	}
}
