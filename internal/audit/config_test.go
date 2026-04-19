package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestDefaultConfigValues(t *testing.T) {
	c := audit.DefaultConfig()
	if c.Enabled {
		t.Error("expected disabled by default")
	}
	if c.FilePath == "" {
		t.Error("expected non-empty default file path")
	}
}

func TestValidateRejectsEnabledWithEmptyPath(t *testing.T) {
	c := audit.Config{Enabled: true, FilePath: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected validation error")
	}
}

func TestValidateAcceptsDisabledWithEmptyPath(t *testing.T) {
	c := audit.Config{Enabled: false, FilePath: ""}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	a, err := audit.NewFromConfig(audit.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	if a != nil {
		t.Error("expected nil auditor when disabled")
	}
}

func TestNewFromConfigEnabledCreatesFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "audit.jsonl")
	c := audit.Config{Enabled: true, FilePath: tmp}
	a, err := audit.NewFromConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	if a == nil {
		t.Fatal("expected non-nil auditor")
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
