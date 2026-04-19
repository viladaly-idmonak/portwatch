package tagger

import (
	"fmt"
	"sync"
)

// Tagger assigns human-readable labels to ports.
type Tagger struct {
	mu   sync.RWMutex
	tags map[uint16]string
}

// New returns a Tagger pre-loaded with the given tags map.
func New(initial map[uint16]string) *Tagger {
	t := &Tagger{tags: make(map[uint16]string, len(initial))}
	for k, v := range initial {
		t.tags[k] = v
	}
	return t
}

// Set assigns a label to a port, overwriting any existing tag.
func (t *Tagger) Set(port uint16, label string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tags[port] = label
}

// Get returns the label for a port and whether it was found.
func (t *Tagger) Get(port uint16) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.tags[port]
	return v, ok
}

// Label returns the tag for a port, or a default string if not tagged.
func (t *Tagger) Label(port uint16) string {
	if v, ok := t.Get(port); ok {
		return v
	}
	return fmt.Sprintf("port/%d", port)
}

// Delete removes a tag for a port.
func (t *Tagger) Delete(port uint16) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tags, port)
}

// All returns a copy of all current tags.
func (t *Tagger) All() map[uint16]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[uint16]string, len(t.tags))
	for k, v := range t.tags {
		out[k] = v
	}
	return out
}
