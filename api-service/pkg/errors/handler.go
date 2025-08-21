package errors

import (
	"api-service/pkg/i18n"
	"api-service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		switch err := recovered.(type) {
		case string:
			HandleError(c, NewAppErrorWithI18n(CodeInternalError, err, "error.unknown_error"))
		case error:
			HandleError(c, WrapError(err, CodeInternalError, "error.internal_error"))
		default:
			HandleError(c, ErrInternalError)
		}
		c.Abort()
	})
}

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, err error) {
	// 获取请求语言
	lang := getLanguageFromContext(c)

	appErr, ok := err.(*AppError)
	if !ok {
		// 标准错误
		message := i18n.T("error.internal_error", lang)
		response.Error(c, http.StatusInternalServerError, message, err.Error())
		return
	}

	// 自定义应用错误
	message := appErr.Message

	// 如果有i18n键，使用翻译后的消息
	if appErr.I18nKey != "" {
		translatedMsg := i18n.T(appErr.I18nKey, lang)
		if translatedMsg != appErr.I18nKey { // 翻译成功
			message = translatedMsg
		}
	}

	response.Error(c, appErr.HTTPStatus, message, appErr.Details)
} // getLanguageFromContext 从gin上下文获取语言
func getLanguageFromContext(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return i18n.DefaultLanguage
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
