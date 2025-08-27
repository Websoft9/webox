package service

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/service"
	"api-service/pkg/logger"
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// AuthConfigServiceImpl implements the AuthConfigService interface
type AuthConfigServiceImpl struct {
	authConfigManager *config.AuthConfigManager
	logger            logger.Logger
}

// NewAuthConfigService creates a new authentication configuration service instance
func NewAuthConfigService(authConfigManager *config.AuthConfigManager, logger logger.Logger) service.AuthConfigService {
	return &AuthConfigServiceImpl{
		authConfigManager: authConfigManager,
		logger:            logger,
	}
}

// GetAuthConfig retrieves the current authentication configuration
func (s *AuthConfigServiceImpl) GetAuthConfig(ctx context.Context) (*response.AuthConfigResponse, error) {
	s.logger.InfoContext(ctx, "Getting authentication configuration")

	authConfig := s.authConfigManager.GetConfig()
	if authConfig == nil {
		s.logger.ErrorContext(ctx, "Authentication configuration not found")
		return nil, fmt.Errorf("authentication configuration not found")
	}

	// Convert internal config to response format
	resp := &response.AuthConfigResponse{
		APIAuth: response.APIAuthResponse{
			TokenAuthEnabled: authConfig.APIAuth.TokenAuth.Enabled,
			OAuth2Enabled:    authConfig.APIAuth.OAuth2.Enabled,
			JWTConfig: response.JWTConfigResponse{
				Algorithm:        authConfig.APIAuth.TokenAuth.Algorithm,
				ExpiresIn:        authConfig.APIAuth.TokenAuth.ExpiresIn,
				RefreshExpiresIn: authConfig.APIAuth.TokenAuth.RefreshExpiresIn,
				AutoRefresh:      authConfig.APIAuth.TokenAuth.AutoRefresh,
			},
		},
		UserAuth: response.UserAuthResponse{
			OAuth2Enabled:          authConfig.UserAuth.OAuth2.Enabled,
			OAuth2Providers:        s.convertOAuth2ProvidersToResponse(authConfig.UserAuth.OAuth2.Providers),
			TwoFactorEnabled:       authConfig.UserAuth.TwoFactor.Enabled,
			TwoFactorMethods:       []string{},
			TwoFactorRequiredRoles: authConfig.UserAuth.TwoFactor.RequiredRoles,
			PasswordPolicy: response.PasswordPolicyResponse{
				MinLength:           authConfig.UserAuth.PasswordPolicy.MinLength,
				MaxLength:           authConfig.UserAuth.PasswordPolicy.MaxLength,
				RequireUppercase:    authConfig.UserAuth.PasswordPolicy.RequireUppercase,
				RequireLowercase:    authConfig.UserAuth.PasswordPolicy.RequireLowercase,
				RequireNumbers:      authConfig.UserAuth.PasswordPolicy.RequireNumbers,
				RequireSymbols:      authConfig.UserAuth.PasswordPolicy.RequireSymbols,
				PasswordHistory:     authConfig.UserAuth.PasswordPolicy.PasswordHistory,
				PasswordExpiresDays: authConfig.UserAuth.PasswordPolicy.PasswordExpiresDays,
			},
			LoginSecurity: response.LoginSecurityResponse{
				MaxLoginAttempts:     authConfig.UserAuth.LoginSecurity.MaxLoginAttempts,
				LockoutDuration:      authConfig.UserAuth.LoginSecurity.LockoutDuration,
				IPWhitelistEnabled:   authConfig.UserAuth.LoginSecurity.IPWhitelistEnabled,
				IPWhitelist:          authConfig.UserAuth.LoginSecurity.IPWhitelist,
				LoginTimeRestriction: authConfig.UserAuth.LoginSecurity.LoginTimeRestriction,
				AllowedLoginHours:    authConfig.UserAuth.LoginSecurity.AllowedLoginHours,
			},
		},
		SessionConfig: response.SessionConfigResponse{
			Timeout:               authConfig.SessionConfig.Timeout,
			MaxConcurrentSessions: authConfig.SessionConfig.MaxConcurrentSessions,
			RememberMeEnabled:     authConfig.SessionConfig.RememberMeEnabled,
			RememberMeDuration:    authConfig.SessionConfig.RememberMeDuration,
		},
	}

	// Set two-factor authentication methods
	if authConfig.UserAuth.TwoFactor.Methods.TOTP.Enabled {
		resp.UserAuth.TwoFactorMethods = append(resp.UserAuth.TwoFactorMethods, "TOTP")
	}
	if authConfig.UserAuth.TwoFactor.Methods.Email.Enabled {
		resp.UserAuth.TwoFactorMethods = append(resp.UserAuth.TwoFactorMethods, "EMAIL")
	}

	s.logger.InfoContext(ctx, "Successfully retrieved authentication configuration")
	return resp, nil
}

