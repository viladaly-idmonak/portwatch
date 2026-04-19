package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/scanner"
)

func diff(opened, closed []scanner.Port) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func port(p uint16, proto string) scanner.Port {
	return scanner.Port{Port: p, Proto: proto}
}

func TestRecordOpened(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	if err := a.Record(diff([]scanner.Port{port(80, "tcp")}, nil)); err != nil {
		t.Fatal(err)
	}
	var e audit.Entry
	if err := json.NewDecoder(&buf).Decode(&e); err != nil {
		t.Fatal(err)
	}
	if e.Event != "opened" || e.Port != 80 || e.Proto != "tcp" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestRecordClosed(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	if err := a.Record(diff(nil, []scanner.Port{port(443, "tcp")})); err != nil {
		t.Fatal(err)
	}
	var e audit.Entry
	if err := json.NewDecoder(&buf).Decode(&e); err != nil {
		t.Fatal(err)
	}
	if e.Event != "closed" || e.Port != 443 {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestRecordEmptyDiffWritesNothing(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	if err := a.Record(diff(nil, nil)); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty buffer, got %d bytes", buf.Len())
	}
}

func TestRecordMultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	_ = a.Record(diff([]scanner.Port{port(22, "tcp"), port(8080, "tcp")}, []scanner.Port{port(3000, "tcp")}))
	dec := json.NewDecoder(&buf)
	count := 0
	for dec.More() {
		var e audit.Entry
		if err := dec.Decode(&e); err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 3 {
		t.Fatalf("expected 3 entries, got %d", count)
	}
}
