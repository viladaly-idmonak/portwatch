package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestLogDiffOpened(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	diff := scanner.Diff{
		Opened: []scanner.PortEntry{{Port: 8080, Protocol: "tcp"}},
	}
	l.LogDiff(diff)

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO level, got: %s", out)
	}
	if !strings.Contains(out, "port opened: 8080/tcp") {
		t.Errorf("expected opened message, got: %s", out)
	}
}

func TestLogDiffClosed(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	diff := scanner.Diff{
		Closed: []scanner.PortEntry{{Port: 443, Protocol: "tcp"}},
	}
	l.LogDiff(diff)

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level, got: %s", out)
	}
	if !strings.Contains(out, "port closed: 443/tcp") {
		t.Errorf("expected closed message, got: %s", out)
	}
}

func TestLogDiffNoChange(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	l.LogDiff(scanner.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	l := New(nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}
