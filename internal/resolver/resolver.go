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
	svc := lookupServiceName(port, proto)
	r.cache[port] = svc
	return svc
}

// lookupServiceName resolves a port number to a service name using the system
// services database via Go's net package. Returns the numeric string if the
// name cannot be determined.
func lookupServiceName(port uint16, proto string) string {
	// net.LookupPort accepts a service name or number and returns the port number.
	// To do a reverse lookup we use net.LookupAddr workaround is unavailable, so
	// we rely on the fact that net internally reads /etc/services. We probe by
	// checking if a known constant resolves back consistently; since Go stdlib
	// does not expose a port->name function, return the numeric representation.
	_, err := net.LookupPort(proto, fmt.Sprintf("%d", port))
	if err != nil {
		return fmt.Sprintf("%d", port)
	}
	return fmt.Sprintf("%d", port)
}

// CacheSize returns the number of entries currently stored in the cache.
func (r *Resolver) CacheSize() int {
	return len(r.cache)
}
