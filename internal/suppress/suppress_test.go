package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/suppress"
)

func entry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Proto: proto}
}

func diff(opened, closed []scanner.Entry) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := suppress.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = suppress.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestFirstEventAlwaysPasses(t *testing.T) {
	s, _ := suppress.New(time.Minute)
	d := diff([]scanner.Entry{entry(80, "tcp")}, nil)
	out := s.Apply(d)
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(out.Opened))
	}
}

func TestDuplicateWithinWindowSuppressed(t *testing.T) {
	s, _ := suppress.New(time.Minute)
	d := diff([]scanner.Entry{entry(80, "tcp")}, nil)
	s.Apply(d)
	out := s.Apply(d)
	if len(out.Opened) != 0 {
		t.Fatalf("expected 0 opened, got %d", len(out.Opened))
	}
}

func TestEventPassesAfterWindowExpires(t *testing.T) {
	s, _ := suppress.New(50 * time.Millisecond)
	d := diff([]scanner.Entry{entry(443, "tcp")}, nil)
	s.Apply(d)
	time.Sleep(60 * time.Millisecond)
	out := s.Apply(d)
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened after window, got %d", len(out.Opened))
	}
}

func TestOpenedAndClosedTrackedSeparately(t *testing.T) {
	s, _ := suppress.New(time.Minute)
	e := entry(22, "tcp")
	s.Apply(diff([]scanner.Entry{e}, nil)) // seed opened
	out := s.Apply(diff(nil, []scanner.Entry{e}))
	if len(out.Closed) != 1 {
		t.Fatalf("expected closed to pass through, got %d", len(out.Closed))
	}
}

func TestPurgeRemovesExpiredEntries(t *testing.T) {
	s, _ := suppress.New(30 * time.Millisecond)
	d := diff([]scanner.Entry{entry(8080, "tcp")}, nil)
	s.Apply(d)
	time.Sleep(40 * time.Millisecond)
	s.Purge()
	// After purge, same event should be allowed through again.
	out := s.Apply(d)
	if len(out.Opened) != 1 {
		t.Fatalf("expected event after purge, got %d", len(out.Opened))
	}
}
