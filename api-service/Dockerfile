# Websoft9 API Service Dockerfile
# Multi-stage build for optimized image size and security

# Define InfluxDB version variable
ARG INFLUXDB_VERSION=2.7.12

# Build stage
FROM golang:1.24-alpine AS builder

# Declare variables
ARG INFLUXDB_VERSION

# Install necessary tools
RUN apk add --no-cache git bash gcc musl-dev ca-certificates curl

# Create non-privileged user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Download and extract InfluxDB binary
RUN curl -L https://dl.influxdata.com/influxdb/releases/v${INFLUXDB_VERSION}/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}_linux_amd64.tar.gz | tar -xz

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
ARG VERSION=dev
ARG COMMIT=1.0.0
ARG BUILD_TIME=2025-08-01T00:00:00Z

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
    -a -installsuffix cgo \
    -o api-service \
    .

# Runtime stage
FROM alpine:latest

LABEL version="0.1.3"

# Declare variables
ARG INFLUXDB_VERSION

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates bash sqlite curl supervisor redis

# Create non-privileged user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser/

# Copy binary files from build stage
COPY --from=builder /app/api-service .

# Copy configuration files (if exists)
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/init_db.sh ./scripts/init_influxdb.sh ./scripts/healthcheck.sh

# Copy influxd from build stage
COPY --from=builder /app/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}/usr/bin/influxd /usr/local/bin/

# Create data directories
RUN mkdir -p /home/appuser/data /home/appuser/logs /home/appuser/data/influxdb && chown -R appuser:appuser /home/appuser/

# Copy supervisord configuration file and health check script
COPY supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# Switch to non-privileged user
USER appuser

# Expose ports
EXPOSE 8080 9090 6379 8086

# Start application
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]
