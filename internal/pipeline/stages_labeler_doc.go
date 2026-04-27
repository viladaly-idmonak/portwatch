// Package pipeline provides composable processing stages for port diff events.
//
// # WithLabeler Stage
//
// WithLabeler annotates each entry in a [scanner.Diff] with a human-readable
// label derived from a [labeler.Labeler] rule set. Entries whose port/protocol
// combination matches a rule receive a "label" key in their Meta map. Entries
// with no matching rule are passed through unchanged.
//
// Usage:
//
//	l, _ := labeler.NewFromConfig(cfg)
//	pipeline.New(ctx, stages..., pipeline.WithLabeler(l))
//
// If the provided labeler is nil the stage is a no-op.
package pipeline
