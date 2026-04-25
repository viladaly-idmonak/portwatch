package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/redactor"
	"github.com/user/portwatch/internal/scanner"
)

// WithRedactor returns a Stage that masks configured meta keys in every
// diff that passes through the pipeline. Useful before forwarding diffs
// to external sinks such as alerters or notifiers.
//
// Example:
//
//	r := redactor.New([]string{"api_key", "token"})
//	p := pipeline.New(pipeline.WithRedactor(r))
func WithRedactor(r *redactor.Redactor) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		return r.Apply(diff), nil
	}
}
