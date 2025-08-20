package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse 标准响应格式
type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Type    string        `json:"type"`
	Code    string        `json:"code"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Value   string `json:"value,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, message, errMsg string) {
	errorType := getErrorType(httpCode)
	errorCode := getErrorCode(httpCode)

	response := StandardResponse{
		Code:    httpCode,
		Message: message,
		Error: &ErrorInfo{
			Type: errorType,
			Code: errorCode,
			Details: []ErrorDetail{
				{
					Message: errMsg,
				},
			},
		},
	}

	c.JSON(httpCode, response)
}

// ValidationError 参数验证错误响应
func ValidationError(c *gin.Context, message string, details []ErrorDetail) {
	response := StandardResponse{
		Code:    http.StatusBadRequest,
		Message: message,
		Error: &ErrorInfo{
			Type:    "VALIDATION_ERROR",
			Code:    "INVALID_PARAMETER",
			Details: details,
		},
	}

	c.JSON(http.StatusBadRequest, response)
}

// getErrorType 根据HTTP状态码获取错误类型
func getErrorType(httpCode int) string {
	switch httpCode {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "AUTHENTICATION_ERROR"
	case http.StatusForbidden:
		return "AUTHORIZATION_ERROR"
	case http.StatusNotFound:
		return "NOT_FOUND_ERROR"
	case http.StatusConflict:
		return "CONFLICT_ERROR"
	case http.StatusUnprocessableEntity:
		return "BUSINESS_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_ERROR"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	case http.StatusBadGateway:
		return "GATEWAY_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "UNKNOWN_ERROR"
	}
}

// getErrorCode 根据HTTP状态码获取错误代码
func getErrorCode(httpCode int) string {
	switch httpCode {
	case http.StatusBadRequest:
		return "INVALID_PARAMETER"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "BUSINESS_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	case http.StatusBadGateway:
		return "GATEWAY_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "UNKNOWN_ERROR"
	}
}

// 保持向后兼容性的别名
func ErrorResponse(c *gin.Context, code int, message, errMsg string) {
	Error(c, code, message, errMsg)
}
