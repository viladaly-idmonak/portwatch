package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/scanner"
)

// WithAlerter returns a Stage that dispatches webhook alerts for port changes.
func WithAlerter(a *alerter.Alerter) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		if err := a.Notify(ctx, diff); err != nil {
			return diff, err
		}
		return diff, nil
	}
}
