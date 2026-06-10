// Package main is the entry point for the Redis Sync Daemon.
//
// This daemon runs alongside a per-DC Redis Sentinel cluster and selectively
// replicates a subset of keys (idempotency cache + rate-limit counter deltas)
// to peer DCs over an mTLS HTTPS channel.
//
// See docs/22-redis-cross-dc-sync.md for the architecture rationale.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/config"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/health"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/metrics"
	"github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/peer"
	redisclient "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/redis"
	syncengine "github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon/internal/sync"
)

func main() {
	configPath := flag.String("config", "/etc/redis-sync-daemon/config.yaml", "Path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}
	log.Printf("redis-sync-daemon starting (local=%s peers=%d patterns=%d)",
		cfg.Local.SentinelMaster, len(cfg.Peers), len(cfg.Sync.Patterns))

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Prometheus metrics
	met := metrics.New()

	// Connect to local Redis (Sentinel-aware)
	localRedis, err := redisclient.NewSentinelClient(rootCtx, cfg.Local)
	if err != nil {
		log.Fatalf("local Redis connect failed: %v", err)
	}
	defer localRedis.Close()

	// HTTP clients to each peer DC
	peerClients := make([]*peer.Client, 0, len(cfg.Peers))
	for _, p := range cfg.Peers {
		c, err := peer.NewClient(p, met)
		if err != nil {
			log.Fatalf("peer client setup failed for %s: %v", p.Name, err)
		}
		peerClients = append(peerClients, c)
	}

	// Sync engine: orchestrates replicator + aggregator goroutines
	engine := syncengine.NewEngine(cfg, localRedis, peerClients, met)
	go engine.Run(rootCtx)

	// Peer HTTP server: receives sync messages from other DCs
	peerServer := peer.NewServer(cfg.Server, localRedis, met)
	go func() {
		if err := peerServer.ListenAndServe(rootCtx); err != nil {
			log.Printf("peer server stopped: %v", err)
			cancel()
		}
	}()

	// Health + metrics servers
	healthSrv := health.NewServer(cfg.Server.HealthAddr, localRedis, engine)
	go healthSrv.Run(rootCtx)
	metricsSrv := metrics.NewServer(cfg.Server.MetricsAddr)
	go metricsSrv.Run(rootCtx)

	// Wait for signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("received %s, shutting down", sig)

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	cancel()
	_ = peerServer.Shutdown(shutdownCtx)
	_ = healthSrv.Shutdown(shutdownCtx)
	_ = metricsSrv.Shutdown(shutdownCtx)
	log.Println("shutdown complete")
}
