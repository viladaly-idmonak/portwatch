package ratelimit

import (
	"testing"
	"time"
)

func TestAllowFirstEventPasses(t *testing.T) {
	l := New(100 * time.Millisecond)
	if !l.Allow("tcp:8080:opened") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllowDuplicateWithinWindowBlocked(t *testing.T) {
	l := New(200 * time.Millisecond)
	l.Allow("tcp:8080:opened")
	if l.Allow("tcp:8080:opened") {
		t.Fatal("expected duplicate within cooldown to be blocked")
	}
}

func TestAllowDuplicateAfterWindowPasses(t *testing.T) {
	l := New(50 * time.Millisecond)
	l.Allow("tcp:9090:closed")
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("tcp:9090:closed") {
		t.Fatal("expected event after cooldown to be allowed")
	}
}

func TestAllowDistinctKeysBothPass(t *testing.T) {
	l := New(500 * time.Millisecond)
	if !l.Allow("tcp:80:opened") {
		t.Fatal("expected first key to pass")
	}
	if !l.Allow("tcp:443:opened") {
		t.Fatal("expected distinct key to pass")
	}
}

func TestPurgeRemovesExpiredEntries(t *testing.T) {
	l := New(30 * time.Millisecond)
	l.Allow("tcp:22:opened")
	time.Sleep(40 * time.Millisecond)
	l.Purge()

	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.seen["tcp:22:opened"]; ok {
		t.Fatal("expected expired entry to be purged")
	}
}
