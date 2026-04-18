package health

// Config holds configuration for the health server.
type Config struct {
	// Enabled controls whether the health HTTP server starts.
	Enabled bool `yaml:"enabled" json:"enabled"`
	// Addr is the TCP address to listen on, e.g. ":9090".
	Addr string `yaml:"addr" json:"addr"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Addr:    ":9090",
	}
}

// Validate checks that the Config fields are valid.
// It returns an error if Enabled is true but Addr is empty after applying defaults.
func (c *Config) Validate() error {
	if c.Enabled && c.Addr == "" {
		return fmt.Errorf("health: addr must not be empty when enabled")
	}
	return nil
}

// NewFromConfig constructs a Server from a Config.
// Returns nil if the server is disabled.
func NewFromConfig(cfg Config) *Server {
	if !cfg.Enabled {
		return nil
	}
	addr := cfg.Addr
	if addr == "" {
		addr = DefaultConfig().Addr
	}
	return New(addr)
}
