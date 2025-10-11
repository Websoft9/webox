package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// RoleService role service interface
type RoleService interface {
	// Basic CRUD operations
	CreateRole(ctx context.Context, req *request.CreateRoleRequest, createdBy uint) (*response.RoleResponse, error)
	GetRole(ctx context.Context, id uint) (*response.RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req *request.UpdateRoleRequest, updatedBy uint) (*response.RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error

	// Query operations
	ListRoles(ctx context.Context, req *request.ListRolesRequest) (*common.PaginationResponse, error)
	GetRoleWithPermissions(ctx context.Context, id uint) (*response.RoleResponse, error)
	GetRoleUsers(ctx context.Context, id uint, req *common.PaginationRequest) (*common.PaginationResponse, error)

	// Permission management
	AssignPermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest, grantedBy uint) error
	RemovePermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest) error

	// Batch operations
	BatchUpdateRoleStatus(ctx context.Context, ids []uint, status int) error
}

// PermissionService permission service interface
type PermissionService interface {
	// Basic CRUD operations
	CreatePermission(ctx context.Context, req *request.CreatePermissionRequest, createdBy uint) (*response.PermissionResponse, error)
	GetPermission(ctx context.Context, id uint) (*response.PermissionResponse, error)
	UpdatePermission(ctx context.Context, id uint, req *request.UpdatePermissionRequest, updatedBy uint) (*response.PermissionResponse, error)
	DeletePermission(ctx context.Context, id uint) error

	// Query operations
	ListPermissions(ctx context.Context, req *request.ListPermissionsRequest) (*common.PaginationResponse, error)
	GetPermissionTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*response.PermissionTreeResponse, error)
	GetPermissionRoles(ctx context.Context, id uint, req *common.PaginationRequest) (*common.PaginationResponse, error)

	// Permission verification
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*response.PermissionResponse, error)

	// Batch operations
	BatchUpdatePermissionStatus(ctx context.Context, ids []uint, status int) error
}

// APITokenService API token service interface
type APITokenService interface {
	// Token management
	RefreshUserAPIToken(ctx context.Context, userID uint) (*response.APITokenResponse, error)
	RevokeAPITokenByToken(ctx context.Context, token string, userID uint) error

	// Maintenance operations
	CleanExpiredTokens(ctx context.Context) error
	CheckTokenIsExists(ctx context.Context, token string) bool
}

// AuthConfigService authentication configuration service interface
type AuthConfigService interface {
	// Configuration management
	GetAuthConfig(ctx context.Context) (*response.AuthConfigResponse, error)
	UpdateAuthConfig(ctx context.Context, req *request.UpdateAuthConfigRequest) error

	// OAuth2 provider management
	GetOAuth2Providers(ctx context.Context) ([]*response.OAuth2ProviderResponse, error)
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
