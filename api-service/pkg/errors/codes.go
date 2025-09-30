package errors

import (
	"net/http"
)

// AppError represents a custom application error type
type ErrorCode int

// HTTPCode represents the HTTP status code for an error
type HTTPCode int

// Error code constants definition based on API documentation
// Error codes are organized by category with specific ranges for easy identification
const (
	// Success codes (0)
	CodeSuccess ErrorCode = 200 // Operation completed successfully

	// Authentication related error codes (1000-1999)
	CodeInvalidCredentials      ErrorCode = 1001 // Username or password incorrect
	CodeTokenExpired            ErrorCode = 1002 // Token has expired
	CodeInvalidToken            ErrorCode = 1003 // Invalid token
	CodeAccountDisabled         ErrorCode = 1004 // Account has been disabled
	CodeAccountLocked           ErrorCode = 1005 // Account has been locked
	CodePasswordTooWeak         ErrorCode = 1006 // Password strength insufficient
	CodeInvalidVerificationCode ErrorCode = 1007 // Verification code incorrect
	CodeLoginAttemptsExceeded   ErrorCode = 1008 // Too many login failures
	CodeEmailAlreadyExists      ErrorCode = 1009 // Email already exists
	CodeTokenAlreadyUsed        ErrorCode = 1010 // Token already used
	CodeUsernameSuported        ErrorCode = 1011 // Username supported
	CodeEmailSuported           ErrorCode = 1012 // Email supported

	// Permission related error codes (2000-2999)
	CodeInsufficientPermissions       ErrorCode = 2001 // Insufficient permissions
	CodeResourceAccessDenied          ErrorCode = 2002 // Resource access denied
	CodeOperationPermissionDenied     ErrorCode = 2003 // Operation permission insufficient
	CodeRolePermissionDenied          ErrorCode = 2004 // Role permission insufficient
	CodeResourceGroupPermissionDenied ErrorCode = 2005 // Resource group permission insufficient
	CodeAccessDenied                  ErrorCode = 2006 // Access denied

	// Parameter validation error codes (3000-3999)
	CodeValidationFailed         ErrorCode = 3000 // General validation error
	CodeRequiredParameterMissing ErrorCode = 3001 // Required parameter missing
	CodeInvalidParameterFormat   ErrorCode = 3002 // Parameter format error
	CodeParameterOutOfRange      ErrorCode = 3003 // Parameter value out of range
	CodeInvalidParameterLength   ErrorCode = 3004 // Parameter length does not meet requirements
	CodeInvalidEmailFormat       ErrorCode = 3005 // Email format error
	CodeInvalidPhoneFormat       ErrorCode = 3006 // Phone number format error
	CodeInvalidURLFormat         ErrorCode = 3007 // URL format error
	CodeInvalidDateFormat        ErrorCode = 3008 // Date format error
	CodeEmailNotVerified         ErrorCode = 3009 // Date format error

	// Resource related error codes (4000-4999)
	CodeRecordNotFound             ErrorCode = 4000 // Record does not exist
	CodeResourceNotFound           ErrorCode = 4001 // Resource does not exist
	CodeResourceAlreadyExists      ErrorCode = 4002 // Resource already exists
	CodeResourceStateNotAllowed    ErrorCode = 4003 // Resource state does not allow operation
	CodeResourceDependencyConflict ErrorCode = 4004 // Resource dependency conflict
	CodeResourceQuotaInsufficient  ErrorCode = 4005 // Resource quota insufficient
	CodeResourceInUse              ErrorCode = 4006 // Resource is in use
	CodeRecordQueryFailed          ErrorCode = 4007 // Record query failed
	CodeRecordCreateFailed         ErrorCode = 4008 // Record creation failed
	CodeRecordUpdateFailed         ErrorCode = 4009 // Record update failed
	CodeRecordDeleteFailed         ErrorCode = 4010 // Record deletion failed
	CodeRecordIsDisabled           ErrorCode = 4011 // Record is disabled
	CodeRecordNoAffected           ErrorCode = 4012 // Record is not affected
	CodeRecordDeleteDenied         ErrorCode = 4013 // Record delete denied

	// Business logic error codes (5000-5999)
	CodeServerOffline                  = 5001 // Server offline, cannot operate
	CodeAppDeploymentFailed            = 5002 // Application deployment failed
	CodeWorkflowExecutionFailed        = 5003 // Workflow execution failed
	CodeCertificateRequestFailed       = 5004 // Certificate request failed
	CodeBackupOperationFailed          = 5005 // Backup operation failed
	CodeMonitoringDataCollectionFailed = 5006 // Monitoring data collection failed
	CodeAppPublishFailed               = 5007 // Application publish failed
	CodeAppOfflineFailed               = 5008 // Application offline failed
	CodeHealthCheckFailed              = 5009 // Health check failed
	CodeGatewayConfigUpdateFailed      = 5010 // Gateway configuration update failed
	CodeUserAlreadyExists              = 5011 // Username already exists
	CodePermissionInvalid              = 5012 // Permission invalid
	CodeEncryptFailed                  = 5013 // Data encryption failed
	CodeDecryptFailed                  = 5014 // Data decryption failed

	// System related error codes (6000-6999)
	CodeInternalError                ErrorCode = 6001 // System Internal Error
	CodeCacheServiceUnavailable      ErrorCode = 6002 // Cache service unavailable
	CodeFilesystemError              ErrorCode = 6003 // Filesystem error
	CodeNetworkTimeout               ErrorCode = 6004 // Network connection timeout
	CodeThirdPartyServiceUnavailable ErrorCode = 6005 // Third party service unavailable
	CodeSystemMaintenance            ErrorCode = 6006 // System under maintenance
	CodeDatabaseConnectionFailed     ErrorCode = 6007 // Database connection failed
)

