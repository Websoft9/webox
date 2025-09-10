package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// RoleRepository role repository interface
type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id uint) (*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	List(ctx context.Context, req *request.ListRolesRequest) ([]*model.Role, int64, error)
	GetWithPermissions(ctx context.Context, id uint) (*model.Role, error)
	GetWithUsers(ctx context.Context, id uint) (*model.Role, error)

	// Permission association operations
	AssignPermissionsWithTx(ctx context.Context, tx *gorm.DB, roleID uint, permissionIDs []uint, grantedBy uint) error
	RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error)

	// User association operations
	GetUsers(ctx context.Context, roleID uint, page, pageSize int) ([]model.User, int64, error)

	// Statistics operations
	CountPermissions(ctx context.Context, roleID uint) (int64, error)
	CountUsers(ctx context.Context, roleID uint) (int64, error)

	// Batch operations
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// Transaction operations
	CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error
}

// PermissionRepository permission repository interface
type PermissionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, permission *model.Permission) error
	GetByID(ctx context.Context, id uint) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	List(ctx context.Context, req *request.ListPermissionsRequest) ([]*model.Permission, int64, error)
	GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error)
	GetWithRoles(ctx context.Context, id uint) (*model.Permission, error)
	GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error)

	// Role association operations
	GetRoles(ctx context.Context, permissionID uint, page, pageSize int) ([]model.Role, int64, error)

	// Statistics operations
	CountRoles(ctx context.Context, permissionID uint) (int64, error)

	// Batch operations
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error
	GetByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error)

	// User permission query
	GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error)
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)

	// Transaction operations
	CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error
}

// APITokenRepository API token repository interface
type APITokenRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, token *model.APIToken) error
	GetByID(ctx context.Context, id uint) (*model.APIToken, error)
	GetByToken(ctx context.Context, tokenHash string) (*model.APIToken, error)
	Update(ctx context.Context, token *model.APIToken) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	List(ctx context.Context, req *request.ListAPITokensRequest) ([]*model.APIToken, int64, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error)

	// Batch operations
	BatchDelete(ctx context.Context, ids []uint) error

	// Token management
	UpdateLastUsed(ctx context.Context, id uint, ip string) error
	CleanExpiredTokens(ctx context.Context) error

	// Transaction operations
	CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error
}

// UserTwoFactorRepository user two-factor authentication repository interface
type UserTwoFactorRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, twoFactor *model.UserTwoFactor) error
	GetByUserIDAndMethod(ctx context.Context, userID uint, method string) (*model.UserTwoFactor, error)
	Update(ctx context.Context, twoFactor *model.UserTwoFactor) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	GetByUserID(ctx context.Context, userID uint) ([]*model.UserTwoFactor, error)

	// Two-factor authentication management
	EnableMethod(ctx context.Context, userID uint, method string, secret string) error
	DisableMethod(ctx context.Context, userID uint, method string) error
	IsMethodEnabled(ctx context.Context, userID uint, method string) (bool, error)

	// Transaction operations
	CreateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error
}
