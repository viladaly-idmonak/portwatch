package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Truncator limits the number of entries in a Diff that are passed downstream.
// When a Diff exceeds MaxEntries, the excess entries are dropped and a
// truncation annotation is added to the first retained entry's Meta map.
type Truncator struct {
	MaxEntries int
}

// NewTruncator returns a Truncator that retains at most maxEntries opened and
// maxEntries closed entries per Diff. maxEntries must be >= 1.
func NewTruncator(maxEntries int) (*Truncator, error) {
	if maxEntries < 1 {
		return nil, fmt.Errorf("truncator: maxEntries must be >= 1, got %d", maxEntries)
	}
	return &Truncator{MaxEntries: maxEntries}, nil
}

// Apply trims d.Opened and d.Closed to at most MaxEntries each.
// If either slice is trimmed, a "truncated" key is set on the first entry's
// Meta map with the original count as the value.
func (t *Truncator) Apply(d scanner.Diff) scanner.Diff {
	d.Opened = t.trim(d.Opened)
	d.Closed = t.trim(d.Closed)
	return d
}

func (t *Truncator) trim(entries []scanner.Entry) []scanner.Entry {
	if len(entries) <= t.MaxEntries {
		return entries
	}
	original := len(entries)
	entries = entries[:t.MaxEntries]
	if entries[0].Meta == nil {
		entries[0].Meta = make(map[string]string)
	}
	entries[0].Meta["truncated"] = fmt.Sprintf("%d", original)
	return entries
}

// WithTruncator returns a Stage that applies t to every non-empty Diff.
func WithTruncator(t *Truncator) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		return t.Apply(d), nil
	}
}
