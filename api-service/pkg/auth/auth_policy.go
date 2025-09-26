package auth

import (
	"api-service/internal/config"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// AuthPolicy provides unified authentication policy validation and security checks
type AuthPolicy struct {
	authConfigManager *config.AuthConfigManager
	logger            logger.Logger
}

// NewAuthPolicy creates a new authentication policy manager
func NewAuthPolicy(authConfigManager *config.AuthConfigManager, logger logger.Logger) *AuthPolicy {
	return &AuthPolicy{
		authConfigManager: authConfigManager,
		logger:            logger,
	}
}

// Constants for validation
const (
	MinPasswordRequirements = 2   // Minimum number of password requirements to meet
	EmailPartsCount         = 2   // Number of parts in email address after @ split
	TimeFormatParts         = 2   // Expected number of parts in HH:MM format
	HoursToMinutesConv      = 100 // Conversion factor from hours to time int format

	// Username constraints
	UsernameMinLength = 3  // Minimum username length
	UsernameMaxLength = 20 // Maximum username length

	// Login methods
	LoginMethodUsername = "username"
	LoginMethodEmail    = "email"

	// Default TTL durations
	DefaultAttemptsCounterTTL = 60 * time.Minute // 60 minutes for attempts counter
)

// Allowed email domain whitelist (if empty, all domains are allowed)
var allowedEmailDomains = []string{
	// "company.com",
	// "websoft9.com",
}

// NewValidator creates a new validator instance
func NewValidator() *validator.Validate {
	return validator.New()
}

// ========== Username Validation ==========

// ValidateUsername validates username format and constraints according to security policies
// It checks for emptiness, length limits, character composition, and starting character requirements
func (ap *AuthPolicy) ValidateUsername(username string) error {
	ap.logger.Debug("Starting username validation", logger.String("username", username))

	// Check if empty - username is required for all operations
	if username == "" {
		ap.logger.Debug("Username validation failed: empty username")
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Check length constraints - enforce minimum security requirements
	if len(username) < UsernameMinLength || len(username) > UsernameMaxLength {
		ap.logger.Debug("Username validation failed: invalid length",
			logger.Int("length", len(username)),
			logger.Int("min_length", UsernameMinLength),
			logger.Int("max_length", UsernameMaxLength))
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Check character composition - only allow alphanumeric and underscores for security
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		ap.logger.Debug("Username validation failed: invalid characters detected", logger.String("username", username))
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Check starting character requirement - must start with letter for consistency
	if !unicode.IsLetter(rune(username[0])) {
		ap.logger.Debug("Username validation failed: does not start with letter",
			logger.String("first_char", string(username[0])))
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	ap.logger.Debug("Username validation successful", logger.String("username", username))
	return nil
}

// ========== Email Validation ==========

// ValidateEmail validates email format and domain constraints based on security policies
// It performs RFC-compliant format validation and optional domain whitelist checking
func (ap *AuthPolicy) ValidateEmail(email string) error {
	ap.logger.Debug("Starting email validation", logger.String("email", email))

	// Check if empty - email is required when email login is enabled
	if email == "" {
		ap.logger.Debug("Email validation failed: empty email")
		return errors.ErrInvalidEmailFormat
	}

	// Basic RFC-compliant format validation using regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		ap.logger.Debug("Email validation failed: invalid format", logger.String("email", email))
		return errors.ErrInvalidEmailFormat
	}

	// Domain whitelist validation - enforces organizational email policies if configured
	if len(allowedEmailDomains) > 0 {
		ap.logger.Debug("Checking email domain against whitelist",
			logger.String("email", email),
			logger.Any("allowed_domains", allowedEmailDomains))

		parts := strings.Split(email, "@")
		if len(parts) != EmailPartsCount {
			ap.logger.Debug("Email validation failed: invalid email structure", logger.Int("parts_count", len(parts)))
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
			ap.logger.Debug("Email validation failed: domain not in whitelist",
				logger.String("domain", domain),
				logger.Any("allowed_domains", allowedEmailDomains))
			return errors.NewAppError(errors.CodeValidationFailed)
		}

		ap.logger.Debug("Email domain whitelist check passed", logger.String("domain", domain))
	}

	ap.logger.Debug("Email validation successful", logger.String("email", email))
	return nil
}

// IsEmail checks if the given string is a valid email format
func (ap *AuthPolicy) IsEmail(input string) bool {
	if input == "" {
		return false
	}

	// Basic email format check
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(input)
}

// ValidateUsernameOrEmail validates that input is either a valid username or email
// It performs automatic detection and applies appropriate validation rules
func (ap *AuthPolicy) ValidateUsernameOrEmail(input string) error {
	ap.logger.Debug("Starting username or email validation", logger.String("input", input))

	if input == "" {
		ap.logger.Debug("Validation failed: empty input")
		return errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	// Determine input type by checking email pattern
	if ap.IsEmail(input) {
		ap.logger.Debug("Input detected as email, applying email validation", logger.String("input", input))
		return ap.ValidateEmail(input)
	}

	// Input is not email format, validate as username
	ap.logger.Debug("Input detected as username, applying username validation", logger.String("input", input))
	return ap.ValidateUsername(input)
}

// ========== Password Policy Validation ==========

// ValidatePasswordWithPolicy validates password based on the provided policy configuration
func (ap *AuthPolicy) ValidatePasswordWithPolicy(ctx context.Context, password string) error {
	policy := ap.authConfigManager.GetPasswordPolicy()

	if password == "" {
		ap.logger.WarnContext(ctx, "Password cannot be empty")
		return errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	// Check password length constraints
	if err := ap.validatePasswordLength(ctx, password, policy); err != nil {
		return err
	}

	// Check password complexity requirements
	return ap.validatePasswordComplexity(ctx, password, policy)
}

// validatePasswordLength validates password length against policy
func (ap *AuthPolicy) validatePasswordLength(ctx context.Context, password string, policy *config.PasswordPolicyConfig) error {
	passwordLen := len(password)

	if passwordLen < policy.MinLength {
		ap.logger.WarnContext(ctx, "Password too short", logger.Int("min_length", policy.MinLength), logger.Int("actual_length", passwordLen))
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_too_short")
	}

	if passwordLen > policy.MaxLength {
		ap.logger.WarnContext(ctx, "Password too long", logger.Int("max_length", policy.MaxLength), logger.Int("actual_length", passwordLen))
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_too_long")
	}

	return nil
}

// validatePasswordComplexity validates password complexity requirements
func (ap *AuthPolicy) validatePasswordComplexity(ctx context.Context, password string, policy *config.PasswordPolicyConfig) error {
	// Analyze password character types
	hasUpper, hasLower, hasNumber, hasSymbol := ap.analyzePasswordCharacters(password)

	// Check each requirement individually
	if policy.RequireUppercase && !hasUpper {
		ap.logger.WarnContext(ctx, "Password missing uppercase character")
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_require_uppercase")
	}

	if policy.RequireLowercase && !hasLower {
		ap.logger.WarnContext(ctx, "Password missing lowercase character")
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_require_lowercase")
	}

	if policy.RequireNumbers && !hasNumber {
		ap.logger.WarnContext(ctx, "Password missing number")
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_require_numbers")
	}

	if policy.RequireSymbols && !hasSymbol {
		ap.logger.WarnContext(ctx, "Password missing symbol")
		return errors.NewAppErrorWithI18n(errors.CodeValidationFailed, "auth.password_require_symbols")
	}

	ap.logger.InfoContext(ctx, "Password policy validation successful")
	return nil
}

// analyzePasswordCharacters determines which character types are present in password
func (ap *AuthPolicy) analyzePasswordCharacters(password string) (hasUpper, hasLower, hasNumber, hasSymbol bool) {
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}
	return hasUpper, hasLower, hasNumber, hasSymbol
}

// ========== Login Method Validation ==========

// ValidateLoginMethod checks if the input login method is allowed based on configuration
func (ap *AuthPolicy) ValidateLoginMethod(ctx context.Context, username string) error {
	// Check configured login methods and validate input accordingly
	loginMethods := ap.authConfigManager.GetConfig().UserAuth.BasicAuth.LoginMethods
	// If LoginMethods is empty or nil, default to both username and email
	if len(loginMethods) == 0 {
		loginMethods = []string{LoginMethodUsername, LoginMethodEmail}
	}

	// Determine if input is email or username
	isEmailInput := ap.IsEmail(username)
	inputType := LoginMethodUsername
	if isEmailInput {
		inputType = LoginMethodEmail
	}

	// Check if the input type is allowed
	allowedMethod := false
	for _, method := range loginMethods {
		if strings.EqualFold(method, inputType) {
			allowedMethod = true
			break
		}
	}

	if !allowedMethod {
		ap.logger.WarnContext(ctx, "Login method not allowed",
			logger.String("input_type", inputType),
			logger.Any("allowed_methods", loginMethods),
			logger.String("username", username))
		loginErr := errors.ErrEmailSuported
		if len(loginMethods) > 0 && loginMethods[0] == LoginMethodUsername {
			loginErr = errors.ErrUsernameSuported
		}
		return loginErr
	}

	// Validate username or email format based on input type
	if isEmailInput {
		if err := ap.ValidateEmail(username); err != nil {
			ap.logger.WarnContext(ctx, "Invalid email format",
				logger.String("email", username), logger.ErrorField(err))
			return errors.ErrInvalidCredentials
		}
	} else {
		if err := ap.ValidateUsernameOrEmail(username); err != nil {
			ap.logger.WarnContext(ctx, "Invalid username format",
				logger.String("username", username), logger.ErrorField(err))
			return errors.ErrInvalidCredentials
		}
	}

	return nil
}

// ========== Login Security Checks ==========

// CheckLoginSecurity performs comprehensive login security checks in sequential order
// This includes account lockout status, IP whitelist, time restrictions, and login method validation
// Each check is designed to fail fast for better performance and security
func (ap *AuthPolicy) CheckLoginSecurity(ctx context.Context, username, clientIP string) error {
	ap.logger.DebugContext(ctx, "Starting comprehensive login security checks",
		logger.String("username", username),
		logger.String("client_ip", clientIP))

	// 1. Check account lockout status first - fastest check to prevent unnecessary processing
	ap.logger.DebugContext(ctx, "Checking account lockout status", logger.String("username", username))
	if err := ap.checkAccountLockout(ctx, username); err != nil {
		ap.logger.DebugContext(ctx, "Account lockout check failed", logger.String("username", username), logger.ErrorField(err))
		return err
	}

	// 2. Check IP whitelist if enabled - network-level security
	ap.logger.DebugContext(ctx, "Checking IP whitelist",
		logger.String("username", username),
		logger.String("client_ip", clientIP))
	if err := ap.checkIPWhitelist(ctx, username, clientIP); err != nil {
		ap.logger.DebugContext(ctx, "IP whitelist check failed",
			logger.String("username", username),
			logger.String("client_ip", clientIP),
			logger.ErrorField(err))
		return err
	}

	// 3. Check login time restrictions if enabled - temporal access control
	ap.logger.DebugContext(ctx, "Checking login time restrictions", logger.String("username", username))
	if err := ap.checkLoginTimeRestrictions(ctx, username); err != nil {
		ap.logger.DebugContext(ctx, "Login time restriction check failed",
			logger.String("username", username),
			logger.ErrorField(err))
		return err
	}

	// 4. Validate login method and input format - application-level validation
	ap.logger.DebugContext(ctx, "Validating login method and format", logger.String("username", username))
	if err := ap.ValidateLoginMethod(ctx, username); err != nil {
		ap.logger.DebugContext(ctx, "Login method validation failed",
			logger.String("username", username),
			logger.ErrorField(err))
		return err
	}

	ap.logger.DebugContext(ctx, "All login security checks passed successfully",
		logger.String("username", username),
		logger.String("client_ip", clientIP))
	return nil
}

// ValidatePassword validates password and updates login attempts accordingly
func (ap *AuthPolicy) ValidatePassword(ctx context.Context, username string, passwordValid bool) error {
	if !passwordValid {
		// Record password validation failure
		ap.recordLoginFailure(ctx, username, "password_error")
		return errors.NewAppError(errors.CodeInvalidCredentials)
	}

	// Password is valid, clear login attempts counter
	ap.clearLoginAttempts(ctx, username)
	return nil
}

// checkIPWhitelist validates client IP against configured whitelist for network-level access control
// Supports both individual IP addresses and CIDR notation for network ranges
func (ap *AuthPolicy) checkIPWhitelist(ctx context.Context, username, clientIP string) error {
	security := ap.authConfigManager.GetLoginSecurity()
	ap.logger.DebugContext(ctx, "Checking IP whitelist configuration",
		logger.Bool("whitelist_enabled", security.IPWhitelistEnabled),
		logger.Int("whitelist_count", len(security.IPWhitelist)))

	// Skip IP whitelist check if disabled or no whitelist configured
	if !security.IPWhitelistEnabled || len(security.IPWhitelist) == 0 {
		ap.logger.DebugContext(ctx, "IP whitelist check skipped - disabled or empty whitelist")
		return nil
	}

	ap.logger.DebugContext(ctx, "Validating client IP against whitelist",
		logger.String("client_ip", clientIP),
		logger.Any("whitelist", security.IPWhitelist))

	// Check each whitelisted IP/CIDR against client IP
	allowed := false
	for _, allowedIP := range security.IPWhitelist {
		if ap.isIPAllowed(ctx, allowedIP, clientIP) {
			ap.logger.DebugContext(ctx, "Client IP matched whitelist entry",
				logger.String("client_ip", clientIP),
				logger.String("matched_entry", allowedIP))
			allowed = true
			break
		}
	}

	// Deny access if IP not in whitelist
	if !allowed {
		// Record security violation for audit and lockout logic
		ap.logger.DebugContext(ctx, "Recording IP whitelist violation",
			logger.String("username", username),
			logger.String("client_ip", clientIP))
		ap.recordLoginFailure(ctx, username, "ip_whitelist_error")
		ap.logger.WarnContext(ctx, "Access denied - IP address not in whitelist",
			logger.String("client_ip", clientIP),
			logger.String("username", username))
		return errors.NewAppErrorWithI18n(errors.CodeAccessDenied, "auth.ip_not_allowed")
	}

	ap.logger.DebugContext(ctx, "IP whitelist validation successful",
		logger.String("client_ip", clientIP))
	return nil
}

// isIPAllowed checks if client IP matches the allowed IP entry with support for CIDR notation
// Returns true if the client IP is within the allowed range/address
func (ap *AuthPolicy) isIPAllowed(ctx context.Context, allowedIP, clientIP string) bool {
	ap.logger.DebugContext(ctx, "Checking IP match",
		logger.String("allowed_ip", allowedIP),
		logger.String("client_ip", clientIP))

	// Handle CIDR notation for network ranges (e.g., 192.168.1.0/24)
	if strings.Contains(allowedIP, "/") {
		ap.logger.DebugContext(ctx, "Processing CIDR notation", logger.String("cidr", allowedIP))
		_, network, err := net.ParseCIDR(allowedIP)
		if err != nil {
			ap.logger.WarnContext(ctx, "Invalid CIDR format in IP whitelist",
				logger.String("cidr", allowedIP),
				logger.ErrorField(err))
			return false
		}
		contains := network.Contains(net.ParseIP(clientIP))
		ap.logger.DebugContext(ctx, "CIDR range check result",
			logger.String("cidr", allowedIP),
			logger.String("client_ip", clientIP),
			logger.Bool("contains", contains))
		return contains
	}

	// Direct IP address comparison for exact matches
	match := clientIP == allowedIP
	ap.logger.DebugContext(ctx, "Direct IP comparison result",
		logger.String("allowed_ip", allowedIP),
		logger.String("client_ip", clientIP),
		logger.Bool("match", match))
	return match
}

// checkLoginTimeRestrictions validates login time against configured restrictions
func (ap *AuthPolicy) checkLoginTimeRestrictions(ctx context.Context, username string) error {
	security := ap.authConfigManager.GetLoginSecurity()
	if !security.LoginTimeRestriction || security.AllowedLoginHours == "" {
		return nil
	}

	now := time.Now()
	currentTime := now.Hour()*HoursToMinutesConv + now.Minute()

	// Parse allowed time range (format: "HH:MM-HH:MM")
	timeParts := strings.Split(security.AllowedLoginHours, "-")
	if len(timeParts) != TimeFormatParts {
		return nil
	}

	startTime, startErr := ap.parseTimeString(timeParts[0])
	endTime, endErr := ap.parseTimeString(timeParts[1])

	if startErr != nil || endErr != nil {
		ap.logger.WarnContext(ctx, "Invalid time format in login restrictions", logger.String("allowed_hours", security.AllowedLoginHours))
		return nil
	}

	if currentTime < startTime || currentTime > endTime {
		// Record time restriction validation failure
		ap.recordLoginFailure(ctx, username, "time_restriction_error")
		ap.logger.WarnContext(ctx, "Login attempt outside allowed hours",
			logger.Int("current_time", currentTime),
			logger.String("allowed_hours", security.AllowedLoginHours))
		return errors.NewAppErrorWithI18n(errors.CodeAccessDenied, "auth.login_time_restricted")
	}

	return nil
}

// parseTimeString parses time string in HH:MM format to minutes since midnight
func (ap *AuthPolicy) parseTimeString(timeStr string) (int, error) {
	timeParts := strings.Split(strings.TrimSpace(timeStr), ":")
	if len(timeParts) != TimeFormatParts {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time: %02d:%02d", hour, minute)
	}

	return hour*HoursToMinutesConv + minute, nil
}

// ========== Redis-based Login Security Methods ==========

// checkAccountLockout checks if the account is currently locked due to failed login attempts
func (ap *AuthPolicy) checkAccountLockout(ctx context.Context, username string) error {
	lockerKey := redis.FormatRedisKey(redis.RK_LOGIN_LOCKER, username)

	// Check if lockout key exists
	exists, err := redis.Exists(ctx, lockerKey)
	if err != nil {
		ap.logger.ErrorContext(ctx, "Failed to check account lockout status",
			logger.String("username", username),
			logger.ErrorField(err))
		// Continue without Redis check
		return nil
	}

	if exists > 0 {
		// Account is locked, get remaining TTL
		ttl, err := redis.TTL(ctx, lockerKey)
		if err != nil {
			ap.logger.ErrorContext(ctx, "Failed to get lockout TTL",
				logger.String("username", username),
				logger.ErrorField(err))
		}

		ap.logger.WarnContext(ctx, "Account is locked due to too many failed login attempts",
			logger.String("username", username),
			logger.Duration("remaining_lockout_time", ttl))

		return errors.NewAppErrorWithI18n(errors.CodeLoginAttemptsExceeded, "auth.account_locked")
	}

	return nil
}

// recordLoginFailure records a login failure and handles progressive lockout logic
// It implements a fail-safe mechanism with automatic lockout after maximum attempts
func (ap *AuthPolicy) recordLoginFailure(ctx context.Context, username, failureType string) {
	security := ap.authConfigManager.GetLoginSecurity()
	attemptsKey := redis.FormatRedisKey(redis.RK_LOGIN_ATTEMPTS, username)
	lockerKey := redis.FormatRedisKey(redis.RK_LOGIN_LOCKER, username)

	ap.logger.DebugContext(ctx, "Recording login failure",
		logger.String("username", username),
		logger.String("failure_type", failureType),
		logger.Int("max_attempts_allowed", security.MaxLoginAttempts))

	// Atomically increment failure counter in Redis
	currentAttempts, err := redis.IncrBy(ctx, attemptsKey, 1)
	if err != nil {
		ap.logger.ErrorContext(ctx, "Failed to increment login attempts counter - continuing without Redis tracking",
			logger.String("username", username),
			logger.String("failure_type", failureType),
			logger.ErrorField(err))
		return
	}

	ap.logger.DebugContext(ctx, "Login attempts counter incremented",
		logger.String("username", username),
		logger.Int64("current_attempts", currentAttempts))

	// Set TTL for attempts counter on first failure to auto-expire tracking
	if currentAttempts == 1 {
		ap.logger.DebugContext(ctx, "Setting TTL for new attempts counter",
			logger.String("username", username),
			logger.Duration("ttl", DefaultAttemptsCounterTTL))
		if err := redis.Expire(ctx, attemptsKey, DefaultAttemptsCounterTTL); err != nil {
			ap.logger.ErrorContext(ctx, "Failed to set TTL for login attempts counter",
				logger.String("username", username),
				logger.ErrorField(err))
		}
	}

	ap.logger.InfoContext(ctx, "Login failure recorded in security system",
		logger.String("username", username),
		logger.String("failure_type", failureType),
		logger.Int64("current_attempts", currentAttempts),
		logger.Int("max_attempts", security.MaxLoginAttempts))

	// Trigger account lockout if maximum attempts threshold reached
	if currentAttempts >= int64(security.MaxLoginAttempts) {
		ap.logger.DebugContext(ctx, "Maximum login attempts reached - initiating account lockout",
			logger.String("username", username),
			logger.Int64("failed_attempts", currentAttempts),
			logger.Int("max_attempts", security.MaxLoginAttempts))

		// Calculate lockout duration and create lockout record
		lockoutDuration := time.Duration(security.LockoutDuration) * time.Second
		createdTime := time.Now().Format(time.RFC3339)

		// Set account lockout flag with expiration
		if err := redis.Set(ctx, lockerKey, createdTime, lockoutDuration); err != nil {
			ap.logger.ErrorContext(ctx, "Critical failure - unable to set account lockout in Redis",
				logger.String("username", username),
				logger.Duration("lockout_duration", lockoutDuration),
				logger.ErrorField(err))
			return
		}

		ap.logger.DebugContext(ctx, "Account lockout flag set successfully",
			logger.String("username", username),
			logger.String("lockout_key", lockerKey),
			logger.Duration("lockout_duration", lockoutDuration))

		// Clean up attempts counter since account is now locked
		if _, err := redis.Del(ctx, attemptsKey); err != nil {
			ap.logger.WarnContext(ctx, "Failed to clear attempts counter after lockout - not critical",
				logger.String("username", username),
				logger.ErrorField(err))
		} else {
			ap.logger.DebugContext(ctx, "Attempts counter cleared after lockout",
				logger.String("username", username))
		}

		// Log security event for audit trail
		ap.logger.WarnContext(ctx, "SECURITY EVENT: Account automatically locked due to excessive failed login attempts",
			logger.String("username", username),
			logger.Int64("failed_attempts", currentAttempts),
			logger.Duration("lockout_duration", lockoutDuration),
			logger.String("failure_type", failureType))
	}
}

// clearLoginAttempts clears the login attempts counter for a user
func (ap *AuthPolicy) clearLoginAttempts(ctx context.Context, username string) {
	attemptsKey := redis.FormatRedisKey(redis.RK_LOGIN_ATTEMPTS, username)

	// Delete the attempts counter
	deleted, err := redis.Del(ctx, attemptsKey)
	if err != nil {
		ap.logger.ErrorContext(ctx, "Failed to clear login attempts counter",
			logger.String("username", username),
			logger.ErrorField(err))
		return
	}

	if deleted > 0 {
		ap.logger.InfoContext(ctx, "Login attempts counter cleared after successful login",
			logger.String("username", username))
	}
}

// GetLoginAttempts returns the current number of login attempts for a user
func (ap *AuthPolicy) GetLoginAttempts(ctx context.Context, username string) (int64, error) {
	attemptsKey := redis.FormatRedisKey(redis.RK_LOGIN_ATTEMPTS, username)

	attempts, err := redis.Get(ctx, attemptsKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			// Key doesn't exist, return 0 attempts
			return 0, nil
		}
		return 0, err
	}

	// Convert string to int64
	count, err := strconv.ParseInt(attempts, 10, 64)
	if err != nil {
		ap.logger.ErrorContext(ctx, "Failed to parse login attempts count",
			logger.String("username", username),
			logger.String("attempts_value", attempts),
			logger.ErrorField(err))
		return 0, err
	}

	return count, nil
}

// IsAccountLocked checks if an account is currently locked
func (ap *AuthPolicy) IsAccountLocked(ctx context.Context, username string) (bool, time.Duration, error) {
	lockerKey := redis.FormatRedisKey(redis.RK_LOGIN_LOCKER, username)

	exists, err := redis.Exists(ctx, lockerKey)
	if err != nil {
		return false, 0, err
	}

	if exists == 0 {
		return false, 0, nil
	}

	// Get remaining TTL
	ttl, err := redis.TTL(ctx, lockerKey)
	if err != nil {
		return true, 0, err
	}

	return true, ttl, nil
}

// ========== Global Package Functions for Backward Compatibility ==========

// Global instance for backward compatibility
var globalAuthPolicy *AuthPolicy

// InitGlobalAuthPolicy initializes the global auth policy instance
func InitGlobalAuthPolicy(authConfigManager *config.AuthConfigManager, logger logger.Logger) {
	globalAuthPolicy = NewAuthPolicy(authConfigManager, logger)
}

// Package-level functions for backward compatibility

// ValidateEmail validates email format using global policy
func ValidateEmail(email string) error {
	if globalAuthPolicy == nil {
		// Fallback to basic email validation
		if email == "" {
			return errors.ErrInvalidEmailFormat
		}
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return errors.ErrInvalidEmailFormat
		}
		return nil
	}
	return globalAuthPolicy.ValidateEmail(email)
}

// IsEmail checks if string is email using global policy
func IsEmail(input string) bool {
	if globalAuthPolicy == nil {
		// Fallback to basic check
		if input == "" {
			return false
		}
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		return emailRegex.MatchString(input)
	}
	return globalAuthPolicy.IsEmail(input)
}

// ValidateUsernameOrEmail validates username or email using global policy
func ValidateUsernameOrEmail(input string) error {
	if globalAuthPolicy == nil {
		return errors.NewAppError(errors.CodeRequiredParameterMissing)
	}
	return globalAuthPolicy.ValidateUsernameOrEmail(input)
}
