package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// WithSnapshot adds a stage that persists the current port state after each scan diff.
func WithSnapshot(path string) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}
		// Load existing snapshot, apply diff, and save.
		current, err := snapshot.Load(path)
		if err != nil {
			return d, err
		}
		portSet := make(map[uint16]bool, len(current))
		for _, p := range current {
			portSet[p] = true
		}
		for _, p := range d.Opened {
			portSet[p.Port] = true
		}
		for _, p := range d.Closed {
			delete(portSet, p.Port)
		}
		updated := make([]uint16, 0, len(portSet))
		for p := range portSet {
			updated = append(updated, p)
		}
		if err := snapshot.Save(path, updated); err != nil {
			return d, err
		}
		return d, nil
	}
}
