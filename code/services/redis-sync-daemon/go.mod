module github.com/kantheti73/MulesoftArchitecture/code/services/redis-sync-daemon

go 1.22

require (
	github.com/go-redis/redis/v9 v9.5.1
	github.com/prometheus/client_golang v1.18.0
	gopkg.in/yaml.v3 v3.0.1
)

// Minimal dependency set: Redis client + Prometheus metrics + YAML config.
// Avoid pulling heavier frameworks (gin/echo/etc.) — the HTTP surface is tiny
// and net/http stdlib is sufficient.
