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

// NewRoleRepository creates a role repository instance
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a role
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return errors.Wrap(err, "failed to create role")
	}
	return nil
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("status != -1").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role by ID")
	}
	return &role, nil
}

// GetByCode retrieves a role by code
func (r *roleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ? AND status != -1", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get role by code")
	}
	return &role, nil
}

// Update updates a role
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	// Check if role exists and is not deleted or disabled
	var existingRole model.Role
	if err := r.db.WithContext(ctx).Where("status != -1").First(&existingRole, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return errors.Wrap(err, "failed to get role for update")
	}

	// Check if role is disabled
	if existingRole.Status == 0 {
		return errors.New("cannot update disabled role")
	}

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

// Delete deletes a role (soft delete)
func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	// Check if role exists and is not already deleted
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return errors.Wrap(err, "failed to get role")
	}

	// Check if already deleted
	if role.Status == -1 {
		return errors.New("role already deleted")
	}

	// Check if there are associated users
	var userCount int64
	if err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return errors.Wrap(err, "failed to check role users")
	}

	if userCount > 0 {
		return errors.New("cannot delete role with associated users")
	}

	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	// Perform soft delete by setting status to -1
	result := r.db.WithContext(ctx).Model(&model.Role{}).
		Where("id = ?", id).
		Update("status", -1)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete role")
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
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
		return nil, 0, errors.Wrap(err, "failed to count roles")
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&roles).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to list roles")
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
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role with permissions")
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
			return nil, errors.New("role not found")
		}
		return nil, errors.Wrap(err, "failed to get role with users")
	}

	return &role, nil
}

// AssignPermissions assigns permissions to a role
func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, grantedBy uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing permission associations
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return errors.Wrap(err, "failed to remove existing permissions")
		}

		// Create new permission associations
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

// RemovePermissions removes role permissions
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

// GetPermissions retrieves role permissions
func (r *roleRepository) GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND role_permissions.status != -1 AND permissions.status != -1", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get role permissions")
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
		return 0, errors.Wrap(err, "failed to count role permissions")
	}

	return count, nil
}

// CountUsers counts the number of users for a role
func (r *roleRepository) CountUsers(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("role_id = ? AND status = 1", roleID).Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, "failed to count role users")
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
		return errors.Wrap(result.Error, "failed to batch update role status")
	}

	return nil
}

// CreateWithTx creates a role within a transaction
func (r *roleRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	if err := tx.WithContext(ctx).Create(role).Error; err != nil {
		return errors.Wrap(err, "failed to create role in transaction")
	}
	return nil
}

// UpdateWithTx updates a role within a transaction
func (r *roleRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	// Check if role exists and is not deleted or disabled
	var existingRole model.Role
	if err := tx.WithContext(ctx).Where("status != -1").First(&existingRole, role.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return errors.Wrap(err, "failed to get role for update in transaction")
	}

	// Check if role is disabled
	if existingRole.Status == 0 {
		return errors.New("cannot update disabled role")
	}

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
