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
	MinPasswordRequirements = 2 // Minimum number of password requirements to meet
	EmailPartsCount         = 2 // Number of parts in email address after @ split
)

// Allowed email domain whitelist (if empty, all domains are allowed)
var allowedEmailDomains = []string{
	// "company.com",
	// "websoft9.com",
}

// ValidateUsername validate username
func ValidateUsername(username string) error {
	// Check if empty
	if username == "" {
		return errors.NewAppError(errors.CodeInvalidParameterFormat, "Username cannot be empty")
	}

	// Check length
	if len(username) < 3 || len(username) > 20 {
		return errors.NewAppError(errors.CodeInvalidParameterFormat, "Username length must be between 3-20 characters")
	}

	// Check character rules: only allow letters, numbers, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.NewAppError(errors.CodeInvalidParameterFormat, "Username can only contain letters, numbers, and underscores")
	}

	// Check if starts with a letter
	if !unicode.IsLetter(rune(username[0])) {
		return errors.NewAppError(errors.CodeInvalidParameterFormat, "Username must start with a letter")
	}

	return nil
}

// ValidatePassword validate password strength
func ValidatePassword(password string) error {
	// Check if empty
	if password == "" {
		return errors.NewAppError(errors.CodeRequiredParameterMissing, "Password cannot be empty")
	}

	// Check length
	if len(password) < 6 || len(password) > 50 {
		return errors.NewAppError(errors.CodeRequiredParameterMissing, "Password length must be between 6-50 characters")
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
		return errors.ErrInvalidEmailFormat
	}

	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmailFormat
	}

	// Domain whitelist validation (if whitelist is configured)
	if len(allowedEmailDomains) > 0 {
		parts := strings.Split(email, "@")
		if len(parts) != EmailPartsCount {
			return errors.ErrInvalidEmailFormat
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
			return errors.NewAppError(errors.CodeValidationFailed, "This email domain is not allowed")
		}
	}

	return nil
}

// IsEmail checks if the given string is a valid email format
func IsEmail(input string) bool {
	if input == "" {
		return false
	}

	// Basic email format check
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(input)
}

// ValidateUsernameOrEmail validates that input is either a valid username or email
func ValidateUsernameOrEmail(input string) error {
	if input == "" {
		return errors.NewAppError(errors.CodeRequiredParameterMissing, "Username or email cannot be empty")
	}

	// Check if it's an email
	if IsEmail(input) {
		return ValidateEmail(input)
	}

	// If not an email, validate as username
	return ValidateUsername(input)
}
