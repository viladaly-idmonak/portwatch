package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/anomaly"
	"github.com/user/portwatch/internal/scanner"
)

const metaAnomalous = "anomalous"

// WithAnomaly returns a Stage that records each opened or closed entry with
// the Detector and annotates it with meta key "anomalous" = "true" when the
// event frequency within the configured window exceeds the threshold.
//
// A nil detector is a no-op: the diff passes through unchanged.
func WithAnomaly(d *anomaly.Detector) Stage {
	if d == nil {
		return func(_ context.Context, diff scanner.Diff) (scanner.Diff, error) {
			return diff, nil
		}
	}
	return func(_ context.Context, diff scanner.Diff) (scanner.Diff, error) {
		if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
			return diff, nil
		}
		out := scanner.Diff{
			Opened: annotateAnomaly(d, diff.Opened),
			Closed: annotateAnomaly(d, diff.Closed),
		}
		return out, nil
	}
}

func annotateAnomaly(d *anomaly.Detector, entries []scanner.Entry) []scanner.Entry {
	result := make([]scanner.Entry, len(entries))
	for i, e := range entries {
		if d.Record(e) {
			if e.Meta == nil {
				e.Meta = make(map[string]string)
			} else {
				m := make(map[string]string, len(e.Meta)+1)
				for k, v := range e.Meta {
					m[k] = v
				}
				e.Meta = m
			}
			e.Meta[metaAnomalous] = "true"
		}
		result[i] = e
	}
	return result
}
