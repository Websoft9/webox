package repository

import (
	"api-service/internal/model"
	"context"
)

// ServiceConfigRepository defines the interface for service configuration data access
type ServiceConfigRepository interface {
	// Create creates a new service configuration
	Create(ctx context.Context, config *model.ServiceConfig) error

	// GetByCode retrieves a service configuration by code
	GetByCode(ctx context.Context, code string) (*model.ServiceConfig, error)

	// GetByID retrieves a service configuration by ID
	GetByID(ctx context.Context, id uint) (*model.ServiceConfig, error)

	// List retrieves all service configurations
	List(ctx context.Context) ([]*model.ServiceConfig, error)

	// ListByOwner retrieves service configurations by owner ID
	ListByOwner(ctx context.Context, ownerID uint) ([]*model.ServiceConfig, error)

	// ListByCategory retrieves service configurations by category
	ListByCategory(ctx context.Context, category string) ([]*model.ServiceConfig, error)

	// Update updates an existing service configuration
	Update(ctx context.Context, config *model.ServiceConfig) error

	// UpdateValue updates the value of a service configuration by code
	UpdateValue(ctx context.Context, code, value string) error

	// Delete deletes a service configuration by ID
	Delete(ctx context.Context, id uint) error
}
