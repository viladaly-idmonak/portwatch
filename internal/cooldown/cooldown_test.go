package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func TestNewInvalidDurationReturnsError(t *testing.T) {
	_, err := cooldown.New(0)
	if err == nil {
		t.Fatal("expected error for zero duration, got nil")
	}
	_, err = cooldown.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration, got nil")
	}
}

func TestFirstCallAlwaysAllowed(t *testing.T) {
	c, _ := cooldown.New(time.Second)
	if !c.Allow("key1") {
		t.Fatal("first call should always be allowed")
	}
}

func TestSecondCallWithinCooldownBlocked(t *testing.T) {
	c, _ := cooldown.New(time.Hour)
	c.Allow("key1")
	if c.Allow("key1") {
		t.Fatal("second call within cooldown window should be blocked")
	}
}

func TestCallAfterCooldownAllowed(t *testing.T) {
	c, _ := cooldown.New(50 * time.Millisecond)
	c.Allow("key1")
	time.Sleep(60 * time.Millisecond)
	if !c.Allow("key1") {
		t.Fatal("call after cooldown period should be allowed")
	}
}

func TestDistinctKeysAreIndependent(t *testing.T) {
	c, _ := cooldown.New(time.Hour)
	if !c.Allow("a") {
		t.Fatal("first call for key 'a' should be allowed")
	}
	if !c.Allow("b") {
		t.Fatal("first call for key 'b' should be allowed")
	}
	if c.Allow("a") {
		t.Fatal("second call for key 'a' should be blocked")
	}
}

func TestResetAllowsImmediateCall(t *testing.T) {
	c, _ := cooldown.New(time.Hour)
	c.Allow("key1")
	c.Reset("key1")
	if !c.Allow("key1") {
		t.Fatal("call after Reset should be allowed immediately")
	}
}

func TestPurgeRemovesExpiredEntries(t *testing.T) {
	c, _ := cooldown.New(50 * time.Millisecond)
	c.Allow("key1")
	c.Allow("key2")
	time.Sleep(60 * time.Millisecond)
	c.Purge()
	// After purge both keys should be gone, so next Allow should succeed.
	if !c.Allow("key1") {
		t.Fatal("key1 should be allowed after purge")
	}
	if !c.Allow("key2") {
		t.Fatal("key2 should be allowed after purge")
	}
}
