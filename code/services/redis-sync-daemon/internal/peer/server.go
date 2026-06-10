// Package peer implements the HTTPS server (receives sync messages from peer DCs)
// and HTTPS client (sends sync messages to peer DCs).
package peer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
)

// SyncMessage is the wire format exchanged between peer daemons.
type SyncMessage struct {
	Kind     string `json:"kind"`     // "set" | "counter-delta"
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`     // for "set"
	Delta    int64  `json:"delta,omitempty"`     // for "counter-delta"
	TTLSec   int64  `json:"ttl_sec,omitempty"`
	OriginDC string `json:"origin_dc"`
	SentAt   int64  `json:"sent_at_unix_ms"`
}

// Server receives SyncMessages from peers and applies them to local Redis.
type Server struct {
	cfg   config.ServerConfig
	local *redisclient.Client
	met   *metrics.Metrics
	srv   *http.Server
}

func NewServer(cfg config.ServerConfig, local *redisclient.Client, met *metrics.Metrics) *Server {
	return &Server{cfg: cfg, local: local, met: met}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/sync", s.handleSync)

	// mTLS — require client certs signed by the same internal CA.
	// Cert paths are read from env vars for the scaffold; production
	// reads from config + Vault.
	caCert, err := os.ReadFile(os.Getenv("PEER_SERVER_CA_CERT_PATH"))
	if err != nil {
		return fmt.Errorf("read CA cert: %w", err)
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	srvCert, err := tls.LoadX509KeyPair(
		os.Getenv("PEER_SERVER_CERT_PATH"),
		os.Getenv("PEER_SERVER_KEY_PATH"),
	)
	if err != nil {
		return fmt.Errorf("load server cert/key: %w", err)
	}

	s.srv = &http.Server{
		Addr:    s.cfg.ListenAddr,
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{srvCert},
			ClientCAs:    caPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS12,
		},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("peer server: listening on %s (mTLS)", s.cfg.ListenAddr)
	err = s.srv.ListenAndServeTLS("", "")
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

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg SyncMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		s.met.SyncMessagesReceived.WithLabelValues(msg.Kind, "bad-request").Inc()
		return
	}

	ctx := r.Context()
	var err error
	switch msg.Kind {
	case "set":
		_, err = s.local.ApplyReplicatedSet(ctx, msg.Key, msg.Value, msg.TTLSec)
	case "counter-delta":
		_, err = s.local.ApplyCounterDelta(ctx, msg.Key, msg.Delta, msg.TTLSec)
	default:
		http.Error(w, "unknown kind", http.StatusBadRequest)
		s.met.SyncMessagesReceived.WithLabelValues(msg.Kind, "unknown-kind").Inc()
		return
	}

	if err != nil {
		log.Printf("peer server: apply %s for %s failed: %v", msg.Kind, msg.Key, err)
		s.met.SyncMessagesReceived.WithLabelValues(msg.Kind, "apply-error").Inc()
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Record observed lag (sender clock vs ours)
	if msg.SentAt > 0 {
		lag := float64(time.Now().UnixMilli()-msg.SentAt) / 1000.0
		s.met.ObservedLagSeconds.WithLabelValues(msg.OriginDC, msg.Kind).Observe(lag)
	}

	s.met.SyncMessagesReceived.WithLabelValues(msg.Kind, "ok").Inc()
	w.WriteHeader(http.StatusNoContent)
}
