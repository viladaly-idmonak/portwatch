package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func dedupEntry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto}
}

func TestWithDedupPassesFirstEvent(t *testing.T) {
	d := pipeline.NewDedup(5*time.Second, 100)
	stage := pipeline.WithDedup(d)

	in := scanner.Diff{Opened: []scanner.Entry{dedupEntry(80, "tcp")}}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Errorf("expected 1 opened, got %d", len(out.Opened))
	}
}

func TestWithDedupSuppressesDuplicateOpened(t *testing.T) {
	d := pipeline.NewDedup(5*time.Second, 100)
	stage := pipeline.WithDedup(d)
	ctx := context.Background()

	in := scanner.Diff{Opened: []scanner.Entry{dedupEntry(80, "tcp")}}
	_, _ = stage(ctx, in)

	out, err := stage(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 {
		t.Errorf("expected duplicate to be suppressed, got %d opened", len(out.Opened))
	}
}

func TestWithDedupAllowsStateChange(t *testing.T) {
	d := pipeline.NewDedup(5*time.Second, 100)
	stage := pipeline.WithDedup(d)
	ctx := context.Background()

	opened := scanner.Diff{Opened: []scanner.Entry{dedupEntry(443, "tcp")}}
	_, _ = stage(ctx, opened)

	closed := scanner.Diff{Closed: []scanner.Entry{dedupEntry(443, "tcp")}}
	out, err := stage(ctx, closed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Errorf("expected state change to pass through, got %d closed", len(out.Closed))
	}
}

func TestWithDedupEmptyDiffPassesThrough(t *testing.T) {
	d := pipeline.NewDedup(5*time.Second, 100)
	stage := pipeline.WithDedup(d)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Error("expected empty diff to pass through unchanged")
	}
}

func TestWithDedupIntegratesWithPipeline(t *testing.T) {
	d := pipeline.NewDedup(5*time.Second, 100)
	p := pipeline.New(pipeline.WithDedup(d))
	ctx := context.Background()

	in := scanner.Diff{Opened: []scanner.Entry{dedupEntry(8080, "tcp")}}
	out, err := p.Run(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Errorf("expected 1 opened on first run, got %d", len(out.Opened))
	}

	out, err = p.Run(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error on second run: %v", err)
	}
	if len(out.Opened) != 0 {
		t.Errorf("expected duplicate suppressed on second run, got %d", len(out.Opened))
	}
}
