package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// CodeLength is the length of the random string in resource code
	CodeLength = 10
	// CodeCharset is the character set for generating resource codes
	CodeCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
	// MaxRetryAttempts is the maximum number of retry attempts for generating unique codes
	MaxRetryAttempts = 3
)

// GenerateCode generates a unique resource code with format "resource_" + 10-character random string
// The random string consists of lowercase letters (a-z) and digits (0-9)
// Uses cryptographically secure random number generator
func GenerateCode(resourceType string) (string, error) {
	randomStr, err := generateRandomString(CodeLength, resourceType)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return fmt.Sprintf("%s_%s", resourceType, randomStr), nil
}

// generateRandomString generates a cryptographically secure random string
// of the specified length using the provided character set
func generateRandomString(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	if charset == "" {
		return "", fmt.Errorf("charset cannot be empty")
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// GenerateCodeWithRetry generates a unique resource code with retry mechanism
// It attempts to generate a code up to maxRetries times
// The uniqueCheck function should return true if the code is unique
func GenerateCodeWithRetry(resourceType string, uniqueCheck func(string) (bool, error)) (string, error) {
	var lastErr error

	for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
		code, err := GenerateCode(resourceType)
		if err != nil {
			lastErr = err
			continue
		}

		// Check if code is unique
		isUnique, err := uniqueCheck(code)
		if err != nil {
			lastErr = err
			continue
		}

		if isUnique {
			return code, nil
		}

		lastErr = fmt.Errorf("generated code is not unique")
	}

	return "", fmt.Errorf("failed to generate unique resource code after %d attempts: %w", MaxRetryAttempts, lastErr)
}
