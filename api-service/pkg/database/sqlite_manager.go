package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"api-service/pkg/database/errors"
	"api-service/pkg/logger"

	"gorm.io/gorm"
)

const (
	// SQLite lock-related constants
	DefaultRetryAttempts   = 5
	DefaultRetryDelay      = 100 * time.Millisecond
	DefaultBusyTimeout     = 30 * time.Second
	DefaultTransactionMode = "IMMEDIATE"
	DefaultMaxRetryDelay   = 5 * time.Second
)

// SQLiteManager SQLite database connection manager
type SQLiteManager struct {
	db                 *gorm.DB
	transactionMutex   sync.RWMutex
	activeTransactions int
	retryHandler       *errors.SQLiteRetryHandler
	busyTimeout        time.Duration
}

// NewSQLiteManager creates a new SQLite manager
func NewSQLiteManager(db *gorm.DB) *SQLiteManager {
	manager := &SQLiteManager{
		db:           db,
		retryHandler: errors.NewSQLiteRetryHandler(errors.DefaultRetryConfig()),
		busyTimeout:  DefaultBusyTimeout,
	}

	// Initialize SQLite-specific optimization settings
	logger.Debug("Initializing SQLite manager", logger.String("manager_type", "sqlite"))
	manager.initializePragmas()
	logger.Debug("SQLite manager initialized successfully")

	return manager
}

// initializePragmas initializes SQLite Pragma settings
func (sm *SQLiteManager) initializePragmas() {
	sqlDB, err := sm.db.DB()
	if err != nil {
		return
	}

	logger.Debug("Initializing SQLite pragmas")
	pragmas := []string{
		"PRAGMA journal_mode = WAL",        // Write-Ahead Logging mode
		"PRAGMA synchronous = NORMAL",      // Normal sync mode, balancing performance and safety
		"PRAGMA cache_size = -64000",       // 64MB cache
		"PRAGMA foreign_keys = ON",         // Enable foreign key constraints
		"PRAGMA temp_store = MEMORY",       // Use memory for temporary storage
		"PRAGMA mmap_size = 268435456",     // 256MB memory mapping
		"PRAGMA wal_autocheckpoint = 1000", // Automatic checkpoint
		fmt.Sprintf("PRAGMA busy_timeout = %d", int(sm.busyTimeout.Milliseconds())),
	}

	for _, pragma := range pragmas {
		logger.Debug("Executing SQLite pragma", logger.String("pragma", pragma))
		if _, err := sqlDB.Exec(pragma); err != nil {
			// Log but don't interrupt, some pragmas might not be supported
			logger.Warn("Failed to execute pragma", logger.String("pragma", pragma), logger.ErrorField(err))
		} else {
			logger.Debug("Pragma executed successfully", logger.String("pragma", pragma))
		}
	}
}

// WithRetry executes operations with retry mechanism
func (sm *SQLiteManager) WithRetry(ctx context.Context, operation func() error, operationName string) error {
	logger.Debug("Starting retry operation", logger.String("operation", operationName), logger.Int("max_attempts", DefaultRetryAttempts))
	var lastErr error

	for attempt := 1; attempt <= DefaultRetryAttempts; attempt++ {
		logger.Debug("Executing operation attempt", logger.String("operation", operationName), logger.Int("attempt", attempt))
		if attempt > 1 {
			// Use retry handler to calculate delay
			logger.Debug("Checking if retry is needed", logger.String("operation", operationName), logger.ErrorField(lastErr))
			shouldRetry, delay := sm.retryHandler.ShouldRetry(lastErr, operationName, attempt-1)
			if !shouldRetry {
				logger.Debug("Retry limit reached or error not retryable", logger.String("operation", operationName))
				return sm.retryHandler.WrapError(lastErr, operationName, attempt-1)
			}

			logger.Debug("Waiting before retry", logger.String("operation", operationName), logger.Duration("delay", delay))
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				logger.Debug("Operation canceled by context", logger.String("operation", operationName))
				return ctx.Err()
			}
		}

		err := operation()
		if err == nil {
			logger.Debug("Operation completed successfully", logger.String("operation", operationName), logger.Int("attempt", attempt))
			return nil
		}

		logger.Debug("Operation failed", logger.String("operation", operationName), logger.Int("attempt", attempt), logger.ErrorField(err))
		lastErr = err
	}

	// Wrap final error
	logger.Error("All retry attempts failed", logger.String("operation", operationName), logger.Int("total_attempts", DefaultRetryAttempts), logger.ErrorField(lastErr))
	return sm.retryHandler.WrapError(lastErr, operationName, DefaultRetryAttempts)
}

