package pipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeFingerprint(t *testing.T) *fingerprint.Hasher {
	t.Helper()
	h, err := fingerprint.New("testsalt")
	if err != nil {
		t.Fatalf("fingerprint.New: %v", err)
	}
	return h
}

func TestWithFingerprintEmptyDiffPassesThrough(t *testing.T) {
	h := makeFingerprint(t)
	stage := pipeline.WithFingerprint(h)
	diff := scanner.Diff{}
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatalf("expected empty diff, got %+v", out)
	}
}

func TestWithFingerprintAnnotatesOpenedEntries(t *testing.T) {
	h := makeFingerprint(t)
	stage := pipeline.WithFingerprint(h)
	diff := scanner.Diff{
		Opened: []scanner.Entry{{Port: 443, Proto: "tcp"}},
	}
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if out.Opened[0].Meta == nil {
		t.Fatal("expected Meta to be set after fingerprinting")
	}
	if _, ok := out.Opened[0].Meta["fingerprint"]; !ok {
		t.Fatal("expected 'fingerprint' key in Meta")
	}
}

func TestWithFingerprintAnnotatesClosedEntries(t *testing.T) {
	h := makeFingerprint(t)
	stage := pipeline.WithFingerprint(h)
	diff := scanner.Diff{
		Closed: []scanner.Entry{{Port: 80, Proto: "tcp"}},
	}
	out, err := stage(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %d", len(out.Closed))
	}
	if out.Closed[0].Meta == nil {
		t.Fatal("expected Meta to be set after fingerprinting")
	}
	if _, ok := out.Closed[0].Meta["fingerprint"]; !ok {
		t.Fatal("expected 'fingerprint' key in Meta")
	}
}

func TestWithFingerprintIntegratesWithPipeline(t *testing.T) {
	h := makeFingerprint(t)
	p := pipeline.New(
		pipeline.WithFingerprint(h),
	)
	diff := scanner.Diff{
		Opened: []scanner.Entry{{Port: 8443, Proto: "tcp"}},
	}
	out, err := p.Run(t.Context(), diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if _, ok := out.Opened[0].Meta["fingerprint"]; !ok {
		t.Fatal("expected fingerprint meta key")
	}
}
