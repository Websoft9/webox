package database

import (
	"api-service/internal/config"
	"api-service/pkg/database/plugins"
	"api-service/pkg/logger"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// DatabaseManager defines the database manager interface
type DatabaseManager interface {
	// Basic operations
	Query(ctx context.Context, fn func(*gorm.DB) error) error
	Transaction(ctx context.Context, fn func(*gorm.DB) error) error

	// CRUD operations
	Create(ctx context.Context, value any) error
	Update(ctx context.Context, model, updates any) error
	Delete(ctx context.Context, value any, conds ...any) error
	BatchCreate(ctx context.Context, values any, batchSize int) error

	// Management operations
	GetDB() *gorm.DB
	HealthCheck(ctx context.Context) error
	Close() error
}

// DatabaseManagerFactory creates database managers
type DatabaseManagerFactory struct {
	config *config.Config
}

// NewDatabaseManagerFactory creates a database manager factory
func NewDatabaseManagerFactory(cfg *config.Config) *DatabaseManagerFactory {
	return &DatabaseManagerFactory{
		config: cfg,
	}
}

// CreateManager creates the appropriate database manager based on configuration
func (factory *DatabaseManagerFactory) CreateManager(db *gorm.DB) (DatabaseManager, error) {
	logger.Debug("Creating database manager", logger.String("database_type", factory.config.Database.Type))

	// Initialize timezone conversion plugin
	if err := factory.initTimezonePlugin(db); err != nil {
		logger.Error("Failed to initialize timezone plugin", logger.ErrorField(err))
		// Does not affect manager creation, only logs the error
	}

	// Initialize Permission i18n hook
	if err := factory.initPermissionI18nHook(db); err != nil {
		logger.Error("Failed to initialize Permission i18n hook", logger.ErrorField(err))
		// Does not affect manager creation, only logs the error
	}

	switch factory.config.Database.Type {
	case DatabaseTypeSQLite:
		logger.Debug("Creating SQLite manager")
		return NewSQLiteManager(db), nil
	case DatabaseTypeMySQL:
		logger.Debug("Creating MySQL generic manager")
		return NewGenericManager(db), nil
	case DatabaseTypePostgres, DatabaseTypePostgreSQL:
		logger.Debug("Creating PostgreSQL generic manager")
		return NewGenericManager(db), nil
	default:
		logger.Error("Unsupported database type", logger.String("database_type", factory.config.Database.Type))
		return nil, fmt.Errorf("unsupported database type: %s", factory.config.Database.Type)
	}
}

// initTimezonePlugin initializes the timezone conversion plugin
func (factory *DatabaseManagerFactory) initTimezonePlugin(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Create timezone conversion plugin with Redis client
	timezonePlugin := plugins.NewTimezonePlugin(db, plugins.WithLogger(logger.GetDefault()))

	// Register plugin to database
	if err := db.Use(timezonePlugin); err != nil {
		return fmt.Errorf("failed to register timezone plugin: %v", err)
	}

	logger.Info("Timezone conversion plugin initialized successfully")
	return nil
}

// initPermissionI18nHook initializes the Permission i18n hook
func (factory *DatabaseManagerFactory) initPermissionI18nHook(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Register Permission i18n hook
	if err := plugins.RegisterPermissionI18nHook(db); err != nil {
		return fmt.Errorf("failed to register Permission i18n hook: %v", err)
	}

	logger.Info("Permission i18n hook initialized successfully")
	return nil
}

// GenericManager is a generic database manager (for MySQL and PostgreSQL)
type GenericManager struct {
	db *gorm.DB
}

// NewGenericManager creates a generic database manager
func NewGenericManager(db *gorm.DB) *GenericManager {
	logger.Debug("Initializing generic database manager")
	return &GenericManager{
		db: db,
	}
}

// Query executes query operations
func (gm *GenericManager) Query(ctx context.Context, fn func(*gorm.DB) error) error {
	logger.Debug("Executing query operation")
	err := fn(gm.db.WithContext(ctx))
	if err != nil {
		logger.Debug("Query operation failed", logger.ErrorField(err))
	} else {
		logger.Debug("Query operation completed successfully")
	}
	return err
}

// Transaction executes transaction operations
func (gm *GenericManager) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	logger.Debug("Starting database transaction")
	err := gm.db.Transaction(func(tx *gorm.DB) error {
		if ctx != context.Background() {
			tx = tx.WithContext(ctx)
		}
		logger.Debug("Executing transaction function")
		return fn(tx)
	})
	if err != nil {
		logger.Debug("Transaction failed", logger.ErrorField(err))
	} else {
		logger.Debug("Transaction completed successfully")
	}
	return err
}

// Create executes create operations
func (gm *GenericManager) Create(ctx context.Context, value any) error {
	return gm.db.WithContext(ctx).Create(value).Error
}

// Update executes update operations
func (gm *GenericManager) Update(ctx context.Context, model, updates any) error {
	return gm.db.WithContext(ctx).Model(model).Updates(updates).Error
}

// Delete executes delete operations
func (gm *GenericManager) Delete(ctx context.Context, value any, conds ...any) error {
	return gm.db.WithContext(ctx).Delete(value, conds...).Error
}

// BatchCreate executes batch create operations
func (gm *GenericManager) BatchCreate(ctx context.Context, values any, batchSize int) error {
	return gm.db.WithContext(ctx).CreateInBatches(values, batchSize).Error
}

// GetDB gets the raw database connection
func (gm *GenericManager) GetDB() *gorm.DB {
	return gm.db
}

// HealthCheck performs health check
func (gm *GenericManager) HealthCheck(ctx context.Context) error {
	logger.Debug("Starting database health check")
	sqlDB, err := gm.db.DB()
	if err != nil {
		logger.Debug("Failed to get raw database connection", logger.ErrorField(err))
		return err
	}
	err = sqlDB.PingContext(ctx)
	if err != nil {
		logger.Debug("Database ping failed", logger.ErrorField(err))
	} else {
		logger.Debug("Database ping successful")
	}
	return err
}

// Close closes the database connection
func (gm *GenericManager) Close() error {
	logger.Debug("Closing database connection")
	sqlDB, err := gm.db.DB()
	if err != nil {
		logger.Error("Failed to get raw database connection for closing", logger.ErrorField(err))
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		logger.Error("Failed to close database connection", logger.ErrorField(err))
	} else {
		logger.Debug("Database connection closed successfully")
	}
	return err
}
