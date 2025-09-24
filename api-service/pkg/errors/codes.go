package errors

// Error code constants definition based on API documentation
// Error codes are organized by category with specific ranges for easy identification
const (
	// Success codes (0)
	CodeSuccess = 0 // Operation completed successfully

	// Authentication related error codes (1000-1999)
	CodeInvalidCredentials      = 1001 // Username or password incorrect
	CodeTokenExpired            = 1002 // Token has expired
	CodeInvalidToken            = 1003 // Invalid token
	CodeAccountDisabled         = 1004 // Account has been disabled
	CodeAccountLocked           = 1005 // Account has been locked
	CodePasswordTooWeak         = 1006 // Password strength insufficient
	CodeInvalidVerificationCode = 1007 // Verification code incorrect
	CodeLoginAttemptsExceeded   = 1008 // Too many login failures
	CodeEmailAlreadyExists      = 1009 // Email already exists
	CodeTokenAlreadyUsed        = 1010 // Token already used
	CodeUsernameSuported        = 1011 // Username supported
	CodeEmailSuported           = 1012 // Email supported

	// Permission related error codes (2000-2999)
	CodeInsufficientPermissions       = 2001 // Insufficient permissions
	CodeResourceAccessDenied          = 2002 // Resource access denied
	CodeOperationPermissionDenied     = 2003 // Operation permission insufficient
	CodeRolePermissionDenied          = 2004 // Role permission insufficient
	CodeResourceGroupPermissionDenied = 2005 // Resource group permission insufficient
	CodeAccessDenied                  = 2006 // Access denied

	// Parameter validation error codes (3000-3999)
	CodeValidationFailed         = 3000 // General validation error
	CodeRequiredParameterMissing = 3001 // Required parameter missing
	CodeInvalidParameterFormat   = 3002 // Parameter format error
	CodeParameterOutOfRange      = 3003 // Parameter value out of range
	CodeInvalidParameterLength   = 3004 // Parameter length does not meet requirements
	CodeInvalidEmailFormat       = 3005 // Email format error
	CodeInvalidPhoneFormat       = 3006 // Phone number format error
	CodeInvalidURLFormat         = 3007 // URL format error
	CodeInvalidDateFormat        = 3008 // Date format error
	CodeEmailNotVerified         = 3009 // Date format error

	// Resource related error codes (4000-4999)
	CodeRecordNotFound             = 4000 // Record does not exist
	CodeResourceNotFound           = 4001 // Resource does not exist
	CodeResourceAlreadyExists      = 4002 // Resource already exists
	CodeResourceStateNotAllowed    = 4003 // Resource state does not allow operation
	CodeResourceDependencyConflict = 4004 // Resource dependency conflict
	CodeResourceQuotaInsufficient  = 4005 // Resource quota insufficient
	CodeResourceInUse              = 4006 // Resource is in use
	CodeRecordQueryFailed          = 4007 // Record query failed
	CodeRecordCreateFailed         = 4008 // Record creation failed
	CodeRecordUpdateFailed         = 4009 // Record update failed
	CodeRecordDeleteFailed         = 4010 // Record deletion failed
	CodeRecordIsDisabled           = 4011 // Record is disabled
	CodeRecordNoAffected           = 4012 // Record is not affected
	CodeRecordDeleteDenied         = 4013 // Record delete denied

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
	CodeInternalError                = 6001 // System Internal Error
	CodeCacheServiceUnavailable      = 6002 // Cache service unavailable
	CodeFilesystemError              = 6003 // Filesystem error
	CodeNetworkTimeout               = 6004 // Network connection timeout
	CodeThirdPartyServiceUnavailable = 6005 // Third party service unavailable
	CodeSystemMaintenance            = 6006 // System under maintenance
	CodeDatabaseConnectionFailed     = 6007 // Database connection failed
)

