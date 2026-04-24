package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/circuitbreaker"
	"github.com/user/portwatch/internal/scanner"
)

// WithCircuitBreaker wraps a downstream stage with circuit-breaker protection.
// If the circuit is open, the stage is skipped and a warning is emitted to errFn.
// On success the circuit is closed; on error the failure is recorded.
//
// Pass nil for cb to disable protection (stage always runs).
func WithCircuitBreaker(cb *circuitbreaker.CircuitBreaker, stage Stage, errFn func(error)) Stage {
	if cb == nil {
		return stage
	}
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if err := cb.Allow(); err != nil {
			if errFn != nil {
				errFn(fmt.Errorf("pipeline: circuit open, skipping stage: %w", err))
			}
			return d, nil // pass diff through without running stage
		}
		out, err := stage(ctx, d)
		if err != nil {
			cb.RecordFailure()
			return out, err
		}
		cb.RecordSuccess()
		return out, nil
	}
}
