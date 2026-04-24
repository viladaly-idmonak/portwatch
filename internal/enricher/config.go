package enricher

import "fmt"

// Config holds enricher configuration.
type Config struct {
	// Enabled controls whether enrichment is applied.
	Enabled bool
	// Fields is the static set of key-value pairs to attach.
	Fields map[string]string
}

// DefaultConfig returns a disabled enricher config with no fields.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Fields:  map[string]string{},
	}
}

// Validate returns an error if the config is invalid.
func (c Config) Validate() error {
	for k := range c.Fields {
		if k == "" {
			return fmt.Errorf("enricher: field key must not be empty")
		}
	}
	return nil
}

// NewFromConfig returns a new Enricher when enabled, or nil when disabled.
// Returns an error if validation fails.
func NewFromConfig(c Config) (*Enricher, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if !c.Enabled {
		return nil, nil
	}
	return New(c.Fields), nil
}
