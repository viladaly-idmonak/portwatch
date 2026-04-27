package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Timeouter wraps a pipeline stage with a per-invocation deadline.
type Timeouter struct {
	timeout time.Duration
}

// NewTimeouter creates a Timeouter with the given deadline.
func NewTimeouter(timeout time.Duration) (*Timeouter, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive, got %v", timeout)
	}
	return &Timeouter{timeout: timeout}, nil
}

// Wrap returns a Stage that runs the next stage with a bounded context.
// If the downstream stage exceeds the deadline the diff is passed through
// unchanged and the context error is returned so the pipeline can decide
// whether to halt.
func (t *Timeouter) Wrap(next Stage) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		ctx, cancel := context.WithTimeout(ctx, t.timeout)
		defer cancel()

		type result struct {
			d   scanner.Diff
			err error
		}

		ch := make(chan result, 1)
		go func() {
			out, err := next(ctx, d)
			ch <- result{out, err}
		}()

		select {
		case r := <-ch:
			return r.d, r.err
		case <-ctx.Done():
			return d, fmt.Errorf("stage timeout after %v: %w", t.timeout, ctx.Err())
		}
	}
}

// WithTimeout wraps the provided stage so that it must complete within
// timeout. On expiry the original diff is forwarded and an error is returned.
func WithTimeout(timeout time.Duration, next Stage) (Stage, error) {
	t, err := NewTimeouter(timeout)
	if err != nil {
		return nil, err
	}
	return t.Wrap(next), nil
}
