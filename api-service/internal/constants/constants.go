package constants

import "time"

// HTTP 状态码常量
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

// 应用状态常量
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

// 部署状态常量
const (
	DeploymentStatusPending  = "PENDING"
	DeploymentStatusRunning  = "RUNNING"
	DeploymentStatusSuccess  = "SUCCESS"
	DeploymentStatusFailed   = "FAILED"
	DeploymentStatusCanceled = "CANCELED"
)

// 时间相关常量
const (
	DefaultJWTExpireTime       = 3600 // 1小时
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

// 文件权限常量
const (
	DefaultDirPerm  = 0755
	DefaultFilePerm = 0644
)

// 默认值常量
const (
	DefaultBatchSize     = 100
	DefaultMaxRetries    = 3
	DefaultPort          = "8080"
	DefaultGinMode       = "release"
	DefaultLogLevel      = "info"
	DefaultMaxHeaderSize = 1 << 20 // 1MB
)

// 测试数据常量
const (
	TestCPUUsage    = 45.5
	TestMemoryUsage = 78.2
	TestDiskUsage   = 65.0
)

// 优先级常量
const (
	PriorityHigh   = 1
	PriorityMedium = 2
	PriorityLow    = 3
)

// 健康检查相关常量
const (
	HealthCheckTimeout  = 5 * time.Second
	HealthCheckInterval = 30 * time.Second
)
