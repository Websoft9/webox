package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
)

// ResourceGroupRepository defines the interface for resource group data access
type ResourceGroupRepository interface {
	// Create creates a new resource group
	Create(ctx context.Context, rg *model.ResourceGroup) error

	// GetByID retrieves a resource group by ID
	GetByID(ctx context.Context, id uint) (*model.ResourceGroup, error)

	// GetByCode retrieves a resource group by code
	GetByCode(ctx context.Context, code string) (*model.ResourceGroup, error)

	// GetList retrieves a paginated list of resource groups with filters
	GetList(ctx context.Context, req *request.GetResourceGroupListRequest, ownerID uint) ([]*model.ResourceGroup, int64, error)

	// Update updates an existing resource group
	Update(ctx context.Context, rg *model.ResourceGroup) error

	// Delete deletes a resource group by ID
	Delete(ctx context.Context, id uint) error

	// CheckNameExists checks if a resource group name exists in a project
	CheckNameExists(ctx context.Context, projectID uint, name string, excludeID *uint) (bool, error)

	// GetDefaultResourceGroup retrieves the default resource group for a project
	GetDefaultResourceGroup(ctx context.Context, projectID uint) (*model.ResourceGroup, error)

	// GetProjectIDFromResourceCode retrieves the project_id from a resource code by querying the resource
	GetProjectIDFromResourceCode(ctx context.Context, resourceCode string) (uint, error)

	// GetResourcesByGroupID retrieves resources in a resource group
	GetResourcesByGroupID(ctx context.Context, groupID uint, req *request.GetResourceGroupResourcesRequest) ([]response.ResourceItemResponse, int64, error)

	// MoveResourcesToGroup moves multiple resources to a target resource group (or default group if targetGroupID is nil)
	MoveResourcesToGroup(ctx context.Context, resourceCodes []string, targetGroupID *uint, projectID uint) error

	// GetResourceStatistics gets resource statistics by project or resource group
	GetResourceStatistics(ctx context.Context, projectID *uint, resourceGroupID *uint) (map[string]int, error)
}
