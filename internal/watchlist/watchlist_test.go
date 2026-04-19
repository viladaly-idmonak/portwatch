package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func TestAddAndContains(t *testing.T) {
	w := watchlist.New()
	e := watchlist.Entry{Port: 8080, Protocol: "tcp", Note: "http-alt"}
	if err := w.Add(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Contains("tcp", 8080) {
		t.Error("expected port 8080/tcp to be in watchlist")
	}
}

func TestAddInvalidProtocolReturnsError(t *testing.T) {
	w := watchlist.New()
	err := w.Add(watchlist.Entry{Port: 80, Protocol: "icmp"})
	if err == nil {
		t.Error("expected error for invalid protocol")
	}
}

func TestRemoveDeletesEntry(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(watchlist.Entry{Port: 443, Protocol: "tcp"})
	w.Remove("tcp", 443)
	if w.Contains("tcp", 443) {
		t.Error("expected port 443/tcp to be removed")
	}
}

func TestRemoveNonExistentIsNoop(t *testing.T) {
	w := watchlist.New()
	w.Remove("tcp", 9999) // should not panic
}

func TestAllReturnsCopy(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(watchlist.Entry{Port: 22, Protocol: "tcp", Note: "ssh"})
	_ = w.Add(watchlist.Entry{Port: 53, Protocol: "udp", Note: "dns"})
	all := w.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the slice should not affect the watchlist.
	all[0] = watchlist.Entry{Port: 1, Protocol: "tcp"}
	if w.Len() != 2 {
		t.Error("watchlist length changed after mutating snapshot")
	}
}

func TestLenReflectsCount(t *testing.T) {
	w := watchlist.New()
	if w.Len() != 0 {
		t.Error("expected empty watchlist")
	}
	_ = w.Add(watchlist.Entry{Port: 80, Protocol: "tcp"})
	if w.Len() != 1 {
		t.Errorf("expected 1, got %d", w.Len())
	}
}
