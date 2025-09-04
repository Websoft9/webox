package validator

import (
	"api-service/pkg/errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// New creates a new validator instance
func New() *validator.Validate {
	return validator.New()
}

// Password validation related constants
const (
	MinPasswordRequirements = 2   // Minimum number of password requirements to meet
	EmailPartsCount         = 2   // Number of parts in email address after @ split
	MaxApplications         = 10  // Maximum number of applications
	MaxWorkflows            = 5   // Maximum number of workflows
	MaxFiles                = 100 // Maximum number of files
)

// System reserved usernames
var reservedUsernames = map[string]bool{
	"admin":     true,
	"root":      true,
	"system":    true,
	"websoft9":  true,
	"api":       true,
	"www":       true,
	"ftp":       true,
	"mail":      true,
	"test":      true,
	"guest":     true,
	"anonymous": true,
}

// Allowed email domain whitelist (if empty, all domains are allowed)
var allowedEmailDomains = []string{
	// "company.com",
	// "websoft9.com",
}

// ValidateUsername validate username
func ValidateUsername(username string) error {
	// Check if empty
	if username == "" {
		return errors.NewAppError(errors.CodeValidationError, "Username cannot be empty")
	}

	// Check length
	if len(username) < 3 || len(username) > 20 {
		return errors.NewAppError(errors.CodeValidationError, "Username length must be between 3-20 characters")
	}

	// Check character rules: only allow letters, numbers, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.NewAppError(errors.CodeValidationError, "Username can only contain letters, numbers, and underscores")
	}

	// Check if starts with a letter
	if !unicode.IsLetter(rune(username[0])) {
		return errors.NewAppError(errors.CodeValidationError, "Username must start with a letter")
	}

	// Check if it's a reserved username
	if reservedUsernames[strings.ToLower(username)] {
		return errors.ErrUsernameReserved
	}

	return nil
}

// ValidatePassword validate password strength
func ValidatePassword(password string) error {
	// Check if empty
	if password == "" {
		return errors.NewAppError(errors.CodeValidationError, "Password cannot be empty")
	}

	// Check length
	if len(password) < 6 || len(password) > 50 {
		return errors.NewAppError(errors.CodeValidationError, "Password length must be between 6-50 characters")
	}

	// Check password complexity
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// At least contains two of: uppercase letters, lowercase letters, numbers
	requirements := 0
	if hasUpper {
		requirements++
	}
	if hasLower {
		requirements++
	}
	if hasNumber {
		requirements++
	}
	if hasSpecial {
		requirements++
	}

	if requirements < MinPasswordRequirements {
		return errors.ErrPasswordTooWeak
	}

	return nil
}

// ValidateEmail validate email format and domain
func ValidateEmail(email string) error {
	// Check if empty
	if email == "" {
		return errors.NewAppError(errors.CodeValidationError, "Email cannot be empty")
	}

	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmail
	}

	// Domain whitelist validation (if whitelist is configured)
	if len(allowedEmailDomains) > 0 {
		parts := strings.Split(email, "@")
		if len(parts) != EmailPartsCount {
			return errors.ErrInvalidEmail
		}

		domain := strings.ToLower(parts[1])
		allowed := false
		for _, allowedDomain := range allowedEmailDomains {
			if strings.EqualFold(domain, allowedDomain) {
				allowed = true
				break
			}
		}

		if !allowed {
			return errors.NewAppError(errors.CodeValidationError, "This email domain is not allowed")
		}
	}

	return nil
}

// ValidateUserStatus validate user status transition
func ValidateUserStatus(currentStatus, newStatus string) error {
	// Define allowed status transitions
	allowedTransitions := map[string][]string{
		"inactive": {"active", "banned"},
		"active":   {"inactive", "banned"},
		"banned":   {"inactive"},
	}

	validStatuses := []string{"active", "inactive", "banned"}

	// Check if the new status is valid
	isValidStatus := false
	for _, status := range validStatuses {
		if newStatus == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return errors.NewAppError(errors.CodeValidationError, "Invalid user status")
	}

	// Check if the status transition is allowed
	if allowedNextStatuses, exists := allowedTransitions[currentStatus]; exists {
		for _, allowedStatus := range allowedNextStatuses {
			if newStatus == allowedStatus {
				return nil
			}
		}
		return errors.NewAppError(errors.CodeValidationError, "Status transition not allowed")
	}

	return errors.NewAppError(errors.CodeValidationError, "Current status does not support transition")
}

// ValidateUserPermission validate user permission
func ValidateUserPermission(userRole, requiredPermission string) error {
	// Define role permission mapping
	rolePermissions := map[string][]string{
		"admin": {
			"user:create", "user:read", "user:update", "user:delete",
			"app:create", "app:read", "app:update", "app:delete",
		},
		"user":  {"user:read", "app:create", "app:read", "app:update"},
		"guest": {"user:read", "app:read"},
	}

	permissions, exists := rolePermissions[userRole]
	if !exists {
		return errors.NewAppError(errors.CodeValidationError, "Invalid user role")
	}

	for _, permission := range permissions {
		if permission == requiredPermission {
			return nil
		}
	}

	return errors.ErrForbidden
}

// ValidateResourceQuota validate resource quota
func ValidateResourceQuota(userID uint, resourceType string, currentCount int) error {
	// Define resource quota limits
	quotaLimits := map[string]int{
		"applications": MaxApplications,
		"workflows":    MaxWorkflows,
		"files":        MaxFiles,
	}

	limit, exists := quotaLimits[resourceType]
	if !exists {
		return errors.NewAppError(errors.CodeValidationError, "Unknown resource type")
	}

	if currentCount >= limit {
		return errors.ErrUserQuotaExceeded
	}

	return nil
}
