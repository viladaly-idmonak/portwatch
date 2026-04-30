package pipeline

import (
	"context"
	"fmt"

	"github.com/yourorg/portwatch/internal/decay"
	"github.com/yourorg/portwatch/internal/scanner"
)

const decayScoreMetaKey = "decay_score"

// WithDecay returns a Stage that annotates each entry in the diff with its
// current decayed score. Opened entries increment the score by 1.0; closed
// entries decrement it by 1.0 (floor 0). The score is stored in entry.Meta
// under the key "decay_score".
//
// A nil decayer is a no-op: the diff passes through unchanged.
func WithDecay(d *decay.Decayer) Stage {
	if d == nil {
		return func(_ context.Context, diff scanner.Diff) (scanner.Diff, error) {
			return diff, nil
		}
	}
	return func(_ context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		out := scanner.Diff{
			Opened: make([]scanner.Entry, len(diff.Opened)),
			Closed: make([]scanner.Entry, len(diff.Closed)),
		}
		for i, e := range diff.Opened {
			k := decayKey(e)
			score := d.Add(k, 1.0)
			out.Opened[i] = annotateDecayScore(e, score)
		}
		for i, e := range diff.Closed {
			k := decayKey(e)
			score := d.Add(k, -1.0)
			if score < 0 {
				score = 0
			}
			out.Closed[i] = annotateDecayScore(e, score)
		}
		return out, nil
	}
}

func decayKey(e scanner.Entry) string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

func annotateDecayScore(e scanner.Entry, score float64) scanner.Entry {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[decayScoreMetaKey] = fmt.Sprintf("%.4f", score)
	return e
}
