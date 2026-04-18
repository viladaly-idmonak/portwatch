package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated diff events within a cooldown window.
// If the same port+state combination is seen within the window, it is dropped.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the event identified by key should be allowed through.
// Repeated keys within the cooldown window return false.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if last, ok := l.seen[key]; ok && now.Sub(last) < l.cooldown {
		return false
	}
	l.seen[key] = now
	return true
}

// Purge removes expired entries from the seen map to prevent unbounded growth.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for k, t := range l.seen {
		if now.Sub(t) >= l.cooldown {
			delete(l.seen, k)
		}
	}
}
