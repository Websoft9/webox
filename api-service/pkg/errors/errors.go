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

// NewAppErrorWithMessage creates a new application error with automatic i18n key resolution
// The i18nKey is automatically determined based on the error code
func NewAppErrorWithMessage(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    getI18nKeyByCode(code),
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

	// Authentication related errors (1000-1999)
	CodeInvalidCredentials:      http.StatusUnauthorized,
	CodeTokenExpired:            http.StatusUnauthorized,
	CodeInvalidToken:            http.StatusUnauthorized,
	CodeAccountDisabled:         http.StatusForbidden,
	CodeAccountLocked:           http.StatusForbidden,
	CodePasswordTooWeak:         http.StatusBadRequest,
	CodeInvalidVerificationCode: http.StatusBadRequest,
	CodeLoginAttemptsExceeded:   http.StatusTooManyRequests,
	CodeEmailAlreadyExists:      http.StatusBadRequest,
	CodeTokenAlreadyUsed:        http.StatusUnauthorized,

	// Permission related errors (2000-2999)
	CodeInsufficientPermissions:       http.StatusForbidden,
	CodeResourceAccessDenied:          http.StatusForbidden,
	CodeOperationPermissionDenied:     http.StatusForbidden,
	CodeRolePermissionDenied:          http.StatusForbidden,
	CodeResourceGroupPermissionDenied: http.StatusForbidden,

	// Parameter validation errors (3000-3999)
	CodeValidationFailed:         http.StatusBadRequest,
	CodeRequiredParameterMissing: http.StatusBadRequest,
	CodeInvalidParameterFormat:   http.StatusBadRequest,
	CodeParameterOutOfRange:      http.StatusBadRequest,
	CodeInvalidParameterLength:   http.StatusBadRequest,
	CodeInvalidEmailFormat:       http.StatusBadRequest,
	CodeInvalidPhoneFormat:       http.StatusBadRequest,
	CodeInvalidURLFormat:         http.StatusBadRequest,
	CodeInvalidDateFormat:        http.StatusBadRequest,
	CodeEmailNotVerified:         http.StatusBadRequest,

	// Resource related errors (4000-4999)
	CodeRecordNotFound:             http.StatusNotFound,
	CodeResourceNotFound:           http.StatusNotFound,
	CodeResourceAlreadyExists:      http.StatusConflict,
	CodeResourceStateNotAllowed:    http.StatusConflict,
	CodeResourceDependencyConflict: http.StatusConflict,
	CodeResourceQuotaInsufficient:  http.StatusConflict,
	CodeResourceInUse:              http.StatusConflict,
	CodeRecordQueryFailed:          http.StatusConflict,
	CodeRecordCreateFailed:         http.StatusConflict,
	CodeRecordUpdateFailed:         http.StatusConflict,
	CodeRecordDeleteFailed:         http.StatusConflict,
	CodeRecordIsDisabled:           http.StatusConflict,
	CodeRecordNoAffected:           http.StatusConflict,
	CodeRecordDeleteDenied:         http.StatusConflict,

	// Business logic errors (5000-5999)
	CodeServerOffline:                  http.StatusServiceUnavailable,
	CodeAppDeploymentFailed:            http.StatusUnprocessableEntity,
	CodeWorkflowExecutionFailed:        http.StatusUnprocessableEntity,
	CodeCertificateRequestFailed:       http.StatusUnprocessableEntity,
	CodeBackupOperationFailed:          http.StatusUnprocessableEntity,
	CodeMonitoringDataCollectionFailed: http.StatusServiceUnavailable,
	CodeAppPublishFailed:               http.StatusUnprocessableEntity,
	CodeAppOfflineFailed:               http.StatusUnprocessableEntity,
	CodeHealthCheckFailed:              http.StatusServiceUnavailable,
	CodeGatewayConfigUpdateFailed:      http.StatusUnprocessableEntity,
	CodeUserAlreadyExists:              http.StatusConflict,
	CodePermissionInvalid:              http.StatusUnprocessableEntity,

	// System related errors (6000-6999)
	CodeInternalError:                http.StatusInternalServerError,
	CodeDatabaseConnectionFailed:     http.StatusInternalServerError,
	CodeCacheServiceUnavailable:      http.StatusServiceUnavailable,
	CodeFilesystemError:              http.StatusInternalServerError,
	CodeNetworkTimeout:               http.StatusRequestTimeout,
	CodeThirdPartyServiceUnavailable: http.StatusServiceUnavailable,
	CodeSystemMaintenance:            http.StatusServiceUnavailable,
}

