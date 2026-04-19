// Package pipeline provides composable processing stages for port diff events.
//
// WithLogger stage
//
// WithLogger wraps a [logger.Logger] as a pipeline stage. It logs every
// non-empty diff and passes the diff through unmodified, allowing downstream
// stages to continue processing.
//
// Example usage:
//
//	p := pipeline.New(
//		pipeline.WithFilter(f),
//		pipeline.WithLogger(log),
//		pipeline.WithSnapshot(path),
//	)
package pipeline
