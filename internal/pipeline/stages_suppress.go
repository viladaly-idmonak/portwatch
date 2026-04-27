package pipeline

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/suppress"
)

// WithSuppress returns a Stage that drops repeated port+state events seen
// within the given suppression window. This is useful for noisy environments
// where the same port flaps repeatedly in a short period.
//
// Events that have not been seen before, or whose last occurrence falls
// outside the window, are passed through unchanged.
func WithSuppress(window time.Duration) (Stage, error) {
	s, err := suppress.New(window)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		return s.Apply(diff), nil
	}, nil
}
