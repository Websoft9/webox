package constants

import "time"

// HTTP status code constants
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
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
	ModuleProjectOverview = "PROJECT_OVERVIEW" // 项目总览看板
	ModuleMonitorOverview = "MONITOR_OVERVIEW" // 监控总览看板
	ModuleAppNavigation   = "APP_NAVIGATION"   // 应用快捷导航
	ModuleDashboard       = "DASHBOARD"        // 项目仪表盘
	ModuleMonitorBoard    = "MONITOR_BOARD"    // 项目监控看板
	ModuleTaskBoard       = "TASK_BOARD"       // 项目任务看板
	ModuleResourceBoard   = "RESOURCE_BOARD"   // 项目资源看板
	ModuleProjectFolder   = "PROJECT_FOLDER"   // 项目文件夹管理
	ModulePersonalFolder  = "PERSONAL_FOLDER"  // 个人文件夹管理
	ModuleApplication     = "APPLICATION"      // 应用管理
	ModuleWorkflow        = "WORKFLOW"         // 工作流管理
	ModuleJob             = "JOB"              // 任务管理
	ModuleResourceGroup   = "RESOURCE_GROUP"   // 资源组管理
	ModuleServer          = "SERVER"           // 服务器管理
	ModuleSecret          = "SECRET"           // 密钥管理
	ModuleDatabase        = "DATABASE"         // 数据库管理
	ModuleGateway         = "GATEWAY"          // 应用网关管理
	ModuleCertificate     = "CERTIFICATE"      // 证书管理
	ModuleCloudResource   = "CLOUD_RESOURCE"   // 云资源管理
	ModuleProjectTeam     = "PROJECT_TEAM"     // 项目团队管理
	ModuleProjectSetting  = "PROJECT_SETTING"  // 项目设置
	ModuleMarketplace     = "MARKETPLACE"      // 应用市场
	ModuleWishlist        = "WISHLIST"         // 应用心愿单
	ModuleProject         = "PROJECT"          // 项目管理
	ModulePlatformSetting = "PLATFORM_SETTING" // 平台设置
	ModuleRole            = "ROLE"             // 角色管理
	ModulePermission      = "PERMISSION"       // 权限管理
	ModuleAuth            = "AUTH"             // 认证管理
	ModuleUser            = "USER"             // 用户管理
	ModuleNotification    = "NOTIFICATION"     // 告警通知管理
	ModuleProfile         = "PROFILE"          // 个人中心
	ModuleAuditLog        = "AUDIT_LOG"        // 审计日志
)

// GetModuleTableName returns the database table name for a module
func GetModuleTableName(module string) string {
	moduleToTableName := map[string]string{
		ModuleUser:        "users",
		ModuleRole:        "roles",
		ModulePermission:  "permissions",
		ModuleAuth:        "user_roles",
		ModuleAuditLog:    "audit_logs",
		ModuleApplication: "applications",
		ModuleProject:     "projects",
		ModuleServer:      "servers",
		ModuleDatabase:    "databases",
		ModuleSecret:      "secrets",
		ModuleCertificate: "certificates",
	}

	if tableName, exists := moduleToTableName[module]; exists {
		return tableName
	}
	return ""
}

// Error message constants
const (
	ErrInvalidCode        = "invalid code"
	ErrTokenNotFound      = "token not found"
	ErrPermissionNotFound = "permission not found"
	ErrRoleNotFound       = "role not found"
)
