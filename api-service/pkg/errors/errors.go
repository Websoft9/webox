package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用程序自定义错误类型
type AppError struct {
	Code       int    `json:"code"`    // 业务错误码
	Message    string `json:"message"` // 错误消息
	Details    string `json:"details"` // 详细信息
	HTTPStatus int    `json:"-"`       // HTTP状态码
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// NewAppErrorWithDetails 创建带详细信息的应用错误
func NewAppErrorWithDetails(code int, message, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// WrapError 包装标准error为AppError
func WrapError(err error, code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    err.Error(),
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// HTTP状态码映射表
var codeToHTTPStatus = map[int]int{
	CodeSuccess: http.StatusOK,
	// 通用错误
	CodeInvalidRequest:  http.StatusBadRequest,
	CodeValidationError: http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	// 用户相关错误
	CodeUserNotFound:       http.StatusNotFound,
	CodeUserAlreadyExists:  http.StatusConflict,
	CodeEmailAlreadyExists: http.StatusConflict,
	CodeInvalidCredentials: http.StatusUnauthorized,
	CodeInvalidPassword:    http.StatusUnauthorized,
	CodeUserInactive:       http.StatusForbidden,
	// 应用相关错误
	CodeAppNotFound:      http.StatusNotFound,
	CodeAppAlreadyExists: http.StatusConflict,
}

// 根据业务错误码映射HTTP状态码
func getHTTPStatusByCode(code int) int {
	if status, exists := codeToHTTPStatus[code]; exists {
		return status
	}
	return getDefaultStatusByCodeRange(code)
}

// 根据错误码范围获取默认HTTP状态码
func getDefaultStatusByCodeRange(code int) int {
	switch {
	case code >= 10000 && code < 20000:
		return http.StatusInternalServerError
	case code >= 20000 && code < 30000:
		return http.StatusBadRequest
	case code >= 30000 && code < 40000:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// 预定义的常见错误
var (
	ErrInternalError   = NewAppError(CodeInternalError, CodeMessages[CodeInternalError])
	ErrInvalidRequest  = NewAppError(CodeInvalidRequest, CodeMessages[CodeInvalidRequest])
	ErrUnauthorized    = NewAppError(CodeUnauthorized, CodeMessages[CodeUnauthorized])
	ErrForbidden       = NewAppError(CodeForbidden, CodeMessages[CodeForbidden])
	ErrNotFound        = NewAppError(CodeNotFound, CodeMessages[CodeNotFound])
	ErrValidationError = NewAppError(CodeValidationError, CodeMessages[CodeValidationError])

	// 用户相关错误
	ErrUserNotFound       = NewAppError(CodeUserNotFound, CodeMessages[CodeUserNotFound])
	ErrUserAlreadyExists  = NewAppError(CodeUserAlreadyExists, CodeMessages[CodeUserAlreadyExists])
	ErrInvalidCredentials = NewAppError(CodeInvalidCredentials, CodeMessages[CodeInvalidCredentials])
	ErrUserInactive       = NewAppError(CodeUserInactive, CodeMessages[CodeUserInactive])
	ErrInvalidPassword    = NewAppError(CodeInvalidPassword, CodeMessages[CodeInvalidPassword])
	ErrPasswordTooWeak    = NewAppError(CodePasswordTooWeak, CodeMessages[CodePasswordTooWeak])
	ErrEmailAlreadyExists = NewAppError(CodeEmailAlreadyExists, CodeMessages[CodeEmailAlreadyExists])
	ErrInvalidEmail       = NewAppError(CodeInvalidEmail, CodeMessages[CodeInvalidEmail])
	ErrUsernameReserved   = NewAppError(CodeUsernameReserved, CodeMessages[CodeUsernameReserved])
	ErrUserQuotaExceeded  = NewAppError(CodeUserQuotaExceeded, CodeMessages[CodeUserQuotaExceeded])
)
