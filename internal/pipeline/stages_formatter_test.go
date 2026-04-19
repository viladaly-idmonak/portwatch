package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/formatter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeFormatter(t *testing.T, format string) *formatter.Formatter {
	t.Helper()
	f, err := formatter.New(format)
	if err != nil {
		t.Fatalf("formatter.New: %v", err)
	}
	return f
}

func TestWithFormatterEmptyDiffPassesThrough(t *testing.T) {
	f := makeFormatter(t, "text")
	stage := pipeline.WithFormatter(f)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithFormatterTextFormatPassesDiff(t *testing.T) {
	f := makeFormatter(t, "text")
	stage := pipeline.WithFormatter(f)
	d := openDiff(9200)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(out.Opened))
	}
}

func TestWithFormatterJSONFormatPassesDiff(t *testing.T) {
	f := makeFormatter(t, "json")
	stage := pipeline.WithFormatter(f)
	d := openDiff(8080)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(out.Opened))
	}
}

func TestWithFormatterIntegratesWithPipeline(t *testing.T) {
	f := makeFormatter(t, "text")
	p := pipeline.New(pipeline.WithFormatter(f))
	d := openDiff(3000)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(out.Opened))
	}
}
