COPY . .

# Building Applications
ARG VERSION=dev
ARG COMMIT=1.0.0
ARG BUILD_TIME=2025-08-01T00:00:00Z

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
    -a -installsuffix cgo \
    -o api-service \
    ./cmd/server

# Runtime
FROM alpine:latest

LABEL version="0.1.3"

# Declare variables
ARG INFLUXDB_VERSION

# Install CA Certificates and Time Zone Data
RUN apk --no-cache add ca-certificates bash sqlite curl supervisor redis

# Create non-privileged user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser/

# Copy binary files from the build phase
COPY --from=builder /app/api-service .

# Copy configuration file (if it exists)
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/init_db.sh ./scripts/init_influxdb.sh ./scripts/healthcheck.sh

# Copy influxdb from the build phase
COPY --from=builder /app/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}/usr/bin/influxd /usr/local/bin/

# Create Data Directory
RUN mkdir -p /home/appuser/data /home/appuser/logs /home/appuser/data/influxdb && chown -R appuser:appuser /home/appuser/

# Copy the supervisord configuration file and health check script
COPY supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# Switch to a non-privileged user
USER appuser

# Exposed Ports
EXPOSE 8080 9090 6379 8086

# Launch the application
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]