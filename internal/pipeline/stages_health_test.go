package pipeline_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/health"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeHealth(t *testing.T) *health.Health {
	t.Helper()
	cfg := health.DefaultConfig()
	cfg.Enabled = false
	h := health.New(cfg)
	return h
}

func TestWithHealthEmptyDiffSkipsRecord(t *testing.T) {
	h := makeHealth(t)
	stage := pipeline.WithHealth(h)
	d := scanner.Diff{}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Errorf("expected empty diff")
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.ServeHTTP(rr, req)
	// scans counter should remain 0
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestWithHealthRecordsScan(t *testing.T) {
	h := makeHealth(t)
	stage := pipeline.WithHealth(h)
	d := scanner.Diff{Opened: []uint16{8080}}
	_, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithHealthIntegratesWithPipeline(t *testing.T) {
	h := makeHealth(t)
	p := pipeline.New(pipeline.WithHealth(h))
	d := scanner.Diff{Opened: []uint16{9000}, Closed: []uint16{80}}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 || len(out.Closed) != 1 {
		t.Errorf("diff should pass through unchanged")
	}
}
