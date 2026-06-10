package sync

import (
	"context"
	"log"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v9"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/peer"
	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
)

// Replicator handles SETNX-semantics keys (idempotency cache).
//
// It subscribes to Redis keyspace notifications, filters by prefix, fetches
// the current value + TTL, and fans out to peer DCs via HTTP. Replication is
// at-least-once and idempotent: peer applies via SETNX so duplicate deliveries
// are no-ops.
type Replicator struct {
	local           *redisclient.Client
	peers           []*peer.Client
	pattern         config.PatternConfig
	met             *metrics.Metrics
	keyspaceChannel string
}

func NewReplicator(local *redisclient.Client, peers []*peer.Client, p config.PatternConfig, met *metrics.Metrics, ch string) *Replicator {
	return &Replicator{local: local, peers: peers, pattern: p, met: met, keyspaceChannel: ch}
}

// Run blocks until ctx is cancelled.
func (r *Replicator) Run(ctx context.Context) {
	log.Printf("replicator: starting for prefix=%q semantics=setnx", r.pattern.Prefix)

	pubsub := r.local.Raw().PSubscribe(ctx, r.keyspaceChannel)
	defer pubsub.Close()

	ch := pubsub.Channel(goredis.WithChannelSize(1000))

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				log.Printf("replicator: pubsub channel closed for prefix=%q", r.pattern.Prefix)
				return
			}
			r.handleEvent(ctx, msg.Channel, msg.Payload)
		}
	}
}

func (r *Replicator) handleEvent(ctx context.Context, channel, eventType string) {
	// Channel format: __keyspace@<db>__:<key>
	idx := strings.Index(channel, ":")
	if idx < 0 {
		return
	}
	key := channel[idx+1:]
	if !strings.HasPrefix(key, r.pattern.Prefix) {
		return
	}

	switch eventType {
	case "set", "setex":
		// New / refreshed key — fetch + replicate
		r.replicateSet(ctx, key)
	case "expired", "del":
		// Local expiry; peer will expire independently via its own TTL.
		// Skip explicit delete propagation to avoid race conditions.
	default:
		// Other events (e.g., DEL, RENAME) are noise for idempotency keys
	}
}

func (r *Replicator) replicateSet(ctx context.Context, key string) {
	value, err := r.local.Raw().Get(ctx, key).Result()
	if err == goredis.Nil {
		return // gone before we got to it; fine
	}
	if err != nil {
		log.Printf("replicator: GET %s failed: %v", key, err)
		r.met.ReplicationErrors.WithLabelValues(r.pattern.Prefix, "local-get").Inc()
		return
	}

	var ttlSec int64
	if r.pattern.TTLPreserve {
		ttl, err := r.local.Raw().TTL(ctx, key).Result()
		if err == nil && ttl > 0 {
			ttlSec = int64(ttl.Seconds())
		}
	}

	start := time.Now()
	for _, p := range r.peers {
		if err := p.SendSet(ctx, key, value, ttlSec); err != nil {
			log.Printf("replicator: peer %s SendSet failed for %s: %v", p.Name(), key, err)
			r.met.ReplicationErrors.WithLabelValues(r.pattern.Prefix, "peer-send").Inc()
			continue
		}
		r.met.KeysReplicated.WithLabelValues(r.pattern.Prefix, p.Name()).Inc()
	}
	r.met.ReplicationLatency.WithLabelValues(r.pattern.Prefix).Observe(time.Since(start).Seconds())
}
