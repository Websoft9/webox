package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// resourceGroupService implements ResourceGroupService
type resourceGroupService struct {
	repo   interfaceRepo.ResourceGroupRepository
	logger logger.Logger
}

// NewResourceGroupService creates a new resource group service
func NewResourceGroupService(
	repo interfaceRepo.ResourceGroupRepository,
	logger logger.Logger,
) (interfaceService.ResourceGroupService, error) {
	return &resourceGroupService{
		repo:   repo,
		logger: logger,
	}, nil
}

// CreateResourceGroup creates a new resource group
func (s *resourceGroupService) CreateResourceGroup(
	ctx context.Context,
	req *request.CreateResourceGroupRequest,
	ownerID uint,
) (*response.ResourceGroupResponse, error) {
	// Check if name already exists in the project
	exists, err := s.repo.CheckNameExists(ctx, req.ProjectID, req.Name, nil)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check resource group name existence",
			logger.String("name", req.Name),
			logger.Uint("project_id", req.ProjectID),
			logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Set default sort order if not provided
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	// Create resource group model
	rg := &model.ResourceGroup{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		SortOrder:   sortOrder,
	}

	// Save to database (code will be generated in repository)
	if err := s.repo.Create(ctx, rg); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create resource group",
			logger.String("name", req.Name),
			logger.Uint("project_id", req.ProjectID),
			logger.Uint("owner_id", ownerID),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Resource group created successfully",
		logger.Uint("resource_group_id", rg.ID),
		logger.String("code", rg.Code),
		logger.String("name", rg.Name),
		logger.Uint("project_id", req.ProjectID),
		logger.Uint("owner_id", ownerID))

	return s.toResponse(rg), nil
}

// GetResourceGroup retrieves a resource group by ID
func (s *resourceGroupService) GetResourceGroup(ctx context.Context, id uint) (*response.ResourceGroupDetailResponse, error) {
	rg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(rg), nil
}

// GetResourceGroupList retrieves a paginated list of resource groups
func (s *resourceGroupService) GetResourceGroupList(ctx context.Context, req *request.GetResourceGroupListRequest, ownerID uint) (*common.PaginationResponse, error) {
	resourceGroups, total, err := s.repo.GetList(ctx, req, ownerID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(resourceGroups))
	for i, rg := range resourceGroups {
		items[i] = s.toResponse(rg)
	}

	return common.NewPaginationResponse(req.Page, req.PageSize, total, items), nil
}

// UpdateResourceGroup updates an existing resource group
func (s *resourceGroupService) UpdateResourceGroup(
	ctx context.Context,
	id uint,
	req *request.UpdateResourceGroupRequest,
	ownerID uint,
) (*response.ResourceGroupResponse, error) {
	rg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		// Check if new name already exists in the project (excluding current record)
		exists, err := s.repo.CheckNameExists(ctx, rg.ProjectID, *req.Name, &id)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to check resource group name existence",
				logger.String("name", *req.Name),
				logger.Uint("project_id", rg.ProjectID),
				logger.ErrorField(err))
			return nil, err
		}
		if exists {
			return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
		}
		rg.Name = *req.Name
	}
	if req.Description != nil {
		rg.Description = req.Description
	}
	if req.SortOrder != nil {
		rg.SortOrder = *req.SortOrder
	}

	// Save to database
	if err := s.repo.Update(ctx, rg); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update resource group",
			logger.Uint("resource_group_id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Resource group updated successfully",
		logger.Uint("resource_group_id", id),
		logger.String("name", rg.Name))

	return s.toResponse(rg), nil
}

// DeleteResourceGroup deletes a resource group
func (s *resourceGroupService) DeleteResourceGroup(ctx context.Context, id uint) error {
	rg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if resource group is in use (has associated resources)
	// For now, we allow deletion. In the future, add validation here.

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete resource group",
			logger.Uint("resource_group_id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Resource group deleted successfully",
		logger.Uint("resource_group_id", id),
		logger.String("name", rg.Name))

	return nil
}

