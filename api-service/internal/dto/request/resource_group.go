package request

import (
	"api-service/internal/dto/common"
	"errors"
)

// CreateResourceGroupRequest represents the request to create a resource group
type CreateResourceGroupRequest struct {
	ProjectID   uint    `json:"project_id" validate:"required"`           // Project ID
	Name        string  `json:"name" validate:"required,min=3,max=64"`    // Resource group name
	Description *string `json:"description" validate:"omitempty,max=500"` // Resource group description
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`    // Sort order
}

// UpdateResourceGroupRequest represents the request to update a resource group
type UpdateResourceGroupRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=64"`   // Resource group name
	Description *string `json:"description" validate:"omitempty,max=500"` // Resource group description
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`    // Sort order
}

// GetResourceGroupListRequest represents the request to get resource group list
type GetResourceGroupListRequest struct {
	common.BaseListRequest
	ProjectID *uint `form:"project_id" binding:"omitempty" json:"project_id"` // Filter by project ID
}

// GetResourceGroupResourcesRequest represents the request to get resources in a resource group
type GetResourceGroupResourcesRequest struct {
	common.BaseListRequest
	ResourceType *string `form:"resource_type" binding:"omitempty" json:"resource_type"` // Filter by resource type
}

// BatchManageResourcesRequest represents the request to add or remove resources from a resource group
type BatchManageResourcesRequest struct {
	ResourceCodes []string `json:"resource_codes" validate:"required,min=1,dive,required"` // Array of resource codes
}

// MoveResourcesToGroupRequest represents the request to move resources to a specific resource group
// The resource_group_id field is required, but can be null to move resources to the default group
type MoveResourcesToGroupRequest struct {
	ResourceCodes   []string `json:"resource_codes" validate:"required,min=1,dive,required"` // Array of resource codes (format: {resource_type}_{id})
	ResourceGroupID *uint    `json:"resource_group_id"`                                      // Target resource group ID (required field, null for default group)
}

// GetResourceStatisticsRequest represents the request to get resource statistics
type GetResourceStatisticsRequest struct {
	ProjectID       *uint `form:"project_id" binding:"omitempty" json:"project_id"`               // Filter by project ID
	ResourceGroupID *uint `form:"resource_group_id" binding:"omitempty" json:"resource_group_id"` // Filter by resource group ID
}

// Validate validates the GetResourceStatisticsRequest
func (r *GetResourceStatisticsRequest) Validate() error {
	// Mutual exclusivity check: project_id and resource_group_id cannot be provided together
	if r.ProjectID != nil && r.ResourceGroupID != nil {
		return errors.New("project_id and resource_group_id are mutually exclusive")
	}
	return nil
}
