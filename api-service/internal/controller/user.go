package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController user controller
type UserController struct {
	userService service.UserService
	logger      logger.Logger
	i18n        *i18n.I18n
}

// NewUserController create new user controller
func NewUserController(userService service.UserService, logger logger.Logger, i18n *i18n.I18n) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
		i18n:        i18n,
	}
}

// bindAndValidateRequest bind and validate request parameters
func (c *UserController) bindAndValidateRequest(ctx *gin.Context, req interface{}, action string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.logger.WarnContext(ctx, action+" request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return false
	}
	return true
}

// handleUserIDBasedRequest handle requests that need user ID from URL parameter
func (c *UserController) handleUserIDBasedRequest(
	ctx *gin.Context,
	req interface{},
	action string,
	serviceFunc func(context.Context, uint, interface{}) error,
	successMessageKey string,
) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return
	}

	if !c.bindAndValidateRequest(ctx, req, action) {
		return
	}

	err = serviceFunc(ctx, uint(userID), req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to "+action, logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User "+action+" successful", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, c.i18n.T(ctx, successMessageKey), nil)
}

// ListUsers get user list (admin function)
// @Summary List users
// @Description Get paginated list of users (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query int false "User status filter" Enums(0,1)
// @Param keyword query string false "Search keyword"
// @Param gender query int false "Gender filter" Enums(0,1,2)
// @Param language query string false "Language filter"
// @Success 200 {object} response.APIResponse{data=response.UserListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	var req request.UserListRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "List users request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, total, err := c.userService.ListUsers(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user list", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User list retrieved successfully", logger.Int64("total", total))

	// Return paginated response
	result := map[string]interface{}{
		"users":       users.Users,
		"total":       users.Total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": (int(total) + req.PageSize - 1) / req.PageSize,
	}

	pkg_response.Success(ctx, c.i18n.T(ctx, "user.list_get_success"), result)
}

// CreateUser create user (admin function)
// @Summary Create user
// @Description Create a new user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UserCreateRequest true "User creation request"
// @Success 200 {object} response.APIResponse{data=response.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}
	var req request.UserCreateRequest
	if !c.bindAndValidateRequest(ctx, &req, "create user") {
		return
	}

	result, err := c.userService.CreateUser(ctx, currentUserID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create user", logger.String("username", req.Username), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User created successfully", logger.String("username", req.Username),
		logger.Uint("user_id", result.ID))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.created_success"), result)
}

// GetUser get user details (admin function)
// @Summary Get user details
// @Description Get detailed information of a specific user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse{data=response.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return
	}

	user, err := c.userService.GetUser(ctx, uint(userID))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user details", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User details retrieved successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "common.success"), user)
}

// UpdateUser update user (admin function)
// @Summary Update user
// @Description Update a user's information (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body request.UserUpdateRequest true "User update request"
// @Success 200 {object} response.APIResponse{data=response.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return
	}

	var req request.UserUpdateRequest
	if !c.bindAndValidateRequest(ctx, &req, "update user") {
		return
	}

	result, err := c.userService.UpdateUser(ctx, currentUserID.(uint), uint(userID), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User updated successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.updated_success"), result)
}

// DeleteUser delete user (admin function)
// @Summary Delete user
// @Description Delete a user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.ErrValidationFailed)
		return
	}

	err = c.userService.DeleteUser(ctx, uint(userID))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete user", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User deleted successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.deleted_success"), nil)
}

// UpdateUserStatus update user status (admin function)
// @Summary Update user status
// @Description Update user account status (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body request.UserUpdateStatusRequest true "User status update request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id}/status [put]
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	var req request.UserUpdateStatusRequest
	c.handleUserIDBasedRequest(ctx, &req, "update user status",
		func(ctx context.Context, userID uint, r interface{}) error {
			return c.userService.UpdateUserStatus(ctx, userID, r.(*request.UserUpdateStatusRequest))
		}, "user.status_update_success")
}

// UpdateUserPassword update user password (admin function)
// @Summary Update user password
// @Description Update user password (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body request.UserPasswordUpdateRequest true "User password update request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id}/password [put]
func (c *UserController) UpdateUserPassword(ctx *gin.Context) {
	var req request.UserPasswordUpdateRequest
	c.handleUserIDBasedRequest(ctx, &req, "update user password",
		func(ctx context.Context, userID uint, r interface{}) error {
			return c.userService.UpdateUserPassword(ctx, userID, r.(*request.UserPasswordUpdateRequest))
		}, "user.password_update_success")
}
