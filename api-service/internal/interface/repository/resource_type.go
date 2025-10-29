package repository

import (
	"context"

	"api-service/internal/model"
)

// ResourceTypeRepository defines the interface for resource type data access
type ResourceTypeRepository interface {
	// GetByCode retrieves a resource type by code
	GetByCode(ctx context.Context, code string) (*model.ResourceType, error)

	// List retrieves all resource types
	List(ctx context.Context) ([]*model.ResourceType, error)
}
