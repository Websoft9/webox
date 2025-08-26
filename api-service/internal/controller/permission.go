package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PermissionController permission controller
type PermissionController struct {
	permissionService service.PermissionService
	validator         *validator.Validate
	logger            logger.Logger
	i18n              *i18n.I18n
}

// NewPermissionController creates a new permission controller instance
func NewPermissionController(
	permissionService service.PermissionService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *PermissionController {
	return &PermissionController{
		permissionService: permissionService,
		validator:         validator,
		logger:            logger,
		i18n:              i18n,
	}
}

// CreatePermission creates a new permission
// @Summary Create permission
// @Description Create a new permission
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param request body request.CreatePermissionRequest true "Create permission request"
// @Success 201 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/permissions [post]
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req request.CreatePermissionRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Create permission
	permission, err := c.permissionService.CreatePermission(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to create permission", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.create_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": c.i18n.T(ctx, "permission.create_success"),
		"data":    permission,
	})
}

// GetPermission gets permission details
// @Summary Get permission details
// @Description Get permission details by ID
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [get]
func (c *PermissionController) GetPermission(ctx *gin.Context) {
	// Get permission ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_permission_id"),
		})
		return
	}

	// Get permission
	permission, err := c.permissionService.GetPermission(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "permission not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "permission.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get permission", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.get_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    permission,
	})
}

// UpdatePermission updates permission
// @Summary Update permission
// @Description Update permission information
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param request body request.UpdatePermissionRequest true "Update permission request"
// @Success 200 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [put]
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	// Get permission ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_permission_id"),
		})
		return
	}

	var req request.UpdatePermissionRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Update permission
	permission, err := c.permissionService.UpdatePermission(ctx.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == "permission not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "permission.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to update permission", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.update_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "permission.update_success"),
		"data":    permission,
	})
}

// DeletePermission deletes permission
// @Summary Delete permission
// @Description Delete specified permission
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [delete]
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	// Get permission ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_permission_id"),
		})
		return
	}

	// Delete permission
	err = c.permissionService.DeletePermission(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "permission not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "permission.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to delete permission", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.delete_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "permission.delete_success"),
	})
}

// ListPermissions gets permission list
// @Summary Get permission list
// @Description Get paginated permission list
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param module query string false "Permission module"
// @Param scope query string false "Permission scope" Enums(platform, project)
// @Param status query int false "Permission status" Enums(0, 1)
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.PermissionListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/permissions [get]
func (c *PermissionController) ListPermissions(ctx *gin.Context) {
	var req request.ListPermissionsRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_query_parameters"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Query validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.query_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get permission list
	permissions, err := c.permissionService.ListPermissions(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to list permissions", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.list_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    permissions,
	})
}

// GetPermissionTree gets permission tree
// @Summary Get permission tree
// @Description Get permission tree structure
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param scope query string false "Permission scope" Enums(platform, project)
// @Param status query int false "Permission status" Enums(0, 1)
// @Success 200 {object} response.APIResponse{data=response.PermissionTreeResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/permissions/tree [get]
func (c *PermissionController) GetPermissionTree(ctx *gin.Context) {
	var req request.PermissionTreeRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_query_parameters"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Query validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.query_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get permission tree
	tree, err := c.permissionService.GetPermissionTree(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get permission tree", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.tree_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    tree,
	})
}

// GetPermissionRoles gets roles associated with permission
// @Summary Get roles associated with permission
// @Description Get list of roles associated with specified permission
// @Tags Permission Management
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id}/roles [get]
func (c *PermissionController) GetPermissionRoles(ctx *gin.Context) {
	// Get permission ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_permission_id"),
		})
		return
	}

	// Get pagination parameters
	page := 1
	pageSize := 20

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Get permission roles
	roles, err := c.permissionService.GetPermissionRoles(ctx.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get permission roles", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "permission.roles_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    roles,
	})
}
