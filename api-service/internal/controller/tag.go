package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"
	"strconv"

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

// getCurrentUserID64 gets current user ID as uint64
func (tc *TagController) getCurrentUserID64(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return uint64(id)
	}
	if id, ok := userID.(uint64); ok {
		return id
	}
	return 0
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
func (tc *TagController) CreateTag(c *gin.Context) {
	var req request.TagCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		tc.logger.WarnContext(c, "Create tag request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	tag, err := tc.tagService.CreateTag(c, &req, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to create tag", logger.String("tag_name", req.Name), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tag created successfully", logger.Uint("tag_id", uint(tag.ID)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.create.success"), tag)
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
func (tc *TagController) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		tc.logger.WarnContext(c, "Invalid tag ID parameter", logger.String("id", idStr))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	tag, err := tc.tagService.GetTag(c, id)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to get tag", logger.Uint("tag_id", uint(id)), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	pkg_response.Success(c, tc.i18n.T(c, "tag.get.success"), tag)
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
func (tc *TagController) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		tc.logger.WarnContext(c, "Invalid tag ID parameter", logger.String("id", idStr))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	var req request.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tc.logger.WarnContext(c, "Update tag request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	tag, err := tc.tagService.UpdateTag(c, id, &req, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to update tag", logger.Uint("tag_id", uint(id)), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tag updated successfully", logger.Uint("tag_id", uint(id)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.update.success"), tag)
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
func (tc *TagController) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		tc.logger.WarnContext(c, "Invalid tag ID parameter", logger.String("id", idStr))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	err = tc.tagService.DeleteTag(c, id, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to delete tag", logger.Uint("tag_id", uint(id)), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tag deleted successfully", logger.Uint("tag_id", uint(id)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.delete.success"), nil)
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
func (tc *TagController) ListTags(c *gin.Context) {
	var req request.TagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.logger.WarnContext(c, "List tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	tags, err := tc.tagService.ListTags(c, &req)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to list tags", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	pkg_response.Success(c, tc.i18n.T(c, "tag.list.success"), tags)
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
func (tc *TagController) AssignTags(c *gin.Context) {
	var req request.TagAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tc.logger.WarnContext(c, "Assign tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	result, err := tc.tagService.AssignTags(c, &req, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to assign tags", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tags assigned successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagNames)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.assign.success"), result)
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
func (tc *TagController) ReplaceTags(c *gin.Context) {
	var req request.TagAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tc.logger.WarnContext(c, "Replace tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	result, err := tc.tagService.ReplaceTags(c, &req, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to replace tags", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tags replaced successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagNames)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.replace.success"), result)
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
func (tc *TagController) UnassignTags(c *gin.Context) {
	var req request.TagUnassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tc.logger.WarnContext(c, "Unassign tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	userID := tc.getCurrentUserID64(c)
	if userID == 0 {
		errors.HandleError(c, errors.ErrInvalidToken)
		return
	}

	result, err := tc.tagService.UnassignTags(c, &req, userID)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to unassign tags", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tags unassigned successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(req.TagIDs)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.unassign.success"), result)
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
func (tc *TagController) SearchTags(c *gin.Context) {
	var req request.TagNameSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.logger.WarnContext(c, "Search tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	tags, err := tc.tagService.SearchTags(c, &req)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to search tags", logger.String("query", req.Q), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Tags searched successfully",
		logger.String("query", req.Q),
		logger.Int("result_count", len(tags)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.search.success"), tags)
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
func (tc *TagController) GetResourceTags(c *gin.Context) {
	var req request.TaggingListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.logger.WarnContext(c, "Get resource tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	tags, err := tc.tagService.GetResourceTags(c, &req)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to get resource tags", logger.Uint("resource_id", uint(req.ResourceID)), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Resource tags retrieved successfully",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Int("tag_count", len(tags)))
	pkg_response.Success(c, tc.i18n.T(c, "tag.resource.success"), tags)
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
func (tc *TagController) SearchResourcesByTags(c *gin.Context) {
	var req request.TagSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tc.logger.WarnContext(c, "Search resources by tags request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(c, errors.ErrValidationFailed)
		return
	}

	result, err := tc.tagService.SearchResourcesByTags(c, &req)
	if err != nil {
		tc.logger.ErrorContext(c, "Failed to search resources by tags", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	tc.logger.InfoContext(c, "Resources searched by tags successfully",
		logger.Int("tag_id_count", len(req.TagIDs)),
		logger.Int("tag_name_count", len(req.TagNames)),
		logger.String("operation", req.Operation))
	pkg_response.Success(c, tc.i18n.T(c, "tag.search.resources.success"), result)
}
