package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeNotifier(t *testing.T, buf *bytes.Buffer) *notifier.Notifier {
	t.Helper()
	cfg := notifier.DefaultConfig()
	cfg.Output = buf
	n, err := notifier.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("notifier.NewFromConfig: %v", err)
	}
	return n
}

func TestWithNotifierEmptyDiffNoOutput(t *testing.T) {
	var buf bytes.Buffer
	n := makeNotifier(t, &buf)
	stage := pipeline.WithNotifier(n)
	_, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestWithNotifierOpenedWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	n := makeNotifier(t, &buf)
	stage := pipeline.WithNotifier(n)
	d := scanner.Diff{Opened: []scanner.Port{{Proto: "tcp", Port: 8080}}}
	_, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Errorf("expected port 8080 in output, got %q", buf.String())
	}
}

func TestWithNotifierIntegratesWithPipeline(t *testing.T) {
	var buf bytes.Buffer
	n := makeNotifier(t, &buf)
	p := pipeline.New(pipeline.WithNotifier(n))
	d := scanner.Diff{Closed: []scanner.Port{{Proto: "udp", Port: 53}}}
	_, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "53") {
		t.Errorf("expected port 53 in output, got %q", buf.String())
	}
}
