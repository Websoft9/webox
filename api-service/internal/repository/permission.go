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

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限存储实例
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

// Create 创建权限
func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		return errors.Wrap(err, "failed to create permission")
	}
	return nil
}

// GetByID 根据ID获取权限
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, errors.Wrap(err, "failed to get permission by ID")
	}
	return &permission, nil
}

// GetByCode 根据代码获取权限
func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get permission by code")
	}
	return &permission, nil
}

// Update 更新权限
func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	result := r.db.WithContext(ctx).Model(permission).
		Omit("created_at", "code", "scope", "module", "action", "resource").
		Updates(permission)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update permission")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// Delete 删除权限
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	// 检查是否有关联角色
	var roleCount int64
	if err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("permission_id = ?", id).Count(&roleCount).Error; err != nil {
		return errors.Wrap(err, "failed to check permission roles")
	}

	if roleCount > 0 {
		return errors.New("cannot delete permission with associated roles")
	}

	// 检查是否为系统权限
	var permission model.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return errors.Wrap(err, "failed to get permission")
	}

	if permission.IsSystem {
		return errors.New("cannot delete system permission")
	}

	// 检查是否有子权限
	var childCount int64
	if err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return errors.Wrap(err, "failed to check child permissions")
	}

	if childCount > 0 {
		return errors.New("cannot delete permission with child permissions")
	}

	// 开始事务删除
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除权限角色关联
		if err := tx.Where("permission_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return errors.Wrap(err, "failed to delete permission roles")
		}

		// 删除权限
		if err := tx.Delete(&model.Permission{}, id).Error; err != nil {
			return errors.Wrap(err, "failed to delete permission")
		}

		return nil
	})
}

// List 获取权限列表
func (r *permissionRepository) List(ctx context.Context, req *request.ListPermissionsRequest) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Permission{})

	// 构建查询条件
	if req.Search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}

	if req.Scope != "" {
		query = query.Where("scope = ?", req.Scope)
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
		return nil, 0, errors.Wrap(err, "failed to count permissions")
	}

	// 分页查询
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("module ASC, sort_order ASC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&permissions).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list permissions")
	}

	// 获取统计信息
	for _, permission := range permissions {
		// 获取角色数量
		var roleCount int64
		r.db.WithContext(ctx).Model(&model.RolePermission{}).
			Where("permission_id = ? AND status = 1", permission.ID).Count(&roleCount)
		permission.RoleCount = roleCount
	}

	return permissions, total, nil
}

// GetTree 获取权限树
func (r *permissionRepository) GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error) {
	var permissions []*model.Permission

	query := r.db.WithContext(ctx).Model(&model.Permission{})

	// 构建查询条件
	if req.Scope != "" {
		query = query.Where("scope = ?", req.Scope)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		query = query.Where("status = 1") // 默认只显示启用的权限
	}

	// 获取所有权限
	err := query.Order("module ASC, sort_order ASC, created_at ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get permission tree")
	}

	// 构建树形结构
	return r.buildPermissionTree(permissions), nil
}

// buildPermissionTree 构建权限树
func (r *permissionRepository) buildPermissionTree(permissions []*model.Permission) []*model.Permission {
	// 创建ID到权限的映射
	permissionMap := make(map[uint]*model.Permission)
	for _, perm := range permissions {
		permissionMap[perm.ID] = perm
	}

	// 构建树形结构
	var roots []*model.Permission
	for _, perm := range permissions {
		if perm.ParentID == nil {
			// 根节点
			roots = append(roots, perm)
		} else {
			// 子节点
			if parent, exists := permissionMap[*perm.ParentID]; exists {
				parent.Children = append(parent.Children, *perm)
			}
		}
	}

	return roots
}

// GetWithRoles 获取权限及其角色
func (r *permissionRepository) GetWithRoles(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).
		Preload("Roles", "status = 1").
		First(&permission, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, errors.Wrap(err, "failed to get permission with roles")
	}

	return &permission, nil
}

// GetChildren 获取子权限
func (r *permissionRepository) GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND status = 1", parentID).
		Order("sort_order ASC, created_at ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get child permissions")
	}

	return permissions, nil
}

// GetRoles 获取权限关联的角色
func (r *permissionRepository) GetRoles(ctx context.Context, permissionID uint, page, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	// 获取总数
	countQuery := r.db.WithContext(ctx).Model(&model.Role{}).
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Where("role_permissions.permission_id = ? AND role_permissions.status = 1 AND roles.status = 1", permissionID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to count permission roles")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Where("role_permissions.permission_id = ? AND role_permissions.status = 1 AND roles.status = 1", permissionID).
		Offset(offset).Limit(pageSize).
		Find(&roles).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get permission roles")
	}

	return roles, total, nil
}

// CountRoles 统计权限角色数量
func (r *permissionRepository) CountRoles(ctx context.Context, permissionID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("permission_id = ? AND status = 1", permissionID).Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "failed to count permission roles")
	}

	return count, nil
}

// BatchUpdateStatus 批量更新状态
func (r *permissionRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("id IN ? AND is_system = false", ids).
		Update("status", status)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to batch update permission status")
	}

	return nil
}

// GetByIDs 根据ID列表获取权限
func (r *permissionRepository) GetByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error) {
	if len(ids) == 0 {
		return []*model.Permission{}, nil
	}

	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Where("id IN ? AND status = 1", ids).Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get permissions by IDs")
	}

	return permissions, nil
}

// GetUserPermissions 获取用户权限
func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.status = 1 AND role_permissions.status = 1 AND permissions.status = 1", userID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user permissions")
	}

	return permissions, nil
}

// CheckUserPermission 检查用户权限
func (r *permissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Where("user_roles.status = 1 AND role_permissions.status = 1 AND permissions.status = 1").
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check user permission")
	}

	return count > 0, nil
}

// CreateWithTx 在事务中创建权限
func (r *permissionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	if err := tx.WithContext(ctx).Create(permission).Error; err != nil {
		return errors.Wrap(err, "failed to create permission in transaction")
	}
	return nil
}

// UpdateWithTx 在事务中更新权限
func (r *permissionRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	result := tx.WithContext(ctx).Model(permission).
		Omit("created_at", "code", "scope", "module", "action", "resource").
		Updates(permission)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update permission in transaction")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected in transaction")
	}

	return nil
}