// CodeToI18nKey maps error codes to their i18n message keys
// These keys should correspond to entries in the i18n locale files
var CodeToI18nKey = map[ErrorCode]string{
	// Success codes
	CodeSuccess: "common.success",

	// Authentication related errors (1000-1999)
	CodeInvalidCredentials:      "auth.invalid_credentials",
	CodeTokenExpired:            "auth.token_expired",
	CodeInvalidToken:            "auth.token_invalid",
	CodeAccountDisabled:         "auth.account_disabled",
	CodeAccountLocked:           "auth.account_locked",
	CodePasswordTooWeak:         "auth.password_too_weak",
	CodeInvalidVerificationCode: "auth.verification_code_invalid",
	CodeLoginAttemptsExceeded:   "auth.login_attempts_exceeded",
	CodeEmailAlreadyExists:      "auth.email_already_exists",
	CodeTokenAlreadyUsed:        "auth.token_already_used",
	CodeUsernameSuported:        "auth.username_supported",
	CodeEmailSuported:           "auth.email_supported",

	// Permission related errors (2000-2999)
	CodeInsufficientPermissions:       "auth.permission_denied",
	CodeResourceAccessDenied:          "auth.resource_access_denied",
	CodeOperationPermissionDenied:     "auth.operation_permission_denied",
	CodeRolePermissionDenied:          "auth.role_permission_denied",
	CodeResourceGroupPermissionDenied: "auth.resource_group_permission_denied",
	CodeAccessDenied:                  "auth.access_denied",

	// Parameter validation errors (3000-3999)
	CodeValidationFailed:         "validation.validation_failed",
	CodeRequiredParameterMissing: "validation.required_parameter_missing",
	CodeInvalidParameterFormat:   "validation.invalid_parameter_format",
	CodeParameterOutOfRange:      "validation.parameter_out_of_range",
	CodeInvalidParameterLength:   "validation.invalid_parameter_length",
	CodeInvalidEmailFormat:       "validation.invalid_email_format",
	CodeInvalidPhoneFormat:       "validation.invalid_phone_format",
	CodeInvalidURLFormat:         "validation.invalid_url_format",
	CodeInvalidDateFormat:        "validation.invalid_date_format",
	CodeEmailNotVerified:         "validation.email_not_verified",

	// Resource related errors (4000-4999)
	CodeRecordNotFound:             "resource.record_not_found",
	CodeResourceNotFound:           "resource.not_found",
	CodeResourceAlreadyExists:      "resource.already_exists",
	CodeResourceStateNotAllowed:    "resource.state_not_allowed",
	CodeResourceDependencyConflict: "resource.dependency_conflict",
	CodeResourceQuotaInsufficient:  "resource.quota_insufficient",
	CodeResourceInUse:              "resource.in_use",
	CodeRecordQueryFailed:          "resource.record_query_failed",
	CodeRecordCreateFailed:         "resource.record_create_failed",
	CodeRecordUpdateFailed:         "resource.record_update_failed",
	CodeRecordDeleteFailed:         "resource.record_delete_failed",
	CodeRecordIsDisabled:           "resource.record_is_disabled",
	CodeRecordNoAffected:           "resource.record_no_affected",
	CodeRecordDeleteDenied:         "resource.record_delete_denied",

	// Business logic errors (5000-5999)
	CodeServerOffline:                  "business.server_offline",
	CodeAppDeploymentFailed:            "business.app_deployment_failed",
	CodeWorkflowExecutionFailed:        "business.workflow_execution_failed",
	CodeCertificateRequestFailed:       "business.certificate_request_failed",
	CodeBackupOperationFailed:          "business.backup_operation_failed",
	CodeMonitoringDataCollectionFailed: "business.monitoring_data_collection_failed",
	CodeAppPublishFailed:               "business.app_publish_failed",
	CodeAppOfflineFailed:               "business.app_offline_failed",
	CodeHealthCheckFailed:              "business.health_check_failed",
	CodeGatewayConfigUpdateFailed:      "business.gateway_config_update_failed",
	CodeUserAlreadyExists:              "business.user_already_exists",
	CodePermissionInvalid:              "business.permission_invalid",
	CodeEncryptFailed:                  "business.encrypt_failed",
	CodeDecryptFailed:                  "business.decrypt_failed",

	// System related errors (6000-6999)
	CodeInternalError:                "system.internal_error",
	CodeCacheServiceUnavailable:      "system.cache_service_unavailable",
	CodeFilesystemError:              "system.filesystem_error",
	CodeNetworkTimeout:               "system.network_timeout",
	CodeThirdPartyServiceUnavailable: "system.third_party_service_unavailable",
	CodeSystemMaintenance:            "system.system_maintenance",
	CodeDatabaseConnectionFailed:     "system.database_connection_failed",
}

// codeToHTTPStatus maps business error codes to HTTP status codes
// This mapping ensures consistent HTTP responses for different error types
var CodeToHTTPStatus = map[ErrorCode]HTTPCode{
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
	CodeUsernameSuported:        http.StatusBadRequest,
	CodeEmailSuported:           http.StatusBadRequest,

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
