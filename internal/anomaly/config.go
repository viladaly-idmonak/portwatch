package anomaly

import (
	"fmt"
	"time"
)

// Config holds configuration for the anomaly Detector.
type Config struct {
	Enabled   bool
	Window    time.Duration
	Threshold int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:   false,
		Window:    time.Minute,
		Threshold: 5,
	}
}

// Validate returns an error if the Config contains invalid values.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return fmt.Errorf("anomaly: window must be positive")
	}
	if c.Threshold <= 0 {
		return fmt.Errorf("anomaly: threshold must be positive")
	}
	return nil
}

// NewFromConfig constructs a Detector from cfg, returning nil when disabled.
func NewFromConfig(cfg Config) (*Detector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	return New(cfg.Window, cfg.Threshold)
}
