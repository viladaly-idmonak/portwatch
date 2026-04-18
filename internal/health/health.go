package health

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	Scans     int64     `json:"scans_total"`
	StartedAt time.Time `json:"started_at"`
}

// Server exposes a lightweight HTTP health endpoint.
type Server struct {
	addr      string
	scans     atomic.Int64
	startedAt time.Time
	server    *http.Server
}

// New creates a new health Server listening on addr (e.g. ":9090").
func New(addr string) *Server {
	s := &Server{addr: addr, startedAt: time.Now()}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// RecordScan increments the scan counter.
func (s *Server) RecordScan() {
	s.scans.Add(1)
}

// Start begins serving in a background goroutine.
func (s *Server) Start() {
	go func() { _ = s.server.ListenAndServe() }()
}

// Stop shuts down the HTTP server.
func (s *Server) Stop() {
	_ = s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := Status{
		OK:        true,
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
		Scans:     s.scans.Load(),
		StartedAt: s.startedAt,
	}
	_ = json.NewEncoder(w).Encode(status)
}
