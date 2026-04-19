package pipeline

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/scanner"
)

// WithDebounce wraps a pipeline stage so that rapid successive diffs are
// coalesced: the stage is only invoked once the diff stream has been quiet
// for at least `wait` duration.
func WithDebounce(wait time.Duration) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		var last scanner.Diff
		d := debounce.New(wait, func() {
			last = diff
		})
		d.Call()
		d.Flush()
		if last.Opened == nil && last.Closed == nil {
			return diff, nil
		}
		return last, nil
	}
}
