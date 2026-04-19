package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/health"
	"github.com/user/portwatch/internal/scanner"
)

// WithHealth returns a Stage that records each scan event into the health tracker.
func WithHealth(h *health.Health) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}
		h.RecordScan()
		return d, nil
	}
}
