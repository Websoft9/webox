package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"

	"github.com/gin-gonic/gin"
)

// AuditLogService audit log business logic interface
type AuditLogService interface {
	// Record audit log (internal use)
	RecordLog(ctx context.Context, req *request.CreateAuditLogRequest) error

	// Middleware support methods
	ShouldSkipAudit(method, path string) bool
	RecordAuditFromRequest(backgroundCtx context.Context, ginCtx *gin.Context, responseBody []byte, responseTime int) error
	SubmitAuditJob(ginCtx *gin.Context, responseBody []byte, responseTime int) bool

	// Query operations
	GetAuditLog(ctx context.Context, id uint) (*response.AuditLogResponse, error)
	ListAuditLogs(ctx context.Context, req *request.ListAuditLogRequest) (*common.PaginationResponse, error)

	// Export functionality
	ExportAuditLogs(ctx context.Context, ginCtx *gin.Context, req *request.ExportAuditLogRequest) ([]byte, string, error)

	// System maintenance (internal use)
	CleanupExpiredLogs(ctx context.Context, retentionDays int) (int64, error)
}
