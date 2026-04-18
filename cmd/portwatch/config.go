package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval    time.Duration `yaml:"interval"`
	Protocols   []string      `yaml:"protocols"`
	Include     []int         `yaml:"include"`
	Exclude     []int         `yaml:"exclude"`
	LogFile     string        `yaml:"log_file"`
	Format      string        `yaml:"format"`
	SnapshotPath string       `yaml:"snapshot_path"`
	RateWindow  time.Duration `yaml:"rate_window"`
}

func loadConfig(path string) (*Config, error) {
	cfg := &Config{
		Interval:     2 * time.Second,
		Protocols:    []string{"tcp", "udp"},
		Format:       "text",
		SnapshotPath: "",
		RateWindow:   5 * time.Second,
	}
	if path == "" {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	for _, p := range cfg.Protocols {
		if p != "tcp" && p != "udp" {
			return nil, fmt.Errorf("invalid protocol %q: must be tcp or udp", p)
		}
	}
	return cfg, nil
}
