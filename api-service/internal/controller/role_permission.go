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

// RolePermissionController handles role & permission management operations
type RolePermissionController struct {
	roleService       service.RoleService
	permissionService service.PermissionService
	validator         *validator.Validate
	logger            logger.Logger
	i18n              *i18n.I18n
}

// NewRolePermissionController creates a new role & permission controller instance
func NewRolePermissionController(
	roleService service.RoleService,
	permissionService service.PermissionService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *RolePermissionController {
	return &RolePermissionController{
		roleService:       roleService,
		permissionService: permissionService,
		validator:         validator,
		logger:            logger,
		i18n:              i18n,
	}
}

// CreateRole creates a new role
// @Summary Create role
// @Description Create a new role
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreateRoleRequest true "Create role request"
// @Success 201 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/roles [post]
func (c *RolePermissionController) CreateRole(ctx *gin.Context) {
	var req request.CreateRoleRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	// Get current user ID
	userID, ok := GetUserID(ctx, c.i18n)
	if !ok {
		return
	}

	// Create role
	role, err := c.roleService.CreateRole(ctx.Request.Context(), &req, userID)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}

	ResponseOKWithData(ctx, role, "role.created_success", c.i18n)
}

// GetRole gets role details
// @Summary Get role details
// @Description Get role details by ID
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [get]
func (c *RolePermissionController) GetRole(ctx *gin.Context) {
	// Get role ID
	id, ok := ParseIDParam(ctx, "id", "validation.invalid_role_id", c.i18n)
	if !ok {
		return
	}

	// Get role
	role, err := c.roleService.GetRole(ctx.Request.Context(), id)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}

	ResponseOKWithData(ctx, role, "common.success", c.i18n)
}

// ListRoles gets role list
// @Summary Get role list
// @Description Get paginated role list
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param status query int false "Role status" Enums(-1, 0, 1)
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/roles [get]
func (c *RolePermissionController) ListRoles(ctx *gin.Context) {
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, roles, "common.success", c.i18n)
}

// UpdateRole updates role
// @Summary Update role
// @Description Update role information
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.UpdateRoleRequest true "Update role request"
// @Success 200 {object} response.APIResponse{data=response.RoleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [put]
func (c *RolePermissionController) UpdateRole(ctx *gin.Context) {
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
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(bindErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   bindErr.Error(),
		})
		return
	}

	// Validate request parameters
	if validateErr := c.validator.Struct(&req); validateErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(validateErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   validateErr.Error(),
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}

	ResponseOKWithData(ctx, role, "role.update_success", c.i18n)
}

// DeleteRole deletes role
// @Summary Delete role
// @Description Delete specified role
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id} [delete]
func (c *RolePermissionController) DeleteRole(ctx *gin.Context) {
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOK(ctx, "role.delete_success", c.i18n)
}

// AssignPermissions assigns permissions to role
// @Summary Assign permissions to role
// @Description Assign permissions to specified role
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.RolePermissionRequest true "Role permission request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/permissions [post]
func (c *RolePermissionController) AssignPermissions(ctx *gin.Context) {
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
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(bindErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   bindErr.Error(),
		})
		return
	}

	// Validate request parameters
	if validateErr := c.validator.Struct(&req); validateErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(validateErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   validateErr.Error(),
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOK(ctx, "role.assign_permissions_success", c.i18n)
}

// RemovePermissions removes permissions from role
// @Summary Remove permissions from role
// @Description Remove permissions from specified role
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body request.RolePermissionRequest true "Role permission request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/permissions [delete]
func (c *RolePermissionController) RemovePermissions(ctx *gin.Context) {
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
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(bindErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   bindErr.Error(),
		})
		return
	}

	// Validate request parameters
	if validateErr := c.validator.Struct(&req); validateErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(validateErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   validateErr.Error(),
		})
		return
	}

	// Remove permissions
	err = c.roleService.RemovePermissions(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOK(ctx, "role.remove_permissions_success", c.i18n)
}