// CodeToI18nKey maps error codes to their i18n message keys
// These keys should correspond to entries in the i18n locale files
var CodeToI18nKey = map[int]string{
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

// CodeMessages maps error codes to their default English messages
// These messages serve as fallbacks when internationalization is not available
var CodeMessages = map[int]string{
	// Success codes
	CodeSuccess: "Success",

	// Authentication related errors (1000-1999)
	CodeInvalidCredentials:      "Username or password incorrect",
	CodeTokenExpired:            "Token has expired",
	CodeInvalidToken:            "Invalid token",
	CodeAccountDisabled:         "Account has been disabled",
	CodeAccountLocked:           "Account has been locked",
	CodePasswordTooWeak:         "Password strength insufficient",
	CodeInvalidVerificationCode: "Verification code incorrect",
	CodeLoginAttemptsExceeded:   "Too many login failures",
	CodeEmailAlreadyExists:      "Email already exists",
	CodeTokenAlreadyUsed:        "Token has already been used",
	CodeUsernameSuported:        "Only username login is supported",
	CodeEmailSuported:           "Only email login is supported",

	// Permission related errors (2000-2999)
	CodeInsufficientPermissions:       "Insufficient permissions",
	CodeResourceAccessDenied:          "Resource access denied",
	CodeOperationPermissionDenied:     "Operation permission insufficient",
	CodeRolePermissionDenied:          "Role permission insufficient",
	CodeResourceGroupPermissionDenied: "Resource group permission insufficient",
	CodeAccessDenied:                  "Access denied",

	// Parameter validation errors (3000-3999)
	CodeValidationFailed:         "Validation failed",
	CodeRequiredParameterMissing: "Required parameter missing",
	CodeInvalidParameterFormat:   "Parameter format error",
	CodeParameterOutOfRange:      "Parameter value out of range",
	CodeInvalidParameterLength:   "Parameter length does not meet requirements",
	CodeInvalidEmailFormat:       "Email format error",
	CodeInvalidPhoneFormat:       "Phone number format error",
	CodeInvalidURLFormat:         "URL format error",
	CodeInvalidDateFormat:        "Date format error",
	CodeEmailNotVerified:         "Email not verified",

	// Resource related errors (4000-4999)
	CodeRecordNotFound:             "Record does not exist",
	CodeResourceNotFound:           "Resource does not exist",
	CodeResourceAlreadyExists:      "Resource already exists",
	CodeResourceStateNotAllowed:    "Resource state does not allow operation",
	CodeResourceDependencyConflict: "Resource dependency conflict",
	CodeResourceQuotaInsufficient:  "Resource quota insufficient",
	CodeResourceInUse:              "Resource is in use",
	CodeRecordQueryFailed:          "Record query failed",
	CodeRecordCreateFailed:         "Record create failed",
	CodeRecordUpdateFailed:         "Record update failed",
	CodeRecordDeleteFailed:         "Record delete failed",
	CodeRecordIsDisabled:           "Record is disabled",
	CodeRecordNoAffected:           "No rows affected",
	CodeRecordDeleteDenied:         "Record delete denied",

	// Business logic errors (5000-5999)
	CodeServerOffline:                  "Server offline, cannot operate",
	CodeAppDeploymentFailed:            "Application deployment failed",
	CodeWorkflowExecutionFailed:        "Workflow execution failed",
	CodeCertificateRequestFailed:       "Certificate request failed",
	CodeBackupOperationFailed:          "Backup operation failed",
	CodeMonitoringDataCollectionFailed: "Monitoring data collection failed",
	CodeAppPublishFailed:               "Application publish failed",
	CodeAppOfflineFailed:               "Application offline failed",
	CodeHealthCheckFailed:              "Health check failed",
	CodeGatewayConfigUpdateFailed:      "Gateway configuration update failed",
	CodeUserAlreadyExists:              "Username already exists",
	CodePermissionInvalid:              "Permission invalid",

	// System related errors (6000-6999)
	CodeInternalError:                "System Internal Error",
	CodeCacheServiceUnavailable:      "Cache service unavailable",
	CodeFilesystemError:              "Filesystem error",
	CodeNetworkTimeout:               "Network connection timeout",
	CodeThirdPartyServiceUnavailable: "Third party service unavailable",
	CodeSystemMaintenance:            "System under maintenance",
	CodeDatabaseConnectionFailed:     "Database connection failed",
}
