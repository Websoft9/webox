package controller

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	securityInterface "api-service/internal/interface/service"
	serviceImpl "api-service/internal/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"net/http"
	"strconv"
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
	i18n              *i18n.I18n
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
	i18n *i18n.I18n,
) *SecurityController {
	return &SecurityController{
		apiTokenService:   apiTokenService,
		authConfigService: authConfigService,
		oauth2Service:     oauth2Service,
		twoFactorService:  twoFactorService,
		validator:         validator,
		logger:            logger,
		i18n:              i18n,
	}
}

// CreateAPIToken creates a new API token
// @Summary Create API token
// @Description Create a new API token
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param request body request.CreateAPITokenRequest true "Create API token request"
// @Success 201 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/api-tokens [post]
func (c *SecurityController) CreateAPIToken(ctx *gin.Context) {
	var req request.CreateAPITokenRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Create API token
	token, err := c.apiTokenService.CreateAPIToken(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to create API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.create_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": c.i18n.T(ctx, "api_token.create_success"),
		"data":    token,
	})
}

// GetAPIToken gets API token details
// @Summary Get API token details
// @Description Get API token details by ID
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id} [get]
func (c *SecurityController) GetAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Get API token
	token, err := c.apiTokenService.GetAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == constants.ErrTokenNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.get_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    token,
	})
}

// UpdateAPIToken updates API token
// @Summary Update API token
// @Description Update API token information
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Param request body request.UpdateAPITokenRequest true "Update API token request"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id} [put]
func (c *SecurityController) UpdateAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	var req request.UpdateAPITokenRequest

	// Bind request parameters
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(bindErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   bindErr.Error(),
		})
		return
	}

	// Validate request parameters
	if validateErr := c.validator.Struct(&req); validateErr != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(validateErr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   validateErr.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Update API token
	token, err := c.apiTokenService.UpdateAPIToken(ctx.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == constants.ErrTokenNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to update API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.update_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.update_success"),
		"data":    token,
	})
}

// RevokeAPIToken revokes API token
// @Summary Revoke API token
// @Description Revoke specified API token
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id}/revoke [post]
func (c *SecurityController) RevokeAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Revoke API token
	err = c.apiTokenService.RevokeAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == constants.ErrTokenNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to revoke API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.revoke_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.revoke_success"),
	})
}

// RefreshAPIToken refreshes API token
// @Summary Refresh API token
// @Description Refresh API token to extend expiration
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id}/refresh [post]
func (c *SecurityController) RefreshAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Refresh API token
	token, err := c.apiTokenService.RefreshAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == constants.ErrTokenNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to refresh API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.refresh_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.refresh_success"),
		"data":    token,
	})
}

// ListAPITokens gets API token list
// @Summary Get API token list
// @Description Get paginated API token list
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.APITokenListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/api-tokens [get]
func (c *SecurityController) ListAPITokens(ctx *gin.Context) {
	var req request.ListAPITokensRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_query_parameters"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Query validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.query_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Get API token list
	tokens, err := c.apiTokenService.ListAPITokens(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to list API tokens", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.list_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    tokens,
	})
}

// ValidateAPIToken validates API token
// @Summary Validate API token
// @Description Validate API token and return token information
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param request body request.ValidateAPITokenRequest true "Validate API token request"
// @Success 200 {object} response.APIResponse{data=response.APITokenValidationResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/api-tokens/validate [post]
func (c *SecurityController) ValidateAPIToken(ctx *gin.Context) {
	var req request.ValidateAPITokenRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Validate API token
	validation, err := c.apiTokenService.ValidateAPIToken(ctx.Request.Context(), req.Token)
	if err != nil {
		if err.Error() == "invalid token" || err.Error() == "token expired" || err.Error() == "token revoked" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": c.i18n.T(ctx, "api_token.invalid_token"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to validate API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.validate_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.validate_success"),
		"data":    validation,
	})
}

