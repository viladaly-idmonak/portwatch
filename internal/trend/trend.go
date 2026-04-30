package trend

import (
	"sync"
	"time"
)

// Direction represents the direction of a port activity trend.
type Direction int

const (
	DirectionNeutral Direction = iota
	DirectionRising
	DirectionFalling
)

// String returns a human-readable label for the direction.
func (d Direction) String() string {
	switch d {
	case DirectionRising:
		return "rising"
	case DirectionFalling:
		return "falling"
	default:
		return "neutral"
	}
}

// bucket holds open/close event counts in a time window.
type bucket struct {
	at      time.Time
	opened int
	closed int
}

// Tracker accumulates port events over a sliding window and computes a trend
// direction based on the ratio of opened-to-closed events.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	buckets []bucket
}

// New returns a Tracker with the given sliding window duration.
func New(window time.Duration) (*Tracker, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Tracker{window: window}, nil
}

// Record adds an event to the tracker.
func (t *Tracker) Record(opened, closed int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	t.buckets = append(t.buckets, bucket{
		at:     time.Now(),
		opened: opened,
		closed: closed,
	})
}

// Trend returns the current trend direction based on accumulated events.
func (t *Tracker) Trend() Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	var totalOpened, totalClosed int
	for _, b := range t.buckets {
		totalOpened += b.opened
		totalClosed += b.closed
	}
	if totalOpened == 0 && totalClosed == 0 {
		return DirectionNeutral
	}
	if totalOpened > totalClosed {
		return DirectionRising
	}
	if totalClosed > totalOpened {
		return DirectionFalling
	}
	return DirectionNeutral
}

// evict removes buckets older than the window. Must be called with t.mu held.
func (t *Tracker) evict() {
	cutoff := time.Now().Add(-t.window)
	i := 0
	for i < len(t.buckets) && t.buckets[i].at.Before(cutoff) {
		i++
	}
	t.buckets = t.buckets[i:]
}
