package controller

import (
	resp "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	securityInterface "api-service/internal/interface/service"
	serviceImpl "api-service/internal/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SecurityController Security management controller
type SecurityController struct {
	apiTokenService   securityInterface.APITokenService
	authConfigService securityInterface.AuthConfigService
	oauth2Service     *serviceImpl.OAuth2Service
	twoFactorService  securityInterface.TwoFactorService
	validator         *validator.Validate
	logger            logger.Logger
}

// OAuth2CallbackRequest represents OAuth2 callback request
type OAuth2CallbackRequest struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required"`
	Error string `form:"error"`
}

// OAuth2AuthorizeResponse represents OAuth2 authorization response
type OAuth2AuthorizeResponse struct {
	AuthorizeURL string `json:"authorize_url"`
	State        string `json:"state"`
	Provider     string `json:"provider"`
}

// OAuth2CallbackResponse represents OAuth2 callback response
type OAuth2CallbackResponse struct {
	AccessToken  string                       `json:"access_token"`
	RefreshToken string                       `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time                    `json:"expires_at"`
	TokenType    string                       `json:"token_type"`
	User         *response.UserSimpleResponse `json:"user"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// NewSecurityController creates a new Security management controller instance
func NewSecurityController(
	apiTokenService securityInterface.APITokenService,
	authConfigService securityInterface.AuthConfigService,
	oauth2Service *serviceImpl.OAuth2Service,
	twoFactorService securityInterface.TwoFactorService,
	validator *validator.Validate,
	logger logger.Logger,
) *SecurityController {
	return &SecurityController{
		apiTokenService:   apiTokenService,
		authConfigService: authConfigService,
		oauth2Service:     oauth2Service,
		twoFactorService:  twoFactorService,
		validator:         validator,
		logger:            logger,
	}
}

// RevokeAPIToken revokes API token
// @Summary Revoke API token
// @Description Revoke specified API token
// @Tags API Token Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.RevokeAPITokenRequest true "Revoke API token request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/revoke [post]
func (c *SecurityController) RevokeAPIToken(ctx *gin.Context) {
	var req request.RevokeAPITokenRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Revoke API token by token string
	err := c.apiTokenService.RevokeAPITokenByToken(ctx.Request.Context(), req.Token, userID)
	if err != nil {
		if errors.Is(err, errors.ErrRecordNotFound) {
			resp.BadRequest(ctx, err)
		} else {
			resp.WithError(ctx, err)
		}
		return
	}
	resp.Success(ctx)
}

// RefreshAPIToken refreshes API token
// @Summary Refresh API token
// @Description Refresh API token to extend expiration
// @Tags API Token Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/refresh [get]
func (c *SecurityController) RefreshAPIToken(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Refresh API token for current user
	token, err := c.apiTokenService.RefreshUserAPIToken(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, token)
}

// GetAuthConfig gets authentication config
// @Summary Get authentication config
// @Description Get current authentication configuration
// @Tags Authentication Config
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.AuthConfigResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth-config [get]
func (c *SecurityController) GetAuthConfig(ctx *gin.Context) {
	// Get auth config
	config, err := c.authConfigService.GetAuthConfig(ctx.Request.Context())
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, config)
}

// UpdateAuthConfig updates authentication config
// @Summary Update authentication config
// @Description Update authentication configuration
// @Tags Authentication Config
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.UpdateAuthConfigRequest true "Update auth config request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth-config [put]
func (c *SecurityController) UpdateAuthConfig(ctx *gin.Context) {
	var req request.UpdateAuthConfigRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Update auth config
	err := c.authConfigService.UpdateAuthConfig(ctx.Request.Context(), &req)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// GetOAuth2Providers gets OAuth2 providers
// @Summary Get OAuth2 providers
// @Description Get OAuth2 provider configurations
// @Tags Authentication Config
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]resp.OAuth2ProviderResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth-config/oauth2-providers [get]
func (c *SecurityController) GetOAuth2Providers(ctx *gin.Context) {
	// Get OAuth2 providers
	providers, err := c.authConfigService.GetOAuth2Providers(ctx.Request.Context())
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, providers)
}

// EnableTOTP enables TOTP two-factor authentication
// @Summary Enable TOTP 2FA
// @Description Enable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TOTPSetupResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/enable [post]
func (c *SecurityController) EnableTOTP(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Enable TOTP
	setup, err := c.twoFactorService.EnableTOTP(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, setup)
}

// ConfirmTOTP confirms TOTP setup
// @Summary Confirm TOTP setup
// @Description Confirm TOTP setup with verification code
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.ConfirmTOTPRequest true "Confirm TOTP request"
// @Success 200 {object} response.APIResponse{data=response.TOTPConfirmResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/confirm [post]
func (c *SecurityController) ConfirmTOTP(ctx *gin.Context) {
	var req request.ConfirmTOTPRequest

	// Bind request parameters
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Confirm TOTP
	result, err := c.twoFactorService.ConfirmTOTP(ctx.Request.Context(), userID, req.Code)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, result)
}

