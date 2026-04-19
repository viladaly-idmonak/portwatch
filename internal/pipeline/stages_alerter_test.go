package pipeline_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeAlerter(t *testing.T, srv *httptest.Server) *alerter.Alerter {
	t.Helper()
	a, err := alerter.New(alerter.Config{
		OnOpen:  srv.URL,
		OnClose: srv.URL,
	})
	if err != nil {
		t.Fatalf("alerter.New: %v", err)
	}
	return a
}

func TestWithAlerterEmptyDiffNoRequest(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	a := makeAlerter(t, srv)
	stage := pipeline.WithAlerter(a)
	_, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty diff")
	}
}

func TestWithAlerterOpenedSendsRequest(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := makeAlerter(t, srv)
	stage := pipeline.WithAlerter(a)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}}}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected HTTP request to be made")
	}
	if len(out.Opened) != 1 {
		t.Errorf("expected diff to pass through, got %+v", out)
	}
}

func TestWithAlerterIntegratesWithPipeline(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := makeAlerter(t, srv)
	p := pipeline.New(pipeline.WithAlerter(a))
	d := scanner.Diff{Closed: []scanner.Entry{{Port: 443, Proto: "tcp"}}}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Errorf("expected closed entry to pass through")
	}
	if hits == 0 {
		t.Error("expected alerter to fire")
	}
}
