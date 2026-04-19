package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/scanner"
)

// WithNotifier returns a Stage that sends notifications for port changes.
func WithNotifier(n *notifier.Notifier) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		if err := n.Notify(diff); err != nil {
			return diff, err
		}
		return diff, nil
	}
}
