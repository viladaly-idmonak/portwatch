package enricher_test

import (
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/scanner"
)

func entry(port uint16) scanner.Entry {
	return scanner.Entry{Port: port, Proto: "tcp"}
}

func diff(opened, closed []scanner.Entry) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func TestApplyEmptyDiffPassesThrough(t *testing.T) {
	e := enricher.New(map[string]string{"host": "box1"})
	out := e.Apply(scanner.Diff{})
	if len(out.Opened)+len(out.Closed) != 0 {
		t.Fatal("expected empty diff")
	}
}

func TestApplyAttachesFields(t *testing.T) {
	e := enricher.New(map[string]string{"env": "prod", "host": "box1"})
	out := e.Apply(diff([]scanner.Entry{entry(80)}, nil))
	if out.Opened[0].Meta["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", out.Opened[0].Meta["env"])
	}
	if out.Opened[0].Meta["host"] != "box1" {
		t.Errorf("expected host=box1, got %q", out.Opened[0].Meta["host"])
	}
}

func TestApplyExistingMetaKeyNotOverwritten(t *testing.T) {
	e := enricher.New(map[string]string{"env": "prod"})
	en := entry(443)
	en.Meta = map[string]string{"env": "staging"}
	out := e.Apply(diff([]scanner.Entry{en}, nil))
	if out.Opened[0].Meta["env"] != "staging" {
		t.Errorf("existing key should not be overwritten, got %q", out.Opened[0].Meta["env"])
	}
}

func TestApplyClosedEntries(t *testing.T) {
	e := enricher.New(map[string]string{"region": "us-east"})
	out := e.Apply(diff(nil, []scanner.Entry{entry(22)}))
	if out.Closed[0].Meta["region"] != "us-east" {
		t.Errorf("expected region=us-east on closed entry")
	}
}

func TestSetAndDelete(t *testing.T) {
	e := enricher.New(nil)
	e.Set("k", "v")
	if e.Fields()["k"] != "v" {
		t.Fatal("Set did not store field")
	}
	e.Delete("k")
	if _, ok := e.Fields()["k"]; ok {
		t.Fatal("Delete did not remove field")
	}
}
