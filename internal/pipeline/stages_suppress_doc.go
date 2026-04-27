// Package pipeline provides a composable stage-based processing pipeline
// for port activity diffs.
//
// # WithSuppress Stage
//
// WithSuppress wraps a [suppress.Suppressor] as a pipeline stage. It drops
// repeated port events (same port, same state) that occur within a
// configurable time window, reducing noise from flapping ports.
//
// Usage:
//
//	s, _ := suppress.New(suppress.DefaultConfig())
//	pipeline.New(ctx, WithSuppress(s))
//
// Entries that fall within the suppression window are silently removed from
// the diff. If the resulting diff is empty the stage returns early without
// forwarding to subsequent stages.
package pipeline
