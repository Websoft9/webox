package repository

import (
	"api-service/internal/model"
	"context"
	"time"
)

// AuditLogRepository audit log data access interface
type AuditLogRepository interface {
	// Basic operations
	Create(ctx context.Context, auditLog *model.AuditLog) error
	GetByID(ctx context.Context, id uint) (*model.AuditLog, error)
	List(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, int64, error)

	// Statistics and export
	GetStatistics(ctx context.Context, filter *StatisticsFilter) (*AuditLogStatistics, error)
	Export(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, error)

	// System maintenance (internal use only)
	CleanupOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)
}

// AuditLogFilter audit log query filter conditions
type AuditLogFilter struct {
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	UserID       *uint      `json:"user_id"`
	Username     string     `json:"username"`
	Action       string     `json:"action"`
	Module       string     `json:"module"`
	ResourceType string     `json:"resource_type"`
	ResourceID   *uint      `json:"resource_id"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	IPAddress    string     `json:"ip_address"`
	Success      *bool      `json:"success"`
}

// StatisticsFilter statistics query filter conditions
type StatisticsFilter struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	GroupBy   string     `json:"group_by"`
}

// AuditLogStatistics audit log statistics information
type AuditLogStatistics struct {
	TotalOperations   int64                `json:"total_operations"`
	SuccessOperations int64                `json:"success_operations"`
	FailedOperations  int64                `json:"failed_operations"`
	SuccessRate       float64              `json:"success_rate"`
	TopUsers          []UserOperationCount `json:"top_users"`
	TopActions        []ActionCount        `json:"top_actions"`
	Timeline          []TimelineCount      `json:"timeline"`
}

// UserOperationCount user operation statistics
type UserOperationCount struct {
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
	OperationCount int64  `json:"operation_count"`
}

// ActionCount action type statistics
type ActionCount struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// TimelineCount timeline statistics
type TimelineCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
