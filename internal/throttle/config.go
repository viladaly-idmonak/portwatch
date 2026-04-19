package throttle

import (
	"fmt"
	"time"
)

// Config holds throttle configuration.
type Config struct {
	// Interval is the minimum duration between allowed events.
	Interval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 500 * time.Millisecond,
	}
}

// Validate checks that the config values are acceptable.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("throttle: interval must be positive, got %s", c.Interval)
	}
	return nil
}

// NewFromConfig constructs a Throttle from a Config, returning an error if
// the config is invalid.
func NewFromConfig(cfg Config) (*Throttle, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Interval), nil
}

// WithInterval returns a copy of the Config with the given interval set.
func (c Config) WithInterval(d time.Duration) Config {
	c.Interval = d
	return c
}
