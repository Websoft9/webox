package request

import "api-service/internal/dto/common"

// CreateRoleRequest creates role request
type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=64"`
	Code          string `json:"code" validate:"required,min=2,max=32"`
	Description   string `json:"description" validate:"max=500"`
	PermissionIDs []uint `json:"permission_ids"`
	SortOrder     int    `json:"sort_order"`
}

// UpdateRoleRequest updates role request
type UpdateRoleRequest struct {
	Name          string `json:"name" validate:"omitempty,min=2,max=64"`
	Description   string `json:"description" validate:"max=500"`
	PermissionIDs []uint `json:"permission_ids"`
	SortOrder     int    `json:"sort_order"`
	Status        int    `json:"status" validate:"oneof=-1 0 1"`
}

// RolePermissionRequest role permission operation request
type RolePermissionRequest struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}

// ListRolesRequest role list query request
type ListRolesRequest struct {
	common.PaginationRequest
	Search string `form:"search"`
	Status *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
	common.TimeRangeRequest
	common.SortRequest
}

// CreatePermissionRequest creates permission request
type CreatePermissionRequest struct {
	ParentID    *uint  `json:"parent_id"`
	Scope       string `json:"scope" validate:"required,oneof=platform project"`
	Name        string `json:"name" validate:"required,min=2,max=64"`
	Code        string `json:"code" validate:"required,min=2,max=64"`
	Module      string `json:"module" validate:"required,min=2,max=32"`
	Action      string `json:"action" validate:"required,min=2,max=32"`
	Resource    string `json:"resource" validate:"max=64"`
	Description string `json:"description" validate:"max=500"`
	IsMenu      bool   `json:"is_menu"`
	SortOrder   int    `json:"sort_order"`
}

// UpdatePermissionRequest updates permission request
type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=64"`
	Description string `json:"description" validate:"max=500"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status" validate:"oneof=-1 0 1"`
}

// ListPermissionsRequest permission list query request
type ListPermissionsRequest struct {
	common.PaginationRequest
	Search string `form:"search"`
	Module string `form:"module"`
	Scope  string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
	common.TimeRangeRequest
	common.SortRequest
}

// PermissionTreeRequest permission tree query request
type PermissionTreeRequest struct {
	Scope  string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
}

// UpdateAuthConfigRequest updates authentication config request
type UpdateAuthConfigRequest struct {
	APIAuth       *APIAuthRequest       `json:"api_auth,omitempty"`
	UserAuth      *UserAuthRequest      `json:"user_auth,omitempty"`
	SessionConfig *SessionConfigRequest `json:"session_config,omitempty"`
}

// APIAuthRequest API authentication request
type APIAuthRequest struct {
	OAuth2    OAuth2Request     `json:"oauth2,omitempty"`
	TokenAuth *TokenAuthRequest `json:"token_auth,omitempty"`
}

// TokenAuthRequest Token configuration request
type TokenAuthRequest struct {
	Algorithm        *string `json:"algorithm,omitempty"`
	Secret           *string `json:"secret,omitempty"`
	ExpiresIn        *int    `json:"expires_in,omitempty"`
	RefreshExpiresIn *int    `json:"refresh_expires_in,omitempty"`
	AutoRefresh      *bool   `json:"auto_refresh,omitempty"`
}

// OAuth2Request OAuth2 API authentication configuration request
type OAuth2Request struct {
	Enabled           *bool     `json:"enabled,omitempty"`
	DefaultScopes     *[]string `json:"default_scopes,omitempty"`
	TokenEndpoint     *string   `json:"token_endpoint,omitempty"`
	AuthorizeEndpoint *string   `json:"authorize_endpoint,omitempty"`
}

// UserAuthRequest user authentication request
type UserAuthRequest struct {
	// Basic authentication configuration
	BasicAuth *BasicAuthRequest `json:"basic_auth,omitempty"`
	// Email authentication configuration
	EmailAuth *EmailAuthRequest `json:"email_auth,omitempty"`
	// OAuth2 configuration
	OAuth2 *OAuth2LoginConfigRequest `json:"oauth2,omitempty"`
	// Two-factor authentication configuration
	TwoFactor *TwoFactorRequest `json:"two_factor,omitempty"`
	// Password policy configuration
	PasswordPolicy *PasswordPolicyRequest `json:"password_policy,omitempty"`
	// Login security configuration
	LoginSecurity *LoginSecurityRequest `json:"login_security,omitempty"`
}

// SessionConfigRequest session configuration request
type SessionConfigRequest struct {
	Timeout               *int  `json:"timeout,omitempty"`
	MaxConcurrentSessions *int  `json:"max_concurrent_sessions,omitempty"`
	RememberMeEnabled     *bool `json:"remember_me_enabled,omitempty"`
	RememberMeDuration    *int  `json:"remember_me_duration,omitempty"`
}

// BasicAuthRequest basic authentication request
type BasicAuthRequest struct {
	LoginMethods *[]string `json:"login_methods,omitempty"`
}

