package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

func makeRateLimit(window time.Duration) *ratelimit.RateLimit {
	return ratelimit.New(window)
}

func TestWithRateLimitPerPortAllowsFirst(t *testing.T) {
	rl := makeRateLimit(5 * time.Second)
	stage := pipeline.WithRateLimitPerPort(rl)

	d := openDiff(8080)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(out.Opened))
	}
}

func TestWithRateLimitPerPortBlocksDuplicate(t *testing.T) {
	rl := makeRateLimit(5 * time.Second)
	stage := pipeline.WithRateLimitPerPort(rl)

	d := openDiff(8080)
	_, _ = stage(context.Background(), d)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 {
		t.Fatalf("expected 0 opened after dedup, got %d", len(out.Opened))
	}
}

func TestWithRateLimitPerPortEmptyDiffPassesThrough(t *testing.T) {
	rl := makeRateLimit(5 * time.Second)
	stage := pipeline.WithRateLimitPerPort(rl)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestWithRateLimitPerPortIntegratesWithPipeline(t *testing.T) {
	rl := makeRateLimit(5 * time.Second)
	p := pipeline.New(pipeline.WithRateLimitPerPort(rl))

	d := openDiff(9090)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(out.Opened))
	}
}
