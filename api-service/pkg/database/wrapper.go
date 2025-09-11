package database

import (
	"api-service/internal/config"
	"api-service/pkg/logger"
	"context"
	"sync"

	"gorm.io/gorm"
)

// DBWrapper wraps database connection
type DBWrapper struct {
	manager DatabaseManager
	db      *gorm.DB
	config  *config.Config
	mu      sync.RWMutex
}

var (
	dbInstance *DBWrapper
	dbOnce     sync.Once
)

// InitDBWrapper initializes database wrapper (singleton pattern)
func InitDBWrapper(cfg *config.Config) (*DBWrapper, error) {
	logger.Debug("Initializing database wrapper", logger.String("database_type", cfg.Database.Type))
	var err error
	dbOnce.Do(func() {
		dbInstance, err = newDBWrapper(cfg)
	})

	if err != nil {
		// Reset to allow re-initialization
		logger.Error("Failed to initialize database wrapper", logger.ErrorField(err))
		dbOnce = sync.Once{}
		dbInstance = nil
	} else {
		logger.Debug("Database wrapper initialized successfully")
	}

	return dbInstance, err
}

// GetDBWrapper gets the database wrapper instance
func GetDBWrapper() *DBWrapper {
	return dbInstance
}

// newDBWrapper creates a new database wrapper
func newDBWrapper(cfg *config.Config) (*DBWrapper, error) {
	logger.Debug("Creating new database wrapper")
	// Initialize raw database connection
	db, err := InitDB(cfg)
	if err != nil {
		logger.Error("Failed to initialize raw database connection", logger.ErrorField(err))
		return nil, err
	}

	// Create database manager factory
	logger.Debug("Creating database manager factory")
	factory := NewDatabaseManagerFactory(cfg)

	// Create the appropriate manager based on database type
	logger.Debug("Creating database manager", logger.String("database_type", cfg.Database.Type))
	manager, err := factory.CreateManager(db)
	if err != nil {
		logger.Error("Failed to create database manager", logger.ErrorField(err))
		return nil, err
	}

	logger.Debug("Database wrapper created successfully")
	return &DBWrapper{
		manager: manager,
		db:      db,
		config:  cfg,
	}, nil
}

// GetManager gets the database manager
func (w *DBWrapper) GetManager() DatabaseManager {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.manager
}

// GetDB gets the raw GORM database connection (compatibility method)
func (w *DBWrapper) GetDB() *gorm.DB {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.db
}

// Query executes query operations
func (w *DBWrapper) Query(ctx context.Context, fn func(*gorm.DB) error) error {
	return w.manager.Query(ctx, fn)
}

// Transaction executes transaction operations
func (w *DBWrapper) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return w.manager.Transaction(ctx, fn)
}

// Create creates records
func (w *DBWrapper) Create(ctx context.Context, value any) error {
	return w.manager.Create(ctx, value)
}

// Update updates records
func (w *DBWrapper) Update(ctx context.Context, model, updates any) error {
	return w.manager.Update(ctx, model, updates)
}

// Delete deletes records
func (w *DBWrapper) Delete(ctx context.Context, value any, conds ...any) error {
	return w.manager.Delete(ctx, value, conds...)
}

// BatchCreate creates records in batches
func (w *DBWrapper) BatchCreate(ctx context.Context, values any, batchSize int) error {
	return w.manager.BatchCreate(ctx, values, batchSize)
}

// HealthCheck performs health check
func (w *DBWrapper) HealthCheck(ctx context.Context) error {
	return w.manager.HealthCheck(ctx)
}

// Close closes the database connection
func (w *DBWrapper) Close() error {
	logger.Debug("Closing database wrapper")
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.manager != nil {
		err := w.manager.Close()
		if err != nil {
			logger.Error("Failed to close database manager", logger.ErrorField(err))
		} else {
			logger.Debug("Database wrapper closed successfully")
		}
		return err
	}
	return nil
}

// GetConfig gets the database configuration
func (w *DBWrapper) GetConfig() *config.Config {
	return w.config
}

// IsSQLite checks if the database is SQLite
func (w *DBWrapper) IsSQLite() bool {
	return w.config.Database.Type == "sqlite"
}

// GetActiveTransactions gets the number of active transactions (SQLite only)
func (w *DBWrapper) GetActiveTransactions() int {
	if sqliteManager, ok := w.manager.(*SQLiteManager); ok {
		return sqliteManager.GetActiveTransactions()
	}
	return 0
}

// ResetConnection resets the database connection (for testing or connection recovery)
func (w *DBWrapper) ResetConnection() error {
	logger.Debug("Resetting database connection")
	w.mu.Lock()
	defer w.mu.Unlock()

	// Close existing connection
	if w.manager != nil {
		logger.Debug("Closing existing connection")
		if err := w.manager.Close(); err != nil {
			logger.Warn("Failed to close database manager", logger.ErrorField(err))
		}
	}

	// Re-initialize connection
	logger.Debug("Re-initializing database connection")
	db, err := InitDB(w.config)
	if err != nil {
		logger.Error("Failed to re-initialize database connection", logger.ErrorField(err))
		return err
	}

	// Create new manager
	logger.Debug("Creating new database manager")
	factory := NewDatabaseManagerFactory(w.config)
	manager, err := factory.CreateManager(db)
	if err != nil {
		logger.Error("Failed to create new database manager", logger.ErrorField(err))
		return err
	}

	w.db = db
	w.manager = manager
	logger.Debug("Database connection reset successfully")
	return nil
}
