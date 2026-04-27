package history

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func seedHistory(t *testing.T) *History {
	t.Helper()
	h, _ := New("", 100)
	_ = h.Record(diff([]scanner.PortEntry{port(80, "tcp"), port(443, "tcp")}, nil))
	_ = h.Record(diff(nil, []scanner.PortEntry{port(80, "tcp")}))
	return h
}

func TestQueryNoFilterReturnsAll(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{})
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
}

func TestQueryByPort(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Port: 80})
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestQueryByState(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{State: "closed"})
	if len(results) != 1 || results[0].Port != 80 {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestQueryBySince(t *testing.T) {
	h := seedHistory(t)
	future := time.Now().UTC().Add(time.Hour)
	results := h.Query(Filter{Since: &future})
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

func TestQueryByProto(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Proto: "tcp"})
	if len(results) != 3 {
		t.Fatalf("expected 3, got %d", len(results))
	}
}

func TestQueryByPortAndState(t *testing.T) {
	h := seedHistory(t)
	results := h.Query(Filter{Port: 80, State: "open"})
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Port != 80 || results[0].State != "open" {
		t.Fatalf("unexpected result: %+v", results[0])
	}
}
