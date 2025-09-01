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

// NewPermissionRepository creates a permission repository instance
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

// Create creates a permission
func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		return errors.Wrap(err, "failed to create permission")
	}
	return nil
}

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("status != -1").First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, errors.Wrap(err, "failed to get permission by ID")
	}
	return &permission, nil
}

// GetByCode retrieves a permission by code
func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ? AND status != -1", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get permission by code")
	}
	return &permission, nil
}

// Update updates a permission
func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	// Check if permission exists and is not deleted or disabled
	var existingPermission model.Permission
	if err := r.db.WithContext(ctx).Where("status != -1").First(&existingPermission, permission.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return errors.Wrap(err, "failed to get permission for update")
	}

	// Check if permission is disabled
	if existingPermission.Status == 0 {
		return errors.New("cannot update disabled permission")
	}

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

// Delete deletes a permission (soft delete)
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	// Check if permission exists and is not already deleted
	var permission model.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return errors.Wrap(err, "failed to get permission")
	}

	// Check if already deleted
	if permission.Status == -1 {
		return errors.New("permission already deleted")
	}

	// Check if there are associated roles
	var roleCount int64
	if err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("permission_id = ?", id).Count(&roleCount).Error; err != nil {
		return errors.Wrap(err, "failed to check permission roles")
	}

	if roleCount > 0 {
		return errors.New("cannot delete permission with associated roles")
	}

	if permission.IsSystem {
		return errors.New("cannot delete system permission")
	}

	// Check if there are child permissions that are not deleted
	var childCount int64
	if err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("parent_id = ? AND status != -1", id).Count(&childCount).Error; err != nil {
		return errors.Wrap(err, "failed to check child permissions")
	}

	if childCount > 0 {
		return errors.New("cannot delete permission with active child permissions")
	}

	// Perform soft delete by setting status to -1
	result := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("id = ?", id).
		Update("status", -1)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete permission")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// List retrieves a list of permissions
func (r *permissionRepository) List(ctx context.Context, req *request.ListPermissionsRequest) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Permission{}).Where("status != -1")

	// Build query conditions
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

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "failed to count permissions")
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("module ASC, sort_order ASC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&permissions).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list permissions")
	}

	// Get statistical information
	for _, permission := range permissions {
		// Get role count
		var roleCount int64
		r.db.WithContext(ctx).Model(&model.RolePermission{}).
			Where("permission_id = ? AND status = 1", permission.ID).Count(&roleCount)
		permission.RoleCount = roleCount
	}

	return permissions, total, nil
}

// GetTree retrieves the permission tree
func (r *permissionRepository) GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error) {
	var permissions []*model.Permission

	query := r.db.WithContext(ctx).Model(&model.Permission{}).Where("status != -1")

	// Build query conditions
	if req.Scope != "" {
		query = query.Where("scope = ?", req.Scope)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		query = query.Where("status = 1") // Show only enabled permissions by default
	}

	// Get all permissions
	err := query.Order("module ASC, sort_order ASC, created_at ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get permission tree")
	}

	// Build tree structure
	return r.buildPermissionTree(permissions), nil
}

// buildPermissionTree builds the permission tree
func (r *permissionRepository) buildPermissionTree(permissions []*model.Permission) []*model.Permission {
	// Create ID to permission mapping
	permissionMap := make(map[uint]*model.Permission)
	for _, perm := range permissions {
		permissionMap[perm.ID] = perm
	}

	// Build tree structure
	var roots []*model.Permission
	for _, perm := range permissions {
		if perm.ParentID == nil {
			// Root node
			roots = append(roots, perm)
		} else {
			// Child node
			if parent, exists := permissionMap[*perm.ParentID]; exists {
				parent.Children = append(parent.Children, *perm)
			}
		}
	}

	return roots
}

// GetWithRoles retrieves a permission with its roles
func (r *permissionRepository) GetWithRoles(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).
		Preload("Roles", "status != -1").
		First(&permission, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, errors.Wrap(err, "failed to get permission with roles")
	}

	return &permission, nil
}

// GetChildren retrieves child permissions
func (r *permissionRepository) GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND status != -1", parentID).
		Order("sort_order ASC, created_at ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get child permissions")
	}

	return permissions, nil
}

// GetRoles retrieves roles associated with a permission
func (r *permissionRepository) GetRoles(ctx context.Context, permissionID uint, page, pageSize int) ([]model.Role, int64, error) {
	return GetRolesByPermissionID(ctx, r.db, permissionID, page, pageSize)
}

// CountRoles counts the number of roles for a permission
func (r *permissionRepository) CountRoles(ctx context.Context, permissionID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("permission_id = ? AND status = 1", permissionID).Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "failed to count permission roles")
	}

	return count, nil
}

// BatchUpdateStatus updates status in batch
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

// GetByIDs retrieves permissions by ID list
func (r *permissionRepository) GetByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error) {
	if len(ids) == 0 {
		return []*model.Permission{}, nil
	}

	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Where("id IN ? AND status != -1", ids).Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get permissions by IDs")
	}

	return permissions, nil
}

// GetUserPermissions retrieves permissions for a user
func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.status != -1 AND role_permissions.status != -1 AND permissions.status != -1", userID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user permissions")
	}

	return permissions, nil
}

// CheckUserPermission checks if a user has a specific permission
func (r *permissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Where("user_roles.status = 1 AND role_permissions.status != -1 AND permissions.status != -1").
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check user permission")
	}

	return count > 0, nil
}

// CreateWithTx creates a permission within a transaction
func (r *permissionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	if err := tx.WithContext(ctx).Create(permission).Error; err != nil {
		return errors.Wrap(err, "failed to create permission in transaction")
	}
	return nil
}

// UpdateWithTx updates a permission within a transaction
func (r *permissionRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	// Check if permission exists and is not deleted or disabled
	var existingPermission model.Permission
	if err := tx.WithContext(ctx).Where("status != -1").First(&existingPermission, permission.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return errors.Wrap(err, "failed to get permission for update in transaction")
	}

	// Check if permission is disabled
	if existingPermission.Status == 0 {
		return errors.New("cannot update disabled permission")
	}

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
