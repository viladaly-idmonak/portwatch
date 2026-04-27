package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

// WithFingerprint is a pipeline stage that annotates each entry in the diff
// with a stable content hash produced by the provided Hasher. The hash is
// stored in the entry's Meta map under the key "fingerprint".
func WithFingerprint(h *fingerprint.Hasher) Stage {
	return func(ctx context.Context, d scanner.Diff) (scanner.Diff, error) {
		if len(d.Opened) == 0 && len(d.Closed) == 0 {
			return d, nil
		}

		annotate := func(entries []scanner.Entry) []scanner.Entry {
			out := make([]scanner.Entry, len(entries))
			for i, e := range entries {
				fp := h.Compute(e)
				if e.Meta == nil {
					e.Meta = make(map[string]string)
				}
				e.Meta["fingerprint"] = fp.Hash
				out[i] = e
			}
			return out
		}

		return scanner.Diff{
			Opened: annotate(d.Opened),
			Closed: annotate(d.Closed),
		}, nil
	}
}
