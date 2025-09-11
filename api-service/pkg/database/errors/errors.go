package errors

import (
	"api-service/pkg/errors"
	"fmt"
	"strings"
	"time"
)

// Retry delay constants
const (
	// Delay increment factor for locked errors (milliseconds)
	LockedDelayFactor = 100
	// Delay increment factor for busy errors (milliseconds)
	BusyDelayFactor = 50
	// Delay increment factor for transaction errors (milliseconds)
	TransactionDelayFactor = 200
	// Default delay time (milliseconds)
	DefaultDelayMs = 100

	// Default retry configuration constants
	DefaultMaxAttempts     = 5
	DefaultBaseDelayMs     = 100
	DefaultMaxDelaySeconds = 5
	DefaultBackoffFactor   = 2.0
)

// SQLiteErrorType represents SQLite error types
type SQLiteErrorType int

const (
	SQLiteErrorUnknown SQLiteErrorType = iota
	SQLiteErrorLocked
	SQLiteErrorBusy
	SQLiteErrorCorrupt
	SQLiteErrorConstraint
	SQLiteErrorTransaction
)

// SQLiteError represents SQLite-specific errors
type SQLiteError struct {
	Type          SQLiteErrorType
	OriginalError error
	Operation     string
	Timestamp     time.Time
	AttemptCount  int
}

// Error implements the error interface
func (e *SQLiteError) Error() string {
	return fmt.Sprintf("SQLite %s error in operation '%s' (attempt %d): %v",
		e.getTypeString(), e.Operation, e.AttemptCount, e.OriginalError)
}

// Unwrap supports errors.Unwrap
func (e *SQLiteError) Unwrap() error {
	return e.OriginalError
}

// getTypeString returns the error type string
func (e *SQLiteError) getTypeString() string {
	switch e.Type {
	case SQLiteErrorLocked:
		return "LOCKED"
	case SQLiteErrorBusy:
		return "BUSY"
	case SQLiteErrorCorrupt:
		return "CORRUPT"
	case SQLiteErrorConstraint:
		return "CONSTRAINT"
	case SQLiteErrorTransaction:
		return "TRANSACTION"
	default:
		return "UNKNOWN"
	}
}

// IsRetryable determines if the error is retryable
func (e *SQLiteError) IsRetryable() bool {
	switch e.Type {
	case SQLiteErrorLocked, SQLiteErrorBusy, SQLiteErrorTransaction:
		return true
	default:
		return false
	}
}

// GetRetryDelay returns the retry delay duration
func (e *SQLiteError) GetRetryDelay() time.Duration {
	switch e.Type {
	case SQLiteErrorLocked:
		return time.Duration(e.AttemptCount*LockedDelayFactor) * time.Millisecond // 100ms, 200ms, 300ms...
	case SQLiteErrorBusy:
		return time.Duration(e.AttemptCount*BusyDelayFactor) * time.Millisecond // 50ms, 100ms, 150ms...
	case SQLiteErrorTransaction:
		return time.Duration(e.AttemptCount*TransactionDelayFactor) * time.Millisecond // 200ms, 400ms, 600ms...
	default:
		return DefaultDelayMs * time.Millisecond
	}
}

// SQLiteErrorClassifier classifies SQLite errors
type SQLiteErrorClassifier struct{}

// NewSQLiteErrorClassifier creates a new SQLite error classifier
func NewSQLiteErrorClassifier() *SQLiteErrorClassifier {
	return &SQLiteErrorClassifier{}
}

// ClassifyError classifies SQLite errors
func (c *SQLiteErrorClassifier) ClassifyError(err error, operation string, attemptCount int) *SQLiteError {
	if err == nil {
		return nil
	}

	errorMsg := strings.ToLower(err.Error())
	errorType := c.determineErrorType(errorMsg)

	return &SQLiteError{
		Type:          errorType,
		OriginalError: err,
		Operation:     operation,
		Timestamp:     time.Now(),
		AttemptCount:  attemptCount,
	}
}

