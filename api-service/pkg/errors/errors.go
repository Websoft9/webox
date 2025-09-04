package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error type
// It provides structured error information including business error codes,
// HTTP status codes, and internationalization support
type AppError struct {
	Code       int    `json:"code"`    // Business error code for client identification
	Message    string `json:"message"` // Human-readable error message
	Details    string `json:"details"` // Additional error details for debugging
	HTTPStatus int    `json:"-"`       // HTTP status code for API responses
	I18nKey    string `json:"-"`       // Internationalization key for localized messages
}

// Error implements the error interface
// It returns a formatted string representation of the error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// GetI18nKey returns the internationalization key for this error
// This key can be used to retrieve localized error messages
func (e *AppError) GetI18nKey() string {
	return e.I18nKey
}

// NewAppError creates a new application error with the specified code and message
// The HTTP status code is automatically determined based on the error code
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// NewAppErrorWithI18n creates a new application error with internationalization support
// The i18nKey parameter allows for localized error messages
func NewAppErrorWithI18n(code int, message, i18nKey string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    i18nKey,
	}
}

// NewAppErrorWithDetails creates a new application error with additional details
// The details parameter provides extra context for debugging purposes
func NewAppErrorWithDetails(code int, message, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// WrapError wraps a standard Go error into an AppError
// This is useful for converting system errors into structured application errors
func WrapError(err error, code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    err.Error(),
		HTTPStatus: getHTTPStatusByCode(code),
	}
}

// codeToHTTPStatus maps business error codes to HTTP status codes
// This mapping ensures consistent HTTP responses for different error types
var codeToHTTPStatus = map[int]int{
	CodeSuccess: http.StatusOK,

	// General errors (10000-19999)
	CodeInvalidRequest:  http.StatusBadRequest,
	CodeValidationError: http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,

	// User-related errors (20000-29999)
	CodeUserNotFound:       http.StatusNotFound,
	CodeUserAlreadyExists:  http.StatusConflict,
	CodeEmailAlreadyExists: http.StatusConflict,
	CodeInvalidCredentials: http.StatusUnauthorized,
	CodeInvalidPassword:    http.StatusUnauthorized,
	CodeUserInactive:       http.StatusForbidden,

	// Application-related errors (30000-39999)
	CodeAppNotFound:      http.StatusNotFound,
	CodeAppAlreadyExists: http.StatusConflict,
}

// getHTTPStatusByCode maps business error codes to appropriate HTTP status codes
// It first checks the explicit mapping table, then falls back to range-based defaults
func getHTTPStatusByCode(code int) int {
	if status, exists := codeToHTTPStatus[code]; exists {
		return status
	}
	return getDefaultStatusByCodeRange(code)
}

// getDefaultStatusByCodeRange returns default HTTP status codes based on error code ranges
// This provides a fallback mechanism for unmapped error codes
func getDefaultStatusByCodeRange(code int) int {
	switch {
	case code >= 10000 && code < 20000: // System errors
		return http.StatusInternalServerError
	case code >= 20000 && code < 30000: // User errors
		return http.StatusBadRequest
	case code >= 30000 && code < 40000: // Application errors
		return http.StatusBadRequest
	default: // Unknown error codes
		return http.StatusInternalServerError
	}
}

// Predefined common errors with internationalization support
// These errors can be reused throughout the application for consistency
var (
	// General system errors
	ErrInternalError   = NewAppErrorWithI18n(CodeInternalError, CodeMessages[CodeInternalError], "error.internal_error")
	ErrInvalidRequest  = NewAppErrorWithI18n(CodeInvalidRequest, CodeMessages[CodeInvalidRequest], "common.invalid_request")
	ErrUnauthorized    = NewAppErrorWithI18n(CodeUnauthorized, CodeMessages[CodeUnauthorized], "auth.unauthorized")
	ErrForbidden       = NewAppErrorWithI18n(CodeForbidden, CodeMessages[CodeForbidden], "common.forbidden")
	ErrNotFound        = NewAppErrorWithI18n(CodeNotFound, CodeMessages[CodeNotFound], "common.not_found")
	ErrValidationError = NewAppErrorWithI18n(CodeValidationError, CodeMessages[CodeValidationError], "common.validation_failed")

	// User-related errors with specific business logic
	ErrUserNotFound       = NewAppErrorWithI18n(CodeUserNotFound, CodeMessages[CodeUserNotFound], "user.not_found")
	ErrUserAlreadyExists  = NewAppErrorWithI18n(CodeUserAlreadyExists, CodeMessages[CodeUserAlreadyExists], "user.already_exists")
	ErrInvalidCredentials = NewAppErrorWithI18n(CodeInvalidCredentials, CodeMessages[CodeInvalidCredentials], "user.invalid_credentials")
	ErrUserInactive       = NewAppErrorWithI18n(CodeUserInactive, CodeMessages[CodeUserInactive], "auth.permission_denied")
	ErrInvalidPassword    = NewAppErrorWithI18n(CodeInvalidPassword, CodeMessages[CodeInvalidPassword], "user.password_required")
	ErrPasswordTooWeak    = NewAppErrorWithI18n(CodePasswordTooWeak, CodeMessages[CodePasswordTooWeak], "user.password_too_short")
	ErrEmailAlreadyExists = NewAppErrorWithI18n(CodeEmailAlreadyExists, CodeMessages[CodeEmailAlreadyExists], "user.already_exists")
	ErrInvalidEmail       = NewAppErrorWithI18n(CodeInvalidEmail, CodeMessages[CodeInvalidEmail], "user.invalid_email")
	ErrUsernameReserved   = NewAppErrorWithI18n(CodeUsernameReserved, CodeMessages[CodeUsernameReserved], "user.username_required")
	ErrUserQuotaExceeded  = NewAppErrorWithI18n(CodeUserQuotaExceeded, CodeMessages[CodeUserQuotaExceeded], "common.forbidden")
	ErrEmailNotVerified   = NewAppErrorWithI18n(CodeEmailNotVerified, CodeMessages[CodeEmailNotVerified], "user.email_not_verified")

	// Authentication and token-related errors
	ErrInvalidToken     = NewAppErrorWithI18n(CodeInvalidToken, CodeMessages[CodeInvalidToken], "auth.invalid_token")
	ErrTokenExpired     = NewAppErrorWithI18n(CodeTokenExpired, CodeMessages[CodeTokenExpired], "auth.token_expired")
	ErrTokenAlreadyUsed = NewAppErrorWithI18n(CodeTokenAlreadyUsed, CodeMessages[CodeTokenAlreadyUsed], "auth.token_already_used")
)

// Is checks if an error matches a target error
// It provides enhanced error matching for AppError types by comparing error codes
// and falls back to standard error unwrapping for other error types
func Is(err, target error) bool {
	if err == target {
		return true
	}

	// Check if both errors are AppError types and compare their codes
	if appErr, ok := err.(*AppError); ok {
		if targetAppErr, ok := target.(*AppError); ok {
			return appErr.Code == targetAppErr.Code
		}
	}

	// Fall back to standard library errors.Is behavior with unwrapping
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
