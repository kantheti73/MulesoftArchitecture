package sync

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/peer"
	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
)

// Aggregator handles counter-delta semantics keys (rate-limit counters).
//
// On a periodic timer, it scans the local key space for matching prefix,
// computes the delta since the last cycle, and ships the delta to all peers.
// Peers apply via INCRBY so concurrent updates on both sides add cleanly.
type Aggregator struct {
	local    *redisclient.Client
	peers    []*peer.Client
	pattern  config.PatternConfig
	met      *metrics.Metrics

	prevSnapshot map[string]int64 // last-observed values, per key
}

func NewAggregator(local *redisclient.Client, peers []*peer.Client, p config.PatternConfig, met *metrics.Metrics) *Aggregator {
	return &Aggregator{
		local:        local,
		peers:        peers,
		pattern:      p,
		met:          met,
		prevSnapshot: make(map[string]int64),
	}
}

// Run blocks until ctx is cancelled.
func (a *Aggregator) Run(ctx context.Context) {
	interval := time.Duration(a.pattern.AggregationIntervalSeconds) * time.Second
	log.Printf("aggregator: starting for prefix=%q interval=%s", a.pattern.Prefix, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.cycle(ctx)
		}
	}
}

func (a *Aggregator) cycle(ctx context.Context) {
	start := time.Now()
	pattern := a.pattern.Prefix + "*"

	// SCAN avoids blocking the Redis server (vs KEYS which is dangerous in prod)
	iter := a.local.Raw().Scan(ctx, 0, pattern, 1000).Iterator()
	currentSnapshot := make(map[string]int64)

	for iter.Next(ctx) {
		key := iter.Val()
		val, err := a.local.Raw().Get(ctx, key).Int64()
		if err != nil {
			continue
		}
		currentSnapshot[key] = val

		prev := a.prevSnapshot[key]
		delta := val - prev
		if delta <= 0 {
			continue // no progress to ship; counter may have rolled at TTL
		}

		var ttlSec int64
		if a.pattern.TTLPreserve {
			ttl, err := a.local.Raw().TTL(ctx, key).Result()
			if err == nil && ttl > 0 {
				ttlSec = int64(ttl.Seconds())
			}
		}

		for _, p := range a.peers {
			// Strip the local prefix and re-prefix consistently so peers
			// use the same key
			suffix := strings.TrimPrefix(key, a.pattern.Prefix)
			peerKey := a.pattern.Prefix + suffix

			if err := p.SendCounterDelta(ctx, peerKey, delta, ttlSec); err != nil {
				log.Printf("aggregator: peer %s SendCounterDelta failed for %s: %v", p.Name(), key, err)
				a.met.ReplicationErrors.WithLabelValues(a.pattern.Prefix, "peer-send").Inc()
				continue
			}
			a.met.KeysReplicated.WithLabelValues(a.pattern.Prefix, p.Name()).Inc()
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("aggregator: SCAN failed: %v", err)
		a.met.ReplicationErrors.WithLabelValues(a.pattern.Prefix, "scan").Inc()
	}

	a.prevSnapshot = currentSnapshot
	a.met.AggregationCycleDuration.WithLabelValues(a.pattern.Prefix).Observe(time.Since(start).Seconds())
}
