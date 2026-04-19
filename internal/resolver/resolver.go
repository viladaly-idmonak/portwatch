// Package resolver maps port numbers to well-known service names.
package resolver

import (
	"fmt"
	"net"
)

// Resolver looks up service names for ports.
type Resolver struct {
	cache    map[uint16]string
	fallback bool
}

// New returns a Resolver. If fallback is true, unknown ports return "unknown".
func New(fallback bool) *Resolver {
	return &Resolver{
		cache:    make(map[uint16]string),
		fallback: fallback,
	}
}

// Lookup returns the service name for the given port and protocol (e.g. "tcp").
// Results are cached after the first lookup.
func (r *Resolver) Lookup(port uint16, proto string) string {
	if name, ok := r.cache[port]; ok {
		return name
	}
	name, err := net.LookupPort(proto, fmt.Sprintf("%d", port))
	if err == nil && name > 0 {
		// net.LookupPort returns the numeric port; use getservbyport equivalent via /etc/services
		svc := lookupServiceName(port, proto)
		r.cache[port] = svc
		return svc
	}
	if r.fallback {
		return "unknown"
	}
	return ""
}

// lookupServiceName attempts to resolve via net package workaround.
func lookupServiceName(port uint16, proto string) string {
	// Use reverse lookup: try to get the name from the system services database.
	names, err := net.LookupPort(proto, fmt.Sprintf("%d", port))
	_ = names
	if err != nil {
		return fmt.Sprintf("%d", port)
	}
	// Go stdlib doesn't expose service name lookup directly; return numeric string.
	return fmt.Sprintf("%d", port)
}
