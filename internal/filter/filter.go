package filter

import "github.com/user/portwatch/internal/scanner"

// Config holds filtering rules for port events.
type Config struct {
	// IncludePorts restricts alerts to these ports. Empty means all ports.
	IncludePorts map[uint16]struct{}
	// ExcludePorts suppresses alerts for these ports.
	ExcludePorts map[uint16]struct{}
}

// Filter applies inclusion/exclusion rules to a Diff.
type Filter struct {
	cfg Config
}

// New creates a Filter from the given Config.
func New(cfg Config) *Filter {
	return &Filter{cfg: cfg}
}

// NewFromSlices builds a Filter from plain slices of port numbers.
func NewFromSlices(include, exclude []uint16) *Filter {
	cfg := Config{
		IncludePorts: toSet(include),
		ExcludePorts: toSet(exclude),
	}
	return New(cfg)
}

// Apply returns a new Diff containing only the ports that pass the filter.
func (f *Filter) Apply(d scanner.Diff) scanner.Diff {
	return scanner.Diff{
		Opened: f.filterPorts(d.Opened),
		Closed: f.filterPorts(d.Closed),
	}
}

func (f *Filter) filterPorts(ports []uint16) []uint16 {
	var out []uint16
	for _, p := range ports {
		if len(f.cfg.ExcludePorts) > 0 {
			if _, excluded := f.cfg.ExcludePorts[p]; excluded {
				continue
			}
		}
		if len(f.cfg.IncludePorts) > 0 {
			if _, included := f.cfg.IncludePorts[p]; !included {
				continue
			}
		}
		out = append(out, p)
	}
	return out
}

func toSet(ports []uint16) map[uint16]struct{} {
	s := make(map[uint16]struct{}, len(ports))
	for _, p := range ports {
		s[p] = struct{}{}
	}
	return s
}
