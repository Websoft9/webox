package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"

	"gorm.io/gorm"
)

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	db             *gorm.DB
	logger         logger.Logger
}

// NewRoleService creates a new role service instance
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	db *gorm.DB,
	logger logger.Logger,
) service.RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		db:             db,
		logger:         logger,
	}
}

// CreateRole creates a new role
func (s *roleService) CreateRole(ctx context.Context, req *request.CreateRoleRequest, createdBy uint) (*response.RoleResponse, error) {
	s.logger.InfoContext(ctx, "Creating role",
		logger.String("operation", "CreateRole"),
		logger.String("role_code", req.Code),
		logger.Uint("created_by", createdBy))

	// Validate role data
	if err := s.validateRoleForCreation(ctx, req); err != nil {
		return nil, err
	}

	// Create and save role
	role := s.buildRoleEntity(req)
	if err := s.createRoleWithPermissions(ctx, role, req.PermissionIDs, createdBy); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Role created successfully", logger.Uint("role_id", role.ID))
	return s.GetRole(ctx, role.ID)
}

// validateRoleForCreation validates role data before creation
func (s *roleService) validateRoleForCreation(ctx context.Context, req *request.CreateRoleRequest) error {
	// Check if role code already exists
	existingRole, err := s.roleRepo.GetByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, errors.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing role", logger.ErrorField(err))
		return err
	}
	if existingRole != nil {
		s.logger.WarnContext(ctx, "Role code already exists", logger.String("role_code", req.Code))
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Validate permission IDs if provided
	return s.validatePermissionIDs(ctx, req.PermissionIDs)
}

// validatePermissionIDs validates that all permission IDs exist
func (s *roleService) validatePermissionIDs(ctx context.Context, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	permissions, err := s.permissionRepo.GetByIDs(ctx, permissionIDs)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to validate permissions", logger.ErrorField(err))
		return err
	}
	if len(permissions) != len(permissionIDs) {
		s.logger.WarnContext(ctx, "Some permission IDs are invalid",
			logger.Int("requested", len(permissionIDs)),
			logger.Int("found", len(permissions)))
		return errors.NewAppError(errors.CodePermissionInvalid)
	}
	return nil
}

// buildRoleEntity creates a role entity from request
func (s *roleService) buildRoleEntity(req *request.CreateRoleRequest) *model.Role {
	return &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsSystem:    false,
		SortOrder:   req.SortOrder,
		Status:      1,
	}
}

// createRoleWithPermissions creates role and assigns permissions in transaction
func (s *roleService) createRoleWithPermissions(ctx context.Context, role *model.Role, permissionIDs []uint, createdBy uint) error {
	return s.saveRoleWithPermissions(ctx, role, permissionIDs, createdBy, true)
}

// GetRole gets a role
func (s *roleService) GetRole(ctx context.Context, id uint) (*response.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get statistical information
	permCount, _ := s.roleRepo.CountPermissions(ctx, id)
	userCount, _ := s.roleRepo.CountUsers(ctx, id)

	role.PermissionCount = permCount
	role.UserCount = userCount

	return response.ConvertToRoleResponse(role), nil
}

// UpdateRole updates an existing role
func (s *roleService) UpdateRole(ctx context.Context, id uint, req *request.UpdateRoleRequest, updatedBy uint) (*response.RoleResponse, error) {
	s.logger.InfoContext(ctx, "Updating role",
		logger.String("operation", "UpdateRole"),
		logger.Uint("role_id", id),
		logger.Uint("updated_by", updatedBy))

	// Get and validate existing role
	role, err := s.getRoleForUpdate(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate permission IDs if provided
	if err := s.validatePermissionIDs(ctx, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Update role information
	s.updateRoleFields(role, req)

	// Save changes in transaction
	if err := s.updateRoleWithPermissions(ctx, role, req.PermissionIDs, updatedBy); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Role updated successfully", logger.Uint("role_id", id))
	return s.GetRole(ctx, id)
}

// getRoleForUpdate gets role and validates it can be updated
func (s *roleService) getRoleForUpdate(ctx context.Context, id uint) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get role for update", logger.ErrorField(err))
		return nil, err
	}

	// Check if it's a system role
	if role.IsSystem {
		s.logger.WarnContext(ctx, "Attempt to update system role", logger.Uint("role_id", id))
		return nil, errors.NewAppError(errors.CodeRecordUpdateFailed)
	}

	return role, nil
}

// updateRoleFields updates role fields from request
func (s *roleService) updateRoleFields(role *model.Role, req *request.UpdateRoleRequest) {
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.SortOrder = req.SortOrder
	role.Status = req.Status
}

// updateRoleWithPermissions updates role and permissions in transaction
func (s *roleService) updateRoleWithPermissions(ctx context.Context, role *model.Role, permissionIDs []uint, updatedBy uint) error {
	return s.saveRoleWithPermissions(ctx, role, permissionIDs, updatedBy, false)
}

