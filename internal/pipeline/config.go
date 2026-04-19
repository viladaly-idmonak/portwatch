package pipeline

import (
	"errors"
	"time"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DebounceDelay:   200 * time.Millisecond,
		ThrottleEnabled: false,
		ThrottleInterval: time.Second,
		MetricsEnabled:  true,
	}
}

// Config holds pipeline-level tunables.
type Config struct {
	DebounceDelay    time.Duration
	ThrottleEnabled  bool
	ThrottleInterval time.Duration
	MetricsEnabled   bool
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	if c.DebounceDelay < 0 {
		return errors.New("pipeline: debounce delay must be non-negative")
	}
	if c.ThrottleEnabled && c.ThrottleInterval <= 0 {
		return errors.New("pipeline: throttle interval must be positive when throttle is enabled")
	}
	return nil
}
