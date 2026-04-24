// Package sampling provides a stage that probabilistically samples port events,
// forwarding only a configured fraction of diffs to downstream pipeline stages.
package sampling

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Sampler forwards diffs with probability Rate (0.0–1.0).
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New creates a Sampler with the given sample rate.
// rate must be in the range (0.0, 1.0].
func New(rate float64) (*Sampler, error) {
	if rate <= 0 || rate > 1 {
		return nil, fmt.Errorf("sampling: rate must be in (0, 1], got %v", rate)
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(rand.Int63())),
	}, nil
}

// Sample returns true if the event should be forwarded based on the configured rate.
func (s *Sampler) Sample() bool {
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Apply filters diff entries probabilistically, returning a new diff containing
// only the sampled entries. If no entries survive sampling, an empty diff is returned.
func (s *Sampler) Apply(d scanner.Diff) scanner.Diff {
	out := scanner.Diff{}
	for _, e := range d.Opened {
		if s.Sample() {
			out.Opened = append(out.Opened, e)
		}
	}
	for _, e := range d.Closed {
		if s.Sample() {
			out.Closed = append(out.Closed, e)
		}
	}
	return out
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rate
}

// SetRate updates the sample rate at runtime. Returns an error if the new rate is invalid.
func (s *Sampler) SetRate(rate float64) error {
	if rate <= 0 || rate > 1 {
		return fmt.Errorf("sampling: rate must be in (0, 1], got %v", rate)
	}
	s.mu.Lock()
	s.rate = rate
	s.mu.Unlock()
	return nil
}
