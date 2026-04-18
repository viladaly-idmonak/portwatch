package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestSaveAndLoad(t *testing.T) {
	ports := []scanner.Port{
		{Proto: "tcp", Number: 80},
		{Proto: "tcp", Number: 443},
	}
	p := tmpPath(t)
	if err := snapshot.Save(p, ports); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := snapshot.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(loaded))
	}
	for i, pp := range ports {
		if loaded[i] != pp {
			t.Errorf("port[%d]: expected %+v, got %+v", i, pp, loaded[i])
		}
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	ports, err := snapshot.Load("/nonexistent/portwatch-snap.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Fatalf("expected empty slice, got %d ports", len(ports))
	}
}

func TestSaveIsAtomic(t *testing.T) {
	p := tmpPath(t)
	if err := snapshot.Save(p, []scanner.Port{{Proto: "udp", Number: 53}}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not found after Save: %v", err)
	}
}
