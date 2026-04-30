package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/trend"
)

// WithTrend returns a Stage that records each diff into the Tracker and
// annotates every entry with the current trend direction ("rising",
// "falling", or "neutral"). A nil tracker is a no-op.
func WithTrend(t *trend.Tracker) Stage {
	if t == nil {
		return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}

	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		t.Record(len(d.Opened), len(d.Closed))
		direction := t.Direction()

		annotate := func(entries []scanner.Entry) []scanner.Entry {
			out := make([]scanner.Entry, len(entries))
			for i, e := range entries {
				if e.Meta == nil {
					e.Meta = make(map[string]string)
				}
				e.Meta["trend"] = direction
				out[i] = e
			}
			return out
		}

		return scanner.Diff{
			Opened: annotate(d.Opened),
			Closed: annotate(d.Closed),
		}, nil
	}
}
