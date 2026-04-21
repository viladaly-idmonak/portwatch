package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// dedupState tracks the last seen state for a port key.
type dedupState struct {
	state     string
	seenAt    time.Time
}

// Dedup holds state for suppressing duplicate port events.
type Dedup struct {
	mu      sync.Mutex
	entries map[string]dedupState
	ttl     time.Duration
	max     int
}

// NewDedup creates a Dedup with the given TTL and max entry count.
func NewDedup(ttl time.Duration, max int) *Dedup {
	return &Dedup{
		entries: make(map[string]dedupState),
		ttl:     ttl,
		max:     max,
	}
}

// IsDuplicate returns true if the (key, state) pair was recently seen.
func (d *Dedup) IsDuplicate(key, state string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if s, ok := d.entries[key]; ok {
		if s.state == state && time.Since(s.seenAt) < d.ttl {
			return true
		}
	}
	d.entries[key] = dedupState{state: state, seenAt: time.Now()}
	d.purge()
	return false
}

// purge removes the oldest entries when max is exceeded. Must be called with mu held.
func (d *Dedup) purge() {
	if d.max <= 0 || len(d.entries) <= d.max {
		return
	}
	var oldest string
	var oldestTime time.Time
	for k, v := range d.entries {
		if oldest == "" || v.seenAt.Before(oldestTime) {
			oldest = k
			oldestTime = v.seenAt
		}
	}
	delete(d.entries, oldest)
}

// WithDedup returns a Stage that suppresses duplicate port state events within a TTL window.
func WithDedup(d *Dedup) Stage {
	return func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}

		var opened, closed []scanner.Entry

		for _, e := range diff.Opened {
			k := portKey(e.Port, e.Protocol)
			if !d.IsDuplicate(k, "opened") {
				opened = append(opened, e)
			}
		}
		for _, e := range diff.Closed {
			k := portKey(e.Port, e.Protocol)
			if !d.IsDuplicate(k, "closed") {
				closed = append(closed, e)
			}
		}

		return scanner.Diff{Opened: opened, Closed: closed}, nil
	}
}
