package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

// WithFilter returns a Stage that applies a filter.Filter to the diff,
// removing ports that do not pass the filter rules.
func WithFilter(f *filter.Filter) Stage {
	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		return f.Apply(d), nil
	}
}

// WithRateLimit returns a Stage that drops individual port events that are
// within the rate-limit window for their key.
func WithRateLimit(rl *ratelimit.RateLimit) Stage {
	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		filtered := scanner.Diff{}
		for _, p := range d.Opened {
			key := portKey(p, "opened")
			if rl.Allow(key) {
				filtered.Opened = append(filtered.Opened, p)
			}
		}
		for _, p := range d.Closed {
			key := portKey(p, "closed")
			if rl.Allow(key) {
				filtered.Closed = append(filtered.Closed, p)
			}
		}
		return filtered, nil
	}
}

func portKey(port uint16, state string) string {
	return state + ":" + uint16ToStr(port)
}

func uint16ToStr(n uint16) string {
	if n == 0 {
		return "0"
	}
	buf := [5]byte{}
	i := 5
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
