package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/logger"
)

type notificationChannelRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewNotificationChannelRepository creates a new notification channel repository instance
func NewNotificationChannelRepository(db *gorm.DB, logger logger.Logger) repository.NotificationChannelRepository {
	return &notificationChannelRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new notification channel
func (r *notificationChannelRepository) Create(ctx context.Context, channel *model.NotificationChannelConfig) error {
	r.logger.InfoContext(ctx, "Creating notification channel",
		logger.String("operation", "Create"),
		logger.String("code", channel.Code))

	if err := r.db.WithContext(ctx).Create(channel).Error; err != nil {
		r.logger.ErrorContext(ctx, "Failed to create notification channel", logger.ErrorField(err))
		return err
	}

	r.logger.InfoContext(ctx, "Notification channel created successfully",
		logger.Uint("id", channel.ID),
		logger.String("code", channel.Code))

	return nil
}

// GetByID retrieves a notification channel by ID
func (r *notificationChannelRepository) GetByID(ctx context.Context, id uint) (*model.NotificationChannelConfig, error) {
	r.logger.InfoContext(ctx, "Getting notification channel by ID",
		logger.String("operation", "GetByID"),
		logger.Uint("id", id))

	var channel model.NotificationChannelConfig
	if err := r.db.WithContext(ctx).First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.InfoContext(ctx, "Notification channel not found", logger.Uint("id", id))
		} else {
			r.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		}
		return nil, err
	}

	return &channel, nil
}

// GetByCode retrieves a notification channel by code
func (r *notificationChannelRepository) GetByCode(ctx context.Context, code string) (*model.NotificationChannelConfig, error) {
	r.logger.InfoContext(ctx, "Getting notification channel by code",
		logger.String("operation", "GetByCode"),
		logger.String("code", code))

	var channel model.NotificationChannelConfig
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&channel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.InfoContext(ctx, "Notification channel not found", logger.String("code", code))
		} else {
			r.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		}
		return nil, err
	}

	return &channel, nil
}

// Update updates a notification channel
func (r *notificationChannelRepository) Update(ctx context.Context, channel *model.NotificationChannelConfig) error {
	r.logger.InfoContext(ctx, "Updating notification channel",
		logger.String("operation", "Update"),
		logger.Uint("id", channel.ID),
		logger.String("code", channel.Code))

	if err := r.db.WithContext(ctx).Save(channel).Error; err != nil {
		r.logger.ErrorContext(ctx, "Failed to update notification channel", logger.ErrorField(err))
		return err
	}

	r.logger.InfoContext(ctx, "Notification channel updated successfully",
		logger.Uint("id", channel.ID),
		logger.String("code", channel.Code))

	return nil
}

// Delete hard deletes a notification channel by ID
func (r *notificationChannelRepository) Delete(ctx context.Context, id uint) error {
	r.logger.InfoContext(ctx, "Hard deleting notification channel",
		logger.String("operation", "Delete"),
		logger.Uint("id", id))

	// Use hard delete instead of soft delete
	if err := r.db.WithContext(ctx).Delete(&model.NotificationChannelConfig{}, id).Error; err != nil {
		r.logger.ErrorContext(ctx, "Failed to delete notification channel", logger.ErrorField(err))
		return err
	}

	r.logger.InfoContext(ctx, "Notification channel deleted successfully", logger.Uint("id", id))
	return nil
}

// GetList retrieves notification channels with pagination and filtering
func (r *notificationChannelRepository) GetList(ctx context.Context, req *request.GetNotificationChannelListRequest) ([]*model.NotificationChannelConfig, int64, error) {
	r.logger.InfoContext(ctx, "Getting notification channels list",
		logger.String("operation", "GetList"),
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize))

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
		r.logger.ErrorContext(ctx, "Failed to count notification channels", logger.ErrorField(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Get results with ordering
	var channels []*model.NotificationChannelConfig
	if err := query.Order("created_at DESC").Find(&channels).Error; err != nil {
		r.logger.ErrorContext(ctx, "Failed to get notification channels list", logger.ErrorField(err))
		return nil, 0, err
	}

	r.logger.InfoContext(ctx, "Notification channels list retrieved successfully",
		logger.Int64("total", total),
		logger.Int("count", len(channels)))

	return channels, total, nil
}

// CheckCodeExists checks if a channel code already exists (excluding a specific ID)
func (r *notificationChannelRepository) CheckCodeExists(ctx context.Context, code string, excludeID ...uint) (bool, error) {
	r.logger.InfoContext(ctx, "Checking if notification channel code exists",
		logger.String("operation", "CheckCodeExists"),
		logger.String("code", code))

	query := r.db.WithContext(ctx).Model(&model.NotificationChannelConfig{}).Where("code = ?", code)

	// Exclude specific ID if provided
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		r.logger.ErrorContext(ctx, "Failed to check code existence", logger.ErrorField(err))
		return false, err
	}

	exists := count > 0
	r.logger.InfoContext(ctx, "Code existence check completed",
		logger.String("code", code),
		logger.Bool("exists", exists))

	return exists, nil
}
