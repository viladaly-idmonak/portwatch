package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/scanner"
)

// WithAudit returns a Stage that writes every diff to the provided Auditor.
// If a is nil the stage is a no-op pass-through.
func WithAudit(a *audit.Auditor) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if a == nil {
			return d, nil
		}
		if err := a.Record(d); err != nil {
			return d, err
		}
		return d, nil
	}
}
