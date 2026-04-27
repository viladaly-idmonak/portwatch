// Package pipeline provides composable stages for processing port diff events.
//
// # WithFingerprint Stage
//
// WithFingerprint annotates each entry in a [scanner.Diff] with a stable
// cryptographic fingerprint derived from the entry's port, protocol, and an
// optional salt. The fingerprint is stored in the entry's Meta map under the
// key "fingerprint".
//
// This is useful for correlating events across multiple runs, deduplicating
// alerts in downstream systems, or building audit trails that survive restarts.
//
// Example usage:
//
//	hasher, _ := fingerprint.New("my-secret-salt")
//	p := pipeline.New(
//		pipeline.WithFingerprint(hasher),
//	)
package pipeline
