package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	// Password policy constants
	DefaultMinPasswordLength = 8
	DefaultMaxPasswordLength = 128

	// Token generation constants
	APITokenRandomBytesSize = 32
)

// HashPassword encrypts password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPasswordHash verifies password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// PasswordPolicy password policy
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	MaxLength        int  `json:"max_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireLowercase bool `json:"require_lowercase"`
	RequireNumbers   bool `json:"require_numbers"`
	RequireSymbols   bool `json:"require_symbols"`
}

// DefaultPasswordPolicy default password policy
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:        DefaultMinPasswordLength,
	MaxLength:        DefaultMaxPasswordLength,
	RequireUppercase: true,
	RequireLowercase: true,
	RequireNumbers:   true,
	RequireSymbols:   false,
}

// ValidatePassword validates if password meets policy requirements
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = &DefaultPasswordPolicy
	}

	// Check length
	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}
	if len(password) > policy.MaxLength {
		return fmt.Errorf("password must be no more than %d characters long", policy.MaxLength)
	}

	// Check uppercase letters
	if policy.RequireUppercase {
		matched, _ := regexp.MatchString(`[A-Z]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one uppercase letter")
		}
	}

	// Check lowercase letters
	if policy.RequireLowercase {
		matched, _ := regexp.MatchString(`[a-z]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one lowercase letter")
		}
	}

	// Check numbers
	if policy.RequireNumbers {
		matched, _ := regexp.MatchString(`[0-9]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one number")
		}
	}

	// Check special characters
	if policy.RequireSymbols {
		matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one special character")
		}
	}

	return nil
}

// GenerateRandomPassword generates random password
func GenerateRandomPassword(length int, includeSymbols bool) (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers   = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	charset := lowercase + uppercase + numbers
	if includeSymbols {
		charset += symbols
	}

	password := make([]byte, length)
	for i := range password {
		randomIndex, err := randomInt(len(charset))
		if err != nil {
			return "", fmt.Errorf("failed to generate random password: %w", err)
		}
		password[i] = charset[randomIndex]
	}

	return string(password), nil
}

// randomInt generates random integer
func randomInt(maxVal int) (int, error) {
	bytes := make([]byte, 1)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, err
	}
	return int(bytes[0]) % maxVal, nil
}

// GenerateAPIToken generates API token
func GenerateAPIToken(prefix string) (token, hashedToken string, err error) {
	// Generate 32 bytes of random data
	randomBytes := make([]byte, APITokenRandomBytesSize)
	if _, err = rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert to hexadecimal string
	tokenSuffix := hex.EncodeToString(randomBytes)

	// Combine complete token
	token = fmt.Sprintf("%s_%s", prefix, tokenSuffix)

	// Generate token hash for storage
	hash := sha256.Sum256([]byte(token))
	hashedToken = hex.EncodeToString(hash[:])

	return token, hashedToken, nil
}

// HashAPIToken hashes API token
func HashAPIToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
