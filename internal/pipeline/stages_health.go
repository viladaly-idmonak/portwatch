package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/health"
	"github.com/user/portwatch/internal/scanner"
)

// WithHealth returns a Stage that records each scan event into the health tracker.
// Only scans that produced changes (opened or closed ports) are recorded.
func WithHealth(h *health.Health) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}
		h.RecordScan()
		return d, nil
	}
}

// WithHealthAlways returns a Stage that records every scan into the health
// tracker regardless of whether any ports changed. Use this when you want the
// health tracker to reflect scan activity even during quiet periods.
func WithHealthAlways(h *health.Health) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		h.RecordScan()
		return d, nil
	}
}
