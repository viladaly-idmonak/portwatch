package main

import (
	"encoding/json"
	"errors"
	"os"
)

// Config holds runtime configuration for portwatch.
type Config struct {
	Hosts        []string `json:"hosts"`
	Protocol     string   `json:"protocol"`
	IntervalSecs int      `json:"interval_secs"`
}

const defaultConfigPath = "portwatch.json"

// loadConfig reads config from portwatch.json if present, otherwise returns defaults.
func loadConfig() (*Config, error) {
	cfg := &Config{
		Hosts:        []string{"127.0.0.1"},
		Protocol:     "tcp",
		IntervalSecs: 5,
	}

	f, err := os.Open(defaultConfigPath)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	if len(cfg.Hosts) == 0 {
		return nil, errors.New("config: hosts must not be empty")
	}
	if cfg.IntervalSecs <= 0 {
		return nil, errors.New("config: interval_secs must be positive")
	}
	if cfg.Protocol != "tcp" && cfg.Protocol != "udp" {
		return nil, errors.New("config: protocol must be tcp or udp")
	}
	return cfg, nil
}
