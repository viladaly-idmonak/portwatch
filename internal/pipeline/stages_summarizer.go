package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summarizer"
)

// WithSummarizer returns a Stage that records each diff into the provided
// Summarizer. The diff passes through unchanged.
func WithSummarizer(s *summarizer.Summarizer) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		s.Record(d)
		return d, nil
	}
}
