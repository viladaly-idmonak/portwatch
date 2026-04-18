package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestNewValidFormats(t *testing.T) {
	for _, f := range []string{"text", "json", "TEXT", "JSON"} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNewInvalidFormat(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestFormatDiffTextOpened(t *testing.T) {
	f, _ := New("text")
	d := scanner.Diff{Opened: []scanner.Port{{Proto: "tcp", Port: 8080}}}
	out := f.FormatDiff(d, fixedTime)
	if !strings.Contains(out, "OPENED") || !strings.Contains(out, "tcp/8080") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDiffTextClosed(t *testing.T) {
	f, _ := New("text")
	d := scanner.Diff{Closed: []scanner.Port{{Proto: "udp", Port: 53}}}
	out := f.FormatDiff(d, fixedTime)
	if !strings.Contains(out, "CLOSED") || !strings.Contains(out, "udp/53") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDiffJSON(t *testing.T) {
	f, _ := New("json")
	d := scanner.Diff{
		Opened: []scanner.Port{{Proto: "tcp", Port: 443}},
		Closed: []scanner.Port{{Proto: "tcp", Port: 80}},
	}
	out := f.FormatDiff(d, fixedTime)
	if !strings.Contains(out, `"event":"opened"`) {
		t.Errorf("missing opened event in JSON output: %q", out)
	}
	if !strings.Contains(out, `"event":"closed"`) {
		t.Errorf("missing closed event in JSON output: %q", out)
	}
	if !strings.Contains(out, `"port":443`) {
		t.Errorf("missing port 443 in JSON output: %q", out)
	}
}

func TestFormatDiffEmptyIsEmpty(t *testing.T) {
	f, _ := New("text")
	out := f.FormatDiff(scanner.Diff{}, fixedTime)
	if out != "" {
		t.Errorf("expected empty output for empty diff, got %q", out)
	}
}
