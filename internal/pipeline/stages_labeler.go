package pipeline

import (
	"context"

	"github.com/yourorg/portwatch/internal/labeler"
	"github.com/yourorg/portwatch/internal/scanner"
)

// WithLabeler returns a Stage that annotates each diff entry with a human-readable
// label derived from the port/protocol pair. Entries for unknown ports are left
// unannotated. A nil labeler is a no-op.
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
	result := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		label := lb.Label(e.Port, e.Proto)
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
