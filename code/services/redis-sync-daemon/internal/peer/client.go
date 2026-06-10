package peer

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
)

// Client sends SyncMessages to one peer DC's sync daemon over mTLS HTTPS.
type Client struct {
	cfg    config.PeerEndpoint
	met    *metrics.Metrics
	httpc  *http.Client
	origin string
}

func NewClient(cfg config.PeerEndpoint, met *metrics.Metrics) (*Client, error) {
	caCert, err := os.ReadFile(cfg.MTLS.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("parse CA cert")
	}
	clientCert, err := tls.LoadX509KeyPair(cfg.MTLS.ClientCertPath, cfg.MTLS.ClientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load client cert/key: %w", err)
	}

	origin := os.Getenv("DAEMON_DC_NAME")
	if origin == "" {
		origin = "unknown"
	}

	return &Client{
		cfg:    cfg,
		met:    met,
		origin: origin,
		httpc: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caPool,
					Certificates: []tls.Certificate{clientCert},
					MinVersion:   tls.VersionTLS12,
				},
				MaxIdleConns:        10,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
			Timeout: cfg.Timeout(),
		},
	}, nil
}

func (c *Client) Name() string { return c.cfg.Name }

func (c *Client) SendSet(ctx context.Context, key, value string, ttlSec int64) error {
	msg := SyncMessage{
		Kind:     "set",
		Key:      key,
		Value:    value,
		TTLSec:   ttlSec,
		OriginDC: c.origin,
		SentAt:   time.Now().UnixMilli(),
	}
	return c.send(ctx, msg)
}

func (c *Client) SendCounterDelta(ctx context.Context, key string, delta, ttlSec int64) error {
	msg := SyncMessage{
		Kind:     "counter-delta",
		Key:      key,
		Delta:    delta,
		TTLSec:   ttlSec,
		OriginDC: c.origin,
		SentAt:   time.Now().UnixMilli(),
	}
	return c.send(ctx, msg)
}

func (c *Client) send(ctx context.Context, msg SyncMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	url := c.cfg.Endpoint + "/v1/sync"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.httpc.Do(req)
	if err != nil {
		c.met.PeerCallErrors.WithLabelValues(c.cfg.Name, "transport").Inc()
		return fmt.Errorf("POST: %w", err)
	}
	defer resp.Body.Close()
	c.met.PeerCallLatency.WithLabelValues(c.cfg.Name).Observe(time.Since(start).Seconds())

	if resp.StatusCode/100 != 2 {
		c.met.PeerCallErrors.WithLabelValues(c.cfg.Name, fmt.Sprintf("%d", resp.StatusCode)).Inc()
		return fmt.Errorf("peer returned status %d", resp.StatusCode)
	}
	return nil
}
