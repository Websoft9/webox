package plugins

import (
	"api-service/internal/constants"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"api-service/pkg/utils"
	"context"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// TimezonePlugin is a GORM plugin for automatic timezone conversion
// It converts time fields in query results from UTC to user's preferred timezone
type TimezonePlugin struct {
	name              string          // Plugin name identifier
	db                *gorm.DB        // Database connection for timezone queries
	convertibleFields map[string]bool // Map of field names that should be converted
	logger            logger.Logger   // Logger instance for debugging and monitoring
}

// TimezonePluginOption defines configuration options for the timezone plugin
type TimezonePluginOption func(*TimezonePlugin)

// NewTimezonePlugin creates a new timezone conversion plugin instance
// It initializes the plugin with default settings and applies any provided options
func NewTimezonePlugin(db *gorm.DB, opts ...TimezonePluginOption) *TimezonePlugin {
	plugin := &TimezonePlugin{
		name:              "timezone_converter",
		db:                db,
		convertibleFields: make(map[string]bool),
		logger:            logger.GetDefault(),
	}

	// Initialize convertible fields mapping from constants
	// These fields will be automatically converted during query results processing
	for _, field := range constants.GetTimezoneConvertibleFields() {
		plugin.convertibleFields[field] = true
	}

	// Apply configuration options
	for _, opt := range opts {
		opt(plugin)
	}

	return plugin
}

// WithLogger sets a custom logger for the timezone plugin
// This allows using a specific logger instance instead of the default one
func WithLogger(log logger.Logger) TimezonePluginOption {
	return func(p *TimezonePlugin) {
		p.logger = log
	}
}

// WithCustomFields adds custom field names to the convertible fields list
// This allows extending the plugin to handle additional time fields beyond the defaults
func WithCustomFields(fields ...string) TimezonePluginOption {
	return func(p *TimezonePlugin) {
		for _, field := range fields {
			p.convertibleFields[field] = true
		}
	}
}

// Name returns the plugin name identifier
// This is required by the GORM plugin interface
func (p *TimezonePlugin) Name() string {
	return p.name
}

// Initialize initializes the timezone plugin with GORM
// This method is called by GORM when the plugin is registered
// It sets up callbacks to automatically convert timezones after query operations
func (p *TimezonePlugin) Initialize(db *gorm.DB) error {
	// Register callback for after query operations
	// This handles SELECT queries and converts time fields in the results
	err := db.Callback().Query().After("gorm:after_query").Register("timezone:convert_after_query", p.afterQueryCallback)
	if err != nil {
		p.logger.Error("Failed to register after query callback", logger.ErrorField(err))
		return err
	}

	// Register callback for after find operations
	// This handles Find, First, Take and similar operations
	err = db.Callback().Query().After("gorm:after_find").Register("timezone:convert_after_find", p.afterQueryCallback)
	if err != nil {
		p.logger.Error("Failed to register after find callback", logger.ErrorField(err))
		return err
	}

	p.logger.Info("Timezone plugin initialized successfully")
	return nil
}

// afterQueryCallback is the main callback function executed after query operations
// It automatically converts time fields from UTC to the user's preferred timezone
func (p *TimezonePlugin) afterQueryCallback(db *gorm.DB) {
	p.logger.Debug("Timezone plugin callback triggered")

	// Check if timezone conversion should be skipped for this query
	// This includes write operations, count queries, and explicitly skipped queries
	if p.shouldSkipConversion(db) {
		return
	}

	// Retrieve the user's preferred timezone from context and database
	// Falls back to system timezone if user timezone is not available
	timezone := p.getUserTimezone(db.Statement.Context)
	p.logger.DebugContext(db.Statement.Context, "Retrieved timezone for conversion",
		logger.String(constants.UserTimezone, timezone))

	// Skip conversion if timezone is empty or already UTC
	if timezone == "" || timezone == constants.DefaultTimeZone {
		p.logger.DebugContext(db.Statement.Context, "No timezone conversion needed",
			logger.String("reason", "timezone is empty or UTC"))
		return
	}

	// Load the target timezone location
	// This validates the timezone string and creates a time.Location object
	targetTZ, err := time.LoadLocation(timezone)
	if err != nil {
		p.logger.WarnContext(db.Statement.Context, "Invalid timezone, using UTC",
			logger.String(constants.UserTimezone, timezone),
			logger.ErrorField(err))
		return
	}

	// Convert time fields in the query results to the target timezone
	p.logger.DebugContext(db.Statement.Context, "Starting timezone conversion",
		logger.String("targetTimezone", targetTZ.String()))
	p.convertTimezoneInResult(db, targetTZ)
}

// shouldSkipConversion determines whether timezone conversion should be skipped
// It checks for explicit skip flags, write operations, and aggregate queries
func (p *TimezonePlugin) shouldSkipConversion(db *gorm.DB) bool {
	// Check for explicit skip flag set by SkipTimezoneConversion helper
	if skip, exists := db.Get("timezone:skip"); exists && skip.(bool) {
		p.logger.Debug("Skipping: timezone:skip flag is set")
		return true
	}

	// Skip if there's no SQL statement (shouldn't happen in normal cases)
	if db.Statement.SQL.String() == "" {
		p.logger.Debug("Skipping: empty SQL statement")
		return true
	}

	// Analyze the SQL statement to determine if conversion is needed
	sql := strings.ToUpper(strings.TrimSpace(db.Statement.SQL.String()))
	sqlToLog := sql
	if len(sql) > constants.MaxSQLLogLength {
		sqlToLog = sql[:constants.MaxSQLLogLength] + "..."
	}
	p.logger.Debug("Checking SQL statement", logger.String("sql", sqlToLog))

	// Skip write operations (INSERT, UPDATE, DELETE)
	// These operations don't return data that needs timezone conversion
	if strings.HasPrefix(sql, "INSERT") || strings.HasPrefix(sql, "UPDATE") || strings.HasPrefix(sql, "DELETE") {
		p.logger.Debug("Skipping: write operation detected")
		return true
	}

	// Skip COUNT queries as they don't return time field data
	if strings.Contains(sql, "SELECT COUNT(") || strings.Contains(sql, "SELECT COUNT *") {
		p.logger.Debug("Skipping: COUNT query detected")
		return true
	}

	// Skip other aggregate function queries that don't return individual records
	if strings.Contains(sql, "SELECT SUM(") || strings.Contains(sql, "SELECT AVG(") ||
		strings.Contains(sql, "SELECT MAX(") || strings.Contains(sql, "SELECT MIN(") {
		p.logger.Debug("Skipping: aggregate function query detected")
		return true
	}

	p.logger.Debug("Not skipping conversion - proceeding with timezone conversion")
	return false
}

// getUserTimezone retrieves the user's preferred timezone from Redis cache first, then database
// It implements a fallback strategy: Redis cache -> user preference -> system setting -> UTC
func (p *TimezonePlugin) getUserTimezone(ctx context.Context) string {
	if ctx == nil || p.db == nil {
		return ""
	}

	// Extract user ID from the request context
	userID, exists := utils.GetUserIDFromContext(ctx)
	if !exists {
		p.logger.DebugContext(ctx, "No user ID found in context, using system timezone")
		return p.getSystemTimezone(ctx)
	}

	p.logger.DebugContext(ctx, "Found user ID in context for timezone lookup",
		logger.Uint("userID", userID))

	// Priority 1: Try to get timezone from Redis cache
	redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID)
	timezone, err := redis.HGet(ctx, redisKey, constants.UserTimezone)

	if err == nil && timezone != "" {
		p.logger.DebugContext(ctx, "Found user timezone in Redis cache",
			logger.Uint("userID", userID),
			logger.String(constants.UserTimezone, timezone),
			logger.String("redisKey", redisKey))
		return timezone
	}

	if err != nil && err.Error() != redis.RedisNilError {
		p.logger.DebugContext(ctx, "Failed to get timezone from Redis cache",
			logger.Uint("userID", userID),
			logger.String("redisKey", redisKey),
			logger.ErrorField(err))
	}

	// Priority 2: Query user timezone preference from user_profiles table
	// This allows each user to have their own timezone setting
	var userProfile struct {
		ConfigValue string `gorm:"column:config_value"`
	}

	err = SkipTimezoneConversion(p.db).WithContext(ctx).
		Table("user_profile").
		Select("config_value").
		Where("user_id = ? AND category = ? AND config_key = ?", userID, constants.UserCategory, constants.UserTimezone).
		First(&userProfile).Error

	if err == nil && userProfile.ConfigValue != "" {
		p.logger.DebugContext(ctx, "Found user timezone preference in database",
			logger.Uint("userID", userID),
			logger.String(constants.UserTimezone, userProfile.ConfigValue))
		return userProfile.ConfigValue
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		p.logger.DebugContext(ctx, "Failed to get user timezone preference from database",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
	}

	// Priority 3: Fall back to system timezone setting
	return p.getSystemTimezone(ctx)
}

// getSystemTimezone retrieves the system-wide timezone setting from configuration
// This serves as a fallback when user-specific timezone is not available
func (p *TimezonePlugin) getSystemTimezone(ctx context.Context) string {
	if p.db == nil {
		return constants.DefaultTimeZone
	}

	// Query system timezone configuration from system_configs table
	var systemConfig struct {
		ConfigValue string `gorm:"column:config_value"`
	}

	err := SkipTimezoneConversion(p.db).WithContext(ctx).
		Table("system_configs").
		Select("config_value").
		Where("config_key = ?", "system.timezone").
		First(&systemConfig).Error

	if err == nil && systemConfig.ConfigValue != "" {
		p.logger.DebugContext(ctx, "Found system timezone setting",
			logger.String(constants.UserTimezone, systemConfig.ConfigValue))
		return systemConfig.ConfigValue
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		p.logger.DebugContext(ctx, "Failed to get system timezone setting",
			logger.ErrorField(err))
	}

	// Final fallback to UTC timezone
	p.logger.DebugContext(ctx, "Using default timezone",
		logger.String("defaultTimezone", constants.DefaultTimeZone))
	return constants.DefaultTimeZone
}

// convertTimezoneInResult processes the query result and converts time fields
// It handles both single struct results and slice results (collections)
func (p *TimezonePlugin) convertTimezoneInResult(db *gorm.DB, targetTZ *time.Location) {
	if db.Statement.Dest == nil {
		p.logger.Debug("No destination to convert - Statement.Dest is nil")
		return
	}

	// Get the reflection value of the destination
	destValue := reflect.ValueOf(db.Statement.Dest)
	p.logger.Debug("Processing result for timezone conversion",
		logger.String("destType", destValue.Type().String()),
		logger.String("destKind", destValue.Kind().String()))

	// Dereference pointer if necessary
	if destValue.Kind() == reflect.Ptr {
		destValue = destValue.Elem()
		p.logger.Debug("Dereferenced pointer",
			logger.String("actualType", destValue.Type().String()),
			logger.String("actualKind", destValue.Kind().String()))
	}

	// Handle different result types
	switch destValue.Kind() {
	case reflect.Slice:
		// Handle slice results (typically from Find operations)
		// Iterate through each item in the slice and convert its time fields
		p.logger.Debug("Processing slice result", logger.Int("length", destValue.Len()))
		for i := 0; i < destValue.Len(); i++ {
			item := destValue.Index(i)
			p.logger.Debug("Converting slice item", logger.Int("index", i))
			p.convertStructTimezone(item, targetTZ)
		}
	case reflect.Struct:
		// Handle single struct results (typically from First, Take, etc.)
		p.logger.Debug("Processing single struct result")
		p.convertStructTimezone(destValue, targetTZ)
	default:
		p.logger.Debug("Unsupported result type for timezone conversion",
			logger.String("kind", destValue.Kind().String()))
	}
}

// convertStructTimezone processes a single struct and converts its time fields
// It handles embedded structs recursively and respects field visibility rules
func (p *TimezonePlugin) convertStructTimezone(structValue reflect.Value, targetTZ *time.Location) {
	p.logger.Debug("Starting struct timezone conversion",
		logger.String("structType", structValue.Type().String()),
		logger.String("structKind", structValue.Kind().String()),
		logger.Bool("canSet", structValue.CanSet()))

	// Dereference pointer if the struct value is a pointer
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return
		}
		structValue = structValue.Elem()
		p.logger.Debug("Dereferenced struct pointer",
			logger.String("actualType", structValue.Type().String()),
			logger.String("actualKind", structValue.Kind().String()),
			logger.Bool("canSet", structValue.CanSet()))
	}

	// Ensure we're working with a struct
	if structValue.Kind() != reflect.Struct {
		p.logger.Debug("Value is not a struct after dereferencing", logger.String("kind", structValue.Kind().String()))
		return
	}

	// Check if the struct can be modified
	if !structValue.CanSet() {
		p.logger.Debug("Struct cannot be set, skipping conversion")
		return
	}

	structType := structValue.Type()
	p.logger.Debug("Processing struct fields", logger.Int("numFields", structValue.NumField()))

	// Iterate through all fields in the struct
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		p.logger.Debug("Examining field",
			logger.Int("index", i),
			logger.String("fieldName", fieldType.Name),
			logger.String("fieldType", field.Type().String()),
			logger.String("jsonTag", fieldType.Tag.Get("json")),
			logger.Bool("anonymous", fieldType.Anonymous))

		// Skip unexported fields (cannot be modified)
		if !field.CanSet() {
			p.logger.Debug("Skipping field - cannot set", logger.String("fieldName", fieldType.Name))
			continue
		}

		// Handle embedded structs (like BaseModel) recursively
		if fieldType.Anonymous && field.Kind() == reflect.Struct {
			p.convertStructTimezone(field, targetTZ)
			continue
		}

		// Determine if this field should be converted
		shouldConvert := false
		fieldName := ""

		// Priority 1: Check JSON tag (preferred for API responses)
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			tagName := strings.Split(jsonTag, ",")[0]
			if tagName != "-" && p.isConvertibleField(tagName) {
				shouldConvert = true
				fieldName = tagName
			}
		}

		// Priority 2: Check struct field name if JSON tag doesn't match
		if !shouldConvert && p.isConvertibleField(fieldType.Name) {
			shouldConvert = true
			fieldName = fieldType.Name
		}

		if !shouldConvert {
			continue
		}

		p.logger.DebugContext(context.Background(), "Converting time field",
			logger.String("fieldName", fieldName),
			logger.String("structFieldName", fieldType.Name))

		// Convert the time field to target timezone
		p.convertTimeField(field, targetTZ)
	}
}

