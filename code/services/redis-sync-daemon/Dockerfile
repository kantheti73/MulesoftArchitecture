# Multi-stage build: small final image, no Go toolchain shipped.

# ---------- builder ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

# Dependencies cache layer
COPY go.mod ./
RUN go mod download

# Source
COPY cmd ./cmd
COPY internal ./internal

# Build static binary (no CGO)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w" \
    -o /out/redis-sync-daemon ./cmd/redis-sync-daemon

# ---------- runtime ----------
FROM gcr.io/distroless/static:nonroot AS runtime

# Distroless has no shell; daemon must be configured via env + mounted config

COPY --from=builder /out/redis-sync-daemon /usr/local/bin/redis-sync-daemon

# Documented expected mounts:
#   /etc/redis-sync-daemon/config.yaml
#   /etc/ssl/*.crt|key

EXPOSE 8080 8443 9090

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/redis-sync-daemon"]
CMD ["--config", "/etc/redis-sync-daemon/config.yaml"]
