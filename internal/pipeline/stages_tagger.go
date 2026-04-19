package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

// WithTagger is a pipeline stage that annotates each port entry in the diff
// with a human-readable tag by embedding it in the port's metadata label.
// The diff itself is passed through unchanged; tagging is advisory/logging use.
func WithTagger(tr *tagger.Tagger) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}
		if tr == nil {
			return d, nil
		}
		// Annotate opened ports — store tag in a side-channel log if needed.
		// For now we enrich the diff entries' Tag field when present.
		annotate := func(entries []scanner.Entry) []scanner.Entry {
			out := make([]scanner.Entry, len(entries))
			for i, e := range entries {
				e.Tag = tr.Label(e.Port)
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
