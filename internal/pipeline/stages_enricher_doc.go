// Package pipeline provides a composable stage-based processing pipeline
// for port activity diffs.
//
// WithEnricher attaches static or dynamic metadata fields to each entry in a
// scanner.Diff before passing it downstream. Fields are sourced from the
// enricher.Enricher configuration and will not overwrite keys that already
// exist in an entry's Meta map.
//
// Typical usage:
//
//	pipeline.New(
//		pipeline.WithEnricher(e),
//		pipeline.WithLogger(l),
//	)
package pipeline