// determineErrorType determines the error type
func (c *SQLiteErrorClassifier) determineErrorType(errorMsg string) SQLiteErrorType {
	// Database locked errors
	if containsAny(errorMsg, []string{
		"database is locked",
		"database table is locked",
		"sqlite_locked",
	}) {
		return SQLiteErrorLocked
	}

	// Database busy errors
	if containsAny(errorMsg, []string{
		"database is busy",
		"sqlite_busy",
		"sqlite_busy_recovery",
		"sqlite_busy_snapshot",
		"sqlite_busy_timeout",
	}) {
		return SQLiteErrorBusy
	}

	// Database corruption errors
	if containsAny(errorMsg, []string{
		"database disk image is malformed",
		"database corruption",
		"sqlite_corrupt",
		"sqlite_notadb",
	}) {
		return SQLiteErrorCorrupt
	}

	// Constraint errors
	if containsAny(errorMsg, []string{
		"constraint failed",
		"unique constraint",
		"foreign key constraint",
		"check constraint",
		"not null constraint",
		"sqlite_constraint",
	}) {
		return SQLiteErrorConstraint
	}

	// Transaction errors
	if containsAny(errorMsg, []string{
		"cannot start a transaction within a transaction",
		"transaction",
		"sqlite_misuse",
	}) {
		return SQLiteErrorTransaction
	}

	return SQLiteErrorUnknown
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []SQLiteErrorType
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   DefaultMaxAttempts,
		BaseDelay:     DefaultBaseDelayMs * time.Millisecond,
		MaxDelay:      DefaultMaxDelaySeconds * time.Second,
		BackoffFactor: DefaultBackoffFactor,
		RetryableErrors: []SQLiteErrorType{
			SQLiteErrorLocked,
			SQLiteErrorBusy,
			SQLiteErrorTransaction,
		},
	}
}

// SQLiteRetryHandler handles SQLite retry logic
type SQLiteRetryHandler struct {
	classifier *SQLiteErrorClassifier
	config     *RetryConfig
}

// NewSQLiteRetryHandler creates a new SQLite retry handler
func NewSQLiteRetryHandler(config *RetryConfig) *SQLiteRetryHandler {
	if config == nil {
		config = DefaultRetryConfig()
	}

	return &SQLiteRetryHandler{
		classifier: NewSQLiteErrorClassifier(),
		config:     config,
	}
}

// ShouldRetry determines if the operation should be retried
func (h *SQLiteRetryHandler) ShouldRetry(err error, operation string, attemptCount int) (bool, time.Duration) {
	if err == nil || attemptCount >= h.config.MaxAttempts {
		return false, 0
	}

	sqliteErr := h.classifier.ClassifyError(err, operation, attemptCount)
	if sqliteErr == nil {
		return false, 0
	}

	// Check if error type is retryable
	retryable := false
	for _, retryableType := range h.config.RetryableErrors {
		if sqliteErr.Type == retryableType {
			retryable = true
			break
		}
	}

	if !retryable {
		return false, 0
	}

	// Calculate retry delay
	delay := h.calculateDelay(attemptCount)
	return true, delay
}

// calculateDelay calculates the retry delay
func (h *SQLiteRetryHandler) calculateDelay(attemptCount int) time.Duration {
	delay := h.config.BaseDelay

	// Exponential backoff
	for i := 1; i < attemptCount; i++ {
		delay = time.Duration(float64(delay) * h.config.BackoffFactor)
	}

	// Limit maximum delay
	if delay > h.config.MaxDelay {
		delay = h.config.MaxDelay
	}

	return delay
}

// WrapError wraps error information
func (h *SQLiteRetryHandler) WrapError(err error, operation string, attemptCount int) error {
	if err == nil {
		return nil
	}

	sqliteErr := h.classifier.ClassifyError(err, operation, attemptCount)
	if sqliteErr != nil {
		return sqliteErr
	}

	// If not a SQLite-specific error, return original error
	return err
}

// containsAny checks if string contains any of the specified substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// Predefined common errors
var (
	ErrDatabaseLocked     = errors.NewAppError(errors.CodeInternalError, "database is locked")
	ErrDatabaseBusy       = errors.NewAppError(errors.CodeInternalError, "database is busy")
	ErrDatabaseCorrupt    = errors.NewAppError(errors.CodeInternalError, "database is corrupt")
	ErrTransactionActive  = errors.NewAppError(errors.CodeInternalError, "transaction already active")
	ErrMaxRetriesExceeded = errors.NewAppError(errors.CodeInternalError, "maximum retry attempts exceeded")
)