// GetRoleUsers gets users associated with role
// @Summary Get users associated with role
// @Description Get list of users associated with specified role
// @Tags Role Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/roles/{id}/users [get]
func (c *RolePermissionController) GetRoleUsers(ctx *gin.Context) {
	// Get role ID
	id, ok := ParseIDParam(ctx, "id", "validation.invalid_role_id", c.i18n)
	if !ok {
		return
	}

	// Get pagination parameters
	page, pageSize := GetPaginationParams(ctx)

	// Get role users
	users, err := c.roleService.GetRoleUsers(ctx.Request.Context(), id, page, pageSize)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, users, "common.success", c.i18n)
}

// CreatePermission creates a new permission
// @Summary Create permission
// @Description Create a new permission
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreatePermissionRequest true "Create permission request"
// @Success 201 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/permissions [post]
func (c *RolePermissionController) CreatePermission(ctx *gin.Context) {
	var req request.CreatePermissionRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	// Get current user ID
	userID, ok := GetUserID(ctx, c.i18n)
	if !ok {
		return
	}

	// Create permission
	permission, err := c.permissionService.CreatePermission(ctx.Request.Context(), &req, userID)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, permission, "permission.created_success", c.i18n)
}

// GetPermission gets permission details
// @Summary Get permission details
// @Description Get permission details by ID
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [get]
func (c *RolePermissionController) GetPermission(ctx *gin.Context) {
	// Get permission ID
	id, ok := ParseIDParam(ctx, "id", "validation.invalid_permission_id", c.i18n)
	if !ok {
		return
	}

	// Get permission
	permission, err := c.permissionService.GetPermission(ctx.Request.Context(), id)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, permission, "common.success", c.i18n)
}

// UpdatePermission updates permission
// @Summary Update permission
// @Description Update permission information
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param request body request.UpdatePermissionRequest true "Update permission request"
// @Success 200 {object} response.APIResponse{data=response.PermissionResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [put]
func (c *RolePermissionController) UpdatePermission(ctx *gin.Context) {
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
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(bindErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   bindErr.Error(),
		})
		return
	}

	// Validate request parameters
	if validateErr := c.validator.Struct(&req); validateErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(validateErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   validateErr.Error(),
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, permission, "permission.update_success", c.i18n)
}

// DeletePermission deletes permission
// @Summary Delete permission
// @Description Delete specified permission
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id} [delete]
func (c *RolePermissionController) DeletePermission(ctx *gin.Context) {
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOK(ctx, "permission.delete_success", c.i18n)
}

// ListPermissions gets permission list
// @Summary Get permission list
// @Description Get paginated permission list
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param module query string false "Permission module"
// @Param scope query string false "Permission scope" Enums(platform, project)
// @Param status query int false "Permission status" Enums(-1, 0, 1)
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.PermissionListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/permissions [get]
func (c *RolePermissionController) ListPermissions(ctx *gin.Context) {
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, permissions, "common.success", c.i18n)
}

// GetPermissionTree gets permission tree
// @Summary Get permission tree
// @Description Get permission tree structure
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param scope query string false "Permission scope" Enums(platform, project)
// @Param status query int false "Permission status" Enums(-1, 0, 1)
// @Success 200 {object} response.APIResponse{data=response.PermissionTreeResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/permissions/tree [get]
func (c *RolePermissionController) GetPermissionTree(ctx *gin.Context) {
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
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, tree, "common.success", c.i18n)
}

// GetPermissionRoles gets roles associated with permission
// @Summary Get roles associated with permission
// @Description Get list of roles associated with specified permission
// @Tags Permission Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.APIResponse{data=response.RoleListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/permissions/{id}/roles [get]
func (c *RolePermissionController) GetPermissionRoles(ctx *gin.Context) {
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
		if p, parseErr := strconv.Atoi(pageStr); parseErr == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, parseErr := strconv.Atoi(pageSizeStr); parseErr == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Get permission roles
	roles, err := c.permissionService.GetPermissionRoles(ctx.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		ResponseWithError(ctx, err, c.logger, c.i18n)
		return
	}
	ResponseOKWithData(ctx, roles, "common.success", c.i18n)
}
