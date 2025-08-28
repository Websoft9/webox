package controller

import (
	"github.com/gin-gonic/gin"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/internal/middleware"
	"api-service/pkg/errors"
	pkg_response "api-service/pkg/response"
)

// ProfileController handles profile-related HTTP requests
type ProfileController struct {
	profileService service.ProfileService
}

// NewProfileController creates a new profile controller instance
func NewProfileController(profileService service.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

// GetTwoFactorStatus godoc
// @Summary Get two-factor authentication status
// @Description Get current user two-factor authentication status
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.TwoFactorStatusResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/two-factor [get]
func (pc *ProfileController) GetTwoFactorStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	status, err := pc.profileService.GetTwoFactorStatus(ctx.Request.Context(), userIDUint)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), status)
}

// EnableTwoFactor godoc
// @Summary Enable two-factor authentication
// @Description Enable two-factor authentication for current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TwoFactorEnableRequest true "Two-factor enable request"
// @Success 200 {object} response.Response{data=response.TwoFactorEnableResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/two-factor [post]
func (pc *ProfileController) EnableTwoFactor(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.TwoFactorEnableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	result, err := pc.profileService.EnableTwoFactor(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), result)
}

// VerifyTwoFactor godoc
// @Summary Verify two-factor authentication
// @Description Verify and activate two-factor authentication
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TwoFactorVerifyRequest true "Two-factor verify request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/two-factor/verify [post]
func (pc *ProfileController) VerifyTwoFactor(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.TwoFactorVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.VerifyTwoFactor(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// DisableTwoFactor godoc
// @Summary Disable two-factor authentication
// @Description Disable two-factor authentication for current user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TwoFactorDisableRequest true "Two-factor disable request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/two-factor [delete]
func (pc *ProfileController) DisableTwoFactor(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.TwoFactorDisableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.DisableTwoFactor(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// GetNotificationSettings godoc
// @Summary Get notification settings
// @Description Get current user notification settings
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.NotificationSettingsResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/settings/notification [get]
func (pc *ProfileController) GetNotificationSettings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	settings, err := pc.profileService.GetNotificationSettings(ctx.Request.Context(), userIDUint)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), settings)
}

// UpdateNotificationSettings godoc
// @Summary Update notification settings
// @Description Update current user notification settings
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.NotificationSettingsRequest true "Notification settings request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/settings/notification [put]
func (pc *ProfileController) UpdateNotificationSettings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.NotificationSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.UpdateNotificationSettings(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// GetSecuritySettings godoc
// @Summary Get security settings
// @Description Get current user security settings
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.SecuritySettingsResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/settings/security [get]
func (pc *ProfileController) GetSecuritySettings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	settings, err := pc.profileService.GetSecuritySettings(ctx.Request.Context(), userIDUint)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), settings)
}

// UpdateSecuritySettings godoc
// @Summary Update security settings
// @Description Update current user security settings
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SecuritySettingsRequest true "Security settings request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/settings/security [put]
func (pc *ProfileController) UpdateSecuritySettings(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.SecuritySettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.UpdateSecuritySettings(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=response.ProfileResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile [get]
func (pc *ProfileController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	profile, err := pc.profileService.GetProfile(ctx.Request.Context(), userIDUint)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ProfileUpdateRequest true "Profile update request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile [put]
func (pc *ProfileController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.UpdateProfile(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user password
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.PasswordChangeRequest true "Password change request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/password [put]
func (pc *ProfileController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.PasswordChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInvalidRequest, "Invalid request format"))
		return
	}

	err := pc.profileService.ChangePassword(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), nil)
}

// GetLoginHistory godoc
// @Summary Get login history
// @Description Get current user login history
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param status query string false "Login status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=response.LoginHistoryListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/profile/login-history [get]
func (pc *ProfileController) GetLoginHistory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, "User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeInternalError, "Invalid user ID format"))
		return
	}

	var req request.LoginHistoryListRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "Invalid request format"))
		return
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	history, err := pc.profileService.GetLoginHistory(ctx.Request.Context(), userIDUint, &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, middleware.T(ctx, "common.success"), history)
}
