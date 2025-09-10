package response

import (
	"api-service/internal/model"
	"time"
)

// RoleResponse role response
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
	Permissions     []PermissionResponse `json:"permissions,omitempty"`
	Users           []UserSimpleResponse `json:"users,omitempty"`
}

// RoleListResponse role list response
type RoleListResponse struct {
	Items      []RoleResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// PermissionResponse permission response
type PermissionResponse struct {
	ID          uint                 `json:"id"`
	ParentCode  string               `json:"parent_code,omitempty"`
	Scope       string               `json:"scope"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Module      string               `json:"module"`
	Action      string               `json:"action"`
	Resource    string               `json:"resource,omitempty"`
	Element     string               `json:"element,omitempty"`
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

// PermissionListResponse permission list response
type PermissionListResponse struct {
	Items      []PermissionResponse `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// PermissionTreeResponse represents permission tree structure
type PermissionTreeResponse = PermissionResponse

// APITokenResponse API token response
type APITokenResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Token       string     `json:"token,omitempty"` // only return complete token when creates
	UserID      uint       `json:"user_id"`
	Username    string     `json:"username"`
	Scopes      []string   `json:"scopes"`
	Description string     `json:"description,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// APITokenListResponse API token list response
type APITokenListResponse struct {
	Items      []APITokenResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// AuthConfigResponse authentication config response
type AuthConfigResponse struct {
	APIAuth       APIAuthResponse       `json:"api_auth"`
	UserAuth      UserAuthResponse      `json:"user_auth"`
	SessionConfig SessionConfigResponse `json:"session_config"`
}

// APIAuthResponse API authentication response
type APIAuthResponse struct {
	TokenAuthEnabled bool              `json:"token_auth_enabled"`
	OAuth2Enabled    bool              `json:"oauth2_enabled"`
	JWTConfig        JWTConfigResponse `json:"jwt_config"`
}

// JWTConfigResponse JWT configuration response
type JWTConfigResponse struct {
	Algorithm        string `json:"algorithm"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	AutoRefresh      bool   `json:"auto_refresh"`
}

// UserAuthResponse user authentication response
type UserAuthResponse struct {
	OAuth2Enabled          bool                      `json:"oauth2_enabled"`
	OAuth2Providers        []*OAuth2ProviderResponse `json:"oauth2_providers"`
	TwoFactorEnabled       bool                      `json:"two_factor_enabled"`
	TwoFactorMethods       []string                  `json:"two_factor_methods"`
	TwoFactorRequiredRoles []string                  `json:"two_factor_required_roles"`
	PasswordPolicy         PasswordPolicyResponse    `json:"password_policy"`
	LoginSecurity          LoginSecurityResponse     `json:"login_security"`
}

// SessionConfigResponse session configuration response
type SessionConfigResponse struct {
	Timeout               int  `json:"timeout"`
	MaxConcurrentSessions int  `json:"max_concurrent_sessions"`
	RememberMeEnabled     bool `json:"remember_me_enabled"`
	RememberMeDuration    int  `json:"remember_me_duration"`
}

// OAuth2ProviderResponse OAuth2 provider response
type OAuth2ProviderResponse struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret,omitempty"` // sensitive information replaced with ***
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	Enabled      bool              `json:"enabled"`
	AutoRegister bool              `json:"auto_register"`
	UserMapping  map[string]string `json:"user_mapping"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// PasswordPolicyResponse password policy response
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

// LoginSecurityResponse login security response
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

// TwoFactorMethodResponse two-factor authentication method response
type TwoFactorMethodResponse struct {
	Type       string `json:"type"`
	Enabled    bool   `json:"enabled"`
	Configured bool   `json:"configured"`
	Email      string `json:"email,omitempty"`
}

// UserSimpleResponse user simple response
type UserSimpleResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RoleSimpleResponse role simple response
type RoleSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ConvertToRoleResponse converts to role response
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

	// converts permission list
	if len(role.Permissions) > 0 {
		resp.Permissions = make([]PermissionResponse, len(role.Permissions))
		for i := range role.Permissions {
			resp.Permissions[i] = *ConvertToPermissionResponse(&role.Permissions[i])
		}
	}

	// converts user list
	if len(role.Users) > 0 {
		resp.Users = make([]UserSimpleResponse, len(role.Users))
		for i := range role.Users {
			resp.Users[i] = UserSimpleResponse{
				ID:       role.Users[i].ID,
				Username: role.Users[i].Username,
				Email:    role.Users[i].Email,
			}
		}
	}

	return resp
}

// ConvertToPermissionResponse converts to permission response
func ConvertToPermissionResponse(perm *model.Permission) *PermissionResponse {
	resp := &PermissionResponse{
		ID:          perm.ID,
		ParentCode:  perm.ParentCode,
		Scope:       perm.Scope,
		Name:        perm.Name,
		Code:        perm.Code,
		Module:      perm.Module,
		Action:      perm.Action,
		Resource:    perm.Resource,
		Element:     perm.Element,
		Description: perm.Description,
		IsSystem:    perm.IsSystem,
		IsMenu:      perm.IsMenu,
		SortOrder:   perm.SortOrder,
		Status:      perm.Status,
		RoleCount:   perm.RoleCount,
		CreatedAt:   perm.CreatedAt,
		UpdatedAt:   perm.UpdatedAt,
	}

	// converts child permissions
	if len(perm.Children) > 0 {
		resp.Children = make([]PermissionResponse, len(perm.Children))
		for i := range perm.Children {
			resp.Children[i] = *ConvertToPermissionResponse(perm.Children[i])
		}
	}

	// converts role list
	if len(perm.Roles) > 0 {
		resp.Roles = make([]RoleSimpleResponse, len(perm.Roles))
		for i := range perm.Roles {
			resp.Roles[i] = RoleSimpleResponse{
				ID:   perm.Roles[i].ID,
				Name: perm.Roles[i].Name,
				Code: perm.Roles[i].Code,
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
		CreatedAt:   token.CreatedAt,
	}

	// Convert scopes from JSON to string slice
	resp.Scopes = extractScopes(token.Scopes)

	return resp
}

// extractScopes safely extracts scopes from JSON data
func extractScopes(scopesJSON map[string]interface{}) []string {
	if scopesJSON == nil {
		return []string{}
	}

	scopesData, exists := scopesJSON["scopes"]
	if !exists {
		return []string{}
	}

	scopesSlice, ok := scopesData.([]interface{})
	if !ok {
		return []string{}
	}

	scopes := make([]string, 0, len(scopesSlice))
	for _, scope := range scopesSlice {
		if s, ok := scope.(string); ok {
			scopes = append(scopes, s)
		}
	}

	return scopes
}

// ConvertToPermissionTreeResponse converts permissions to tree structure
func ConvertToPermissionTreeResponse(permissions []*model.Permission) []*PermissionTreeResponse {
	tree := make([]*PermissionTreeResponse, len(permissions))
	for i, perm := range permissions {
		resp := ConvertToPermissionResponse(perm)
		tree[i] = &PermissionTreeResponse{
			ID:          resp.ID,
			Name:        resp.Name,
			Code:        resp.Code,
			Scope:       resp.Scope,
			Module:      resp.Module,
			Action:      resp.Action,
			Resource:    resp.Resource,
			Element:     resp.Element,
			Description: resp.Description,
			ParentCode:  resp.ParentCode,
			IsSystem:    resp.IsSystem,
			IsMenu:      resp.IsMenu,
			SortOrder:   resp.SortOrder,
			Status:      resp.Status,
			Children:    resp.Children,
			Roles:       resp.Roles,
			CreatedAt:   resp.CreatedAt,
			UpdatedAt:   resp.UpdatedAt,
		}
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

// APITokenValidationResponse API token validation response
type APITokenValidationResponse struct {
	Valid     bool      `json:"valid"`
	UserID    uint      `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Scopes    []string  `json:"scopes,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// TOTPSetupResponse TOTP setup response
type TOTPSetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// TOTPConfirmResponse TOTP confirm response
type TOTPConfirmResponse struct {
	Enabled     bool     `json:"enabled"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFactorVerificationResponse two-factor authentication verify response
type TwoFactorVerificationResponse struct {
	Valid  bool   `json:"valid"`
	UserID uint   `json:"user_id,omitempty"`
	Method string `json:"method,omitempty"`
}

// BackupCodesResponse backup codes response
type BackupCodesResponse struct {
	Codes []string `json:"codes"`
}
