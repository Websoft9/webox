package plugins

import (
	"api-service/internal/constants"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"api-service/pkg/utils"
	"context"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// PermissionI18nHook is a GORM hook for automatic i18n translation of Permission.Name field
// It translates the Name field in Permission query results based on user's language preference
type PermissionI18nHook struct {
	db     *gorm.DB      // Database connection for language preference queries
	logger logger.Logger // Logger instance for debugging and monitoring
}

// NewPermissionI18nHook creates a new Permission i18n hook instance
// It initializes the hook with database connection and default logger
func NewPermissionI18nHook(db *gorm.DB) *PermissionI18nHook {
	return &PermissionI18nHook{
		db:     db,
		logger: logger.GetDefault(),
	}
}

// RegisterPermissionI18nHook registers the Permission i18n hook with GORM
// This should be called during database initialization to enable automatic translation
func RegisterPermissionI18nHook(db *gorm.DB) error {
	hook := NewPermissionI18nHook(db)

	// Register callback for after query operations
	// This handles SELECT queries and translates Name fields in the results
	err := db.Callback().Query().After("gorm:after_query").Register("permission:i18n_translate", hook.afterQueryCallback)
	if err != nil {
		hook.logger.Error("Failed to register after query callback", logger.ErrorField(err))
		return err
	}

	// Register callback for after find operations
	// This handles Find, First, Take and similar operations
	err = db.Callback().Query().After("gorm:after_find").Register("permission:i18n_translate_find", hook.afterQueryCallback)
	if err != nil {
		hook.logger.Error("Failed to register after find callback", logger.ErrorField(err))
		return err
	}

	hook.logger.Info("Permission i18n hook initialized successfully")
	return nil
}

// afterQueryCallback is the main callback function executed after query operations
// It automatically translates Permission.Name field based on user's language preference
func (h *PermissionI18nHook) afterQueryCallback(db *gorm.DB) {
	// Check if this is a Permission model query
	// Only process queries for the permissions table
	if !h.isPermissionQuery(db) {
		return
	}

	// Check if translation should be skipped for this query
	// This includes write operations, count queries, and explicitly skipped queries
	if h.shouldSkipTranslation(db) {
		return
	}

	// Retrieve the user's preferred language from context and database
	// Falls back to system language if user language is not available
	language := h.getUserLanguage(db.Statement.Context)

	// Skip translation if language is empty or default (en-US)
	if language == "" || language == constants.DefaultLanguage {
		return
	}

	// Translate Permission.Name fields in the query results
	h.translatePermissionNames(db, language)
}

// isPermissionQuery checks if the current query is for Permission model
// It examines the GORM schema to determine if we're querying the permissions table
func (h *PermissionI18nHook) isPermissionQuery(db *gorm.DB) bool {
	if db.Statement.Schema == nil {
		return false
	}

	tableName := db.Statement.Schema.Table
	isPermissionTable := tableName == "permissions"

	return isPermissionTable
}

// shouldSkipTranslation determines whether i18n translation should be skipped
// It checks for explicit skip flags, write operations, and aggregate queries
func (h *PermissionI18nHook) shouldSkipTranslation(db *gorm.DB) bool {
	// Check for explicit skip flag set by SkipPermissionI18n helper
	if skip, exists := db.Get("permission:skip_i18n"); exists && skip.(bool) {
		h.logger.Debug("Skipping: permission:skip_i18n flag is set")
		return true
	}

	// Skip if there's no SQL statement (shouldn't happen in normal cases)
	if db.Statement.SQL.String() == "" {
		h.logger.Debug("Skipping: empty SQL statement")
		return true
	}

	// Analyze the SQL statement to determine if translation is needed
	sql := strings.ToUpper(strings.TrimSpace(db.Statement.SQL.String()))
	sqlToLog := sql
	if len(sql) > constants.MaxSQLLogLength {
		sqlToLog = sql[:constants.MaxSQLLogLength] + "..."
	}
	h.logger.Debug("Checking SQL statement", logger.String("sql", sqlToLog))

	// Skip write operations (INSERT, UPDATE, DELETE)
	// These operations don't return data that needs translation
	if strings.HasPrefix(sql, "INSERT") || strings.HasPrefix(sql, "UPDATE") || strings.HasPrefix(sql, "DELETE") {
		h.logger.Debug("Skipping: write operation detected")
		return true
	}

	// Skip COUNT queries as they don't return Name field data
	if strings.Contains(sql, "SELECT COUNT(") || strings.Contains(sql, "SELECT COUNT *") {
		h.logger.Debug("Skipping: COUNT query detected")
		return true
	}

	// Skip other aggregate function queries that don't return individual records
	if strings.Contains(sql, "SELECT SUM(") || strings.Contains(sql, "SELECT AVG(") ||
		strings.Contains(sql, "SELECT MAX(") || strings.Contains(sql, "SELECT MIN(") {
		h.logger.Debug("Skipping: aggregate function query detected")
		return true
	}

	h.logger.Debug("Not skipping translation - proceeding with Permission.Name translation")
	return false
}

// getUserLanguage retrieves the user's preferred language from Redis cache first, then database
// It implements a fallback strategy: Redis cache -> user preference -> system setting -> default
func (h *PermissionI18nHook) getUserLanguage(ctx context.Context) string {
	if ctx == nil || h.db == nil {
		return constants.DefaultLanguage
	}

	// Extract user ID from the request context
	userID, exists := utils.GetUserIDFromContext(ctx)
	if !exists {
		return h.getSystemLanguage(ctx)
	}

	// Priority 1: Try to get language from Redis cache
	redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID)
	language, err := redis.HGet(ctx, redisKey, constants.UserLanguage)

	if err == nil && language != "" {
		h.logger.DebugContext(ctx, "Found user language in Redis cache",
			logger.Uint("userID", userID),
			logger.String(constants.UserLanguage, language),
			logger.String("redisKey", redisKey))
		return language
	}

	if err != nil && err.Error() != redis.RedisNilError {
		h.logger.DebugContext(ctx, "Failed to get language from Redis cache",
			logger.Uint("userID", userID),
			logger.String("redisKey", redisKey),
			logger.ErrorField(err))
	}

	// Priority 2: Query user language preference from user_profile table
	// This allows each user to have their own language setting
	var userProfile struct {
		ConfigValue string `gorm:"column:config_value"`
	}

	err = h.db.WithContext(ctx).
		Table("user_profile").
		Select("config_value").
		Where("user_id = ? AND category = ? AND config_key = ?", userID, constants.UserCategory, constants.UserLanguage).
		First(&userProfile).Error

	if err == nil && userProfile.ConfigValue != "" {
		h.logger.DebugContext(ctx, "Found user language preference in database",
			logger.Uint("userID", userID),
			logger.String(constants.UserLanguage, userProfile.ConfigValue))
		return userProfile.ConfigValue
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		h.logger.DebugContext(ctx, "Failed to get user language preference from database",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
	}

	// Priority 3: Fall back to system language setting
	return h.getSystemLanguage(ctx)
}

// getSystemLanguage retrieves the system-wide language setting from configuration
// This serves as a fallback when user-specific language is not available
func (h *PermissionI18nHook) getSystemLanguage(ctx context.Context) string {
	if h.db == nil {
		h.logger.DebugContext(ctx, "Database connection is nil, using default language")
		return constants.DefaultLanguage
	}

	// Query system language configuration from system_configs table
	var systemConfig struct {
		ConfigValue string `gorm:"column:config_value"`
	}

	err := h.db.WithContext(ctx).
		Table("system_configs").
		Select("config_value").
		Where("config_key = ?", "system.language").
		First(&systemConfig).Error

	if err == nil && systemConfig.ConfigValue != "" {
		return systemConfig.ConfigValue
	}

	// Final fallback to default language
	return constants.DefaultLanguage
}

// translatePermissionNames processes the query result and translates Permission.Name fields
// It handles both single struct results and slice results (collections)
func (h *PermissionI18nHook) translatePermissionNames(db *gorm.DB, language string) {
	if db.Statement.Dest == nil {
		h.logger.Debug("No destination to translate - Statement.Dest is nil")
		return
	}

	// Get the reflection value of the destination
	destValue := reflect.ValueOf(db.Statement.Dest)
	h.logger.Debug("Processing result for Permission.Name translation",
		logger.String("destType", destValue.Type().String()),
		logger.String("destKind", destValue.Kind().String()))

	// Dereference pointer if necessary
	if destValue.Kind() == reflect.Ptr {
		destValue = destValue.Elem()
		h.logger.Debug("Dereferenced pointer",
			logger.String("actualType", destValue.Type().String()),
			logger.String("actualKind", destValue.Kind().String()))
	}

	// Handle different result types
	switch destValue.Kind() {
	case reflect.Slice:
		// Handle slice results (typically from Find operations)
		// Iterate through each item in the slice and translate its Name field
		h.logger.Debug("Processing slice result", logger.Int("length", destValue.Len()))
		for i := 0; i < destValue.Len(); i++ {
			item := destValue.Index(i)
			h.logger.Debug("Translating slice item", logger.Int("index", i))
			h.translatePermissionName(item, language)
		}
	case reflect.Struct:
		// Handle single struct results (typically from First, Take, etc.)
		h.logger.Debug("Processing single struct result")
		h.translatePermissionName(destValue, language)
	default:
		h.logger.Debug("Unsupported result type for Permission.Name translation",
			logger.String("kind", destValue.Kind().String()))
	}
}

// translatePermissionName processes a single Permission struct and translates its Name field
// It handles embedded structs recursively and respects field visibility rules
func (h *PermissionI18nHook) translatePermissionName(structValue reflect.Value, language string) {
	h.logger.Debug("Starting Permission struct Name field translation",
		logger.String("structType", structValue.Type().String()),
		logger.String("structKind", structValue.Kind().String()),
		logger.Bool("canSet", structValue.CanSet()))

	// Safely handle pointer types by dereferencing them
	for structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			h.logger.Debug("Encountered nil pointer, skipping translation")
			return
		}
		structValue = structValue.Elem()
		h.logger.Debug("Dereferenced struct pointer",
			logger.String("actualType", structValue.Type().String()),
			logger.String("actualKind", structValue.Kind().String()))
	}

	// Ensure we're working with a struct
	if structValue.Kind() != reflect.Struct {
		h.logger.Debug("Value is not a struct after dereferencing",
			logger.String("kind", structValue.Kind().String()))
		return
	}

	// Try to find the Name field in the Permission struct
	nameField := structValue.FieldByName("Name")
	if !nameField.IsValid() {
		h.logger.Debug("Name field not found in struct")
		return
	}

	// Check if it's a string field and can be modified
	if nameField.Kind() != reflect.String || !nameField.CanSet() {
		h.logger.Debug("Name field is not a settable string field",
			logger.String("fieldKind", nameField.Kind().String()),
			logger.Bool("canSet", nameField.CanSet()))
		return
	}

	// Get the original Name field value
	originalName := nameField.String()
	if originalName == "" {
		h.logger.Debug("Name field is empty, skipping translation")
		return
	}

	h.logger.Debug("Found Permission.Name field for translation",
		logger.String("originalName", originalName),
		logger.String("targetLanguage", language))

	// Translate the value with safety check to prevent panics
	var translatedName string
	func() {
		defer func() {
			if r := recover(); r != nil {
				// If translation panics, use original name as fallback
				h.logger.Warn("Translation panicked, using original name",
					logger.String("originalName", originalName),
					logger.Any("panicValue", r))
				translatedName = originalName
			}
		}()
		translatedName = i18n.T(originalName, language)
	}()

	// Set the translated value back to the Name field
	nameField.SetString(translatedName)

	// Log the successful translation
	h.logger.Debug("Successfully translated Permission.Name field",
		logger.String("original", originalName),
		logger.String("translated", translatedName),
		logger.String(constants.UserLanguage, language))
}

// SkipPermissionI18n is a helper function to skip i18n translation for specific Permission queries
// Usage: SkipPermissionI18n(db).Where(...).Find(&permissions)
func SkipPermissionI18n(db *gorm.DB) *gorm.DB {
	return db.Set("permission:skip_i18n", true)
}
