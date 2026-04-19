package pipeline_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func snapshotPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestWithSnapshotSavesOpenedPorts(t *testing.T) {
	p := snapshotPath(t)
	stage := pipeline.WithSnapshot(p)
	d := openDiff(8080, 9090)
	_, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ports, err := snapshot.Load(p)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
}

func TestWithSnapshotRemovesClosedPorts(t *testing.T) {
	p := snapshotPath(t)
	// seed with two ports
	_ = snapshot.Save(p, []uint16{8080, 9090})

	stage := pipeline.WithSnapshot(p)
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
	_, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ports, _ := snapshot.Load(p)
	if len(ports) != 1 || ports[0] != 9090 {
		t.Fatalf("expected [9090], got %v", ports)
	}
}

func TestWithSnapshotSkipsEmptyDiff(t *testing.T) {
	p := snapshotPath(t)
	stage := pipeline.WithSnapshot(p)
	_, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// file should not be created
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("expected no snapshot file for empty diff")
	}
}

func TestWithSnapshotIntegratesWithPipeline(t *testing.T) {
	p := snapshotPath(t)
	pl := pipeline.New(pipeline.WithSnapshot(p))
	d := openDiff(3000, 4000, 5000)
	_, err := pl.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	ports, _ := snapshot.Load(p)
	sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })
	expected := []uint16{3000, 4000, 5000}
	for i, want := range expected {
		if ports[i] != want {
			t.Errorf("port[%d]: want %d, got %d", i, want, ports[i])
		}
	}
}
