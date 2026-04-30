package quota

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a key has exceeded its allowed budget.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Quota tracks a rolling event budget per key over a fixed window.
type Quota struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	window  time.Duration
	max     int
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New creates a Quota that allows at most max events per key within window.
func New(max int, window time.Duration) (*Quota, error) {
	if max <= 0 {
		return nil, errors.New("quota: max must be positive")
	}
	if window <= 0 {
		return nil, errors.New("quota: window must be positive")
	}
	return &Quota{
		buckets: make(map[string]*bucket),
		window:  window,
		max:     max,
	}, nil
}

// Allow reports whether the event for key is within quota.
// It increments the counter and returns ErrQuotaExceeded when the budget is exhausted.
func (q *Quota) Allow(key string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	b, ok := q.buckets[key]
	if !ok || now.After(b.windowEnd) {
		q.buckets[key] = &bucket{count: 1, windowEnd: now.Add(q.window)}
		return nil
	}
	if b.count >= q.max {
		return ErrQuotaExceeded
	}
	b.count++
	return nil
}

// Remaining returns how many events key may still emit in the current window.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	b, ok := q.buckets[key]
	if !ok || now.After(b.windowEnd) {
		return q.max
	}
	r := q.max - b.count
	if r < 0 {
		return 0
	}
	return r
}

// Purge removes expired buckets to reclaim memory.
func (q *Quota) Purge() {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	for k, b := range q.buckets {
		if now.After(b.windowEnd) {
			delete(q.buckets, k)
		}
	}
}
