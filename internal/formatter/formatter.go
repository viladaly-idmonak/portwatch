package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls the output format style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter formats port diff events for output.
type Formatter struct {
	format Format
}

// New returns a Formatter for the given format string.
func New(format string) (*Formatter, error) {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatText, FormatJSON:
		return &Formatter{format: f}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %q (want \"text\" or \"json\")", format)
	}
}

// formatEntry returns a single formatted line for a port event.
func (f *Formatter) formatEntry(ts, event string, p scanner.Port) string {
	if f.format == FormatJSON {
		return fmt.Sprintf(`{"time":%q,"event":%q,"proto":%q,"port":%d}`+"\n", ts, event, p.Proto, p.Port)
	}
	return fmt.Sprintf("%s %s %s/%d\n", ts, strings.ToUpper(event), p.Proto, p.Port)
}

// FormatDiff returns a formatted string representing the diff at the given time.
func (f *Formatter) FormatDiff(d scanner.Diff, t time.Time) string {
	ts := t.Format(time.RFC3339)
	var sb strings.Builder
	for _, p := range d.Opened {
		sb.WriteString(f.formatEntry(ts, "opened", p))
	}
	for _, p := range d.Closed {
		sb.WriteString(f.formatEntry(ts, "closed", p))
	}
	return sb.String()
}
