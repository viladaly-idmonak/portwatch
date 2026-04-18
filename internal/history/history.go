package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a single port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	State     string    `json:"state"` // "opened" or "closed"
}

// History maintains an in-memory ring buffer of recent events and persists them.
type History struct {
	mu      sync.Mutex
	entries []Entry
	max     int
	path    string
}

// New creates a History with a max capacity, loading any existing entries from path.
func New(path string, max int) (*History, error) {
	h := &History{path: path, max: max}
	if err := h.load(); err != nil {
		return nil, err
	}
	return h, nil
}

// Record appends opened/closed events from a Diff to the history.
func (h *History) Record(diff scanner.Diff) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now().UTC()
	for _, p := range diff.Opened {
		h.append(Entry{Timestamp: now, Port: p.Port, Proto: p.Proto, State: "opened"})
	}
	for _, p := range diff.Closed {
		h.append(Entry{Timestamp: now, Port: p.Port, Proto: p.Proto, State: "closed"})
	}
	return h.save()
}

// Entries returns a copy of current entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) append(e Entry) {
	if len(h.entries) >= h.max {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

func (h *History) save() error {
	if h.path == "" {
		return nil
	}
	data, err := json.Marshal(h.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}

func (h *History) load() error {
	if h.path == "" {
		return nil
	}
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}
