package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

// WithRateLimitPerPort filters diff entries using a per-port rate limiter,
// dropping events that occur within the configured deduplication window.
func WithRateLimitPerPort(rl *ratelimit.RateLimit) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		filtered := scanner.Diff{}

		for _, p := range d.Opened {
			key := "open:" + uint16ToStr(p.Port)
			if rl.Allow(key) {
				filtered.Opened = append(filtered.Opened, p)
			}
		}

		for _, p := range d.Closed {
			key := "close:" + uint16ToStr(p.Port)
			if rl.Allow(key) {
				filtered.Closed = append(filtered.Closed, p)
			}
		}

		return filtered, nil
	}
}
