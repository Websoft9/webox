package errors

// Error code constants definition
// Error codes are organized by category with specific ranges for easy identification
const (
	// General error codes (10000-19999)
	// These codes represent system-level and common application errors
	CodeSuccess         = 0     // Operation completed successfully
	CodeInternalError   = 10001 // Internal server error
	CodeInvalidRequest  = 10002 // Invalid request parameters
	CodeUnauthorized    = 10003 // Authentication required
	CodeForbidden       = 10004 // Access denied
	CodeNotFound        = 10005 // Resource not found
	CodeValidationError = 10006 // Data validation failed

	// User-related error codes (20000-29999)
	// These codes represent user management and authentication errors
	CodeUserNotFound       = 20001 // User does not exist
	CodeUserAlreadyExists  = 20002 // Username already taken
	CodeInvalidCredentials = 20003 // Invalid username or password
	CodeUserInactive       = 20004 // User account is inactive
	CodeInvalidPassword    = 20005 // Password is incorrect
	CodePasswordTooWeak    = 20006 // Password does not meet requirements
	CodeEmailAlreadyExists = 20007 // Email address already registered
	CodeInvalidEmail       = 20008 // Email format is invalid
	CodeUsernameReserved   = 20009 // Username is reserved by system
	CodeUserQuotaExceeded  = 20010 // User quota limit exceeded
	CodeEmailNotVerified   = 20011 // Email address not verified
	CodeInvalidToken       = 20012 // Invalid verification token
	CodeTokenExpired       = 20013 // Verification token expired
	CodeTokenAlreadyUsed   = 20014 // Verification token already used

	// Application-related error codes (30000-39999)
	// These codes represent application deployment and management errors
	CodeAppNotFound              = 30001 // Application not found
	CodeAppAlreadyExists         = 30002 // Application already exists
	CodeAppDeployFailed          = 30003 // Application deployment failed
	CodeAppPortConflict          = 30004 // Port conflict during deployment
	CodeAppResourcesInsufficient = 30005 // Insufficient resources for deployment
)

// CodeMessages maps error codes to their default English messages
// These messages serve as fallbacks when internationalization is not available
var CodeMessages = map[int]string{
	// General error messages
	CodeSuccess:         "Success",
	CodeInternalError:   "Internal server error",
	CodeInvalidRequest:  "Invalid request parameters",
	CodeUnauthorized:    "Unauthorized access",
	CodeForbidden:       "Access forbidden",
	CodeNotFound:        "Resource not found",
	CodeValidationError: "Data validation failed",

	// User-related error messages
	CodeUserNotFound:       "User not found",
	CodeUserAlreadyExists:  "Username already exists",
	CodeInvalidCredentials: "Invalid username or password",
	CodeUserInactive:       "User account is inactive",
	CodeInvalidPassword:    "Invalid password",
	CodePasswordTooWeak:    "Password is too weak",
	CodeEmailAlreadyExists: "Email address already exists",
	CodeInvalidEmail:       "Invalid email format",
	CodeUsernameReserved:   "Username is reserved by system",
	CodeUserQuotaExceeded:  "User quota exceeded",
	CodeEmailNotVerified:   "Email address not verified",
	CodeInvalidToken:       "Invalid verification token",
	CodeTokenExpired:       "Verification token expired",
	CodeTokenAlreadyUsed:   "Verification token already used",

	// Application-related error messages
	CodeAppNotFound:              "Application not found",
	CodeAppAlreadyExists:         "Application already exists",
	CodeAppDeployFailed:          "Application deployment failed",
	CodeAppPortConflict:          "Port conflict",
	CodeAppResourcesInsufficient: "Insufficient resources",
}
