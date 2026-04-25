// Package pipeline provides composable pipeline stages for processing port diffs.
//
// # WithTruncator
//
// WithTruncator limits the number of entries in the Opened and Closed slices of
// a [scanner.Diff] to a configurable maximum. This prevents downstream stages
// from being overwhelmed during large bursts of port activity — for example,
// when a host restarts and hundreds of ports appear simultaneously.
//
// Entries are truncated from the tail; the first max entries are preserved so
// that higher-priority (lower-numbered) ports are more likely to be retained
// when combined with a sort stage.
//
// Example usage:
//
//	p := pipeline.New()
//	p.Use(pipeline.WithTruncator(pipeline.NewTruncator(50)))
//
// A zero or negative max disables truncation and passes all entries through.
package pipeline
