package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/labeler"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeLabeler(t *testing.T) *labeler.Labeler {
	t.Helper()
	rules := []labeler.Rule{
		{Port: 80, Protocol: "tcp", Label: "http"},
		{Port: 443, Protocol: "tcp", Label: "https"},
		{Port: 22, Protocol: "tcp", Label: "ssh"},
	}
	lb, err := labeler.New(rules)
	if err != nil {
		t.Fatalf("makeLabeler: %v", err)
	}
	return lb
}

func TestWithLabelerEmptyDiffPassesThrough(t *testing.T) {
	stage := pipeline.WithLabeler(makeLabeler(t))
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithLabelerAnnotatesKnownPort(t *testing.T) {
	stage := pipeline.WithLabeler(makeLabeler(t))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if got := out.Opened[0].Meta["label"]; got != "http" {
		t.Errorf("expected label=http, got %q", got)
	}
}

func TestWithLabelerUnknownPortNoLabel(t *testing.T) {
	stage := pipeline.WithLabeler(makeLabeler(t))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 9999, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Opened[0].Meta["label"]; ok {
		t.Error("expected no label for unknown port")
	}
}

func TestWithLabelerNilLabelerIsNoop(t *testing.T) {
	stage := pipeline.WithLabeler(nil)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Opened[0].Meta != nil {
		t.Error("expected nil meta when labeler is nil")
	}
}

func TestWithLabelerIntegratesWithPipeline(t *testing.T) {
	p := pipeline.New(
		pipeline.WithLabeler(makeLabeler(t)),
	)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 443, Proto: "tcp"}},
		Closed: []scanner.Entry{{Port: 22, Proto: "tcp"}},
	}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Opened[0].Meta["label"]; got != "https" {
		t.Errorf("opened: expected https, got %q", got)
	}
	if got := out.Closed[0].Meta["label"]; got != "ssh" {
		t.Errorf("closed: expected ssh, got %q", got)
	}
}
