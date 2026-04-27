package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeLimiter(t *testing.T, max int) *pipeline.Limiter {
	t.Helper()
	l, err := pipeline.NewLimiter(max)
	if err != nil {
		t.Fatalf("NewLimiter: %v", err)
	}
	return l
}

func TestWithLimiterEmptyDiffPassesThrough(t *testing.T) {
	l := makeLimiter(t, 5)
	stage := pipeline.WithLimiter(l)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithLimiterAllowsUnderLimit(t *testing.T) {
	l := makeLimiter(t, 10)
	stage := pipeline.WithLimiter(l)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80}, {Port: 443}},
		Closed: []scanner.Entry{{Port: 8080}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 2 || len(out.Closed) != 1 {
		t.Fatalf("expected 2 opened, 1 closed; got %d opened, %d closed", len(out.Opened), len(out.Closed))
	}
}

func TestWithLimiterTruncatesOpened(t *testing.T) {
	l := makeLimiter(t, 2)
	stage := pipeline.WithLimiter(l)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80}, {Port: 443}, {Port: 8080}},
	}
	out, _ := stage(context.Background(), d)
	if len(out.Opened) != 2 {
		t.Fatalf("expected 2 opened entries, got %d", len(out.Opened))
	}
}

func TestWithLimiterRemainingCapacityGivenToClosed(t *testing.T) {
	l := makeLimiter(t, 3)
	stage := pipeline.WithLimiter(l)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80}, {Port: 443}},
		Closed: []scanner.Entry{{Port: 8080}, {Port: 9090}},
	}
	out, _ := stage(context.Background(), d)
	if len(out.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(out.Opened))
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed (remaining capacity), got %d", len(out.Closed))
	}
}

func TestNewLimiterInvalidMaxReturnsError(t *testing.T) {
	_, err := pipeline.NewLimiter(0)
	if err == nil {
		t.Fatal("expected error for maxEntries=0")
	}
}

func TestWithLimiterIntegratesWithPipeline(t *testing.T) {
	l := makeLimiter(t, 1)
	p := pipeline.New(
		pipeline.WithLimiter(l),
	)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 22}, {Port: 80}},
	}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry after limit, got %d", len(out.Opened))
	}
}
