package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
)

func TestDefaultConfigValues(t *testing.T) {
	c := fingerprint.DefaultConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Salt != "" {
		t.Errorf("expected empty Salt, got %q", c.Salt)
	}
}

func TestValidateRejectsSaltTooLong(t *testing.T) {
	c := fingerprint.Config{
		Enabled: true,
		Salt:    strings.Repeat("x", 65),
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for salt > 64 chars")
	}
}

func TestValidateAcceptsDisabledWithLongSalt(t *testing.T) {
	c := fingerprint.Config{
		Enabled: false,
		Salt:    strings.Repeat("x", 65),
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error for disabled config: %v", err)
	}
}

func TestNewFromConfigDisabledReturnsNil(t *testing.T) {
	h, err := fingerprint.NewFromConfig(fingerprint.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h != nil {
		t.Fatal("expected nil Hasher when disabled")
	}
}

func TestNewFromConfigEnabledCreatesHasher(t *testing.T) {
	c := fingerprint.Config{Enabled: true, Salt: "testsalt"}
	h, err := fingerprint.NewFromConfig(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Hasher when enabled")
	}
}

func TestNewFromConfigInvalidReturnsError(t *testing.T) {
	c := fingerprint.Config{Enabled: true, Salt: strings.Repeat("z", 65)}
	_, err := fingerprint.NewFromConfig(c)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
