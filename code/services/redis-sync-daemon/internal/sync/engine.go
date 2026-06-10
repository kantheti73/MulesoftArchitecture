// Package sync orchestrates the replicator + aggregator goroutines.
package sync

import (
	"context"
	"log"
	"sync"

	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/peer"
	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
)

// Engine coordinates replicators (per-key fan-out) and aggregators
// (periodic counter convergence) across all configured patterns.
type Engine struct {
	cfg         *config.Config
	local       *redisclient.Client
	peers       []*peer.Client
	met         *metrics.Metrics

	mu          sync.RWMutex
	ready       bool
	lagSeconds  float64 // most-recent sync lag observed; surfaced via /readyz
}

// NewEngine constructs the engine.
func NewEngine(cfg *config.Config, local *redisclient.Client, peers []*peer.Client, met *metrics.Metrics) *Engine {
	return &Engine{cfg: cfg, local: local, peers: peers, met: met}
}

// Run blocks until ctx is cancelled; spawns one goroutine per (pattern × semantics).
func (e *Engine) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, pattern := range e.cfg.Sync.Patterns {
		switch pattern.Semantics {
		case "setnx":
			wg.Add(1)
			go func(p config.PatternConfig) {
				defer wg.Done()
				rep := NewReplicator(e.local, e.peers, p, e.met, e.cfg.Sync.KeyspaceChannel)
				rep.Run(ctx)
			}(pattern)
		case "counter-delta":
			wg.Add(1)
			go func(p config.PatternConfig) {
				defer wg.Done()
				agg := NewAggregator(e.local, e.peers, p, e.met)
				agg.Run(ctx)
			}(pattern)
		default:
			log.Printf("engine: unknown semantics %q for pattern %q; skipping", pattern.Semantics, pattern.Prefix)
		}
	}
	e.mu.Lock()
	e.ready = true
	e.mu.Unlock()

	wg.Wait()
	log.Println("engine: all sync goroutines exited")
}

// Ready reports whether the engine is past initial setup and serving traffic.
func (e *Engine) Ready() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.ready
}

// ObservedLag returns the most-recent sync lag in seconds (for /readyz + alerting).
func (e *Engine) ObservedLag() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lagSeconds
}

// recordLag is called by replicators/aggregators when they measure round-trip lag.
func (e *Engine) recordLag(seconds float64) {
	e.mu.Lock()
	e.lagSeconds = seconds
	e.mu.Unlock()
}