// Transaction executes transaction operations
func (sm *SQLiteManager) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return sm.WithRetry(ctx, func() error {
		return sm.executeTransaction(ctx, fn)
	}, "transaction")
}

// executeTransaction executes a single transaction
func (sm *SQLiteManager) executeTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	// Acquire transaction lock
	logger.Debug("Acquiring transaction lock")
	sm.transactionMutex.Lock()
	sm.activeTransactions++
	logger.Debug("Transaction started", logger.Int("active_transactions", sm.activeTransactions))
	sm.transactionMutex.Unlock()

	defer func() {
		sm.transactionMutex.Lock()
		sm.activeTransactions--
		logger.Debug("Transaction completed", logger.Int("active_transactions", sm.activeTransactions))
		sm.transactionMutex.Unlock()
	}()

	// Use IMMEDIATE transaction mode to acquire write lock immediately
	logger.Debug("Starting database transaction with IMMEDIATE mode")
	return sm.db.Transaction(func(tx *gorm.DB) error {
		// Set transaction timeout
		if ctx != context.Background() {
			tx = tx.WithContext(ctx)
		}

		// Execute user operation
		logger.Debug("Executing user transaction function")
		err := fn(tx)
		if err != nil {
			logger.Debug("Transaction function failed", logger.ErrorField(err))
		} else {
			logger.Debug("Transaction function completed successfully")
		}
		return err
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

// Query executes query operations
func (sm *SQLiteManager) Query(ctx context.Context, fn func(*gorm.DB) error) error {
	return sm.WithRetry(ctx, func() error {
		return fn(sm.db.WithContext(ctx))
	}, "query")
}

// Create executes create operations
func (sm *SQLiteManager) Create(ctx context.Context, value any) error {
	return sm.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Create(value).Error
	})
}

// Update executes update operations
func (sm *SQLiteManager) Update(ctx context.Context, model, updates any) error {
	return sm.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Model(model).Updates(updates).Error
	})
}

// Delete executes delete operations
func (sm *SQLiteManager) Delete(ctx context.Context, value any, conds ...any) error {
	return sm.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Delete(value, conds...).Error
	})
}

// BatchCreate executes batch create operations
func (sm *SQLiteManager) BatchCreate(ctx context.Context, values any, batchSize int) error {
	return sm.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.CreateInBatches(values, batchSize).Error
	})
}

// GetDB gets the raw database connection
func (sm *SQLiteManager) GetDB() *gorm.DB {
	return sm.db
}

// GetActiveTransactions gets the current number of active transactions
func (sm *SQLiteManager) GetActiveTransactions() int {
	sm.transactionMutex.RLock()
	defer sm.transactionMutex.RUnlock()
	return sm.activeTransactions
}

// SetRetryHandler sets the retry handler
func (sm *SQLiteManager) SetRetryHandler(handler *errors.SQLiteRetryHandler) {
	sm.retryHandler = handler
}

// HealthCheck performs health check
func (sm *SQLiteManager) HealthCheck(ctx context.Context) error {
	logger.Debug("Starting SQLite health check")
	return sm.WithRetry(ctx, func() error {
		sqlDB, err := sm.db.DB()
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
	}, "health_check")
}

// Close closes the database connection
func (sm *SQLiteManager) Close() error {
	logger.Debug("Closing SQLite database connection")
	sqlDB, err := sm.db.DB()
	if err != nil {
		logger.Error("Failed to get raw database connection for closing", logger.ErrorField(err))
		return err
	}

	// Execute WAL checkpoint to ensure all data is written to disk
	logger.Debug("Executing WAL checkpoint before closing")
	_, _ = sqlDB.Exec("PRAGMA wal_checkpoint(TRUNCATE)")

	err = sqlDB.Close()
	if err != nil {
		logger.Error("Failed to close database connection", logger.ErrorField(err))
	} else {
		logger.Debug("SQLite database connection closed successfully")
	}
	return err
}
