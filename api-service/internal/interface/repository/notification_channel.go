package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// NotificationChannelRepository defines the interface for notification channel repository
type NotificationChannelRepository interface {
	// Create creates a new notification channel
	Create(ctx context.Context, channel *model.NotificationChannelConfig) error

	// GetByID retrieves a notification channel by ID
	GetByID(ctx context.Context, id uint) (*model.NotificationChannelConfig, error)

	// GetByCode retrieves a notification channel by code
	GetByCode(ctx context.Context, code string) (*model.NotificationChannelConfig, error)

	// Update updates a notification channel
	Update(ctx context.Context, channel *model.NotificationChannelConfig) error

	// Delete soft deletes a notification channel by ID
	Delete(ctx context.Context, id uint) error

	// GetList retrieves notification channels with pagination and filtering
	GetList(ctx context.Context, req *request.GetNotificationChannelListRequest) ([]*model.NotificationChannelConfig, int64, error)

	// CheckCodeExists checks if a channel code already exists (excluding a specific ID)
	CheckCodeExists(ctx context.Context, code string, excludeID ...uint) (bool, error)
}
