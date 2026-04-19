package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func baselineDiff(opened, closed []uint16) scanner.Diff {
	d := scanner.Diff{}
	for _, p := range opened {
		d.Opened = append(d.Opened, scanner.Port{Number: p, Protocol: "tcp"})
	}
	for _, p := range closed {
		d.Closed = append(d.Closed, scanner.Port{Number: p, Protocol: "tcp"})
	}
	return d
}

func TestWithBaselineEmptyDiffPassesThrough(t *testing.T) {
	b := baseline.New(nil)
	stage := pipeline.WithBaseline(b)
	d := scanner.Diff{}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Errorf("expected empty diff, got %+v", out)
	}
}

func TestWithBaselineFiltersExpectedPorts(t *testing.T) {
	entries := []baseline.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	b := baseline.New(entries)
	stage := pipeline.WithBaseline(b)

	// port 80 is expected, port 9000 is unexpected
	d := baselineDiff([]uint16{80, 9000}, nil)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 || out.Opened[0].Number != 9000 {
		t.Errorf("expected only port 9000 as deviation, got %+v", out.Opened)
	}
}

func TestWithBaselineIntegratesWithPipeline(t *testing.T) {
	entries := []baseline.Entry{
		{Port: 22, Protocol: "tcp"},
	}
	b := baseline.New(entries)
	p := pipeline.New(pipeline.WithBaseline(b))

	d := baselineDiff([]uint16{22, 8080}, nil)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 || out.Opened[0].Number != 8080 {
		t.Errorf("expected port 8080 deviation only, got %+v", out.Opened)
	}
}
