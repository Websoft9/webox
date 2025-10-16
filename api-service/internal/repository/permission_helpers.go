package repository

import (
	"context"

	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// GetRolesByPermissionID gets roles associated with a permission with pagination
func GetRolesByPermissionID(
	ctx context.Context,
	db *gorm.DB,
	permissionID uint,
	page, pageSize int,
) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	// Count query
	countQuery := db.WithContext(ctx).Model(&model.Role{}).
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_code  = permissions.code").
		Where("permissions.id = ? AND role_permissions.status != -1 AND roles.status != -1", permissionID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Paginated query
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_code  = permissions.code").
		Where("permissions.id = ? AND role_permissions.status != -1 AND roles.status != -1", permissionID).
		Offset(offset).Limit(pageSize).
		Find(&roles).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return roles, total, nil
}

// GetUsersByRoleID gets users associated with a role with pagination
func GetUsersByRoleID(
	ctx context.Context,
	db *gorm.DB,
	roleID uint,
	page, pageSize int,
) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Count query
	countQuery := db.WithContext(ctx).Model(&model.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND user_roles.status != -1 AND users.status != -1", roleID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Paginated query
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND user_roles.status != -1 AND users.status != -1", roleID).
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return users, total, nil
}
