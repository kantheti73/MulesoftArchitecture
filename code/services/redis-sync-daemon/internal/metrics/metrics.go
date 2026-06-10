// Package metrics defines Prometheus metrics + the /metrics HTTP server.
package metrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds the Prometheus collectors for the daemon.
type Metrics struct {
	KeysReplicated           *prometheus.CounterVec
	ReplicationErrors        *prometheus.CounterVec
	ReplicationLatency       *prometheus.HistogramVec
	AggregationCycleDuration *prometheus.HistogramVec
	PeerCallLatency          *prometheus.HistogramVec
	PeerCallErrors           *prometheus.CounterVec
	SyncMessagesReceived     *prometheus.CounterVec
	ObservedLagSeconds       *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		KeysReplicated: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "redis_sync_keys_replicated_total",
			Help: "Total keys replicated out, by prefix and peer DC",
		}, []string{"prefix", "peer"}),

		ReplicationErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "redis_sync_replication_errors_total",
			Help: "Replication errors, by prefix and stage",
		}, []string{"prefix", "stage"}),

		ReplicationLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "redis_sync_replication_latency_seconds",
			Help:    "Time from local SET observation to all peers acknowledged",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		}, []string{"prefix"}),

		AggregationCycleDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "redis_sync_aggregation_cycle_seconds",
			Help:    "Duration of one counter-aggregation cycle",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		}, []string{"prefix"}),

		PeerCallLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "redis_sync_peer_call_latency_seconds",
			Help:    "HTTPS round-trip to peer DC",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		}, []string{"peer"}),

		PeerCallErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "redis_sync_peer_call_errors_total",
			Help: "HTTPS errors calling peer DC, by peer and status",
		}, []string{"peer", "status"}),

		SyncMessagesReceived: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "redis_sync_messages_received_total",
			Help: "Sync messages received from peer DCs",
		}, []string{"kind", "outcome"}),

		ObservedLagSeconds: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "redis_sync_observed_lag_seconds",
			Help:    "Wall-clock lag between peer SentAt and our receipt",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30},
		}, []string{"origin_dc", "kind"}),
	}
}

// Server is the /metrics HTTP server.
type Server struct {
	addr string
	srv  *http.Server
}

func NewServer(addr string) *Server { return &Server{addr: addr} }

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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
