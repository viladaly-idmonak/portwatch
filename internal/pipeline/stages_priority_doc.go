// Package pipeline provides a composable stage-based processing pipeline
// for port diff events.
//
// # WithPriority Stage
//
// WithPriority annotates each entry in a [scanner.Diff] with a numeric
// priority score stored in the entry's Meta map under the key "priority".
//
// Priority rules are defined as a mapping of port numbers to integer scores.
// Ports not present in the rules map receive the default priority of 0.
//
// Example usage:
//
//	rules := map[uint16]int{
//		22:  100, // SSH — high priority
//		80:  50,  // HTTP — medium priority
//		443: 50,  // HTTPS — medium priority
//	}
//
//	p, err := pipeline.NewPrioritizer(rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	pl := pipeline.New(
//		pipeline.WithPriority(p),
//	)
//
// Passing nil as the prioritizer is safe and results in a no-op stage that
// passes the diff through unchanged.
package pipeline
