package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"context"
	"time"
)

// AuditLogRepository audit log data access interface
type AuditLogRepository interface {
	// Basic operations
	Create(ctx context.Context, auditLog *model.AuditLog) error
	GetByID(ctx context.Context, id uint) (*model.AuditLog, error)
	List(ctx context.Context, filter *request.AuditLogFilter) ([]*model.AuditLog, int64, error)

	// Statistics and export
	GetStatistics(ctx context.Context, filter *request.StatisticsFilter) (*response.AuditLogStatistics, error)
	Export(ctx context.Context, filter *request.AuditLogFilter) ([]*model.AuditLog, error)

	// System maintenance (internal use only)
	CleanupOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)
}
