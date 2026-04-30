package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/decay"
	"github.com/yourorg/portwatch/internal/pipeline"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makeDecayer(t *testing.T) *decay.Decayer {
	t.Helper()
	d, err := decay.New(time.Hour)
	if err != nil {
		t.Fatalf("decay.New: %v", err)
	}
	return d
}

func TestWithDecayEmptyDiffPassesThrough(t *testing.T) {
	d := makeDecayer(t)
	stage := pipeline.WithDecay(d)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestWithDecayNilIsNoop(t *testing.T) {
	stage := pipeline.WithDecay(nil)
	in := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Protocol: "tcp"}},
	}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Opened[0].Meta != nil {
		t.Fatal("nil decayer should not annotate entries")
	}
}

func TestWithDecayAnnotatesOpenedEntry(t *testing.T) {
	d := makeDecayer(t)
	stage := pipeline.WithDecay(d)
	in := scanner.Diff{
		Opened: []scanner.Entry{{Port: 443, Protocol: "tcp"}},
	}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if out.Opened[0].Meta["decay_score"] == "" {
		t.Error("expected decay_score meta key to be set")
	}
}

func TestWithDecayAnnotatesClosedEntry(t *testing.T) {
	d := makeDecayer(t)
	stage := pipeline.WithDecay(d)
	// Pre-populate a score.
	d.Add("tcp:22", 3.0)
	in := scanner.Diff{
		Closed: []scanner.Entry{{Port: 22, Protocol: "tcp"}},
	}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Closed[0].Meta["decay_score"] == "" {
		t.Error("expected decay_score meta key to be set on closed entry")
	}
}

func TestWithDecayIntegratesWithPipeline(t *testing.T) {
	d := makeDecayer(t)
	p := pipeline.New(
		pipeline.WithDecay(d),
	)
	in := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Protocol: "tcp"}},
	}
	out, err := p.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if out.Opened[0].Meta["decay_score"] == "" {
		t.Error("expected decay_score to be present after pipeline run")
	}
}
