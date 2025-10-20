package database

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/pkg/database/plugins"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Database type constants
const (
	DatabaseTypeSQLite     = "sqlite"
	DatabaseTypeMySQL      = "mysql"
	DatabaseTypePostgres   = "postgres"
	DatabaseTypePostgreSQL = "postgresql"
)

// parseLogLevel converts string log level to gorm logger level
func parseLogLevel(level string) gormlogger.LogLevel {
	switch strings.ToLower(level) {
	case "error":
		return gormlogger.Error
	case "warn", "warning":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	case "debug":
		return gormlogger.Info // GORM doesn't have debug level, use Info
	default:
		return gormlogger.Info // Default to Info level
	}
}

// DatabaseConfig contains database connection parameters
type DatabaseConnectionConfig struct {
	Type            string
	Path            string
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	ConnectTimeout  int
	Charset         string
	Timezone        string
	LogLevel        string
}

// InitDB initializes database connection based on configuration
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Convert config to connection config
	dbConfig := &DatabaseConnectionConfig{
		Type:            cfg.Database.Type,
		Path:            cfg.Database.Path,
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Database:        cfg.Database.Database,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		SSLMode:         cfg.Database.SSLMode,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnectTimeout:  cfg.Database.ConnectTimeout,
		Charset:         cfg.Database.Charset,
		Timezone:        cfg.Database.Timezone,
		LogLevel:        cfg.Server.Log.LogLevel,
	}

	// Initialize database connection based on type
	switch strings.ToLower(cfg.Database.Type) {
	case DatabaseTypeSQLite:
		db, err = initSQLite(dbConfig)
	case DatabaseTypeMySQL:
		db, err = initMySQL(dbConfig)
	case DatabaseTypePostgres, DatabaseTypePostgreSQL:
		db, err = initPostgreSQL(dbConfig)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s database: %v", cfg.Database.Type, err)
	}

	// Register custom serializers
	schema.RegisterSerializer("datetime", plugins.DateTimeSerializer{})

	// Configure connection pool (SQLite already configured in initSQLite)
	if cfg.Database.Type != DatabaseTypeSQLite {
		if err := configureConnectionPool(db, dbConfig); err != nil {
			return nil, fmt.Errorf("failed to configure connection pool: %v", err)
		}
	}

	return db, nil
}

// initSQLite initializes basic SQLite database connection
func initSQLite(cfg *DatabaseConnectionConfig) (*gorm.DB, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dbDir, constants.DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Basic SQLite DSN (optimizations will be handled by SQLiteManager)
	dsn := cfg.Path + "?_foreign_keys=ON"
	// Add timezone parameter if configured
	if cfg.Timezone != "" {
		// SQLite uses _loc parameter for timezone
		dsn += fmt.Sprintf("&_loc=%s", url.PathEscape(cfg.Timezone))
	}

	// Basic GORM configuration
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
		// Enable prepared statements to improve performance
		PrepareStmt: true,
		// Disable default transaction, manually manage transactions for better control
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %v", err)
	}

	// Basic connection pool configuration (SQLiteManager will optimize further)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(1) // SQLite single writer constraint
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return db, nil
}

// initMySQL initializes MySQL database connection
func initMySQL(cfg *DatabaseConnectionConfig) (*gorm.DB, error) {
	// Build MySQL DSN
	dsn := buildMySQLDSN(cfg)

	// Open MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %v", err)
	}

	// Set session timezone if configured
	if cfg.Timezone != "" {
		// Convert timezone format: UTC -> +00:00, or use timezone name directly
		timezoneValue := cfg.Timezone
		if cfg.Timezone == "UTC" {
			timezoneValue = "+00:00"
		}
		if err := db.Exec("SET time_zone = ?", timezoneValue).Error; err != nil {
			return nil, fmt.Errorf("failed to set MySQL timezone: %v", err)
		}
	}
	return db, nil
}

// initPostgreSQL initializes PostgreSQL database connection
func initPostgreSQL(cfg *DatabaseConnectionConfig) (*gorm.DB, error) {
	// Build PostgreSQL DSN
	dsn := buildPostgreSQLDSN(cfg)

	// Open PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %v", err)
	}

	// Set session timezone if configured
	if cfg.Timezone != "" {
		if err := db.Exec("SET timezone = ?", cfg.Timezone).Error; err != nil {
			return nil, fmt.Errorf("failed to set PostgreSQL timezone: %v", err)
		}
	}
	return db, nil
}

// buildMySQLDSN builds MySQL data source name
func buildMySQLDSN(cfg *DatabaseConnectionConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// Add connection parameters
	params := []string{}

	if cfg.Charset != "" {
		params = append(params, fmt.Sprintf("charset=%s", cfg.Charset))
	}

	if cfg.Timezone != "" {
		params = append(params, fmt.Sprintf("loc=%s", cfg.Timezone))
	}

	params = append(params, "parseTime=True")

	if cfg.ConnectTimeout > 0 {
		params = append(params, fmt.Sprintf("timeout=%ds", cfg.ConnectTimeout))
	}

	// Add SSL mode for MySQL
	if cfg.SSLMode != "" && cfg.SSLMode != "disable" {
		if cfg.SSLMode == "require" {
			params = append(params, "tls=true")
		} else {
			params = append(params, fmt.Sprintf("tls=%s", cfg.SSLMode))
		}
	}

	if len(params) > 0 {
		dsn += "?" + strings.Join(params, "&")
	}

	return dsn
}

// buildPostgreSQLDSN builds PostgreSQL data source name
func buildPostgreSQLDSN(cfg *DatabaseConnectionConfig) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Database,
	)

	if cfg.Password != "" {
		dsn += fmt.Sprintf(" password=%s", cfg.Password)
	}

	if cfg.SSLMode != "" {
		dsn += fmt.Sprintf(" sslmode=%s", cfg.SSLMode)
	}

	if cfg.ConnectTimeout > 0 {
		dsn += fmt.Sprintf(" connect_timeout=%d", cfg.ConnectTimeout)
	}

	// Add timezone support
	if cfg.Timezone != "" {
		dsn += fmt.Sprintf(" TimeZone=%s", cfg.Timezone)
	}

	return dsn
}

// configureConnectionPool configures database connection pool
func configureConnectionPool(db *gorm.DB, cfg *DatabaseConnectionConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Set maximum number of open connections
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	// Set maximum number of idle connections
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	// Set maximum connection lifetime
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	return nil
}

// TestDatabaseConnection tests database connection
func TestDatabaseConnection(cfg *config.Config) error {
	db, err := InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get underlying sql.DB for ping test
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Ping database to test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	return nil
}

// GetDatabaseInfo returns database connection information
func GetDatabaseInfo(cfg *config.Config) (map[string]any, error) {
	db, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	stats := sqlDB.Stats()

	info := map[string]any{
		"type":           cfg.Database.Type,
		"host":           cfg.Database.Host,
		"port":           cfg.Database.Port,
		"database":       cfg.Database.Database,
		"max_open_conns": stats.MaxOpenConnections,
		"open_conns":     stats.OpenConnections,
		"in_use":         stats.InUse,
		"idle":           stats.Idle,
	}

	// Add database-specific information
	switch strings.ToLower(cfg.Database.Type) {
	case DatabaseTypeSQLite:
		info["path"] = cfg.Database.Path
	case DatabaseTypeMySQL:
		info["charset"] = cfg.Database.Charset
		info["timezone"] = cfg.Database.Timezone
	case DatabaseTypePostgres, DatabaseTypePostgreSQL:
		info["ssl_mode"] = cfg.Database.SSLMode
	}

	return info, nil
}
