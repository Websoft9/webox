package controller

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/response"
	"api-service/internal/service"
	"api-service/pkg/auth"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// OAuth2Controller handles OAuth2 authentication endpoints
type OAuth2Controller struct {
	oauth2Service *service.OAuth2Service
	validator     *validator.Validate
	logger        logger.Logger
	i18n          *i18n.I18n
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

// NewOAuth2Controller creates a new OAuth2 controller instance
func NewOAuth2Controller(
	oauth2Service *service.OAuth2Service,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *OAuth2Controller {
	return &OAuth2Controller{
		oauth2Service: oauth2Service,
		validator:     validator,
		logger:        logger,
		i18n:          i18n,
	}
}

// Authorize initiates OAuth2 authorization flow
// @Summary Initiate OAuth2 authorization
// @Description Start OAuth2 authorization flow with specified provider
// @Tags OAuth2
// @Accept json
// @Produce json
// @Param provider path string true "OAuth2 provider name"
// @Success 200 {object} response.APIResponse{data=OAuth2AuthorizeResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/oauth2/{provider}/authorize [get]
func (c *OAuth2Controller) Authorize(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ResponseBadRequest(ctx, fmt.Errorf("provider parameter is required"), "oauth2.provider_required", c.i18n)
		return
	}

	// Check if provider is enabled
	if !c.oauth2Service.IsProviderEnabled(provider) {
		ResponseBadRequest(ctx, fmt.Errorf("provider '%s' is not enabled", provider), "oauth2.provider_disabled", c.i18n)
		return
	}

	// Generate state parameter for CSRF protection
	state, err := c.generateState()
	if err != nil {
		ResponseInternalError(ctx, err, "oauth2.state_generation_failed", c.logger, c.i18n)
		return
	}

	// Store state in session or cache (for production, use Redis or similar)
	// For now, we'll just return it and expect the frontend to handle it
	ctx.SetCookie("oauth2_state", state, constants.OAuth2CookieMaxAge, "/", "", false, true)

	// Generate authorization URL
	authorizeURL, err := c.oauth2Service.GetAuthorizationURL(ctx.Request.Context(), provider, state)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to generate OAuth2 authorization URL", logger.ErrorField(err))
		ResponseInternalError(ctx, err, "oauth2.authorization_url_failed", c.logger, c.i18n)
		return
	}

	response := &OAuth2AuthorizeResponse{
		AuthorizeURL: authorizeURL,
		State:        state,
		Provider:     provider,
	}

	ResponseOK(ctx, response, "oauth2.authorization_url_generated", c.i18n)
}

// Callback handles OAuth2 callback from provider
// @Summary Handle OAuth2 callback
// @Description Process OAuth2 callback and complete authentication
// @Tags OAuth2
// @Accept json
// @Produce json
// @Param provider path string true "OAuth2 provider name"
// @Param code query string true "Authorization code from provider"
// @Param state query string true "State parameter for CSRF protection"
// @Param error query string false "Error from provider"
// @Success 200 {object} response.APIResponse{data=OAuth2CallbackResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/oauth2/{provider}/callback [get]
func (c *OAuth2Controller) Callback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ResponseBadRequest(ctx, fmt.Errorf("provider parameter is required"), "oauth2.provider_required", c.i18n)
		return
	}

	var req OAuth2CallbackRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ResponseBadRequest(ctx, err, "validation.invalid_request_format", c.i18n)
		return
	}

	// Check for OAuth2 error
	if req.Error != "" {
		c.logger.WarnContext(ctx.Request.Context(), "OAuth2 provider returned error",
			logger.String("provider", provider),
			logger.String("error", req.Error))
		ResponseBadRequest(ctx, fmt.Errorf("OAuth2 error: %s", req.Error), "oauth2.provider_error", c.i18n)
		return
	}

	// Validate state parameter
	storedState, err := ctx.Cookie("oauth2_state")
	if err != nil {
		ResponseBadRequest(ctx, fmt.Errorf("OAuth2 state not found"), "oauth2.state_not_found", c.i18n)
		return
	}

	if validateErr := c.oauth2Service.ValidateState(ctx.Request.Context(), req.State, storedState); validateErr != nil {
		ResponseBadRequest(ctx, validateErr, "oauth2.invalid_state", c.i18n)
		return
	}

	// Clear state cookie
	ctx.SetCookie("oauth2_state", "", -1, "/", "", false, true)

	// Exchange code for token
	tokenResp, err := c.oauth2Service.ExchangeCodeForToken(ctx.Request.Context(), provider, req.Code)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to exchange OAuth2 code for token", logger.ErrorField(err))
		ResponseInternalError(ctx, err, "oauth2.token_exchange_failed", c.logger, c.i18n)
		return
	}

	// Get user info from provider
	userInfo, err := c.oauth2Service.GetUserInfo(ctx.Request.Context(), provider, tokenResp.AccessToken)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get OAuth2 user info", logger.ErrorField(err))
		ResponseInternalError(ctx, err, "oauth2.user_info_failed", c.logger, c.i18n)
		return
	}

	// Find or create user
	user := c.findOrCreateUser(userInfo)

	// Generate JWT token for the user
	jwtAuth := auth.GetGlobalJWT()
	if jwtAuth == nil {
		ResponseInternalError(ctx, fmt.Errorf("JWT service not initialized"), "oauth2.jwt_not_initialized", c.logger, c.i18n)
		return
	}

	accessToken, expiresAt, err := jwtAuth.GenerateTokenWithUserInfo(user.ID, user.Username, "user")
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to generate JWT token", logger.ErrorField(err))
		ResponseInternalError(ctx, err, "oauth2.jwt_generation_failed", c.logger, c.i18n)
		return
	}

	callbackResp := &OAuth2CallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
		User:         user,
	}

	c.logger.InfoContext(ctx.Request.Context(), "OAuth2 authentication successful",
		logger.String("provider", provider),
		logger.Uint("user_id", user.ID))

	ResponseOK(ctx, callbackResp, "oauth2.authentication_successful", c.i18n)
}

