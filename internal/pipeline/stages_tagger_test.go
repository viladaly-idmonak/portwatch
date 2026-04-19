package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func makeTagger() *tagger.Tagger {
	return tagger.New(map[uint16]string{80: "http", 443: "https"})
}

func TestWithTaggerEmptyDiffPassesThrough(t *testing.T) {
	tr := makeTagger()
	stage := pipeline.WithTagger(tr)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestWithTaggerAnnotatesKnownPort(t *testing.T) {
	tr := makeTagger()
	stage := pipeline.WithTagger(tr)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 80, Proto: "tcp"}}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	if out.Opened[0].Tag != "http" {
		t.Fatalf("expected http, got %s", out.Opened[0].Tag)
	}
}

func TestWithTaggerFallsBackForUnknownPort(t *testing.T) {
	tr := makeTagger()
	stage := pipeline.WithTagger(tr)
	d := scanner.Diff{Closed: []scanner.Entry{{Port: 9999, Proto: "tcp"}}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	if out.Closed[0].Tag != "port/9999" {
		t.Fatalf("unexpected tag: %s", out.Closed[0].Tag)
	}
}

func TestWithTaggerIntegratesWithPipeline(t *testing.T) {
	tr := tagger.New(map[uint16]string{22: "ssh"})
	p := pipeline.New(pipeline.WithTagger(tr))
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 22, Proto: "tcp"}}}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	if out.Opened[0].Tag != "ssh" {
		t.Fatalf("expected ssh tag, got %s", out.Opened[0].Tag)
	}
}
