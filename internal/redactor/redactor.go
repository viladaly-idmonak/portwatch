// Package redactor masks sensitive metadata fields in diff entries
// before they are forwarded to external sinks (notifiers, alerters, audit logs).
package redactor

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

const redactedValue = "[REDACTED]"

// Redactor masks configured meta keys in diff entries.
type Redactor struct {
	keys map[string]struct{}
}

// New returns a Redactor that will mask the given meta keys.
// Key matching is case-insensitive.
func New(keys []string) *Redactor {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return &Redactor{keys: m}
}

// Apply returns a copy of diff with sensitive meta values replaced.
// The original diff is not modified.
func (r *Redactor) Apply(diff scanner.Diff) scanner.Diff {
	if len(r.keys) == 0 {
		return diff
	}
	return scanner.Diff{
		Opened: r.redactEntries(diff.Opened),
		Closed: r.redactEntries(diff.Closed),
	}
}

func (r *Redactor) redactEntries(entries []scanner.Entry) []scanner.Entry {
	if len(entries) == 0 {
		return entries
	}
	out := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		out[i] = r.redactEntry(e)
	}
	return out
}

func (r *Redactor) redactEntry(e scanner.Entry) scanner.Entry {
	if len(e.Meta) == 0 {
		return e
	}
	newMeta := make(map[string]string, len(e.Meta))
	for k, v := range e.Meta {
		if _, masked := r.keys[strings.ToLower(k)]; masked {
			newMeta[k] = redactedValue
		} else {
			newMeta[k] = v
		}
	}
	e.Meta = newMeta
	return e
}
