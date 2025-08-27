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
	DefaultJWTExpireTime       = 3600 // 1 hour
	TokenExpireHours           = 24   // 24 hours
	DefaultShutdownTimeout     = 30 * time.Second
	DefaultReadTimeout         = 30 * time.Second
	DefaultWriteTimeout        = 30 * time.Second
	DefaultReadHeaderTimeout   = 10 * time.Second
	DefaultIdleTimeout         = 120 * time.Second
	DefaultHealthCheckInterval = 30 * time.Second
	DefaultRetryInterval       = 10 * time.Second
	DefaultBatchProcessDelay   = 2 * time.Second
	DefaultMetricsInterval     = 60 * time.Second
)

// File permission constants
const (
	DefaultDirPerm  = 0755
	DefaultFilePerm = 0644
)

// Default value constants
const (
	DefaultBatchSize     = 100
	DefaultMaxRetries    = 3
	DefaultPort          = "8080"
	DefaultGinMode       = "release"
	DefaultLogLevel      = "info"
	DefaultMaxHeaderSize = 1 << 20 // 1MB
)

// Test data constants
const (
	TestCPUUsage    = 45.5
	TestMemoryUsage = 78.2
	TestDiskUsage   = 65.0
)

// Priority constants
const (
	PriorityHigh   = 1
	PriorityMedium = 2
	PriorityLow    = 3
)

// Health check related constants
const (
	HealthCheckTimeout  = 5 * time.Second
	HealthCheckInterval = 30 * time.Second
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

// Error message constants
const (
	ErrInvalidCode        = "invalid code"
	ErrTokenNotFound      = "token not found"
	ErrPermissionNotFound = "permission not found"
	ErrRoleNotFound       = "role not found"
)
