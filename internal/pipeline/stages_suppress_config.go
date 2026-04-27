package pipeline

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

// SuppressConfig holds pipeline-level configuration for the suppress stage.
type SuppressConfig struct {
	// Enabled controls whether the suppress stage is active.
	Enabled bool

	// Window is the duration during which duplicate events are suppressed.
	Window time.Duration
}

// DefaultSuppressConfig returns a SuppressConfig with sensible defaults.
func DefaultSuppressConfig() SuppressConfig {
	return SuppressConfig{
		Enabled: false,
		Window:  30 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c SuppressConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return fmt.Errorf("suppress: window must be positive, got %s", c.Window)
	}
	return nil
}

// NewSuppressFromConfig constructs a suppress.Suppressor from the pipeline
// SuppressConfig. Returns nil when the stage is disabled.
func NewSuppressFromConfig(c SuppressConfig) (*suppress.Suppressor, error) {
	if !c.Enabled {
		return nil, nil
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	sc := suppress.DefaultConfig()
	sc.Window = c.Window
	return suppress.New(sc)
}