// EmailAuthRequest email authentication request
type EmailAuthRequest struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn *int  `json:"expires_in,omitempty"`
}

// OAuth2LoginConfigRequest OAuth2 login configuration request
type OAuth2LoginConfigRequest struct {
	Enabled      *bool                    `json:"enabled,omitempty"`
	AutoRegister *bool                    `json:"auto_register,omitempty"`
	DefaultRole  *string                  `json:"default_role,omitempty"`
	Providers    *[]OAuth2ProviderRequest `json:"providers,omitempty"`
}

// TwoFactorRequest two-factor authentication request
type TwoFactorRequest struct {
	Enabled       *bool                    `json:"enabled,omitempty"`
	RequiredRoles *[]string                `json:"required_roles,omitempty"`
	Methods       *TwoFactorMethodsRequest `json:"methods,omitempty"`
}

// TwoFactorMethodsRequest two-factor methods request
type TwoFactorMethodsRequest struct {
	TOTP  *TOTPMethodRequest  `json:"totp,omitempty"`
	Email *EmailMethodRequest `json:"email,omitempty"`
}

// TOTPMethodRequest TOTP method configuration request
type TOTPMethodRequest struct {
	Enabled          *bool   `json:"enabled,omitempty"`
	Issuer           *string `json:"issuer,omitempty"`
	Algorithm        *string `json:"algorithm,omitempty"`
	Digits           *int    `json:"digits,omitempty"`
	Period           *int    `json:"period,omitempty"`
	BackupCodesCount *int    `json:"backup_codes_count,omitempty"`
}

// EmailMethodRequest email method configuration request
type EmailMethodRequest struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	CodeLength *int    `json:"code_length,omitempty"`
	ExpiresIn  *int    `json:"expires_in,omitempty"`
	RateLimit  *int    `json:"rate_limit,omitempty"`
	Template   *string `json:"template,omitempty"`
}

// OAuth2ProviderRequest OAuth2 provider request
type OAuth2ProviderRequest struct {
	Name         *string            `json:"name,omitempty"`
	Enabled      *bool              `json:"enabled,omitempty"`
	ClientID     *string            `json:"client_id,omitempty"`
	ClientSecret *string            `json:"client_secret,omitempty"`
	RedirectURI  *string            `json:"redirect_uri,omitempty"`
	Scopes       *[]string          `json:"scopes,omitempty"`
	AuthorizeURL *string            `json:"authorize_url,omitempty"`
	TokenURL     *string            `json:"token_url,omitempty"`
	UserInfoURL  *string            `json:"user_info_url,omitempty"`
	AutoRegister *bool              `json:"auto_register,omitempty"`
	UserMapping  *map[string]string `json:"user_mapping,omitempty"`
	SortOrder    *int               `json:"sort_order,omitempty"`
}

// PasswordPolicyRequest password policy request
type PasswordPolicyRequest struct {
	MinLength           *int  `json:"min_length,omitempty" validate:"omitempty,min=1,max=128"`
	MaxLength           *int  `json:"max_length,omitempty" validate:"omitempty,min=1,max=128"`
	RequireUppercase    *bool `json:"require_uppercase,omitempty"`
	RequireLowercase    *bool `json:"require_lowercase,omitempty"`
	RequireNumbers      *bool `json:"require_numbers,omitempty"`
	RequireSymbols      *bool `json:"require_symbols,omitempty"`
	PasswordExpiresDays *int  `json:"password_expires_days,omitempty" validate:"omitempty,min=0,max=365"`
}

// LoginSecurityRequest login security request
type LoginSecurityRequest struct {
	MaxLoginAttempts     *int      `json:"max_login_attempts,omitempty" validate:"omitempty,min=1,max=20"`
	LockoutDuration      *int      `json:"lockout_duration,omitempty" validate:"omitempty,min=60,max=86400"`
	IPWhitelistEnabled   *bool     `json:"ip_whitelist_enabled,omitempty"`
	IPWhitelist          *[]string `json:"ip_whitelist,omitempty"`
	LoginTimeRestriction *bool     `json:"login_time_restriction,omitempty"`
	AllowedLoginHours    *string   `json:"allowed_login_hours,omitempty"`
}

// EnableTwoFactorRequest enables two-factor authentication request
type EnableTwoFactorRequest struct {
	Method string `json:"method" validate:"required,oneof=TOTP EMAIL"`
	Code   string `json:"code" validate:"required"`
}

// RevokeAPITokenRequest API token revoke request
type RevokeAPITokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// ConfirmTOTPRequest TOTP confirmation request
type ConfirmTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// DisableTOTPRequest disables TOTP request
type DisableTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// EnableEmailTwoFactorRequest enables email two-factor authentication request
type EnableEmailTwoFactorRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyTwoFactorRequest verifies two-factor authentication request
type VerifyTwoFactorRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Method string `json:"method" validate:"required,oneof=totp email backup"`
}
