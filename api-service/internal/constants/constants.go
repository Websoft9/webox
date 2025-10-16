package constants

import "time"

// HTTP status code constants
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
)

// Standard error codes mapping to HTTP status codes
const (
	// Success response
	SUCCESS = 200

	// Client error responses
	VALIDATION_ERROR     = 400 // Parameter validation error
	AUTHENTICATION_ERROR = 401 // Authentication failed
	AUTHORIZATION_ERROR  = 403 // Insufficient permissions
	NOT_FOUND_ERROR      = 404 // Resource not found
	CONFLICT_ERROR       = 409 // Resource conflict
	BUSINESS_ERROR       = 422 // Business logic error
	RATE_LIMIT_ERROR     = 429 // Request rate limit exceeded

	// Server error responses
	INTERNAL_ERROR      = 500 // Internal server error
	GATEWAY_ERROR       = 502 // Gateway error
	SERVICE_UNAVAILABLE = 503 // Service unavailable
)

// Application status constants
const (
	AppStatusPending   = "pending"
	AppStatusRunning   = "running"
	AppStatusStopped   = "stopped"
	AppStatusFailed    = "failed"
	AppStatusSuccess   = "success"
	AppStatusHealthy   = "healthy"
	AppStatusUnhealthy = "unhealthy"
	AppStatusCanceled  = "canceled"
)

// System config type constants
const (
	ConfigTypeString  = "STRING"
	ConfigTypeBoolean = "BOOLEAN"
	ConfigTypeNumber  = "NUMBER"
	ConfigTypeJSON    = "JSON"
)

// System config category constants
const (
	CategoryBasic    = "basic"
	CategorySMTP     = "smtp"
	CategorySMS      = "sms"
	CategorySecurity = "security"
	CategorySystem   = "system"
	CategoryLicense  = "license"
)

// Deployment status constants
const (
	DeploymentStatusPending  = "PENDING"
	DeploymentStatusRunning  = "RUNNING"
	DeploymentStatusSuccess  = "SUCCESS"
	DeploymentStatusFailed   = "FAILED"
	DeploymentStatusCanceled = "CANCELED"
)

// Time-related constants
const (
	DefaultJWTExpireTime = 3600 // 1 hour
	TokenExpireHours     = 24   // 24 hours
)

// File permission constants
const (
	DefaultDirPerm  = 0755
	DefaultFilePerm = 0644
)

// Default value constants
const (
	DefaultPort              = "8080"
	DefaultGinMode           = "release"
	DefaultLogLevel          = "info"
	DefaultShutdownTimeout   = 30 * time.Second
	DefaultReadTimeout       = 30 * time.Second
	DefaultWriteTimeout      = 30 * time.Second
	DefaultReadHeaderTimeout = 10 * time.Second
	DefaultIdleTimeout       = 120 * time.Second
	DefaultMaxHeaderBytes    = 1 << 20 // 1MB
)

// Log configuration constants
const (
	DefaultLogMaxSize    = 10 // 10 MB
	DefaultLogMaxBackups = 5  // 5 backup files
	DefaultLogMaxAge     = 30 // 30 days
)

// Database configuration constants
const (
	DefaultMySQLPort       = 3306
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
	DefaultConnMaxLifetime = 3600 // seconds
	DefaultConnectTimeout  = 30   // seconds
)

// Logger directory permissions
const (
	DefaultLogDirPerm = 0755
)

// OAuth2 constants
const (
	OAuth2StateLength   = 16
	OAuth2CookieMaxAge  = 600 // 10 minutes
	HTTPClientTimeout   = 30  // seconds
	ProviderGoogle      = "google"
	ProviderGitHub      = "github"
	ProviderAzureAD     = "azure_ad"
	UserMappingEmail    = "email"
	UserMappingUsername = "username"
)

// Tag management constants
const (
	TagMaxNameLength          = 128
	TagMaxBatchSize           = 50
	TagMaxColorLength         = 16
	TagMaxDescLength          = 500
	TagSearchOpAND            = "AND"
	TagSearchOpOR             = "OR"
	TagAssignStatusCreated    = "created"
	TagAssignStatusAssociated = "associated"
	DefaultTagPageSize        = 20
	MaxTagPageSize            = 100
)

// Auth config constants
const (
	DefaultTokenExpiresIn        = 3600  // 1 hour
	DefaultRefreshTokenExpiresIn = 86400 // 24 hours
	PasswordMinLength            = 8
	PasswordMaxLength            = 128
	PasswordHistoryCount         = 5
	PasswordExpiresDays          = 90
	MaxLoginAttempts             = 5
	LockoutDuration              = 1800 // 30 minutes
	TOTPDigits                   = 6
	TOTPPeriod                   = 30
	BackupCodesCount             = 10
	SessionTimeout               = 3600 // 1 hour
	MaxConcurrentSessions        = 5
	RememberMeDuration           = 2592000 // 30 days
	DirPerm                      = 0755
	RefreshTokenMultiplier       = 24
	PasswordMinLengthCheck       = 8
	TimeRangePartsCount          = 2

	// 2FA method constants
	TwoFactorMethodTOTP   = "totp"
	TwoFactorMethodEmail  = "email"
	TwoFactorMethodBackup = "backup"

	// SMTPDefaultPort is the default SMTP port
	SMTPDefaultPort = 587
)

