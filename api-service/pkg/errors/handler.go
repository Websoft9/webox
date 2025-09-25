package errors

import (
	"api-service/pkg/i18n"
	"api-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware provides unified error handling middleware for Gin
// It catches panics and converts them into structured error responses
// This middleware ensures consistent error formatting across the application
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Handle different types of recovered values
		switch err := recovered.(type) {
		case string:
			// Handle string panics as internal errors
			HandleError(c, NewAppErrorWithI18n(CodeInternalError, err, "error.unknown_error"))
		case error:
			// Wrap standard errors as internal errors
			HandleError(c, WrapError(err, CodeInternalError, "error.internal_error"))
		default:
			// Handle unknown panic types
			HandleError(c, ErrInternalError)
		}
		c.Abort()
	})
}

// HandleError provides unified error handling for HTTP responses
// It processes both standard Go errors and custom AppErrors,
// applying internationalization when available
func HandleError(c *gin.Context, err error) {
	// Extract user language preference from redis
	lang := utils.GetUserLangFromRedis(c)

	// Check if the error is a custom AppError
	appErr, ok := err.(*AppError)
	if !ok {
		// Handle standard Go errors as internal server errors
		message := i18n.T("error.internal_error", lang)
		sendErrorResponse(c, http.StatusInternalServerError, message, err.Error())
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
	sendErrorResponse(c, appErr.Code, message, appErr.Details)
}

// sendErrorResponse sends a structured error response without importing the response package
// This prevents circular import dependencies
func sendErrorResponse(c *gin.Context, statusCode int, message, details string) {
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": message,
		"error":   details,
	})
}

// IsAppError checks if an error is an AppError and returns it
// This utility function helps with type assertion and error handling
// Returns the AppError and a boolean indicating success
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
