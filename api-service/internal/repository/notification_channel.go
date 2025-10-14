package repository

import (
	"context"
	"fmt"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"

	"gorm.io/gorm"
)

type notificationChannelRepository struct {
	db *gorm.DB
}

// NewNotificationChannelRepository creates a new notification channel repository instance
func NewNotificationChannelRepository(db *gorm.DB, logger logger.Logger) repository.NotificationChannelRepository {
	return &notificationChannelRepository{
		db: db,
	}
}

// Create creates a new notification channel
func (r *notificationChannelRepository) Create(ctx context.Context, channel *model.NotificationChannelConfig) error {
	if err := r.db.WithContext(ctx).Create(channel).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByID retrieves a notification channel by ID
func (r *notificationChannelRepository) GetByID(ctx context.Context, id uint) (*model.NotificationChannelConfig, error) {
	var channel model.NotificationChannelConfig
	if err := r.db.WithContext(ctx).First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &channel, nil
}

// GetByCode retrieves a notification channel by code
func (r *notificationChannelRepository) GetByCode(ctx context.Context, code string) (*model.NotificationChannelConfig, error) {
	var channel model.NotificationChannelConfig
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&channel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &channel, nil
}

// Update updates a notification channel
func (r *notificationChannelRepository) Update(ctx context.Context, channel *model.NotificationChannelConfig) error {
	if err := r.db.WithContext(ctx).Save(channel).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// Delete hard deletes a notification channel by ID
func (r *notificationChannelRepository) Delete(ctx context.Context, id uint) error {
	// Use hard delete instead of soft delete
	if err := r.db.WithContext(ctx).Delete(&model.NotificationChannelConfig{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	return nil
}

// GetList retrieves notification channels with pagination and filtering
func (r *notificationChannelRepository) GetList(ctx context.Context, req *request.GetNotificationChannelListRequest) ([]*model.NotificationChannelConfig, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.NotificationChannelConfig{})

	// Apply filters
	if req.ChannelType != "" {
		query = query.Where("channel_type = ?", req.ChannelType)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.OwnerID != nil {
		query = query.Where("owner_id = ?", *req.OwnerID)
	}

	if req.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", req.Search)
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Get results with ordering
	var channels []*model.NotificationChannelConfig
	if err := query.Order("created_at DESC").Find(&channels).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return channels, total, nil
}

// CheckCodeExists checks if a channel code already exists (excluding a specific ID)
func (r *notificationChannelRepository) CheckCodeExists(ctx context.Context, code string, excludeID ...uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&model.NotificationChannelConfig{}).Where("code = ?", code)

	// Exclude specific ID if provided
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	exists := count > 0
	return exists, nil
}
