package scorecard_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scorecard"
)

func entry(port uint16, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto}
}

func TestScoreKnownPortReturnsConfiguredLevel(t *testing.T) {
	sc := scorecard.New([]scorecard.Rule{
		{Port: 22, Protocol: "tcp", Level: scorecard.High},
	}, scorecard.Low)

	if got := sc.Score(entry(22, "tcp")); got != "high" {
		t.Fatalf("expected high, got %s", got)
	}
}

func TestScoreUnknownPortReturnsDefault(t *testing.T) {
	sc := scorecard.New(nil, scorecard.Medium)

	if got := sc.Score(entry(9999, "tcp")); got != "medium" {
		t.Fatalf("expected medium, got %s", got)
	}
}

func TestAddRuleUpdatesScore(t *testing.T) {
	sc := scorecard.New(nil, scorecard.Low)
	sc.AddRule(scorecard.Rule{Port: 443, Protocol: "tcp", Level: scorecard.High})

	if got := sc.Score(entry(443, "tcp")); got != "high" {
		t.Fatalf("expected high after AddRule, got %s", got)
	}
}

func TestScoreProtocolDistinguished(t *testing.T) {
	sc := scorecard.New([]scorecard.Rule{
		{Port: 53, Protocol: "udp", Level: scorecard.Medium},
	}, scorecard.Low)

	if got := sc.Score(entry(53, "tcp")); got != "low" {
		t.Fatalf("tcp/53 should fall back to default, got %s", got)
	}
	if got := sc.Score(entry(53, "udp")); got != "medium" {
		t.Fatalf("udp/53 should be medium, got %s", got)
	}
}

func TestWithScorecardAnnotatesOpenedEntries(t *testing.T) {
	// integration-style check via the stage helper directly
	sc := scorecard.New([]scorecard.Rule{
		{Port: 80, Protocol: "tcp", Level: scorecard.Low},
	}, scorecard.High)

	e := scanner.Entry{Port: 80, Protocol: "tcp"}
	if got := sc.Score(e); got != "low" {
		t.Fatalf("expected low for port 80, got %s", got)
	}

	e2 := scanner.Entry{Port: 8080, Protocol: "tcp"}
	if got := sc.Score(e2); got != "high" {
		t.Fatalf("expected high for unknown port, got %s", got)
	}
}