// getHTTPStatusByCode maps business error codes to appropriate HTTP status codes
// It first checks the explicit mapping table, then falls back to range-based defaults
func getHTTPStatusByCode(code int) int {
	if status, exists := codeToHTTPStatus[code]; exists {
		return status
	}
	return getDefaultStatusByCodeRange(code)
}

// getI18nKeyByCode retrieves the internationalization key for a given error code
// If no specific mapping exists, it returns a default unknown error key
func getI18nKeyByCode(code int) string {
	if key, exists := CodeToI18nKey[code]; exists {
		return key
	}
	return "system.unknown_error"
}

// getDefaultStatusByCodeRange returns default HTTP status codes based on error code ranges
// This provides a fallback mechanism for unmapped error codes
func getDefaultStatusByCodeRange(code int) int {
	switch {
	case code >= 1000 && code < 2000: // Authentication errors
		return http.StatusUnauthorized
	case code >= 2000 && code < 3000: // Permission errors
		return http.StatusForbidden
	case code >= 3000 && code < 4000: // Validation errors
		return http.StatusBadRequest
	case code >= 4000 && code < 5000: // Resource errors
		return http.StatusNotFound
	case code >= 5000 && code < 6000: // Business logic errors
		return http.StatusUnprocessableEntity
	case code >= 6000 && code < 7000: // System errors
		return http.StatusInternalServerError
	default: // Unknown error codes
		return http.StatusInternalServerError
	}
}

