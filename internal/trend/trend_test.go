package trend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/trend"
)

func TestNewInvalidWindowReturnsError(t *testing.T) {
	_, err := trend.New(0)
	if err == nil {
		t.Fatal("expected error for zero window, got nil")
	}
	_, err = trend.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window, got nil")
	}
}

func TestNewValidWindowNoError(t *testing.T) {
	tr, err := trend.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestEmptyTrackerIsNeutral(t *testing.T) {
	tr, _ := trend.New(time.Minute)
	if got := tr.Trend(); got != trend.DirectionNeutral {
		t.Errorf("expected neutral, got %s", got)
	}
}

func TestMoreOpenedIsRising(t *testing.T) {
	tr, _ := trend.New(time.Minute)
	tr.Record(5, 1)
	if got := tr.Trend(); got != trend.DirectionRising {
		t.Errorf("expected rising, got %s", got)
	}
}

func TestMoreClosedIsFalling(t *testing.T) {
	tr, _ := trend.New(time.Minute)
	tr.Record(1, 8)
	if got := tr.Trend(); got != trend.DirectionFalling {
		t.Errorf("expected falling, got %s", got)
	}
}

func TestEqualOpenedAndClosedIsNeutral(t *testing.T) {
	tr, _ := trend.New(time.Minute)
	tr.Record(3, 3)
	if got := tr.Trend(); got != trend.DirectionNeutral {
		t.Errorf("expected neutral, got %s", got)
	}
}

func TestAccumulatesAcrossMultipleRecords(t *testing.T) {
	tr, _ := trend.New(time.Minute)
	tr.Record(2, 1)
	tr.Record(1, 3)
	// total opened=3, closed=4 → falling
	if got := tr.Trend(); got != trend.DirectionFalling {
		t.Errorf("expected falling, got %s", got)
	}
}

func TestEvictsExpiredBuckets(t *testing.T) {
	tr, _ := trend.New(50 * time.Millisecond)
	tr.Record(10, 0) // strongly rising
	time.Sleep(80 * time.Millisecond)
	tr.Record(0, 1) // now only falling event in window
	if got := tr.Trend(); got != trend.DirectionFalling {
		t.Errorf("expected falling after eviction, got %s", got)
	}
}

func TestDirectionStringLabels(t *testing.T) {
	cases := []struct {
		d    trend.Direction
		want string
	}{
		{trend.DirectionNeutral, "neutral"},
		{trend.DirectionRising, "rising"},
		{trend.DirectionFalling, "falling"},
	}
	for _, tc := range cases {
		if got := tc.d.String(); got != tc.want {
			t.Errorf("Direction(%d).String() = %q, want %q", tc.d, got, tc.want)
		}
	}
}