// UpdateAuthConfig updates the authentication configuration
func (s *AuthConfigServiceImpl) UpdateAuthConfig(ctx context.Context, req *request.UpdateAuthConfigRequest) error {
	s.logger.InfoContext(ctx, "Updating authentication configuration")

	// Get current configuration
	currentConfig := s.authConfigManager.GetConfig()
	if currentConfig == nil {
		return fmt.Errorf("current authentication configuration not found")
	}

	// Create updated configuration
	updatedConfig := *currentConfig

	// Update different configuration sections
	s.updateAPIAuthConfig(&updatedConfig, req.APIAuth)
	s.updateUserAuthConfig(&updatedConfig, req.UserAuth)
	s.updateSessionConfig(&updatedConfig, req.SessionConfig)

	// Update configuration in manager
	if err := s.authConfigManager.UpdateConfig(&updatedConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update authentication configuration", logger.ErrorField(err))
		return fmt.Errorf("failed to update authentication configuration: %w", err)
	}

	s.logger.InfoContext(ctx, "Successfully updated authentication configuration")
	return nil
}

// updateAPIAuthConfig updates API authentication configuration section
func (s *AuthConfigServiceImpl) updateAPIAuthConfig(config *config.AuthConfig, apiAuth *request.APIAuthRequest) {
	if apiAuth == nil {
		return
	}

	if apiAuth.TokenAuthEnabled != nil {
		config.APIAuth.TokenAuth.Enabled = *apiAuth.TokenAuthEnabled
	}
	if apiAuth.OAuth2Enabled != nil {
		config.APIAuth.OAuth2.Enabled = *apiAuth.OAuth2Enabled
	}

	s.updateJWTConfig(&config.APIAuth.TokenAuth, apiAuth.JWTConfig)
}

// updateJWTConfig updates JWT configuration
func (s *AuthConfigServiceImpl) updateJWTConfig(tokenAuth *config.TokenAuthConfig, jwtConfig *request.JWTConfigRequest) {
	if jwtConfig == nil {
		return
	}

	if jwtConfig.Algorithm != nil {
		tokenAuth.Algorithm = *jwtConfig.Algorithm
	}
	if jwtConfig.ExpiresIn != nil {
		tokenAuth.ExpiresIn = *jwtConfig.ExpiresIn
	}
	if jwtConfig.RefreshExpiresIn != nil {
		tokenAuth.RefreshExpiresIn = *jwtConfig.RefreshExpiresIn
	}
	if jwtConfig.AutoRefresh != nil {
		tokenAuth.AutoRefresh = *jwtConfig.AutoRefresh
	}
}

// updateUserAuthConfig updates user authentication configuration section
func (s *AuthConfigServiceImpl) updateUserAuthConfig(config *config.AuthConfig, userAuth *request.UserAuthRequest) {
	if userAuth == nil {
		return
	}

	if userAuth.OAuth2Enabled != nil {
		config.UserAuth.OAuth2.Enabled = *userAuth.OAuth2Enabled
	}
	if userAuth.TwoFactorEnabled != nil {
		config.UserAuth.TwoFactor.Enabled = *userAuth.TwoFactorEnabled
	}

	s.updateTwoFactorMethods(&config.UserAuth.TwoFactor, userAuth.TwoFactorMethods)

	if userAuth.TwoFactorRequiredRoles != nil {
		config.UserAuth.TwoFactor.RequiredRoles = *userAuth.TwoFactorRequiredRoles
	}
	if userAuth.PasswordPolicy != nil {
		s.updatePasswordPolicy(&config.UserAuth.PasswordPolicy, userAuth.PasswordPolicy)
	}
	if userAuth.LoginSecurity != nil {
		s.updateLoginSecurity(&config.UserAuth.LoginSecurity, userAuth.LoginSecurity)
	}
}

// updateTwoFactorMethods updates two-factor authentication methods
func (s *AuthConfigServiceImpl) updateTwoFactorMethods(twoFactor *config.TwoFactorConfig, methods *[]string) {
	if methods == nil {
		return
	}

	// Update two-factor methods based on request
	twoFactor.Methods.TOTP.Enabled = contains(*methods, "TOTP")
	twoFactor.Methods.Email.Enabled = contains(*methods, "EMAIL")
}

