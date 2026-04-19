package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
)

// WithMetrics returns a Stage that records opened/closed port counts.
func WithMetrics(m *metrics.Metrics) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened)+len(diff.Closed) == 0 {
			return diff, nil
		}
		m.RecordOpened(len(diff.Opened))
		m.RecordClosed(len(diff.Closed))
		return diff, nil
	}
}
