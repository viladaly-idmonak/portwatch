// Package pipeline provides a composable processing pipeline for port diff events.
//
// The WithHealth stage integrates the health subsystem into the pipeline so that
// every non-empty diff (i.e. a scan that detected a change) is recorded via
// health.Health.RecordScan. This allows the /healthz endpoint to reflect live
// scan activity and surface staleness when the watcher stops producing events.
//
// Stage Ordering
//
// Stages are applied in the order they are passed to pipeline.New. The recommended
// ordering is:
//
//  1. WithFilter   – drop unwanted ports before any further processing
//  2. WithHealth   – record scan activity after filtering so that only meaningful
//                    changes update the health timestamp
//  3. WithSnapshot – persist the final, filtered state to disk
//
// Placing WithHealth before WithSnapshot ensures that a successful write to disk
// is not required for the health check to be updated, keeping the two concerns
// independent.
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
