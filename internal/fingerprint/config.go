package fingerprint

import "fmt"

// Config holds configuration for the Hasher.
type Config struct {
	// Enabled controls whether fingerprinting is active.
	Enabled bool

	// Salt is mixed into every hash to namespace results per deployment.
	Salt string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Salt:    "",
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if c.Enabled && len(c.Salt) > 64 {
		return fmt.Errorf("fingerprint: salt must be 64 characters or fewer, got %d", len(c.Salt))
	}
	return nil
}

// NewFromConfig returns a Hasher when enabled, or nil when disabled.
func NewFromConfig(c Config) (*Hasher, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if !c.Enabled {
		return nil, nil
	}
	return New(c.Salt), nil
}