// isConvertibleField checks if a field name is in the convertible fields list
// This determines whether a field should undergo timezone conversion
func (p *TimezonePlugin) isConvertibleField(fieldName string) bool {
	return p.convertibleFields[fieldName]
}

// convertTimeField converts a time field value to the target timezone
// It handles both time.Time and *time.Time field types
func (p *TimezonePlugin) convertTimeField(field reflect.Value, targetTZ *time.Location) {
	switch field.Kind() {
	case reflect.Struct:
		// Handle time.Time type (value type)
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if timeVal := field.Interface().(time.Time); !timeVal.IsZero() {
				// Convert the time to the target timezone
				convertedTime := timeVal.In(targetTZ)
				p.logger.Debug("Converting time.Time field",
					logger.String("original", timeVal.String()),
					logger.String("converted", convertedTime.String()),
					logger.String("targetTZ", targetTZ.String()))
				field.Set(reflect.ValueOf(convertedTime))
			}
		}
	case reflect.Ptr:
		// Handle *time.Time type (pointer type)
		if field.Type() == reflect.TypeOf((*time.Time)(nil)) {
			if !field.IsNil() {
				timePtr := field.Interface().(*time.Time)
				if !timePtr.IsZero() {
					// Convert the time to the target timezone
					convertedTime := timePtr.In(targetTZ)
					p.logger.Debug("Converting *time.Time field",
						logger.String("original", timePtr.String()),
						logger.String("converted", convertedTime.String()),
						logger.String("targetTZ", targetTZ.String()))
					field.Set(reflect.ValueOf(&convertedTime))
				}
			}
		}
	}
}

// SkipTimezoneConversion is a helper function to skip timezone conversion for specific queries
// Usage: SkipTimezoneConversion(db).Where(...).Find(&result)
// This is useful for internal queries that should remain in UTC
func SkipTimezoneConversion(db *gorm.DB) *gorm.DB {
	return db.Set("timezone:skip", true)
}