// toResponse converts model to response DTO
func (s *resourceGroupService) toResponse(rg *model.ResourceGroup) *response.ResourceGroupResponse {
	return &response.ResourceGroupResponse{
		ID:          rg.ID,
		ProjectID:   rg.ProjectID,
		Name:        rg.Name,
		Code:        rg.Code,
		Description: rg.Description,
		OwnerID:     rg.OwnerID,
		IsDefault:   rg.IsDefault,
		SortOrder:   rg.SortOrder,
		CreatedAt:   rg.CreatedAt,
		UpdatedAt:   rg.UpdatedAt,
	}
}

// toDetailResponse converts model to detail response DTO
func (s *resourceGroupService) toDetailResponse(rg *model.ResourceGroup) *response.ResourceGroupDetailResponse {
	return &response.ResourceGroupDetailResponse{
		ResourceGroupResponse: *s.toResponse(rg),
	}
}

// GetResourcesByGroupID retrieves resources in a resource group
func (s *resourceGroupService) GetResourcesByGroupID(
	ctx context.Context,
	groupID uint,
	req *request.GetResourceGroupResourcesRequest,
) (*response.ResourceGroupResourcesResponse, error) {
	// Get resources
	resources, total, err := s.repo.GetResourcesByGroupID(ctx, groupID, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get resources by group ID",
			logger.Uint("resource_group_id", groupID),
			logger.ErrorField(err))
		return nil, err
	}

	return &response.ResourceGroupResourcesResponse{
		Total:     total,
		Resources: resources,
	}, nil
}

// MoveResourcesToGroup moves multiple resources to a specific resource group or default group
// If resource_group_id is null, resources are moved to the project's default resource group
func (s *resourceGroupService) MoveResourcesToGroup(
	ctx context.Context,
	req *request.MoveResourcesToGroupRequest,
	ownerID uint,
) error {
	var targetGroupID *uint
	var projectID uint

	// Check if target resource group is specified (not nil)
	if req.ResourceGroupID != nil {
		// Validate target resource group existence
		targetGroup, err := s.repo.GetByID(ctx, *req.ResourceGroupID)
		if err != nil {
			return err
		}

		targetGroupID = req.ResourceGroupID
		projectID = targetGroup.ProjectID
	} else {
		// resource_group_id is null, move to default group
		// Get project_id from the first resource code
		firstProjectID, err := s.repo.GetProjectIDFromResourceCode(ctx, req.ResourceCodes[0])
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to get project_id from resource code",
				logger.String("resource_code", req.ResourceCodes[0]),
				logger.ErrorField(err))
			return err
		}

		projectID = firstProjectID
		targetGroupID = nil // Repository layer will get the default group automatically
	}

	// Move resources to target group (or default group if targetGroupID is nil)
	if err := s.repo.MoveResourcesToGroup(ctx, req.ResourceCodes, targetGroupID, projectID); err != nil {
		groupID := uint(0)
		if targetGroupID != nil {
			groupID = *targetGroupID
		}
		s.logger.ErrorContext(ctx, "Failed to move resources to group",
			logger.Uint("target_group_id", groupID),
			logger.Int("resource_count", len(req.ResourceCodes)),
			logger.ErrorField(err))
		return err
	}

	groupID := uint(0)
	if targetGroupID != nil {
		groupID = *targetGroupID
	}
	s.logger.InfoContext(ctx, "Resources moved to group successfully",
		logger.Uint("target_group_id", groupID),
		logger.Int("resource_count", len(req.ResourceCodes)))

	return nil
}

// GetResourceStatistics gets resource statistics by project or resource group
func (s *resourceGroupService) GetResourceStatistics(
	ctx context.Context,
	req *request.GetResourceStatisticsRequest,
) (*response.ResourceStatisticsResponse, error) {
	// Validate: project_id and resource_group_id are mutually exclusive
	if req.ProjectID != nil && req.ResourceGroupID != nil {
		return nil, errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Get statistics
	stats, err := s.repo.GetResourceStatistics(ctx, req.ProjectID, req.ResourceGroupID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get resource statistics",
			logger.ErrorField(err))
		return nil, err
	}

	// Calculate total
	total := int64(0)
	for _, count := range stats {
		total += int64(count)
	}

	return &response.ResourceStatisticsResponse{
		Total:         total,
		ResourceTypes: stats,
	}, nil
}
