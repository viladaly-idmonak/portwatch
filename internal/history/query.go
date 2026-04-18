package history

import "time"

// Filter holds optional criteria for querying history entries.
type Filter struct {
	Since *time.Time
	Until *time.Time
	Port  int    // 0 means any
	Proto string // "" means any
	State string // "" means any
}

// Query returns entries matching all non-zero filter fields.
func (h *History) Query(f Filter) []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	var out []Entry
	for _, e := range h.entries {
		if f.Since != nil && e.Timestamp.Before(*f.Since) {
			continue
		}
		if f.Until != nil && e.Timestamp.After(*f.Until) {
			continue
		}
		if f.Port != 0 && e.Port != f.Port {
			continue
		}
		if f.Proto != "" && e.Proto != f.Proto {
			continue
		}
		if f.State != "" && e.State != f.State {
			continue
		}
		out = append(out, e)
	}
	return out
}
