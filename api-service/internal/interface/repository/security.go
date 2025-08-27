package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// RoleRepository 角色存储接口
type RoleRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id uint) (*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error

	// 查询操作
	List(ctx context.Context, req *request.ListRolesRequest) ([]*model.Role, int64, error)
	GetWithPermissions(ctx context.Context, id uint) (*model.Role, error)
	GetWithUsers(ctx context.Context, id uint) (*model.Role, error)

	// 权限关联操作
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, grantedBy uint) error
	RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error)

	// 用户关联操作
	GetUsers(ctx context.Context, roleID uint, page, pageSize int) ([]model.User, int64, error)

	// 统计操作
	CountPermissions(ctx context.Context, roleID uint) (int64, error)
	CountUsers(ctx context.Context, roleID uint) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// 事务操作
	CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error
}

// PermissionRepository 权限存储接口
type PermissionRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, permission *model.Permission) error
	GetByID(ctx context.Context, id uint) (*model.Permission, error)
	GetByCode(ctx context.Context, code string) (*model.Permission, error)
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error

	// 查询操作
	List(ctx context.Context, req *request.ListPermissionsRequest) ([]*model.Permission, int64, error)
	GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error)
	GetWithRoles(ctx context.Context, id uint) (*model.Permission, error)
	GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error)

	// 角色关联操作
	GetRoles(ctx context.Context, permissionID uint, page, pageSize int) ([]model.Role, int64, error)

	// 统计操作
	CountRoles(ctx context.Context, permissionID uint) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error
	GetByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error)

	// 用户权限查询
	GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error)
	CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error)

	// 事务操作
	CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error
}

// APITokenRepository API令牌存储接口
type APITokenRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, token *model.APIToken) error
	GetByID(ctx context.Context, id uint) (*model.APIToken, error)
	GetByToken(ctx context.Context, tokenHash string) (*model.APIToken, error)
	Update(ctx context.Context, token *model.APIToken) error
	Delete(ctx context.Context, id uint) error

	// 查询操作
	List(ctx context.Context, req *request.ListAPITokensRequest) ([]*model.APIToken, int64, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error)

	// 批量操作
	BatchDelete(ctx context.Context, ids []uint) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// 令牌管理
	UpdateLastUsed(ctx context.Context, id uint, ip string) error
	CleanExpiredTokens(ctx context.Context) error

	// 事务操作
	CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error
}

// UserTwoFactorRepository 用户双因子认证存储接口
type UserTwoFactorRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, twoFactor *model.UserTwoFactor) error
	GetByUserIDAndMethod(ctx context.Context, userID uint, method string) (*model.UserTwoFactor, error)
	Update(ctx context.Context, twoFactor *model.UserTwoFactor) error
	Delete(ctx context.Context, id uint) error

	// 查询操作
	GetByUserID(ctx context.Context, userID uint) ([]*model.UserTwoFactor, error)

	// 双因子认证管理
	EnableMethod(ctx context.Context, userID uint, method string, secret string) error
	DisableMethod(ctx context.Context, userID uint, method string) error
	IsMethodEnabled(ctx context.Context, userID uint, method string) (bool, error)

	// 事务操作
	CreateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error
}
