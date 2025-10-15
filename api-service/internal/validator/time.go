package validator

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"api-service/internal/constants"

	"github.com/go-playground/validator/v10"
)

const hoursPerDay = 24

// ValidateTimeRange custom validator for time range business rules
// Validation rules:
// 1. Time format must be RFC3339 with URL decoding support (validates individual fields)
// 2. Start time must be before end time (validates range)
// 3. Time range cannot exceed the configured maximum duration (validates range)
func ValidateTimeRange(fl validator.FieldLevel) bool {
	// First, validate the current field's RFC3339 format
	fieldValue := fl.Field().String()
	if fieldValue != "" {
		// Parse with URL decoding support (similar to rfc3339 validator)
		_, err := parseTimeWithURLDecoding(fieldValue)
		if err != nil {
			return false // Invalid RFC3339 format
		}
	}

	// Get parent struct for range validation
	structValue := fl.Parent()
	if !structValue.IsValid() {
		return true // Skip range validation if struct is invalid
	}

	// Get StartTime and EndTime fields
	startTimeField := structValue.FieldByName("StartTime")
	endTimeField := structValue.FieldByName("EndTime")

	if !startTimeField.IsValid() || !endTimeField.IsValid() {
		return true // Skip range validation if fields don't exist
	}

	startTimeStr := startTimeField.String()
	endTimeStr := endTimeField.String()

	// Skip range validation if either field is empty
	if startTimeStr == "" || endTimeStr == "" {
		return true
	}

	// Parse start time with URL decoding support
	startTime, err := parseTimeWithURLDecoding(startTimeStr)
	if err != nil {
		return false
	}

	// Parse end time with URL decoding support
	endTime, err := parseTimeWithURLDecoding(endTimeStr)
	if err != nil {
		return false
	}

	// Validate: start time must be before end time
	if !startTime.Before(endTime) {
		return false
	}

	// Validate: time range cannot exceed maximum allowed duration
	maxDuration := GetMaxTimeRangeDuration()
	return endTime.Sub(startTime) <= maxDuration
}

// parseTimeWithURLDecoding parses time string with URL decoding support
func parseTimeWithURLDecoding(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	// Handle URL encoding:
	// 1. URL decode to handle %3A (colon) and other encoded characters
	// 2. Replace spaces with '+' since URL query parameters convert '+' to space
	decodedStr, err := url.QueryUnescape(dateStr)
	if err != nil {
		decodedStr = dateStr
	}

	// In URL query parameters, '+' becomes space, so we need to convert back
	// Check if the string looks like a datetime with spaces instead of '+'
	if strings.Contains(decodedStr, " ") && strings.Count(decodedStr, " ") == 1 {
		// Replace the space with '+' for timezone offset
		parts := strings.Split(decodedStr, " ")
		if len(parts) == 2 && len(parts[1]) >= 5 {
			// Check if the second part looks like timezone offset (e.g., "08:00")
			if matched, _ := regexp.MatchString(`^\d{2}:\d{2}$`, parts[1]); matched {
				decodedStr = parts[0] + "+" + parts[1]
			}
		}
	}

	return time.Parse(time.RFC3339, decodedStr)
}

// GetMaxTimeRangeDuration returns the maximum allowed time range duration
// This can be customized based on business requirements or configuration
func GetMaxTimeRangeDuration() time.Duration {
	return time.Duration(constants.MaxTimeRangeDays) * hoursPerDay * time.Hour
}

// ValidateCustomTimeRange validates time range with custom maximum duration
func ValidateCustomTimeRange(startTime, endTime time.Time, maxDays int) bool {
	if startTime.After(endTime) || startTime.Equal(endTime) {
		return false
	}

	maxDuration := time.Duration(maxDays) * hoursPerDay * time.Hour
	return endTime.Sub(startTime) <= maxDuration
}
