package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
)

// WithLogger returns a Stage that logs every diff using the provided logger.
func WithLogger(l *logger.Logger) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened)+len(diff.Closed) == 0 {
			return diff, nil
		}
		l.Log(diff)
		return diff, nil
	}
}
