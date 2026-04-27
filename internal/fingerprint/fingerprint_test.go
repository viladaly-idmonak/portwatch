package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func entry(port uint16, proto, state string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto, State: state}
}

func TestComputeReturnsSameHashForSameInput(t *testing.T) {
	h := fingerprint.New("salt")
	a := h.Compute(entry(80, "tcp", "open"))
	b := h.Compute(entry(80, "tcp", "open"))
	if a.Hash != b.Hash {
		t.Fatalf("expected identical hashes, got %s and %s", a.Hash, b.Hash)
	}
}

func TestComputeDiffersForDifferentPorts(t *testing.T) {
	h := fingerprint.New("")
	a := h.Compute(entry(80, "tcp", "open"))
	b := h.Compute(entry(443, "tcp", "open"))
	if a.Hash == b.Hash {
		t.Fatal("expected different hashes for different ports")
	}
}

func TestComputeDiffersForDifferentSalts(t *testing.T) {
	a := fingerprint.New("alpha").Compute(entry(80, "tcp", "open"))
	b := fingerprint.New("beta").Compute(entry(80, "tcp", "open"))
	if a.Hash == b.Hash {
		t.Fatal("expected different hashes for different salts")
	}
}

func TestComputeAllReturnsOnePerEntry(t *testing.T) {
	h := fingerprint.New("")
	entries := []scanner.Entry{
		entry(22, "tcp", "open"),
		entry(80, "tcp", "open"),
		entry(443, "tcp", "open"),
	}
	fps := h.ComputeAll(entries)
	if len(fps) != len(entries) {
		t.Fatalf("expected %d fingerprints, got %d", len(entries), len(fps))
	}
}

func TestEqualReturnsTrueForMatchingHash(t *testing.T) {
	h := fingerprint.New("")
	a := h.Compute(entry(8080, "tcp", "open"))
	b := h.Compute(entry(8080, "tcp", "open"))
	if !fingerprint.Equal(a, b) {
		t.Fatal("expected Equal to return true")
	}
}

func TestEqualReturnsFalseForMismatch(t *testing.T) {
	h := fingerprint.New("")
	a := h.Compute(entry(8080, "tcp", "open"))
	b := h.Compute(entry(8080, "udp", "open"))
	if fingerprint.Equal(a, b) {
		t.Fatal("expected Equal to return false")
	}
}
