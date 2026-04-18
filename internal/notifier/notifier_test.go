package notifier

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestNotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	n, err := New("localhost", "", &buf)
	if err != nil {
		t.Fatal(err)
	}
	diff := scanner.Diff{
		Opened: []scanner.Port{{Number: 8080, Proto: "tcp"}},
	}
	if err := n.Notify(diff); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "OPEN") || !strings.Contains(buf.String(), "8080") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestNotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	n, _ := New("host1", "", &buf)
	diff := scanner.Diff{
		Closed: []scanner.Port{{Number: 22, Proto: "tcp"}},
	}
	_ = n.Notify(diff)
	if !strings.Contains(buf.String(), "CLOSE") || !strings.Contains(buf.String(), "22") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestNotifyEmptyDiffNoOutput(t *testing.T) {
	var buf bytes.Buffer
	n, _ := New("host1", "", &buf)
	_ = n.Notify(scanner.Diff{})
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestNotifyCustomTemplate(t *testing.T) {
	var buf bytes.Buffer
	n, err := New("myhost", "{{.State}}|{{.Port}}|{{.Host}}\n", &buf)
	if err != nil {
		t.Fatal(err)
	}
	_ = n.Notify(scanner.Diff{Opened: []scanner.Port{{Number: 443, Proto: "tcp"}}})
	if buf.String() != "OPEN|443|myhost\n" {
		t.Errorf("unexpected: %q", buf.String())
	}
}

func TestNewInvalidTemplate(t *testing.T) {
	_, err := New("h", "{{.Unclosed", nil)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}
