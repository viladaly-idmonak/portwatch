package audit

import "fmt"

// Config holds configuration for the Auditor.
type Config struct {
	Enabled  bool   `toml:"enabled"`
	FilePath string `toml:"file_path"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:  false,
		FilePath: "portwatch-audit.jsonl",
	}
}

// Validate checks that the config is consistent.
func (c Config) Validate() error {
	if c.Enabled && c.FilePath == "" {
		return fmt.Errorf("audit: file_path must not be empty when enabled")
	}
	return nil
}

// NewFromConfig constructs an Auditor from config, or returns nil when disabled.
func NewFromConfig(c Config) (*Auditor, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if !c.Enabled {
		return nil, nil
	}
	return NewToFile(c.FilePath)
}
