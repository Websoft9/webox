package validator

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// credentialNamePattern is the regex pattern for credential name validation
var credentialNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)

// ValidateCredentialName validates if a credential name follows the required format
// Credential names must:
// - Be 3-50 characters long
// - Contain only alphanumeric characters and underscores
// - No special characters or spaces allowed
func ValidateCredentialName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" {
		return true // omitempty will handle empty strings
	}
	return credentialNamePattern.MatchString(name)
}

// RegisterCredentialValidators registers all credential-related validation rules
// This function should be called during application initialization
//
// Parameters:
//   - validatorInstance: The validator.Validate instance to register the rules with
//
// Returns:
//   - error: Any error encountered during registration
func RegisterCredentialValidators(validatorInstance *validator.Validate) error {
	// Register the alphanum_underscore validation rule for credential names
	err := validatorInstance.RegisterValidation("alphanum_underscore", ValidateCredentialName)
	if err != nil {
		return fmt.Errorf("failed to register alphanum_underscore validator: %w", err)
	}

	return nil
}
