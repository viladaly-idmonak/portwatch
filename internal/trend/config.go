package trend

import "fmt"

// Config holds configuration for the Tracker.
type Config struct {
	// Enabled controls whether trend tracking is active.
	Enabled bool

	// Window is the number of recent scans to consider when computing trend.
	Window int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Window:  10,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return fmt.Errorf("trend: window must be positive, got %d", c.Window)
	}
	return nil
}

// NewFromConfig returns a new Tracker built from cfg, or nil when disabled.
func NewFromConfig(cfg Config) (*Tracker, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return New(cfg.Window)
}
