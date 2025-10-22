# Websoft9 API Service Dockerfile
# Multi-stage build for optimized image size and security

# Define InfluxDB version variable
ARG INFLUXDB_VERSION=2.7.12

# Build stage
FROM golang:1.24-alpine AS builder

# Declare variables
ARG INFLUXDB_VERSION

# Install necessary tools
#RUN apk add --no-cache git bash gcc musl-dev ca-certificates curl
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache git bash gcc musl-dev ca-certificates curl

# Create non-privileged user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Download and extract InfluxDB binary
RUN curl -L https://dl.influxdata.com/influxdb/releases/v${INFLUXDB_VERSION}/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}_linux_amd64.tar.gz | tar -xz

# Set Go module proxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Copy go mod files
COPY go.mod go.sum ./

# Set Go proxy for faster dependency resolution
ENV GOPROXY=https://goproxy.cn,direct

# Download dependencies
RUN go mod download

# Download swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy Source code
COPY . .

# Generate swagger docs
RUN swag init -g cmd/server/main.go

# Building Application
ARG VERSION=dev
ARG COMMIT=1.0.0
ARG BUILD_TIME=2025-08-01T00:00:00Z

# Build both binaries in a single RUN command to leverage build cache
RUN CGO_ENABLED=1 GOOS=linux go build \
  -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
  -o api-service \
  ./cmd/server && \
  CGO_ENABLED=1 GOOS=linux go build \
  -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
  -o create-user \
  ./tools/create-user

# Runtime
FROM alpine:latest

LABEL version="0.1.6"

# Declare variables
ARG INFLUXDB_VERSION

# Install CA Certificates and Time zone Data
RUN apk --no-cache add ca-certificates bash sqlite curl supervisor redis tzdata

# Set timezone to UTC
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create non-privileged user and prepare directories
RUN adduser -D -s /bin/sh appuser && \
  mkdir -p /home/appuser/data /home/appuser/logs /home/appuser/data/influxdb && \
  chown -R appuser:appuser /home/appuser

# Set working directory
WORKDIR /home/appuser/

# Copy binary files from the build phase (as appuser)
COPY --from=builder --chown=appuser:appuser /app/api-service .

# Copy configuration files (as appuser)
COPY --from=builder --chown=appuser:appuser /app/configs ./configs
COPY --from=builder --chown=appuser:appuser /app/scripts ./scripts
COPY --from=builder --chown=appuser:appuser /app/create-user ./scripts

# Make scripts executable
RUN chmod +x ./scripts/*.sh

# Copy influxdb from the build phase
COPY --from=builder /app/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}/usr/bin/influxd /usr/local/bin/

# Copy the supervisord configuration file and health check script
COPY supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# Switch to a non-privileged user
USER appuser

# Exposed Ports
EXPOSE 8080 9090 6379 8086

# Launch the application
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]