// saveRoleWithPermissions saves role and assigns permissions in transaction
// isCreate determines whether to create or update the role
func (s *roleService) saveRoleWithPermissions(ctx context.Context, role *model.Role, permissionIDs []uint, operatorID uint, isCreate bool) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save role (create or update)
		var err error
		if isCreate {
			err = s.roleRepo.CreateWithTx(ctx, tx, role)
		} else {
			err = s.roleRepo.UpdateWithTx(ctx, tx, role)
		}
		if err != nil {
			return err
		}

		// Assign permissions if provided
		if len(permissionIDs) > 0 {
			if err := s.roleRepo.AssignPermissionsWithTx(ctx, tx, role.ID, permissionIDs, operatorID); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteRole deletes a role
func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Deleting role",
		logger.String("operation", "DeleteRole"),
		logger.Uint("role_id", id))

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete role", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Role deleted successfully", logger.Uint("role_id", id))
	return nil
}

// ListRoles gets role list
func (s *roleService) ListRoles(ctx context.Context, req *request.ListRolesRequest) (*common.PaginationResponse, error) {
	roles, total, err := s.roleRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert response
	items := make([]response.RoleResponse, len(roles))
	for i, role := range roles {
		items[i] = *response.ConvertToRoleResponse(role)
	}

	return common.NewPaginationResponse(
		req.GetPage(),
		req.GetPageSize(),
		total,
		items,
	), nil
}

// GetRoleWithPermissions gets role with its permissions
func (s *roleService) GetRoleWithPermissions(ctx context.Context, id uint) (*response.RoleResponse, error) {
	role, err := s.roleRepo.GetWithPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.ConvertToRoleResponse(role), nil
}

// GetRoleUsers gets role users
func (s *roleService) GetRoleUsers(ctx context.Context, id uint, req *common.PaginationRequest) (*common.PaginationResponse, error) {
	users, total, err := s.roleRepo.GetUsers(ctx, id, req.GetOffset(), req.GetPageSize())
	if err != nil {
		return nil, err
	}

	// Convert to simple user response
	items := make([]response.RoleResponse, len(users))
	for i := range users {
		items[i] = response.RoleResponse{
			ID:   users[i].ID,
			Name: users[i].Username,
			Code: users[i].Email,
		}
	}

	return common.NewPaginationResponse(
		req.GetPage(),
		req.GetPageSize(),
		total,
		items,
	), nil
}

// AssignPermissions assigns permissions to a role
func (s *roleService) AssignPermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest, grantedBy uint) error {
	s.logger.InfoContext(ctx, "Assigning permissions to role",
		logger.String("operation", "AssignPermissions"),
		logger.Uint("role_id", roleID),
		logger.Uint("granted_by", grantedBy))

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if role exists
		role, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to get role for permission assignment", logger.ErrorField(err))
			return err
		}

		// Check if it's a system role
		if role.IsSystem {
			s.logger.WarnContext(ctx, "Attempt to modify system role permissions", logger.Uint("role_id", roleID))
			return errors.NewAppError(errors.CodeRecordUpdateFailed)
		}

		// Validate permission IDs
		permissions, err := s.permissionRepo.GetByIDs(ctx, req.PermissionIDs)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to validate permissions", logger.ErrorField(err))
			return err
		}
		if len(permissions) != len(req.PermissionIDs) {
			s.logger.WarnContext(ctx, "Some permission IDs are invalid during assignment",
				logger.Int("requested", len(req.PermissionIDs)),
				logger.Int("found", len(permissions)))
			return errors.NewAppError(errors.CodePermissionInvalid)
		}

		// Assign permissions
		if err := s.roleRepo.AssignPermissionsWithTx(ctx, tx, roleID, req.PermissionIDs, grantedBy); err != nil {
			s.logger.ErrorContext(ctx, "Failed to assign permissions", logger.ErrorField(err))
			return err
		}

		s.logger.InfoContext(ctx, "Permissions assigned successfully",
			logger.Uint("role_id", roleID),
			logger.Int("permission_count", len(req.PermissionIDs)))
		return nil
	})
}

// RemovePermissions removes permissions from a role
func (s *roleService) RemovePermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest) error {
	s.logger.InfoContext(ctx, "Removing permissions from role",
		logger.String("operation", "RemovePermissions"),
		logger.Uint("role_id", roleID))

	// Check if role exists
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get role for permission removal", logger.ErrorField(err))
		return err
	}

	// Check if it's a system role
	if role.IsSystem {
		s.logger.WarnContext(ctx, "Attempt to modify system role permissions", logger.Uint("role_id", roleID))
		return errors.NewAppError(errors.CodeRecordDeleteFailed)
	}

	// Remove permissions
	if err := s.roleRepo.RemovePermissions(ctx, roleID, req.PermissionIDs); err != nil {
		s.logger.ErrorContext(ctx, "Failed to remove permissions", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Permissions removed successfully",
		logger.Uint("role_id", roleID),
		logger.Int("permission_count", len(req.PermissionIDs)))
	return nil
}

// BatchUpdateRoleStatus updates status for multiple roles
func (s *roleService) BatchUpdateRoleStatus(ctx context.Context, ids []uint, status int) error {
	s.logger.InfoContext(ctx, "Batch updating role status",
		logger.String("operation", "BatchUpdateRoleStatus"),
		logger.Any("role_ids", ids),
		logger.Int("status", status))

	if err := s.roleRepo.BatchUpdateStatus(ctx, ids, status); err != nil {
		s.logger.ErrorContext(ctx, "Failed to batch update role status", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Role status updated successfully",
		logger.Int("updated_count", len(ids)))
	return nil
}
