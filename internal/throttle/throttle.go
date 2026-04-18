package throttle

import (
	"sync"
	"time"
)

// Throttle limits how frequently a scan cycle can be triggered.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
}

// New creates a Throttle that enforces a minimum interval between calls.
func New(interval time.Duration) *Throttle {
	return &Throttle{interval: interval}
}

// Allow returns true if enough time has passed since the last allowed call.
// If allowed, it updates the internal timestamp.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if now.Sub(t.last) < t.interval {
		return false
	}
	t.last = now
	return true
}

// Reset clears the last-seen timestamp, allowing the next call immediately.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = time.Time{}
}

// SetInterval updates the throttle interval at runtime.
func (t *Throttle) SetInterval(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.interval = d
}
