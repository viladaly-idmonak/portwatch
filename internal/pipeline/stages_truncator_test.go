package pipeline_test

import (
	"testing"

	"github.com/example/portwatch/internal/pipeline"
	"github.com/example/portwatch/internal/scanner"
)

func makeTruncator(max int) *pipeline.Truncator {
	return pipeline.NewTruncator(max)
}

func openEntries(n int) []scanner.Entry {
	entries := make([]scanner.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = scanner.Entry{Port: uint16(8000 + i), Proto: "tcp"}
	}
	return entries
}

func TestWithTruncatorEmptyDiffPassesThrough(t *testing.T) {
	p := pipeline.New()
	p.Use(pipeline.WithTruncator(makeTruncator(5)))

	out, err := p.Run(t.Context(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithTruncatorAllowsUnderLimit(t *testing.T) {
	p := pipeline.New()
	p.Use(pipeline.WithTruncator(makeTruncator(5)))

	d := scanner.Diff{Opened: openEntries(3)}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 3 {
		t.Fatalf("expected 3 opened, got %d", len(out.Opened))
	}
}

func TestWithTruncatorTruncatesOpened(t *testing.T) {
	p := pipeline.New()
	p.Use(pipeline.WithTruncator(makeTruncator(4)))

	d := scanner.Diff{Opened: openEntries(10)}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 4 {
		t.Fatalf("expected 4 opened after truncation, got %d", len(out.Opened))
	}
}

func TestWithTruncatorTruncatesClosed(t *testing.T) {
	p := pipeline.New()
	p.Use(pipeline.WithTruncator(makeTruncator(2)))

	entries := openEntries(6)
	d := scanner.Diff{Closed: entries}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 2 {
		t.Fatalf("expected 2 closed after truncation, got %d", len(out.Closed))
	}
}

func TestWithTruncatorIntegratesWithPipeline(t *testing.T) {
	p := pipeline.New()
	p.Use(pipeline.WithTruncator(makeTruncator(3)))

	d := scanner.Diff{
		Opened: openEntries(7),
		Closed: openEntries(5),
	}
	out, err := p.Run(t.Context(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 3 {
		t.Errorf("expected 3 opened, got %d", len(out.Opened))
	}
	if len(out.Closed) != 3 {
		t.Errorf("expected 3 closed, got %d", len(out.Closed))
	}
}
