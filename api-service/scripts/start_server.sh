#!/bin/bash

# Websoft9 API Service Container Startup Script
# This script orchestrates the complete startup sequence for the API service container
# including database initialization, user setup, and service startup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_DIR="/home/appuser/logs"
DATA_DIR="/home/appuser/data"

# Service check configuration
MAX_RETRIES=30
RETRY_INTERVAL=2

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function to handle errors and exit
handle_error() {
    local step="$1"
    local error_msg="$2"
    print_error "Failed at step: $step"
    print_error "$error_msg"
    print_error "Container startup aborted"
    exit 1
}

# Function to ensure required directories exist
ensure_directories() {
    print_info "Ensuring required directories exist..."

    local dirs=("$LOG_DIR" "$DATA_DIR" "$DATA_DIR/influxdb")

    for dir in "${dirs[@]}"; do
        if [[ ! -d "$dir" ]]; then
            print_info "Creating directory: $dir"
            mkdir -p "$dir" || handle_error "Directory Creation" "Failed to create directory: $dir"
        fi
    done

    print_info "All required directories are ready"
}

# Function to wait for database to be ready
wait_for_database() {
    print_info "Checking database readiness..."

    local db_type="${WEBSOFT9_DB_TYPE:-sqlite}"
    local db_host="${WEBSOFT9_DB_HOST:-localhost}"
    local db_port="${WEBSOFT9_DB_PORT}"

    case "$db_type" in
        sqlite)
            # For SQLite, just check if directory is writable
            if [[ -w "$DATA_DIR" ]]; then
                print_info "SQLite database directory is ready"
                return 0
            else
                handle_error "Database Check" "SQLite database directory is not writable: $DATA_DIR"
            fi
            ;;
        mysql)
            print_info "Skip MySQL database check"
            return 0
            ;;
        postgres)
            print_info "Skip PostgreSQL database check"
            return 0
            ;;
        *)
            handle_error "Database Check" "Unsupported database type: $db_type"
            ;;
    esac
}



# Function to initialize InfluxDB
initialize_influxdb() {
    local influxdb_url="$1"
    local flag_file="$DATA_DIR/.influxdb_initialized"

    # Check if already initialized
    if [[ -f "$flag_file" ]]; then
        print_info "InfluxDB already initialized, skipping..."
        return 0
    fi

    # Set default values
    local username="${WEBSOFT9_INFLUXDB_USERNAME:-admin}"
    local password="${WEBSOFT9_INFLUXDB_PASSWORD:-websoft9}"
    local org="${WEBSOFT9_INFLUXDB_ORG:-websoft9}"
    local bucket="${WEBSOFT9_INFLUXDB_BUCKET:-metrics}"
    local token="${WEBSOFT9_INFLUXDB_TOKEN:-mytoken}"

    print_info "Initializing InfluxDB with org: $org, bucket: $bucket"

    # Initialize InfluxDB (only works on first run)
    local response=$(curl -s -w "\n%{http_code}" -XPOST "${influxdb_url}/api/v2/setup" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\":\"${username}\",
            \"password\":\"${password}\",
            \"org\":\"${org}\",
            \"bucket\":\"${bucket}\",
            \"token\":\"${token}\",
            \"retentionPeriodSeconds\":0
        }")

    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')

    # HTTP 201 = success, 422 = already initialized
    if [[ "$http_code" == "201" ]]; then
        print_info "InfluxDB setup completed successfully"
        touch "$flag_file"
        return 0
    elif [[ "$http_code" == "422" ]]; then
        print_info "InfluxDB already initialized (onboarding already completed)"
        touch "$flag_file"
        return 0
    else
        print_warn "InfluxDB initialization returned HTTP $http_code"
        print_warn "Response: $body"
        print_warn "Continuing anyway, you may need to manually initialize InfluxDB"
        return 0
    fi
}

