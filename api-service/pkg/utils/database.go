package utils

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database type constants
const (
	DatabaseTypeSQLite     = "sqlite"
	DatabaseTypeMySQL      = "mysql"
	DatabaseTypePostgres   = "postgres"
	DatabaseTypePostgreSQL = "postgresql"
)

// DatabaseConfig contains database connection parameters
type DatabaseConnectionConfig struct {
	Type            string
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
}

// InitDB initializes database connection based on configuration
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Convert config to connection config
	dbConfig := &DatabaseConnectionConfig{
		Type:            cfg.Database.Type,
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
	}

	// Initialize database connection based on type
	switch strings.ToLower(cfg.Database.Type) {
	case DatabaseTypeSQLite:
		db, err = initSQLite(cfg.Database.Path)
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

	// Configure connection pool for MySQL and PostgreSQL
	if cfg.Database.Type != DatabaseTypeSQLite {
		if err := configureConnectionPool(db, dbConfig); err != nil {
			return nil, fmt.Errorf("failed to configure connection pool: %v", err)
		}
	}

	return db, nil
}

// initSQLite initializes SQLite database connection
func initSQLite(dbPath string) (*gorm.DB, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, constants.DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %v", err)
	}

	// Enable foreign key constraints for SQLite
	db.Exec("PRAGMA foreign_keys = ON")

	return db, nil
}

// initMySQL initializes MySQL database connection
func initMySQL(cfg *DatabaseConnectionConfig) (*gorm.DB, error) {
	// Build MySQL DSN
	dsn := buildMySQLDSN(cfg)

	// Open MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %v", err)
	}

	return db, nil
}

// initPostgreSQL initializes PostgreSQL database connection
func initPostgreSQL(cfg *DatabaseConnectionConfig) (*gorm.DB, error) {
	// Build PostgreSQL DSN
	dsn := buildPostgreSQLDSN(cfg)

	// Open PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %v", err)
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
func GetDatabaseInfo(cfg *config.Config) (map[string]interface{}, error) {
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

	info := map[string]interface{}{
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

// InitInfluxDB initializes InfluxDB connection (unchanged)
func InitInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	client := influxdb2.NewClient(cfg.InfluxDB.URL, cfg.InfluxDB.Token)
	return client, nil
}
