package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func dedupEntry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Proto: proto}
}

func TestWithDedupPassesFirstEvent(t *testing.T) {
	stage := pipeline.WithDedup()
	d := scanner.Diff{
		Opened: []scanner.Entry{dedupEntry(80, "tcp")},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(out.Opened))
	}
}

func TestWithDedupSuppressesDuplicateOpened(t *testing.T) {
	stage := pipeline.WithDedup()
	d := scanner.Diff{Opened: []scanner.Entry{dedupEntry(80, "tcp")}}

	// first call — should pass
	if _, err := stage(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	// second call with same state — should be suppressed
	outBackground(), d)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Opened) != 0 {
		t.Fatalf("expected duplicate to be suppressed, got %d entries", len(out.n}

func TestWithDedupAllowsStateChange(t *testing.T) {
	stage := pipeline.WithDedup()
	e := dedupEntry(443, "tcp")

	// open it
	if _, err := stage(context.Background(), scanner.Diff{Opened: []scanner.Entry{e}}); err != nil {
		t.Fatal(err)
	}
	// close it — different state, must pass
	out, err := stage(context.Background(), scanner.Diff{Closed: []scanner.Entry{e}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected state change to pass, got %d closed entries", len(out.Closed))
	}
}

func TestWithDedupEmptyDiffPassesThrough(t *testing.T) {
	stage := pipeline.WithDedup()
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff to pass through unchanged")
	}
}

func TestWithDedupIntegratesWithPipeline(t *testing.T) {
	p := pipeline.New(pipeline.WithDedup())
	d := scanner.Diff{Opened: []scanner.Entry{dedupEntry(8080, "tcp")}}

	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened via pipeline, got %d", len(out.Opened))
	}
}
