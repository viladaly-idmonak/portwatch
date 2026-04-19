// Package watchlist maintains a user-defined set of ports that should always
// be monitored, regardless of the active filter rules.
package watchlist

import (
	"fmt"
	"sync"
)

// Entry represents a single watched port.
type Entry struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
	Note     string
}

// Watchlist holds a thread-safe collection of port entries.
type Watchlist struct {
	mu      sync.RWMutex
	entries map[string]Entry // key: "proto:port"
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{entries: make(map[string]Entry)}
}

func key(proto string, port uint16) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Add inserts or replaces an entry.
func (w *Watchlist) Add(e Entry) error {
	if e.Protocol != "tcp" && e.Protocol != "udp" {
		return fmt.Errorf("watchlist: invalid protocol %q", e.Protocol)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[key(e.Protocol, e.Port)] = e
	return nil
}

// Remove deletes an entry. It is a no-op if the entry does not exist.
func (w *Watchlist) Remove(proto string, port uint16) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, key(proto, port))
}

// Contains reports whether the given port/protocol pair is watched.
func (w *Watchlist) Contains(proto string, port uint16) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.entries[key(proto, port)]
	return ok
}

// All returns a snapshot of all current entries.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of watched entries.
func (w *Watchlist) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.entries)
}
