package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"

	"github.com/gin-gonic/gin"
)

// UserAuthController handles user authentication related requests
type UserAuthController struct {
	userAuthService service.UserAuthService
	logger          logger.Logger
}

// NewUserAuthController creates a new user authentication controller
func NewUserAuthController(
	userAuthService service.UserAuthService,
	logger logger.Logger,
) *UserAuthController {
	return &UserAuthController{
		userAuthService: userAuthService,
		logger:          logger,
	}
}

// bindAndValidateRequest binds and validates request parameters
func (c *UserAuthController) bindAndValidateRequest(ctx *gin.Context, req interface{}, action string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.logger.WarnContext(ctx, action+" request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.ErrValidationFailed)
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
	response.OKWithData(ctx, result, successMessageKey)
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
	response.OKWithData(ctx, result, "user.login_success")
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
	response.OK(ctx, "user.forgot_password_success")
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
		errors.HandleError(ctx, errors.ErrValidationFailed)
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
	response.OKWithData(ctx, map[string]string{"token": token}, "common.success")
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
	response.OK(ctx, "user.reset_password_success")
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
		errors.HandleError(ctx, errors.ErrValidationFailed)
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
	response.OK(ctx, "user.email_verify_success")
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
	response.OK(ctx, "user.verification_email_sent")
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
	response.OKWithData(ctx, result, "user.oauth2_login_success")
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
		errors.HandleError(ctx, errors.ErrInvalidToken)
		return
	}

	// Extract token from header (remove "Bearer " prefix)
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.logger.WarnContext(ctx, "Invalid authorization header format during logout")
		errors.HandleError(ctx, errors.ErrInvalidToken)
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
	response.OK(ctx, "user.logout_success")
}