// BatchRevokeAPITokens batch revokes API tokens
// @Summary Batch revoke API tokens
// @Description Batch revoke multiple API tokens
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param request body request.BatchRevokeAPITokensRequest true "Batch revoke API tokens request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/api-tokens [delete]
func (c *SecurityController) BatchRevokeAPITokens(ctx *gin.Context) {
	var req request.BatchRevokeAPITokensRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Batch revoke API tokens
	err := c.apiTokenService.BatchRevokeAPITokens(ctx.Request.Context(), req.IDs, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to batch revoke API tokens", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.batch_revoke_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.batch_revoke_success"),
	})
}

// GetAuthConfig gets authentication config
// @Summary Get authentication config
// @Description Get current authentication configuration
// @Tags Authentication Config
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.AuthConfigResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth-config [get]
func (c *SecurityController) GetAuthConfig(ctx *gin.Context) {
	// Get auth config
	config, err := c.authConfigService.GetAuthConfig(ctx.Request.Context())
	if err != nil {
		ResponseInternalError(ctx, err, "auth_config.get_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, config, "common.success", c.i18n)
}

// UpdateAuthConfig updates authentication config
// @Summary Update authentication config
// @Description Update authentication configuration
// @Tags Authentication Config
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
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	// Update auth config
	err := c.authConfigService.UpdateAuthConfig(ctx.Request.Context(), &req)
	if err != nil {
		ResponseInternalError(ctx, err, "auth_config.update_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, nil, "auth_config.update_success", c.i18n)
}

// GetOAuth2Providers gets OAuth2 providers
// @Summary Get OAuth2 providers
// @Description Get OAuth2 provider configurations
// @Tags Authentication Config
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]response.OAuth2ProviderResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth-config/oauth2-providers [get]
func (c *SecurityController) GetOAuth2Providers(ctx *gin.Context) {
	// Get OAuth2 providers
	providers, err := c.authConfigService.GetOAuth2Providers(ctx.Request.Context())
	if err != nil {
		ResponseInternalError(ctx, err, "auth_config.oauth2_providers_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, providers, "common.success", c.i18n)
}

// EnableTOTP enables TOTP two-factor authentication
// @Summary Enable TOTP 2FA
// @Description Enable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TOTPSetupResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/enable [post]
func (c *SecurityController) EnableTOTP(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Enable TOTP
	setup, err := c.twoFactorService.EnableTOTP(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to enable TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.enable_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.enable_totp_success"),
		"data":    setup,
	})
}

// ConfirmTOTP confirms TOTP setup
// @Summary Confirm TOTP setup
// @Description Confirm TOTP setup with verification code
// @Tags Two-Factor Authentication
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Confirm TOTP
	result, err := c.twoFactorService.ConfirmTOTP(ctx.Request.Context(), userID.(uint), req.Code)
	if err != nil {
		if err.Error() == constants.ErrInvalidCode {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to confirm TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.confirm_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.confirm_totp_success"),
		"data":    result,
	})
}

// DisableTOTP disables TOTP two-factor authentication
// @Summary Disable TOTP 2FA
// @Description Disable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Disable TOTP
	err := c.twoFactorService.DisableTOTP(ctx.Request.Context(), userID.(uint), req.Code)
	if err != nil {
		if err.Error() == constants.ErrInvalidCode {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to disable TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.disable_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.disable_totp_success"),
	})
}

// EnableEmailTwoFactor enables email two-factor authentication
// @Summary Enable Email 2FA
// @Description Enable email two-factor authentication for user
// @Tags Two-Factor Authentication
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Enable email two-factor
	err := c.twoFactorService.EnableEmailTwoFactor(ctx.Request.Context(), userID.(uint), req.Email)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to enable email 2FA", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.enable_email_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.enable_email_success"),
	})
}

