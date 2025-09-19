#!/bin/bash

# Websoft9 Database Initialization Script
# This script helps initialize the database with the appropriate SQL script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
DB_TYPE="sqlite"
DB_HOST="localhost"
DB_PORT=""
DB_NAME="websoft9"
DB_USER=""
DB_PASS=""
SQLITE_PATH="./data/websoft9.db"
INIT_SQL_FILE=""
FORCE=false

# Path to the flag file
FLAG_FILE="./data/.websoft9_db_initialized"

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --type TYPE        Database type (sqlite|mysql|postgres) [default: sqlite]"
    echo "  -h, --host HOST        Database host [default: localhost]"
    echo "  -P, --port PORT        Database port"
    echo "  -d, --database NAME    Database name [default: websoft9]"
    echo "  -u, --user USER        Database username"
    echo "  -p, --password PASS    Database password"
    echo "  -f, --file PATH        SQLite database file path [default: ./data/websoft9.db]"
    echo "      --force            Force initialization, ignore existing flag file and remove SQLite database"
    echo "      --init SQL_FILE    Initialize database and import data from specified SQL file"
    echo "      --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Initialize SQLite database"
    echo "  $0 -t mysql -u root -p password      # Initialize MySQL database"
    echo "  $0 -t sqlite -f /path/to/db.sqlite   # Initialize SQLite with custom path"
    echo "  $0 --init /path/to/data.sql          # Initialize and import data from SQL file"
    echo "  $0 -t mysql -u root -p pwd --init data.sql  # Initialize MySQL and import data"
    echo "  $0 --force                           # Force re-initialization (ignore flag file)"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            DB_TYPE="$2"
            shift 2
            ;;
        -h|--host)
            DB_HOST="$2"
            shift 2
            ;;
        -P|--port)
            DB_PORT="$2"
            shift 2
            ;;
        -d|--database)
            DB_NAME="$2"
            shift 2
            ;;
        -u|--user)
            DB_USER="$2"
            shift 2
            ;;
        -p|--password)
            DB_PASS="$2"
            shift 2
            ;;
        -f|--file)
            SQLITE_PATH="$2"
            shift 2
            ;;
        --force)
            FORCE=true
            shift 1
            ;;
        --init)
            INIT_SQL_FILE="$2"
            shift 2
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check if the script has already been executed
if [[ -f "$FLAG_FILE" && "$FORCE" != "true" ]]; then
    print_info "Database initialization already completed. Skipping..."
    print_info "Use --force to reinitialize the database."
    exit 0
fi

# Handle force mode
if [[ "$FORCE" == "true" ]]; then
    print_warn "Force mode enabled - will reinitialize database"

    # Remove flag file if it exists
    if [[ -f "$FLAG_FILE" ]]; then
        print_info "Removing existing flag file: $FLAG_FILE"
        rm -f "$FLAG_FILE"
    fi

    # For SQLite, remove existing database file
    if [[ "$DB_TYPE" == "sqlite" && -f "$SQLITE_PATH" ]]; then
        print_warn "Removing existing SQLite database: $SQLITE_PATH"
        rm -f "$SQLITE_PATH"
    fi
fi

# Validate database type
if [[ "$DB_TYPE" != "sqlite" && "$DB_TYPE" != "mysql" && "$DB_TYPE" != "postgres" ]]; then
    print_error "Unsupported database type: $DB_TYPE"
    print_info "Supported types: sqlite, mysql, postgres"
    exit 1
fi

# Validate init SQL file if specified
if [[ -n "$INIT_SQL_FILE" ]]; then
    if [[ ! -f "$INIT_SQL_FILE" ]]; then
        print_error "SQL file not found: $INIT_SQL_FILE"
        exit 1
    fi
    if [[ ! -r "$INIT_SQL_FILE" ]]; then
        print_error "SQL file is not readable: $INIT_SQL_FILE"
        exit 1
    fi
    print_info "SQL file for data import: $INIT_SQL_FILE"
fi

# Backup and rollback functions
create_backup() {
    local backup_path="$1"
    local source_path="$2"

    if [[ "$DB_TYPE" == "sqlite" ]]; then
        if [[ -f "$source_path" ]]; then
            print_info "Creating backup of existing database..."
            cp "$source_path" "$backup_path"
            return $?
        fi
    fi
    return 0
}

