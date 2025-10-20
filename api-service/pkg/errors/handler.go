package errors

import (
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware provides unified error handling middleware for Gin
// It catches panics and converts them into structured error responses
// This middleware ensures consistent error formatting across the application
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Handle different types of recovered values
		switch err := recovered.(type) {
		case string:
			// Handle string panics as internal errors
			log.ErrorContext(c, err)
			processes(c, NewAppErrorWithI18nDetails(CodeInternalError, "error.unknown_error", err))
		case error:
			// Wrap standard errors as internal errors
			log.ErrorContext(c, "system panic", logger.ErrorField(err))
			processes(c, NewAppErrorWrapError(err, CodeInternalError))
		default:
			// Handle unknown panic types
			log.ErrorContext(c, "system panic", logger.Any("error", recovered))
			processes(c, ErrInternalError)
		}
		c.Abort()
	})
}

// Processes provides unified error handling for HTTP responses
// It processes both standard Go errors and custom AppErrors,
// applying internationalization when available
func processes(c *gin.Context, err error) {
	// Extract user language preference from redis
	lang := utils.GetUserLangFromRedis(c)

	// Check if the error is a custom AppError
	appErr, ok := err.(*AppError)
	if !ok {
		// Handle standard Go errors as internal server errors
		message := i18n.T("error.internal_error", lang)
		sendErrorResponse(c, http.StatusInternalServerError, CodeInternalError, message, utils.FormatErrorWithStack(err))
		return
	}

	// Process custom application errors
	message := appErr.Message

	// Apply internationalization if i18n key is available
	if appErr.I18nKey != "" {
		translatedMsg := i18n.T(appErr.I18nKey, lang)
		if translatedMsg != appErr.I18nKey { // Translation was successful
			message = translatedMsg
		}
	}

	// Send structured error response
	sendErrorResponse(c, appErr.HTTPStatus, appErr.Code, message, appErr.Details)
}

// sendErrorResponse sends a structured error response without importing the response package
// This prevents circular import dependencies
func sendErrorResponse(c *gin.Context, statusCode HTTPCode, errorCode ErrorCode, message, details string) {
	c.JSON(int(statusCode), gin.H{
		"code":    int(errorCode),
		"message": message,
		"error":   details,
	})
}
