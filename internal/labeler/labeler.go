// Package labeler assigns human-readable labels to ports based on
// configurable rules, falling back to the resolver for known services.
package labeler

import "fmt"

// Rule maps a port/protocol pair to a custom label.
type Rule struct {
	Port     uint16
	Protocol string
	Label    string
}

// Labeler assigns labels to ports.
type Labeler struct {
	rules map[string]string
}

// New creates a Labeler from the provided rules. Protocol is
// normalised to lowercase. Duplicate keys use the last value.
func New(rules []Rule) *Labeler {
	m := make(map[string]string, len(rules))
	for _, r := range rules {
		if r.Label == "" || r.Protocol == "" {
			continue
		}
		m[key(r.Port, r.Protocol)] = r.Label
	}
	return &Labeler{rules: m}
}

// Label returns the custom label for port/protocol, or an empty
// string when no rule matches.
func (l *Labeler) Label(port uint16, protocol string) string {
	return l.rules[key(port, protocol)]
}

// LabelOrDefault returns the custom label when present, otherwise
// returns the provided fallback string.
func (l *Labeler) LabelOrDefault(port uint16, protocol, fallback string) string {
	if v := l.Label(port, protocol); v != "" {
		return v
	}
	return fallback
}

// All returns a copy of all registered rules.
func (l *Labeler) All() map[string]string {
	out := make(map[string]string, len(l.rules))
	for k, v := range l.rules {
		out[k] = v
	}
	return out
}

func key(port uint16, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
