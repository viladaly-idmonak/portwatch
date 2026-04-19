package summarizer

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Summary holds aggregated port activity over a time window.
type Summary struct {
	From    time.Time
	To      time.Time
	Opened  []scanner.Port
	Closed  []scanner.Port
}

func (s Summary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Summary [%s -> %s]\n", s.From.Format(time.RFC3339), s.To.Format(time.RFC3339))
	fmt.Fprintf(&b, "  Opened (%d):", len(s.Opened))
	for _, p := range s.Opened {
		fmt.Fprintf(&b, " %d/%s", p.Number, p.Protocol)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "  Closed (%d):", len(s.Closed))
	for _, p := range s.Closed {
		fmt.Fprintf(&b, " %d/%s", p.Number, p.Protocol)
	}
	b.WriteString("\n")
	return b.String()
}

// Summarizer accumulates diffs and produces summaries.
type Summarizer struct {
	start   time.Time
	opened  map[scanner.Port]struct{}
	closed  map[scanner.Port]struct{}
}

// New returns a new Summarizer anchored at now.
func New() *Summarizer {
	return &Summarizer{
		start:  time.Now(),
		opened: make(map[scanner.Port]struct{}),
		closed: make(map[scanner.Port]struct{}),
	}
}

// Record accumulates a diff into the running summary.
func (s *Summarizer) Record(d scanner.Diff) {
	for _, p := range d.Opened {
		s.opened[p] = struct{}{}
		delete(s.closed, p)
	}
	for _, p := range d.Closed {
		s.closed[p] = struct{}{}
		delete(s.opened, p)
	}
}

// Flush returns the current summary and resets state.
func (s *Summarizer) Flush() Summary {
	now := time.Now()
	sum := Summary{
		From: s.start,
		To:   now,
	}
	for p := range s.opened {
		sum.Opened = append(sum.Opened, p)
	}
	for p := range s.closed {
		sum.Closed = append(sum.Closed, p)
	}
	s.start = now
	s.opened = make(map[scanner.Port]struct{})
	s.closed = make(map[scanner.Port]struct{})
	return sum
}
