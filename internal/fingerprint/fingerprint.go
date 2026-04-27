// Package fingerprint identifies ports by computing a stable hash
// over their key attributes, enabling change detection across scans.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint holds the computed hash and source attributes for a port entry.
type Fingerprint struct {
	Port     uint16
	Protocol string
	State    string
	Hash     string
}

// Hasher computes fingerprints for port entries.
type Hasher struct {
	salt string
}

// New returns a Hasher. An optional salt can be provided to namespace hashes.
func New(salt string) *Hasher {
	return &Hasher{salt: salt}
}

// Compute derives a Fingerprint from a scanner.Entry.
func (h *Hasher) Compute(e scanner.Entry) Fingerprint {
	raw := fmt.Sprintf("%s:%d:%s:%s", h.salt, e.Port, e.Protocol, e.State)
	sum := sha256.Sum256([]byte(raw))
	return Fingerprint{
		Port:     e.Port,
		Protocol: e.Protocol,
		State:    e.State,
		Hash:     hex.EncodeToString(sum[:8]),
	}
}

// ComputeAll returns fingerprints for every entry in the slice.
func (h *Hasher) ComputeAll(entries []scanner.Entry) []Fingerprint {
	out := make([]Fingerprint, 0, len(entries))
	for _, e := range entries {
		out = append(out, h.Compute(e))
	}
	return out
}

// Equal reports whether two fingerprints share the same hash.
func Equal(a, b Fingerprint) bool {
	return a.Hash == b.Hash
}
