package notifier

import (
	"fmt"
	"io"
	"os"
)

// Config holds notifier configuration sourced from the main config file.
type Config struct {
	Host     string `yaml:"host"`
	Template string `yaml:"template"`
	Output   string `yaml:"output"` // "stdout", "stderr", or a file path
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Template: "",
		Output:   "stdout",
	}
}

// NewFromConfig constructs a Notifier from a Config.
func NewFromConfig(cfg Config) (*Notifier, error) {
	out, err := resolveOutput(cfg.Output)
	if err != nil {
		return nil, err
	}
	return New(cfg.Host, cfg.Template, out)
}

func resolveOutput(output string) (io.Writer, error) {
	switch output {
	case "", "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("notifier: open output file: %w", err)
		}
		return f, nil
	}
}
