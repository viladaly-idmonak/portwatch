// Package scorecard assigns a simple textual risk score to port entries based
// on configurable rules. Scores are one of "low", "medium", or "high".
package scorecard

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents a risk score level.
type Level string

const (
	Low    Level = "low"
	Medium Level = "medium"
	High   Level = "high"
)

// Rule maps a port/protocol pair to a fixed risk level.
type Rule struct {
	Port     uint16
	Protocol string
	Level    Level
}

// Scorecard holds a set of rules and returns a risk score for any entry.
type Scorecard struct {
	mu      sync.RWMutex
	rules   map[string]Level
	default_ Level
}

// New creates a Scorecard with the given rules and a fallback default level.
func New(rules []Rule, defaultLevel Level) *Scorecard {
	sc := &Scorecard{
		rules:   make(map[string]Level, len(rules)),
		default_: defaultLevel,
	}
	for _, r := range rules {
		sc.rules[key(r.Port, r.Protocol)] = r.Level
	}
	return sc
}

// Score returns the risk level string for the given entry.
func (sc *Scorecard) Score(e scanner.Entry) string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if lvl, ok := sc.rules[key(e.Port, e.Protocol)]; ok {
		return string(lvl)
	}
	return string(sc.default_)
}

// AddRule inserts or replaces a rule at runtime.
func (sc *Scorecard) AddRule(r Rule) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.rules[key(r.Port, r.Protocol)] = r.Level
}

func key(port uint16, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
