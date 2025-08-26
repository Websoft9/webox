package request

import "time"

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=64"`
	Code          string `json:"code" validate:"required,min=2,max=32"`
	Description   string `json:"description" validate:"max=500"`
	PermissionIDs []uint `json:"permission_ids"`
	SortOrder     int    `json:"sort_order"`
	Status        int    `json:"status" validate:"oneof=0 1"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name          string `json:"name" validate:"omitempty,min=2,max=64"`
	Description   string `json:"description" validate:"max=500"`
	PermissionIDs []uint `json:"permission_ids"`
	SortOrder     int    `json:"sort_order"`
	Status        int    `json:"status" validate:"oneof=0 1"`
}

// RolePermissionRequest 角色权限操作请求
type RolePermissionRequest struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}

// ListRolesRequest 角色列表查询请求
type ListRolesRequest struct {
	PaginationRequest
	Search    string `form:"search"`
	Status    *int   `form:"status" validate:"omitempty,oneof=0 1"`
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
	Status      int    `json:"status" validate:"oneof=0 1"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=64"`
	Description string `json:"description" validate:"max=500"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status" validate:"oneof=0 1"`
}

// ListPermissionsRequest 权限列表查询请求
type ListPermissionsRequest struct {
	PaginationRequest
	Search    string `form:"search"`
	Module    string `form:"module"`
	Scope     string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status    *int   `form:"status" validate:"omitempty,oneof=0 1"`
	StartTime string `form:"start_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime   string `form:"end_time" validate:"omitempty,datetime=2006-01-02 15:04:05"`
}

// PermissionTreeRequest 权限树查询请求
type PermissionTreeRequest struct {
	Scope  string `form:"scope" validate:"omitempty,oneof=platform project"`
	Status *int   `form:"status" validate:"omitempty,oneof=0 1"`
}

// CreateAPITokenRequest 创建API令牌请求
type CreateAPITokenRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=64"`
	Description string     `json:"description" validate:"max=500"`
	Scopes      []string   `json:"scopes" validate:"required,min=1"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// UpdateAPITokenRequest 更新API令牌请求
type UpdateAPITokenRequest struct {
	Name        string     `json:"name" validate:"omitempty,min=2,max=64"`
	Description string     `json:"description" validate:"max=500"`
	Scopes      []string   `json:"scopes"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Status      int        `json:"status" validate:"oneof=0 1"`
}

// ListAPITokensRequest API令牌列表查询请求
type ListAPITokensRequest struct {
	PaginationRequest
	Search  string `form:"search"`
	UserID  *uint  `form:"user_id"`
	Status  *int   `form:"status" validate:"omitempty,oneof=0 1"`
	Expired *bool  `form:"expired"`
}

// BatchRevokeAPITokensRequest 批量撤销API令牌请求
type BatchRevokeAPITokensRequest struct {
	IDs []uint `json:"ids" validate:"required,min=1"`
}

// UpdateAuthConfigRequest 更新认证配置请求
type UpdateAuthConfigRequest struct {
	APIAuth struct {
		TokenAuthEnabled bool `json:"token_auth_enabled"`
		OAuth2Enabled    bool `json:"oauth2_enabled"`
		JWTConfig        struct {
			Algorithm        string `json:"algorithm"`
			ExpiresIn        int    `json:"expires_in"`
			RefreshExpiresIn int    `json:"refresh_expires_in"`
			AutoRefresh      bool   `json:"auto_refresh"`
		} `json:"jwt_config"`
	} `json:"api_auth"`

	UserAuth struct {
		OAuth2Enabled          bool                    `json:"oauth2_enabled"`
		OAuth2Providers        []OAuth2ProviderRequest `json:"oauth2_providers"`
		TwoFactorEnabled       bool                    `json:"two_factor_enabled"`
		TwoFactorMethods       []string                `json:"two_factor_methods"`
		TwoFactorRequiredRoles []string                `json:"two_factor_required_roles"`
		PasswordPolicy         PasswordPolicyRequest   `json:"password_policy"`
		LoginSecurity          LoginSecurityRequest    `json:"login_security"`
	} `json:"user_auth"`

	SessionConfig struct {
		Timeout               int  `json:"timeout"`
		MaxConcurrentSessions int  `json:"max_concurrent_sessions"`
		RememberMeEnabled     bool `json:"remember_me_enabled"`
		RememberMeDuration    int  `json:"remember_me_duration"`
	} `json:"session_config"`
}

// OAuth2ProviderRequest OAuth2提供商请求
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

// PasswordPolicyRequest 密码策略请求
type PasswordPolicyRequest struct {
	MinLength           int  `json:"min_length" validate:"min=1,max=128"`
	MaxLength           int  `json:"max_length" validate:"min=1,max=128"`
	RequireUppercase    bool `json:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers"`
	RequireSymbols      bool `json:"require_symbols"`
	PasswordHistory     int  `json:"password_history" validate:"min=0,max=20"`
	PasswordExpiresDays int  `json:"password_expires_days" validate:"min=0,max=365"`
}

// LoginSecurityRequest 登录安全请求
type LoginSecurityRequest struct {
	MaxLoginAttempts     int      `json:"max_login_attempts" validate:"min=1,max=20"`
	LockoutDuration      int      `json:"lockout_duration" validate:"min=60,max=86400"`
	IPWhitelistEnabled   bool     `json:"ip_whitelist_enabled"`
	IPWhitelist          []string `json:"ip_whitelist"`
	LoginTimeRestriction bool     `json:"login_time_restriction"`
	AllowedLoginHours    string   `json:"allowed_login_hours"`
}

// EnableTwoFactorRequest 启用双因子认证请求
type EnableTwoFactorRequest struct {
	Method string `json:"method" validate:"required,oneof=TOTP EMAIL"`
	Code   string `json:"code" validate:"required"`
}

// PaginationRequest 分页请求基础结构
type PaginationRequest struct {
	Page     int `form:"page" validate:"min=1"`
	PageSize int `form:"page_size" validate:"min=1,max=100"`
}

// GetPage 获取页码，默认为1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页大小，默认为20
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// ValidateAPITokenRequest API令牌验证请求
type ValidateAPITokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// ConfirmTOTPRequest TOTP确认请求
type ConfirmTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// DisableTOTPRequest 禁用TOTP请求
type DisableTOTPRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// EnableEmailTwoFactorRequest 启用邮箱双因子认证请求
type EnableEmailTwoFactorRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyTwoFactorRequest 验证双因子认证请求
type VerifyTwoFactorRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required"`
	Method string `json:"method" validate:"required,oneof=totp email backup"`
}
