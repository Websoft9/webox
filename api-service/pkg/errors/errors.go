package errors

import (
	"api-service/pkg/i18n"
	"api-service/pkg/utils"
	"net/http"
)

// AppError represents a custom application error type
// It provides structured error information including business error codes,
// HTTP status codes, and internationalization support
type AppError struct {
	Code       ErrorCode `json:"code"`    // Business error code for client identification
	Message    string    `json:"message"` // Human-readable error message
	Details    string    `json:"details"` // Additional error details for debugging
	HTTPStatus HTTPCode  `json:"-"`       // HTTP status code for API responses
	I18nKey    string    `json:"-"`       // Internationalization key for localized messages
}

// Error implements the error interface
// It returns a formatted string representation of the error
func (e *AppError) Error() string {
	// Use i18n key or message
	message := e.Message
	if message == "" && e.I18nKey != "" {
		// Try to translate using default language if i18n is initialized
		if i18n.Bundle == nil {
			_ = i18n.Init()
		}
		message = i18n.T(e.I18nKey, i18n.DefaultLanguage)
	}

	if e.Details != "" {
		return e.Details
	}
	return message
}

// GetI18nKey returns the internationalization key for this error
// This key can be used to retrieve localized error messages
func (e *AppError) GetI18nKey() string {
	return e.I18nKey
}

// NewAppErrorWithI18n creates a new application error with internationalization support
// The i18nKey parameter allows for localized error messages
func NewAppErrorWithI18n(code ErrorCode, i18nKey string) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    i18nKey,
	}
}

// NewAppErrorWithI18nDetails creates a new application error with internationalization support and additional details
func NewAppErrorWithI18nDetails(code ErrorCode, i18nKey, details string) *AppError {
	return &AppError{
		Code:       code,
		Details:    details,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    i18nKey,
	}
}

// NewAppError creates a new application error with the specified code and message
// The HTTP status code is automatically determined based on the error code
func NewAppError(code ErrorCode) *AppError {
	return NewAppErrorWithI18n(code, getI18nKeyByCode(code))
}

// NewAppErrorWithDetails creates a new application error with additional details
// The details parameter provides extra context for debugging purposes
func NewAppErrorWithDetails(code ErrorCode, details string) *AppError {
	return &AppError{
		Code:       code,
		Details:    details,
		HTTPStatus: getHTTPStatusByCode(code),
		I18nKey:    getI18nKeyByCode(code),
	}
}

// WrapError wraps a standard Go error into an AppError
// This is useful for converting system errors into structured application errors
func NewAppErrorWrapError(err error, code ErrorCode) *AppError {
	return NewAppErrorWithDetails(code, utils.FormatErrorWithStack(err))
}

// getHTTPStatusByCode maps business error codes to appropriate HTTP status codes
// It first checks the explicit mapping table, then falls back to range-based defaults
func getHTTPStatusByCode(code ErrorCode) HTTPCode {
	if status, exists := CodeToHTTPStatus[code]; exists {
		return status
	}
	return getDefaultStatusByCodeRange(code)
}

// getI18nKeyByCode retrieves the internationalization key for a given error code
// If no specific mapping exists, it returns a default unknown error key
func getI18nKeyByCode(code ErrorCode) string {
	if key, exists := CodeToI18nKey[code]; exists {
		return key
	}
	return "system.unknown_error"
}

