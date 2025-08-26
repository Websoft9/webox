package response

import (
	"api-service/internal/model"
	"time"
)

// RoleResponse 角色响应
type RoleResponse struct {
	ID              uint                 `json:"id"`
	Name            string               `json:"name"`
	Code            string               `json:"code"`
	Description     string               `json:"description"`
	IsSystem        bool                 `json:"is_system"`
	SortOrder       int                  `json:"sort_order"`
	Status          int                  `json:"status"`
	PermissionCount int64                `json:"permission_count,omitempty"`
	UserCount       int64                `json:"user_count,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	CreatedBy       string               `json:"created_by,omitempty"`
	Permissions     []PermissionResponse `json:"permissions,omitempty"`
	Users           []UserSimpleResponse `json:"users,omitempty"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Items      []RoleResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint                 `json:"id"`
	ParentID    *uint                `json:"parent_id,omitempty"`
	Scope       string               `json:"scope"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Module      string               `json:"module"`
	Action      string               `json:"action"`
	Resource    string               `json:"resource,omitempty"`
	Description string               `json:"description,omitempty"`
	IsSystem    bool                 `json:"is_system"`
	IsMenu      bool                 `json:"is_menu"`
	SortOrder   int                  `json:"sort_order"`
	Status      int                  `json:"status"`
	RoleCount   int64                `json:"role_count,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Children    []PermissionResponse `json:"children,omitempty"`
	Roles       []RoleSimpleResponse `json:"roles,omitempty"`
}

// PermissionListResponse 权限列表响应
type PermissionListResponse struct {
	Items      []PermissionResponse `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// PermissionTreeResponse represents permission tree structure
type PermissionTreeResponse = PermissionResponse

// APITokenResponse API令牌响应
type APITokenResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Token       string     `json:"token,omitempty"` // 只在创建时返回完整token
	UserID      uint       `json:"user_id"`
	Username    string     `json:"username"`
	Scopes      []string   `json:"scopes"`
	Description string     `json:"description,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Status      int        `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

// APITokenListResponse API令牌列表响应
type APITokenListResponse struct {
	Items      []APITokenResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// AuthConfigResponse 认证配置响应
type AuthConfigResponse struct {
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
		OAuth2Enabled          bool                     `json:"oauth2_enabled"`
		OAuth2Providers        []OAuth2ProviderResponse `json:"oauth2_providers"`
		TwoFactorEnabled       bool                     `json:"two_factor_enabled"`
		TwoFactorMethods       []string                 `json:"two_factor_methods"`
		TwoFactorRequiredRoles []string                 `json:"two_factor_required_roles"`
		PasswordPolicy         PasswordPolicyResponse   `json:"password_policy"`
		LoginSecurity          LoginSecurityResponse    `json:"login_security"`
	} `json:"user_auth"`

	SessionConfig struct {
		Timeout               int  `json:"timeout"`
		MaxConcurrentSessions int  `json:"max_concurrent_sessions"`
		RememberMeEnabled     bool `json:"remember_me_enabled"`
		RememberMeDuration    int  `json:"remember_me_duration"`
	} `json:"session_config"`
}

// OAuth2ProviderResponse OAuth2提供商响应
type OAuth2ProviderResponse struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret,omitempty"` // 敏感信息用***替代
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	Enabled      bool              `json:"enabled"`
	AutoRegister bool              `json:"auto_register"`
	UserMapping  map[string]string `json:"user_mapping"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// PasswordPolicyResponse 密码策略响应
type PasswordPolicyResponse struct {
	MinLength           int  `json:"min_length"`
	MaxLength           int  `json:"max_length"`
	RequireUppercase    bool `json:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers"`
	RequireSymbols      bool `json:"require_symbols"`
	PasswordHistory     int  `json:"password_history"`
	PasswordExpiresDays int  `json:"password_expires_days"`
}

// LoginSecurityResponse 登录安全响应
type LoginSecurityResponse struct {
	MaxLoginAttempts     int      `json:"max_login_attempts"`
	LockoutDuration      int      `json:"lockout_duration"`
	IPWhitelistEnabled   bool     `json:"ip_whitelist_enabled"`
	IPWhitelist          []string `json:"ip_whitelist"`
	LoginTimeRestriction bool     `json:"login_time_restriction"`
	AllowedLoginHours    string   `json:"allowed_login_hours"`
}

// TwoFactorStatusResponse represents two-factor authentication status
type TwoFactorStatusResponse struct {
	Enabled     bool                      `json:"enabled"`
	Methods     []TwoFactorMethodResponse `json:"methods"`
	BackupCodes int                       `json:"backup_codes"`
}

// TOTPSecretResponse represents TOTP secret generation response
type TOTPSecretResponse struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFactorMethodResponse 双因子认证方法响应
type TwoFactorMethodResponse struct {
	Type       string `json:"type"`
	Enabled    bool   `json:"enabled"`
	Configured bool   `json:"configured"`
	Email      string `json:"email,omitempty"`
}

