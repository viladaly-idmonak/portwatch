package pipeline_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/suppress"
)

func makeSuppress(t *testing.T, window time.Duration) *suppress.Suppressor {
	t.Helper()
	s, err := suppress.New(window)
	if err != nil {
		t.Fatalf("suppress.New: %v", err)
	}
	return s
}

func TestWithSuppressEmptyDiffPassesThrough(t *testing.T) {
	s := makeSuppress(t, 5*time.Second)
	stage := pipeline.WithSuppress(s)
	diff := scanner.Diff{}
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithSuppressFirstEventPasses(t *testing.T) {
	s := makeSuppress(t, 5*time.Second)
	stage := pipeline.WithSuppress(s)
	diff := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}

func TestWithSuppressDuplicateWithinWindowBlocked(t *testing.T) {
	s := makeSuppress(t, 5*time.Second)
	stage := pipeline.WithSuppress(s)
	diff := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
	// first pass
	_, _ = stage(t.Context(), diff)
	// second pass within window
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 {
		t.Fatalf("expected 0 opened entries (suppressed), got %d", len(out.Opened))
	}
}

func TestWithSuppressIntegratesWithPipeline(t *testing.T) {
	s := makeSuppress(t, 5*time.Second)
	p := pipeline.New(
		pipeline.WithSuppress(s),
	)
	diff := scanner.Diff{
		Opened: []scanner.Entry{{Port: 9090, Proto: "tcp"}},
	}
	out, err := p.Run(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}
