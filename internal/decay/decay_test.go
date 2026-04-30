package decay_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/decay"
)

func TestNewInvalidHalfLifeReturnsError(t *testing.T) {
	_, err := decay.New(0)
	if err == nil {
		t.Fatal("expected error for zero half-life")
	}
	_, err = decay.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative half-life")
	}
}

func TestNewValidHalfLifeNoError(t *testing.T) {
	d, err := decay.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Decayer")
	}
}

func TestAddAccumulatesScore(t *testing.T) {
	d, _ := decay.New(time.Hour) // very long half-life: negligible decay
	score := d.Add("tcp:80", 1.0)
	if score < 0.99 {
		t.Fatalf("expected score near 1.0, got %f", score)
	}
	score = d.Add("tcp:80", 1.0)
	if score < 1.9 {
		t.Fatalf("expected score near 2.0, got %f", score)
	}
}

func TestScoreDecaysOverTime(t *testing.T) {
	halfLife := 50 * time.Millisecond
	d, _ := decay.New(halfLife)
	d.Add("tcp:443", 8.0)

	time.Sleep(halfLife)

	s := d.Score("tcp:443")
	// After one half-life the score should be ~4.0; allow ±1.5 for timing jitter.
	if s < 2.5 || s > 5.5 {
		t.Fatalf("expected score ~4.0 after one half-life, got %f", s)
	}
}

func TestScoreUnknownKeyReturnsZero(t *testing.T) {
	d, _ := decay.New(time.Minute)
	if s := d.Score("udp:53"); s != 0 {
		t.Fatalf("expected 0 for unknown key, got %f", s)
	}
}

func TestDeleteRemovesEntry(t *testing.T) {
	d, _ := decay.New(time.Minute)
	d.Add("tcp:22", 5.0)
	d.Delete("tcp:22")
	if s := d.Score("tcp:22"); s != 0 {
		t.Fatalf("expected 0 after delete, got %f", s)
	}
}

func TestDistinctKeysAreIndependent(t *testing.T) {
	d, _ := decay.New(time.Hour)
	d.Add("tcp:80", 3.0)
	d.Add("udp:53", 7.0)

	if s := d.Score("tcp:80"); s < 2.9 {
		t.Fatalf("tcp:80 score unexpected: %f", s)
	}
	if s := d.Score("udp:53"); s < 6.9 {
		t.Fatalf("udp:53 score unexpected: %f", s)
	}
}
