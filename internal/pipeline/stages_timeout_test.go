package pipeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func slowStage(delay time.Duration) pipeline.Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		select {
		case <-time.After(delay):
			return d, nil
		case <-ctx.Done():
			return d, ctx.Err()
		}
	}
}

func TestWithTimeoutFastStageCompletes(t *testing.T) {
	d := diff([]uint16{80}, nil)
	stage, err := pipeline.WithTimeout(100*time.Millisecond, slowStage(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}

func TestWithTimeoutSlowStageReturnsError(t *testing.T) {
	d := diff([]uint16{443}, nil)
	stage, err := pipeline.WithTimeout(20*time.Millisecond, slowStage(200*time.Millisecond))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, stageErr := stage(context.Background(), d)
	if stageErr == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(stageErr, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", stageErr)
	}
}

func TestWithTimeoutZeroDurationReturnsError(t *testing.T) {
	_, err := pipeline.WithTimeout(0, slowStage(0))
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestWithTimeoutEmptyDiffPassesThrough(t *testing.T) {
	stage, err := pipeline.WithTimeout(50*time.Millisecond, func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		return d, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestWithTimeoutIntegratesWithPipeline(t *testing.T) {
	d := diff([]uint16{22}, nil)
	inner := func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		return d, nil
	}
	timeoutStage, err := pipeline.WithTimeout(100*time.Millisecond, inner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := pipeline.New(timeoutStage)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}
