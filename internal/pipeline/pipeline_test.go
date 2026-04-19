package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func diff(opened, closed []uint16) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func TestEmptyPipelinePassesThrough(t *testing.T) {
	p := pipeline.New()
	in := diff([]uint16{80}, nil)
	out, err := p.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 || out.Opened[0] != 80 {
		t.Fatalf("expected diff unchanged, got %+v", out)
	}
}

func TestStagesAppliedInOrder(t *testing.T) {
	order := []int{}
	makeStage := func(n int) pipeline.Stage {
		return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
			order = append(order, n)
			return d, nil
		}
	}
	p := pipeline.New(makeStage(1), makeStage(2), makeStage(3))
	p.Run(context.Background(), diff(nil, nil))
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestStageErrorStopsPipeline(t *testing.T) {
	ran := false
	errStage := func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		return d, errors.New("boom")
	}
	afterStage := func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		ran = true
		return d, nil
	}
	p := pipeline.New(errStage, afterStage)
	_, err := p.Run(context.Background(), diff(nil, nil))
	if err == nil {
		t.Fatal("expected error")
	}
	if ran {
		t.Fatal("stage after error should not run")
	}
}

func TestContextCancellationStopsPipeline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ran := false
	stage := func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		ran = true
		return d, nil
	}
	p := pipeline.New(stage)
	_, err := p.Run(ctx, diff(nil, nil))
	if err == nil {
		t.Fatal("expected context error")
	}
	if ran {
		t.Fatal("stage should not run after cancel")
	}
}

func TestLenReturnsStageCount(t *testing.T) {
	noop := func(_ context.Context, d scanner.Diff) (scanner.Diff, error) { return d, nil }
	p := pipeline.New(noop, noop)
	if p.Len() != 2 {
		t.Fatalf("expected 2, got %d", p.Len())
	}
}
