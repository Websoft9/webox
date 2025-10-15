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

// UserProfileController handles requests related to user profile
type UserProfileController struct {
	profileService service.UserProfileService
	logger         logger.Logger
	i18n           *i18n.I18n
	validator      *validator.Validate
}

// NewUserProfileController creates a new user profile controller
func NewUserProfileController(
	profileService service.UserProfileService,
	logger logger.Logger,
	i18n *i18n.I18n,
	validator *validator.Validate,
) *UserProfileController {
	return &UserProfileController{
		profileService: profileService,
		logger:         logger,
		i18n:           i18n,
		validator:      validator,
	}
}

// GetProfile handles requests to get user profile information
// @Summary Get user profile
// @Description Get user profile information by user ID
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=response.UserProfileResponse}
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile [get]
func (c *UserProfileController) GetProfile(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	// Call service layer to get user profile
	profile, err := c.profileService.GetUserProfile(ctx, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get user profile", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile retrieved successfully")
	response.SuccessWithData(ctx, profile)
}

// UpdateProfile handles requests to update user profile information
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body request.UserProfileUpdateRequest true "Profile update data"
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=response.UserProfileResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile [put]
func (c *UserProfileController) UpdateProfile(ctx *gin.Context) {
	// Parse request parameters
	var req request.UserProfileUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	// Call service layer to update user profile
	profile, err := c.profileService.UpdateUserProfile(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update user profile", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User profile updated successfully")
	response.SuccessWithData(ctx, profile)
}

// ChangePassword handles requests to change user password
// @Summary Change user password
// @Description Change user's password by user ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body request.ProfileChangePasswordRequest true "Password change request"
// @Security BearerAuth
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/password [put]
func (c *UserProfileController) ChangePassword(ctx *gin.Context) {
	// Parse request parameters
	var req request.ProfileChangePasswordRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	// Call service layer to change password
	err := c.profileService.ChangeProfilePassword(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to change user password", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", currentUserID))
	response.Success(ctx)
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
// @Success 200 {object} common.APIResponse{data=response.LoginHistoryResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/login-history [get]
func (c *UserProfileController) GetLoginHistories(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	// Bind request parameters
	var req request.LoginHistoryRequest
	// Bind and validate request
	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get login history
	result, err := c.profileService.GetLoginHistories(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get login histories", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Login histories retrieved successfully", logger.Uint("user_id", currentUserID))
	response.SuccessWithData(ctx, result)
}

// GetNotificationSettings gets notification settings
// @Summary Get notification settings
// @Description Get notification settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=response.NotificationSettingsResponse}
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/notification-settings [get]
func (c *UserProfileController) GetNotificationSettings(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	settings, err := c.profileService.GetNotificationSettings(ctx, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get notification settings", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Notification settings retrieved successfully", logger.Uint("user_id", currentUserID))
	response.SuccessWithData(ctx, settings)
}

// UpdateNotificationSettings updates notification settings
// @Summary Update notification settings
// @Description Update notification settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.NotificationSettingsRequest true "Notification settings request"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/notification-settings [put]
func (c *UserProfileController) UpdateNotificationSettings(ctx *gin.Context) {
	// Parse request parameters
	var req request.NotificationSettingsRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	err := c.profileService.UpdateNotificationSettings(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update notification settings", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Notification settings updated successfully", logger.Uint("user_id", currentUserID))
	response.Success(ctx)
}

// GetSecuritySettings gets security settings
// @Summary Get security settings
// @Description Get security settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=response.SecuritySettingsResponse}
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/security-settings [get]
func (c *UserProfileController) GetSecuritySettings(ctx *gin.Context) {
	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	settings, err := c.profileService.GetSecuritySettings(ctx, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get security settings", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Security settings retrieved successfully", logger.Uint("user_id", currentUserID))
	response.SuccessWithData(ctx, settings)
}

// UpdateSecuritySettings updates security settings
// @Summary Update security settings
// @Description Update security settings for the current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SecuritySettingsRequest true "Security settings request"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/profile/security-settings [put]
func (c *UserProfileController) UpdateSecuritySettings(ctx *gin.Context) {
	// Parse request parameters
	var req request.SecuritySettingsRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get profile request", logger.Uint("userID", currentUserID))

	err := c.profileService.UpdateSecuritySettings(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update security settings", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Security settings updated successfully", logger.Uint("user_id", currentUserID))
	response.Success(ctx)
}
