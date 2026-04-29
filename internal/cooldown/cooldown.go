// Package cooldown provides a per-key cooldown mechanism that suppresses
// repeated events for a configurable duration after the first occurrence.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-seen time for arbitrary string keys and reports
// whether enough time has elapsed since the previous event.
type Cooldown struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	duration time.Duration
	now      func() time.Time
}

// New creates a Cooldown with the given suppression duration.
// Events for a key are allowed only when duration has elapsed since the last
// allowed event for that key.
func New(d time.Duration) (*Cooldown, error) {
	if d <= 0 {
		return nil, ErrInvalidDuration
	}
	return &Cooldown{
		lastSeen: make(map[string]time.Time),
		duration: d,
		now:      time.Now,
	}, nil
}

// Allow returns true if the key is not currently in cooldown.
// Calling Allow with a key that is allowed resets the cooldown timer for that
// key.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if last, ok := c.lastSeen[key]; ok {
		if now.Sub(last) < c.duration {
			return false
		}
	}
	c.lastSeen[key] = now
	return true
}

// Reset removes the cooldown entry for key, allowing the next event through
// immediately regardless of when the last event occurred.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastSeen, key)
}

// Purge removes all entries whose cooldown period has fully elapsed, freeing
// memory for keys that are no longer active.
func (c *Cooldown) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	for k, last := range c.lastSeen {
		if now.Sub(last) >= c.duration {
			delete(c.lastSeen, k)
		}
	}
}

// ErrInvalidDuration is returned when a non-positive duration is provided.
var ErrInvalidDuration = errInvalidDuration("cooldown duration must be positive")

type errInvalidDuration string

func (e errInvalidDuration) Error() string { return string(e) }
