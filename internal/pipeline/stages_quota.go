package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/portwatch/internal/quota"
	"github.com/yourorg/portwatch/internal/scanner"
)

// Quoter wraps a quota.Quota to filter diff entries that exceed their budget.
type Quoter struct {
	q      *quota.Quota
	keyFn  func(scanner.Entry) string
}

// NewQuoter creates a Quoter allowing max events per key within window.
func NewQuoter(max int, window time.Duration) (*Quoter, error) {
	q, err := quota.New(max, window)
	if err != nil {
		return nil, fmt.Errorf("stages_quota: %w", err)
	}
	return &Quoter{
		q:     q,
		keyFn: func(e scanner.Entry) string { return quotaKey(e) },
	}, nil
}

func quotaKey(e scanner.Entry) string {
	return fmt.Sprintf("%d/%s", e.Port, e.Protocol)
}

// WithQuota returns a Stage that drops entries whose per-port quota is exhausted.
func WithQuota(qr *Quoter) Stage {
	if qr == nil {
		return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		out := scanner.Diff{}
		for _, e := range d.Opened {
			if err := qr.q.Allow(qr.keyFn(e)); err == nil {
				out.Opened = append(out.Opened, e)
			}
		}
		for _, e := range d.Closed {
			if err := qr.q.Allow(qr.keyFn(e)); err == nil {
				out.Closed = append(out.Closed, e)
			}
		}
		return out, nil
	}
}
