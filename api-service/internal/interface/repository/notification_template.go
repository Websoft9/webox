package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// NotificationTemplateRepository defines the interface for notification template data access
type NotificationTemplateRepository interface {
	// Create creates a new notification template
	Create(ctx context.Context, template *model.NotificationTemplate) error

	// GetByID retrieves a notification template by ID
	GetByID(ctx context.Context, id uint) (*model.NotificationTemplate, error)

	// GetList retrieves a paginated list of notification templates with filters
	GetList(ctx context.Context, req *request.GetNotificationTemplateListRequest) ([]*model.NotificationTemplate, int64, error)

	// Update updates an existing notification template
	Update(ctx context.Context, template *model.NotificationTemplate) error

	// Delete deletes a notification template by ID
	Delete(ctx context.Context, id uint) error

	// ExistsByName checks if a template with the given name exists
	ExistsByName(ctx context.Context, name string, excludeID ...uint) (bool, error)
}
