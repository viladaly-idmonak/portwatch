package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"` // "opened" | "closed"
	Port      uint16    `json:"port"`
	Proto     string    `json:"proto"`
}

// Auditor writes structured audit records to a destination.
type Auditor struct {
	w   io.Writer
	enc *json.Encoder
}

// New returns an Auditor writing to w.
func New(w io.Writer) *Auditor {
	enc := json.NewEncoder(w)
	return &Auditor{w: w, enc: enc}
}

// NewToFile opens (or creates) a file and returns an Auditor writing to it.
func NewToFile(path string) (*Auditor, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open file: %w", err)
	}
	return New(f), nil
}

// Record writes audit entries for every port change in diff.
func (a *Auditor) Record(diff scanner.Diff) error {
	now := time.Now().UTC()
	for _, p := range diff.Opened {
		if err := a.enc.Encode(Entry{Timestamp: now, Event: "opened", Port: p.Port, Proto: p.Proto}); err != nil {
			return err
		}
	}
	for _, p := range diff.Closed {
		if err := a.enc.Encode(Entry{Timestamp: now, Event: "closed", Port: p.Port, Proto: p.Proto}); err != nil {
			return err
		}
	}
	return nil
}
