# Websoft9 API Service Dockerfile
# Multi-stage build: compile with golang, run with slim base

# ============================================
# Build stage
# ============================================
FROM golang:1.24-alpine3.22 AS builder

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME
ARG GOPROXY=https://goproxy.cn,direct
ARG GONOPROXY=
ARG GONOSUMDB=*

ENV GOPROXY=${GOPROXY} \
  GONOPROXY=${GONOPROXY} \
  GONOSUMDB=${GONOSUMDB} \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

# Install build dependencies for CGO (required by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Install swag for swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Copy go mod files for better layer caching
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init -g cmd/server/main.go -o ./docs

# Build create-user tool (rarely changes, better cache)
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go build -ldflags="-w -s" -trimpath -o create-user ./tools/create-user

# Build api-service with parallel compilation
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go build -p $(nproc) \
  -ldflags="-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
  -trimpath -o api-service ./cmd/server

# ============================================
# Runtime stage
# ============================================
FROM websoft9-base:latest

LABEL version="0.1.8" \
  maintainer="Websoft9" \
  description="Websoft9 API Service"

WORKDIR /home/appuser

# Copy supervisord config (needs root permission)
COPY --from=builder /app/supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# Copy binaries with correct permissions
COPY --from=builder --chown=appuser:appuser --chmod=755 /app/api-service ./
COPY --from=builder --chown=appuser:appuser /app/docs ./docs/
COPY --from=builder --chown=appuser:appuser /app/configs ./configs/
COPY --from=builder --chown=appuser:appuser /app/scripts ./scripts/
COPY --from=builder --chown=appuser:appuser --chmod=755 /app/create-user ./scripts/create-user

# Set executable permissions for scripts
RUN chmod +x ./scripts/*.sh

USER appuser

EXPOSE 8080 9090

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]