# Step 1: Start and initialize InfluxDB
start_influxdb() {
    print_step "Step 1: Starting InfluxDB service..."

    local influxdb_url="${WEBSOFT9_INFLUXDB_URL:-}"
    local use_external_influxdb=false

    # Check if external InfluxDB is configured
    if [[ -n "$influxdb_url" ]]; then
        print_info "Detected external InfluxDB configuration: $influxdb_url"
        print_info "Testing connection to external InfluxDB..."

        # Test external InfluxDB connection
        if curl -s "${influxdb_url}/health" | grep -q '"status":"pass"'; then
            print_info "External InfluxDB is accessible and ready"
            use_external_influxdb=true

            # Initialize external InfluxDB
            initialize_influxdb "$influxdb_url"
            return 0
        else
            print_warn "External InfluxDB is not accessible at $influxdb_url"
            print_warn "Falling back to local InfluxDB startup"
        fi
    fi

    # Start local InfluxDB if external InfluxDB is not configured or not accessible
    if [[ "$use_external_influxdb" == "false" ]]; then
        print_info "Starting local InfluxDB service..."

        local influxdb_cmd="/usr/local/influxdb/bin/influxd"
        local influxdb_opts="--bolt-path $DATA_DIR/influxdb/influxd.bolt --engine-path $DATA_DIR/influxdb/engine"

        # Check if InfluxDB binary exists
        if [[ ! -f "$influxdb_cmd" ]]; then
            handle_error "InfluxDB Startup" "InfluxDB binary not found: $influxdb_cmd"
        fi

        # Start InfluxDB in background
        print_info "Starting local InfluxDB..."
        nohup $influxdb_cmd $influxdb_opts > "$LOG_DIR/influxdb.log" 2>&1 &
        local influxdb_pid=$!

        print_info "InfluxDB started with PID: $influxdb_pid"

        # Wait for InfluxDB to be ready
        print_info "Waiting for local InfluxDB to be ready..."
        local retries=0
        while [[ $retries -lt $MAX_RETRIES ]]; do
            if curl -s http://localhost:8086/health | grep -q '"status":"pass"'; then
                print_info "Local InfluxDB is ready"
                break
            fi
            retries=$((retries + 1))
            sleep $RETRY_INTERVAL
        done

        if [[ $retries -ge $MAX_RETRIES ]]; then
            handle_error "InfluxDB Startup" "Local InfluxDB failed to start within timeout"
        fi

        # Initialize local InfluxDB
        initialize_influxdb "http://localhost:8086"
        return 0
    fi
}

