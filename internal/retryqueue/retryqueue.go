// Package retryqueue provides a bounded, in-memory queue that retries failed
// pipeline diffs with exponential back-off.
package retryqueue

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Handler is a function that processes a diff and returns an error if it fails.
type Handler func(ctx context.Context, diff scanner.Diff) error

// entry holds a diff together with retry bookkeeping data.
type entry struct {
	diff    scanner.Diff
	attempt int
	nextAt  time.Time
}

// Queue retries failed diffs with exponential back-off.
type Queue struct {
	cfg     Config
	handler Handler
	mu      sync.Mutex
	items   []entry
}

// New creates a Queue with the given config and downstream handler.
func New(cfg Config, handler Handler) (*Queue, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Queue{cfg: cfg, handler: handler}, nil
}

// Enqueue adds a diff to the retry queue if capacity allows.
// It is safe to call from multiple goroutines.
func (q *Queue) Enqueue(diff scanner.Diff) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.cfg.MaxSize {
		return false
	}
	q.items = append(q.items, entry{diff: diff, attempt: 0, nextAt: time.Now()})
	return true
}

// Flush processes all due entries once, removing successes and updating
// back-off for failures. It returns the number of entries still pending.
func (q *Queue) Flush(ctx context.Context) int {
	q.mu.Lock()
	due := make([]entry, 0, len(q.items))
	deferred := q.items[:0]
	now := time.Now()
	for _, e := range q.items {
		if now.Before(e.nextAt) {
			deferred = append(deferred, e)
		} else {
			due = append(due, e)
		}
	}
	q.items = deferredCopy(deferred)
	q.mu.Unlock()

	for _, e := range due {
		if err := q.handler(ctx, e.diff); err != nil {
			e.attempt++
			if e.attempt <= q.cfg.MaxAttempts {
				delay := q.cfg.BaseDelay * time.Duration(1<<(e.attempt-1))
				if delay > q.cfg.MaxDelay {
					delay = q.cfg.MaxDelay
				}
				e.nextAt = time.Now().Add(delay)
				q.mu.Lock()
				q.items = append(q.items, e)
				q.mu.Unlock()
			}
		}
	}

	q.mu.Lock()
	n := len(q.items)
	q.mu.Unlock()
	return n
}

// Len returns the current number of pending entries.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func deferredCopy(src []entry) []entry {
	out := make([]entry, len(src))
	copy(out, src)
	return out
}
