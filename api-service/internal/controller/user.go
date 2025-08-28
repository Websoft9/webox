package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/internal/middleware"
	"api-service/pkg/errors"
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
}

// NewUserController create new user controller
func NewUserController(userService service.UserService, logger logger.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// bindAndValidateRequest bind and validate request parameters
func (c *UserController) bindAndValidateRequest(ctx *gin.Context, req interface{}, action string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.logger.WarnContext(ctx, action+" request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, middleware.T(ctx, "common.validation_failed")))
		return false
	}
	return true
}

// handleUserAuth handle user authentication related requests
func (c *UserController) handleUserAuth(
	ctx *gin.Context,
	req interface{},
	action string,
	serviceFunc func(context.Context, interface{}) (interface{}, error),
	successMessageKey string,
) {
	if !c.bindAndValidateRequest(ctx, req, action) {
		return
	}

	result, err := serviceFunc(ctx, req)
	if err != nil {
		if action == "login" {
			c.logger.WarnContext(ctx, "User "+action+" failed", logger.ErrorField(err))
		} else {
			c.logger.ErrorContext(ctx, "User "+action+" failed", logger.ErrorField(err))
		}
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User "+action+" successful")
	pkg_response.Success(ctx, middleware.T(ctx, successMessageKey), result)
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
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, middleware.T(ctx, "common.invalid_request")))
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
	pkg_response.Success(ctx, middleware.T(ctx, successMessageKey), nil)
}

// Register user registration
func (c *UserController) Register(ctx *gin.Context) {
	var req request.UserRegisterRequest
	c.handleUserAuth(ctx, &req, "registration", func(ctx context.Context, r interface{}) (interface{}, error) {
		return c.userService.Register(ctx, r.(*request.UserRegisterRequest))
	}, "user.register_success")
}

// Login user login
func (c *UserController) Login(ctx *gin.Context) {
	var req request.UserLoginRequest
	c.handleUserAuth(ctx, &req, "login", func(ctx context.Context, r interface{}) (interface{}, error) {
		return c.userService.Login(ctx, r.(*request.UserLoginRequest))
	}, "user.login_success")
}

// ChangePassword change user password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := c.getCurrentUserID(ctx)
	if userID == 0 {
		errors.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req request.UserChangePasswordRequest
	if !c.bindAndValidateRequest(ctx, &req, "change password") {
		return
	}

	err := c.userService.ChangePassword(ctx, userID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to change user password", logger.Uint("user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("user_id", userID))
	pkg_response.Success(ctx, middleware.T(ctx, "user.password_change_success"), nil)
}

// ListUsers get user list (admin function)
func (c *UserController) ListUsers(ctx *gin.Context) {
	var req request.UserListRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "List users request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, middleware.T(ctx, "common.validation_failed")))
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

	pkg_response.Success(ctx, middleware.T(ctx, "user.list_get_success"), result)
}

// CreateUser create user (admin function)
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req request.UserCreateRequest
	if !c.bindAndValidateRequest(ctx, &req, "create user") {
		return
	}

	result, err := c.userService.CreateUser(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create user", logger.String("username", req.Username), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User created successfully", logger.String("username", req.Username),
		logger.Uint("user_id", result.ID))
	pkg_response.Success(ctx, middleware.T(ctx, "user.created_success"), result)
}

// GetUser get user details (admin function)
func (c *UserController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, middleware.T(ctx, "common.invalid_request")))
		return
	}

	user, err := c.userService.GetUser(ctx, uint(userID))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user details", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User details retrieved successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), user)
}

// UpdateUser update user (admin function)
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, middleware.T(ctx, "common.invalid_request")))
		return
	}

	var req request.UserUpdateRequest
	if !c.bindAndValidateRequest(ctx, &req, "update user") {
		return
	}

	result, err := c.userService.UpdateUser(ctx, uint(userID), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User updated successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, middleware.T(ctx, "user.updated_success"), result)
}

// DeleteUser delete user (admin function)
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid user ID parameter", logger.String("user_id", userIDStr))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, middleware.T(ctx, "common.invalid_request")))
		return
	}

	err = c.userService.DeleteUser(ctx, uint(userID))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete user", logger.Uint("user_id", uint(userID)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User deleted successfully", logger.Uint("user_id", uint(userID)))
	pkg_response.Success(ctx, middleware.T(ctx, "user.deleted_success"), nil)
}

// UpdateUserStatus update user status (admin function)
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	var req request.UserUpdateStatusRequest
	c.handleUserIDBasedRequest(ctx, &req, "update user status",
		func(ctx context.Context, userID uint, r interface{}) error {
			return c.userService.UpdateUserStatus(ctx, userID, r.(*request.UserUpdateStatusRequest))
		}, "user.status_update_success")
}

// UpdateUserPassword update user password (admin function)
func (c *UserController) UpdateUserPassword(ctx *gin.Context) {
	var req request.UserPasswordUpdateRequest
	c.handleUserIDBasedRequest(ctx, &req, "update user password",
		func(ctx context.Context, userID uint, r interface{}) error {
			return c.userService.UpdateUserPassword(ctx, userID, r.(*request.UserPasswordUpdateRequest))
		}, "user.password_update_success")
}

// getCurrentUserID get current user ID from context
func (c *UserController) getCurrentUserID(ctx *gin.Context) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.WarnContext(ctx, "User ID not found in context")
		return 0
	}

	id, ok := userID.(uint)
	if !ok {
		c.logger.WarnContext(ctx, "Invalid user ID type in context", logger.Any("user_id", userID))
		return 0
	}

	return id
}
