package pipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/classifier"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeClassifier(t *testing.T) *classifier.Classifier {
	t.Helper()
	c, err := classifier.New("unknown", []classifier.Rule{
		{Port: 80, Protocol: "tcp", Class: "web"},
		{Port: 443, Protocol: "tcp", Class: "web"},
		{Port: 22, Protocol: "tcp", Class: "admin"},
	})
	if err != nil {
		t.Fatalf("makeClassifier: %v", err)
	}
	return c
}

func TestWithClassifierEmptyDiffPassesThrough(t *testing.T) {
	c := makeClassifier(t)
	p := pipeline.New(pipeline.WithClassifier(c))
	out, err := p.Run(t.Context(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithClassifierAnnotatesOpenedEntry(t *testing.T) {
	c := makeClassifier(t)
	p := pipeline.New(pipeline.WithClassifier(c))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Protocol: "tcp"}},
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	got := out.Opened[0].Meta["class"]
	if got != "web" {
		t.Errorf("expected class=web, got %q", got)
	}
}

func TestWithClassifierAnnotatesClosedEntry(t *testing.T) {
	c := makeClassifier(t)
	p := pipeline.New(pipeline.WithClassifier(c))
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 22, Protocol: "tcp"}},
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %d", len(out.Closed))
	}
	got := out.Closed[0].Meta["class"]
	if got != "admin" {
		t.Errorf("expected class=admin, got %q", got)
	}
}

func TestWithClassifierUnknownPortUsesDefault(t *testing.T) {
	c := makeClassifier(t)
	p := pipeline.New(pipeline.WithClassifier(c))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 9999, Protocol: "tcp"}},
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.Opened[0].Meta["class"]
	if got != "unknown" {
		t.Errorf("expected class=unknown, got %q", got)
	}
}

func TestWithClassifierNilIsNoop(t *testing.T) {
	p := pipeline.New(pipeline.WithClassifier(nil))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Protocol: "tcp"}},
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Opened[0].Meta["class"]; ok {
		t.Error("expected no class annotation when classifier is nil")
	}
}

func TestWithClassifierIntegratesWithPipeline(t *testing.T) {
	c := makeClassifier(t)
	p := pipeline.New(
		pipeline.WithClassifier(c),
		pipeline.WithClassifier(c), // idempotent: second pass should not overwrite
	)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 443, Protocol: "tcp"}},
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Opened[0].Meta["class"]; got != "web" {
		t.Errorf("expected class=web, got %q", got)
	}
}
