package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

// WithWatchlist returns a Stage that ensures any port in the watchlist is
// always present in the Opened set when it appears, and never silently dropped
// by upstream stages. Ports not in the watchlist are passed through unchanged.
//
// Concretely: if a watched port is in diff.Closed it is moved back to
// diff.Opened, signalling that the port is expected to be open.
func WithWatchlist(wl *watchlist.Watchlist) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}

		var filtered []scanner.Entry
		for _, e := range d.Closed {
			if wl.Contains(e.Protocol, e.Port) {
				// Re-surface as opened so downstream stages can alert on it.
				d.Opened = append(d.Opened, e)
			} else {
				filtered = append(filtered, e)
			}
		}
		d.Closed = filtered
		return d, nil
	}
}
