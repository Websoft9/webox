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
	ModuleTag              = "标签管理"
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
		ModuleAdminSetting:   "system_configs",
		ModuleTag:            "tags",
		ModuleSystem:         "system",
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
		"system-configs":    ModuleAdminSetting,
		"tags":              ModuleTag,
		"system":            ModuleSystem,
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

// Secret key type constants
const (
	SecretTypeSSHKey      = "ssh_key"     // SSH密钥类型
	SecretTypePassword    = "password"    // 密码类型
	SecretTypeAPIKey      = "api_key"     // API密钥类型
	SecretTypeDatabase    = "database"    // 数据库凭据类型
	SecretTypeCertificate = "certificate" // 证书类型
)

// audit-logs constants
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
	TimeFieldCreatedAt = "created_at"
	TimeFieldUpdatedAt = "updated_at"
)

// GetTimezoneConvertibleFields returns a list of all time field names that should be converted to user's timezone
func GetTimezoneConvertibleFields() []string {
	return []string{
		// Common timestamp fields
		TimeFieldCreatedAt,
		TimeFieldUpdatedAt,
	}
}

// ===== 服务器管理数据字典与常量定义 (设计文档第8章) =====

// 服务器状态常量 (7.1.1 服务器基础信息)
const (
	ServerStatusOnline      = "online"      // 在线状态
	ServerStatusOffline     = "offline"     // 离线状态
	ServerStatusError       = "error"       // 错误状态
	ServerStatusUnreachable = "unreachable" // 不可达状态
	ServerStatusPending     = "pending"     // 待检测状态
	ServerStatusMaintenance = "maintenance" // 维护状态
)

// SSH状态常量
const (
	SSHStatusConnected    = "connected"    // SSH连接正常
	SSHStatusDisconnected = "disconnected" // SSH连接断开
	SSHStatusAuthFailed   = "auth_failed"  // SSH认证失败
	SSHStatusTimeout      = "timeout"      // SSH连接超时
)

// Agent状态常量
const (
	AgentStatusOnline   = "online"   // Agent在线
	AgentStatusOffline  = "offline"  // Agent离线
	AgentStatusError    = "error"    // Agent错误
	AgentStatusUpdating = "updating" // Agent更新中
)

// Docker状态常量
const (
	DockerStatusRunning = "running" // Docker运行中
	DockerStatusStopped = "stopped" // Docker已停止
	DockerStatusError   = "error"   // Docker错误
)

// 服务器操作类型常量 (7.1.2 服务器批量操作)
const (
	ServerActionRestart     = "restart"      // 重启服务器
	ServerActionShutdown    = "shutdown"     // 关闭服务器
	ServerActionReboot      = "reboot"       // 重新启动
	ServerActionUpdate      = "update"       // 更新系统
	ServerActionMaintenance = "maintenance"  // 进入维护模式
	ServerActionOnline      = "online"       // 上线
	ServerActionOffline     = "offline"      // 下线
	ServerActionHealthCheck = "health_check" // 健康检查
)

// 服务器检查类型常量 (7.1.3 服务器状态检查)
const (
	CheckTypeSSH    = "ssh"    // SSH连通性检查
	CheckTypeAgent  = "agent"  // Agent状态检查
	CheckTypeDocker = "docker" // Docker状态检查
	CheckTypeAll    = "all"    // 全部检查
)

// 操作系统类型常量
const (
	OSTypeLinux   = "linux"   // Linux系统
	OSTypeWindows = "windows" // Windows系统
	OSTypeMacOS   = "macos"   // macOS系统
	OSTypeUnknown = "unknown" // 未知系统
)

// Linux发行版常量
const (
	OSDistroUbuntu      = "ubuntu"   // Ubuntu
	OSDistroCentOS      = "centos"   // CentOS
	OSDistroRHEL        = "rhel"     // Red Hat Enterprise Linux
	OSDistroDebian      = "debian"   // Debian
	OSDistroFedora      = "fedora"   // Fedora
	OSDistroOpenSUSE    = "opensuse" // openSUSE
	OSDistroArch        = "arch"     // Arch Linux
	OSDistroAlpine      = "alpine"   // Alpine Linux
	OSDistroAmazonLinux = "amazon"   // Amazon Linux
)

// CPU架构常量
const (
	ArchitectureX86_64  = "x86_64"  // x86_64架构
	ArchitectureARM64   = "arm64"   // ARM64架构
	ArchitectureAARCH64 = "aarch64" // AARCH64架构
	ArchitectureARM     = "arm"     // ARM架构
	ArchitectureUnknown = "unknown" // 未知架构
)

// 文件操作常量 (7.1.4 服务器文件管理)
const (
	FileOperationUpload   = "upload"   // 文件上传
	FileOperationDownload = "download" // 文件下载
	FileOperationDelete   = "delete"   // 文件删除
)

// 文件大小限制常量
const (
	MaxFileUploadSize   = 100 * 1024 * 1024 // 100MB最大上传文件大小
	MaxFileDownloadSize = 500 * 1024 * 1024 // 500MB最大下载文件大小
)

// 网络连接超时常量
const (
	DefaultSSHTimeout        = 30 // SSH连接超时（秒）
	DefaultAgentTimeout      = 10 // Agent通信超时（秒）
	DefaultDockerTimeout     = 15 // Docker检查超时（秒）
	DefaultStatusCheckRetry  = 3  // 状态检查重试次数
	DefaultHeartbeatInterval = 60 // 心跳间隔（秒）
)

