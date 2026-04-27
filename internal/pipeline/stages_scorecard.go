package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scorecard"
)

// WithScorecard returns a Stage that records a risk score for every opened
// port entry. The score is written into the entry's Meta map under the key
// "score". Closed entries and empty diffs are passed through unchanged.
func WithScorecard(sc *scorecard.Scorecard) Stage {
	if sc == nil {
		return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}
	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		out := scanner.Diff{
			Opened: make([]scanner.Entry, len(d.Opened)),
			Closed: d.Closed,
		}
		for i, e := range d.Opened {
			if e.Meta == nil {
				e.Meta = make(map[string]string)
			}
			e.Meta["score"] = sc.Score(e)
			out.Opened[i] = e
		}
		return out, nil
	}
}
