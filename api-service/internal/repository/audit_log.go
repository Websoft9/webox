package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
)

const (
	// Default page size
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	// ErrRecordNotFound is returned when a record is not found
	ErrRecordNotFound = errors.New("record not found")
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
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetByID retrieves an audit log by ID
func (r *auditLogRepository) GetByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	var auditLog model.AuditLog

	err := r.db.WithContext(ctx).First(&auditLog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get audit log by ID: %w", err)
	}

	return &auditLog, nil
}

// List retrieves audit logs with pagination and filtering
func (r *auditLogRepository) List(ctx context.Context, filter *repository.AuditLogFilter) ([]*model.AuditLog, int64, error) {
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

	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}

	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}

	if filter.IPAddress != "" {
		// Search by IP address directly (no encryption)
		query = query.Where("ip_address = ?", filter.IPAddress)
	}

	if filter.Success != nil {
		query = query.Where("success = ?", *filter.Success)
	}

	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}

	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// Execute query with pagination and ordering
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&auditLogs).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	return auditLogs, total, nil
}

// GetStatistics retrieves audit log statistics
func (r *auditLogRepository) GetStatistics(ctx context.Context, filter *repository.StatisticsFilter) (*repository.AuditLogStatistics, error) {
	stats := &repository.AuditLogStatistics{}

	// Build base query
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// Apply time filters
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// Get total operations
	if err := query.Count(&stats.TotalOperations).Error; err != nil {
		return nil, fmt.Errorf("failed to count total operations: %w", err)
	}

	// Get success operations
	if err := query.Where("success = ?", true).Count(&stats.SuccessOperations).Error; err != nil {
		return nil, fmt.Errorf("failed to count success operations: %w", err)
	}

	// Calculate failed operations
	stats.FailedOperations = stats.TotalOperations - stats.SuccessOperations

	// Get unique users count
	var uniqueUsers int64
	if err := query.Select("COUNT(DISTINCT user_id)").Where("user_id IS NOT NULL").Count(&uniqueUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique users: %w", err)
	}

	// Get unique IPs count
	var uniqueIPs int64
	if err := query.Select("COUNT(DISTINCT ip_address)").Where("ip_address != '' AND ip_address IS NOT NULL").Count(&uniqueIPs).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique IPs: %w", err)
	}

	// Get top users - using Raw SQL to ensure proper field mapping
	var topUsers []repository.UserOperationCount
	if err := r.db.WithContext(ctx).Raw(`
		SELECT user_id, username, COUNT(*) as operation_count
		FROM audit_logs 
		WHERE user_id IS NOT NULL 
			AND (? IS NULL OR created_at >= ?)
			AND (? IS NULL OR created_at <= ?)
		GROUP BY user_id, username 
		ORDER BY operation_count DESC 
		LIMIT 10
	`, filter.StartTime, filter.StartTime, filter.EndTime, filter.EndTime).Scan(&topUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to get top users: %w", err)
	}
	stats.TopUsers = topUsers

	// Get top actions - using Raw SQL to ensure proper field mapping
	var topActions []repository.ActionCount
	if err := r.db.WithContext(ctx).Raw(`
		SELECT action, COUNT(*) as count
		FROM audit_logs 
		WHERE (? IS NULL OR created_at >= ?)
			AND (? IS NULL OR created_at <= ?)
		GROUP BY action 
		ORDER BY count DESC 
		LIMIT 10
	`, filter.StartTime, filter.StartTime, filter.EndTime, filter.EndTime).Scan(&topActions).Error; err != nil {
		return nil, fmt.Errorf("failed to get top actions: %w", err)
	}
	stats.TopActions = topActions

	// Get timeline data based on GroupBy
	var timeline []repository.TimelineCount
	var dateFormat string

	switch filter.GroupBy {
	case "hour":
		dateFormat = "strftime('%Y-%m-%d %H:00', created_at)"
	case "week":
		dateFormat = "strftime('%Y-W%W', created_at)"
	case "month":
		dateFormat = "strftime('%Y-%m', created_at)"
	default: // day
		dateFormat = "strftime('%Y-%m-%d', created_at)"
	}

	if err := query.Select(fmt.Sprintf("%s as date, COUNT(*) as count", dateFormat)).
		Group("date").
		Order("date").
		Scan(&timeline).Error; err != nil {
		return nil, fmt.Errorf("failed to get timeline data: %w", err)
	}
	stats.Timeline = timeline

	return stats, nil
}

// Export exports audit logs for external processing
func (r *auditLogRepository) Export(ctx context.Context, filter *repository.AuditLogFilter) ([]*model.AuditLog, error) {
	var auditLogs []*model.AuditLog

	// Build query
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// Apply filters (similar to List method but without pagination)
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}

	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}

	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}

	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// Execute query with ordering (no pagination for export)
	err := query.Order("created_at DESC").Find(&auditLogs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to export audit logs: %w", err)
	}

	return auditLogs, nil
}

// CleanupOldLogs deletes audit logs older than the specified date
func (r *auditLogRepository) CleanupOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", beforeDate).Delete(&model.AuditLog{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup old audit logs: %w", result.Error)
	}

	return result.RowsAffected, nil
}
