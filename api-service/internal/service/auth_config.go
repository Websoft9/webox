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
			OAuth2: response.OAuth2Response{
				Enabled:           authConfig.APIAuth.OAuth2.Enabled,
				DefaultScopes:     authConfig.APIAuth.OAuth2.DefaultScopes,
				TokenEndpoint:     authConfig.APIAuth.OAuth2.TokenEndpoint,
				AuthorizeEndpoint: authConfig.APIAuth.OAuth2.AuthorizeEndpoint,
			},
			TokenAuth: response.TokenAuthResponse{
				Algorithm:        authConfig.APIAuth.TokenAuth.Algorithm,
				Secret:           maskSensitiveValue(authConfig.APIAuth.TokenAuth.Secret),
				ExpiresIn:        authConfig.APIAuth.TokenAuth.ExpiresIn,
				RefreshExpiresIn: authConfig.APIAuth.TokenAuth.RefreshExpiresIn,
				AutoRefresh:      authConfig.APIAuth.TokenAuth.AutoRefresh,
			},
		},
		UserAuth: response.UserAuthResponse{
			BasicAuth: response.BasicAuthResponse{
				LoginMethods: authConfig.UserAuth.BasicAuth.LoginMethods,
			},
			EmailAuth: response.EmailAuthResponse{
				Enabled:   authConfig.UserAuth.EmailAuth.Enabled,
				ExpiresIn: authConfig.UserAuth.EmailAuth.ExpiresIn,
			},
			PasswordPolicy: response.PasswordPolicyResponse{
				MinLength:           authConfig.UserAuth.PasswordPolicy.MinLength,
				MaxLength:           authConfig.UserAuth.PasswordPolicy.MaxLength,
				RequireUppercase:    authConfig.UserAuth.PasswordPolicy.RequireUppercase,
				RequireLowercase:    authConfig.UserAuth.PasswordPolicy.RequireLowercase,
				RequireNumbers:      authConfig.UserAuth.PasswordPolicy.RequireNumbers,
				RequireSymbols:      authConfig.UserAuth.PasswordPolicy.RequireSymbols,
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
			OAuth2: response.OAuth2LoginResponse{
				Enabled:      authConfig.UserAuth.OAuth2.Enabled,
				AutoRegister: authConfig.UserAuth.OAuth2.AutoRegister,
				DefaultRole:  authConfig.UserAuth.OAuth2.DefaultRole,
				Providers:    s.convertOAuth2ProvidersToResponse(authConfig.UserAuth.OAuth2.Providers),
			},
			TwoFactor: response.TwoFactorAuthResponse{
				Enabled:       authConfig.UserAuth.TwoFactor.Enabled,
				RequiredRoles: authConfig.UserAuth.TwoFactor.RequiredRoles,
				Methods: response.TwoFactorMethodsResponse{
					TOTP: response.TOTPMethodResponse{
						Enabled:          authConfig.UserAuth.TwoFactor.Methods.TOTP.Enabled,
						Issuer:           authConfig.UserAuth.TwoFactor.Methods.TOTP.Issuer,
						Algorithm:        authConfig.UserAuth.TwoFactor.Methods.TOTP.Algorithm,
						Digits:           authConfig.UserAuth.TwoFactor.Methods.TOTP.Digits,
						Period:           authConfig.UserAuth.TwoFactor.Methods.TOTP.Period,
						BackupCodesCount: authConfig.UserAuth.TwoFactor.Methods.TOTP.BackupCodesCount,
					},
					Email: response.EmailMethodResponse{
						Enabled:    authConfig.UserAuth.TwoFactor.Methods.Email.Enabled,
						CodeLength: authConfig.UserAuth.TwoFactor.Methods.Email.CodeLength,
						ExpiresIn:  authConfig.UserAuth.TwoFactor.Methods.Email.ExpiresIn,
						RateLimit:  authConfig.UserAuth.TwoFactor.Methods.Email.RateLimit,
						Template:   authConfig.UserAuth.TwoFactor.Methods.Email.Template,
					},
				},
			},
		},
		SessionConfig: response.SessionConfigResponse{
			Timeout:               authConfig.SessionConfig.Timeout,
			MaxConcurrentSessions: authConfig.SessionConfig.MaxConcurrentSessions,
			RememberMeEnabled:     authConfig.SessionConfig.RememberMeEnabled,
			RememberMeDuration:    authConfig.SessionConfig.RememberMeDuration,
		},
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

	// Update OAuth2 configuration
	if apiAuth.OAuth2.Enabled != nil {
		config.APIAuth.OAuth2.Enabled = *apiAuth.OAuth2.Enabled
	}
	if apiAuth.OAuth2.DefaultScopes != nil {
		config.APIAuth.OAuth2.DefaultScopes = *apiAuth.OAuth2.DefaultScopes
	}
	if apiAuth.OAuth2.TokenEndpoint != nil {
		config.APIAuth.OAuth2.TokenEndpoint = *apiAuth.OAuth2.TokenEndpoint
	}
	if apiAuth.OAuth2.AuthorizeEndpoint != nil {
		config.APIAuth.OAuth2.AuthorizeEndpoint = *apiAuth.OAuth2.AuthorizeEndpoint
	}

	s.updateTokenAuthConfig(&config.APIAuth.TokenAuth, apiAuth.TokenAuth)
}

