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
	I18nKey    string `json:"-"`       // i18n消息键
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// GetI18nKey 获取i18n键
func (e *AppError) GetI18nKey() string {
	return e.I18nKey
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// NewAppErrorWithI18n 创建带i18n键的应用错误
func NewAppErrorWithI18n(code int, message, i18nKey string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    i18nKey,
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
	ErrInternalError  = NewAppErrorWithI18n(CodeInternalError, CodeMessages[CodeInternalError], "error.internal_error")
	ErrInvalidRequest = NewAppErrorWithI18n(CodeInvalidRequest, CodeMessages[CodeInvalidRequest],
		"common.invalid_request")
	ErrUnauthorized    = NewAppErrorWithI18n(CodeUnauthorized, CodeMessages[CodeUnauthorized], "auth.unauthorized")
	ErrForbidden       = NewAppErrorWithI18n(CodeForbidden, CodeMessages[CodeForbidden], "common.forbidden")
	ErrNotFound        = NewAppErrorWithI18n(CodeNotFound, CodeMessages[CodeNotFound], "common.not_found")
	ErrValidationError = NewAppErrorWithI18n(CodeValidationError, CodeMessages[CodeValidationError],
		"common.validation_failed")

	// 用户相关错误
	ErrUserNotFound      = NewAppErrorWithI18n(CodeUserNotFound, CodeMessages[CodeUserNotFound], "user.not_found")
	ErrUserAlreadyExists = NewAppErrorWithI18n(CodeUserAlreadyExists, CodeMessages[CodeUserAlreadyExists],
		"user.already_exists")
	ErrInvalidCredentials = NewAppErrorWithI18n(CodeInvalidCredentials, CodeMessages[CodeInvalidCredentials],
		"user.invalid_credentials")
	ErrUserInactive    = NewAppErrorWithI18n(CodeUserInactive, CodeMessages[CodeUserInactive], "auth.permission_denied")
	ErrInvalidPassword = NewAppErrorWithI18n(CodeInvalidPassword, CodeMessages[CodeInvalidPassword],
		"user.password_required")
	ErrPasswordTooWeak = NewAppErrorWithI18n(CodePasswordTooWeak, CodeMessages[CodePasswordTooWeak],
		"user.password_too_short")
	ErrEmailAlreadyExists = NewAppErrorWithI18n(CodeEmailAlreadyExists, CodeMessages[CodeEmailAlreadyExists],
		"user.already_exists")
	ErrInvalidEmail     = NewAppErrorWithI18n(CodeInvalidEmail, CodeMessages[CodeInvalidEmail], "user.invalid_email")
	ErrUsernameReserved = NewAppErrorWithI18n(CodeUsernameReserved, CodeMessages[CodeUsernameReserved],
		"user.username_required")
	ErrUserQuotaExceeded = NewAppErrorWithI18n(CodeUserQuotaExceeded, CodeMessages[CodeUserQuotaExceeded],
		"common.forbidden")
)

// Is 检查错误是否匹配
func Is(err, target error) bool {
	if err == target {
		return true
	}

	// 检查是否为AppError类型
	if appErr, ok := err.(*AppError); ok {
		if targetAppErr, ok := target.(*AppError); ok {
			return appErr.Code == targetAppErr.Code
		}
	}

	// 回退到标准库的errors.Is行为
	for {
		if err == target {
			return true
		}
		if x, ok := err.(interface{ Unwrap() error }); ok {
			err = x.Unwrap()
			if err == nil {
				return false
			}
		} else {
			return false
		}
	}
}
