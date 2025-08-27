package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthConfigController authentication config controller
type AuthConfigController struct {
	authConfigService service.AuthConfigService
	validator         *validator.Validate
	logger            logger.Logger
	i18n              *i18n.I18n
}

// NewAuthConfigController creates a new auth config controller instance
func NewAuthConfigController(
	authConfigService service.AuthConfigService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *AuthConfigController {
	return &AuthConfigController{
		authConfigService: authConfigService,
		validator:         validator,
		logger:            logger,
		i18n:              i18n,
	}
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
func (c *AuthConfigController) GetAuthConfig(ctx *gin.Context) {
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
func (c *AuthConfigController) UpdateAuthConfig(ctx *gin.Context) {
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
func (c *AuthConfigController) GetOAuth2Providers(ctx *gin.Context) {
	// Get OAuth2 providers
	providers, err := c.authConfigService.GetOAuth2Providers(ctx.Request.Context())
	if err != nil {
		ResponseInternalError(ctx, err, "auth_config.oauth2_providers_failed", c.logger, c.i18n)
		return
	}

	ResponseOK(ctx, providers, "common.success", c.i18n)
}