// DisableEmailTwoFactor disables email two-factor authentication
// @Summary Disable Email 2FA
// @Description Disable email two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/email/disable [post]
func (c *SecurityController) DisableEmailTwoFactor(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Disable email two-factor
	err := c.twoFactorService.DisableEmailTwoFactor(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to disable email 2FA", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.disable_email_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.disable_email_success"),
	})
}

// SendEmailCode sends email verification code
// @Summary Send email verification code
// @Description Send verification code to user's email for 2FA
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/email/send-code [post]
func (c *SecurityController) SendEmailCode(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Send email code
	err := c.twoFactorService.SendEmailCode(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to send email code", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.send_email_code_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.send_email_code_success"),
	})
}

// VerifyTwoFactor verifies two-factor authentication code
// @Summary Verify 2FA code
// @Description Verify two-factor authentication code
// @Tags Two-Factor Authentication
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Verify two-factor code
	result, err := c.twoFactorService.VerifyTwoFactor(ctx.Request.Context(), req.UserID, req.Code, req.Method)
	if err != nil {
		if err.Error() == "invalid code" || err.Error() == "code expired" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to verify 2FA code", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.verify_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.verify_success"),
		"data":    result,
	})
}

// GetTwoFactorStatus gets user's two-factor authentication status
// @Summary Get 2FA status
// @Description Get user's two-factor authentication status
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TwoFactorStatusResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/two-factor [get]
func (c *SecurityController) GetTwoFactorStatus(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Get two-factor status
	status, err := c.twoFactorService.GetTwoFactorStatus(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get 2FA status", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.status_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    status,
	})
}

// GenerateBackupCodes generates backup codes for two-factor authentication
// @Summary Generate backup codes
// @Description Generate backup codes for two-factor authentication recovery
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.BackupCodesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/backup-codes [post]
func (c *SecurityController) GenerateBackupCodes(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Generate backup codes
	codes, err := c.twoFactorService.GenerateBackupCodes(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to generate backup codes", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.backup_codes_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.backup_codes_success"),
		"data":    codes,
	})
}

// DisableTwoFactor disables two-factor authentication
// @Summary Disable two-factor authentication
// @Description Disable two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/disable [post]
func (c *SecurityController) DisableTwoFactor(ctx *gin.Context) {
	// Get user ID from path parameter
	userID, ok := ParseIDParam(ctx, "user_id", "validation.invalid_user_id", c.i18n)
	if !ok {
		return
	}

	// Get current user ID for authorization check
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	// Check if user can manage this account (self or admin)
	if userID != currentUserID.(uint) {
		// TODO: Add admin permission check here
		ResponseForbidden(ctx, "auth.insufficient_permissions", c.i18n)
		return
	}

	// Disable two-factor authentication
	err := c.twoFactorService.DisableEmailTwoFactor(ctx.Request.Context(), userID)
	if err != nil {
		ResponseInternalError(ctx, err, "two_factor.disable_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, nil, "two_factor.disable_success", c.i18n)
}

// GenerateTOTPSecret generates TOTP secret
// @Summary Generate TOTP secret
// @Description Generate TOTP secret for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TOTPSecretResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/two-factor/totp/generate [post]
func (c *SecurityController) GenerateTOTPSecret(ctx *gin.Context) {
	// Get user ID from path parameter
	userID, ok := ParseIDParam(ctx, "user_id", "validation.invalid_user_id", c.i18n)
	if !ok {
		return
	}

	// Get current user ID for authorization check
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	// Check if user can manage this account (self or admin)
	if userID != currentUserID.(uint) {
		// TODO: Add admin permission check here
		ResponseForbidden(ctx, "auth.insufficient_permissions", c.i18n)
		return
	}

	// Generate TOTP secret
	setup, err := c.twoFactorService.EnableTOTP(ctx.Request.Context(), userID)
	if err != nil {
		ResponseInternalError(ctx, err, "two_factor.totp_generate_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, setup, "two_factor.totp_generate_success", c.i18n)
}
