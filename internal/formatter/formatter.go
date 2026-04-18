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

// FormatDiff returns a formatted string representing the diff at the given time.
func (f *Formatter) FormatDiff(d scanner.Diff, t time.Time) string {
	ts := t.Format(time.RFC3339)
	var sb strings.Builder
	for _, p := range d.Opened {
		if f.format == FormatJSON {
			sb.WriteString(fmt.Sprintf(`{"time":%q,"event":"opened","proto":%q,"port":%d}`+"\n", ts, p.Proto, p.Port))
		} else {
			sb.WriteString(fmt.Sprintf("%s OPENED %s/%d\n", ts, p.Proto, p.Port))
		}
	}
	for _, p := range d.Closed {
		if f.format == FormatJSON {
			sb.WriteString(fmt.Sprintf(`{"time":%q,"event":"closed","proto":%q,"port":%d}`+"\n", ts, p.Proto, p.Port))
		} else {
			sb.WriteString(fmt.Sprintf("%s CLOSED %s/%d\n", ts, p.Proto, p.Port))
		}
	}
	return sb.String()
}
