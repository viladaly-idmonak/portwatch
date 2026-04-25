// Package pipeline provides a composable, stage-based processing pipeline
// for port activity diffs.
//
// # Watchlist Stage
//
// WithWatchlist promotes ports in the watchlist that appear in the closed
// set to the opened set, ensuring watched ports are always surfaced as
// active events regardless of transient close/reopen cycles.
//
// Ports not present in the watchlist pass through unchanged. This stage
// is intended to be placed early in the pipeline, before filtering or
// rate-limiting stages.
//
// Example usage:
//
//	wl := watchlist.New()
//	_ = wl.Add(80, "tcp")
//	_ = wl.Add(443, "tcp")
//	p := pipeline.New()
//	p.Use(pipeline.WithWatchlist(wl))
package pipeline
