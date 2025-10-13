package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers all custom validators for request validation
func RegisterCustomValidators(validatorInstance *validator.Validate) error {
	// Register time range validator
	err := validatorInstance.RegisterValidation("time_range", ValidateTimeRange)
	if err != nil {
		return fmt.Errorf("failed to register time_range validator: %w", err)
	}

	return nil
}
