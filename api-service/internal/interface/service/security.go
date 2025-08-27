package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// RoleService 角色服务接口
type RoleService interface {
	// 基础CRUD操作
	CreateRole(ctx context.Context, req *request.CreateRoleRequest, createdBy uint) (*response.RoleResponse, error)
	GetRole(ctx context.Context, id uint) (*response.RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req *request.UpdateRoleRequest, updatedBy uint) (*response.RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error

	// 查询操作
	ListRoles(ctx context.Context, req *request.ListRolesRequest) (*response.RoleListResponse, error)
	GetRoleWithPermissions(ctx context.Context, id uint) (*response.RoleResponse, error)
	GetRoleUsers(ctx context.Context, id uint, page, pageSize int) (*response.RoleListResponse, error)

	// 权限管理
	AssignPermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest, grantedBy uint) error
	RemovePermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest) error

	// 批量操作
	BatchUpdateRoleStatus(ctx context.Context, ids []uint, status int) error

	// 初始化系统角色
	InitializeSystemRoles(ctx context.Context) error
}

// PermissionService 权限服务接口
type PermissionService interface {
	// 基础CRUD操作
	CreatePermission(ctx context.Context, req *request.CreatePermissionRequest, createdBy uint) (*response.PermissionResponse, error)
	GetPermission(ctx context.Context, id uint) (*response.PermissionResponse, error)
	UpdatePermission(ctx context.Context, id uint, req *request.UpdatePermissionRequest, updatedBy uint) (*response.PermissionResponse, error)
	DeletePermission(ctx context.Context, id uint) error

	// 查询操作
	ListPermissions(ctx context.Context, req *request.ListPermissionsRequest) (*response.PermissionListResponse, error)
	GetPermissionTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*response.PermissionTreeResponse, error)
	GetPermissionRoles(ctx context.Context, id uint, page, pageSize int) (*response.RoleListResponse, error)

	// 权限验证
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*response.PermissionResponse, error)

	// 批量操作
	BatchUpdatePermissionStatus(ctx context.Context, ids []uint, status int) error

	// 初始化系统权限
	InitializeSystemPermissions(ctx context.Context) error
}

// APITokenService API token service interface
type APITokenService interface {
	// Basic CRUD operations
	CreateAPIToken(ctx context.Context, req *request.CreateAPITokenRequest, userID uint) (*response.APITokenResponse, error)
	GetAPIToken(ctx context.Context, id uint, userID uint) (*response.APITokenResponse, error)
	UpdateAPIToken(ctx context.Context, id uint, req *request.UpdateAPITokenRequest, userID uint) (*response.APITokenResponse, error)
	RevokeAPIToken(ctx context.Context, id uint, userID uint) error

	// Query operations
	ListAPITokens(ctx context.Context, req *request.ListAPITokensRequest, userID uint) (*response.APITokenListResponse, error)

	// Token management
	RefreshAPIToken(ctx context.Context, id uint, userID uint) (*response.APITokenResponse, error)
	ValidateAPIToken(ctx context.Context, token string) (*response.APITokenValidationResponse, error)
	BatchRevokeAPITokens(ctx context.Context, ids []uint, userID uint) error

	// Maintenance operations
	CleanExpiredTokens(ctx context.Context) error
}

// AuthConfigService 认证配置服务接口
type AuthConfigService interface {
	// 配置管理
	GetAuthConfig(ctx context.Context) (*response.AuthConfigResponse, error)
	UpdateAuthConfig(ctx context.Context, req *request.UpdateAuthConfigRequest) error

	// OAuth2提供商管理
	GetOAuth2Providers(ctx context.Context) ([]*response.OAuth2ProviderResponse, error)

	// 密码策略验证
	ValidatePassword(password string) error

	// 登录安全检查
	CheckLoginSecurity(ctx context.Context, userID uint, ip string) error
	RecordLoginAttempt(ctx context.Context, userID uint, success bool, ip string) error
}

// TwoFactorService two-factor authentication service interface
type TwoFactorService interface {
	// Two-factor authentication management
	GetTwoFactorStatus(ctx context.Context, userID uint) (*response.TwoFactorStatusResponse, error)
	VerifyTwoFactor(ctx context.Context, userID uint, code, method string) (*response.TwoFactorVerificationResponse, error)

	// TOTP management
	EnableTOTP(ctx context.Context, userID uint) (*response.TOTPSetupResponse, error)
	ConfirmTOTP(ctx context.Context, userID uint, code string) (*response.TOTPConfirmResponse, error)
	DisableTOTP(ctx context.Context, userID uint, code string) error

	// Email 2FA management
	EnableEmailTwoFactor(ctx context.Context, userID uint, email string) error
	DisableEmailTwoFactor(ctx context.Context, userID uint) error
	SendEmailCode(ctx context.Context, userID uint) error

	// Backup codes management
	GenerateBackupCodes(ctx context.Context, userID uint) (*response.BackupCodesResponse, error)
}
