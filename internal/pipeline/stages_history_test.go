package pipeline_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func historyPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestWithHistoryRecordsOpened(t *testing.T) {
	h := history.New(historyPath(t), 100)
	stage := pipeline.WithHistory(h)

	d := openDiff(9090)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected diff to pass through unchanged")
	}

	entries := h.Query(history.QueryOptions{})
	if len(entries) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(entries))
	}
}

func TestWithHistorySkipsEmptyDiff(t *testing.T) {
	h := history.New(historyPath(t), 100)
	stage := pipeline.WithHistory(h)

	_, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := h.Query(history.QueryOptions{})
	if len(entries) != 0 {
		t.Fatalf("expected no history entries for empty diff, got %d", len(entries))
	}
}

func TestWithHistoryIntegratesWithPipeline(t *testing.T) {
	h := history.New(historyPath(t), 100)
	p := pipeline.New(pipeline.WithHistory(h))

	for _, port := range []uint16{8080, 8081} {
		_, err := p.Run(context.Background(), openDiff(port))
		if err != nil {
			t.Fatalf("pipeline error: %v", err)
		}
	}

	entries := h.Query(history.QueryOptions{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(entries))
	}
	_ = os.Remove(historyPath(t))
}