// GetProviders returns list of enabled OAuth2 providers
// @Summary Get OAuth2 providers
// @Description Get list of enabled OAuth2 providers
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]response.OAuth2ProviderResponse}
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/oauth2/providers [get]
func (c *OAuth2Controller) GetProviders(ctx *gin.Context) {
	providers := c.oauth2Service.GetEnabledProviders()

	// Convert to response format
	resp := make([]*response.OAuth2ProviderResponse, 0, len(providers))
	for i := range providers {
		provider := &providers[i]
		resp = append(resp, &response.OAuth2ProviderResponse{
			Name:         provider.Name,
			Provider:     c.findProviderKey(provider),
			Enabled:      provider.Enabled,
			AutoRegister: provider.AutoRegister,
			UserMapping:  provider.UserMapping,
		})
	}

	ResponseOK(ctx, resp, "oauth2.providers_retrieved", c.i18n)
}

// RefreshToken refreshes OAuth2 access token
// @Summary Refresh OAuth2 token
// @Description Refresh OAuth2 access token using refresh token
// @Tags OAuth2
// @Accept json
// @Produce json
// @Param provider path string true "OAuth2 provider name"
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.APIResponse{data=service.OAuth2TokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/oauth2/{provider}/refresh [post]
func (c *OAuth2Controller) RefreshToken(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider == "" {
		ResponseBadRequest(ctx, fmt.Errorf("provider parameter is required"), "oauth2.provider_required", c.i18n)
		return
	}

	var req RefreshTokenRequest
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	// Refresh token
	tokenResp, err := c.oauth2Service.RefreshToken(ctx.Request.Context(), provider, req.RefreshToken)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to refresh OAuth2 token", logger.ErrorField(err))
		ResponseInternalError(ctx, err, "oauth2.token_refresh_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, tokenResp, "oauth2.token_refreshed", c.i18n)
}

// Helper methods

// generateState generates a secure random state parameter
func (c *OAuth2Controller) generateState() (string, error) {
	bytes := make([]byte, constants.OAuth2StateLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// findOrCreateUser finds existing user or creates new one based on OAuth2 info
func (c *OAuth2Controller) findOrCreateUser(userInfo *service.OAuth2UserInfo) *response.UserSimpleResponse {
	// This is a placeholder implementation
	// In a real implementation, you would:
	// 1. Check if user exists by OAuth2 provider ID
	// 2. Check if user exists by email
	// 3. Create new user if not found and auto-registration is enabled
	// 4. Link OAuth2 account to existing user

	// For now, return a mock user
	return &response.UserSimpleResponse{
		ID:       1,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}
}

// findProviderKey finds provider key for a given provider config
func (c *OAuth2Controller) findProviderKey(providerConfig *config.OAuth2ProviderConfig) string {
	// This would typically be stored in the provider config or determined by name
	switch providerConfig.Name {
	case "Google":
		return "google"
	case "GitHub":
		return "github"
	case "企业微信":
		return "wechat_work"
	case "钉钉":
		return "dingtalk"
	case "GitLab":
		return "gitlab"
	case "Azure AD":
		return "azure_ad"
	default:
		return strings.ToLower(strings.ReplaceAll(providerConfig.Name, " ", "_"))
	}
}
