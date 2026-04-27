// Package pipeline provides the processing pipeline for port diff events.
//
// # Timeout Stage
//
// WithTimeout wraps any downstream Stage with a per-invocation deadline.
// If the wrapped stage does not return within the configured duration the
// original [scanner.Diff] is forwarded unchanged and a wrapped
// [context.DeadlineExceeded] error is returned, allowing the caller to
// decide whether to halt the pipeline or continue.
//
// Usage:
//
//	stage, err := pipeline.WithTimeout(2*time.Second, myExpensiveStage)
//	if err != nil {
//		log.Fatal(err)
//	}
//	p := pipeline.New(stage)
//
// Configuration-driven construction is available via [NewTimeouterFromConfig].
// When [TimeoutConfig.Enabled] is false the function returns nil so callers
// can skip wiring the stage entirely.
package pipeline
