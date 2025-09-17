package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"context"
)

// SystemConfigRepository defines the interface for system configuration data access
type SystemConfigRepository interface {
	// Create creates a new system configuration
	Create(ctx context.Context, config *model.SystemConfig) error

	// GetByKey retrieves a system configuration by key
	GetByKey(ctx context.Context, key string) (*model.SystemConfig, error)

	// GetByID retrieves a system configuration by ID
	GetByID(ctx context.Context, id uint) (*model.SystemConfig, error)

	// List retrieves system configurations with filtering
	List(ctx context.Context, filter *request.ListSystemConfigsFilter) ([]*model.SystemConfig, error)

	// Update updates an existing system configuration
	Update(ctx context.Context, config *model.SystemConfig) error

	// UpdateValue updates the value of a system configuration by key
	UpdateValue(ctx context.Context, key, value string) error

	// Delete deletes a system configuration by ID
	Delete(ctx context.Context, id uint) error
}
