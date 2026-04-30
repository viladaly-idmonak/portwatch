// Package decay implements a time-based score decay mechanism that reduces
// the significance of port events over time, useful for deprioritising
// long-standing open ports in favour of recently changed ones.
package decay

import (
	"math"
	"sync"
	"time"
)

// ErrInvalidHalfLife is returned when a non-positive half-life is provided.
var ErrInvalidHalfLife = errInvalidHalfLife{}

type errInvalidHalfLife struct{}

func (errInvalidHalfLife) Error() string { return "decay: half-life must be greater than zero" }

// Entry holds the current score and the last time it was updated.
type Entry struct {
	Score     float64
	UpdatedAt time.Time
}

// Decayer applies exponential decay to per-key scores.
type Decayer struct {
	mu       sync.Mutex
	entries  map[string]*Entry
	halfLife time.Duration
}

// New creates a Decayer with the given half-life duration.
func New(halfLife time.Duration) (*Decayer, error) {
	if halfLife <= 0 {
		return nil, ErrInvalidHalfLife
	}
	return &Decayer{
		entries:  make(map[string]*Entry),
		halfLife: halfLife,
	}, nil
}

// Add increments the score for key by delta after applying decay since the
// last update. Returns the new score.
func (d *Decayer) Add(key string, delta float64) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	e, ok := d.entries[key]
	if !ok {
		e = &Entry{Score: 0, UpdatedAt: now}
		d.entries[key] = e
	}
	e.Score = d.decayed(e.Score, e.UpdatedAt, now) + delta
	e.UpdatedAt = now
	return e.Score
}

// Score returns the current decayed score for key without modifying state.
func (d *Decayer) Score(key string) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	e, ok := d.entries[key]
	if !ok {
		return 0
	}
	return d.decayed(e.Score, e.UpdatedAt, time.Now())
}

// Delete removes the tracked entry for key.
func (d *Decayer) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, key)
}

// decayed computes score * 0.5^(elapsed/halfLife).
func (d *Decayer) decayed(score float64, from, to time.Time) float64 {
	elapsed := to.Sub(from)
	if elapsed <= 0 {
		return score
	}
	exponent := float64(elapsed) / float64(d.halfLife)
	return score * math.Pow(0.5, exponent)
}
