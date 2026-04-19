// Package pipeline provides a composable stage-based processing pipeline
// for port scan diffs.
//
// WithBaseline wraps a [baseline.Baseline] as a pipeline stage. It filters
// the incoming [scanner.Diff] so that only deviations from the established
// baseline are forwarded downstream.
//
// Ports listed in the baseline are considered "expected" and are stripped
// from the Opened set. Ports that disappear but were expected are stripped
// from the Closed set. Only genuinely unexpected changes survive.
//
// Usage:
//
//	b := baseline.NewFromScan(previousScan)
//	p := pipeline.New(
//		pipeline.WithBaseline(b),
//		pipeline.WithNotifier(n),
//	)
package pipeline
