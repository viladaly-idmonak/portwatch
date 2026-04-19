package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summarizer"
)

func TestWithSummarizerEmptyDiffSkips(t *testing.T) {
	s := summarizer.New()
	stage := pipeline.WithSummarizer(s)
	d, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatal(err)
	}
	sum := s.Flush()
	if len(sum.Opened) != 0 || len(sum.Closed) != 0 {
		t.Fatalf("expected no records for empty diff, got %+v", sum)
	}
	_ = d
}

func TestWithSummarizerRecordsOpened(t *testing.T) {
	s := summarizer.New()
	stage := pipeline.WithSummarizer(s)
	input := scanner.Diff{Opened: []scanner.Port{{Number: 80, Protocol: "tcp"}}}
	_, err := stage(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	sum := s.Flush()
	if len(sum.Opened) != 1 {
		t.Fatalf("expected 1 opened in summary, got %d", len(sum.Opened))
	}
}

func TestWithSummarizerIntegratesWithPipeline(t *testing.T) {
	s := summarizer.New()
	p := pipeline.New(pipeline.WithSummarizer(s))
	input := scanner.Diff{
		Opened: []scanner.Port{{Number: 443, Protocol: "tcp"}},
		Closed: []scanner.Port{{Number: 22, Protocol: "tcp"}},
	}
	out, err := p.Run(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Opened) != 1 || len(out.Closed) != 1 {
		t.Fatalf("diff should pass through unchanged: %+v", out)
	}
	sum := s.Flush()
	if len(sum.Opened) != 1 || len(sum.Closed) != 1 {
		t.Fatalf("summary should reflect diff: %+v", sum)
	}
}
