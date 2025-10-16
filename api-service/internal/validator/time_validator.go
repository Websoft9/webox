package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"api-service/internal/constants"

	"github.com/go-playground/validator/v10"
)

// ValidateTimeRange validates if a time string can be parsed correctly
// It accepts empty strings (handled by omitempty tag) and validates non-empty strings
// using ParseTimeWithURLDecoding to handle URL-encoded time formats
func ValidateTimeRange(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // omitempty will handle empty strings
	}
	// Attempt to parse the time string with URL decoding support
	_, err := ParseTimeWithURLDecoding(dateStr)
	return err == nil
}

// ParseTimeWithURLDecoding parses a time string with URL decoding support
// It handles URL-encoded characters and converts URL query parameter format to standard time format
//
// The function performs the following steps:
//  1. Returns zero time for empty strings
//  2. URL decodes the string to handle encoded characters like %3A (colon)
//  3. Converts space-separated timezone format to plus-sign format
//     Example: "2024-01-01T12:00:00 08:00" -> "2024-01-01T12:00:00+08:00"
//  4. Parses the time using the default time format from constants
//
// Parameters:
//   - dateStr: The time string to parse, may be URL-encoded
//
// Returns:
//   - time.Time: The parsed time value
//   - error: Any error encountered during parsing
func ParseTimeWithURLDecoding(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	// Step 1: URL decode to handle %3A (colon) and other encoded characters
	// If decoding fails, use the original string
	decodedStr, err := url.QueryUnescape(dateStr)
	if err != nil {
		decodedStr = dateStr
	}

	// Step 2: Handle timezone offset format conversion
	// In URL query parameters, '+' becomes space, so we need to convert back
	// Check if the string looks like a datetime with spaces instead of '+'
	if strings.Contains(decodedStr, " ") && strings.Count(decodedStr, " ") == 1 {
		// Split by space to separate datetime and timezone offset
		parts := strings.Split(decodedStr, " ")
		if len(parts) == 2 && len(parts[1]) >= 5 {
			// Check if the second part looks like timezone offset (e.g., "08:00")
			// Pattern matches: two digits, colon, two digits
			if matched, _ := regexp.MatchString(`^\d{2}:\d{2}$`, parts[1]); matched {
				// Reconstruct with '+' sign for timezone offset
				decodedStr = parts[0] + "+" + parts[1]
			}
		}
	}

	// Step 3: Parse the time string using the default time format
	return time.Parse(constants.DefaultTimeFormat, decodedStr)
}

// RegisterTimeRangeValidator registers the time_range validation rule with the validator instance
// This function should be called during application initialization to make the time_range
// validation tag available for use in struct field validation
//
// Parameters:
//   - validatorInstance: The validator.Validate instance to register the rule with
//
// Returns:
//   - error: Any error encountered during registration
func RegisterTimeRangeValidator(validatorInstance *validator.Validate) error {
	// Register the time_range validation rule
	err := validatorInstance.RegisterValidation("time_range", ValidateTimeRange)
	if err != nil {
		return fmt.Errorf("failed to register time_range validator: %w", err)
	}

	return nil
}
