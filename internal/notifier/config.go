package notifier

import (
	"errors"
	"io"
	"os"
)

// Config holds configuration for the Notifier.
type Config struct {
	// Template is the Go text/template string used to format messages.
	Template string
	// Output is the writer to send notifications to. Defaults to os.Stdout.
	Output io.Writer
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Template: defaultTemplate,
		Output:   os.Stdout,
	}
}

// Validate checks the Config for invalid values.
func (c Config) Validate() error {
	if c.Template == "" {
		return errors.New("notifier: template must not be empty")
	}
	if c.Output == nil {
		return errors.New("notifier: output writer must not be nil")
	}
	return nil
}

// NewFromConfig creates a Notifier from the given Config.
func NewFromConfig(cfg Config) (*Notifier, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Output, cfg.Template)
}

// resolveOutput returns w if non-nil, otherwise os.Stdout.
func resolveOutput(w io.Writer) io.Writer {
	if w != nil {
		return w
	}
	return os.Stdout
}
