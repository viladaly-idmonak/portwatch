package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func makeWatchlist(entries ...watchlist.Entry) *watchlist.Watchlist {
	wl := watchlist.New()
	for _, e := range entries {
		_ = wl.Add(e)
	}
	return wl
}

func TestWithWatchlistEmptyDiffPassesThrough(t *testing.T) {
	wl := makeWatchlist(watchlist.Entry{Port: 22, Protocol: "tcp"})
	stage := pipeline.WithWatchlist(wl)
	out, err := stage(context.Background(), scanner.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Error("expected empty diff")
	}
}

func TestWithWatchlistMovesClosedWatchedPortToOpened(t *testing.T) {
	wl := makeWatchlist(watchlist.Entry{Port: 22, Protocol: "tcp"})
	stage := pipeline.WithWatchlist(wl)
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 22, Protocol: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 0 {
		t.Errorf("expected closed to be empty, got %d entries", len(out.Closed))
	}
	if len(out.Opened) != 1 || out.Opened[0].Port != 22 {
		t.Error("expected watched port to be moved to opened")
	}
}

func TestWithWatchlistUnwatchedClosedPortUnchanged(t *testing.T) {
	wl := makeWatchlist(watchlist.Entry{Port: 22, Protocol: "tcp"})
	stage := pipeline.WithWatchlist(wl)
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 9000, Protocol: "tcp"}},
	}
	out, err := stage(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Closed) != 1 {
		t.Errorf("expected 1 closed entry, got %d", len(out.Closed))
	}
}

func TestWithWatchlistIntegratesWithPipeline(t *testing.T) {
	wl := makeWatchlist(watchlist.Entry{Port: 443, Protocol: "tcp"})
	p := pipeline.New(pipeline.WithWatchlist(wl))
	d := scanner.Diff{
		Closed: []scanner.Entry{{Port: 443, Protocol: "tcp"}},
	}
	out, err := p.Run(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 || out.Opened[0].Port != 443 {
		t.Error("pipeline: expected watched port in opened")
	}
}
