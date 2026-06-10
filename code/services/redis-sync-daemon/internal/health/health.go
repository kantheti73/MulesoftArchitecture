// Package health exposes /healthz (liveness) and /readyz (readiness).
package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
)

// EngineReadinessProvider is the contract the sync engine satisfies.
// Defined here to avoid a cycle with the sync package.
type EngineReadinessProvider interface {
	Ready() bool
	ObservedLag() float64
}

// Server hosts /healthz and /readyz.
type Server struct {
	addr   string
	local  *redisclient.Client
	engine EngineReadinessProvider
	srv    *http.Server
}

func NewServer(addr string, local *redisclient.Client, engine EngineReadinessProvider) *Server {
	return &Server{addr: addr, local: local, engine: engine}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleLive)
	mux.HandleFunc("/readyz", s.handleReady)

	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleLive(w http.ResponseWriter, r *http.Request) {
	// Liveness: process is up + can respond. Don't check Redis here — that
	// belongs in readiness.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"live"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := s.local.Ping(ctx); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"not-ready","reason":"local-redis-ping","error":%q}`, err.Error()), http.StatusServiceUnavailable)
		return
	}
	if !s.engine.Ready() {
		http.Error(w, `{"status":"not-ready","reason":"engine-not-ready"}`, http.StatusServiceUnavailable)
		return
	}

	// Alarm on excessive lag, but still report ready
	lag := s.engine.ObservedLag()
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"status":"ready","observed_lag_seconds":%.2f}`, lag)
}
