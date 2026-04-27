package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Limiter caps the total number of entries (opened + closed) that pass through
// a pipeline stage. Entries beyond the cap are silently dropped.
type Limiter struct {
	maxEntries int
}

// NewLimiter creates a Limiter that allows at most maxEntries total entries per
// diff. maxEntries must be greater than zero.
func NewLimiter(maxEntries int) (*Limiter, error) {
	if maxEntries <= 0 {
		return nil, fmt.Errorf("pipeline/limiter: maxEntries must be > 0, got %d", maxEntries)
	}
	return &Limiter{maxEntries: maxEntries}, nil
}

// Apply returns a copy of diff with at most maxEntries total entries.
// Opened entries are prioritised; remaining capacity is given to closed entries.
func (l *Limiter) Apply(diff scanner.Diff) scanner.Diff {
	opened := diff.Opened
	if len(opened) > l.maxEntries {
		opened = opened[:l.maxEntries]
	}
	remaining := l.maxEntries - len(opened)
	closed := diff.Closed
	if len(closed) > remaining {
		closed = closed[:remaining]
	}
	return scanner.Diff{Opened: opened, Closed: closed}
}

// WithLimiter returns a Stage that caps the number of entries in a diff.
func WithLimiter(l *Limiter) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		select {
		case <-ctx.Done():
			return scanner.Diff{}, ctx.Err()
		default:
		}
		return l.Apply(diff), nil
	}
}
