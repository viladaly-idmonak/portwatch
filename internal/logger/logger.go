package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the log severity level.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Logger writes port activity events to an output writer.
type Logger struct {
	out io.Writer
}

// New creates a Logger that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// LogDiff writes a human-readable log line for each change in the diff.
func (l *Logger) LogDiff(diff scanner.Diff) {
	for _, p := range diff.Opened {
		l.write(LevelInfo, fmt.Sprintf("port opened: %d/%s", p.Port, p.Protocol))
	}
	for _, p := range diff.Closed {
		l.write(LevelWarn, fmt.Sprintf("port closed: %d/%s", p.Port, p.Protocol))
	}
}

func (l *Logger) write(level Level, msg string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(l.out, "%s [%s] %s\n", timestamp, level, msg)
}