# Step 2: Start Redis
start_redis() {
    print_step "Step 2: Starting Redis service..."

    local redis_host="${WEBSOFT9_REDIS_HOST:-}"
    local redis_port="${WEBSOFT9_REDIS_PORT:-}"
    local redis_password="${WEBSOFT9_REDIS_PASSWORD:-}"
    local use_external_redis=false

    # Check if external Redis is configured
    if [[ -n "$redis_host" ]] && [[ -n "$redis_port" ]]; then
        print_info "Detected external Redis configuration: ${redis_host}:${redis_port}"
        print_info "Testing connection to external Redis..."

        # Test external Redis connection
        local test_cmd="redis-cli -h $redis_host -p $redis_port"
        if [[ -n "$redis_password" ]]; then
            test_cmd="$test_cmd -a $redis_password"
        fi

        if $test_cmd ping 2>/dev/null | grep -q "PONG"; then
            print_info "External Redis is accessible and ready"
            use_external_redis=true
            return 0
        else
            print_warn "External Redis is not accessible at ${redis_host}:${redis_port}"
            print_warn "Falling back to local Redis startup"
        fi
    fi

    # Start local Redis if external Redis is not configured or not accessible
    if [[ "$use_external_redis" == "false" ]]; then
        print_info "Starting local Redis service..."

        local redis_cmd="/usr/local/redis/bin/redis-server"
        local redis_conf="${WEBSOFT9_REDIS_CONF:-/etc/redis.conf}"

        # Check if Redis binary exists
        if [[ ! -f "$redis_cmd" ]]; then
            handle_error "Redis Startup" "Redis binary not found: $redis_cmd"
        fi

        # Check if config file exists
        if [[ -f "$redis_conf" ]]; then
            print_info "Using Redis config file: $redis_conf"

            # Start Redis with config file
            if [[ -n "$redis_password" ]]; then
                print_info "Starting Redis with config file and password override"
                nohup $redis_cmd "$redis_conf" --requirepass "$redis_password" > "$LOG_DIR/redis.log" 2>&1 &
            else
                print_info "Starting Redis with config file (using config file password if set)"
                nohup $redis_cmd "$redis_conf" > "$LOG_DIR/redis.log" 2>&1 &
            fi
        else
            print_warn "Redis config file not found: $redis_conf"
            print_info "Starting Redis with default configuration..."

            # Start Redis with basic options
            local redis_opts="--bind 0.0.0.0 --port 6379"
            if [[ -n "$redis_password" ]]; then
                redis_opts="$redis_opts --requirepass $redis_password"
                print_info "Redis password authentication enabled"
            else
                print_warn "No password set for Redis (not recommended for production)"
            fi

            nohup $redis_cmd $redis_opts > "$LOG_DIR/redis.log" 2>&1 &
        fi

        local redis_pid=$!
        print_info "Redis started with PID: $redis_pid"

        # Wait for Redis to be ready
        print_info "Waiting for local Redis to be ready..."
        local retries=0
        while [[ $retries -lt $MAX_RETRIES ]]; do
            local ping_cmd="redis-cli -h localhost -p 6379"
            if [[ -n "$redis_password" ]]; then
                ping_cmd="$ping_cmd -a $redis_password"
            fi

            if $ping_cmd ping 2>/dev/null | grep -q "PONG"; then
                print_info "Local Redis is ready"
                return 0
            fi
            retries=$((retries + 1))
            sleep $RETRY_INTERVAL
        done

        handle_error "Redis Startup" "Local Redis failed to start within timeout"
    fi
}

# Step 3: Verify all services and start API service
start_api_service() {
    print_step "Step 3: Verifying services and starting API service..."

    # Verify database
    print_info "Verifying database service..."
    wait_for_database

    # Verify InfluxDB
    print_info "Verifying InfluxDB service..."
    local influxdb_url="${WEBSOFT9_INFLUXDB_URL:-http://localhost:8086}"

    if ! curl -s "${influxdb_url}/health" | grep -q '"status":"pass"'; then
        handle_error "Service Verification" "InfluxDB is not running at $influxdb_url"
    fi
    print_info "InfluxDB service is running at $influxdb_url"

    # Verify Redis
    print_info "Verifying Redis service..."
    local redis_host="${WEBSOFT9_REDIS_HOST:-localhost}"
    local redis_port="${WEBSOFT9_REDIS_PORT:-6379}"
    local redis_password="${WEBSOFT9_REDIS_PASSWORD:-}"

    local verify_cmd="redis-cli -h $redis_host -p $redis_port"
    if [[ -n "$redis_password" ]]; then
        verify_cmd="$verify_cmd -a $redis_password"
    fi

    if ! $verify_cmd ping 2>/dev/null | grep -q "PONG"; then
        handle_error "Service Verification" "Redis is not running at ${redis_host}:${redis_port}"
    fi
    print_info "Redis service is running at ${redis_host}:${redis_port}"

    # All services are ready, start API service
    print_info "All services are ready, starting API service..."
    print_info "=========================================="
    print_info "Websoft9 API Service Container Started"
    print_info "=========================================="

    # Execute API service (this will replace the current process)
    exec /home/appuser/api-service
}

# Main execution flow
main() {
    print_info "=========================================="
    print_info "Websoft9 API Service Container Startup"
    print_info "=========================================="

    # Ensure required directories exist
    ensure_directories

    # Execute startup steps in sequence
    start_influxdb
    start_redis
    start_api_service
}

# Trap to handle script interruption
trap 'print_error "Startup interrupted by user"; exit 130' INT TERM

# Run main function
main "$@"
