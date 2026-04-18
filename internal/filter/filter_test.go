package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func TestApplyNoRulesPassesAll(t *testing.T) {
	f := filter.NewFromSlices(nil, nil)
	d := scanner.Diff{Opened: []uint16{80, 443}, Closed: []uint16{8080}}
	out := f.Apply(d)
	assertPorts(t, out.Opened, []uint16{80, 443})
	assertPorts(t, out.Closed, []uint16{8080})
}

func TestApplyIncludeFilter(t *testing.T) {
	f := filter.NewFromSlices([]uint16{443}, nil)
	d := scanner.Diff{Opened: []uint16{80, 443, 8080}, Closed: []uint16{22, 443}}
	out := f.Apply(d)
	assertPorts(t, out.Opened, []uint16{443})
	assertPorts(t, out.Closed, []uint16{443})
}

func TestApplyExcludeFilter(t *testing.T) {
	f := filter.NewFromSlices(nil, []uint16{22, 8080})
	d := scanner.Diff{Opened: []uint16{22, 80, 8080}, Closed: []uint16{8080, 443}}
	out := f.Apply(d)
	assertPorts(t, out.Opened, []uint16{80})
	assertPorts(t, out.Closed, []uint16{443})
}

func TestApplyExcludeTakesPrecedenceOverInclude(t *testing.T) {
	// port 80 is in both include and exclude — exclude wins
	f := filter.NewFromSlices([]uint16{80, 443}, []uint16{80})
	d := scanner.Diff{Opened: []uint16{80, 443}, Closed: []uint16{}}
	out := f.Apply(d)
	assertPorts(t, out.Opened, []uint16{443})
}

func TestApplyEmptyDiff(t *testing.T) {
	f := filter.NewFromSlices([]uint16{80}, []uint16{22})
	d := scanner.Diff{}
	out := f.Apply(d)
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Errorf("expected empty diff, got %+v", out)
	}
}

func assertPorts(t *testing.T, got, want []uint16) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("port count mismatch: got %v, want %v", got, want)
	}
	set := make(map[uint16]struct{}, len(want))
	for _, p := range want {
		set[p] = struct{}{}
	}
	for _, p := range got {
		if _, ok := set[p]; !ok {
			t.Errorf("unexpected port %d in result", p)
		}
	}
}