// updateTokenAuthConfig updates token authentication configuration
func (s *AuthConfigServiceImpl) updateTokenAuthConfig(tokenAuth *config.TokenAuthConfig, tokenConfig *request.TokenAuthRequest) {
	if tokenConfig == nil {
		return
	}

	if tokenConfig.Algorithm != nil {
		tokenAuth.Algorithm = *tokenConfig.Algorithm
	}
	if tokenConfig.Secret != nil {
		tokenAuth.Secret = *tokenConfig.Secret
	}
	if tokenConfig.ExpiresIn != nil {
		tokenAuth.ExpiresIn = *tokenConfig.ExpiresIn
	}
	if tokenConfig.RefreshExpiresIn != nil {
		tokenAuth.RefreshExpiresIn = *tokenConfig.RefreshExpiresIn
	}
	if tokenConfig.AutoRefresh != nil {
		tokenAuth.AutoRefresh = *tokenConfig.AutoRefresh
	}
}

// updateUserAuthConfig updates user authentication configuration section
func (s *AuthConfigServiceImpl) updateUserAuthConfig(config *config.AuthConfig, userAuth *request.UserAuthRequest) {
	if userAuth == nil {
		return
	}

	// Update basic auth configuration
	if userAuth.BasicAuth != nil {
		s.updateBasicAuthConfig(&config.UserAuth.BasicAuth, userAuth.BasicAuth)
	}

	// Update email auth configuration
	if userAuth.EmailAuth != nil {
		s.updateEmailAuthConfig(&config.UserAuth.EmailAuth, userAuth.EmailAuth)
	}

	// Update OAuth2 configuration
	if userAuth.OAuth2 != nil {
		s.updateOAuth2LoginConfig(&config.UserAuth.OAuth2, userAuth.OAuth2)
	}

	// Update two-factor configuration
	if userAuth.TwoFactor != nil {
		s.updateTwoFactorConfig(&config.UserAuth.TwoFactor, userAuth.TwoFactor)
	}
	if userAuth.PasswordPolicy != nil {
		s.updatePasswordPolicy(&config.UserAuth.PasswordPolicy, userAuth.PasswordPolicy)
	}
	if userAuth.LoginSecurity != nil {
		s.updateLoginSecurity(&config.UserAuth.LoginSecurity, userAuth.LoginSecurity)
	}
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

// updateBasicAuthConfig updates basic authentication configuration
func (s *AuthConfigServiceImpl) updateBasicAuthConfig(current *config.BasicAuthConfig, req *request.BasicAuthRequest) {
	if req.LoginMethods != nil {
		current.LoginMethods = *req.LoginMethods
	}
}

// updateEmailAuthConfig updates email authentication configuration
func (s *AuthConfigServiceImpl) updateEmailAuthConfig(current *config.EmailAuthConfig, req *request.EmailAuthRequest) {
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.ExpiresIn != nil {
		current.ExpiresIn = *req.ExpiresIn
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
			Enabled:      provider.Enabled,
			ClientID:     maskSensitiveValue(provider.ClientID),
			ClientSecret: maskSensitiveValue(provider.ClientSecret),
			RedirectURI:  provider.RedirectURI,
			Scopes:       provider.Scopes,
			AuthorizeURL: provider.AuthorizeURL,
			TokenURL:     provider.TokenURL,
			UserInfoURL:  provider.UserInfoURL,
			AutoRegister: provider.AutoRegister,
			UserMapping:  provider.UserMapping,
			SortOrder:    provider.SortOrder,
		})
	}

	s.logger.InfoContext(ctx, "Successfully retrieved OAuth2 providers", logger.Int("count", len(resp)))
	return resp, nil
}

