package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/throttle"
)

func openDiff() scanner.Diff {
	return scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
}

func TestWithDebouncePassesDiff(t *testing.T) {
	p := pipeline.New(pipeline.WithDebounce(10 * time.Millisecond))
	out, err := p.Run(context.Background(), openDiff())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}

func TestWithThrottleAllowsFirst(t *testing.T) {
	cfg := throttle.DefaultConfig()
	tr, _ := throttle.NewFromConfig(cfg)
	p := pipeline.New(pipeline.WithThrottle(tr))
	out, err := p.Run(context.Background(), openDiff())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}

func TestWithThrottleBlocksSecond(t *testing.T) {
	cfg := throttle.DefaultConfig()
	tr, _ := throttle.NewFromConfig(cfg)
	p := pipeline.New(pipeline.WithThrottle(tr))
	p.Run(context.Background(), openDiff()) //nolint
	_, err := p.Run(context.Background(), openDiff())
	if err == nil {
		t.Fatal("expected throttle error on second call")
	}
}

func TestCombinedDebounceAndFilter(t *testing.T) {
	p := pipeline.New(
		pipeline.WithDebounce(5*time.Millisecond),
		pipeline.WithFilter(nil),
	)
	out, err := p.Run(context.Background(), openDiff())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 opened entry, got %d", len(out.Opened))
	}
}
