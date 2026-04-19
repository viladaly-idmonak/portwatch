package geoip_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/geoip"
)

func startGeoServer(t *testing.T, info geoip.Info) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}))
}

func TestLookupPublicIP(t *testing.T) {
	want := geoip.Info{IP: "1.2.3.4", Country: "US", City: "New York", Org: "ExampleISP"}
	srv := startGeoServer(t, want)
	defer srv.Close()

	c := geoip.New(srv.URL)
	got, err := c.Lookup("1.2.3.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Country != want.Country || got.City != want.City {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLookupPrivateIPSkipsHTTP(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	c := geoip.New(srv.URL)
	info, err := c.Lookup("192.168.1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for private IP")
	}
	if info.Country != "private" {
		t.Errorf("expected country=private, got %q", info.Country)
	}
}

func TestLookupCachesResult(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		json.NewEncoder(w).Encode(geoip.Info{IP: "8.8.8.8", Country: "US"})
	}))
	defer srv.Close()

	c := geoip.New(srv.URL)
	for i := 0; i < 3; i++ {
		_, err := c.Lookup("8.8.8.8")
		if err != nil {
			t.Fatalf("lookup %d failed: %v", i, err)
		}
	}
	if hits != 1 {
		t.Errorf("expected 1 HTTP call, got %d", hits)
	}
}
