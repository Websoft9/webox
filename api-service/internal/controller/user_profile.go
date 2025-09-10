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
	var req request.UserProfileRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// 验证userid参数是否存在
	if req.UserID == 0 {
		c.logger.WarnContext(ctx, "Missing required user ID parameter")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.missing_parameter")))
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request",
		logger.Uint("requestedUserID", req.UserID))

	// 调用服务层获取用户资料
	profile, err := c.profileService.GetUserProfile(ctx, &req)
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

	// 验证userid参数是否存在
	if req.UserID == 0 {
		c.logger.WarnContext(ctx, "Missing required user ID parameter")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.missing_parameter")))
		return
	}

	c.logger.InfoContext(ctx, "Handling update profile request", logger.Uint("userID", req.UserID))

	// 调用服务层更新用户资料
	profile, err := c.profileService.UpdateUserProfile(ctx, req.UserID, &req)
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

	// 验证userid参数是否存在
	if req.UserID == 0 {
		c.logger.WarnContext(ctx, "Missing required user ID parameter")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.missing_parameter")))
		return
	}

	// 获取当前登录的用户ID，用于权限检查（可选）
	currentUserID := ctx.GetUint("user_id")
	if currentUserID == 0 {
		c.logger.WarnContext(ctx, "Missing user ID in context")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "auth.unauthorized")))
		return
	}

	// 如果不是当前用户本人，则需要进行权限检查（可选，取决于业务需求）
	// 这里省略了权限检查的代码，如果需要，可以添加相关逻辑

	c.logger.InfoContext(ctx, "Handling change password request",
		logger.Uint("requestedUserID", req.UserID),
		logger.Uint("currentUserID", currentUserID))

	// 调用服务层修改密码
	// 注意第二个参数其实没有使用了，因为我们在服务层使用的是请求中的UserID
	err := c.profileService.ChangeProfilePassword(ctx, 0, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to change user password", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", req.UserID))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.password_change_success"), nil)
}
