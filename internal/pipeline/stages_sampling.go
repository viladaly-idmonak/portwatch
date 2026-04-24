package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/sampling"
	"github.com/user/portwatch/internal/scanner"
)

// WithSampling returns a Stage that probabilistically forwards diffs based on
// the configured sample rate. When s is nil the stage is a no-op pass-through,
// allowing callers to omit sampling by passing a nil *sampling.Sampler.
//
// Example usage:
//
//	sampler, _ := sampling.New(0.1) // forward ~10 % of events
//	p := pipeline.New(pipeline.WithSampling(sampler))
func WithSampling(s *sampling.Sampler) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if s == nil {
			return d, nil
		}
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		return s.Apply(d), nil
	}
}
