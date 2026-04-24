// Package pipeline provides composable processing stages for port diff events.
//
// WithCircuitBreaker wraps an inner pipeline stage function with a circuit
// breaker that trips after a configurable number of consecutive failures.
//
// When the circuit is closed (healthy), calls are forwarded to the inner
// function normally. After MaxFailures consecutive errors the circuit opens
// and subsequent calls return an error immediately without invoking the inner
// function, protecting downstream systems from cascading failures.
//
// After the configured ResetTimeout the circuit transitions to half-open,
// allowing a single probe call through. A success closes the circuit; a
// failure reopens it.
//
// Empty diffs are short-circuited before reaching the inner function, keeping
// the failure counter unaffected by no-op scans.
//
// Example usage:
//
//	cb, _ := circuitbreaker.NewFromConfig(circuitbreaker.DefaultConfig())
//	stage := pipeline.WithCircuitBreaker(cb, myDownstreamStage)
//	p := pipeline.New(stage)
package pipeline
