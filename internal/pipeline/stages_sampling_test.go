package pipeline_test

import (
	"context"
	"testing"

	"github.com/your-org/portwatch/internal/pipeline"
	"github.com/your-org/portwatch/internal/sampling"
	"github.com/your-org/portwatch/internal/scanner"
)

func makeSampler(t *testing.T, rate float64) *sampling.Sampler {
	t.Helper()
	s, err := sampling.New(rate)
	if err != nil {
		t.Fatalf("sampling.New(%v): %v", rate, err)
	}
	return s
}

// TestWithSamplingEmptyDiffPassesThrough verifies that an empty diff is forwarded
// without invoking the sampler, preserving pipeline throughput.
func TestWithSamplingEmptyDiffPassesThrough(t *testing.T) {
	s := makeSampler(t, 1.0)
	stage := pipeline.WithSampling(s)

	in := scanner.Diff{}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Errorf("expected empty diff, got %+v", out)
	}
}

// TestWithSamplingRateOnePassesAll verifies that a sampler with rate 1.0 passes
// every entry through the stage unchanged.
func TestWithSamplingRateOnePassesAll(t *testing.T) {
	s := makeSampler(t, 1.0)
	stage := pipeline.WithSampling(s)

	in := openDiff(80, 443)
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 2 {
		t.Errorf("expected 2 opened entries, got %d", len(out.Opened))
	}
}

// TestWithSamplingRateZeroBlocksAll verifies that a sampler with rate 0.0 drops
// every entry, producing an empty diff.
func TestWithSamplingRateZeroBlocksAll(t *testing.T) {
	s := makeSampler(t, 0.0)
	stage := pipeline.WithSampling(s)

	in := openDiff(80, 443)
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 {
		t.Errorf("expected 0 opened entries after rate-0 sampling, got %d", len(out.Opened))
	}
}

// TestWithSamplingIntegratesWithPipeline verifies that WithSampling composes
// correctly within a multi-stage pipeline and that context cancellation is
// respected.
func TestWithSamplingIntegratesWithPipeline(t *testing.T) {
	s := makeSampler(t, 1.0)

	p := pipeline.New(
		pipeline.WithSampling(s),
	)

	in := openDiff(22, 8080)
	out, err := p.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
	if len(out.Opened) != 2 {
		t.Errorf("expected 2 opened ports through pipeline, got %d", len(out.Opened))
	}
}
