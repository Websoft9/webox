package constants

import "time"

// 应用状态常量
const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
	StatusFailed    = "failed"
	StatusSuccess   = "success"
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusStopped   = "stopped"
)

// 网络相关常量
const (
	DefaultHost = "localhost"
	DefaultPort = "8080"
)

// 时间相关常量
const (
	DefaultStartupDelay        = 2 * time.Second
	DefaultHealthCheckInterval = 30 * time.Second
	DefaultRetryInterval       = 10 * time.Second
	DefaultMetricsInterval     = 60 * time.Second
	DefaultHTTPTimeout         = 10 * time.Second
	DefaultTCPTimeout          = 5 * time.Second
	DefaultCheckInterval       = 30 * time.Second
)

// 测试数据常量 (用于模拟数据)
const (
	TestCPUUsage    = 45.5
	TestMemoryUsage = 78.2
	TestDiskUsage   = 65.0
)

// 默认配置常量
const (
	DefaultBatchSize   = 100
	DefaultMaxRetries  = 3
	DefaultBufferSize  = 1024
	DefaultWorkerCount = 5
)
