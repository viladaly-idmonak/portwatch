package sampling_test

import (
	"testing"

	"github.com/user/portwatch/internal/sampling"
	"github.com/user/portwatch/internal/scanner"
)

func openDiff(ports ...uint16) scanner.Diff {
	d := scanner.Diff{}
	for _, p := range ports {
		d.Opened = append(d.Opened, scanner.Entry{Port: p, Proto: "tcp"})
	}
	return d
}

func TestNewInvalidRateReturnsError(t *testing.T) {
	_, err := sampling.New(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
	_, err = sampling.New(1.5)
	if err == nil {
		t.Fatal("expected error for rate=1.5")
	}
}

func TestNewValidRateNoError(t *testing.T) {
	s, err := sampling.New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestRateOnePassesAll(t *testing.T) {
	s, _ := sampling.New(1.0)
	d := openDiff(80, 443, 8080)
	out := s.Apply(d)
	if len(out.Opened) != 3 {
		t.Fatalf("expected 3 opened, got %d", len(out.Opened))
	}
}

func TestApplyEmptyDiffReturnsEmpty(t *testing.T) {
	s, _ := sampling.New(1.0)
	out := s.Apply(scanner.Diff{})
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestSetRateUpdatesRate(t *testing.T) {
	s, _ := sampling.New(0.5)
	if err := s.SetRate(0.9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.9 {
		t.Fatalf("expected rate 0.9, got %v", s.Rate())
	}
}

func TestSetRateInvalidReturnsError(t *testing.T) {
	s, _ := sampling.New(0.5)
	if err := s.SetRate(-0.1); err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestSampleStatisticallyCorrect(t *testing.T) {
	// With rate=1.0 every sample must pass.
	s, _ := sampling.New(1.0)
	for i := 0; i < 100; i++ {
		if !s.Sample() {
			t.Fatal("rate=1.0 should always sample")
		}
	}
}
