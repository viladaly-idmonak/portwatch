package decay

import "time"

// Config holds configuration for the Decayer.
type Config struct {
	Enabled  bool
	HalfLife time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:  false,
		HalfLife: 5 * time.Minute,
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.HalfLife <= 0 {
		return ErrInvalidHalfLife
	}
	return nil
}

// NewFromConfig creates a Decayer from cfg, returning nil when disabled.
func NewFromConfig(cfg Config) (*Decayer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	return New(cfg.HalfLife)
}
