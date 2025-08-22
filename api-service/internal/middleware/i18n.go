package middleware

import (
	"api-service/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware sets up internationalization for requests
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get language from query parameter first
		lang := c.Query("lang")

		// If not provided, try to get from header
		if lang == "" {
			// Get from custom header
			lang = c.GetHeader("X-Language")
		}

		// If still not provided, try Accept-Language header
		if lang == "" {
			acceptLang := c.GetHeader("Accept-Language")
			lang = i18n.DetectLanguageFromHeader(acceptLang)
		}

		// If still empty, use default
		if lang == "" {
			lang = i18n.DefaultLanguage
		}

		// Normalize language to standard format
		lang = i18n.NormalizeLanguage(lang)

		// Store language in context for later use
		c.Set("language", lang)

		// Add language to response header for client reference
		c.Header("Content-Language", lang)

		c.Next()
	}
}

// GetLanguage gets the language from gin context
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return i18n.DefaultLanguage
}

// T is a helper function to translate messages in gin context
func T(c *gin.Context, key string, templateData ...map[string]interface{}) string {
	lang := GetLanguage(c)
	return i18n.T(key, lang, templateData...)
}
