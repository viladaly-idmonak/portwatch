// Package pipeline provides composable processing stages for port diffs.
package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

// WithBaseline returns a Stage that annotates the diff context and optionally
// filters out ports that are within the established baseline.
// When enforceBaseline is true, only unexpected or missing ports pass through;
// when false, the stage is informational only and the diff is forwarded as-is.
func WithBaseline(b *baseline.Baseline, enforceBaseline bool) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}

		if !enforceBaseline {
			return d, nil
		}

		current := make([]scanner.Entry, 0, len(d.Opened))
		for _, e := range d.Opened {
			current = append(current, e)
		}

		unexpected, _ := b.Deviations(current)
		unexpectedSet := make(map[uint16]struct{}, len(unexpected))
		for _, p := range unexpected {
			unexpectedSet[p] = struct{}{}
		}

		var filteredOpened []scanner.Entry
		for _, e := range d.Opened {
			if _, ok := unexpectedSet[e.Port]; ok {
				filteredOpened = append(filteredOpened, e)
			}
		}

		var filteredClosed []scanner.Entry
		for _, e := range d.Closed {
			if b.Contains(e.Port) {
				filteredClosed = append(filteredClosed, e)
			}
		}

		if len(filteredOpened)+len(filteredClosed) == 0 {
			return scanner.Diff{}, nil
		}

		_ = fmt.Sprintf // keep import
		return scanner.Diff{Opened: filteredOpened, Closed: filteredClosed}, nil
	}
}
