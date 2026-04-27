package retryqueue

import (
	"errors"
	"time"
)

// Config holds tuning parameters for the retry queue.
type Config struct {
	// MaxSize is the maximum number of diffs held in the queue at once.
	MaxSize int
	// MaxAttempts is the total number of delivery attempts before a diff is dropped.
	MaxAttempts int
	// BaseDelay is the initial back-off delay after the first failure.
	BaseDelay time.Duration
	// MaxDelay caps the exponential back-off growth.
	MaxDelay time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSize:     256,
		MaxAttempts: 5,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    30 * time.Second,
	}
}

// Validate returns an error if any field holds an invalid value.
func (c Config) Validate() error {
	if c.MaxSize <= 0 {
		return errors.New("retryqueue: MaxSize must be > 0")
	}
	if c.MaxAttempts <= 0 {
		return errors.New("retryqueue: MaxAttempts must be > 0")
	}
	if c.BaseDelay <= 0 {
		return errors.New("retryqueue: BaseDelay must be > 0")
	}
	if c.MaxDelay < c.BaseDelay {
		return errors.New("retryqueue: MaxDelay must be >= BaseDelay")
	}
	return nil
}

// NewFromConfig validates cfg and returns it unchanged, or an error.
// It exists so callers can follow the same New* pattern used across the project.
func NewFromConfig(cfg Config) (Config, error) {
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
