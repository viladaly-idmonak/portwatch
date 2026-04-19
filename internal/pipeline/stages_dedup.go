package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// WithDedup returns a Stage that suppresses diff entries whose state has not
// changed since the last time they were seen. This prevents duplicate events
// from propagating downstream when the scanner emits the same change twice
// within a short window (e.g. due to a flapping port).
func WithDedup() Stage {
	seen := make(map[string]string) // key -> last state ("opened" | "closed")

	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		var opened, closed []scanner.Entry

		for _, e := range d.Opened {
			k := portKey(e.Port, e.Proto)
			if seen[k] == "opened" {
				continue
			}
			seen[k] = "opened"
			opened = append(opened, e)
		}

		for _, e := range d.Closed {
			k := portKey(e.Port, e.Proto)
			if seen[k] == "closed" {
				continue
			}
			seen[k] = "closed"
			closed = append(closed, e)
		}

		return scanner.Diff{Opened: opened, Closed: closed}, nil
	}
}
