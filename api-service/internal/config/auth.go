package config

import (
	"api-service/internal/constants"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// AuthConfig represents the complete authentication configuration structure
type AuthConfig struct {
	Version       string         `yaml:"version" json:"version" mapstructure:"version"`
	APIAuth       APIAuthConfig  `yaml:"api_auth" json:"api_auth" mapstructure:"api_auth"`
	UserAuth      UserAuthConfig `yaml:"user_auth" json:"user_auth" mapstructure:"user_auth"`
	SessionConfig SessionConfig  `yaml:"session_config" json:"session_config" mapstructure:"session_config"`
}

// APIAuthConfig represents API authentication configuration
type APIAuthConfig struct {
	TokenAuth TokenAuthConfig `yaml:"token_auth" json:"token_auth" mapstructure:"token_auth"`
	OAuth2    OAuth2Config    `yaml:"oauth2" json:"oauth2" mapstructure:"oauth2"`
}

// TokenAuthConfig represents JWT token authentication configuration
type TokenAuthConfig struct {
	Algorithm        string `yaml:"algorithm" json:"algorithm" mapstructure:"algorithm"`
	Secret           string `yaml:"secret" json:"secret" mapstructure:"secret"`
	ExpiresIn        int    `yaml:"expires_in" json:"expires_in" mapstructure:"expires_in"`
	RefreshExpiresIn int    `yaml:"refresh_expires_in" json:"refresh_expires_in" mapstructure:"refresh_expires_in"`
	AutoRefresh      bool   `yaml:"auto_refresh" json:"auto_refresh" mapstructure:"auto_refresh"`
}

// OAuth2Config represents OAuth2 API authentication configuration
type OAuth2Config struct {
	Enabled           bool     `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	DefaultScopes     []string `yaml:"default_scopes" json:"default_scopes" mapstructure:"default_scopes"`
	TokenEndpoint     string   `yaml:"token_endpoint" json:"token_endpoint" mapstructure:"token_endpoint"`
	AuthorizeEndpoint string   `yaml:"authorize_endpoint" json:"authorize_endpoint" mapstructure:"authorize_endpoint"`
}

// UserAuthConfig represents user authentication configuration
type UserAuthConfig struct {
	BasicAuth      BasicAuthConfig      `yaml:"basic_auth" json:"basic_auth" mapstructure:"basic_auth"`
	EmailAuth      EmailAuthConfig      `yaml:"email_auth" json:"email_auth" mapstructure:"email_auth"`
	PasswordPolicy PasswordPolicyConfig `yaml:"password_policy" json:"password_policy" mapstructure:"password_policy"`
	LoginSecurity  LoginSecurityConfig  `yaml:"login_security" json:"login_security" mapstructure:"login_security"`
	OAuth2         UserOAuth2Config     `yaml:"oauth2" json:"oauth2" mapstructure:"oauth2"`
	TwoFactor      TwoFactorConfig      `yaml:"two_factor" json:"two_factor" mapstructure:"two_factor"`
}

// BasicAuthConfig represents basic authentication configuration
type BasicAuthConfig struct {
	LoginMethods []string `yaml:"login_methods" json:"login_methods" mapstructure:"login_methods"`
}

type EmailAuthConfig struct {
	Enabled   bool `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	ExpiresIn int  `yaml:"expires_in" json:"expires_in" mapstructure:"expires_in"`
}

// PasswordPolicyConfig represents password policy configuration
type PasswordPolicyConfig struct {
	MinLength           int  `yaml:"min_length" json:"min_length" mapstructure:"min_length"`
	MaxLength           int  `yaml:"max_length" json:"max_length" mapstructure:"max_length"`
	RequireUppercase    bool `yaml:"require_uppercase" json:"require_uppercase" mapstructure:"require_uppercase"`
	RequireLowercase    bool `yaml:"require_lowercase" json:"require_lowercase" mapstructure:"require_lowercase"`
	RequireNumbers      bool `yaml:"require_numbers" json:"require_numbers" mapstructure:"require_numbers"`
	RequireSymbols      bool `yaml:"require_symbols" json:"require_symbols" mapstructure:"require_symbols"`
	PasswordExpiresDays int  `yaml:"password_expires_days" json:"password_expires_days" mapstructure:"password_expires_days"`
}

// LoginSecurityConfig represents login security configuration
type LoginSecurityConfig struct {
	MaxLoginAttempts     int      `yaml:"max_login_attempts" json:"max_login_attempts" mapstructure:"max_login_attempts"`
	LockoutDuration      int      `yaml:"lockout_duration" json:"lockout_duration" mapstructure:"lockout_duration"`
	IPWhitelistEnabled   bool     `yaml:"ip_whitelist_enabled" json:"ip_whitelist_enabled" mapstructure:"ip_whitelist_enabled"`
	IPWhitelist          []string `yaml:"ip_whitelist" json:"ip_whitelist" mapstructure:"ip_whitelist"`
	LoginTimeRestriction bool     `yaml:"login_time_restriction" json:"login_time_restriction" mapstructure:"login_time_restriction"`
	AllowedLoginHours    string   `yaml:"allowed_login_hours" json:"allowed_login_hours" mapstructure:"allowed_login_hours"`
}

// UserOAuth2Config represents user OAuth2 login configuration
type UserOAuth2Config struct {
	Enabled      bool                            `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	AutoRegister bool                            `yaml:"auto_register" json:"auto_register" mapstructure:"auto_register"`
	DefaultRole  string                          `yaml:"default_role" json:"default_role" mapstructure:"default_role"`
	Providers    map[string]OAuth2ProviderConfig `yaml:"providers" json:"providers" mapstructure:"providers"`
}

// OAuth2ProviderConfig represents OAuth2 provider configuration
type OAuth2ProviderConfig struct {
	Name         string            `yaml:"name" json:"name" mapstructure:"name"`
	Enabled      bool              `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	ClientID     string            `yaml:"client_id" json:"client_id" mapstructure:"client_id"`
	ClientSecret string            `yaml:"client_secret" json:"client_secret" mapstructure:"client_secret"`
	RedirectURI  string            `yaml:"redirect_uri" json:"redirect_uri" mapstructure:"redirect_uri"`
	Scopes       []string          `yaml:"scopes" json:"scopes" mapstructure:"scopes"`
	AuthorizeURL string            `yaml:"authorize_url" json:"authorize_url" mapstructure:"authorize_url"`
	TokenURL     string            `yaml:"token_url" json:"token_url" mapstructure:"token_url"`
	UserInfoURL  string            `yaml:"user_info_url" json:"user_info_url" mapstructure:"user_info_url"`
	AutoRegister bool              `yaml:"auto_register" json:"auto_register" mapstructure:"auto_register"`
	UserMapping  map[string]string `yaml:"user_mapping" json:"user_mapping" mapstructure:"user_mapping"`
	SortOrder    int               `yaml:"sort_order" json:"sort_order" mapstructure:"sort_order"`
}

// TwoFactorConfig represents two-factor authentication configuration
type TwoFactorConfig struct {
	Enabled       bool                   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	RequiredRoles []string               `yaml:"required_roles" json:"required_roles" mapstructure:"required_roles"`
	Methods       TwoFactorMethodsConfig `yaml:"methods" json:"methods" mapstructure:"methods"`
}

// TwoFactorMethodsConfig represents two-factor authentication methods configuration
type TwoFactorMethodsConfig struct {
	TOTP  TOTPConfig  `yaml:"totp" json:"totp" mapstructure:"totp"`
	Email EmailConfig `yaml:"email" json:"email" mapstructure:"email"`
}

// TOTPConfig represents TOTP authentication configuration
type TOTPConfig struct {
	Enabled          bool   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	Issuer           string `yaml:"issuer" json:"issuer" mapstructure:"issuer"`
	Algorithm        string `yaml:"algorithm" json:"algorithm" mapstructure:"algorithm"`
	Digits           int    `yaml:"digits" json:"digits" mapstructure:"digits"`
	Period           int    `yaml:"period" json:"period" mapstructure:"period"`
	BackupCodesCount int    `yaml:"backup_codes_count" json:"backup_codes_count" mapstructure:"backup_codes_count"`
}

// EmailConfig represents email verification code configuration
type EmailConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
	CodeLength int    `yaml:"code_length" json:"code_length" mapstructure:"code_length"`
	ExpiresIn  int    `yaml:"expires_in" json:"expires_in" mapstructure:"expires_in"`
	RateLimit  int    `yaml:"rate_limit" json:"rate_limit" mapstructure:"rate_limit"`
	Template   string `yaml:"template" json:"template" mapstructure:"template"`
}

