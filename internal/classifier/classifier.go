package classifier

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Class represents a risk or trust classification for a port.
type Class string

const (
	ClassTrusted  Class = "trusted"
	ClassSuspect  Class = "suspect"
	ClassCritical Class = "critical"
	ClassUnknown  Class = "unknown"
)

// Rule maps a port+protocol pair to a Class.
type Rule struct {
	Port     uint16
	Protocol string
	Class    Class
}

// Classifier assigns a Class to scanner entries based on configured rules.
type Classifier struct {
	mu      sync.RWMutex
	rules   map[string]Class
	default_ Class
}

// New returns a Classifier with the given rules and a fallback default class.
func New(rules []Rule, defaultClass Class) (*Classifier, error) {
	if defaultClass == "" {
		return nil, fmt.Errorf("classifier: default class must not be empty")
	}
	c := &Classifier{
		rules:    make(map[string]Class, len(rules)),
		default_: defaultClass,
	}
	for _, r := range rules {
		if r.Protocol == "" {
			return nil, fmt.Errorf("classifier: rule for port %d has empty protocol", r.Port)
		}
		if r.Class == "" {
			return nil, fmt.Errorf("classifier: rule for port %d has empty class", r.Port)
		}
		c.rules[key(r.Port, r.Protocol)] = r.Class
	}
	return c, nil
}

// Classify returns the Class for the given entry, falling back to the default.
func (c *Classifier) Classify(e scanner.Entry) Class {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cls, ok := c.rules[key(e.Port, e.Protocol)]; ok {
		return cls
	}
	return c.default_
}

// AddRule inserts or replaces a classification rule at runtime.
func (c *Classifier) AddRule(r Rule) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules[key(r.Port, r.Protocol)] = r.Class
}

func key(port uint16, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
