package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/cooldown"
	"github.com/user/portwatch/internal/scanner"
)

// WithCooldown returns a pipeline stage that suppresses repeated events for the
// same port+protocol+state combination within a cooldown window. Unlike
// WithSuppress (which tracks any state), WithCooldown is keyed on the full
// (port, protocol, state) triple so an open→close→open sequence is never
// silently dropped.
func WithCooldown(cd *cooldown.Cooldown) Stage {
	if cd == nil {
		return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}

	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		filtered := scanner.Diff{}

		for _, e := range d.Opened {
			k := cooldownKey(e, "opened")
			if cd.Allow(k) {
				filtered.Opened = append(filtered.Opened, e)
			}
		}

		for _, e := range d.Closed {
			k := cooldownKey(e, "closed")
			if cd.Allow(k) {
				filtered.Closed = append(filtered.Closed, e)
			}
		}

		return filtered, nil
	}
}

func cooldownKey(e scanner.Entry, state string) string {
	return fmt.Sprintf("%s:%d:%s", e.Protocol, e.Port, state)
}
