package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/circuitbreaker"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeCircuitBreaker(t *testing.T, maxFailures int) *circuitbreaker.CircuitBreaker {
	t.Helper()
	cfg := circuitbreaker.DefaultConfig()
	cfg.MaxFailures = maxFailures
	cb, err := circuitbreaker.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig: %v", err)
	}
	return cb
}

func TestWithCircuitBreakerAllowsWhenClosed(t *testing.T) {
	cb := makeCircuitBreaker(t, 3)
	var called bool
	stage := pipeline.WithCircuitBreaker(cb, func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		called = true
		return d, nil
	})

	d := openDiff(8080, "tcp")
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner fn to be called")
	}
	if len(out.Opened) != 1 {
		t.Errorf("expected 1 opened, got %d", len(out.Opened))
	}
}

func TestWithCircuitBreakerOpensAfterFailures(t *testing.T) {
	cb := makeCircuitBreaker(t, 2)
	failFn := func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		return scanner.Diff{}, errors.New("downstream failure")
	}
	stage := pipeline.WithCircuitBreaker(cb, failFn)

	d := openDiff(9090, "tcp")
	for i := 0; i < 2; i++ {
		_, _ = stage(context.Background(), d)
	}

	// Circuit should now be open; next call should fail fast without calling fn
	var innerCalled bool
	guardedStage := pipeline.WithCircuitBreaker(cb, func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		innerCalled = true
		return d, nil
	})
	_, err := guardedStage(context.Background(), d)
	if err == nil {
		t.Error("expected error from open circuit")
	}
	if innerCalled {
		t.Error("expected inner fn NOT to be called when circuit is open")
	}
}

func TestWithCircuitBreakerEmptyDiffPassesThrough(t *testing.T) {
	cb := makeCircuitBreaker(t, 3)
	var called bool
	stage := pipeline.WithCircuitBreaker(cb, func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		called = true
		return d, nil
	})

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected inner fn NOT to be called for empty diff")
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Error("expected empty diff to pass through unchanged")
	}
}

func TestWithCircuitBreakerIntegratesWithPipeline(t *testing.T) {
	cb := makeCircuitBreaker(t, 5)
	var seen []scanner.Diff

	p := pipeline.New(
		pipeline.WithCircuitBreaker(cb, func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
			seen = append(seen, d)
			return d, nil
		}),
	)

	d := openDiff(443, "tcp")
	if _, err := p.Run(context.Background(), d); err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
	if len(seen) != 1 {
		t.Errorf("expected 1 diff seen, got %d", len(seen))
	}
}
