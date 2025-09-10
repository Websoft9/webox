package request

import "time"

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
	PaginationRequest
	Search    string `form:"search"`
	Status    *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
	StartTime string `form:"start_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime   string `form:"end_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
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
	PaginationRequest
	Search    string `form:"search"`
	Module    string `form:"module"`
	Scope     string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status    *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
	StartTime string `form:"start_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime   string `form:"end_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
}

// PermissionTreeRequest permission tree query request
type PermissionTreeRequest struct {
	Scope  string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status *int   `form:"status" validate:"omitempty,oneof=-1 0 1"`
}

// CreateAPITokenRequest creates API token request
type CreateAPITokenRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=64"`
	Description string     `json:"description" validate:"max=500"`
	Scopes      []string   `json:"scopes" validate:"required,min=1"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// UpdateAPITokenRequest updates API token request
type UpdateAPITokenRequest struct {
	Name        string     `json:"name" validate:"omitempty,min=2,max=64"`
	Description string     `json:"description" validate:"max=500"`
	Scopes      []string   `json:"scopes"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// ListAPITokensRequest API token list query request
type ListAPITokensRequest struct {
	PaginationRequest
	Search  string `form:"search"`
	UserID  *uint  `form:"user_id"`
	Expired *bool  `form:"expired"`
}

// BatchRevokeAPITokensRequest batch revoke API tokens request
type BatchRevokeAPITokensRequest struct {
	IDs []uint `json:"ids" validate:"required,min=1"`
}

// UpdateAuthConfigRequest updates authentication config request
type UpdateAuthConfigRequest struct {
	APIAuth       *APIAuthRequest       `json:"api_auth,omitempty"`
	UserAuth      *UserAuthRequest      `json:"user_auth,omitempty"`
	SessionConfig *SessionConfigRequest `json:"session_config,omitempty"`
}

// APIAuthRequest API authentication request
type APIAuthRequest struct {
	TokenAuthEnabled *bool             `json:"token_auth_enabled,omitempty"`
	OAuth2Enabled    *bool             `json:"oauth2_enabled,omitempty"`
	JWTConfig        *JWTConfigRequest `json:"jwt_config,omitempty"`
}

// JWTConfigRequest JWT configuration request
type JWTConfigRequest struct {
	Algorithm        *string `json:"algorithm,omitempty"`
	ExpiresIn        *int    `json:"expires_in,omitempty"`
	RefreshExpiresIn *int    `json:"refresh_expires_in,omitempty"`
	AutoRefresh      *bool   `json:"auto_refresh,omitempty"`
}

// UserAuthRequest user authentication request
type UserAuthRequest struct {
	OAuth2Enabled          *bool                    `json:"oauth2_enabled,omitempty"`
	OAuth2Providers        *[]OAuth2ProviderRequest `json:"oauth2_providers,omitempty"`
	TwoFactorEnabled       *bool                    `json:"two_factor_enabled,omitempty"`
	TwoFactorMethods       *[]string                `json:"two_factor_methods,omitempty"`
	TwoFactorRequiredRoles *[]string                `json:"two_factor_required_roles,omitempty"`
	PasswordPolicy         *PasswordPolicyRequest   `json:"password_policy,omitempty"`
	LoginSecurity          *LoginSecurityRequest    `json:"login_security,omitempty"`
}

// SessionConfigRequest session configuration request
type SessionConfigRequest struct {
	Timeout               *int  `json:"timeout,omitempty"`
	MaxConcurrentSessions *int  `json:"max_concurrent_sessions,omitempty"`
	RememberMeEnabled     *bool `json:"remember_me_enabled,omitempty"`
	RememberMeDuration    *int  `json:"remember_me_duration,omitempty"`
}

// OAuth2ProviderRequest OAuth2 provider request
type OAuth2ProviderRequest struct {
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	Enabled      bool              `json:"enabled"`
	AutoRegister bool              `json:"auto_register"`
	UserMapping  map[string]string `json:"user_mapping"`
}

// PasswordPolicyRequest password policy request
type PasswordPolicyRequest struct {
	MinLength           *int  `json:"min_length,omitempty" validate:"omitempty,min=1,max=128"`
	MaxLength           *int  `json:"max_length,omitempty" validate:"omitempty,min=1,max=128"`
	RequireUppercase    *bool `json:"require_uppercase,omitempty"`
	RequireLowercase    *bool `json:"require_lowercase,omitempty"`
	RequireNumbers      *bool `json:"require_numbers,omitempty"`
	RequireSymbols      *bool `json:"require_symbols,omitempty"`
	PasswordHistory     *int  `json:"password_history,omitempty" validate:"omitempty,min=0,max=20"`
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

// PaginationRequest pagination request base structure
type PaginationRequest struct {
	Page     int `form:"page" validate:"min=1"`
	PageSize int `form:"page_size" validate:"min=1,max=100"`
}

// GetPage gets page number, defaults to 1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// GetPageSize gets page size, defaults to 20
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		return MaxPageSize
	}
	return p.PageSize
}

// GetOffset gets offset
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// ValidateAPITokenRequest API token validation request
type ValidateAPITokenRequest struct {
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
