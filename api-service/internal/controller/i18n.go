package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// I18nController i18n控制器
type I18nController struct {
	i18nService service.I18nService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewI18nController 创建新的i18n控制器
func NewI18nController(
	i18nService service.I18nService,
	validator *validator.Validate,
	logger logger.Logger,
) *I18nController {
	return &I18nController{
		i18nService: i18nService,
		validator:   validator,
		logger:      logger,
	}
}

// GetLanguages get supported languages
// @Summary Get supported languages
// @Description Get list of supported languages from configuration
// @Tags Internationalization
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=response.SupportedLanguagesResponse}
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/i18n/languages [get]
func (c *I18nController) GetLanguages(ctx *gin.Context) {
	c.logger.InfoContext(ctx.Request.Context(), "Getting supported languages")

	result, err := c.i18nService.GetSupportedLanguages(ctx.Request.Context())
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get supported languages",
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx.Request.Context(), "Successfully retrieved supported languages",
		logger.Int("languages_count", len(result.Languages)))

	response.SuccessWithData(ctx, result)
}

// SwitchLanguage switches user's language preference
// @Summary Switch user language
// @Description Switch user's language preference and update cache
// @Tags Internationalization
// @Accept json
// @Produce json
// @Param request body request.SwitchLanguageRequest true "Language switch request"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Security BearerAuth
// @Router /api/v1/i18n/switch-language [post]
func (c *I18nController) SwitchLanguage(ctx *gin.Context) {
	c.logger.InfoContext(ctx.Request.Context(), "Switching user language")

	// Get user ID from JWT token
	userID, exists := GetUserID(ctx)
	if !exists {
		return // GetUserID already handles the response
	}

	// Parse request
	var req request.SwitchLanguageRequest
	// Bind request parameters
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	c.logger.InfoContext(ctx.Request.Context(), "Switching language for user",
		logger.Uint("user_id", userID),
		logger.String("new_language", req.Language))

	// Switch language
	if err := c.i18nService.SwitchUserLanguage(ctx.Request.Context(), userID, &req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to switch user language",
			logger.ErrorField(err),
			logger.Uint("user_id", userID),
			logger.String("language", req.Language))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx.Request.Context(), "Successfully switched user language",
		logger.Uint("user_id", userID),
		logger.String("language", req.Language))

	response.Success(ctx)
}
