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
	OAuth2StateLength  = 16
	OAuth2CookieMaxAge = 600 // 10 minutes
	HTTPClientTimeout  = 30  // seconds
	ProviderGoogle     = "google"
	ProviderGitHub     = "github"
	ProviderAzureAD    = "azure_ad"
	UserMappingEmail   = "email"
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
)

// Redis token storage constants
const (
	// JWTLeewaySeconds defines the leeway time in seconds for JWT validation
	JWTLeewaySeconds = 5
	// SMTPDefaultPort is the default SMTP port
	SMTPDefaultPort = 587
	// TokenBytes defines the size of random token bytes
	TokenBytes = 32

	// #nosec G101 -- This is not a credential, just a Redis key prefix
	TokenRedisKeyPrefix    = "AUTH:TOKEN:"
	TokenExpirationMinutes = 30 // 30 minutes for verification tokens
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

	// Authentication specific actions
	ActionLogin          = "LOGIN"           // Login operation
	ActionLogout         = "LOGOUT"          // Logout operation
	ActionRegister       = "REGISTER"        // Register operation
	ActionRefresh        = "REFRESH"         // Token refresh operation
	ActionChangePassword = "CHANGE_PASSWORD" // Change password operation
)

// Module constants based on security design document V1.1 section 5.2
const (
	ModulePlatform         = "Websoft9"
	ModuleHome             = "主页"
	ModuleProjectOverview  = "项目总览看板"
	ModuleMonitorOverview  = "监控总览看板"
	ModuleAppNavigation    = "应用快捷导航"
	ModulePersonalFolder   = "个人空间"
	ModuleProject          = "项目"
	ModuleProjectDashboard = "项目仪表盘"
	ModuleProjectFolder    = "项目空间"
	ModuleApplication      = "应用管理"
	ModuleWorkflow         = "工作流管理"
	ModuleJob              = "任务管理"
	ModuleProjectResource  = "资源"
	ModuleResourceGroup    = "资源组管理"
	ModuleServer           = "服务器管理"
	ModuleSecret           = "密钥管理" // #nosec G101 - This is a module name description, not a secret
	ModuleDatabase         = "数据库管理"
	ModuleGateway          = "应用网关管理"
	ModuleCertificate      = "证书管理"
	ModuleCloudResource    = "云资源管理"
	ModuleProjectTeam      = "项目团队管理"
	ModuleProjectSetting   = "项目设置"
	ModuleApps             = "应用"
	ModuleMarketplace      = "应用市场"
	ModuleWishlist         = "应用心愿单"
	ModuleAdminSetting     = "管理员设置"
	ModulePlatformSetting  = "平台管理"
	ModuleProjectManage    = "项目管理"
	ModuleSecurity         = "安全管理"
	ModuleRole             = "角色管理"
	ModulePermission       = "权限管理"
	ModuleAuth             = "认证管理"
	ModuleUser             = "用户管理"
	ModuleNotification     = "告警通知"
	ModuleProfile          = "个人中心"
	ModuleAuditLog         = "审计日志"
	ModuleSystem           = "System"
)

// GetModuleTableName returns the database table name for a module
func GetModuleTableName(module string) string {
	moduleToTableName := map[string]string{
		ModuleUser:           "users",
		ModuleRole:           "roles",
		ModulePermission:     "permissions",
		ModuleAuth:           "user_roles",
		ModuleAuditLog:       "audit_logs",
		ModuleProject:        "projects",
		ModuleProjectTeam:    "project_members",
		ModuleProjectSetting: "project_environments",
		ModuleProjectFolder:  "project_folders",
		ModuleApplication:    "app_instances",
		ModuleApps:           "app_instances",
		ModuleWorkflow:       "workflows",
		ModuleJob:            "workflow_tasks",
		ModuleResourceGroup:  "resource_groups",
		ModuleServer:         "servers",
		ModuleSecret:         "secret_keys",
		ModuleDatabase:       "database_connections",
		ModuleGateway:        "app_gateways",
		ModuleCertificate:    "ssl_certificates",
		ModuleCloudResource:  "cloud_resources",
		ModuleMarketplace:    "app_store_templates",
		ModuleWishlist:       "app_store_wishlists",
		ModuleNotification:   "notifications",
	}

	if tableName, exists := moduleToTableName[module]; exists {
		return tableName
	}
	return ""
}

// GetModuleType determines module name based module
func GetModuleType(module string) string {
	module_type_map := map[string]string{
		"users":             ModuleUser,
		"roles":             ModuleRole,
		"permissions":       ModulePermission,
		"auth":              ModuleAuth,
		"api-tokens":        ModuleAuth,
		"two-factor":        ModuleAuth,
		"auth-config":       ModuleAuth,
		"audit-logs":        ModuleAuditLog,
		"platform":          ModulePlatform,
		"home":              ModuleHome,
		"projects":          ModuleProject,
		"project-teams":     ModuleProjectTeam,
		"project-settings":  ModuleProjectSetting,
		"applications":      ModuleApplication,
		"apps":              ModuleApps,
		"marketplace":       ModuleMarketplace,
		"wishlist":          ModuleWishlist,
		"servers":           ModuleServer,
		"databases":         ModuleDatabase,
		"secrets":           ModuleSecret,
		"certificates":      ModuleCertificate,
		"resource-groups":   ModuleResourceGroup,
		"cloud-resources":   ModuleCloudResource,
		"gateways":          ModuleGateway,
		"workflows":         ModuleWorkflow,
		"jobs":              ModuleJob,
		"notifications":     ModuleNotification,
		"profile":           ModuleProfile,
		"i18n":              ModuleSystem,
		"admin-settings":    ModuleAdminSetting,
		"platform-settings": ModulePlatformSetting,
		"security":          ModuleSecurity,
	}
	if moduleName, exists := module_type_map[module]; exists {
		return moduleName
	}

	return ModuleSystem
}

// Export format constants
const (
	FormatExcel = "excel"
	FormatJSON  = "json"
	FormatCSV   = "csv"
)

// audit-logs constants
const (
	ExcelRowOffset        = 2      // Excel data starts from row 2 (after header)
	ExcelColumnDivisor    = 26     // Excel column calculation divisor (A-Z = 26 letters)
	MaxTimeRangeDays      = 7      // Maximum time range in days for export
	MinPathSegments       = 3      // Minimum path segments for valid API path
	DefaultTimeRangeHours = 7 * 24 // Default time range in hours (7 days)
)

// const (
// 	ErrInvalidCode        = "invalid code"
// 	ErrTokenNotFound      = "token not found"
// 	ErrPermissionNotFound = "permission not found"
// 	ErrRoleNotFound       = "role not found"
// )
