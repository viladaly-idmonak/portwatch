package alerter

import (
	"runtime"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func diffWith(opened, closed []scanner.Entry) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func entry(proto string, port int) scanner.Entry {
	return scanner.Entry{Proto: proto, Port: port}
}

func TestNotifyNoHooksNoErrors(t *testing.T) {
	a := New(Hook{})
	errs := a.Notify(diffWith([]scanner.Entry{entry("tcp", 8080)}, nil))
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifyEmptyDiff(t *testing.T) {
	a := New(Hook{OnOpen: "echo open {port}", OnClose: "echo close {port}"})
	errs := a.Notify(diffWith(nil, nil))
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifyOnOpenHookRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell hook test on windows")
	}
	a := New(Hook{OnOpen: "true"})
	errs := a.Notify(diffWith([]scanner.Entry{entry("tcp", 9000)}, nil))
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifyOnCloseHookRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell hook test on windows")
	}
	a := New(Hook{OnClose: "true"})
	errs := a.Notify(diffWith(nil, []scanner.Entry{entry("udp", 53)}))
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifyBadCommandReturnsError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell hook test on windows")
	}
	a := New(Hook{OnOpen: "false"})
	errs := a.Notify(diffWith([]scanner.Entry{entry("tcp", 443)}, nil))
	if len(errs) == 0 {
		t.Fatal("expected an error from failing hook")
	}
}
