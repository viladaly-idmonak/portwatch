package classifier

import "fmt"

// RuleConfig is the serialisable form of a single classification rule.
type RuleConfig struct {
	Port     uint16 `yaml:"port"     json:"port"`
	Protocol string `yaml:"protocol" json:"protocol"`
	Class    string `yaml:"class"    json:"class"`
}

// Config holds the full classifier configuration.
type Config struct {
	Enabled      bool         `yaml:"enabled"       json:"enabled"`
	DefaultClass string       `yaml:"default_class" json:"default_class"`
	Rules        []RuleConfig `yaml:"rules"         json:"rules"`
}

// DefaultConfig returns a safe default configuration.
func DefaultConfig() Config {
	return Config{
		Enabled:      false,
		DefaultClass: string(ClassUnknown),
		Rules:        nil,
	}
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DefaultClass == "" {
		return fmt.Errorf("classifier: default_class must not be empty when enabled")
	}
	for i, r := range c.Rules {
		if r.Protocol == "" {
			return fmt.Errorf("classifier: rule[%d] has empty protocol", i)
		}
		if r.Class == "" {
			return fmt.Errorf("classifier: rule[%d] has empty class", i)
		}
	}
	return nil
}

// NewFromConfig constructs a Classifier from Config, or returns nil when disabled.
func NewFromConfig(c Config) (*Classifier, error) {
	if !c.Enabled {
		return nil, nil
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	rules := make([]Rule, len(c.Rules))
	for i, r := range c.Rules {
		rules[i] = Rule{Port: r.Port, Protocol: r.Protocol, Class: Class(r.Class)}
	}
	return New(rules, Class(c.DefaultClass))
}
