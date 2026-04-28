// Package pipeline provides a composable stage-based processing pipeline
// for port diff events.
//
// # WithClassifier
//
// WithClassifier annotates each entry in a [scanner.Diff] with a "class"
// metadata field derived from a [classifier.Classifier]. The class is
// determined by matching the entry's port and protocol against a set of
// configured rules; unmatched entries receive the classifier's default class.
//
// If the provided classifier is nil the stage is a no-op and the diff passes
// through unchanged.
//
// Example:
//
//	c, _ := classifier.New("unknown", []classifier.Rule{
//		{Port: 80, Protocol: "tcp", Class: "web"},
//	})
//	p := pipeline.New(pipeline.WithClassifier(c))
//	out, _ := p.Run(ctx, diff)
package pipeline