// Predefined common errors with internationalization support
// These errors can be reused throughout the application for consistency
var (
	// Authentication related errors (1000-1999)
	ErrInvalidCredentials      = NewAppErrorWithI18n(CodeInvalidCredentials, CodeMessages[CodeInvalidCredentials], CodeToI18nKey[CodeInvalidCredentials])
	ErrTokenExpired            = NewAppErrorWithI18n(CodeTokenExpired, CodeMessages[CodeTokenExpired], CodeToI18nKey[CodeTokenExpired])
	ErrInvalidToken            = NewAppErrorWithI18n(CodeInvalidToken, CodeMessages[CodeInvalidToken], CodeToI18nKey[CodeInvalidToken])
	ErrAccountDisabled         = NewAppErrorWithI18n(CodeAccountDisabled, CodeMessages[CodeAccountDisabled], CodeToI18nKey[CodeAccountDisabled])
	ErrAccountLocked           = NewAppErrorWithI18n(CodeAccountLocked, CodeMessages[CodeAccountLocked], CodeToI18nKey[CodeAccountLocked])
	ErrPasswordTooWeak         = NewAppErrorWithI18n(CodePasswordTooWeak, CodeMessages[CodePasswordTooWeak], CodeToI18nKey[CodePasswordTooWeak])
	ErrInvalidVerificationCode = NewAppErrorWithI18n(CodeInvalidVerificationCode, CodeMessages[CodeInvalidVerificationCode], CodeToI18nKey[CodeInvalidVerificationCode])
	ErrLoginAttemptsExceeded   = NewAppErrorWithI18n(CodeLoginAttemptsExceeded, CodeMessages[CodeLoginAttemptsExceeded], CodeToI18nKey[CodeLoginAttemptsExceeded])
	ErrEmailAlreadyExists      = NewAppErrorWithI18n(CodeEmailAlreadyExists, CodeMessages[CodeEmailAlreadyExists], CodeToI18nKey[CodeEmailAlreadyExists])
	ErrTokenAlreadyUsed        = NewAppErrorWithI18n(CodeTokenAlreadyUsed, CodeMessages[CodeTokenAlreadyUsed], CodeToI18nKey[CodeTokenAlreadyUsed])

	// Permission related errors (2000-2999)
	ErrInsufficientPermissions       = NewAppErrorWithI18n(CodeInsufficientPermissions, CodeMessages[CodeInsufficientPermissions], CodeToI18nKey[CodeInsufficientPermissions])
	ErrResourceAccessDenied          = NewAppErrorWithI18n(CodeResourceAccessDenied, CodeMessages[CodeResourceAccessDenied], CodeToI18nKey[CodeResourceAccessDenied])
	ErrOperationPermissionDenied     = NewAppErrorWithI18n(CodeOperationPermissionDenied, CodeMessages[CodeOperationPermissionDenied], CodeToI18nKey[CodeOperationPermissionDenied])
	ErrRolePermissionDenied          = NewAppErrorWithI18n(CodeRolePermissionDenied, CodeMessages[CodeRolePermissionDenied], CodeToI18nKey[CodeRolePermissionDenied])
	ErrResourceGroupPermissionDenied = NewAppErrorWithI18n(
		CodeResourceGroupPermissionDenied,
		CodeMessages[CodeResourceGroupPermissionDenied],
		CodeToI18nKey[CodeResourceGroupPermissionDenied],
	)

	// Parameter validation errors (3000-3999)
	ErrValidationFailed         = NewAppErrorWithI18n(CodeValidationFailed, CodeMessages[CodeValidationFailed], CodeToI18nKey[CodeValidationFailed])
	ErrRequiredParameterMissing = NewAppErrorWithI18n(CodeRequiredParameterMissing, CodeMessages[CodeRequiredParameterMissing], CodeToI18nKey[CodeRequiredParameterMissing])
	ErrInvalidParameterFormat   = NewAppErrorWithI18n(CodeInvalidParameterFormat, CodeMessages[CodeInvalidParameterFormat], CodeToI18nKey[CodeInvalidParameterFormat])
	ErrParameterOutOfRange      = NewAppErrorWithI18n(CodeParameterOutOfRange, CodeMessages[CodeParameterOutOfRange], CodeToI18nKey[CodeParameterOutOfRange])
	ErrInvalidParameterLength   = NewAppErrorWithI18n(CodeInvalidParameterLength, CodeMessages[CodeInvalidParameterLength], CodeToI18nKey[CodeInvalidParameterLength])
	ErrInvalidEmailFormat       = NewAppErrorWithI18n(CodeInvalidEmailFormat, CodeMessages[CodeInvalidEmailFormat], CodeToI18nKey[CodeInvalidEmailFormat])
	ErrInvalidPhoneFormat       = NewAppErrorWithI18n(CodeInvalidPhoneFormat, CodeMessages[CodeInvalidPhoneFormat], CodeToI18nKey[CodeInvalidPhoneFormat])
	ErrInvalidURLFormat         = NewAppErrorWithI18n(CodeInvalidURLFormat, CodeMessages[CodeInvalidURLFormat], CodeToI18nKey[CodeInvalidURLFormat])
	ErrInvalidDateFormat        = NewAppErrorWithI18n(CodeInvalidDateFormat, CodeMessages[CodeInvalidDateFormat], CodeToI18nKey[CodeInvalidDateFormat])
	ErrEmailNotVerified         = NewAppErrorWithI18n(CodeEmailNotVerified, CodeMessages[CodeEmailNotVerified], CodeToI18nKey[CodeEmailNotVerified])

	// Resource related errors (4000-4999)
	ErrRecordNotFound             = NewAppErrorWithI18n(CodeRecordNotFound, CodeMessages[CodeRecordNotFound], CodeToI18nKey[CodeRecordNotFound])
	ErrResourceNotFound           = NewAppErrorWithI18n(CodeResourceNotFound, CodeMessages[CodeResourceNotFound], CodeToI18nKey[CodeResourceNotFound])
	ErrResourceAlreadyExists      = NewAppErrorWithI18n(CodeResourceAlreadyExists, CodeMessages[CodeResourceAlreadyExists], CodeToI18nKey[CodeResourceAlreadyExists])
	ErrResourceStateNotAllowed    = NewAppErrorWithI18n(CodeResourceStateNotAllowed, CodeMessages[CodeResourceStateNotAllowed], CodeToI18nKey[CodeResourceStateNotAllowed])
	ErrResourceDependencyConflict = NewAppErrorWithI18n(CodeResourceDependencyConflict, CodeMessages[CodeResourceDependencyConflict], CodeToI18nKey[CodeResourceDependencyConflict])
	ErrResourceQuotaInsufficient  = NewAppErrorWithI18n(CodeResourceQuotaInsufficient, CodeMessages[CodeResourceQuotaInsufficient], CodeToI18nKey[CodeResourceQuotaInsufficient])
	ErrResourceInUse              = NewAppErrorWithI18n(CodeResourceInUse, CodeMessages[CodeResourceInUse], CodeToI18nKey[CodeResourceInUse])
	ErrRecordQueryFailed          = NewAppErrorWithI18n(CodeRecordQueryFailed, CodeMessages[CodeRecordQueryFailed], CodeToI18nKey[CodeRecordQueryFailed])
	ErrRecordCreateFailed         = NewAppErrorWithI18n(CodeRecordCreateFailed, CodeMessages[CodeRecordCreateFailed], CodeToI18nKey[CodeRecordCreateFailed])
	ErrRecordUpdateFailed         = NewAppErrorWithI18n(CodeRecordUpdateFailed, CodeMessages[CodeRecordUpdateFailed], CodeToI18nKey[CodeRecordUpdateFailed])
	ErrRecordDeleteFailed         = NewAppErrorWithI18n(CodeRecordDeleteFailed, CodeMessages[CodeRecordDeleteFailed], CodeToI18nKey[CodeRecordDeleteFailed])
	ErrRecordIsDisabled           = NewAppErrorWithI18n(CodeRecordIsDisabled, CodeMessages[CodeRecordIsDisabled], CodeToI18nKey[CodeRecordIsDisabled])
	ErrRecordNoAffected           = NewAppErrorWithI18n(CodeRecordNoAffected, CodeMessages[CodeRecordNoAffected], CodeToI18nKey[CodeRecordNoAffected])
	ErrRecordDeleteDenied         = NewAppErrorWithI18n(CodeRecordDeleteDenied, CodeMessages[CodeRecordDeleteDenied], CodeToI18nKey[CodeRecordDeleteDenied])

	// Business logic errors (5000-5999)
	ErrServerOffline                  = NewAppErrorWithI18n(CodeServerOffline, CodeMessages[CodeServerOffline], CodeToI18nKey[CodeServerOffline])
	ErrAppDeploymentFailed            = NewAppErrorWithI18n(CodeAppDeploymentFailed, CodeMessages[CodeAppDeploymentFailed], CodeToI18nKey[CodeAppDeploymentFailed])
	ErrWorkflowExecutionFailed        = NewAppErrorWithI18n(CodeWorkflowExecutionFailed, CodeMessages[CodeWorkflowExecutionFailed], CodeToI18nKey[CodeWorkflowExecutionFailed])
	ErrCertificateRequestFailed       = NewAppErrorWithI18n(CodeCertificateRequestFailed, CodeMessages[CodeCertificateRequestFailed], CodeToI18nKey[CodeCertificateRequestFailed])
	ErrBackupOperationFailed          = NewAppErrorWithI18n(CodeBackupOperationFailed, CodeMessages[CodeBackupOperationFailed], CodeToI18nKey[CodeBackupOperationFailed])
	ErrMonitoringDataCollectionFailed = NewAppErrorWithI18n(
		CodeMonitoringDataCollectionFailed,
		CodeMessages[CodeMonitoringDataCollectionFailed],
		CodeToI18nKey[CodeMonitoringDataCollectionFailed],
	)
	ErrAppPublishFailed          = NewAppErrorWithI18n(CodeAppPublishFailed, CodeMessages[CodeAppPublishFailed], CodeToI18nKey[CodeAppPublishFailed])
	ErrAppOfflineFailed          = NewAppErrorWithI18n(CodeAppOfflineFailed, CodeMessages[CodeAppOfflineFailed], CodeToI18nKey[CodeAppOfflineFailed])
	ErrHealthCheckFailed         = NewAppErrorWithI18n(CodeHealthCheckFailed, CodeMessages[CodeHealthCheckFailed], CodeToI18nKey[CodeHealthCheckFailed])
	ErrGatewayConfigUpdateFailed = NewAppErrorWithI18n(CodeGatewayConfigUpdateFailed, CodeMessages[CodeGatewayConfigUpdateFailed], CodeToI18nKey[CodeGatewayConfigUpdateFailed])
	ErrUserAlreadyExists         = NewAppErrorWithI18n(CodeUserAlreadyExists, CodeMessages[CodeUserAlreadyExists], CodeToI18nKey[CodeUserAlreadyExists])
	ErrPermissionInvalid         = NewAppErrorWithI18n(CodePermissionInvalid, CodeMessages[CodePermissionInvalid], CodeToI18nKey[CodePermissionInvalid])

	// System related errors (6000-6999)
	ErrInternalError                = NewAppErrorWithI18n(CodeInternalError, CodeMessages[CodeInternalError], CodeToI18nKey[CodeInternalError])
	ErrCacheServiceUnavailable      = NewAppErrorWithI18n(CodeCacheServiceUnavailable, CodeMessages[CodeCacheServiceUnavailable], CodeToI18nKey[CodeCacheServiceUnavailable])
	ErrFilesystemError              = NewAppErrorWithI18n(CodeFilesystemError, CodeMessages[CodeFilesystemError], CodeToI18nKey[CodeFilesystemError])
	ErrNetworkTimeout               = NewAppErrorWithI18n(CodeNetworkTimeout, CodeMessages[CodeNetworkTimeout], CodeToI18nKey[CodeNetworkTimeout])
	ErrThirdPartyServiceUnavailable = NewAppErrorWithI18n(
		CodeThirdPartyServiceUnavailable,
		CodeMessages[CodeThirdPartyServiceUnavailable],
		CodeToI18nKey[CodeThirdPartyServiceUnavailable],
	)
	ErrSystemMaintenance        = NewAppErrorWithI18n(CodeSystemMaintenance, CodeMessages[CodeSystemMaintenance], CodeToI18nKey[CodeSystemMaintenance])
	ErrDatabaseConnectionFailed = NewAppErrorWithI18n(CodeDatabaseConnectionFailed, CodeMessages[CodeDatabaseConnectionFailed], CodeToI18nKey[CodeDatabaseConnectionFailed])
)

// GetErrCodeAndMessageKey extracts HTTP status code and i18n message key from an error
// This function is used by the controller helpers to provide consistent error responses
func GetErrCodeAndMessageKey(err error) (statusCode int, messageKey string) {
	if appErr, ok := err.(*AppError); ok {
		// For AppError, return the HTTP status and i18n key
		return appErr.HTTPStatus, appErr.I18nKey
	}

	// For standard errors, return default values
	return http.StatusInternalServerError, "system.internal_error"
}

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
