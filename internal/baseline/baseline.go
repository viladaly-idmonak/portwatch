// Package baseline tracks a known-good set of open ports and flags deviations.
package baseline

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline holds a frozen snapshot of expected open ports.
type Baseline struct {
	mu    sync.RWMutex
	ports map[uint16]struct{}
}

// New creates a Baseline from the given port list.
func New(ports []uint16) *Baseline {
	m := make(map[uint16]struct{}, len(ports))
	for _, p := range ports {
		m[p] = struct{}{}
	}
	return &Baseline{ports: m}
}

// NewFromScan captures the current open ports as the baseline.
func NewFromScan(entries []scanner.Entry) *Baseline {
	m := make(map[uint16]struct{}, len(entries))
	for _, e := range entries {
		m[e.Port] = struct{}{}
	}
	return &Baseline{ports: m}
}

// Deviations returns ports that are open but not in the baseline (unexpected)
// and ports that are in the baseline but now closed (missing).
func (b *Baseline) Deviations(current []scanner.Entry) (unexpected, missing []uint16) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	currentSet := make(map[uint16]struct{}, len(current))
	for _, e := range current {
		currentSet[e.Port] = struct{}{}
		if _, ok := b.ports[e.Port]; !ok {
			unexpected = append(unexpected, e.Port)
		}
	}
	for p := range b.ports {
		if _, ok := currentSet[p]; !ok {
			missing = append(missing, p)
		}
	}
	return
}

// Update replaces the baseline with the given port list.
func (b *Baseline) Update(ports []uint16) {
	b.mu.Lock()
	defer b.mu.Unlock()
	m := make(map[uint16]struct{}, len(ports))
	for _, p := range ports {
		m[p] = struct{}{}
	}
	b.ports = m
}

// Contains reports whether port is in the baseline.
func (b *Baseline) Contains(port uint16) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.ports[port]
	return ok
}
