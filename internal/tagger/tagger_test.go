package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/tagger"
)

func TestGetUnknownReturnsDefault(t *testing.T) {
	tr := tagger.New(nil)
	label := tr.Label(8080)
	if label != "port/8080" {
		t.Fatalf("expected port/8080, got %s", label)
	}
}

func TestSetAndGet(t *testing.T) {
	tr := tagger.New(nil)
	tr.Set(443, "https")
	v, ok := tr.Get(443)
	if !ok || v != "https" {
		t.Fatalf("expected https, got %s (ok=%v)", v, ok)
	}
}

func TestLabelReturnsTag(t *testing.T) {
	tr := tagger.New(map[uint16]string{22: "ssh"})
	if tr.Label(22) != "ssh" {
		t.Fatal("expected ssh")
	}
}

func TestDeleteRemovesTag(t *testing.T) {
	tr := tagger.New(map[uint16]string{80: "http"})
	tr.Delete(80)
	_, ok := tr.Get(80)
	if ok {
		t.Fatal("expected tag to be removed")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	tr := tagger.New(map[uint16]string{53: "dns", 22: "ssh"})
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(all))
	}
	// mutating the copy should not affect the tagger
	all[9999] = "nope"
	if _, ok := tr.Get(9999); ok {
		t.Fatal("copy mutation affected tagger")
	}
}
