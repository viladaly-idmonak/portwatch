package labeler

import "fmt"

// RuleConfig is the serialisable form of a labeling rule, suitable
// for embedding in a YAML/TOML configuration file.
type RuleConfig struct {
	Port     uint16 `yaml:"port"     toml:"port"`
	Protocol string `yaml:"protocol" toml:"protocol"`
	Label    string `yaml:"label"    toml:"label"`
}

// Config holds all labeler configuration.
type Config struct {
	Enabled bool         `yaml:"enabled" toml:"enabled"`
	Rules   []RuleConfig `yaml:"rules"   toml:"rules"`
}

// DefaultConfig returns a safe zero-value Config.
func DefaultConfig() Config {
	return Config{Enabled: false}
}

// Validate checks that every rule has a non-zero port, a non-empty
// protocol and a non-empty label.
func (c Config) Validate() error {
	for i, r := range c.Rules {
		if r.Port == 0 {
			return fmt.Errorf("labeler: rule[%d]: port must be > 0", i)
		}
		if r.Protocol == "" {
			return fmt.Errorf("labeler: rule[%d]: protocol must not be empty", i)
		}
		if r.Label == "" {
			return fmt.Errorf("labeler: rule[%d]: label must not be empty", i)
		}
	}
	return nil
}

// NewFromConfig validates cfg and, when enabled, returns a configured
// Labeler. Returns nil when disabled.
func NewFromConfig(cfg Config) (*Labeler, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	rules := make([]Rule, len(cfg.Rules))
	for i, r := range cfg.Rules {
		rules[i] = Rule{Port: r.Port, Protocol: r.Protocol, Label: r.Label}
	}
	return New(rules), nil
}
