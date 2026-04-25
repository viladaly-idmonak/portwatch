package redactor_test

import (
	"testing"

	"github.com/user/portwatch/internal/redactor"
	"github.com/user/portwatch/internal/scanner"
)

func entry(meta map[string]string) scanner.Entry {
	return scanner.Entry{Port: 80, Proto: "tcp", Meta: meta}
}

func TestNoKeysPassesDiffUnchanged(t *testing.T) {
	r := redactor.New(nil)
	d := scanner.Diff{
		Opened: []scanner.Entry{entry(map[string]string{"token": "secret"})},
	}
	out := r.Apply(d)
	if out.Opened[0].Meta["token"] != "secret" {
		t.Fatalf("expected value to be unchanged, got %q", out.Opened[0].Meta["token"])
	}
}

func TestRedactsMaskedKey(t *testing.T) {
	r := redactor.New([]string{"token"})
	d := scanner.Diff{
		Opened: []scanner.Entry{entry(map[string]string{"token": "abc123", "host": "localhost"})},
	}
	out := r.Apply(d)
	if out.Opened[0].Meta["token"] != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", out.Opened[0].Meta["token"])
	}
	if out.Opened[0].Meta["host"] != "localhost" {
		t.Fatalf("expected host to be unchanged, got %q", out.Opened[0].Meta["host"])
	}
}

func TestRedactionIsCaseInsensitive(t *testing.T) {
	r := redactor.New([]string{"API_KEY"})
	d := scanner.Diff{
		Opened: []scanner.Entry{entry(map[string]string{"api_key": "hunter2"})},
	}
	out := r.Apply(d)
	if out.Opened[0].Meta["api_key"] != "[REDACTED]" {
		t.Fatalf("expected redaction, got %q", out.Opened[0].Meta["api_key"])
	}
}

func TestOriginalDiffNotMutated(t *testing.T) {
	r := redactor.New([]string{"secret"})
	orig := scanner.Diff{
		Opened: []scanner.Entry{entry(map[string]string{"secret": "value"})},
	}
	_ = r.Apply(orig)
	if orig.Opened[0].Meta["secret"] != "value" {
		t.Fatal("original diff was mutated")
	}
}

func TestRedactsClosedEntries(t *testing.T) {
	r := redactor.New([]string{"pw"})
	d := scanner.Diff{
		Closed: []scanner.Entry{entry(map[string]string{"pw": "pass", "user": "admin"})},
	}
	out := r.Apply(d)
	if out.Closed[0].Meta["pw"] != "[REDACTED]" {
		t.Fatalf("expected closed entry to be redacted, got %q", out.Closed[0].Meta["pw"])
	}
	if out.Closed[0].Meta["user"] != "admin" {
		t.Fatalf("expected user to be unchanged, got %q", out.Closed[0].Meta["user"])
	}
}

func TestEmptyDiffPassesThrough(t *testing.T) {
	r := redactor.New([]string{"secret"})
	out := r.Apply(scanner.Diff{})
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Fatal("expected empty diff to pass through unchanged")
	}
}
