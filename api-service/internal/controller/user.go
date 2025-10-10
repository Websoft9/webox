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

// UserController user controller
type UserController struct {
	userService service.UserService
	logger      logger.Logger
	i18n        *i18n.I18n
	validator   *validator.Validate
}

// NewUserController create new user controller
func NewUserController(
	userService service.UserService,
	logger logger.Logger,
	i18n *i18n.I18n,
	validator *validator.Validate,
) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
		i18n:        i18n,
		validator:   validator,
	}
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
// @Success 200 {object} common.APIResponse{data=response.UserListResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	var req request.UserListRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.userService.ListUsers(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user list", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User list retrieved successfully")

	response.SuccessWithData(ctx, result)
}

// CreateUser create user (admin function)
// @Summary Create user
// @Description Create a new user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UserCreateRequest true "User creation request"
// @Success 200 {object} common.APIResponse{data=response.UserResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}
	var req request.UserCreateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.userService.CreateUser(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create user", logger.String("username", req.Username), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User created successfully", logger.String("username", req.Username),
		logger.Uint("user_id", result.ID))
	response.SuccessWithData(ctx, result)
}

// GetUser get user details (admin function)
// @Summary Get user details
// @Description Get detailed information of a specific user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} common.APIResponse{data=response.UserResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	// Get ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	user, err := c.userService.GetUser(ctx, id)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user details", logger.Uint("user_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User details retrieved successfully", logger.Uint("user_id", id))
	response.SuccessWithData(ctx, user)
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
// @Success 200 {object} common.APIResponse{data=response.UserResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}
	// Get ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}
	var req request.UserUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.userService.UpdateUser(ctx, currentUserID, id, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user", logger.Uint("user_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User updated successfully", logger.Uint("user_id", id))
	response.SuccessWithData(ctx, result)
}

// DeleteUser delete user (admin function)
// @Summary Delete user
// @Description Delete a user account (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Get ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	err := c.userService.DeleteUser(ctx, id)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete user", logger.Uint("user_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User deleted successfully", logger.Uint("user_id", id))
	response.Success(ctx)
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
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users/{id}/status [put]
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	// Get ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	var req request.UserUpdateStatusRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.userService.UpdateUserStatus(ctx, id, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Fail to update user status", logger.Uint("user_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User status successful", logger.Uint("user_id", id))
	response.Success(ctx)
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
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/users/{id}/password [put]
func (c *UserController) UpdateUserPassword(ctx *gin.Context) {
	// Get ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	var req request.UserPasswordUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.userService.UpdateUserPassword(ctx, id, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Fail to update user password", logger.Uint("user_id", id), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password successful", logger.Uint("user_id", id))
	response.Success(ctx)
}
