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

// RoleController handles role management operations
type RoleController struct {
	roleService service.RoleService
	validator   *validator.Validate
	logger      logger.Logger
	i18n        *i18n.I18n
}

// NewRoleController creates a new role controller instance
func NewRoleController(
	roleService service.RoleService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *RoleController {
	return &RoleController{
		roleService: roleService,
		validator:   validator,
		logger:      logger,
		i18n:        i18n,
	}
}

// CreateRole creates a new role
// @Summary Create role
// @Description Create a new role
// @Tags Role Management
// @Accept json
// @Produce json
// @Param request body request.CreateRoleRequest true "Create role request"
// @Success 201 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req request.CreateRoleRequest

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

	// Create role
	role, err := c.roleService.CreateRole(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to create role", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.create_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": c.i18n.T(ctx, "role.create_success"),
		"data":    role,
	})
}

// GetRole gets role details
// @Summary Get role details
// @Description Get role details by ID
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [get]
func (c *RoleController) GetRole(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
		})
		return
	}

	// Get role
	role, err := c.roleService.GetRole(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "role.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get role", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.get_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    role,
	})
}

// ListRoles gets role list
// @Summary Get role list
// @Description Get paginated role list
// @Tags Role Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param status query int false "Role status" Enums(0, 1)
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/roles [get]
func (c *RoleController) ListRoles(ctx *gin.Context) {
	var req request.ListRolesRequest

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

	// Get role list
	roles, err := c.roleService.ListRoles(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to list roles", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.list_failed"),
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

// UpdateRole updates role
// @Summary Update role
// @Description Update role information
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.UpdateRoleRequest true "Update role request"
// @Success 200 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
		})
		return
	}

	var req request.UpdateRoleRequest

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

	// Update role
	role, err := c.roleService.UpdateRole(ctx.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "role.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to update role", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.update_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "role.update_success"),
		"data":    role,
	})
}

// DeleteRole deletes role
// @Summary Delete role
// @Description Delete specified role
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
		})
		return
	}

	// Delete role
	err = c.roleService.DeleteRole(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "role.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to delete role", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.delete_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "role.delete_success"),
	})
}

// AssignPermissions assigns permissions to role
// @Summary Assign permissions to role
// @Description Assign permissions to specified role
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.RolePermissionRequest true "Role permission request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/permissions [post]
func (c *RoleController) AssignPermissions(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
		})
		return
	}

	var req request.RolePermissionRequest

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

	// Assign permissions
	err = c.roleService.AssignPermissions(ctx.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to assign permissions", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.assign_permissions_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "role.assign_permissions_success"),
	})
}

// RemovePermissions removes permissions from role
// @Summary Remove permissions from role
// @Description Remove permissions from specified role
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.RolePermissionRequest true "Role permission request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/permissions [delete]
func (c *RoleController) RemovePermissions(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
		})
		return
	}

	var req request.RolePermissionRequest

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

	// Remove permissions
	err = c.roleService.RemovePermissions(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to remove permissions", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.remove_permissions_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "role.remove_permissions_success"),
	})
}

// GetRoleUsers gets users associated with role
// @Summary Get users associated with role
// @Description Get list of users associated with specified role
// @Tags Role Management
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/users [get]
func (c *RoleController) GetRoleUsers(ctx *gin.Context) {
	// Get role ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_role_id"),
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

	// Get role users
	users, err := c.roleService.GetRoleUsers(ctx.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get role users", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "role.users_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    users,
	})
}
