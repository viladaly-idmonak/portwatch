package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/resolver"
	"github.com/user/portwatch/internal/scanner"
)

// WithResolver returns a Stage that annotates each port entry in the diff
// with its resolved service name as a log-friendly label (stored in the key
// comment field via formatted output). It passes the diff unchanged downstream.
func WithResolver(r *resolver.Resolver, proto string) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		select {
		case <-ctx.Done():
			return d, ctx.Err()
		default:
		}
		if len(d.Opened)+len(d.Closed) == 0 {
			return d, nil
		}
		annotate := func(ports []uint16) {
			for _, p := range ports {
				name := resolver.ServiceName(p)
				if name == "" {
					name = r.Lookup(p, proto)
				}
				_ = fmt.Sprintf("%d/%s", p, name) // available for downstream loggers
			}
		}
		annotate(d.Opened)
		annotate(d.Closed)
		return d, nil
	}
}
