package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/scanner"
)

// WithEnricher returns a Stage that attaches static metadata fields to every
// entry in the diff using the provided Enricher. If e is nil the stage is a
// no-op pass-through.
func WithEnricher(e *enricher.Enricher) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if e == nil {
			return d, nil
		}
		select {
		case <-ctx.Done():
			return scanner.Diff{}, ctx.Err()
		default:
		}
		return e.Apply(d), nil
	}
}