// getDefaultStatusByCodeRange returns default HTTP status codes based on error code ranges
// This provides a fallback mechanism for unmapped error codes
func getDefaultStatusByCodeRange(code ErrorCode) HTTPCode {
	switch {
	case code == http.StatusOK: // Default case
		return http.StatusOK
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

// Predefined common errors with internationalization support
// These errors can be reused throughout the application for consistency
var (
	// Authentication related errors (1000-1999)
	ErrInvalidCredentials      = NewAppErrorWithI18n(CodeInvalidCredentials, CodeToI18nKey[CodeInvalidCredentials])
	ErrTokenExpired            = NewAppErrorWithI18n(CodeTokenExpired, CodeToI18nKey[CodeTokenExpired])
	ErrInvalidToken            = NewAppErrorWithI18n(CodeInvalidToken, CodeToI18nKey[CodeInvalidToken])
	ErrAccountDisabled         = NewAppErrorWithI18n(CodeAccountDisabled, CodeToI18nKey[CodeAccountDisabled])
	ErrAccountLocked           = NewAppErrorWithI18n(CodeAccountLocked, CodeToI18nKey[CodeAccountLocked])
	ErrPasswordTooWeak         = NewAppErrorWithI18n(CodePasswordTooWeak, CodeToI18nKey[CodePasswordTooWeak])
	ErrInvalidVerificationCode = NewAppErrorWithI18n(CodeInvalidVerificationCode, CodeToI18nKey[CodeInvalidVerificationCode])
	ErrLoginAttemptsExceeded   = NewAppErrorWithI18n(CodeLoginAttemptsExceeded, CodeToI18nKey[CodeLoginAttemptsExceeded])
	ErrEmailAlreadyExists      = NewAppErrorWithI18n(CodeEmailAlreadyExists, CodeToI18nKey[CodeEmailAlreadyExists])
	ErrTokenAlreadyUsed        = NewAppErrorWithI18n(CodeTokenAlreadyUsed, CodeToI18nKey[CodeTokenAlreadyUsed])
	ErrUsernameSuported        = NewAppErrorWithI18n(CodeUsernameSuported, CodeToI18nKey[CodeUsernameSuported])
	ErrEmailSuported           = NewAppErrorWithI18n(CodeEmailSuported, CodeToI18nKey[CodeEmailSuported])

	// Permission related errors (2000-2999)
	ErrInsufficientPermissions       = NewAppErrorWithI18n(CodeInsufficientPermissions, CodeToI18nKey[CodeInsufficientPermissions])
	ErrResourceAccessDenied          = NewAppErrorWithI18n(CodeResourceAccessDenied, CodeToI18nKey[CodeResourceAccessDenied])
	ErrOperationPermissionDenied     = NewAppErrorWithI18n(CodeOperationPermissionDenied, CodeToI18nKey[CodeOperationPermissionDenied])
	ErrRolePermissionDenied          = NewAppErrorWithI18n(CodeRolePermissionDenied, CodeToI18nKey[CodeRolePermissionDenied])
	ErrResourceGroupPermissionDenied = NewAppErrorWithI18n(CodeResourceGroupPermissionDenied, CodeToI18nKey[CodeResourceGroupPermissionDenied])

	// Parameter validation errors (3000-3999)
	ErrValidationFailed         = NewAppErrorWithI18n(CodeValidationFailed, CodeToI18nKey[CodeValidationFailed])
	ErrRequiredParameterMissing = NewAppErrorWithI18n(CodeRequiredParameterMissing, CodeToI18nKey[CodeRequiredParameterMissing])
	ErrInvalidParameterFormat   = NewAppErrorWithI18n(CodeInvalidParameterFormat, CodeToI18nKey[CodeInvalidParameterFormat])
	ErrParameterOutOfRange      = NewAppErrorWithI18n(CodeParameterOutOfRange, CodeToI18nKey[CodeParameterOutOfRange])
	ErrInvalidParameterLength   = NewAppErrorWithI18n(CodeInvalidParameterLength, CodeToI18nKey[CodeInvalidParameterLength])
	ErrInvalidEmailFormat       = NewAppErrorWithI18n(CodeInvalidEmailFormat, CodeToI18nKey[CodeInvalidEmailFormat])
	ErrInvalidPhoneFormat       = NewAppErrorWithI18n(CodeInvalidPhoneFormat, CodeToI18nKey[CodeInvalidPhoneFormat])
	ErrInvalidURLFormat         = NewAppErrorWithI18n(CodeInvalidURLFormat, CodeToI18nKey[CodeInvalidURLFormat])
	ErrInvalidDateFormat        = NewAppErrorWithI18n(CodeInvalidDateFormat, CodeToI18nKey[CodeInvalidDateFormat])
	ErrEmailNotVerified         = NewAppErrorWithI18n(CodeEmailNotVerified, CodeToI18nKey[CodeEmailNotVerified])

	// Resource related errors (4000-4999)
	ErrRecordNotFound             = NewAppErrorWithI18n(CodeRecordNotFound, CodeToI18nKey[CodeRecordNotFound])
	ErrResourceNotFound           = NewAppErrorWithI18n(CodeResourceNotFound, CodeToI18nKey[CodeResourceNotFound])
	ErrResourceAlreadyExists      = NewAppErrorWithI18n(CodeResourceAlreadyExists, CodeToI18nKey[CodeResourceAlreadyExists])
	ErrResourceStateNotAllowed    = NewAppErrorWithI18n(CodeResourceStateNotAllowed, CodeToI18nKey[CodeResourceStateNotAllowed])
	ErrResourceDependencyConflict = NewAppErrorWithI18n(CodeResourceDependencyConflict, CodeToI18nKey[CodeResourceDependencyConflict])
	ErrResourceQuotaInsufficient  = NewAppErrorWithI18n(CodeResourceQuotaInsufficient, CodeToI18nKey[CodeResourceQuotaInsufficient])
	ErrResourceInUse              = NewAppErrorWithI18n(CodeResourceInUse, CodeToI18nKey[CodeResourceInUse])
	ErrRecordQueryFailed          = NewAppErrorWithI18n(CodeRecordQueryFailed, CodeToI18nKey[CodeRecordQueryFailed])
	ErrRecordCreateFailed         = NewAppErrorWithI18n(CodeRecordCreateFailed, CodeToI18nKey[CodeRecordCreateFailed])
	ErrRecordUpdateFailed         = NewAppErrorWithI18n(CodeRecordUpdateFailed, CodeToI18nKey[CodeRecordUpdateFailed])
	ErrRecordDeleteFailed         = NewAppErrorWithI18n(CodeRecordDeleteFailed, CodeToI18nKey[CodeRecordDeleteFailed])
	ErrRecordIsDisabled           = NewAppErrorWithI18n(CodeRecordIsDisabled, CodeToI18nKey[CodeRecordIsDisabled])
	ErrRecordNoAffected           = NewAppErrorWithI18n(CodeRecordNoAffected, CodeToI18nKey[CodeRecordNoAffected])
	ErrRecordDeleteDenied         = NewAppErrorWithI18n(CodeRecordDeleteDenied, CodeToI18nKey[CodeRecordDeleteDenied])

	// Business logic errors (5000-5999)
	ErrServerOffline                  = NewAppErrorWithI18n(CodeServerOffline, CodeToI18nKey[CodeServerOffline])
	ErrAppDeploymentFailed            = NewAppErrorWithI18n(CodeAppDeploymentFailed, CodeToI18nKey[CodeAppDeploymentFailed])
	ErrWorkflowExecutionFailed        = NewAppErrorWithI18n(CodeWorkflowExecutionFailed, CodeToI18nKey[CodeWorkflowExecutionFailed])
	ErrCertificateRequestFailed       = NewAppErrorWithI18n(CodeCertificateRequestFailed, CodeToI18nKey[CodeCertificateRequestFailed])
	ErrBackupOperationFailed          = NewAppErrorWithI18n(CodeBackupOperationFailed, CodeToI18nKey[CodeBackupOperationFailed])
	ErrMonitoringDataCollectionFailed = NewAppErrorWithI18n(CodeMonitoringDataCollectionFailed, CodeToI18nKey[CodeMonitoringDataCollectionFailed])
	ErrAppPublishFailed               = NewAppErrorWithI18n(CodeAppPublishFailed, CodeToI18nKey[CodeAppPublishFailed])
	ErrAppOfflineFailed               = NewAppErrorWithI18n(CodeAppOfflineFailed, CodeToI18nKey[CodeAppOfflineFailed])
	ErrHealthCheckFailed              = NewAppErrorWithI18n(CodeHealthCheckFailed, CodeToI18nKey[CodeHealthCheckFailed])
	ErrGatewayConfigUpdateFailed      = NewAppErrorWithI18n(CodeGatewayConfigUpdateFailed, CodeToI18nKey[CodeGatewayConfigUpdateFailed])
	ErrUserAlreadyExists              = NewAppErrorWithI18n(CodeUserAlreadyExists, CodeToI18nKey[CodeUserAlreadyExists])
	ErrPermissionInvalid              = NewAppErrorWithI18n(CodePermissionInvalid, CodeToI18nKey[CodePermissionInvalid])

	// Server management specific errors
	ErrServerNotFound      = NewAppErrorWithI18n(CodeResourceNotFound, CodeToI18nKey[CodeResourceNotFound])
	ErrServerNameExists    = NewAppErrorWithI18n(CodeResourceAlreadyExists, CodeToI18nKey[CodeResourceAlreadyExists])
	ErrServerAgentNotFound = NewAppErrorWithI18n(CodeResourceNotFound, CodeToI18nKey[CodeResourceNotFound])

	// System related errors (6000-6999)
	ErrInternalError                = NewAppErrorWithI18n(CodeInternalError, CodeToI18nKey[CodeInternalError])
	ErrCacheServiceUnavailable      = NewAppErrorWithI18n(CodeCacheServiceUnavailable, CodeToI18nKey[CodeCacheServiceUnavailable])
	ErrFilesystemError              = NewAppErrorWithI18n(CodeFilesystemError, CodeToI18nKey[CodeFilesystemError])
	ErrNetworkTimeout               = NewAppErrorWithI18n(CodeNetworkTimeout, CodeToI18nKey[CodeNetworkTimeout])
	ErrThirdPartyServiceUnavailable = NewAppErrorWithI18n(CodeThirdPartyServiceUnavailable, CodeToI18nKey[CodeThirdPartyServiceUnavailable])
	ErrSystemMaintenance            = NewAppErrorWithI18n(CodeSystemMaintenance, CodeToI18nKey[CodeSystemMaintenance])
	ErrDatabaseConnectionFailed     = NewAppErrorWithI18n(CodeDatabaseConnectionFailed, CodeToI18nKey[CodeDatabaseConnectionFailed])
)