// DisableTOTP disables TOTP two-factor authentication
// @Summary Disable TOTP 2FA
// @Description Disable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.DisableTOTPRequest true "Disable TOTP request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/totp/disable [post]
func (c *SecurityController) DisableTOTP(ctx *gin.Context) {
	var req request.DisableTOTPRequest

	// Bind request parameters
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Disable TOTP
	err := c.twoFactorService.DisableTOTP(ctx.Request.Context(), userID, req.Code)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// EnableEmailTwoFactor enables email two-factor authentication
// @Summary Enable Email 2FA
// @Description Enable email two-factor authentication for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.EnableEmailTwoFactorRequest true "Enable email 2FA request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/email/enable [post]
func (c *SecurityController) EnableEmailTwoFactor(ctx *gin.Context) {
	var req request.EnableEmailTwoFactorRequest

	// Bind request parameters
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Enable email two-factor
	err := c.twoFactorService.EnableEmailTwoFactor(ctx.Request.Context(), userID, req.Email)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// DisableEmailTwoFactor disables email two-factor authentication
// @Summary Disable Email 2FA
// @Description Disable email two-factor authentication for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/email/disable [post]
func (c *SecurityController) DisableEmailTwoFactor(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Disable email two-factor
	err := c.twoFactorService.DisableEmailTwoFactor(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// SendEmailCode sends email verification code
// @Summary Send email verification code
// @Description Send verification code to user's email for 2FA
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/email/send-code [post]
func (c *SecurityController) SendEmailCode(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Send email code
	err := c.twoFactorService.SendEmailCode(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// VerifyTwoFactor verifies two-factor authentication code
// @Summary Verify 2FA code
// @Description Verify two-factor authentication code
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.VerifyTwoFactorRequest true "Verify 2FA request"
// @Success 200 {object} response.APIResponse{data=response.TwoFactorVerificationResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/two-factor/verify [post]
func (c *SecurityController) VerifyTwoFactor(ctx *gin.Context) {
	var req request.VerifyTwoFactorRequest

	// Bind request parameters
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Verify two-factor code
	result, err := c.twoFactorService.VerifyTwoFactor(ctx.Request.Context(), req.UserID, req.Code, req.Method)
	if err != nil {
		if err.Error() == "invalid code" || err.Error() == "code expired" {
			resp.Unauthorized(ctx)
			return
		}
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, result)
}

// GetTwoFactorStatus gets user's two-factor authentication status
// @Summary Get 2FA status
// @Description Get user's two-factor authentication status
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TwoFactorStatusResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/two-factor [get]
func (c *SecurityController) GetTwoFactorStatus(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Get two-factor status
	status, err := c.twoFactorService.GetTwoFactorStatus(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, status)
}

// GenerateBackupCodes generates backup codes for two-factor authentication
// @Summary Generate backup codes
// @Description Generate backup codes for two-factor authentication recovery
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.BackupCodesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/backup-codes [post]
func (c *SecurityController) GenerateBackupCodes(ctx *gin.Context) {
	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Generate backup codes
	codes, err := c.twoFactorService.GenerateBackupCodes(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, codes)
}

// DisableTwoFactor disables two-factor authentication
// @Summary Disable two-factor authentication
// @Description Disable two-factor authentication for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/disable [post]
func (c *SecurityController) DisableTwoFactor(ctx *gin.Context) {
	// Get user ID from path parameter
	userID, ok := ParseIDParam(ctx, "user_id")
	if !ok {
		return
	}

	// Get current user ID for authorization check
	currentUserID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Check if user can manage this account (self or admin)
	if userID != currentUserID {
		// TODO: Add admin permission check here
		resp.AccessForbidden(ctx)
		return
	}

	// Disable two-factor authentication
	err := c.twoFactorService.DisableEmailTwoFactor(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.Success(ctx)
}

// GenerateTOTPSecret generates TOTP secret
// @Summary Generate TOTP secret
// @Description Generate TOTP secret for user
// @Tags Two-Factor Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TOTPSecretResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/totp/generate [post]
func (c *SecurityController) GenerateTOTPSecret(ctx *gin.Context) {
	// Get user ID from path parameter
	userID, ok := ParseIDParam(ctx, "user_id")
	if !ok {
		return
	}

	// Get current user ID for authorization check
	currentUserID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Check if user can manage this account (self or admin)
	if userID != currentUserID {
		// TODO: Add admin permission check here
		resp.AccessForbidden(ctx)
		return
	}

	// Generate TOTP secret
	setup, err := c.twoFactorService.EnableTOTP(ctx.Request.Context(), userID)
	if err != nil {
		resp.WithError(ctx, err)
		return
	}
	resp.SuccessWithData(ctx, setup)
}
