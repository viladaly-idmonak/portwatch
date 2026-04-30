package pipeline_test

import (
	"context"
	"testing"

	"github.com/example/portwatch/internal/pipeline"
	"github.com/example/portwatch/internal/scanner"
)

func makePrioritizer(rules map[uint16]int) *pipeline.Prioritizer {
	p, err := pipeline.NewPrioritizer(rules)
	if err != nil {
		panic(err)
	}
	return p
}

func TestWithPriorityEmptyDiffPassesThrough(t *testing.T) {
	p := makePrioritizer(nil)
	stage := pipeline.WithPriority(p)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithPriorityAnnotatesOpenedEntry(t *testing.T) {
	rules := map[uint16]int{80: 10, 443: 20}
	p := makePrioritizer(rules)
	stage := pipeline.WithPriority(p)

	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	prio, ok := out.Opened[0].Meta["priority"]
	if !ok {
		t.Fatal("expected 'priority' key in meta")
	}
	if prio != "10" {
		t.Fatalf("expected priority 10, got %q", prio)
	}
}

func TestWithPriorityAnnotatesClosedEntry(t *testing.T) {
	rules := map[uint16]int{443: 20}
	p := makePrioritizer(rules)
	stage := pipeline.WithPriority(p)

	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 443, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %d", len(out.Closed))
	}
	prio, ok := out.Closed[0].Meta["priority"]
	if !ok {
		t.Fatal("expected 'priority' key in meta")
	}
	if prio != "20" {
		t.Fatalf("expected priority 20, got %q", prio)
	}
}

func TestWithPriorityUnknownPortUsesDefault(t *testing.T) {
	p := makePrioritizer(map[uint16]int{80: 5})
	stage := pipeline.WithPriority(p)

	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 9999, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prio, ok := out.Opened[0].Meta["priority"]
	if !ok {
		t.Fatal("expected 'priority' key in meta")
	}
	if prio != "0" {
		t.Fatalf("expected default priority 0, got %q", prio)
	}
}

func TestWithPriorityNilIsNoop(t *testing.T) {
	stage := pipeline.WithPriority(nil)
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 80, Proto: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Opened[0].Meta["priority"]; ok {
		t.Fatal("expected no priority annotation when prioritizer is nil")
	}
}

func TestWithPriorityIntegratesWithPipeline(t *testing.T) {
	rules := map[uint16]int{22: 100}
	p := makePrioritizer(rules)

	pl := pipeline.New(pipeline.WithPriority(p))
	d := scanner.Diff{
		Opened: []scanner.Entry{{Port: 22, Proto: "tcp"}},
	}
	out, err := pl.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Opened[0].Meta["priority"] != "100" {
		t.Fatalf("expected priority 100, got %q", out.Opened[0].Meta["priority"])
	}
}
