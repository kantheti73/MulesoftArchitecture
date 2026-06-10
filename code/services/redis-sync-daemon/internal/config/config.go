// Package config loads and validates the daemon's YAML configuration.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration parsed from /etc/redis-sync-daemon/config.yaml.
type Config struct {
	Local         LocalRedis      `yaml:"local"`
	Peers         []PeerEndpoint  `yaml:"peers"`
	Sync          SyncConfig      `yaml:"sync"`
	Server        ServerConfig    `yaml:"server"`
	Observability ObservabilityConfig `yaml:"observability"`
}

// LocalRedis describes the local Redis Sentinel cluster the daemon attaches to.
type LocalRedis struct {
	RedisEndpoints     []RedisEndpoint `yaml:"redis_endpoints"`
	SentinelMaster     string          `yaml:"sentinel_master"`
	AuthPasswordSecret string          `yaml:"auth_password_secret"` // vault:// URI
	TLS                TLSConfig       `yaml:"tls"`
}

type RedisEndpoint struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TLSConfig struct {
	Enabled        bool   `yaml:"enabled"`
	CACertPath     string `yaml:"ca_cert_path"`
	ClientCertPath string `yaml:"client_cert_path"`
	ClientKeyPath  string `yaml:"client_key_path"`
}

// PeerEndpoint describes one peer DC's sync daemon endpoint.
type PeerEndpoint struct {
	Name      string     `yaml:"name"`
	Endpoint  string     `yaml:"endpoint"`
	MTLS      TLSConfig  `yaml:"mtls"`
	TimeoutMs int        `yaml:"timeout_ms"`
}

// SyncConfig defines which key patterns to replicate and how.
type SyncConfig struct {
	Patterns        []PatternConfig `yaml:"patterns"`
	KeyspaceChannel string          `yaml:"keyspace_channel"`
}

type PatternConfig struct {
	Prefix                       string `yaml:"prefix"`
	Semantics                    string `yaml:"semantics"` // "setnx" | "counter-delta"
	TTLPreserve                  bool   `yaml:"ttl_preserve"`
	AggregationIntervalSeconds   int    `yaml:"aggregation_interval_seconds"`
}

type ServerConfig struct {
	ListenAddr  string `yaml:"listen_addr"`
	MetricsAddr string `yaml:"metrics_addr"`
	HealthAddr  string `yaml:"health_addr"`
}

type ObservabilityConfig struct {
	PrometheusMetrics bool   `yaml:"prometheus_metrics"`
	LogLevel          string `yaml:"log_level"`
	AuditLogPath      string `yaml:"audit_log_path"`
}

// Load reads and validates the config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	c.applyDefaults()
	return &c, nil
}

func (c *Config) validate() error {
	if len(c.Local.RedisEndpoints) == 0 {
		return fmt.Errorf("local.redis_endpoints must have at least 1 entry")
	}
	if c.Local.SentinelMaster == "" {
		return fmt.Errorf("local.sentinel_master required")
	}
	if len(c.Peers) == 0 {
		return fmt.Errorf("at least one peer required")
	}
	if len(c.Sync.Patterns) == 0 {
		return fmt.Errorf("sync.patterns must have at least 1 entry")
	}
	for _, p := range c.Sync.Patterns {
		if p.Prefix == "" {
			return fmt.Errorf("sync.patterns: prefix required")
		}
		switch p.Semantics {
		case "setnx", "counter-delta":
			// ok
		default:
			return fmt.Errorf("sync.patterns: unknown semantics %q (allowed: setnx, counter-delta)", p.Semantics)
		}
	}
	return nil
}

func (c *Config) applyDefaults() {
	for i := range c.Peers {
		if c.Peers[i].TimeoutMs == 0 {
			c.Peers[i].TimeoutMs = 2000
		}
	}
	for i := range c.Sync.Patterns {
		if c.Sync.Patterns[i].Semantics == "counter-delta" && c.Sync.Patterns[i].AggregationIntervalSeconds == 0 {
			c.Sync.Patterns[i].AggregationIntervalSeconds = 10
		}
	}
	if c.Server.ListenAddr == "" {
		c.Server.ListenAddr = "0.0.0.0:8443"
	}
	if c.Server.MetricsAddr == "" {
		c.Server.MetricsAddr = "0.0.0.0:9090"
	}
	if c.Server.HealthAddr == "" {
		c.Server.HealthAddr = "0.0.0.0:8080"
	}
	if c.Sync.KeyspaceChannel == "" {
		c.Sync.KeyspaceChannel = "__keyspace@0__:*"
	}
}

// PeerTimeout returns the timeout for a given peer as a Duration.
func (p PeerEndpoint) Timeout() time.Duration {
	return time.Duration(p.TimeoutMs) * time.Millisecond
}
