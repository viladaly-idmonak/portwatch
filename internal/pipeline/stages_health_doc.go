// Package pipeline provides a composable processing pipeline for port diff events.
//
// The WithHealth stage integrates the health subsystem into the pipeline so that
// every non-empty diff (i.e. a scan that detected a change) is recorded via
// health.Health.RecordScan. This allows the /healthz endpoint to reflect live
// scan activity and surface staleness when the watcher stops producing events.
//
// Usage:
//
//	h := health.New(cfg)
//	p := pipeline.New(
//		pipeline.WithFilter(f),
//		pipeline.WithHealth(h),
//		pipeline.WithSnapshot(path),
//	)
package pipeline