restore_backup() {
    local backup_path="$1"
    local target_path="$2"

    if [[ "$DB_TYPE" == "sqlite" && -f "$backup_path" ]]; then
        print_warn "Restoring database from backup due to import failure..."
        cp "$backup_path" "$target_path"
        rm -f "$backup_path"
        return $?
    fi
    return 0
}

cleanup_backup() {
    local backup_path="$1"

    if [[ -f "$backup_path" ]]; then
        print_info "Cleaning up backup file..."
        rm -f "$backup_path"
    fi
}

# Import SQL data function
import_sql_data() {
    local sql_file="$1"
    local backup_file="$2"

    print_info "Importing data from SQL file: $sql_file"

    case $DB_TYPE in
        sqlite)
            print_info "Disabling foreign key constraints for SQLite import..."
            {
                echo "PRAGMA foreign_keys=OFF;"
                cat "$sql_file"
                echo "PRAGMA foreign_keys=ON;"
            } | sqlite3 "$SQLITE_PATH"

            if [ $? -ne 0 ]; then
                print_error "Failed to import SQL data from: $sql_file"
                restore_backup "$backup_file" "$SQLITE_PATH"
                return 1
            fi
            ;;
        mysql)
            print_info "Disabling foreign key constraints for MySQL import..."
            {
                echo "SET FOREIGN_KEY_CHECKS=0;"
                cat "$sql_file"
                echo "SET FOREIGN_KEY_CHECKS=1;"
            } | $MYSQL_CMD "$DB_NAME"

            if [ $? -ne 0 ]; then
                print_error "Failed to import SQL data from: $sql_file"
                print_warn "Please manually restore your database if needed"
                # Try to re-enable foreign key checks even on failure
                echo "SET FOREIGN_KEY_CHECKS=1;" | $MYSQL_CMD "$DB_NAME" 2>/dev/null || true
                return 1
            fi
            ;;
        postgres)
            print_info "Disabling foreign key constraints for PostgreSQL import..."
            # Create a temporary SQL file with constraint management
            local temp_sql="/tmp/postgres_import_$$.sql"
            {
                echo "SET session_replication_role = replica;"
                cat "$sql_file"
                echo "SET session_replication_role = DEFAULT;"
            } > "$temp_sql"

            if ! $PSQL_CMD -d "$DB_NAME" -f "$temp_sql"; then
                print_error "Failed to import SQL data from: $sql_file"
                print_warn "Please manually restore your database if needed"
                # Try to re-enable foreign key checks even on failure
                echo "SET session_replication_role = DEFAULT;" | $PSQL_CMD -d "$DB_NAME" 2>/dev/null || true
                rm -f "$temp_sql"
                return 1
            fi

            # Clean up temporary file
            rm -f "$temp_sql"
            ;;
    esac

    print_info "SQL data imported successfully!"
    return 0
}

# Set default ports if not specified
if [[ -z "$DB_PORT" ]]; then
    case $DB_TYPE in
        mysql)
            DB_PORT="3306"
            ;;
        postgres)
            DB_PORT="5432"
            ;;
        sqlite)
            DB_PORT=""
            ;;
    esac
fi

print_info "Initializing Websoft9 database..."
print_info "Database type: $DB_TYPE"

