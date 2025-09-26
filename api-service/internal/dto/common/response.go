package common

import (
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
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

// BuildResponseWithI18n builds a response with localized message
func BuildResponseWithI18n(ctx *gin.Context, success bool, httpCode errors.HTTPCode, bizCode errors.ErrorCode, data any, errMsg string) {
	ctx.JSON(int(httpCode), Response{
		Success: success,
		Code:    int(bizCode),
		Message: i18n.T(errors.CodeToI18nKey[bizCode], utils.GetUserLangFromRedis(ctx)),
		Data:    data,
		Error:   errMsg,
	})
}

// BadRequest sends a bad request response
func BadRequest(ctx *gin.Context, err error) {
	BuildResponseWithI18n(ctx, false, http.StatusBadRequest, http.StatusBadRequest, nil, err.Error())
}

// SuccessWithData sends a success response
func SuccessWithData(ctx *gin.Context, data any) {
	BuildResponseWithI18n(ctx, true, http.StatusOK, http.StatusOK, data, "")
}

// Success sends a success response with a localized message
func Success(ctx *gin.Context) {
	SuccessWithData(ctx, nil)
}

// DataNotFound sends a not found response
func DataNotFound(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusNotFound, errors.CodeRecordNotFound, nil, "")
}

// Unauthorized sends an unauthorized response
func Unauthorized(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusUnauthorized, errors.CodeInsufficientPermissions, nil, "")
}

// AccessForbidden sends a forbidden response
func AccessForbidden(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusForbidden, errors.CodeResourceAccessDenied, nil, "")
}

// InternalError sends an internal server error response
func InternalError(ctx *gin.Context, err error) {
	BuildResponseWithI18n(ctx, false, http.StatusInternalServerError, errors.CodeInternalError, nil, err.Error())
}

// WithError sends an error response using AppError information
// It extracts the HTTP status code and i18n message key from the error
// and logs the error with appropriate context
func WithError(ctx *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if ok {
		BuildResponseWithI18n(ctx, false, appErr.HTTPStatus, appErr.Code, nil, err.Error())
	} else {
		BuildResponseWithI18n(ctx, false, http.StatusInternalServerError, errors.CodeInternalError, nil, err.Error())
	}
}
