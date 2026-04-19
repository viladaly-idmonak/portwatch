package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/resolver"
	"github.com/user/portwatch/internal/scanner"
)

func TestWithResolverPassesDiffUnchanged(t *testing.T) {
	r := resolver.New(true)
	stage := pipeline.WithResolver(r, "tcp")
	d := scanner.Diff{Opened: []uint16{22, 80}, Closed: []uint16{8080}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 2 || len(out.Closed) != 1 {
		t.Fatalf("diff was modified unexpectedly: %+v", out)
	}
}

func TestWithResolverEmptyDiffNoError(t *testing.T) {
	r := resolver.New(false)
	stage := pipeline.WithResolver(r, "tcp")
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestWithResolverIntegratesWithPipeline(t *testing.T) {
	r := resolver.New(true)
	p := pipeline.New(pipeline.WithResolver(r, "tcp"))
	d := scanner.Diff{Opened: []uint16{443}}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(out.Opened))
	}
}
