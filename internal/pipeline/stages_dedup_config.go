package pipeline

import "errors"

// DedupConfig controls deduplication stage behaviour.
type DedupConfig struct {
	// Enabled toggles the dedup stage. When false, WithDedup is a no-op pass-through.
	Enabled bool

	// MaxEntries caps the number of port states held in memory.
	// Zero means unlimited.
	MaxEntries int
}

// DefaultDedupConfig returns a DedupConfig with sensible defaults.
func DefaultDedupConfig() DedupConfig {
	return DedupConfig{
		Enabled:    true,
		MaxEntries: 0,
	}
}

// Validate returns an error if the configuration is invalid.
func (c DedupConfig) Validate() error {
	if c.MaxEntries < 0 {
		return errors.New("dedup: MaxEntries must be >= 0")
	}
	return nil
}

// NewDedupFromConfig returns a DedupConfig derived from the provided values,
// falling back to defaults where the zero value is supplied.
func NewDedupFromConfig(enabled bool, maxEntries int) (DedupConfig, error) {
	cfg := DedupConfig{
		Enabled:    enabled,
		MaxEntries: maxEntries,
	}
	if err := cfg.Validate(); err != nil {
		return DedupConfig{}, err
	}
	return cfg, nil
}