// SessionConfig represents session configuration
type SessionConfig struct {
	Timeout               int  `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
	MaxConcurrentSessions int  `yaml:"max_concurrent_sessions" json:"max_concurrent_sessions" mapstructure:"max_concurrent_sessions"`
	RememberMeEnabled     bool `yaml:"remember_me_enabled" json:"remember_me_enabled" mapstructure:"remember_me_enabled"`
	RememberMeDuration    int  `yaml:"remember_me_duration" json:"remember_me_duration" mapstructure:"remember_me_duration"`
}

// AuthConfigManager manages authentication configuration loading and access
type AuthConfigManager struct {
	config   *AuthConfig
	viper    *viper.Viper
	filePath string
}

// NewAuthConfigManager creates a new authentication configuration manager
func NewAuthConfigManager(configPath string) (*AuthConfigManager, error) {
	manager := &AuthConfigManager{
		viper:    viper.New(),
		filePath: configPath,
	}

	if err := manager.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load auth configuration: %w", err)
	}

	return manager, nil
}

// loadConfig loads authentication configuration from YAML file
func (m *AuthConfigManager) loadConfig() error {
	// Set default values
	m.setDefaults()

	// Configure viper
	m.viper.SetConfigFile(m.filePath)
	m.viper.SetConfigType("yaml")

	// Enable environment variable substitution
	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("WEBSOFT9")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration file
	if err := m.viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			// Create default config file if it doesn't exist
			if createErr := m.createDefaultConfig(); createErr != nil {
				return fmt.Errorf("failed to create default config: %w", createErr)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal configuration
	var config AuthConfig
	if err := m.viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Process environment variable substitutions
	m.processEnvironmentVariables(&config)

	m.config = &config
	return nil
}

// setDefaults sets default configuration values
func (m *AuthConfigManager) setDefaults() {
	// API authentication defaults
	m.viper.SetDefault("api_auth.token_auth.algorithm", "HS256")
	m.viper.SetDefault("api_auth.token_auth.secret", "Websoft9")
	m.viper.SetDefault("api_auth.token_auth.expires_in", constants.DefaultTokenExpiresIn)
	m.viper.SetDefault("api_auth.token_auth.refresh_expires_in", constants.DefaultRefreshTokenExpiresIn)
	m.viper.SetDefault("api_auth.token_auth.auto_refresh", true)

	// User authentication defaults
	m.viper.SetDefault("user_auth.basic_auth.login_methods", []string{"username", "email"})
	m.viper.SetDefault("user_auth.email_auth.enabled", true)
	m.viper.SetDefault("user_auth.email_auth.expires_in", constants.DefaultTokenExpiresIn)

	// Password policy defaults
	m.viper.SetDefault("user_auth.password_policy.min_length", constants.PasswordMinLength)
	m.viper.SetDefault("user_auth.password_policy.max_length", constants.PasswordMaxLength)
	m.viper.SetDefault("user_auth.password_policy.require_uppercase", true)
	m.viper.SetDefault("user_auth.password_policy.require_lowercase", true)
	m.viper.SetDefault("user_auth.password_policy.require_numbers", true)
	m.viper.SetDefault("user_auth.password_policy.require_symbols", false)
	m.viper.SetDefault("user_auth.password_policy.password_expires_days", constants.PasswordExpiresDays)

	// Login security defaults
	m.viper.SetDefault("user_auth.login_security.max_login_attempts", constants.MaxLoginAttempts)
	m.viper.SetDefault("user_auth.login_security.lockout_duration", constants.LockoutDuration)
	m.viper.SetDefault("user_auth.login_security.ip_whitelist_enabled", false)
	m.viper.SetDefault("user_auth.login_security.login_time_restriction", false)

	// Two-factor authentication defaults
	m.viper.SetDefault("user_auth.two_factor.enabled", true)
	m.viper.SetDefault("user_auth.two_factor.required_roles", []string{"admin"})
	m.viper.SetDefault("user_auth.two_factor.methods.totp.enabled", true)
	m.viper.SetDefault("user_auth.two_factor.methods.totp.issuer", "Websoft9")
	m.viper.SetDefault("user_auth.two_factor.methods.totp.algorithm", "SHA1")
	m.viper.SetDefault("user_auth.two_factor.methods.totp.digits", constants.TOTPDigits)
	m.viper.SetDefault("user_auth.two_factor.methods.totp.period", constants.TOTPPeriod)
	m.viper.SetDefault("user_auth.two_factor.methods.totp.backup_codes_count", constants.BackupCodesCount)

	// Session configuration defaults
	m.viper.SetDefault("session_config.timeout", constants.SessionTimeout)
	m.viper.SetDefault("session_config.max_concurrent_sessions", constants.MaxConcurrentSessions)
	m.viper.SetDefault("session_config.remember_me_enabled", true)
	m.viper.SetDefault("session_config.remember_me_duration", constants.RememberMeDuration)
}

// createDefaultConfig creates a default configuration file
func (m *AuthConfigManager) createDefaultConfig() error {
	// Ensure directory exists
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, constants.DirPerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default configuration with current defaults
	if err := m.viper.WriteConfigAs(m.filePath); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}

// processEnvironmentVariables processes environment variable substitutions in configuration
func (m *AuthConfigManager) processEnvironmentVariables(config *AuthConfig) {
	// Process OAuth2 provider configurations
	for providerKey := range config.UserAuth.OAuth2.Providers {
		provider := config.UserAuth.OAuth2.Providers[providerKey]
		provider.ClientID = m.expandEnvironmentVariables(provider.ClientID)
		provider.ClientSecret = m.expandEnvironmentVariables(provider.ClientSecret)
		provider.RedirectURI = m.expandEnvironmentVariables(provider.RedirectURI)
		provider.AuthorizeURL = m.expandEnvironmentVariables(provider.AuthorizeURL)
		provider.TokenURL = m.expandEnvironmentVariables(provider.TokenURL)
		provider.UserInfoURL = m.expandEnvironmentVariables(provider.UserInfoURL)

		config.UserAuth.OAuth2.Providers[providerKey] = provider
	}
}

// expandEnvironmentVariables expands environment variables in the given string
func (m *AuthConfigManager) expandEnvironmentVariables(s string) string {
	return os.ExpandEnv(s)
}

// GetConfig returns the current authentication configuration
func (m *AuthConfigManager) GetConfig() *AuthConfig {
	return m.config
}

// UpdateConfig updates the authentication configuration and saves to file
func (m *AuthConfigManager) UpdateConfig(config *AuthConfig) error {
	// Validate configuration
	if err := m.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Update viper configuration
	if err := m.updateViperConfig(config); err != nil {
		return fmt.Errorf("failed to update viper config: %w", err)
	}

	// Write configuration to file
	if err := m.viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	m.config = config
	return nil
}

// validateConfig validates the authentication configuration
func (m *AuthConfigManager) validateConfig(config *AuthConfig) error {
	// Validate password policy
	if config.UserAuth.PasswordPolicy.MinLength < 4 || config.UserAuth.PasswordPolicy.MinLength > 128 {
		return fmt.Errorf("password min_length must be between 4 and 128")
	}

	if config.UserAuth.PasswordPolicy.MaxLength < config.UserAuth.PasswordPolicy.MinLength {
		return fmt.Errorf("password max_length must be greater than min_length")
	}

	// Validate token expiration times
	if config.APIAuth.TokenAuth.ExpiresIn <= 0 {
		return fmt.Errorf("token expires_in must be positive")
	}

	if config.APIAuth.TokenAuth.RefreshExpiresIn <= config.APIAuth.TokenAuth.ExpiresIn {
		return fmt.Errorf("refresh_expires_in must be greater than expires_in")
	}

	// Validate OAuth2 providers
	for providerName := range config.UserAuth.OAuth2.Providers {
		provider := config.UserAuth.OAuth2.Providers[providerName]
		if provider.Enabled {
			if provider.ClientID == "" || provider.ClientSecret == "" {
				return fmt.Errorf("OAuth2 provider %s: client_id and client_secret are required", providerName)
			}
			if provider.RedirectURI == "" {
				return fmt.Errorf("OAuth2 provider %s: redirect_uri is required", providerName)
			}
		}
	}

	return nil
}

// updateViperConfig updates viper configuration from struct
func (m *AuthConfigManager) updateViperConfig(config *AuthConfig) error {
	// Convert struct to map for viper
	configMap := make(map[string]interface{})

	// This is a simplified implementation - in production you might want to use
	// reflection or a more sophisticated method to convert struct to map
	configMap["version"] = config.Version
	configMap["api_auth"] = config.APIAuth
	configMap["user_auth"] = config.UserAuth
	configMap["session_config"] = config.SessionConfig

	// Merge with existing viper configuration
	if err := m.viper.MergeConfigMap(configMap); err != nil {
		return err
	}

	return nil
}

// ReloadConfig reloads configuration from file
func (m *AuthConfigManager) ReloadConfig() error {
	return m.loadConfig()
}

// GetOAuth2Provider returns a specific OAuth2 provider configuration
func (m *AuthConfigManager) GetOAuth2Provider(providerName string) (*OAuth2ProviderConfig, bool) {
	provider, exists := m.config.UserAuth.OAuth2.Providers[providerName]
	if !exists {
		return nil, false
	}
	return &provider, true
}

// GetEnabledOAuth2Providers returns all enabled OAuth2 providers
func (m *AuthConfigManager) GetEnabledOAuth2Providers() []OAuth2ProviderConfig {
	var providers []OAuth2ProviderConfig
	for key := range m.config.UserAuth.OAuth2.Providers {
		provider := m.config.UserAuth.OAuth2.Providers[key]
		if provider.Enabled {
			providers = append(providers, provider)
		}
	}
	return providers
}

// IsOAuth2Enabled returns whether OAuth2 authentication is enabled
func (m *AuthConfigManager) IsOAuth2Enabled() bool {
	return m.config.UserAuth.OAuth2.Enabled
}

// IsTwoFactorEnabled returns whether two-factor authentication is enabled
func (m *AuthConfigManager) IsTwoFactorEnabled() bool {
	return m.config.UserAuth.TwoFactor.Enabled
}

// IsTwoFactorRequiredForRole checks if two-factor authentication is required for a role
func (m *AuthConfigManager) IsTwoFactorRequiredForRole(role string) bool {
	for _, requiredRole := range m.config.UserAuth.TwoFactor.RequiredRoles {
		if requiredRole == role {
			return true
		}
	}
	return false
}

// GetPasswordPolicy returns the password policy configuration
func (m *AuthConfigManager) GetPasswordPolicy() *PasswordPolicyConfig {
	return &m.config.UserAuth.PasswordPolicy
}

// GetLoginSecurity returns the login security configuration
func (m *AuthConfigManager) GetLoginSecurity() *LoginSecurityConfig {
	return &m.config.UserAuth.LoginSecurity
}