// HTTP method constants
const (
	HTTPMethodGET    = "GET"
	HTTPMethodPOST   = "POST"
	HTTPMethodPUT    = "PUT"
	HTTPMethodPATCH  = "PATCH"
	HTTPMethodDELETE = "DELETE"
)

// Action operation permission constants (based on security design document V1.1)
const (
	ActionAll    = "*"      // All operations
	ActionCreate = "CREATE" // Create operation
	ActionUpdate = "UPDATE" // Update/Modify operation
	ActionDelete = "DELETE" // Delete operation
	ActionQuery  = "QUERY"  // Query operation
)

// Export format constants
const (
	FormatExcel = "excel"
	FormatJSON  = "json"
	FormatCSV   = "csv"
)

// Time range and audit-logs constants
const (
	ExcelRowOffset        = 2      // Excel data starts from row 2 (after header)
	ExcelColumnDivisor    = 26     // Excel column calculation divisor (A-Z = 26 letters)
	MaxTimeRangeDays      = 7      // Maximum allowed time range in days (configurable)
	MinPathSegments       = 3      // Minimum path segments for valid API path
	DefaultTimeRangeHours = 7 * 24 // Default time range in hours (7 days)
)

// alert constants
const (
	AlertStatusFiring    = "FIRING"
	AlertStatusConfirmed = "CONFIRMED"
	AlertStatusResolved  = "RESOLVED"
)

const (
	// StringTrue represents the string "true"
	StringTrue = "true"

	// StringFalse represents the string "false"
	StringFalse = "false"

	// Session timeout in seconds (30 minutes)
	DefaultSessionTimeoutSeconds = 1800
)

// Time field names constants for timezone conversion middleware
// These field names will be automatically converted to user's timezone in API responses
const (
	// Common timestamp fields
	TimeFieldCreatedAt  = "created_at"
	TimeFieldUpdatedAt  = "updated_at"
	TimeFieldLoginAt    = "last_login_at"
	TimeFieldUsedAt     = "last_used_at"
	TimeFieldExpiresAt  = "expires_at"
	TimeFieldGrantedAt  = "granted_at"
	TimeFieldVerifiedAt = "verified_at"
	TimeFieldSentAt     = "sent_at"
	TimeFieldLogin      = "login_time"
	TimeFieldLogout     = "logout_time"
)

// GetTimezoneConvertibleFields returns a list of all time field names that should be converted to user's timezone
func GetTimezoneConvertibleFields() []string {
	return []string{
		// Common timestamp fields
		TimeFieldCreatedAt,
		TimeFieldUpdatedAt,
		TimeFieldLoginAt,
		TimeFieldUsedAt,
		TimeFieldExpiresAt,
		TimeFieldGrantedAt,
		TimeFieldVerifiedAt,
		TimeFieldSentAt,
		TimeFieldLogin,
		TimeFieldLogout,
	}
}

// User preferences config keys
const (
	UserCategory = "general"
	UserLanguage = "language"
	UserTimezone = "timezone"
)

// System preferences config keys
const (
	SystemLanguage = "system." + UserLanguage
	SystemTimezone = "system." + UserTimezone

	// Default time format for datetime display: "2006-01-02T15:04:05Z07:00"
	DefaultTimeFormat = time.RFC3339

	// Default language code
	DefaultLanguage = "en-US"

	// Default time zone
	DefaultTimeZone = "UTC"

	// Maximum length of SQL log
	MaxSQLLogLength = 100
)

// Secret key file types
const (
	FileTypeKey = ".key" // Private key file
	FileTypePem = ".pem" // PEM encoded certificate or key
	FileTypeRsa = ".rsa" // RSA key file
	FileTypeCrt = ".crt" // Certificate file
)

// Allowed secret file extensions
var AllowedSecretFileExtensions = []string{
	FileTypeKey,
	FileTypePem,
	FileTypeRsa,
	FileTypeCrt,
}

// NotificationRecordStatus enum values
const (
	NotificationStatusPending = "PENDING"
	NotificationStatusSent    = "SENT"
	NotificationStatusFailed  = "FAILED"
	NotificationStatusRetry   = "RETRY"
)

// NotificationChannelType enum values
const (
	NotificationChannelEmail    = "EMAIL"
	NotificationChannelWebhook  = "WEBHOOK"
	NotificationChannelInternal = "INTERNAL"
)
