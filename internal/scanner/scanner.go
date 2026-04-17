package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	SeenAt   time.Time
}

// Snapshot holds all open ports at a point in time.
type Snapshot map[int]PortState

// Scanner scans a range of TCP ports on localhost.
type Scanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
}

// New creates a Scanner with sensible defaults.
func New(start, end int) *Scanner {
	return &Scanner{
		StartPort: start,
		EndPort:   end,
		Timeout:   500 * time.Millisecond,
	}
}

// Scan probes each port in the configured range and returns a Snapshot.
func (s *Scanner) Scan() Snapshot {
	snap := make(Snapshot)
	for port := s.StartPort; port <= s.EndPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err == nil {
			conn.Close()
			snap[port] = PortState{
				Port:     port,
				Protocol: "tcp",
				Open:     true,
				SeenAt:   time.Now(),
			}
		}
	}
	return snap
}

// Diff compares two snapshots and returns newly opened and closed ports.
func Diff(prev, curr Snapshot) (opened, closed []PortState) {
	for port, state := range curr {
		if _, existed := prev[port]; !existed {
			opened = append(opened, state)
		}
	}
	for port, state := range prev {
		if _, exists := curr[port]; !exists {
			closed = append(closed, state)
		}
	}
	return
}
