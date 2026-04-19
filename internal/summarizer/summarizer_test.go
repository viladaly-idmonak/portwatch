package summarizer_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summarizer"
)

func port(n uint16, proto string) scanner.Port {
	return scanner.Port{Number: n, Protocol: proto}
}

func diff(opened, closed []scanner.Port) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func TestFlushEmptyReturnsZeroCounts(t *testing.T) {
	s := summarizer.New()
	sum := s.Flush()
	if len(sum.Opened) != 0 || len(sum.Closed) != 0 {
		t.Fatalf("expected empty summary, got %+v", sum)
	}
}

func TestRecordAccumulatesOpened(t *testing.T) {
	s := summarizer.New()
	s.Record(diff([]scanner.Port{port(80, "tcp"), port(443, "tcp")}, nil))
	sum := s.Flush()
	if len(sum.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(sum.Opened))
	}
}

func TestRecordAccumulatesClosed(t *testing.T) {
	s := summarizer.New()
	s.Record(diff(nil, []scanner.Port{port(22, "tcp")}))
	sum := s.Flush()
	if len(sum.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(sum.Closed))
	}
}

func TestOpenThenCloseDeduplicates(t *testing.T) {
	s := summarizer.New()
	s.Record(diff([]scanner.Port{port(8080, "tcp")}, nil))
	s.Record(diff(nil, []scanner.Port{port(8080, "tcp")}))
	sum := s.Flush()
	if len(sum.Opened) != 0 {
		t.Fatalf("expected opened to be removed after close, got %d", len(sum.Opened))
	}
	if len(sum.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(sum.Closed))
	}
}

func TestFlushResetsState(t *testing.T) {
	s := summarizer.New()
	s.Record(diff([]scanner.Port{port(9090, "tcp")}, nil))
	s.Flush()
	sum := s.Flush()
	if len(sum.Opened) != 0 {
		t.Fatalf("expected empty after second flush, got %d opened", len(sum.Opened))
	}
}

func TestStringContainsPortNumbers(t *testing.T) {
	s := summarizer.New()
	s.Record(diff([]scanner.Port{port(80, "tcp")}, []scanner.Port{port(22, "tcp")}))
	sum := s.Flush()
	str := sum.String()
	if !strings.Contains(str, "80") || !strings.Contains(str, "22") {
		t.Fatalf("expected port numbers in string output: %s", str)
	}
}
