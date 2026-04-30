// Package anomaly detects statistically unusual port activity by tracking
// event frequency and flagging entries that exceed a configured threshold.
package anomaly

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Detector tracks per-port event counts within a sliding window and marks
// entries whose frequency exceeds a configured threshold as anomalous.
type Detector struct {
	mu        sync.Mutex
	counts    map[string][]time.Time
	window    time.Duration
	threshold int
	now       func() time.Time
}

// New returns a Detector that flags ports seen more than threshold times
// within the given window duration.
func New(window time.Duration, threshold int) (*Detector, error) {
	if window <= 0 {
		return nil, fmt.Errorf("anomaly: window must be positive, got %s", window)
	}
	if threshold <= 0 {
		return nil, fmt.Errorf("anomaly: threshold must be positive, got %d", threshold)
	}
	return &Detector{
		counts:    make(map[string][]time.Time),
		window:    window,
		threshold: threshold,
		now:       time.Now,
	}, nil
}

// key returns a stable string key for a scanner entry.
func key(e scanner.Entry) string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

// Record adds an observation for the entry and reports whether the event
// count within the window now exceeds the threshold.
func (d *Detector) Record(e scanner.Entry) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.window)
	k := key(e)

	// Evict observations outside the window.
	filtered := d.counts[k][:0]
	for _, t := range d.counts[k] {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	d.counts[k] = filtered

	return len(filtered) > d.threshold
}

// Reset clears all observations for the given entry.
func (d *Detector) Reset(e scanner.Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.counts, key(e))
}
