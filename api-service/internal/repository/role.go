package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色存储实例
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return errors.Wrap(err, "failed to create role")
	}
	return nil
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role by ID")
	}
	return &role, nil
}

// GetByCode 根据代码获取角色
func (r *roleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get role by code")
	}
	return &role, nil
}

// Update 更新角色
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	result := r.db.WithContext(ctx).Model(role).
		Omit("created_at", "code").
		Updates(role)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update role")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// Delete 删除角色
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	// 检查是否有关联用户
	var userCount int64
	if err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return errors.Wrap(err, "failed to check role users")
	}

	if userCount > 0 {
		return errors.New("cannot delete role with associated users")
	}

	// 检查是否为系统角色
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return errors.Wrap(err, "failed to get role")
	}

	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	// 开始事务删除
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return errors.Wrap(err, "failed to delete role permissions")
		}

		// 删除角色
		if err := tx.Delete(&model.Role{}, id).Error; err != nil {
			return errors.Wrap(err, "failed to delete role")
		}

		return nil
	})
}

// List 获取角色列表
func (r *roleRepository) List(ctx context.Context, req *request.ListRolesRequest) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Role{})

	// 构建查询条件
	if req.Search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.StartTime != "" && req.EndTime != "" {
		startTime, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
		endTime, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)
		query = query.Where("updated_at BETWEEN ? AND ?", startTime, endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to count roles")
	}

	// 分页查询
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&roles).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list roles")
	}

	// 获取统计信息
	for _, role := range roles {
		// 获取权限数量
		var permCount int64
		r.db.WithContext(ctx).Model(&model.RolePermission{}).
			Where("role_id = ? AND status = 1", role.ID).Count(&permCount)
		role.PermissionCount = permCount

		// 获取用户数量
		var userCount int64
		r.db.WithContext(ctx).Model(&model.UserRole{}).
			Where("role_id = ? AND status = 1", role.ID).Count(&userCount)
		role.UserCount = userCount
	}

	return roles, total, nil
}

// GetWithPermissions 获取角色及其权限
func (r *roleRepository) GetWithPermissions(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions", "status = 1").
		First(&role, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role with permissions")
	}

	return &role, nil
}

// GetWithUsers 获取角色及其用户
func (r *roleRepository) GetWithUsers(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Users", "status = 1").
		First(&role, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role with users")
	}

	return &role, nil
}

// AssignPermissions 分配权限给角色
func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, grantedBy uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return errors.Wrap(err, "failed to remove existing permissions")
		}

		// 创建新的权限关联
		rolePermissions := make([]model.RolePermission, len(permissionIDs))
		for i, permID := range permissionIDs {
			rolePermissions[i] = model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
				GrantedBy:    &grantedBy,
				GrantedAt:    time.Now(),
				Status:       1,
			}
		}

		if err := tx.Create(&rolePermissions).Error; err != nil {
			return errors.Wrap(err, "failed to assign permissions")
		}

		return nil
	})
}

// RemovePermissions 移除角色权限
func (r *roleRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&model.RolePermission{})

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to remove permissions")
	}

	return nil
}

// GetPermissions 获取角色权限
func (r *roleRepository) GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND role_permissions.status = 1 AND permissions.status = 1", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get role permissions")
	}

	return permissions, nil
}

// GetUsers 获取角色用户
func (r *roleRepository) GetUsers(ctx context.Context, roleID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 获取总数
	countQuery := r.db.WithContext(ctx).Model(&model.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND user_roles.status = 1 AND users.status = 1", roleID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to count role users")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND user_roles.status = 1 AND users.status = 1", roleID).
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get role users")
	}

	return users, total, nil
}

// CountPermissions 统计角色权限数量
func (r *roleRepository) CountPermissions(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("role_id = ? AND status = 1", roleID).Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "failed to count role permissions")
	}

	return count, nil
}

// CountUsers 统计角色用户数量
func (r *roleRepository) CountUsers(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ? AND status = 1", roleID).Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "failed to count role users")
	}

	return count, nil
}

// BatchUpdateStatus 批量更新状态
func (r *roleRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.Role{}).
		Where("id IN ? AND is_system = false", ids).
		Update("status", status)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to batch update role status")
	}

	return nil
}

// CreateWithTx 在事务中创建角色
func (r *roleRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	if err := tx.WithContext(ctx).Create(role).Error; err != nil {
		return errors.Wrap(err, "failed to create role in transaction")
	}
	return nil
}

// UpdateWithTx 在事务中更新角色
func (r *roleRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	result := tx.WithContext(ctx).Model(role).
		Omit("created_at", "code").
		Updates(role)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update role in transaction")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected in transaction")
	}

	return nil
}
