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

// UserProfileController handles requests related to user profile
type UserProfileController struct {
	profileService service.UserProfileService
	logger         logger.Logger
	i18n           *i18n.I18n
}

// NewUserProfileController creates a new user profile controller
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

// GetProfile handles requests to get user profile information
// @Summary Get user profile
// @Description Get user profile information by user ID
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=response.UserProfileResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile [get]
func (c *UserProfileController) GetProfile(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	// Call service layer to get user profile
	profile, err := c.profileService.GetUserProfile(ctx, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user profile", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile retrieved successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.get_success"), profile)
}

// UpdateProfile handles requests to update user profile information
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
func (c *UserProfileController) UpdateProfile(ctx *gin.Context) {
	// Parse request parameters
	var req request.UserProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	// Call service layer to update user profile
	profile, err := c.profileService.UpdateUserProfile(ctx, userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user profile", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile updated successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.update_success"), profile)
}

// ChangePassword handles requests to change user password
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
	// Parse request parameters
	var req request.ProfileChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	// Call service layer to change password
	err := c.profileService.ChangeProfilePassword(ctx, userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to change user password", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.password_change_success"), nil)
}

// GetLoginHistories gets login history records
// @Summary Get login history records
// @Description Get login history records for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {object} response.APIResponse{data=response.LoginHistoryResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/login-history [get]
func (c *UserProfileController) GetLoginHistories(ctx *gin.Context) {
	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	// Bind request parameters
	var req request.LoginHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "Login history request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get login history
	result, err := c.profileService.GetLoginHistories(ctx, userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get login histories", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Login histories retrieved successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.login_history_get_success"), result)
}

// GetNotificationSettings gets notification settings
// @Summary Get notification settings
// @Description Get notification settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=response.NotificationSettingsResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/notification-settings [get]
func (c *UserProfileController) GetNotificationSettings(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	settings, err := c.profileService.GetNotificationSettings(ctx, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get notification settings", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Notification settings retrieved successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.notification_settings_get_success"), settings)
}

// UpdateNotificationSettings updates notification settings
// @Summary Update notification settings
// @Description Update notification settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.NotificationSettingsRequest true "Notification settings request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/notification-settings [put]
func (c *UserProfileController) UpdateNotificationSettings(ctx *gin.Context) {
	// Parse request parameters
	var req request.NotificationSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	err := c.profileService.UpdateNotificationSettings(ctx, userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update notification settings", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Notification settings updated successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.notification_settings_update_success"), nil)
}

// GetSecuritySettings gets security settings
// @Summary Get security settings
// @Description Get security settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=response.SecuritySettingsResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/security-settings [get]
func (c *UserProfileController) GetSecuritySettings(ctx *gin.Context) {
	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	settings, err := c.profileService.GetSecuritySettings(ctx, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get security settings", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Security settings retrieved successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.security_settings_get_success"), settings)
}

// UpdateSecuritySettings updates security settings
// @Summary Update security settings
// @Description Update security settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SecuritySettingsRequest true "Security settings request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/profile/security-settings [put]
func (c *UserProfileController) UpdateSecuritySettings(ctx *gin.Context) {
	// Parse request parameters
	var req request.SecuritySettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", userID.(uint)))

	err := c.profileService.UpdateSecuritySettings(ctx, userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update security settings", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Security settings updated successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "user_profile.security_settings_update_success"), nil)
}
