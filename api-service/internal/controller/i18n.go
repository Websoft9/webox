package controller

import (
	"api-service/internal/dto/response"
	"api-service/internal/middleware"
	"api-service/pkg/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

// I18nController i18n控制器
type I18nController struct{}

// NewI18nController 创建新的i18n控制器
func NewI18nController() *I18nController {
	return &I18nController{}
}

// GetLanguages get supported languages
// @Summary Get supported languages
// @Description Get list of supported languages
// @Tags Internationalization
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=object}
// @Router /api/v1/i18n/languages [get]
func (c *I18nController) GetLanguages(ctx *gin.Context) {
	languages := make([]map[string]interface{}, 0)

	for _, lang := range i18n.GetSupportedLanguages() {
		info := i18n.GetLanguageInfo(lang)
		languages = append(languages, info)
	}

	response.Success(ctx, middleware.T(ctx, "common.success"), gin.H{
		"languages": languages,
		"default":   i18n.DefaultLanguage,
	})
}

// GetTranslations get translations for specified language
// @Summary Get translations
// @Description Get translations for specified language
// @Tags Internationalization
// @Accept json
// @Produce json
// @Param lang path string true "Language code"
// @Success 200 {object} response.APIResponse{data=object}
// @Router /api/v1/i18n/translations/{lang} [get]
func (c *I18nController) GetTranslations(ctx *gin.Context) {
	lang := ctx.Param("lang")
	if lang == "" {
		lang = middleware.GetLanguage(ctx)
	}

	// 返回一些示例翻译用于测试
	translations := map[string]string{
		"user.not_found":           i18n.T("user.not_found", lang),
		"user.created_success":     i18n.T("user.created_success", lang),
		"user.login_success":       i18n.T("user.login_success", lang),
		"auth.unauthorized":        i18n.T("auth.unauthorized", lang),
		"common.success":           i18n.T("common.success", lang),
		"common.validation_failed": i18n.T("common.validation_failed", lang),
		"error.internal_error":     i18n.T("error.internal_error", lang),
	}

	response.Success(ctx, middleware.T(ctx, "common.success"), gin.H{
		"language":     lang,
		"translations": translations,
	})
}

// TestI18n test internationalization functionality
// @Summary Test i18n functionality
// @Description Test internationalization functionality with current language
// @Tags Internationalization
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=object}
// @Router /api/v1/i18n/test [get]
func (c *I18nController) TestI18n(ctx *gin.Context) {
	lang := middleware.GetLanguage(ctx)

	// 测试各种翻译
	testResults := gin.H{
		"current_language": lang,
		"messages": gin.H{
			"welcome":          middleware.T(ctx, "user.login_success"),
			"error":            middleware.T(ctx, "user.not_found"),
			"validation_error": middleware.T(ctx, "common.validation_failed"),
			"success":          middleware.T(ctx, "common.success"),
		},
		"user_flows": gin.H{
			"register_success": middleware.T(ctx, "user.created_success"),
			"login_success":    middleware.T(ctx, "user.login_success"),
			"logout_success":   middleware.T(ctx, "user.logout_success"),
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": middleware.T(ctx, "common.success"),
		"data":    testResults,
	})
}
