package sampling

import "fmt"

// Config holds configuration for the Sampler.
type Config struct {
	// Enabled controls whether sampling is active. When false, all events pass through.
	Enabled bool `yaml:"enabled"`
	// Rate is the fraction of events to forward (0.0, 1.0].
	Rate float64 `yaml:"rate"`
}

// DefaultConfig returns a Config with sampling disabled.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Rate:    1.0,
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Rate <= 0 || c.Rate > 1 {
		return fmt.Errorf("sampling: rate must be in (0, 1], got %v", c.Rate)
	}
	return nil
}

// NewFromConfig constructs a Sampler from Config.
// Returns nil, nil when sampling is disabled.
func NewFromConfig(c Config) (*Sampler, error) {
	if !c.Enabled {
		return nil, nil
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Rate)
}

// String returns a human-readable description of the Config.
func (c Config) String() string {
	if !c.Enabled {
		return "sampling disabled"
	}
	return fmt.Sprintf("sampling enabled at %.2f%% rate", c.Rate*100)
}
