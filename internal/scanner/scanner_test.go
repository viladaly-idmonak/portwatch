package scanner

import (
	"net"
	"testing"
	"time"
)

// startListener opens a TCP listener on an OS-assigned port and returns it.
func startListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestScanDetectsOpenPort(t *testing.T) {
	ln, port := startListener(t)
	defer ln.Close()

	s := &Scanner{StartPort: port, EndPort: port, Timeout: 200 * time.Millisecond}
	snap := s.Scan()

	if _, ok := snap[port]; !ok {
		t.Errorf("expected port %d to be detected as open", port)
	}
}

func TestScanMissesClosedPort(t *testing.T) {
	ln, port := startListener(t)
	ln.Close() // close immediately

	s := &Scanner{StartPort: port, EndPort: port, Timeout: 200 * time.Millisecond}
	snap := s.Scan()

	if _, ok := snap[port]; ok {
		t.Errorf("expected port %d to be absent from snapshot", port)
	}
}

func TestDiffOpenedAndClosed(t *testing.T) {
	now := time.Now()
	prev := Snapshot{
		8080: {Port: 8080, Protocol: "tcp", Open: true, SeenAt: now},
	}
	curr := Snapshot{
		9090: {Port: 9090, Protocol: "tcp", Open: true, SeenAt: now},
	}

	opened, closed := Diff(prev, curr)

	if len(opened) != 1 || opened[0].Port != 9090 {
		t.Errorf("expected port 9090 in opened, got %v", opened)
	}
	if len(closed) != 1 || closed[0].Port != 8080 {
		t.Errorf("expected port 8080 in closed, got %v", closed)
	}
}

func TestDiffNoChange(t *testing.T) {
	now := time.Now()
	snap := Snapshot{
		443: {Port: 443, Protocol: "tcp", Open: true, SeenAt: now},
	}
	opened, closed := Diff(snap, snap)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no diff for identical snapshots")
	}
}
