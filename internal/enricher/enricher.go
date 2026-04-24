// Package enricher attaches additional metadata fields to diff entries
// based on configurable key-value annotations sourced from a static map.
package enricher

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Enricher annotates diff entries with static metadata.
type Enricher struct {
	mu     sync.RWMutex
	fields map[string]string
}

// New returns an Enricher pre-loaded with the given fields.
func New(fields map[string]string) *Enricher {
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Enricher{fields: copy}
}

// Set adds or updates a metadata field.
func (e *Enricher) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.fields[key] = value
}

// Delete removes a metadata field.
func (e *Enricher) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.fields, key)
}

// Fields returns a snapshot of the current metadata fields.
func (e *Enricher) Fields() map[string]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make(map[string]string, len(e.fields))
	for k, v := range e.fields {
		out[k] = v
	}
	return out
}

// Apply returns a new Diff where every entry has the enricher's fields
// merged into its Meta map. Existing meta keys are not overwritten.
func (e *Enricher) Apply(d scanner.Diff) scanner.Diff {
	if len(d.Opened)+len(d.Closed) == 0 {
		return d
	}
	e.mu.RLock()
	fields := make(map[string]string, len(e.fields))
	for k, v := range e.fields {
		fields[k] = v
	}
	e.mu.RUnlock()

	enrich := func(entries []scanner.Entry) []scanner.Entry {
		out := make([]scanner.Entry, len(entries))
		for i, en := range entries {
			meta := make(map[string]string, len(en.Meta)+len(fields))
			for k, v := range fields {
				meta[k] = v
			}
			for k, v := range en.Meta {
				meta[k] = v // existing keys win
			}
			en.Meta = meta
			out[i] = en
		}
		return out
	}

	return scanner.Diff{
		Opened: enrich(d.Opened),
		Closed: enrich(d.Closed),
	}
}
