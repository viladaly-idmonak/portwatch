// Package pipeline provides a composable, stage-based processing pipeline
// for port activity diffs.
//
// # Dedup Stage
//
// WithDedup suppresses duplicate port state events within a configurable
// time window. It tracks the last-seen state for each (port, protocol) pair
// and drops events that repeat the same state before the window expires.
//
// This is useful when the underlying scanner produces repeated open/closed
// events for the same port due to transient noise or polling overlap.
//
// Example usage:
//
//	dedup, err := pipeline.NewDedupFromConfig(pipeline.DefaultDedupConfig())
//	if err != nil {
//		log.Fatal(err)
//	}
//	p := pipeline.New()
//	p.Use(pipeline.WithDedup(dedup))
package pipeline
