package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/geoip"
	"github.com/user/portwatch/internal/scanner"
)

// WithGeoIP returns a Stage that annotates each opened port entry with
// geolocation metadata stored as a formatted label in the diff entries.
// Closed entries are passed through unchanged.
func WithGeoIP(lookup geoip.Lookup) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		annotated := make([]scanner.Entry, 0, len(d.Opened))
		for _, e := range d.Opened {
			info, err := lookup.Lookup(e.Addr)
			if err != nil {
				// Non-fatal: attach error note and continue.
				e.Meta = fmt.Sprintf("geoip_error=%s", err.Error())
			} else if info.Country != "" {
				e.Meta = fmt.Sprintf("country=%s city=%s org=%s", info.Country, info.City, info.Org)
			}
			annotated = append(annotated, e)
		}
		return scanner.Diff{Opened: annotated, Closed: d.Closed}, nil
	}
}
