// Package pipeline provides composable processing stages for port diff events.
//
// # WithEnricher Stage
//
// WithEnricher attaches static key-value metadata fields to every entry in a
// [scanner.Diff] using an [enricher.Enricher]. Fields are written into each
// entry's Meta map, allowing downstream stages (formatters, alerters, audit
// logs) to access host-level context such as environment, region, or cluster
// without requiring each stage to perform its own lookup.
//
// Existing Meta keys are NOT overwritten; the enricher only fills in keys that
// are absent, preserving any values set by earlier pipeline stages.
//
// Example usage:
//
//	e, _ := enricher.NewFromConfig(cfg)
//	p := pipeline.New(
//		pipeline.WithEnricher(e),
//		pipeline.WithFormatter(f),
//	)
package pipeline
