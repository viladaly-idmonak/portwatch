package health_test

import (
	"testing"

	"github.com/user/portwatch/internal/health"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := health.DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected enabled=true by default")
	}
	if cfg.Addr != ":9090" {
		t.Errorf("unexpected default addr: %s", cfg.Addr)
	}
}

func TestNewFromConfigDisabled(t *testing.T) {
	cfg := health.Config{Enabled: false}
	s := health.NewFromConfig(cfg)
	if s != nil {
		t.Error("expected nil server when disabled")
	}
}

func TestNewFromConfigEnabled(t *testing.T) {
	cfg := health.Config{Enabled: true, Addr: "127.0.0.1:19998"}
	s := health.NewFromConfig(cfg)
	if s == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestNewFromConfigEmptyAddrUsesDefault(t *testing.T) {
	cfg := health.Config{Enabled: true, Addr: ""}
	s := health.NewFromConfig(cfg)
	if s == nil {
		t.Fatal("expected non-nil server with empty addr")
	}
}
