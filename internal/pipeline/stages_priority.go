package pipeline

import (
	"context"

	"github.com/yourorg/portwatch/internal/scanner"
)

// Priority levels assigned to diff entries via meta key "priority".
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// PriorityRule maps a port/protocol pair to a priority level.
type PriorityRule struct {
	Port     uint16
	Proto    string
	Priority string
}

// Prioritizer assigns priority levels to diff entries based on configured rules.
type Prioritizer struct {
	rules   []PriorityRule
	default_ string
}

// NewPrioritizer creates a Prioritizer with the given rules and a fallback default.
func NewPrioritizer(rules []PriorityRule, defaultPriority string) *Prioritizer {
	if defaultPriority == "" {
		defaultPriority = PriorityLow
	}
	return &Prioritizer{rules: rules, default_: defaultPriority}
}

func (p *Prioritizer) assign(e scanner.Entry) string {
	for _, r := range p.rules {
		if r.Port == e.Port && r.Proto == e.Proto {
			return r.Priority
		}
	}
	return p.default_
}

// WithPriority returns a Stage that annotates each entry with a "priority" meta key.
// A nil prioritizer is a no-op.
func WithPriority(p *Prioritizer) Stage {
	if p == nil {
		return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
			return d, nil
		}
	}
	return func(_ context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}
		return scanner.Diff{
			Opened: applyPriority(p, d.Opened),
			Closed: applyPriority(p, d.Closed),
		}, nil
	}
}

func applyPriority(p *Prioritizer, entries []scanner.Entry) []scanner.Entry {
	result := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		meta := make(map[string]string, len(e.Meta)+1)
		for k, v := range e.Meta {
			meta[k] = v
		}
		meta["priority"] = p.assign(e)
		e.Meta = meta
		result[i] = e
	}
	return result
}
