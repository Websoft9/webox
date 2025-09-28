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

// BuildResponseWithI18n builds a standardized API response with internationalized message.
// It takes success status, HTTP status code, error code, response data, and error message,
// then sends a JSON response with localized message based on user's language preference.
func BuildResponseWithI18n(ctx *gin.Context, success bool, httpCode errors.HTTPCode, errorCode errors.ErrorCode, data any, errMsg string) {
	ctx.JSON(int(httpCode), Response{
		Success: success,
		Code:    int(errorCode),
		Message: i18n.T(errors.CodeToI18nKey[errorCode], utils.GetUserLangFromRedis(ctx)),
		Data:    data,
		Error:   errMsg,
	})
}

// BadRequest sends a HTTP 400 Bad Request response with error details.
// Used when client request contains invalid parameters or malformed data.
func BadRequest(ctx *gin.Context, err error) {
	BuildResponseWithI18n(ctx, false, http.StatusBadRequest, http.StatusBadRequest, nil, err.Error())
}

// SuccessWithData sends a HTTP 200 OK response with the provided data payload.
// Used when operation completes successfully and needs to return data to client.
func SuccessWithData(ctx *gin.Context, data any) {
	BuildResponseWithI18n(ctx, true, http.StatusOK, http.StatusOK, data, "")
}

// Success sends a HTTP 200 OK response without data payload.
// Used when operation completes successfully but doesn't need to return specific data.
func Success(ctx *gin.Context) {
	SuccessWithData(ctx, nil)
}

// DataNotFound sends a HTTP 404 Not Found response.
// Used when requested resource or data record doesn't exist in the system.
func DataNotFound(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusNotFound, errors.CodeRecordNotFound, nil, "")
}

// Unauthorized sends a HTTP 401 Unauthorized response.
// Used when client request lacks valid authentication credentials.
func Unauthorized(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusUnauthorized, errors.CodeInsufficientPermissions, nil, "")
}

// AccessForbidden sends a HTTP 403 Forbidden response.
// Used when client has valid credentials but lacks permission to access the resource.
func AccessForbidden(ctx *gin.Context) {
	BuildResponseWithI18n(ctx, false, http.StatusForbidden, errors.CodeResourceAccessDenied, nil, "")
}

// InternalError sends a HTTP 500 Internal Server Error response.
// Used when server encounters an unexpected condition that prevents fulfilling the request.
func InternalError(ctx *gin.Context, err error) {
	BuildResponseWithI18n(ctx, false, http.StatusInternalServerError, errors.CodeInternalError, nil, err.Error())
}

// ServiceUnavailable sends a HTTP 503 Service Unavailable response.
// Used when server is temporarily overloaded or under maintenance.
func ServiceUnavailable(ctx *gin.Context, err error) {
	BuildResponseWithI18n(ctx, false, http.StatusServiceUnavailable, errors.CodeInternalError, nil, err.Error())
}

// WithError sends an error response using AppError information.
// It extracts the HTTP status code and i18n message key from the error,
// falling back to HTTP 500 Internal Server Error for non-AppError types.
func WithError(ctx *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if ok {
		BuildResponseWithI18n(ctx, false, appErr.HTTPStatus, appErr.Code, nil, err.Error())
	} else {
		BuildResponseWithI18n(ctx, false, http.StatusInternalServerError, errors.CodeInternalError, nil, err.Error())
	}
}

// WithErrorCode sends an error response using a specific error code.
// It maps the error code to corresponding HTTP status and sends a localized error message.
func WithErrorCode(ctx *gin.Context, errorCode errors.ErrorCode) {
	BuildResponseWithI18n(ctx, false, errors.CodeToHTTPStatus[errorCode], errorCode, nil, "")
}