// updateSessionConfig updates session configuration section
func (s *AuthConfigServiceImpl) updateSessionConfig(config *config.AuthConfig, sessionConfig *request.SessionConfigRequest) {
	if sessionConfig == nil {
		return
	}

	if sessionConfig.Timeout != nil {
		config.SessionConfig.Timeout = *sessionConfig.Timeout
	}
	if sessionConfig.MaxConcurrentSessions != nil {
		config.SessionConfig.MaxConcurrentSessions = *sessionConfig.MaxConcurrentSessions
	}
	if sessionConfig.RememberMeEnabled != nil {
		config.SessionConfig.RememberMeEnabled = *sessionConfig.RememberMeEnabled
	}
	if sessionConfig.RememberMeDuration != nil {
		config.SessionConfig.RememberMeDuration = *sessionConfig.RememberMeDuration
	}
}

// GetOAuth2Providers retrieves OAuth2 provider configurations
func (s *AuthConfigServiceImpl) GetOAuth2Providers(ctx context.Context) ([]*response.OAuth2ProviderResponse, error) {
	s.logger.InfoContext(ctx, "Getting OAuth2 providers")

	providers := s.authConfigManager.GetEnabledOAuth2Providers()

	resp := make([]*response.OAuth2ProviderResponse, 0, len(providers))
	for i := range providers {
		provider := &providers[i]
		resp = append(resp, &response.OAuth2ProviderResponse{
			Name:         provider.Name,
			Provider:     getProviderKeyByName(s.authConfigManager.GetConfig().UserAuth.OAuth2.Providers, provider.Name),
			ClientID:     maskSensitiveValue(provider.ClientID),
			ClientSecret: maskSensitiveValue(provider.ClientSecret),
			RedirectURI:  provider.RedirectURI,
			Scopes:       provider.Scopes,
			Enabled:      provider.Enabled,
			AutoRegister: provider.AutoRegister,
			UserMapping:  provider.UserMapping,
		})
	}

	s.logger.InfoContext(ctx, "Successfully retrieved OAuth2 providers", logger.Int("count", len(resp)))
	return resp, nil
}

