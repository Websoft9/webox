package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TagController handles tag management operations
type TagController struct {
	tagService service.TagService
	validator  *validator.Validate
	logger     logger.Logger
	i18n       *i18n.I18n
}

// NewTagController creates a new tag controller
func NewTagController(tagService service.TagService, validator *validator.Validate, logger logger.Logger, i18n *i18n.I18n) *TagController {
	return &TagController{
		tagService: tagService,
		validator:  validator,
		logger:     logger,
		i18n:       i18n,
	}
}

// CreateTag creates a new tag
// @Summary Create a new tag
// @Description Create a new tag
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TagCreateRequest true "Tag create request"
// @Success 201 {object} response.APIResponse{data=response.TagResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags [post]
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req request.TagCreateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	tag, err := c.tagService.CreateTag(ctx, &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create tag", logger.String("tag_name", req.Name), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tag created successfully", logger.Uint("tag_id", uint(tag.ID)))
	response.SuccessWithData(ctx, tag)
}

// GetTag retrieves a tag by ID
// @Summary Get tag by ID
// @Description Retrieve a specific tag by its ID
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} response.APIResponse{data=response.TagResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/{id} [get]
func (c *TagController) GetTag(ctx *gin.Context) {
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	tag, err := c.tagService.GetTag(ctx, id)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get tag", logger.Uint("tag_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, tag)
}

// UpdateTag updates an existing tag
// @Summary Update tag
// @Description Update an existing tag
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param request body request.TagUpdateRequest true "Tag update request"
// @Success 200 {object} response.APIResponse{data=response.TagResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/{id} [put]
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}
	var req request.TagUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	tag, err := c.tagService.UpdateTag(ctx, id, &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update tag", logger.Uint("tag_id", uint(id)), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tag updated successfully", logger.Uint("tag_id", uint(id)))
	response.SuccessWithData(ctx, tag)
}

// DeleteTag deletes a tag
// @Summary Delete tag
// @Description Delete a tag by ID
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/{id} [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	err := c.tagService.DeleteTag(ctx, id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete tag", logger.Uint("tag_id", uint(id)), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tag deleted successfully", logger.Uint("tag_id", uint(id)))
	response.Success(ctx)
}

// ListTags lists tags with optional filtering
// @Summary List tags
// @Description List tags with optional search and exclusion filters
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search term"
// @Param excludeIds query string false "Comma-separated list of tag IDs to exclude"
// @Success 200 {object} response.APIResponse{data=[]response.TagResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags [get]
func (c *TagController) ListTags(ctx *gin.Context) {
	var req request.TagListRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	tags, err := c.tagService.ListTags(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to list tags", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, tags)
}

// AssignTags assigns tags to resources
// @Summary Assign tags to resources
// @Description Assign multiple tags to multiple resources
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TagAssignRequest true "Tag assignment request"
// @Success 200 {object} response.APIResponse{data=response.TagAssignResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/assign [post]
func (c *TagController) AssignTags(ctx *gin.Context) {
	var req request.TagAssignRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	result, err := c.tagService.AssignTags(ctx, &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to assign tags", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tags assigned successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagNames)))
	response.SuccessWithData(ctx, result)
}

// ReplaceTags replaces all tags for resources
// @Summary Replace tags for resources
// @Description Replace all tags for multiple resources
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TagAssignRequest true "Tag replacement request"
// @Success 200 {object} response.APIResponse{data=response.TagAssignResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/replace [post]
func (c *TagController) ReplaceTags(ctx *gin.Context) {
	var req request.TagAssignRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	result, err := c.tagService.ReplaceTags(ctx, &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to replace tags", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tags replaced successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagNames)))
	response.SuccessWithData(ctx, result)
}

// UnassignTags removes tags from resources
// @Summary Unassign tags from resources
// @Description Remove multiple tags from multiple resources
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TagUnassignRequest true "Tag unassignment request"
// @Success 200 {object} response.APIResponse{data=response.TagAssignResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/unassign [post]
func (c *TagController) UnassignTags(ctx *gin.Context) {
	var req request.TagUnassignRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	result, err := c.tagService.UnassignTags(ctx, &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to unassign tags", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tags unassigned successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagIDs)))
	response.SuccessWithData(ctx, result)
}

// SearchTags searches for tags
// @Summary Search tags
// @Description Search for tags by name or description
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {object} response.APIResponse{data=[]response.TagResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/search [get]
func (c *TagController) SearchTags(ctx *gin.Context) {
	var req request.TagNameSearchRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	tags, err := c.tagService.SearchTags(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to search tags", logger.String("query", req.Q), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Tags searched successfully",
		logger.String("query", req.Q),
		logger.Int("result_count", len(tags)))
	response.SuccessWithData(ctx, tags)
}

// GetResourceTags gets tags for a specific resource
// @Summary Get resource tags
// @Description Get all tags assigned to a specific resource
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param resourceId query int true "Resource ID"
// @Success 200 {object} response.APIResponse{data=[]response.TagSimpleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/taggings [get]
func (c *TagController) GetResourceTags(ctx *gin.Context) {
	var req request.TaggingListRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	tags, err := c.tagService.GetResourceTags(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get resource tags", logger.Uint("resource_id", uint(req.ResourceID)), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Resource tags retrieved successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(tags)))
	response.SuccessWithData(ctx, tags)
}

// SearchResourcesByTags searches for resources by tags
// @Summary Search resources by tags
// @Description Search for resources that have specific tags
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param tagIds query string false "Comma-separated tag IDs"
// @Param tagNames query string false "Comma-separated tag names"
// @Param operation query string false "AND or OR operation" Enums(AND, OR)
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} response.APIResponse{data=response.TagSearchResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/tags/taggings/search [get]
func (c *TagController) SearchResourcesByTags(ctx *gin.Context) {
	var req request.TagSearchRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.tagService.SearchResourcesByTags(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to search resources by tags", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Resources searched by tags successfully",
		logger.Int("tag_id_count", len(req.TagIDs)),
		logger.Int("tag_name_count", len(req.TagNames)),
		logger.String("operation", req.Operation))
	response.SuccessWithData(ctx, result)
}
