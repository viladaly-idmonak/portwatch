package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func diff(opened, closed []scanner.PortEntry) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func port(p int, proto string) scanner.PortEntry {
	return scanner.PortEntry{Port: p, Proto: proto}
}

func TestRecordStoresEntries(t *testing.T) {
	h, err := New("", 100)
	if err != nil {
		t.Fatal(err)
	}
	_ = h.Record(diff([]scanner.PortEntry{port(8080, "tcp")}, nil))
	entries := h.Entries()
	if len(entries) != 1 || entries[0].Port != 8080 || entries[0].State != "opened" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestRingBufferEvictsOldest(t *testing.T) {
	h, err := New("", 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 4; i++ {
		_ = h.Record(diff([]scanner.PortEntry{port(i, "tcp")}, nil))
	}
	entries := h.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Port != 2 {
		t.Fatalf("expected oldest evicted, got port %d", entries[0].Port)
	}
}

func TestPersistAndReload(t *testing.T) {
	p := tmpPath(t)
	h, _ := New(p, 100)
	_ = h.Record(diff([]scanner.PortEntry{port(443, "tcp")}, nil))

	h2, err := New(p, 100)
	if err != nil {
		t.Fatal(err)
	}
	entries := h2.Entries()
	if len(entries) != 1 || entries[0].Port != 443 {
		t.loaded entries: %+v", entries)
	}
}

func TestMissingFileReturnsEmpty(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
, err := New(p, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Entries()) != 0 {
		t.Fatal("expected empty history")
	}
}

func TestSaveIsValidJSON(t *testing.T) {
	p := tmpPath(t)
	h, _ := New(p, 100)
	_ = h.Record(diff(nil, []scanner.PortEntry{port(22, "tcp")}))
	data, _ := os.ReadFile(p)
	var out []Entry
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
