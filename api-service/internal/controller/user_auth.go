package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"
	"api-service/pkg/utils"
	"context"

	"github.com/gin-gonic/gin"
)

// UserAuthController handles user authentication related requests
type UserAuthController struct {
	userAuthService service.UserAuthService
	logger          logger.Logger
	i18n            *i18n.I18n
}

// NewUserAuthController creates a new user authentication controller
func NewUserAuthController(
	userAuthService service.UserAuthService,
	logger logger.Logger,
	i18n *i18n.I18n,
) *UserAuthController {
	return &UserAuthController{
		userAuthService: userAuthService,
		logger:          logger,
		i18n:            i18n,
	}
}

// bindAndValidateRequest binds and validates request parameters
func (c *UserAuthController) bindAndValidateRequest(ctx *gin.Context, req interface{}, action string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.logger.WarnContext(ctx, action+" request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return false
	}
	return true
}

// handleUserAuth handles user authentication related requests with unified error handling
func (c *UserAuthController) handleUserAuth(
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
	pkg_response.Success(ctx, c.i18n.T(ctx, successMessageKey), result)
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user account with email verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserRegisterRequest true "User registration request"
// @Success 200 {object} response.APIResponse{data=response.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/register [post]
func (c *UserAuthController) Register(ctx *gin.Context) {
	var req request.UserRegisterRequest
	c.handleUserAuth(ctx, &req, "registration", func(ctx context.Context, r interface{}) (interface{}, error) {
		return c.userAuthService.Register(ctx, r.(*request.UserRegisterRequest))
	}, "user.register_success")
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user with email and password, return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserLoginRequest true "User login request"
// @Success 200 {object} response.APIResponse{data=response.UserLoginResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (c *UserAuthController) Login(ctx *gin.Context) {
	var req request.UserLoginRequest
	if !c.bindAndValidateRequest(ctx, &req, "login") {
		return
	}

	clientIP := utils.GetRealIP(ctx)
	result, err := c.userAuthService.Login(ctx, &req, clientIP)
	if err != nil {
		c.logger.WarnContext(ctx, "User login failed", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User login successful")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.login_success"), result)
}

// ForgotPassword handles password reset request
// @Summary Forgot password
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/forgot-password [post]
func (c *UserAuthController) ForgotPassword(ctx *gin.Context) {
	var req request.ForgotPasswordRequest
	if !c.bindAndValidateRequest(ctx, &req, "forgot password") {
		return
	}

	err := c.userAuthService.ForgotPassword(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to process forgot password", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Forgot password processed successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.forgot_password_success"), nil)
}

// ShowResetPasswordForm handles password reset form display (GET request from email link)
// @Summary Show reset password form
// @Description Validate reset token and show password reset form
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Password reset token"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/reset-password [get]
func (c *UserAuthController) ShowResetPasswordForm(ctx *gin.Context) {
	// Get token from query parameter
	token := ctx.Query("token")
	if token == "" {
		c.logger.WarnContext(ctx, "Password reset token missing")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Validate token (without consuming it)
	err := c.userAuthService.ValidateResetToken(ctx, token)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid or expired reset token", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Reset password token validated successfully")
	// Return success with token to allow frontend to show reset form
	pkg_response.Success(ctx, c.i18n.T(ctx, "common.success"), map[string]string{"token": token})
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset user password with token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/reset-password [post]
func (c *UserAuthController) ResetPassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest
	if !c.bindAndValidateRequest(ctx, &req, "reset password") {
		return
	}

	err := c.userAuthService.ResetPassword(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to reset password", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Password reset successful")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.reset_password_success"), nil)
}

// VerifyEmail handles email verification
// @Summary Verify email
// @Description Verify user email with token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/verify-email [get]
func (c *UserAuthController) VerifyEmail(ctx *gin.Context) {
	// Get token from query parameter
	token := ctx.Query("token")
	if token == "" {
		c.logger.WarnContext(ctx, "Email verification token missing")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Create request object with token
	req := request.VerifyEmailRequest{
		Token: token,
	}

	err := c.userAuthService.VerifyEmail(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to verify email", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Email verification successful")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.email_verify_success"), nil)
}

// ResendVerificationEmail handles resending verification email
// @Summary Resend verification email
// @Description Resend email verification link to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/resend-verification [post]
func (c *UserAuthController) ResendVerificationEmail(ctx *gin.Context) {
	var req request.ResendVerificationRequest
	if !c.bindAndValidateRequest(ctx, &req, "resend verification email") {
		return
	}

	err := c.userAuthService.ResendVerificationEmail(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to resend verification email", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Verification email resent successfully")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.verification_email_sent"), nil)
}

// OAuth2Login handles OAuth2 authentication
// @Summary OAuth2 login
// @Description Login with OAuth2 provider
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.OAuth2LoginRequest true "OAuth2 login request"
// @Success 200 {object} response.APIResponse{data=response.UserLoginResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/oauth2/login [post]
func (c *UserAuthController) OAuth2Login(ctx *gin.Context) {
	var req request.OAuth2LoginRequest
	if !c.bindAndValidateRequest(ctx, &req, "OAuth2 login") {
		return
	}

	clientIP := utils.GetRealIP(ctx)
	result, err := c.userAuthService.OAuth2Login(ctx, &req, clientIP)
	if err != nil {
		c.logger.ErrorContext(ctx, "User OAuth2 login failed", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User OAuth2 login successful")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.oauth2_login_success"), result)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user and invalidate JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/logout [post]
func (c *UserAuthController) Logout(ctx *gin.Context) {
	// Extract JWT token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		c.logger.WarnContext(ctx, "Authorization header missing during logout")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "common.unauthorized")))
		return
	}

	// Extract token from header (remove "Bearer " prefix)
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.logger.WarnContext(ctx, "Invalid authorization header format during logout")
		errors.HandleError(ctx, errors.NewAppError(errors.CodeUnauthorized, c.i18n.T(ctx, "common.unauthorized")))
		return
	}

	// Call logout service
	err := c.userAuthService.Logout(ctx, token)
	if err != nil {
		c.logger.ErrorContext(ctx, "User logout failed", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User logout successful")
	pkg_response.Success(ctx, c.i18n.T(ctx, "user.logout_success"), nil)
}