// Helper methods

// convertOAuth2ProvidersToResponse converts internal OAuth2 provider config to response format
func (s *AuthConfigServiceImpl) convertOAuth2ProvidersToResponse(providers map[string]config.OAuth2ProviderConfig) []*response.OAuth2ProviderResponse {
	resp := make([]*response.OAuth2ProviderResponse, 0, len(providers))

	for key := range providers {
		provider := providers[key]
		resp = append(resp, &response.OAuth2ProviderResponse{
			Name:         provider.Name,
			Enabled:      provider.Enabled,
			ClientID:     maskSensitiveValue(provider.ClientID),
			ClientSecret: maskSensitiveValue(provider.ClientSecret),
			RedirectURI:  provider.RedirectURI,
			Scopes:       provider.Scopes,
			AuthorizeURL: provider.AuthorizeURL,
			TokenURL:     provider.TokenURL,
			UserInfoURL:  provider.UserInfoURL,
			AutoRegister: provider.AutoRegister,
			UserMapping:  provider.UserMapping,
			SortOrder:    provider.SortOrder,
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

// updateOAuth2LoginConfig updates OAuth2 login configuration
func (s *AuthConfigServiceImpl) updateOAuth2LoginConfig(current *config.UserOAuth2Config, req *request.OAuth2LoginConfigRequest) {
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.AutoRegister != nil {
		current.AutoRegister = *req.AutoRegister
	}
	if req.DefaultRole != nil {
		current.DefaultRole = *req.DefaultRole
	}
	// TODO: Providers update would require more complex logic to handle provider-specific updates
}

// updateTwoFactorConfig updates two-factor authentication configuration
func (s *AuthConfigServiceImpl) updateTwoFactorConfig(current *config.TwoFactorConfig, req *request.TwoFactorRequest) {
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.RequiredRoles != nil {
		current.RequiredRoles = *req.RequiredRoles
	}
	if req.Methods != nil {
		s.updateTwoFactorMethodsConfig(&current.Methods, req.Methods)
	}
}

// updateTwoFactorMethodsConfig updates two-factor methods configuration
func (s *AuthConfigServiceImpl) updateTwoFactorMethodsConfig(current *config.TwoFactorMethodsConfig, req *request.TwoFactorMethodsRequest) {
	if req.TOTP != nil {
		s.updateTOTPMethodConfig(&current.TOTP, req.TOTP)
	}
	if req.Email != nil {
		s.updateEmailMethodConfig(&current.Email, req.Email)
	}
}

// updateTOTPMethodConfig updates TOTP method configuration
func (s *AuthConfigServiceImpl) updateTOTPMethodConfig(current *config.TOTPConfig, req *request.TOTPMethodRequest) {
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.Issuer != nil {
		current.Issuer = *req.Issuer
	}
	if req.Algorithm != nil {
		current.Algorithm = *req.Algorithm
	}
	if req.Digits != nil {
		current.Digits = *req.Digits
	}
	if req.Period != nil {
		current.Period = *req.Period
	}
	if req.BackupCodesCount != nil {
		current.BackupCodesCount = *req.BackupCodesCount
	}
}

// updateEmailMethodConfig updates email method configuration
func (s *AuthConfigServiceImpl) updateEmailMethodConfig(current *config.EmailConfig, req *request.EmailMethodRequest) {
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.CodeLength != nil {
		current.CodeLength = *req.CodeLength
	}
	if req.ExpiresIn != nil {
		current.ExpiresIn = *req.ExpiresIn
	}
	if req.RateLimit != nil {
		current.RateLimit = *req.RateLimit
	}
	if req.Template != nil {
		current.Template = *req.Template
	}
}
