package baseline_test

import (
	"sort"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func entries(ports ...uint16) []scanner.Entry {
	out := make([]scanner.Entry, len(ports))
	for i, p := range ports {
		out[i] = scanner.Entry{Port: p, Proto: "tcp"}
	}
	return out
}

func sorted(s []uint16) []uint16 {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return s
}

func TestNoDeviationsWhenCurrentMatchesBaseline(t *testing.T) {
	b := baseline.New([]uint16{80, 443})
	u, m := b.Deviations(entries(80, 443))
	if len(u) != 0 || len(m) != 0 {
		t.Fatalf("expected no deviations, got unexpected=%v missing=%v", u, m)
	}
}

func TestUnexpectedPortDetected(t *testing.T) {
	b := baseline.New([]uint16{80})
	u, m := b.Deviations(entries(80, 9999))
	if len(m) != 0 {
		t.Fatalf("expected no missing, got %v", m)
	}
	if len(u) != 1 || u[0] != 9999 {
		t.Fatalf("expected unexpected=[9999], got %v", u)
	}
}

func TestMissingPortDetected(t *testing.T) {
	b := baseline.New([]uint16{80, 443})
	u, m := b.Deviations(entries(80))
	if len(u) != 0 {
		t.Fatalf("expected no unexpected, got %v", u)
	}
	if len(m) != 1 || m[0] != 443 {
		t.Fatalf("expected missing=[443], got %v", m)
	}
}

func TestNewFromScan(t *testing.T) {
	b := baseline.NewFromScan(entries(22, 80, 443))
	for _, p := range []uint16{22, 80, 443} {
		if !b.Contains(p) {
			t.Errorf("expected baseline to contain %d", p)
		}
	}
	if b.Contains(8080) {
		t.Error("expected baseline not to contain 8080")
	}
}

func TestUpdateReplacesBaseline(t *testing.T) {
	b := baseline.New([]uint16{80})
	b.Update([]uint16{443, 8080})
	u, m := b.Deviations(entries(443, 8080))
	if len(u) != 0 || len(m) != 0 {
		t.Fatalf("after update expected no deviations, got u=%v m=%v", u, m)
	}
	if b.Contains(80) {
		t.Error("old port 80 should no longer be in baseline")
	}
}

func TestEmptyBaselineAllCurrentAreUnexpected(t *testing.T) {
	b := baseline.New(nil)
	u, m := b.Deviations(entries(22, 80))
	if len(m) != 0 {
		t.Fatalf("expected no missing, got %v", m)
	}
	if got := sorted(u); len(got) != 2 || got[0] != 22 || got[1] != 80 {
		t.Fatalf("expected unexpected=[22,80], got %v", got)
	}
}
