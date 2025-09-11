package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserProfileController 处理用户个人资料相关的请求
type UserProfileController struct {
	profileService service.UserProfileService
	logger         logger.Logger
	i18n           *i18n.I18n
}

// NewUserProfileController 创建用户个人资料控制器
func NewUserProfileController(
	profileService service.UserProfileService,
	logger logger.Logger,
	i18n *i18n.I18n,
) *UserProfileController {
	return &UserProfileController{
		profileService: profileService,
		logger:         logger,
		i18n:           i18n,
	}
}

// GetProfile 处理获取用户个人资料的请求
// @Summary Get user profile
// @Description Get user profile information by user ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param userid query uint true "User ID to get profile for"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=response.UserProfileResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile [get]
func (c *UserProfileController) GetProfile(ctx *gin.Context) {
	// 解析请求参数
	currentUserID, exists := ctx.Get("user_id")
	if !exists || currentUserID == nil {
		c.logger.WarnContext(ctx, "Missing user ID in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	userID, ok := currentUserID.(uint)
	if !ok {
		c.logger.WarnContext(ctx, "Invalid user ID type in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, c.i18n.T(ctx, "common.internal_error")))
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID))

	// 调用服务层获取用户资料
	profile, err := c.profileService.GetUserProfile(ctx, userID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user profile", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile retrieved successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.get_success"), profile)
}

// UpdateProfile 处理更新用户个人资料的请求
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body request.UserProfileUpdateRequest true "Profile update data"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=response.UserProfileResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile [put]
// UpdateProfile 处理更新用户个人资料的请求
func (c *UserProfileController) UpdateProfile(ctx *gin.Context) {

	// 解析请求参数
	var req request.UserProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// 获取当前登录的用户ID
	currentUserID, exists := ctx.Get("user_id")
	if !exists || currentUserID == nil {
		c.logger.WarnContext(ctx, "Missing user ID in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	userID, ok := currentUserID.(uint)
	if !ok {
		c.logger.WarnContext(ctx, "Invalid user ID type in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, c.i18n.T(ctx, "common.internal_error")))
		return
	}

	c.logger.InfoContext(ctx, "Handling update profile request", logger.Uint("userID", userID))

	// 调用服务层更新用户资料
	profile, err := c.profileService.UpdateUserProfile(ctx, userID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user profile", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile updated successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.update_success"), profile)
}

// ChangePassword 处理修改用户个人密码的请求
// @Summary Change user password
// @Description Change user's password by user ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body request.ProfileChangePasswordRequest true "Password change request"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/password [put]
func (c *UserProfileController) ChangePassword(ctx *gin.Context) {
	// 解析请求参数
	var req request.ProfileChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// 获取当前登录的用户ID
	currentUserID, exists := ctx.Get("user_id")
	if !exists || currentUserID == nil {
		c.logger.WarnContext(ctx, "Missing user ID in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	userID, ok := currentUserID.(uint)
	if !ok {
		c.logger.WarnContext(ctx, "Invalid user ID type in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, c.i18n.T(ctx, "common.internal_error")))
		return
	}

	c.logger.InfoContext(ctx, "Handling update profile request", logger.Uint("userID", userID))

	// 调用服务层修改密码
	err := c.profileService.ChangeProfilePassword(ctx, userID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to change user password", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", userID))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.password_change_success"), nil)
}

// GetLoginHistories 获取登录历史记录
// @Summary 获取登录历史记录
// @Description 获取当前用户的登录历史记录
// @Tags 用户配置文件
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse{data=response.LoginHistoryResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/login-history [get]
func (c *UserProfileController) GetLoginHistories(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.WarnContext(ctx, "User ID not found in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	currentUserID, ok := userID.(uint)
	if !ok {
		c.logger.WarnContext(ctx, "Invalid user ID type in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	// 绑定请求参数
	var req request.LoginHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "Login history request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// 获取登录历史
	result, err := c.profileService.GetLoginHistories(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get login histories", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Login histories retrieved successfully", logger.Uint("user_id", currentUserID))
	pkg_response.Success(ctx, c.i18n.T(ctx, "profile.login_history_get_success"), result)
}
