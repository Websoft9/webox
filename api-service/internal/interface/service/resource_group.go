package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// ResourceGroupService defines the interface for resource group business logic
type ResourceGroupService interface {
	// CreateResourceGroup creates a new resource group
	CreateResourceGroup(ctx context.Context, req *request.CreateResourceGroupRequest, ownerID uint) (*response.ResourceGroupResponse, error)

	// GetResourceGroup retrieves a resource group by ID
	GetResourceGroup(ctx context.Context, id uint) (*response.ResourceGroupDetailResponse, error)

	// GetResourceGroupList retrieves a paginated list of resource groups
	GetResourceGroupList(ctx context.Context, req *request.GetResourceGroupListRequest, ownerID uint) (*common.PaginationResponse, error)

	// UpdateResourceGroup updates an existing resource group
	UpdateResourceGroup(ctx context.Context, id uint, req *request.UpdateResourceGroupRequest, ownerID uint) (*response.ResourceGroupResponse, error)

	// DeleteResourceGroup deletes a resource group
	DeleteResourceGroup(ctx context.Context, id uint) error

	// GetResourcesByGroupID retrieves resources in a resource group
	GetResourcesByGroupID(ctx context.Context, groupID uint, req *request.GetResourceGroupResourcesRequest) (*response.ResourceGroupResourcesResponse, error)

	// MoveResourcesToGroup moves multiple resources to a specific resource group (or default group if resource_group_id is nil)
	MoveResourcesToGroup(ctx context.Context, req *request.MoveResourcesToGroupRequest, ownerID uint) error

	// GetResourceStatistics gets resource statistics by project or resource group
	GetResourceStatistics(ctx context.Context, req *request.GetResourceStatisticsRequest) (*response.ResourceStatisticsResponse, error)
}
