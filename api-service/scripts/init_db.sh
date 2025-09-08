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

ADMIN_USERNAME="${WEBSOFT9_USERNAME:-system@websoft9.com}"
ADMIN_PASSWORD="${WEBSOFT9_PASSWORD:-changeme}"

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
    echo "  --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Initialize SQLite database"
    echo "  $0 -t mysql -u root -p password      # Initialize MySQL database"
    echo "  $0 -t sqlite -f /path/to/db.sqlite   # Initialize SQLite with custom path"
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
if [[ -f "$FLAG_FILE" ]]; then
    print_info "Database initialization already completed. Skipping..."
    exit 0
fi

# Validate database type
if [[ "$DB_TYPE" != "sqlite" && "$DB_TYPE" != "mysql" && "$DB_TYPE" != "postgres" ]]; then
    print_error "Unsupported database type: $DB_TYPE"
    print_info "Supported types: sqlite, mysql, postgres"
    exit 1
fi

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
        
        # Initialize SQLite database
        print_info "Executing SQLite initialization script..."
        if sqlite3 "$SQLITE_PATH" < scripts/init_sqlite.sql; then
            print_info "SQLite database initialized successfully!"
            print_info "Database file: $SQLITE_PATH"
            
            # Create admin user from environment variables
            print_info "Creating admin user: $ADMIN_USERNAME"
            
            # 使用SHA256生成密码哈希
            PASSWORD_HASH=$(echo -n "$ADMIN_PASSWORD" | sha256sum | cut -d ' ' -f 1)
            
            # 获取当前时间戳，格式为 YYYY-MM-DD HH:MM:SS
            CURRENT_TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")
    
            # 插入新用户
            # 明确指定ID为2
            sqlite3 "$SQLITE_PATH" <<EOF
-- 临时关闭外键约束
PRAGMA foreign_keys = OFF;

-- 插入新用户（明确指定ID为2）
INSERT OR IGNORE INTO users (id, username, email, password_hash, nickname, gender, signature, status, timezone, language, created_at, updated_at)
VALUES (2, 'admin2', '$ADMIN_USERNAME', '$PASSWORD_HASH', 'Manager', 0, 'Websoft9 manager', 1, 'Asia/Shanghai', 'zh-CN', '$CURRENT_TIMESTAMP', '$CURRENT_TIMESTAMP');

-- 将新用户关联到超级管理员角色（明确指定ID为2）
INSERT OR IGNORE INTO user_roles (id, user_id, role_id, granted_by, granted_at, status, created_at, updated_at)
VALUES (2, 2, 1, 1, '$CURRENT_TIMESTAMP', 1, '$CURRENT_TIMESTAMP', '$CURRENT_TIMESTAMP');

-- 确保角色拥有所有权限（如果需要）
-- 这部分不需要指定ID，因为role_permissions表可能没有ID字段
INSERT OR IGNORE INTO role_permissions (role_id, permission_code, granted_by, granted_at, status, created_at, updated_at)
SELECT 1, code, 1, '$CURRENT_TIMESTAMP', 1, '$CURRENT_TIMESTAMP', '$CURRENT_TIMESTAMP'
FROM permissions 
WHERE code NOT IN (SELECT permission_code FROM role_permissions WHERE role_id = 1);

-- 重新启用外键约束
PRAGMA foreign_keys = ON;
EOF

            print_info "Admin user created successfully with ID=2 and super_admin role!"
            
        else
            print_error "Failed to initialize SQLite database"
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
print_info ""
print_info "Next steps:"
print_info "1. Update your configuration file (configs/config.yaml)"
print_info "2. Start the Websoft9 service: ./api-service"
print_info "3. Access the web interface at: http://localhost:8080"