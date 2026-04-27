// Package pipeline provides a composable stage-based processing pipeline
// for port activity diffs.
//
// # WithScorecard Stage
//
// WithScorecard annotates each entry in a diff with a risk score level
// based on configurable port/protocol rules. It uses the scorecard package
// to assign levels such as "low", "medium", or "high" to opened and closed
// port entries, storing the result in the entry's Meta map under the key
// "score".
//
// Example usage:
//
//	sc := scorecard.New(scorecard.DefaultLevel, rules)
//	pipeline.New(ctx, diff, WithScorecard(sc))
//
// If the scorecard is nil, the stage is a no-op and the diff passes through
// unchanged.
package pipeline
