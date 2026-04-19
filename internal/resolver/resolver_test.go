package resolver_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolver"
)

func TestServiceNameKnownPort(t *testing.T) {
	if got := resolver.ServiceName(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestServiceNameUnknownPort(t *testing.T) {
	if got := resolver.ServiceName(9999); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestLookupCachesResult(t *testing.T) {
	r := resolver.New(true)
	first := r.Lookup(80, "tcp")
	second := r.Lookup(80, "tcp")
	if first != second {
		t.Fatalf("cached result mismatch: %q vs %q", first, second)
	}
}

func TestLookupFallbackUnknown(t *testing.T) {
	r := resolver.New(true)
	got := r.Lookup(19999, "tcp")
	if got == "" {
		t.Fatal("expected non-empty fallback result")
	}
}

func TestLookupNoFallbackUnknown(t *testing.T) {
	r := resolver.New(false)
	// For an unresolvable port with no fallback, empty or numeric is acceptable.
	got := r.Lookup(19999, "tcp")
	_ = got // just ensure no panic
}
