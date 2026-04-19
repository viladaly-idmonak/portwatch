package pipeline

import (
	"fmt"
	"time"
)

// Config holds optional configuration for building a pipeline from flags/env.
type Config struct {
	// DebounceWait is how long to wait for the diff stream to settle.
	// Zero means debouncing is disabled.
	DebounceWait time.Duration `yaml:"debounce_wait"`

	// ThrottleEnabled controls whether the throttle stage is added.
	ThrottleEnabled bool `yaml:"throttle_enabled"`

	// ThrottleInterval is the minimum duration between allowed diffs.
	ThrottleInterval time.Duration `yaml:"throttle_interval"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DebounceWait:     0,
		ThrottleEnabled:  false,
		ThrottleInterval: 500 * time.Millisecond,
	}
}

// Validate checks the config for invalid combinations.
func (c Config) Validate() error {
	if c.DebounceWait < 0 {
		return fmt.Errorf("pipeline: debounce_wait must be non-negative")
	}
	if c.ThrottleEnabled && c.ThrottleInterval <= 0 {
		return fmt.Errorf("pipeline: throttle_interval must be positive when throttle is enabled")
	}
	return nil
}
