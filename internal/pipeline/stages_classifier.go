package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/classifier"
	"github.com/user/portwatch/internal/scanner"
)

const metaKeyClass = "class"

// WithClassifier returns a Stage that annotates each entry in the diff with
// its classification (e.g. "trusted", "suspect", "critical", "unknown").
// The result is stored in entry.Meta["class"]. If clf is nil the stage is a
// no-op passthrough.
func WithClassifier(clf *classifier.Classifier) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if clf == nil || (len(d.Opened) == 0 && len(d.Closed) == 0) {
			return d, nil
		}
		out := scanner.Diff{
			Opened: annotateClass(clf, d.Opened),
			Closed: annotateClass(clf, d.Closed),
		}
		return out, nil
	}
}

func annotateClass(clf *classifier.Classifier, entries []scanner.Entry) []scanner.Entry {
	result := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		if e.Meta == nil {
			e.Meta = make(map[string]string)
		} else {
			meta := make(map[string]string, len(e.Meta))
			for k, v := range e.Meta {
				meta[k] = v
			}
			e.Meta = meta
		}
		if _, exists := e.Meta[metaKeyClass]; !exists {
			e.Meta[metaKeyClass] = string(clf.Classify(e))
		}
		result[i] = e
	}
	return result
}
