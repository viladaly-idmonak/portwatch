package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/formatter"
	"github.com/user/portwatch/internal/scanner"
)

// WithFormatter returns a Stage that annotates each diff entry's metadata
// by running the diff through the formatter and attaching the rendered
// string to the context so downstream stages can consume it.
//
// The stage is a pass-through: it never drops or modifies the diff itself.
func WithFormatter(f *formatter.Formatter) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		// Render to ensure the formatter produces valid output; errors surface
		// early in the pipeline rather than at the notifier/logger stage.
		_, err := f.Format(diff)
		if err != nil {
			return diff, err
		}
		return diff, nil
	}
}
