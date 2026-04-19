package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func TestWithMetricsRecordsOpenedAndClosed(t *testing.T) {
	m := metrics.New()
	stage := pipeline.WithMetrics(m)

	d := scanner.Diff{
		Opened: []scanner.Port{{Number: 80}, {Number: 443}},
		Closed: []scanner.Port{{Number: 8080}},
	}

	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 2 || len(out.Closed) != 1 {
		t.Fatal("diff should pass through unchanged")
	}

	snap := m.Snapshot()
	if snap.TotalOpened != 2 {
		t.Errorf("expected TotalOpened=2, got %d", snap.TotalOpened)
	}
	if snap.TotalClosed != 1 {
		t.Errorf("expected TotalClosed=1, got %d", snap.TotalClosed)
	}
}

func TestWithMetricsSkipsEmptyDiff(t *testing.T) {
	m := metrics.New()
	stage := pipeline.WithMetrics(m)

	_, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := m.Snapshot()
	if snap.TotalOpened != 0 || snap.TotalClosed != 0 {
		t.Error("expected no metrics recorded for empty diff")
	}
}

func TestWithMetricsIntegratesWithPipeline(t *testing.T) {
	m := metrics.New()
	p := pipeline.New(pipeline.WithMetrics(m))

	d := openDiff(9000)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatal("expected diff to pass through")
	}
	if m.Snapshot().TotalOpened != 1 {
		t.Error("expected TotalOpened=1")
	}
}
