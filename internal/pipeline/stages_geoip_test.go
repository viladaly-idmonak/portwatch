package pipeline_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/geoip"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

type stubLookup struct {
	info geoip.Info
	err  error
}

func (s *stubLookup) Lookup(_ string) (geoip.Info, error) {
	return s.info, s.err
}

func TestWithGeoIPAnnotatesOpenedEntry(t *testing.T) {
	lookup := &stubLookup{info: geoip.Info{Country: "DE", City: "Berlin", Org: "Hetzner"}}
	stage := pipeline.WithGeoIP(lookup)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 80, Addr: "5.9.0.1"}}}

	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry")
	}
	if !strings.Contains(out.Opened[0].Meta, "DE") {
		t.Errorf("expected meta to contain country, got %q", out.Opened[0].Meta)
	}
}

func TestWithGeoIPLookupErrorIsNonFatal(t *testing.T) {
	lookup := &stubLookup{err: errors.New("timeout")}
	stage := pipeline.WithGeoIP(lookup)
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 443, Addr: "1.2.3.4"}}}

	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("stage should not return error on lookup failure")
	}
	if !strings.Contains(out.Opened[0].Meta, "geoip_error") {
		t.Errorf("expected geoip_error in meta, got %q", out.Opened[0].Meta)
	}
}

func TestWithGeoIPEmptyDiffPassesThrough(t *testing.T) {
	lookup := &stubLookup{}
	stage := pipeline.WithGeoIP(lookup)
	d := scanner.Diff{}

	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 0 || len(out.Closed) != 0 {
		t.Error("expected empty diff passthrough")
	}
}

func TestWithGeoIPIntegratesWithPipeline(t *testing.T) {
	lookup := &stubLookup{info: geoip.Info{Country: "FR"}}
	p := pipeline.New(pipeline.WithGeoIP(lookup))
	d := scanner.Diff{Opened: []scanner.Entry{{Port: 22, Addr: "185.1.2.3"}}}

	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if !strings.Contains(out.Opened[0].Meta, "FR") {
		t.Errorf("expected FR in meta, got %q", out.Opened[0].Meta)
	}
}
