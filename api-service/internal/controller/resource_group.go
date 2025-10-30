package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// ResourceGroupController handles resource group HTTP requests
type ResourceGroupController struct {
	service   interfaceService.ResourceGroupService
	validator *validator.Validate
	logger    logger.Logger
}

// NewResourceGroupController creates a new resource group controller
func NewResourceGroupController(
	service interfaceService.ResourceGroupService,
	validator *validator.Validate,
	logger logger.Logger,
) *ResourceGroupController {
	return &ResourceGroupController{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// CreateResourceGroup create a resource group
// @Summary Create resource group
// @Description Create a new resource group for organizing resources within a project
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateResourceGroupRequest true "Resource group configuration"
// @Success 201 {object} common.APIResponse{data=response.ResourceGroupResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups [post]
func (ctrl *ResourceGroupController) CreateResourceGroup(c *gin.Context) {
	var req request.CreateResourceGroupRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context (set by auth middleware)
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	result, err := ctrl.service.CreateResourceGroup(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Resource group created successfully",
		logger.String("name", req.Name),
		logger.Uint("resource_group_id", result.ID))

	response.SuccessWithData(c, result)
}

// GetResourceGroup get resource group by ID
// @Summary Get resource group by ID
// @Description Get detailed information about a specific resource group
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource Group ID"
// @Success 200 {object} common.APIResponse{data=response.ResourceGroupDetailResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups/{id} [get]
func (ctrl *ResourceGroupController) GetResourceGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	// Call service
	result, err := ctrl.service.GetResourceGroup(c.Request.Context(), uint(id))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// GetResourceGroupList get resource group list
// @Summary Get resource group list
// @Description Get a paginated list of resource groups with optional filters. Results are sorted by sort_order (ascending) and then by created_at (descending).
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keyword query string false "Search keyword (name, code, or description)"
// @Param project_id query int false "Filter by project ID"
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups [get]
func (ctrl *ResourceGroupController) GetResourceGroupList(c *gin.Context) {
	var req request.GetResourceGroupListRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	result, err := ctrl.service.GetResourceGroupList(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// UpdateResourceGroup update resource group
// @Summary Update resource group
// @Description Update an existing resource group's information
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource Group ID"
// @Param request body request.UpdateResourceGroupRequest true "Resource group update data"
// @Success 200 {object} common.APIResponse{data=response.ResourceGroupResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups/{id} [put]
func (ctrl *ResourceGroupController) UpdateResourceGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	var req request.UpdateResourceGroupRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	result, err := ctrl.service.UpdateResourceGroup(c.Request.Context(), uint(id), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Resource group updated successfully",
		logger.Uint("resource_group_id", uint(id)))

	response.SuccessWithData(c, result)
}

// DeleteResourceGroup delete resource group
// @Summary Delete resource group
// @Description Delete a resource group by ID
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource Group ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups/{id} [delete]
func (ctrl *ResourceGroupController) DeleteResourceGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	// Call service
	if err := ctrl.service.DeleteResourceGroup(c.Request.Context(), uint(id)); err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Resource group deleted successfully",
		logger.Uint("resource_group_id", uint(id)))

	response.Success(c)
}

// GetResourcesByGroupID get resources in a resource group
// @Summary Get resources in a resource group
// @Description Get all resources associated with a specific resource group, supports filtering by resource type
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource Group ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param resource_type query string false "Filter by resource type (server, database, secret, etc.)"
// @Success 200 {object} common.APIResponse{data=response.ResourceGroupResourcesResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resource-groups/{id}/resources [get]
func (ctrl *ResourceGroupController) GetResourcesByGroupID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	var req request.GetResourceGroupResourcesRequest
	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.GetResourcesByGroupID(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// MoveResourcesToGroup move multiple resources to a specific resource group or default group
// @Summary Move resources to a target group
// @Description Move multiple resources to a specified resource group.
// @Description The resource_group_id field is required but can be null to move resources to the default group.
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.MoveResourcesToGroupRequest true "Resource codes and target group ID (required field, null for default group)"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resources/resource-group [put]
func (ctrl *ResourceGroupController) MoveResourcesToGroup(c *gin.Context) {
	var req request.MoveResourcesToGroupRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	if err := ctrl.service.MoveResourcesToGroup(c.Request.Context(), &req, ownerID.(uint)); err != nil {
		response.WithError(c, err)
		return
	}

	groupID := uint(0)
	if req.ResourceGroupID != nil {
		groupID = *req.ResourceGroupID
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Resources moved to group successfully",
		logger.Uint("target_group_id", groupID),
		logger.Int("resource_count", len(req.ResourceCodes)))

	response.Success(c)
}

// GetResourceStatistics get resource statistics
// @Summary Get resource statistics
// @Description Get resource statistics: 1) No params - returns all projects grouped statistics
// @Description 2) By project ID 3) By resource group ID. Note: project_id and resource_group_id are mutually exclusive
// @Tags Resource Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id query int false "Filter by project ID (mutually exclusive with resource_group_id)"
// @Param resource_group_id query int false "Filter by resource group ID (mutually exclusive with project_id)"
// @Success 200 {object} common.APIResponse{data=response.ResourceStatisticsResponse} "Returns statistics with scope indicator (all_projects/project/resource_group)"
// @Failure 400 {object} common.APIResponse "Bad request - parameters are mutually exclusive"
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/resources/statistics [get]
func (ctrl *ResourceGroupController) GetResourceStatistics(c *gin.Context) {
	var req request.GetResourceStatisticsRequest
	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Validate mutual exclusivity
	if err := req.Validate(); err != nil {
		response.WithErrorAndCode(c, errors.CodeInvalidParameterFormat, err)
		return
	}

	// Call service
	result, err := ctrl.service.GetResourceStatistics(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}
