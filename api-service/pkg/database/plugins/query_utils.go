package plugins

import (
	"api-service/internal/constants"
	"api-service/pkg/logger"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// shouldSkipQuery determines whether a query should be skipped for processing
// It checks for explicit skip flags, write operations, and aggregate queries
func shouldSkipQuery(db *gorm.DB, skipKey string, log logger.Logger) bool {
	// Check for explicit skip flag
	if skip, exists := db.Get(skipKey); exists && skip.(bool) {
		log.Debug("Skipping: skip flag is set", logger.String("skipKey", skipKey))
		return true
	}

	// Skip if there's no SQL statement (shouldn't happen in normal cases)
	if db.Statement.SQL.String() == "" {
		log.Debug("Skipping: empty SQL statement")
		return true
	}

	// Analyze the SQL statement to determine if processing is needed
	sql := strings.ToUpper(strings.TrimSpace(db.Statement.SQL.String()))
	sqlToLog := sql
	if len(sql) > constants.MaxSQLLogLength {
		sqlToLog = sql[:constants.MaxSQLLogLength] + "..."
	}
	log.Debug("Checking SQL statement", logger.String("sql", sqlToLog))

	// Skip write operations (INSERT, UPDATE, DELETE)
	// These operations don't return data that needs processing
	if strings.HasPrefix(sql, "INSERT") || strings.HasPrefix(sql, "UPDATE") || strings.HasPrefix(sql, "DELETE") {
		log.Debug("Skipping: write operation detected")
		return true
	}

	// Skip COUNT queries as they don't return individual records
	if strings.Contains(sql, "SELECT COUNT(") || strings.Contains(sql, "SELECT COUNT *") {
		log.Debug("Skipping: COUNT query detected")
		return true
	}

	// Skip other aggregate function queries that don't return individual records
	if strings.Contains(sql, "SELECT SUM(") || strings.Contains(sql, "SELECT AVG(") ||
		strings.Contains(sql, "SELECT MAX(") || strings.Contains(sql, "SELECT MIN(") {
		log.Debug("Skipping: aggregate function query detected")
		return true
	}

	log.Debug("Not skipping query - proceeding with processing")
	return false
}

// processQueryResult processes query results by applying a function to each item
// It handles both single struct results and slice results (collections)
func processQueryResult(db *gorm.DB, log logger.Logger, processItem func(reflect.Value)) {
	if db.Statement.Dest == nil {
		log.Debug("No destination to process - Statement.Dest is nil")
		return
	}

	// Get the reflection value of the destination
	destValue := reflect.ValueOf(db.Statement.Dest)
	log.Debug("Processing query result",
		logger.String("destType", destValue.Type().String()),
		logger.String("destKind", destValue.Kind().String()))

	// Dereference pointer if necessary
	if destValue.Kind() == reflect.Ptr {
		destValue = destValue.Elem()
		log.Debug("Dereferenced pointer",
			logger.String("actualType", destValue.Type().String()),
			logger.String("actualKind", destValue.Kind().String()))
	}

	// Handle different result types
	switch destValue.Kind() {
	case reflect.Slice:
		// Handle slice results (typically from Find operations)
		// Iterate through each item in the slice and process it
		log.Debug("Processing slice result", logger.Int("length", destValue.Len()))
		for i := 0; i < destValue.Len(); i++ {
			item := destValue.Index(i)
			log.Debug("Processing slice item", logger.Int("index", i))
			processItem(item)
		}
	case reflect.Struct:
		// Handle single struct results (typically from First, Take, etc.)
		log.Debug("Processing single struct result")
		processItem(destValue)
	default:
		log.Debug("Unsupported result type for processing",
			logger.String("kind", destValue.Kind().String()))
	}
}