// 服务器配置项常量 (7.1 服务器管理配置项)
const (
	// 服务器管理模块配置键名
	ServerConfigCategoryServer = "server_management"

	// SSH连接配置
	ServerConfigSSHTimeout       = "server.ssh.timeout"        // SSH连接超时时间
	ServerConfigSSHRetryCount    = "server.ssh.retry_count"    // SSH连接重试次数
	ServerConfigSSHRetryInterval = "server.ssh.retry_interval" // SSH重试间隔
	ServerConfigSSHKeepAlive     = "server.ssh.keep_alive"     // SSH保持连接
	ServerConfigSSHDefaultPort   = "server.ssh.default_port"   // SSH默认端口

	// Agent配置
	ServerConfigAgentTimeout    = "server.agent.timeout"     // Agent通信超时
	ServerConfigAgentHeartbeat  = "server.agent.heartbeat"   // Agent心跳间隔
	ServerConfigAgentRetryCount = "server.agent.retry_count" // Agent重试次数

	// Docker配置
	ServerConfigDockerTimeout    = "server.docker.timeout"     // Docker检查超时
	ServerConfigDockerRetryCount = "server.docker.retry_count" // Docker重试次数

	// 状态检查配置
	ServerConfigStatusCheckInterval = "server.status.check_interval" // 状态检查间隔
	ServerConfigStatusCheckTimeout  = "server.status.check_timeout"  // 状态检查超时
	ServerConfigStatusCacheExpiry   = "server.status.cache_expiry"   // 状态缓存过期时间

	// 文件管理配置
	ServerConfigFileUploadMaxSize    = "server.file.upload_max_size"   // 文件上传最大大小
	ServerConfigFileDownloadMaxSize  = "server.file.download_max_size" // 文件下载最大大小
	ServerConfigFileOperationTimeout = "server.file.operation_timeout" // 文件操作超时

	// 批量操作配置
	ServerConfigBatchMaxCount    = "server.batch.max_count"   // 批量操作最大数量
	ServerConfigBatchTimeout     = "server.batch.timeout"     // 批量操作超时
	ServerConfigBatchConcurrency = "server.batch.concurrency" // 批量操作并发数
)

// 服务器管理默认配置值
const (
	// SSH默认配置
	DefaultSSHTimeoutValue  = 30 // 30秒
	DefaultSSHRetryCount    = 3  // 3次重试
	DefaultSSHRetryInterval = 5  // 5秒间隔
	DefaultSSHKeepAlive     = true
	DefaultSSHPort          = 22 // 默认SSH端口

	// Agent默认配置
	DefaultAgentTimeoutValue   = 10 // 10秒
	DefaultAgentHeartbeatValue = 60 // 60秒心跳
	DefaultAgentRetryCount     = 3  // 3次重试

	// Docker默认配置
	DefaultDockerTimeoutValue = 15 // 15秒
	DefaultDockerRetryCount   = 2  // 2次重试

	// 状态检查默认配置
	DefaultStatusCheckInterval = 300 // 5分钟检查间隔
	DefaultStatusCheckTimeout  = 60  // 60秒检查超时
	DefaultStatusCacheExpiry   = 180 // 3分钟缓存过期

	// 文件管理默认配置
	DefaultFileUploadMaxSize    = 100 * 1024 * 1024 // 100MB
	DefaultFileDownloadMaxSize  = 500 * 1024 * 1024 // 500MB
	DefaultFileOperationTimeout = 300               // 5分钟文件操作超时

	// 批量操作默认配置
	DefaultBatchMaxCount    = 50  // 最大50台服务器批量操作
	DefaultBatchTimeout     = 600 // 10分钟批量操作超时
	DefaultBatchConcurrency = 5   // 5个并发操作
)

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

// ==========================================
// File Management Security Constants
// ==========================================

// File path security mode
const (
	FilePathModeBlacklist = "blacklist" // 黑名单模式（默认）
	FilePathModeWhitelist = "whitelist" // 白名单模式
)

// File extension security mode
const (
	FileExtModeBlacklist = "blacklist" // 黑名单模式（默认）
	FileExtModeWhitelist = "whitelist" // 白名单模式
)

// Global forbidden paths (always forbidden, cannot be overridden)
var GlobalForbiddenPaths = []string{
	"/etc/passwd",
	"/etc/shadow",
	"/etc/sudoers",
	"/etc/sudoers.d/*",
	"/root/.ssh/*",
	"/home/*/.ssh/*",
	"/proc/*",
	"/sys/*",
	"~/.ssh/*",
	"/boot/*",
	"/dev/*",
}

// Global forbidden extensions (always forbidden, cannot be overridden)
var GlobalForbiddenExtensions = []string{
	".exe",
	".bat",
	".cmd",
	".com",
	".pif",
	".scr",
	".vbs",
	".js",
	".jse",
	".wsf",
	".wsh",
	".msi",
}

// Default allowed paths (whitelist mode)
var DefaultAllowedPaths = []string{
	"/tmp/*",
	"/var/tmp/*",
	"/home/*/uploads/*",
	"/opt/websoft9/data/*",
	"/data/*",
}

// Default allowed extensions (whitelist mode)
var DefaultAllowedExtensions = []string{
	".txt",
	".log",
	".conf",
	".yaml",
	".yml",
	".json",
	".xml",
	".md",
	".pdf",
	".zip",
	".tar",
	".tar.gz",
	".tgz",
	".csv",
	".sql",
}

// Default custom forbidden paths (blacklist mode, can be overridden)
var DefaultCustomForbiddenPaths = []string{
	"/var/log/auth.log",
	"/var/log/secure",
	"/etc/nginx/*",
	"/etc/apache2/*",
}

// Default custom forbidden extensions (blacklist mode, can be overridden)
var DefaultCustomForbiddenExtensions = []string{
	".sh",
	".bash",
	".py",
	".rb",
	".pl",
	".php",
	".jsp",
	".asp",
	".aspx",
}

// File size limits
const (
	MaxFileSizeLimit = 1073741824 // 1GB (hard limit)
)
