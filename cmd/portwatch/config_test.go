package main

import (
	"encoding/json"
	"os"
	"testing"
)

func writeConfig(t *testing.T, cfg Config) {
	t.Helper()
	f, err := os.Create(defaultConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(defaultConfigPath) })
}

func TestLoadConfigDefaults(t *testing.T) {
	os.Remove(defaultConfigPath)
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Protocol != "tcp" || cfg.IntervalSecs != 5 {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	writeConfig(t, Config{Hosts: []string{"10.0.0.1"}, Protocol: "tcp", IntervalSecs: 10})
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSecs != 10 || cfg.Hosts[0] != "10.0.0.1" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfigInvalidProtocol(t *testing.T) {
	writeConfig(t, Config{Hosts: []string{"127.0.0.1"}, Protocol: "icmp", IntervalSecs: 5})
	_, err := loadConfig()
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}