case $DB_TYPE in
    sqlite)
        print_info "SQLite database path: $SQLITE_PATH"

        # Create directory if it doesn't exist
        DB_DIR=$(dirname "$SQLITE_PATH")
        if [[ ! -d "$DB_DIR" ]]; then
            print_info "Creating directory: $DB_DIR"
            mkdir -p "$DB_DIR"
        fi

        # Check if SQLite is available
        if ! command -v sqlite3 &> /dev/null; then
            print_error "sqlite3 command not found. Please install SQLite3."
            exit 1
        fi

        # Create backup if init SQL file is specified and database exists
        BACKUP_FILE=""
        if [[ -n "$INIT_SQL_FILE" ]]; then
            BACKUP_FILE="${SQLITE_PATH}.backup.$(date +%Y%m%d_%H%M%S)"
            create_backup "$BACKUP_FILE" "$SQLITE_PATH"
        fi

        # Initialize SQLite database
        print_info "Executing SQLite initialization script..."
        if sqlite3 "$SQLITE_PATH" < scripts/init_sqlite.sql; then
            print_info "SQLite database initialized successfully!"
            print_info "Database file: $SQLITE_PATH"

            # Import additional data if specified
            if [[ -n "$INIT_SQL_FILE" ]]; then
                if import_sql_data "$INIT_SQL_FILE" "$BACKUP_FILE"; then
                    cleanup_backup "$BACKUP_FILE"
                else
                    print_error "Database initialization completed but SQL import failed"
                    exit 1
                fi
            fi
        else
            print_error "Failed to initialize SQLite database"
            cleanup_backup "$BACKUP_FILE"
            exit 1
        fi
        ;;

    mysql)
        print_info "MySQL connection: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"

        # Check if MySQL client is available
        if ! command -v mysql &> /dev/null; then
            print_error "mysql command not found. Please install MySQL client."
            exit 1
        fi

        # Validate required parameters
        if [[ -z "$DB_USER" ]]; then
            print_error "MySQL username is required. Use -u or --user option."
            exit 1
        fi

        # Prepare MySQL connection parameters
        MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER"
        if [[ -n "$DB_PASS" ]]; then
            MYSQL_CMD="$MYSQL_CMD -p$DB_PASS"
        fi

        # Test MySQL connection
        print_info "Testing MySQL connection..."
        if ! echo "SELECT 1;" | $MYSQL_CMD > /dev/null 2>&1; then
            print_error "Failed to connect to MySQL server"
            print_info "Please check your connection parameters"
            exit 1
        fi

        # Create database if it doesn't exist
        print_info "Creating database if not exists: $DB_NAME"
        echo "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" | $MYSQL_CMD

        # Initialize MySQL database
        print_info "Executing MySQL initialization script..."
        if $MYSQL_CMD "$DB_NAME" < scripts/init_mysql.sql; then
            print_info "MySQL database initialized successfully!"
            print_info "Database: $DB_NAME"

            # Import additional data if specified
            if [[ -n "$INIT_SQL_FILE" ]]; then
                if import_sql_data "$INIT_SQL_FILE" ""; then
                    print_info "Data import completed successfully!"
                else
                    print_error "Database initialization completed but SQL import failed"
                    exit 1
                fi
            fi
        else
            print_error "Failed to initialize MySQL database"
            exit 1
        fi
        ;;

    postgres)
        print_info "PostgreSQL connection: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"

        # Check if PostgreSQL client is available
        if ! command -v psql &> /dev/null; then
            print_error "psql command not found. Please install PostgreSQL client."
            exit 1
        fi

        # Validate required parameters
        if [[ -z "$DB_USER" ]]; then
            print_error "PostgreSQL username is required. Use -u or --user option."
            exit 1
        fi

        # Prepare PostgreSQL connection parameters
        PSQL_CMD="psql -h $DB_HOST -p $DB_PORT -U $DB_USER"
        if [[ -n "$DB_PASS" ]]; then
            export PGPASSWORD="$DB_PASS"
        fi

        # Test PostgreSQL connection
        print_info "Testing PostgreSQL connection..."
        if ! echo "\q" | $PSQL_CMD > /dev/null 2>&1; then
            print_error "Failed to connect to PostgreSQL server"
            print_info "Please check your connection parameters"
            exit 1
        fi

        # Create database if it doesn't exist
        print_info "Creating database if not exists: $DB_NAME"
        echo "CREATE DATABASE \"$DB_NAME\";" | $PSQL_CMD || print_warn "Database $DB_NAME may already exist."

        # Initialize PostgreSQL database
        print_info "Executing PostgreSQL initialization script..."
        if $PSQL_CMD -d "$DB_NAME" -f scripts/init_postgres.sql; then
            print_info "PostgreSQL database initialized successfully!"
            print_info "Database: $DB_NAME"

            # Import additional data if specified
            if [[ -n "$INIT_SQL_FILE" ]]; then
                if import_sql_data "$INIT_SQL_FILE" ""; then
                    print_info "Data import completed successfully!"
                else
                    print_error "Database initialization completed but SQL import failed"
                    exit 1
                fi
            fi
        else
            print_error "Failed to initialize PostgreSQL database"
            exit 1
        fi
        ;;
esac

# At the end of the script, create the flag file
print_info "Creating flag file to indicate initialization is complete: $FLAG_FILE"
mkdir -p "$(dirname "$FLAG_FILE")"
touch "$FLAG_FILE"

print_info "Database initialization completed!"
if [[ -n "$INIT_SQL_FILE" ]]; then
    print_info "SQL data import completed from: $INIT_SQL_FILE"
fi
print_info ""
print_info "Next steps:"
print_info "1. Update your configuration file (configs/config.yaml)"
print_info "2. Start the Websoft9 service: ./api-service"
print_info "3. Access the web interface at: http://localhost:8080"
