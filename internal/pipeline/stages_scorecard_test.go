package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scorecard"
)

func makeScorecard() *scorecard.Scorecard {
	sc := scorecard.New("low", nil)
	sc.AddRule(8080, "tcp", "high")
	sc.AddRule(443, "tcp", "medium")
	return sc
}

func TestWithScorecardEmptyDiffPassesThrough(t *testing.T) {
	sc := makeScorecard()
	d := scanner.Diff{}
	out, err := pipeline.New(context.Background(), d, pipeline.WithScorecard(sc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithScorecardAnnotatesOpenedEntry(t *testing.T) {
	sc := makeScorecard()
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
	out, err := pipeline.New(context.Background(), d, pipeline.WithScorecard(sc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	got := out.Opened[0].Meta["score"]
	if got != "high" {
		t.Errorf("expected score=high, got %q", got)
	}
}

func TestWithScorecardAnnotatesClosedEntry(t *testing.T) {
	sc := makeScorecard()
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 443, Proto: "tcp"}},
	}
	out, err := pipeline.New(context.Background(), d, pipeline.WithScorecard(sc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %d", len(out.Closed))
	}
	got := out.Closed[0].Meta["score"]
	if got != "medium" {
		t.Errorf("expected score=medium, got %q", got)
	}
}

func TestWithScorecardNilIsNoop(t *testing.T) {
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 9999, Proto: "tcp"}},
	}
	out, err := pipeline.New(context.Background(), d, pipeline.WithScorecard(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if _, ok := out.Opened[0].Meta["score"]; ok {
		t.Error("expected no score annotation on nil scorecard")
	}
}

func TestWithScorecardIntegratesWithPipeline(t *testing.T) {
	sc := makeScorecard()
	d := scanner.Diff{
		Opened: []scanner.Entry{
			{Port: 8080, Proto: "tcp"},
			{Port: 22, Proto: "tcp"},
		},
	}
	out, err := pipeline.New(context.Background(), d, pipeline.WithScorecard(sc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 2 {
		t.Fatalf("expected 2 opened entries, got %d", len(out.Opened))
	}
	scores := map[uint16]string{}
	for _, e := range out.Opened {
		scores[e.Port] = e.Meta["score"]
	}
	if scores[8080] != "high" {
		t.Errorf("port 8080: expected high, got %q", scores[8080])
	}
	if scores[22] != "low" {
		t.Errorf("port 22: expected low (default), got %q", scores[22])
	}
}
