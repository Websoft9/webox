package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"math"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	db             *gorm.DB
	logger         logger.Logger
	i18nInstance   *i18n.I18n
}

// NewRoleService creates a new role service instance
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	db *gorm.DB,
	logger logger.Logger,
	i18nInstance *i18n.I18n,
) service.RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		db:             db,
		logger:         logger,
		i18nInstance:   i18nInstance,
	}
}

// CreateRole creates a new role
func (s *roleService) CreateRole(ctx context.Context, req *request.CreateRoleRequest, createdBy uint) (*response.RoleResponse, error) {
	s.logger.InfoContext(ctx, "Creating role",
		logger.String("operation", "CreateRole"),
		logger.String("role_code", req.Code),
		logger.Uint("created_by", createdBy))

	// Check if role code already exists
	existingRole, err := s.roleRepo.GetByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing role", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to check existing role")
	}
	if existingRole != nil {
		s.logger.WarnContext(ctx, "Role code already exists", logger.String("role_code", req.Code))
		return nil, errors.New("role code already exists")
	}

	// Validate permission IDs if provided
	if len(req.PermissionIDs) > 0 {
		permissions, err := s.permissionRepo.GetByIDs(ctx, req.PermissionIDs)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to validate permissions", logger.ErrorField(err))
			return nil, errors.Wrap(err, "failed to validate permissions")
		}
		if len(permissions) != len(req.PermissionIDs) {
			s.logger.WarnContext(ctx, "Some permission IDs are invalid",
				logger.Int("requested", len(req.PermissionIDs)),
				logger.Int("found", len(permissions)))
			return nil, errors.New("some permission IDs are invalid")
		}
	}

	// Create role entity
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsSystem:    false,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		CreatedBy:   &createdBy,
		UpdatedBy:   &createdBy,
	}

	// Use transaction to create role and assign permissions
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create role
		if err := s.roleRepo.CreateWithTx(ctx, tx, role); err != nil {
			return err
		}

		// Assign permissions if provided
		if len(req.PermissionIDs) > 0 {
			if err := s.roleRepo.AssignPermissions(ctx, role.ID, req.PermissionIDs, createdBy); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create role", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to create role")
	}

	s.logger.InfoContext(ctx, "Role created successfully", logger.Uint("role_id", role.ID))

	// Return complete role information
	return s.GetRole(ctx, role.ID)
}

// GetRole 获取角色
func (s *roleService) GetRole(ctx context.Context, id uint) (*response.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取统计信息
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

	// Get existing role
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get role for update", logger.ErrorField(err))
		return nil, err
	}

	// Check if it's a system role
	if role.IsSystem {
		s.logger.WarnContext(ctx, "Attempt to update system role", logger.Uint("role_id", id))
		return nil, errors.New("cannot update system role")
	}

	// Validate permission IDs if provided
	if len(req.PermissionIDs) > 0 {
		permissions, err := s.permissionRepo.GetByIDs(ctx, req.PermissionIDs)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to validate permissions", logger.ErrorField(err))
			return nil, errors.Wrap(err, "failed to validate permissions")
		}
		if len(permissions) != len(req.PermissionIDs) {
			s.logger.WarnContext(ctx, "Some permission IDs are invalid during update",
				logger.Int("requested", len(req.PermissionIDs)),
				logger.Int("found", len(permissions)))
			return nil, errors.New("some permission IDs are invalid")
		}
	}

	// Update role information
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	role.SortOrder = req.SortOrder
	role.Status = req.Status
	role.UpdatedBy = &updatedBy

	// Use transaction to update role and permissions
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update role
		if err := s.roleRepo.UpdateWithTx(ctx, tx, role); err != nil {
			return err
		}

		// Update permission assignments if provided
		if len(req.PermissionIDs) > 0 {
			if err := s.roleRepo.AssignPermissions(ctx, role.ID, req.PermissionIDs, updatedBy); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update role", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to update role")
	}

	s.logger.InfoContext(ctx, "Role updated successfully", logger.Uint("role_id", id))

	return s.GetRole(ctx, id)
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

// ListRoles 获取角色列表
func (s *roleService) ListRoles(ctx context.Context, req *request.ListRolesRequest) (*response.RoleListResponse, error) {
	roles, total, err := s.roleRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换响应
	items := make([]response.RoleResponse, len(roles))
	for i, role := range roles {
		items[i] = *response.ConvertToRoleResponse(role)
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.GetPageSize())))

	return &response.RoleListResponse{
		Items:      items,
		Total:      total,
		Page:       req.GetPage(),
		PageSize:   req.GetPageSize(),
		TotalPages: totalPages,
	}, nil
}

