// Package suppress provides a time-based suppression window that prevents
// repeated events for the same port+state combination from propagating
// through the pipeline within a configurable cooldown period.
package suppress

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// key uniquely identifies a port+state event.
type key struct {
	port  uint16
	proto string
	state string
}

// Suppressor tracks recently seen events and suppresses duplicates within a
// sliding window.
type Suppressor struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[key]time.Time
	nowFunc func() time.Time
}

// New creates a Suppressor with the given suppression window.
func New(window time.Duration) (*Suppressor, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Suppressor{
		window:  window,
		seen:    make(map[key]time.Time),
		nowFunc: time.Now,
	}, nil
}

// Apply filters diff entries that were seen within the suppression window.
// Entries not yet seen, or seen outside the window, are allowed through.
func (s *Suppressor) Apply(diff scanner.Diff) scanner.Diff {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.nowFunc()
	out := scanner.Diff{}

	for _, e := range diff.Opened {
		k := key{port: e.Port, proto: e.Proto, state: "opened"}
		if last, ok := s.seen[k]; !ok || now.Sub(last) > s.window {
			s.seen[k] = now
			out.Opened = append(out.Opened, e)
		}
	}
	for _, e := range diff.Closed {
		k := key{port: e.Port, proto: e.Proto, state: "closed"}
		if last, ok := s.seen[k]; !ok || now.Sub(last) > s.window {
			s.seen[k] = now
			out.Closed = append(out.Closed, e)
		}
	}
	return out
}

// Purge removes entries whose suppression window has expired.
func (s *Suppressor) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.nowFunc()
	for k, t := range s.seen {
		if now.Sub(t) > s.window {
			delete(s.seen, k)
		}
	}
}