// UserSimpleResponse 用户简单响应
type UserSimpleResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RoleSimpleResponse 角色简单响应
type RoleSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ConvertToRoleResponse 转换为角色响应
func ConvertToRoleResponse(role *model.Role) *RoleResponse {
	resp := &RoleResponse{
		ID:              role.ID,
		Name:            role.Name,
		Code:            role.Code,
		Description:     role.Description,
		IsSystem:        role.IsSystem,
		SortOrder:       role.SortOrder,
		Status:          role.Status,
		PermissionCount: role.PermissionCount,
		UserCount:       role.UserCount,
		CreatedAt:       role.CreatedAt,
		UpdatedAt:       role.UpdatedAt,
	}

	// 转换权限列表
	if len(role.Permissions) > 0 {
		resp.Permissions = make([]PermissionResponse, len(role.Permissions))
		for i, perm := range role.Permissions {
			resp.Permissions[i] = *ConvertToPermissionResponse(&perm)
		}
	}

	// 转换用户列表
	if len(role.Users) > 0 {
		resp.Users = make([]UserSimpleResponse, len(role.Users))
		for i, user := range role.Users {
			resp.Users[i] = UserSimpleResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}
		}
	}

	return resp
}

// ConvertToPermissionResponse 转换为权限响应
func ConvertToPermissionResponse(perm *model.Permission) *PermissionResponse {
	resp := &PermissionResponse{
		ID:          perm.ID,
		ParentID:    perm.ParentID,
		Scope:       perm.Scope,
		Name:        perm.Name,
		Code:        perm.Code,
		Module:      perm.Module,
		Action:      perm.Action,
		Resource:    perm.Resource,
		Description: perm.Description,
		IsSystem:    perm.IsSystem,
		IsMenu:      perm.IsMenu,
		SortOrder:   perm.SortOrder,
		Status:      perm.Status,
		RoleCount:   perm.RoleCount,
		CreatedAt:   perm.CreatedAt,
		UpdatedAt:   perm.UpdatedAt,
	}

	// 转换子权限
	if len(perm.Children) > 0 {
		resp.Children = make([]PermissionResponse, len(perm.Children))
		for i, child := range perm.Children {
			resp.Children[i] = *ConvertToPermissionResponse(&child)
		}
	}

	// 转换角色列表
	if len(perm.Roles) > 0 {
		resp.Roles = make([]RoleSimpleResponse, len(perm.Roles))
		for i, role := range perm.Roles {
			resp.Roles[i] = RoleSimpleResponse{
				ID:   role.ID,
				Name: role.Name,
				Code: role.Code,
			}
		}
	}

	return resp
}

// ConvertToAPITokenResponse converts API token model to response
func ConvertToAPITokenResponse(token *model.APIToken) *APITokenResponse {
	resp := &APITokenResponse{
		ID:          token.ID,
		Name:        token.Name,
		Token:       "***", // Mask token for security
		UserID:      token.UserID,
		Description: token.Description,
		LastUsedAt:  token.LastUsedAt,
		ExpiresAt:   token.ExpiresAt,
		Status:      token.Status,
		CreatedAt:   token.CreatedAt,
	}

	// Convert scopes from JSON to string slice
	if token.Scopes != nil {
		if scopesData, exists := token.Scopes["scopes"]; exists {
			if scopesSlice, ok := scopesData.([]interface{}); ok {
				resp.Scopes = make([]string, len(scopesSlice))
				for i, scope := range scopesSlice {
					if s, ok := scope.(string); ok {
						resp.Scopes[i] = s
					}
				}
			}
		}
	}

	return resp
}

// ConvertToPermissionTreeResponse converts permissions to tree structure
func ConvertToPermissionTreeResponse(permissions []*model.Permission) []*PermissionTreeResponse {
	tree := make([]*PermissionTreeResponse, len(permissions))
	for i, perm := range permissions {
		resp := ConvertToPermissionResponse(perm)
		tree[i] = (*PermissionTreeResponse)(resp)
	}
	return tree
}

// ConvertToTwoFactorStatusResponse converts two-factor status to response
func ConvertToTwoFactorStatusResponse(methods []*model.UserTwoFactor) *TwoFactorStatusResponse {
	resp := &TwoFactorStatusResponse{
		Enabled: false,
		Methods: make([]TwoFactorMethodResponse, 0),
	}

	for _, method := range methods {
		methodResp := TwoFactorMethodResponse{
			Type:       method.Method,
			Enabled:    method.Enabled,
			Configured: true,
		}

		if method.Method == "EMAIL" && method.Email != "" {
			methodResp.Email = method.Email
		}

		resp.Methods = append(resp.Methods, methodResp)

		if method.Enabled {
			resp.Enabled = true
		}
	}

	return resp
}

// APITokenValidationResponse API令牌验证响应
type APITokenValidationResponse struct {
	Valid     bool      `json:"valid"`
	UserID    uint      `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Scopes    []string  `json:"scopes,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// TOTPSetupResponse TOTP设置响应
type TOTPSetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// TOTPConfirmResponse TOTP确认响应
type TOTPConfirmResponse struct {
	Enabled     bool     `json:"enabled"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFactorVerificationResponse 双因子认证验证响应
type TwoFactorVerificationResponse struct {
	Valid  bool   `json:"valid"`
	UserID uint   `json:"user_id,omitempty"`
	Method string `json:"method,omitempty"`
}

// BackupCodesResponse 备用码响应
type BackupCodesResponse struct {
	Codes []string `json:"codes"`
}
