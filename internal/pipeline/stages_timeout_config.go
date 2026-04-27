package pipeline

import (
	"fmt"
	"time"
)

// TimeoutConfig holds configuration for the timeout stage.
type TimeoutConfig struct {
	Enabled bool
	Duration time.Duration
}

// DefaultTimeoutConfig returns a safe default (disabled, 5 s).
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Enabled:  false,
		Duration: 5 * time.Second,
	}
}

// Validate returns an error when the config is logically inconsistent.
func (c TimeoutConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Duration <= 0 {
		return fmt.Errorf("timeout duration must be positive when enabled, got %v", c.Duration)
	}
	return nil
}

// NewTimeouterFromConfig constructs a Timeouter from the supplied config.
// Returns nil when the config is disabled.
func NewTimeouterFromConfig(cfg TimeoutConfig) (*Timeouter, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return NewTimeouter(cfg.Duration)
}
