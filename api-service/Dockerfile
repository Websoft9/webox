# Websoft9 API Service Dockerfile
# Multi-stage build for optimized image size and security

# Build stage - using base image with Go, Redis, and InfluxDB
FROM  websoft9dev/webox-base:latest AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN --mount=type=cache,target=${GOROOT}/pkg/mod \
  go mod download

# Copy Source code
COPY . .

# Generate swagger docs
RUN swag init -g cmd/server/main.go

# Building Application with build cache
ARG VERSION=dev
ARG COMMIT=1.0.0
ARG BUILD_TIME=2025-08-01T00:00:00Z

# Build with cache mount and parallel compilation
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=${GOROOT}/pkg/mod \
  CGO_ENABLED=1 GOOS=linux go build \
  -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
  -trimpath \
  -o api-service \
  ./cmd/server && \
  CGO_ENABLED=1 GOOS=linux go build \
  -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
  -trimpath \
  -o create-user \
  ./tools/create-user

# Runtime - using base image with all required services
FROM alpine:3.22

LABEL version="0.1.7"

# Install only runtime dependencies (no build tools)
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache add \
  ca-certificates \
  bash \
  sqlite \
  supervisor \
  tzdata \
  mysql-client \
  postgresql-client \
  curl && \
  rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

# Set timezone to UTC
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Configure environment variables
ENV REDIS_HOME=/usr/local/redis
ENV INFLUXDB_HOME=/usr/local/influxdb
ENV PATH=$REDIS_HOME/bin:$INFLUXDB_HOME/bin:$PATH

# Copy redis & influxdb from the build phase (as appuser)
COPY --from=builder --chown=appuser:appuser $REDIS_HOME/bin/redis-server $REDIS_HOME/bin/
COPY --from=builder --chown=appuser:appuser $REDIS_HOME/bin/redis-cli $REDIS_HOME/bin/
COPY --from=builder --chown=appuser:appuser /etc/redis.conf /etc/
COPY --from=builder --chown=appuser:appuser $INFLUXDB_HOME/bin/influxd $INFLUXDB_HOME/bin/

# Create necessary directories
RUN adduser -D -s /bin/sh appuser && \
  mkdir -p \
  /var/lib/redis \
  /var/log/redis \
  /var/lib/influxdb \
  /var/log/influxdb \
  /home/appuser/data \
  /home/appuser/logs \
  /home/appuser/data/influxdb && \
  chown -R appuser:appuser \
  /var/lib/redis \
  /var/log/redis \
  /var/lib/influxdb \
  /var/log/influxdb \
  /home/appuser

# Set working directory
WORKDIR /home/appuser/

# Copy binary & configuration files from the build phase (as appuser)
COPY --from=builder --chown=appuser:appuser /app/api-service .
COPY --from=builder --chown=appuser:appuser /app/configs ./configs
COPY --from=builder --chown=appuser:appuser /app/scripts ./scripts
COPY --from=builder --chown=appuser:appuser /app/create-user ./scripts

# Make scripts executable
RUN chmod +x ./scripts/*.sh

# Copy the supervisord configuration file and health check script
COPY supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# Switch to a non-privileged user
USER appuser

# Exposed Ports
EXPOSE 8080 9090 6379 8086

# Launch the application
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]
