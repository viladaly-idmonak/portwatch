package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/labeler"
	"github.com/user/portwatch/internal/scanner"
)

// WithLabeler returns a pipeline stage that annotates each diff entry with a
// human-readable label based on port/protocol rules defined in the labeler.
// Entries that match no rule are left unchanged. The stage is a no-op when
// the provided labeler is nil.
func WithLabeler(lb *labeler.Labeler) Stage {
	if lb == nil {
		return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}

	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		out := scanner.Diff{
			Opened: annotateEntries(lb, d.Opened),
			Closed: annotateEntries(lb, d.Closed),
		}
		return out, nil
	}
}

func annotateEntries(lb *labeler.Labeler, entries []scanner.Entry) []scanner.Entry {
	if len(entries) == 0 {
		return entries
	}
	result := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		label := lb.LabelOrDefault(e.Port, e.Proto, "")
		if label != "" {
			if e.Meta == nil {
				e.Meta = make(map[string]string)
			} else {
				meta := make(map[string]string, len(e.Meta))
				for k, v := range e.Meta {
					meta[k] = v
				}
				e.Meta = meta
			}
			e.Meta["label"] = label
		}
		result[i] = e
	}
	return result
}
