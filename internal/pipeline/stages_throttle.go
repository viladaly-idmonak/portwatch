package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/throttle"
)

// WithThrottle adds a throttle stage that drops diffs that arrive too
// frequently. cfg is obtained via throttle.NewFromConfig or throttle.DefaultConfig.
func WithThrottle(t *throttle.Throttle) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if !t.Allow() {
			return scanner.Diff{}, fmt.Errorf("throttled: diff dropped")
		}
		return diff, nil
	}
}
