package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeEnricher(t *testing.T, fields map[string]string) *enricher.Enricher {
	t.Helper()
	cfg := enricher.DefaultConfig()
	cfg.Enabled = true
	cfg.Fields = fields
	e, err := enricher.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("makeEnricher: %v", err)
	}
	return e
}

func TestWithEnricherEmptyDiffPassesThrough(t *testing.T) {
	e := makeEnricher(t, map[string]string{"env": "test"})
	stage := pipeline.WithEnricher(e)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Errorf("expected empty diff, got %+v", out)
	}
}

func TestWithEnricherAnnotatesOpenedEntries(t *testing.T) {
	e := makeEnricher(t, map[string]string{"env": "prod", "region": "us-east-1"})
	stage := pipeline.WithEnricher(e)

	in := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if out.Opened[0].Meta["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", out.Opened[0].Meta["env"])
	}
	if out.Opened[0].Meta["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", out.Opened[0].Meta["region"])
	}
}

func TestWithEnricherAnnotatesClosedEntries(t *testing.T) {
	e := makeEnricher(t, map[string]string{"tier": "backend"})
	stage := pipeline.WithEnricher(e)

	in := scanner.Diff{
		Closed: []scanner.Entry{{Port: 443, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %d", len(out.Closed))
	}
	if out.Closed[0].Meta["tier"] != "backend" {
		t.Errorf("expected tier=backend, got %q", out.Closed[0].Meta["tier"])
	}
}

func TestWithEnricherIntegratesWithPipeline(t *testing.T) {
	e := makeEnricher(t, map[string]string{"host": "node-1"})
	p := pipeline.New(pipeline.WithEnricher(e))

	in := scanner.Diff{
		Opened: []scanner.Entry{{Port: 22, Proto: "tcp"}},
		Closed: []scanner.Entry{{Port: 9090, Proto: "tcp"}},
	}
	out, err := p.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if out.Opened[0].Meta["host"] != "node-1" {
		t.Errorf("opened entry missing host tag")
	}
	if out.Closed[0].Meta["host"] != "node-1" {
		t.Errorf("closed entry missing host tag")
	}
}
