package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestFirstCallAlwaysAllowed(t *testing.T) {
	th := throttle.New(100 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSecondCallWithinIntervalBlocked(t *testing.T) {
	th := throttle.New(200 * time.Millisecond)
	th.Allow()
	if th.Allow() {
		t.Fatal("expected second immediate call to be blocked")
	}
}

func TestCallAfterIntervalAllowed(t *testing.T) {
	th := throttle.New(20 * time.Millisecond)
	th.Allow()
	time.Sleep(30 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestResetAllowsImmediateCall(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestSetIntervalTakesEffect(t *testing.T) {
	th := throttle.New(500 * time.Millisecond)
	th.Allow()
	th.SetInterval(10 * time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected call to be allowed after interval reduction")
	}
}

func TestConcurrentAllowIsSafe(t *testing.T) {
	th := throttle.New(1 * time.Millisecond)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			th.Allow()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
}
