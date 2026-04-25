// Package pipeline provides composable processing stages for port diff events.
//
// # WithRedactor
//
// WithRedactor wraps a [redactor.Redactor] as a pipeline stage. It scrubs
// sensitive metadata values from every entry in the diff before the event
// continues downstream.
//
// Configured keys are matched case-insensitively. Matched values are replaced
// with the literal string "[REDACTED]". The original diff is never mutated;
// a deep copy is returned.
//
// Example usage:
//
//	r, _ := redactor.New([]string{"token", "apikey", "password"})
//	p := pipeline.New(
//		pipeline.WithRedactor(r),
//		pipeline.WithLogger(log),
//	)
package pipeline
