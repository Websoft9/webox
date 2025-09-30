package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserAuthController handles user authentication related requests
type UserAuthController struct {
	userAuthService service.UserAuthService
	validator       *validator.Validate
	logger          logger.Logger
}

// NewUserAuthController creates a new user authentication controller
func NewUserAuthController(
	userAuthService service.UserAuthService,
	validator *validator.Validate,
	logger logger.Logger,
) *UserAuthController {
	return &UserAuthController{
		userAuthService: userAuthService,
		validator:       validator,
		logger:          logger,
	}
}

// handleLoginRequest handles login-type requests with unified error handling
func (c *UserAuthController) handleLoginRequest(ctx *gin.Context, req interface{}, action string, loginFunc func() (interface{}, error)) {
	if !BindAndValidateRequest(ctx, req, c.validator, c.logger) {
		return
	}

	result, err := loginFunc()
	if err != nil {
		c.logger.WarnContext(ctx, "User "+action+" failed", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User "+action+" successful")
	response.SuccessWithData(ctx, result)
}

// handleUserAuth handles user authentication related requests with unified error handling
func (c *UserAuthController) handleUserAuth(
	ctx *gin.Context,
	req interface{},
	action string,
	serviceFunc func(context.Context, interface{}) (interface{}, error),
) {
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := serviceFunc(ctx, req)
	if err != nil {
		if action == "login" {
			c.logger.WarnContext(ctx, "User "+action+" failed", logger.ErrorField(err))
		} else {
			c.logger.ErrorContext(ctx, "User "+action+" failed", logger.ErrorField(err))
		}
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User "+action+" successful")
	response.SuccessWithData(ctx, result)
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
	})
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
	c.handleLoginRequest(ctx, &req, "login", func() (interface{}, error) {
		clientIP := utils.GetRealIP(ctx)
		return c.userAuthService.Login(ctx, &req, clientIP)
	})
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
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.userAuthService.ForgotPassword(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to process forgot password", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Forgot password processed successfully")
	response.Success(ctx)
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
		response.WithError(ctx, errors.ErrValidationFailed)
		return
	}

	// Validate token (without consuming it)
	err := c.userAuthService.ValidateResetToken(ctx, token)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid or expired reset token", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Reset password token validated successfully")
	// Return success with token to allow frontend to show reset form
	response.SuccessWithData(ctx, map[string]string{"token": token})
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
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.userAuthService.ResetPassword(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to reset password", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Password reset successful")
	response.Success(ctx)
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
		response.WithError(ctx, errors.ErrValidationFailed)
		return
	}

	// Create request object with token
	req := request.VerifyEmailRequest{
		Token: token,
	}

	err := c.userAuthService.VerifyEmail(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to verify email", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Email verification successful")
	response.Success(ctx)
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
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.userAuthService.ResendVerificationEmail(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to resend verification email", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Verification email resent successfully")
	response.Success(ctx)
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
	c.handleLoginRequest(ctx, &req, "OAuth2 login", func() (interface{}, error) {
		clientIP := utils.GetRealIP(ctx)
		return c.userAuthService.OAuth2Login(ctx, &req, clientIP)
	})
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
		response.WithError(ctx, errors.ErrInvalidToken)
		return
	}

	// Extract token from header (remove "Bearer " prefix)
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.logger.WarnContext(ctx, "Invalid authorization header format during logout")
		response.WithError(ctx, errors.ErrInvalidToken)
		return
	}

	// Call logout service
	err := c.userAuthService.Logout(ctx, token)
	if err != nil {
		c.logger.ErrorContext(ctx, "User logout failed", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "User logout successful")
	response.Success(ctx)
}
