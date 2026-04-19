package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

// WithHistory returns a Stage that records each diff entry into the history log.
func WithHistory(h *history.History) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened)+len(diff.Closed) == 0 {
			return diff, nil
		}
		if err := h.Record(diff); err != nil {
			return diff, err
		}
		return diff, nil
	}
}