// GetRoleWithPermissions 获取角色及其权限
func (s *roleService) GetRoleWithPermissions(ctx context.Context, id uint) (*response.RoleResponse, error) {
	role, err := s.roleRepo.GetWithPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.ConvertToRoleResponse(role), nil
}

// GetRoleUsers 获取角色用户
func (s *roleService) GetRoleUsers(ctx context.Context, id uint, page, pageSize int) (*response.RoleListResponse, error) {
	users, total, err := s.roleRepo.GetUsers(ctx, id, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为用户简单响应
	items := make([]response.RoleResponse, len(users))
	for i, user := range users {
		items[i] = response.RoleResponse{
			ID:   user.ID,
			Name: user.Username,
			Code: user.Email,
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &response.RoleListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// AssignPermissions assigns permissions to a role
func (s *roleService) AssignPermissions(ctx context.Context, roleID uint, req *request.RolePermissionRequest, grantedBy uint) error {
	s.logger.InfoContext(ctx, "Assigning permissions to role",
		logger.String("operation", "AssignPermissions"),
		logger.Uint("role_id", roleID),
		logger.Uint("granted_by", grantedBy))

	// Check if role exists
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get role for permission assignment", logger.ErrorField(err))
		return err
	}

	// Check if it's a system role
	if role.IsSystem {
		s.logger.WarnContext(ctx, "Attempt to modify system role permissions", logger.Uint("role_id", roleID))
		return errors.New("cannot modify system role permissions")
	}

	// Validate permission IDs
	permissions, err := s.permissionRepo.GetByIDs(ctx, req.PermissionIDs)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to validate permissions", logger.ErrorField(err))
		return errors.Wrap(err, "failed to validate permissions")
	}
	if len(permissions) != len(req.PermissionIDs) {
		s.logger.WarnContext(ctx, "Some permission IDs are invalid during assignment",
			logger.Int("requested", len(req.PermissionIDs)),
			logger.Int("found", len(permissions)))
		return errors.New("some permission IDs are invalid")
	}

	// Assign permissions
	if err := s.roleRepo.AssignPermissions(ctx, roleID, req.PermissionIDs, grantedBy); err != nil {
		s.logger.ErrorContext(ctx, "Failed to assign permissions", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Permissions assigned successfully",
		logger.Uint("role_id", roleID),
		logger.Int("permission_count", len(req.PermissionIDs)))
	return nil
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
		return errors.New("cannot modify system role permissions")
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

// InitializeSystemRoles initializes system roles
func (s *roleService) InitializeSystemRoles(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Initializing system roles",
		logger.String("operation", "InitializeSystemRoles"))

	systemRoles := []model.Role{
		{
			Name:        "Administrator",
			Code:        "admin",
			Description: "System administrator with full access to all features",
			IsSystem:    true,
			SortOrder:   1,
			Status:      1,
		},
		{
			Name:        "Regular User",
			Code:        "user",
			Description: "Regular user with basic application deployment and resource management permissions",
			IsSystem:    true,
			SortOrder:   2,
			Status:      1,
		},
		{
			Name:        "Operator",
			Code:        "operator",
			Description: "Operations personnel with server management and monitoring permissions",
			IsSystem:    true,
			SortOrder:   3,
			Status:      1,
		},
		{
			Name:        "Developer",
			Code:        "developer",
			Description: "Developer with application development and deployment permissions",
			IsSystem:    true,
			SortOrder:   4,
			Status:      1,
		},
	}

	for _, role := range systemRoles {
		// Check if role already exists
		existingRole, err := s.roleRepo.GetByCode(ctx, role.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.ErrorContext(ctx, "Failed to check existing role",
				logger.ErrorField(err),
				logger.String("role_code", role.Code))
			continue
		}

		if existingRole == nil {
			// Create new role
			if err := s.roleRepo.Create(ctx, &role); err != nil {
				s.logger.ErrorContext(ctx, "Failed to create system role",
					logger.ErrorField(err),
					logger.String("role_code", role.Code))
				continue
			}
			s.logger.InfoContext(ctx, "System role created",
				logger.String("role_code", role.Code))
		}
	}

	s.logger.InfoContext(ctx, "System roles initialized successfully")
	return nil
}
