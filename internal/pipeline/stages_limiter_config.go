package pipeline

import "fmt"

// LimiterConfig holds configuration for the Limiter stage.
type LimiterConfig struct {
	// Enabled controls whether the limiter stage is active.
	Enabled bool `yaml:"enabled" json:"enabled"`
	// MaxEntries is the maximum number of diff entries allowed per scan cycle.
	MaxEntries int `yaml:"max_entries" json:"max_entries"`
}

// DefaultLimiterConfig returns a safe default configuration (disabled).
func DefaultLimiterConfig() LimiterConfig {
	return LimiterConfig{
		Enabled:    false,
		MaxEntries: 100,
	}
}

// Validate returns an error if the configuration is invalid.
func (c LimiterConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxEntries <= 0 {
		return fmt.Errorf("pipeline/limiter: max_entries must be > 0, got %d", c.MaxEntries)
	}
	return nil
}

// NewLimiterFromConfig constructs a Limiter from config, or returns nil when
// the stage is disabled.
func NewLimiterFromConfig(cfg LimiterConfig) (*Limiter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	return NewLimiter(cfg.MaxEntries)
}
