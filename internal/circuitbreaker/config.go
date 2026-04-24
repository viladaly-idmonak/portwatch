package circuitbreaker

import (
	"errors"
	"time"
)

// Config holds configuration for a CircuitBreaker.
type Config struct {
	Enabled      bool          `yaml:"enabled"`
	MaxFailures  int           `yaml:"max_failures"`
	ResetTimeout time.Duration `yaml:"reset_timeout"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:      true,
		MaxFailures:  5,
		ResetTimeout: 30 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxFailures <= 0 {
		return errors.New("circuitbreaker: max_failures must be greater than zero")
	}
	if c.ResetTimeout <= 0 {
		return errors.New("circuitbreaker: reset_timeout must be greater than zero")
	}
	return nil
}

// NewFromConfig creates a CircuitBreaker from cfg, or returns nil if disabled.
func NewFromConfig(cfg Config) (*CircuitBreaker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	return New(cfg.MaxFailures, cfg.ResetTimeout), nil
}
