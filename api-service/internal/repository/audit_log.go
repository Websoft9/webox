package repository

import (
	"api-service/pkg/errors"
	"context"
	"time"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"

	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates audit log repository instance with relational database
func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &auditLogRepository{
		db: db,
	}
}

// Create creates a new audit log record
func (r *auditLogRepository) Create(ctx context.Context, auditLog *model.AuditLog) error {
	// Create the record directly without encryption
	if err := r.db.WithContext(ctx).Create(auditLog).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByID retrieves an audit log by ID
func (r *auditLogRepository) GetByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	var auditLog model.AuditLog

	err := r.db.WithContext(ctx).First(&auditLog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &auditLog, nil
}

// List retrieves audit logs with pagination and filtering
func (r *auditLogRepository) List(ctx context.Context, filter *request.ListAuditLogRequest) ([]*model.AuditLog, int64, error) {
	var auditLogs []*model.AuditLog
	var total int64

	// Build base query
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}

	if filter.IPAddress != "" {
		// Search by IP address directly (no encryption)
		query = query.Where("ip_address = ?", filter.IPAddress)
	}

	startTime, endTime, parseErr := filter.GetTimeRange(true)
	if parseErr != nil {
		return nil, 0, errors.NewAppErrorWrapError(parseErr, errors.CodeRecordQueryFailed)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	if filter.Success != nil {
		query = query.Where("success = ?", *filter.Success)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Paginated query
	offset := filter.GetOffset()
	limit := filter.GetPageSize()

	// Execute query with pagination and ordering
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return auditLogs, total, nil
}

// Export exports audit logs for external processing
func (r *auditLogRepository) Export(ctx context.Context, filter *request.ExportAuditLogRequest) ([]*model.AuditLog, error) {
	var auditLogs []*model.AuditLog

	// Build query
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// Apply filters (similar to List method but without pagination)
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.StartTime != "" && filter.EndTime != "" {
		startTime, _ := time.Parse(constants.DefaultTimeFormat, filter.StartTime)
		endTime, _ := time.Parse(constants.DefaultTimeFormat, filter.EndTime)
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	// Execute query with ordering (no pagination for export)
	err := query.Order("created_at DESC").Find(&auditLogs).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return auditLogs, nil
}

// CleanupOldLogs deletes audit logs older than the specified date
func (r *auditLogRepository) CleanupOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", beforeDate).Delete(&model.AuditLog{})

	if result.Error != nil {
		return 0, errors.NewAppErrorWrapError(result.Error, errors.CodeRecordQueryFailed)
	}

	return result.RowsAffected, nil
}
