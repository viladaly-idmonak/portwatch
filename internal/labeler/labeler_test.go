package labeler_test

import (
	"testing"

	"github.com/user/portwatch/internal/labeler"
)

func rules() []labeler.Rule {
	return []labeler.Rule{
		{Port: 80, Protocol: "tcp", Label: "http"},
		{Port: 443, Protocol: "tcp", Label: "https"},
		{Port: 53, Protocol: "udp", Label: "dns"},
	}
}

func TestLabelKnownPort(t *testing.T) {
	l := labeler.New(rules())
	if got := l.Label(80, "tcp"); got != "http" {
		t.Fatalf("want http, got %q", got)
	}
}

func TestLabelUnknownPortReturnsEmpty(t *testing.T) {
	l := labeler.New(rules())
	if got := l.Label(9999, "tcp"); got != "" {
		t.Fatalf("want empty, got %q", got)
	}
}

func TestLabelOrDefaultFallback(t *testing.T) {
	l := labeler.New(rules())
	if got := l.LabelOrDefault(9999, "tcp", "unknown"); got != "unknown" {
		t.Fatalf("want unknown, got %q", got)
	}
}

func TestLabelOrDefaultUsesLabel(t *testing.T) {
	l := labeler.New(rules())
	if got := l.LabelOrDefault(443, "tcp", "fallback"); got != "https" {
		t.Fatalf("want https, got %q", got)
	}
}

func TestAllReturnsCopy(t *testing.T) {
	l := labeler.New(rules())
	all := l.All()
	if len(all) != 3 {
		t.Fatalf("want 3 rules, got %d", len(all))
	}
	// mutating copy must not affect labeler
	all["80/tcp"] = "mutated"
	if got := l.Label(80, "tcp"); got != "http" {
		t.Fatalf("labeler was mutated: got %q", got)
	}
}

func TestNewSkipsInvalidRules(t *testing.T) {
	invalid := []labeler.Rule{
		{Port: 0, Protocol: "tcp", Label: "bad-port"},
		{Port: 80, Protocol: "", Label: "bad-proto"},
		{Port: 80, Protocol: "tcp", Label: ""},
		{Port: 8080, Protocol: "tcp", Label: "custom"},
	}
	l := labeler.New(invalid)
	if got := l.Label(8080, "tcp"); got != "custom" {
		t.Fatalf("want custom, got %q", got)
	}
	if n := len(l.All()); n != 1 {
		t.Fatalf("want 1 valid rule, got %d", n)
	}
}
