package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeCooldown(t *testing.T, window time.Duration) *cooldown.Cooldown {
	t.Helper()
	cd, err := cooldown.New(window)
	if err != nil {
		t.Fatalf("cooldown.New: %v", err)
	}
	return cd
}

func TestWithCooldownEmptyDiffPassesThrough(t *testing.T) {
	cd := makeCooldown(t, 100*time.Millisecond)
	stage := pipeline.WithCooldown(cd)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithCooldownFirstEventPasses(t *testing.T) {
	cd := makeCooldown(t, 100*time.Millisecond)
	stage := pipeline.WithCooldown(cd)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 80, Protocol: "tcp"}}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}

func TestWithCooldownDuplicateWithinWindowBlocked(t *testing.T) {
	cd := makeCooldown(t, 500*time.Millisecond)
	stage := pipeline.WithCooldown(cd)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 443, Protocol: "tcp"}}}

	// First call should pass.
	out, _ := stage(context.Background(), d)
	if len(out.Opened) != 1 {
		t.Fatalf("first call: expected 1, got %d", len(out.Opened))
	}

	// Second call within window should be blocked.
	out, _ = stage(context.Background(), d)
	if len(out.Opened) != 0 {
		t.Fatalf("second call: expected 0, got %d", len(out.Opened))
	}
}

func TestWithCooldownDistinctStatesAreIndependent(t *testing.T) {
	cd := makeCooldown(t, 500*time.Millisecond)
	stage := pipeline.WithCooldown(cd)

	opened := scanner.Diff{Opened: []scanner.Entry{{Port: 22, Protocol: "tcp"}}}
	closed := scanner.Diff{Closed: []scanner.Entry{{Port: 22, Protocol: "tcp"}}}

	out, _ := stage(context.Background(), opened)
	if len(out.Opened) != 1 {
		t.Fatalf("opened: expected 1, got %d", len(out.Opened))
	}

	out, _ = stage(context.Background(), closed)
	if len(out.Closed) != 1 {
		t.Fatalf("closed: expected 1, got %d", len(out.Closed))
	}
}

func TestWithCooldownNilIsNoop(t *testing.T) {
	stage := pipeline.WithCooldown(nil)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 8080, Protocol: "tcp"}}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("nil cooldown: expected 1, got %d", len(out.Opened))
	}
}

func TestWithCooldownIntegratesWithPipeline(t *testing.T) {
	cd := makeCooldown(t, 500*time.Millisecond)
	p := pipeline.New(pipeline.WithCooldown(cd))
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 9000, Protocol: "udp"}}}

	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("pipeline: expected 1, got %d", len(out.Opened))
	}
}
