package response

import (
	"api-service/internal/constants"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response standard API response structure
type Response struct {
	Success bool   `json:"success" example:"true"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty" example:"Error message"`
}

// APIResponse swagger response wrapper for documentation
type APIResponse struct {
	Success bool   `json:"success" example:"true"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty" example:"Error message"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 修复参数名冲突问题，避免与内置 error 类型冲突
func ErrorResponse(c *gin.Context, code int, message, errMsg string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   errMsg,
	})
}

// Error 保持向后兼容性的别名
func Error(c *gin.Context, code int, message, errMsg string) {
	ErrorResponse(c, code, message, errMsg)
}

// BuildResponseWithI18n builds a response with localized message
func BuildResponseWithI18n(ctx *gin.Context, success bool, httpStatusCode, bizCode int, data any, messageKey, errMsg string) {
	ctx.JSON(httpStatusCode, Response{
		Success: success,
		Code:    bizCode,
		Message: i18n.T(messageKey, utils.GetUserLangFromRedis(ctx)),
		Data:    data,
		Error:   errMsg,
	})
}

// BadRequest sends a bad request response
func BadRequest(ctx *gin.Context, err error, messageKey string) {
	BuildResponseWithI18n(ctx, false, http.StatusBadRequest, http.StatusBadRequest, nil, messageKey, err.Error())
}

// OKWithData sends a success response
func OKWithData(ctx *gin.Context, data any, messageKey string) {
	BuildResponseWithI18n(ctx, true, http.StatusOK, http.StatusOK, data, messageKey, "")
}

// OK sends a success response with a localized message
func OK(ctx *gin.Context, messageKey string) {
	OKWithData(ctx, nil, messageKey)
}

// DataNotFound sends a not found response
func DataNotFound(ctx *gin.Context, messageKey string) {
	BuildResponseWithI18n(ctx, false, http.StatusNotFound, errors.CodeRecordNotFound, nil, messageKey, "")
}

// Unauthorized sends an unauthorized response
func Unauthorized(ctx *gin.Context, messageKey string) {
	BuildResponseWithI18n(ctx, false, http.StatusUnauthorized, errors.CodeInsufficientPermissions, nil, messageKey, "")
}

// AccessForbidden sends a forbidden response
func AccessForbidden(ctx *gin.Context, messageKey string) {
	BuildResponseWithI18n(ctx, false, http.StatusForbidden, errors.CodeResourceAccessDenied, nil, messageKey, "")
}

// InternalError sends an internal server error response
func InternalError(ctx *gin.Context, err error, messageKey string, log logger.Logger) {
	log.ErrorContext(ctx.Request.Context(), "Internal server error", logger.ErrorField(err))
	BuildResponseWithI18n(ctx, false, http.StatusInternalServerError, errors.CodeInternalError, nil, messageKey, err.Error())
}

// WithError sends an error response using AppError information
// It extracts the HTTP status code and i18n message key from the error
// and logs the error with appropriate context
func WithError(ctx *gin.Context, err error, log logger.Logger) {
	httpStatus, bizCode, messageKey := errors.GetErrCodeAndMessageKey(err)
	log.ErrorContext(ctx.Request.Context(), i18n.T(messageKey, constants.DefaultLanguage), logger.ErrorField(err))
	BuildResponseWithI18n(ctx, false, httpStatus, bizCode, nil, messageKey, err.Error())
}
