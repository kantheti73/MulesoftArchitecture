// Package redis wraps the go-redis Sentinel client with TLS + AUTH from config.
package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	goredis "github.com/go-redis/redis/v9"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
)

// Client wraps a Sentinel-failover-aware Redis client.
type Client struct {
	rdb *goredis.Client
	cfg config.LocalRedis
}

// NewSentinelClient connects to the local Redis Sentinel cluster.
//
// The password is resolved from the vault:// URI in config; this scaffold
// reads it from the env var corresponding to the URI path (you would replace
// this with a real Vault SDK call in production — see vault.go TODO).
func NewSentinelClient(ctx context.Context, cfg config.LocalRedis) (*Client, error) {
	addrs := make([]string, 0, len(cfg.RedisEndpoints))
	for _, e := range cfg.RedisEndpoints {
		addrs = append(addrs, fmt.Sprintf("%s:%d", e.Host, e.Port))
	}

	password, err := resolveSecret(cfg.AuthPasswordSecret)
	if err != nil {
		return nil, fmt.Errorf("resolve auth password: %w", err)
	}

	var tlsCfg *tls.Config
	if cfg.TLS.Enabled {
		tlsCfg, err = buildTLSConfig(cfg.TLS)
		if err != nil {
			return nil, fmt.Errorf("build TLS config: %w", err)
		}
	}

	rdb := goredis.NewFailoverClient(&goredis.FailoverOptions{
		MasterName:    cfg.SentinelMaster,
		SentinelAddrs: addrs,
		Password:      password,
		TLSConfig:     tlsCfg,

		// Sentinels also need auth in production
		SentinelPassword: password,

		// Reasonable defaults for a long-lived daemon
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping local Redis: %w", err)
	}
	return &Client{rdb: rdb, cfg: cfg}, nil
}

// Raw returns the underlying go-redis client.
func (c *Client) Raw() *goredis.Client {
	return c.rdb
}

// Close shuts down the connection pool.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// Ping checks liveness.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// ApplyReplicatedSet performs an idempotency-style SETNX with optional TTL.
// Used by the peer server when receiving a sync message from another DC.
func (c *Client) ApplyReplicatedSet(ctx context.Context, key, value string, ttlSec int64) (bool, error) {
	var ok bool
	var err error
	if ttlSec > 0 {
		ok, err = c.rdb.SetNX(ctx, key, value, durationFromSec(ttlSec)).Result()
	} else {
		ok, err = c.rdb.SetNX(ctx, key, value, 0).Result()
	}
	return ok, err
}

// ApplyCounterDelta does INCRBY for a counter key (idempotent over fresh increments).
// The peer sends deltas, not absolute values, so simultaneous increments in both
// DCs both land additively.
func (c *Client) ApplyCounterDelta(ctx context.Context, key string, delta int64, ttlSec int64) (int64, error) {
	newVal, err := c.rdb.IncrBy(ctx, key, delta).Result()
	if err != nil {
		return 0, err
	}
	if ttlSec > 0 {
		// Only set TTL if not already set (preserve original expiry semantics)
		if existing, _ := c.rdb.TTL(ctx, key).Result(); existing.Seconds() < 0 {
			_ = c.rdb.Expire(ctx, key, durationFromSec(ttlSec)).Err()
		}
	}
	return newVal, nil
}

// resolveSecret reads a secret from the vault:// URI.
//
// SCAFFOLD: this returns from environment for local testing. Production
// implementation should use the HashiCorp Vault SDK with the AppRole or
// Kubernetes auth backend per docs/09-onprem-install.md.
func resolveSecret(uri string) (string, error) {
	if uri == "" {
		return "", nil
	}
	if strings.HasPrefix(uri, "vault://") {
		// TODO: implement real Vault client
		// For local testing, fall through to env var lookup
		envName := strings.ReplaceAll(strings.TrimPrefix(uri, "vault://"), "/", "_")
		envName = strings.ToUpper(envName)
		if v := os.Getenv(envName); v != "" {
			return v, nil
		}
		return "", fmt.Errorf("vault:// secret not implemented in scaffold; set env var %s for local testing", envName)
	}
	// Treat as literal (NOT recommended for production)
	return uri, nil
}

func buildTLSConfig(t config.TLSConfig) (*tls.Config, error) {
	caCert, err := os.ReadFile(t.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("parse CA cert")
	}

	cert, err := tls.LoadX509KeyPair(t.ClientCertPath, t.ClientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load client cert/key: %w", err)
	}

	return &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func durationFromSec(seconds int64) (d goredis.DurationParam) {
	// go-redis takes time.Duration; we keep a thin alias for clarity at call sites
	return goredis.DurationParam(seconds * 1_000_000_000)
}
