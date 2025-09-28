package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"
	"time"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a role repository instance
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a role
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("status != -1").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &role, nil
}

// GetByCode retrieves a role by code
func (r *roleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ? AND status != -1", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &role, nil
}

// Update updates a role
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	// Check if role exists and is not deleted or disabled
	var existingRole model.Role
	if err := r.db.WithContext(ctx).Where("status != -1").First(&existingRole, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound)
		}
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Check if role is disabled
	if existingRole.Status == 0 {
		return errors.NewAppError(errors.CodeRecordIsDisabled)
	}

	result := r.db.WithContext(ctx).Model(role).
		Omit("created_at", "code").
		Updates(role)

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}

// Delete deletes a role (soft delete)
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	// Check if role exists and is not already deleted
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound)
		}
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	// Check if already deleted
	if role.Status == -1 {
		return errors.NewAppError(errors.CodeRecordDeleteFailed)
	}

	// Check if there are associated users
	var userCount int64
	if err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	if userCount > 0 {
		return errors.NewAppError(errors.CodeRecordDeleteFailed)
	}

	if role.IsSystem {
		return errors.NewAppError(errors.CodeRecordDeleteDenied)
	}

	// Perform soft delete by setting status to -1
	result := r.db.WithContext(ctx).Model(&model.Role{}).
		Where("id = ?", id).
		Update("status", -1)

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}

// List retrieves a list of roles
func (r *roleRepository) List(ctx context.Context, req *request.ListRolesRequest) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Role{}).Where("status != -1")

	// Build query conditions
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

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&roles).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Get statistical information
	for _, role := range roles {
		// Get permission count
		var permCount int64
		r.db.WithContext(ctx).Model(&model.RolePermission{}).
			Where("role_id = ? AND status = 1", role.ID).Count(&permCount)
		role.PermissionCount = permCount

		// Get user count
		var userCount int64
		r.db.WithContext(ctx).Model(&model.UserRole{}).
			Where("role_id = ? AND status = 1", role.ID).Count(&userCount)
		role.UserCount = userCount
	}

	return roles, total, nil
}

// GetWithPermissions retrieves a role with its permissions
func (r *roleRepository) GetWithPermissions(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions", "status != -1").
		First(&role, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &role, nil
}

// GetWithUsers retrieves a role with its users
func (r *roleRepository) GetWithUsers(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Users", "status != -1").
		First(&role, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &role, nil
}

// AssignPermissionsWithTx assigns permissions to a role within an existing transaction
func (r *roleRepository) AssignPermissionsWithTx(ctx context.Context, tx *gorm.DB, roleID uint, permissionIDs []uint, grantedBy uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	// Delete existing permission associations
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	// Get permission codes first
	var permissions []*model.Permission
	if err := tx.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Create new permission associations
	rolePermissions := make([]model.RolePermission, len(permissions))
	for i, perm := range permissions {
		rolePermissions[i] = model.RolePermission{
			RoleID:         roleID,
			PermissionCode: perm.Code,
			GrantedBy:      &grantedBy,
			GrantedAt:      time.Now(),
			Status:         1,
		}
	}

	if err := tx.Create(&rolePermissions).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// RemovePermissions removes role permissions
func (r *roleRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_code IN ?", roleID, permissionIDs).
		Delete(&model.RolePermission{})

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}

	return nil
}

// GetPermissions retrieves role permissions
func (r *roleRepository) GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.code = role_permissions.permission_code").
		Where("role_permissions.role_id = ? AND role_permissions.status != -1 AND permissions.status != -1", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return permissions, nil
}

// GetUsers retrieves role users
func (r *roleRepository) GetUsers(ctx context.Context, roleID uint, page, pageSize int) ([]model.User, int64, error) {
	return GetUsersByRoleID(ctx, r.db, roleID, page, pageSize)
}

// CountPermissions counts the number of permissions for a role
func (r *roleRepository) CountPermissions(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("role_id = ? AND status = 1", roleID).Count(&count).Error

	if err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count, nil
}

// CountUsers counts the number of users for a role
func (r *roleRepository) CountUsers(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ? AND status = 1", roleID).Count(&count).Error

	if err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count, nil
}

// BatchUpdateStatus updates status in batch
func (r *roleRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.Role{}).
		Where("id IN ? AND is_system = false", ids).
		Update("status", status)

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// CreateWithTx creates a role within a transaction
func (r *roleRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	if err := tx.WithContext(ctx).Create(role).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// UpdateWithTx updates a role within a transaction
func (r *roleRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	// Check if role exists and is not deleted or disabled
	var existingRole model.Role
	if err := tx.WithContext(ctx).Where("status != -1").First(&existingRole, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound)
		}
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Check if role is disabled
	if existingRole.Status == 0 {
		return errors.NewAppError(errors.CodeRecordIsDisabled)
	}

	result := tx.WithContext(ctx).Model(role).
		Omit("created_at", "code").
		Updates(role)

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}