// ValidatePassword validates a password against the configured password policy
func (s *AuthConfigServiceImpl) ValidatePassword(password string) error {
	policy := s.authConfigManager.GetPasswordPolicy()

	// Check minimum length
	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	// Check maximum length
	if len(password) > policy.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", policy.MaxLength)
	}

	// Check for uppercase letters
	if policy.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for lowercase letters
	if policy.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for numbers
	if policy.RequireNumbers && !containsNumber(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	// Check for symbols
	if policy.RequireSymbols && !containsSymbol(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// CheckLoginSecurity checks login security constraints
func (s *AuthConfigServiceImpl) CheckLoginSecurity(ctx context.Context, userID uint, ip string) error {
	security := s.authConfigManager.GetLoginSecurity()

	// Check IP whitelist if enabled
	if security.IPWhitelistEnabled && len(security.IPWhitelist) > 0 {
		if !s.isIPInWhitelist(ip, security.IPWhitelist) {
			s.logger.WarnContext(ctx, "Login attempt from non-whitelisted IP",
				logger.Uint("user_id", userID),
				logger.String("ip", ip))
			return fmt.Errorf("login from this IP address is not allowed")
		}
	}

	// Check login time restrictions if enabled
	if security.LoginTimeRestriction {
		allowed := s.isLoginTimeAllowed(security.AllowedLoginHours)
		if !allowed {
			s.logger.WarnContext(ctx, "Login attempt outside allowed hours",
				logger.Uint("user_id", userID),
				logger.String("allowed_hours", security.AllowedLoginHours))
			return fmt.Errorf("login is not allowed at this time")
		}
	}

	return nil
}

// RecordLoginAttempt records a login attempt for security monitoring
func (s *AuthConfigServiceImpl) RecordLoginAttempt(ctx context.Context, userID uint, success bool, ip string) error {
	// This is a placeholder implementation
	// In a real implementation, you would store login attempts in a database or cache
	s.logger.InfoContext(ctx, "Login attempt recorded",
		logger.Uint("user_id", userID),
		logger.Bool("success", success),
		logger.String("ip", ip))

	return nil
}

// Helper methods

// convertOAuth2ProvidersToResponse converts internal OAuth2 provider config to response format
func (s *AuthConfigServiceImpl) convertOAuth2ProvidersToResponse(providers map[string]config.OAuth2ProviderConfig) []*response.OAuth2ProviderResponse {
	resp := make([]*response.OAuth2ProviderResponse, 0, len(providers))

	for key := range providers {
		provider := providers[key]
		resp = append(resp, &response.OAuth2ProviderResponse{
			Name:         provider.Name,
			Provider:     key,
			ClientID:     maskSensitiveValue(provider.ClientID),
			ClientSecret: maskSensitiveValue(provider.ClientSecret),
			RedirectURI:  provider.RedirectURI,
			Scopes:       provider.Scopes,
			Enabled:      provider.Enabled,
			AutoRegister: provider.AutoRegister,
			UserMapping:  provider.UserMapping,
		})
	}

	return resp
}

// updatePasswordPolicy updates password policy configuration
func (s *AuthConfigServiceImpl) updatePasswordPolicy(current *config.PasswordPolicyConfig, req *request.PasswordPolicyRequest) {
	if req.MinLength != nil {
		current.MinLength = *req.MinLength
	}
	if req.MaxLength != nil {
		current.MaxLength = *req.MaxLength
	}
	if req.RequireUppercase != nil {
		current.RequireUppercase = *req.RequireUppercase
	}
	if req.RequireLowercase != nil {
		current.RequireLowercase = *req.RequireLowercase
	}
	if req.RequireNumbers != nil {
		current.RequireNumbers = *req.RequireNumbers
	}
	if req.RequireSymbols != nil {
		current.RequireSymbols = *req.RequireSymbols
	}
	if req.PasswordHistory != nil {
		current.PasswordHistory = *req.PasswordHistory
	}
	if req.PasswordExpiresDays != nil {
		current.PasswordExpiresDays = *req.PasswordExpiresDays
	}
}

// updateLoginSecurity updates login security configuration
func (s *AuthConfigServiceImpl) updateLoginSecurity(current *config.LoginSecurityConfig, req *request.LoginSecurityRequest) {
	if req.MaxLoginAttempts != nil {
		current.MaxLoginAttempts = *req.MaxLoginAttempts
	}
	if req.LockoutDuration != nil {
		current.LockoutDuration = *req.LockoutDuration
	}
	if req.IPWhitelistEnabled != nil {
		current.IPWhitelistEnabled = *req.IPWhitelistEnabled
	}
	if req.IPWhitelist != nil {
		current.IPWhitelist = *req.IPWhitelist
	}
	if req.LoginTimeRestriction != nil {
		current.LoginTimeRestriction = *req.LoginTimeRestriction
	}
	if req.AllowedLoginHours != nil {
		current.AllowedLoginHours = *req.AllowedLoginHours
	}
}

// Utility functions

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// maskSensitiveValue masks sensitive configuration values
func maskSensitiveValue(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= constants.PasswordMinLengthCheck {
		return "***"
	}
	return value[:4] + "***" + value[len(value)-4:]
}

// getProviderKeyByName finds provider key by name
func getProviderKeyByName(providers map[string]config.OAuth2ProviderConfig, name string) string {
	for key := range providers {
		provider := providers[key]
		if provider.Name == name {
			return key
		}
	}
	return ""
}

// Password validation helper functions

// containsUppercase checks if string contains uppercase letters
func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// containsLowercase checks if string contains lowercase letters
func containsLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// containsNumber checks if string contains numbers
func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

// containsSymbol checks if string contains special characters
func containsSymbol(s string) bool {
	symbolRegex := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)
	return symbolRegex.MatchString(s)
}

// isIPInWhitelist checks if IP address is in whitelist
func (s *AuthConfigServiceImpl) isIPInWhitelist(ip string, whitelist []string) bool {
	for _, allowedIP := range whitelist {
		if ip == allowedIP {
			return true
		}
		// Support CIDR notation matching in future implementation
		// For now, only exact matches are supported
	}
	return false
}

// isLoginTimeAllowed checks if current time is within allowed login hours
// nolint:unparam // This is a placeholder implementation that always returns true
func (s *AuthConfigServiceImpl) isLoginTimeAllowed(allowedHours string) bool {
	// Parse time range format: "09:00-18:00"
	if allowedHours == "" {
		return true // No restriction
	}

	// This is a simplified implementation
	// In production, you would parse the time range and compare with current time
	parts := strings.Split(allowedHours, "-")
	if len(parts) != constants.TimeRangePartsCount {
		return true // Invalid format, allow by default
	}

	// TODO: Implement actual time range checking
	// For now, always return true to maintain backward compatibility
	return true
}
