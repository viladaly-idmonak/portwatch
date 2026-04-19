package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeLogger(buf *bytes.Buffer) *logger.Logger {
	return logger.New(buf)
}

func TestWithLoggerEmptyDiffNoOutput(t *testing.T) {
	var buf bytes.Buffer
	l := makeLogger(&buf)
	stage := pipeline.WithLogger(l)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Error("diff should be unchanged")
	}
}

func TestWithLoggerOpenedWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	l := makeLogger(&buf)
	stage := pipeline.WithLogger(l)

	d := openDiff(9090)
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "9090") {
		t.Errorf("expected port 9090 in log output, got: %s", buf.String())
	}
	if len(out.Opened) != 1 {
		t.Error("diff should pass through unchanged")
	}
}

func TestWithLoggerIntegratesWithPipeline(t *testing.T) {
	var buf bytes.Buffer
	l := makeLogger(&buf)

	p := pipeline.New(pipeline.WithLogger(l))
	d := openDiff(4040)
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Error("expected diff to pass through pipeline")
	}
	if !strings.Contains(buf.String(), "4040") {
		t.Errorf("expected port 4040 in output, got: %s", buf.String())
	}
}
