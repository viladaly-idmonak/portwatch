package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/redactor"
	"github.com/user/portwatch/internal/scanner"
)

func makeRedactor(t *testing.T, keys ...string) *redactor.Redactor {
	t.Helper()
	r, err := redactor.New(keys)
	if err != nil {
		t.Fatalf("redactor.New: %v", err)
	}
	return r
}

func TestWithRedactorEmptyDiffPassesThrough(t *testing.T) {
	r := makeRedactor(t, "secret")
	stage := pipeline.WithRedactor(r)

	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Errorf("expected empty diff, got %+v", out)
	}
}

func TestWithRedactorMasksMetaKey(t *testing.T) {
	r := makeRedactor(t, "token")
	stage := pipeline.WithRedactor(r)

	d := scanner.Diff{
		Opened: []scanner.Entry{
			{Port: 8080, Proto: "tcp", Meta: map[string]string{"token": "supersecret", "env": "prod"}},
		},
	}

	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
	if out.Opened[0].Meta["token"] != "[REDACTED]" {
		t.Errorf("expected token to be redacted, got %q", out.Opened[0].Meta["token"])
	}
	if out.Opened[0].Meta["env"] != "prod" {
		t.Errorf("expected env to be unchanged, got %q", out.Opened[0].Meta["env"])
	}
}

func TestWithRedactorIntegratesWithPipeline(t *testing.T) {
	r := makeRedactor(t, "apikey")

	p := pipeline.New(
		pipeline.WithRedactor(r),
	)

	d := scanner.Diff{
		Opened: []scanner.Entry{
			{Port: 443, Proto: "tcp", Meta: map[string]string{"apikey": "abc123"}},
		},
	}

	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline.Run: %v", err)
	}
	if out.Opened[0].Meta["apikey"] != "[REDACTED]" {
		t.Errorf("expected apikey redacted in pipeline, got %q", out.Opened[0].Meta["apikey"])
	}
}
