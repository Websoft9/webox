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

const (
	RoleManagementSortOrder       = 2
	PermissionManagementSortOrder = 3
)

type permissionService struct {
	permissionRepo repository.PermissionRepository
	db             *gorm.DB
	logger         logger.Logger
	i18n           *i18n.I18n
}

// NewPermissionService creates a new permission service instance
func NewPermissionService(
	permissionRepo repository.PermissionRepository,
	db *gorm.DB,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		db:             db,
		logger:         logger,
		i18n:           i18n,
	}
}

// CreatePermission creates a new permission
func (s *permissionService) CreatePermission(ctx context.Context, req *request.CreatePermissionRequest, createdBy uint) (*response.PermissionResponse, error) {
	s.logger.InfoContext(ctx, "Creating permission",
		logger.String("service", "permission"),
		logger.String("operation", "CreatePermission"),
		logger.String("code", req.Code))

	// Check if permission code already exists
	existing, err := s.permissionRepo.GetByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing permission", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to check existing permission")
	}
	if existing != nil {
		return nil, errors.New("permission code already exists")
	}

	// Create permission entity
	permission := &model.Permission{
		ParentID:    req.ParentID,
		Scope:       req.Scope,
		Name:        req.Name,
		Code:        req.Code,
		Module:      req.Module,
		Action:      req.Action,
		Resource:    req.Resource,
		Description: req.Description,
		IsSystem:    false,
		IsMenu:      req.IsMenu,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		CreatedBy:   &createdBy,
		UpdatedBy:   &createdBy,
	}

	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create permission", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to create permission")
	}

	s.logger.InfoContext(ctx, "Permission created successfully",
		logger.Uint("permission_id", permission.ID))

	return response.ConvertToPermissionResponse(permission), nil
}

// GetPermission gets permission by ID
func (s *permissionService) GetPermission(ctx context.Context, id uint) (*response.PermissionResponse, error) {
	s.logger.InfoContext(ctx, "Getting permission",
		logger.String("service", "permission"),
		logger.String("operation", "GetPermission"),
		logger.Uint("permission_id", id))

	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get permission")
	}

	return response.ConvertToPermissionResponse(permission), nil
}

// UpdatePermission updates permission
func (s *permissionService) UpdatePermission(ctx context.Context, id uint, req *request.UpdatePermissionRequest, updatedBy uint) (*response.PermissionResponse, error) {
	s.logger.InfoContext(ctx, "Updating permission",
		logger.String("service", "permission"),
		logger.String("operation", "UpdatePermission"),
		logger.Uint("permission_id", id))

	// Get existing permission
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get permission")
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return nil, errors.New("cannot update system permission")
	}

	// Update fields
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	permission.SortOrder = req.SortOrder
	permission.Status = req.Status
	permission.UpdatedBy = &updatedBy

	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update permission", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to update permission")
	}

	s.logger.InfoContext(ctx, "Permission updated successfully",
		logger.Uint("permission_id", permission.ID))

	return response.ConvertToPermissionResponse(permission), nil
}

// DeletePermission deletes permission
func (s *permissionService) DeletePermission(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Deleting permission",
		logger.String("service", "permission"),
		logger.String("operation", "DeletePermission"),
		logger.Uint("permission_id", id))

	// Get existing permission
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return errors.Wrap(err, "failed to get permission")
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return errors.New("cannot delete system permission")
	}

	// Check if permission has associated roles
	roleCount, err := s.permissionRepo.CountRoles(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to count permission roles", logger.ErrorField(err))
		return errors.Wrap(err, "failed to count permission roles")
	}
	if roleCount > 0 {
		return errors.New("permission has associated roles and cannot be deleted")
	}

	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete permission", logger.ErrorField(err))
		return errors.Wrap(err, "failed to delete permission")
	}

	s.logger.InfoContext(ctx, "Permission deleted successfully",
		logger.Uint("permission_id", id))

	return nil
}

// ListPermissions lists permissions with pagination
func (s *permissionService) ListPermissions(ctx context.Context, req *request.ListPermissionsRequest) (*response.PermissionListResponse, error) {
	s.logger.InfoContext(ctx, "Listing permissions",
		logger.String("service", "permission"),
		logger.String("operation", "ListPermissions"))

	permissions, total, err := s.permissionRepo.List(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list permissions", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to list permissions")
	}

	// Convert to response format
	items := make([]response.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		items[i] = *response.ConvertToPermissionResponse(permission)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.GetPageSize())))

	return &response.PermissionListResponse{
		Items:      items,
		Total:      total,
		Page:       req.GetPage(),
		PageSize:   req.GetPageSize(),
		TotalPages: totalPages,
	}, nil
}

// GetPermissionTree retrieves permissions in tree structure
func (s *permissionService) GetPermissionTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*response.PermissionTreeResponse, error) {
	s.logger.InfoContext(ctx, "Getting permission tree",
		logger.String("service", "permission"),
		logger.String("operation", "GetPermissionTree"))

	permissions, err := s.permissionRepo.GetTree(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get permission tree", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get permission tree")
	}

	tree := response.ConvertToPermissionTreeResponse(permissions)
	return tree, nil
}

// GetPermissionRoles retrieves roles associated with a permission
func (s *permissionService) GetPermissionRoles(ctx context.Context, id uint, page, pageSize int) (*response.RoleListResponse, error) {
	s.logger.InfoContext(ctx, "Getting permission roles",
		logger.String("service", "permission"),
		logger.String("operation", "GetPermissionRoles"),
		logger.Uint("permission_id", id))

	roles, total, err := s.permissionRepo.GetRoles(ctx, id, page, pageSize)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get permission roles", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get permission roles")
	}

	// Convert to response format
	items := make([]response.RoleResponse, len(roles))
	for i := range roles {
		items[i] = *response.ConvertToRoleResponse(&roles[i])
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &response.RoleListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// CheckUserPermission checks if user has specific permission
func (s *permissionService) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	s.logger.InfoContext(ctx, "Checking user permission",
		logger.String("service", "permission"),
		logger.String("operation", "CheckUserPermission"),
		logger.Uint("user_id", userID),
		logger.String("resource", resource),
		logger.String("action", action))

	hasPermission, err := s.permissionRepo.CheckUserPermission(ctx, userID, resource, action)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check user permission", logger.ErrorField(err))
		return false, errors.Wrap(err, "failed to check user permission")
	}

	return hasPermission, nil
}

// GetUserPermissions gets all permissions for a user
func (s *permissionService) GetUserPermissions(ctx context.Context, userID uint) ([]*response.PermissionResponse, error) {
	s.logger.InfoContext(ctx, "Getting user permissions",
		logger.String("service", "permission"),
		logger.String("operation", "GetUserPermissions"),
		logger.Uint("user_id", userID))

	permissions, err := s.permissionRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user permissions", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get user permissions")
	}

	// Convert to response format
	result := make([]*response.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		result[i] = response.ConvertToPermissionResponse(permission)
	}

	return result, nil
}

// BatchUpdatePermissionStatus updates status for multiple permissions
func (s *permissionService) BatchUpdatePermissionStatus(ctx context.Context, ids []uint, status int) error {
	s.logger.InfoContext(ctx, "Batch updating permission status",
		logger.String("service", "permission"),
		logger.String("operation", "BatchUpdatePermissionStatus"),
		logger.Int("count", len(ids)),
		logger.Int("status", status))

	if err := s.permissionRepo.BatchUpdateStatus(ctx, ids, status); err != nil {
		s.logger.ErrorContext(ctx, "Failed to batch update permission status", logger.ErrorField(err))
		return errors.Wrap(err, "failed to batch update permission status")
	}

	s.logger.InfoContext(ctx, "Permission status updated successfully",
		logger.Int("count", len(ids)))

	return nil
}

// InitializeSystemPermissions initializes system permissions
func (s *permissionService) InitializeSystemPermissions(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Initializing system permissions",
		logger.String("service", "permission"),
		logger.String("operation", "InitializeSystemPermissions"))

	// Define system permissions
	systemPermissions := []*model.Permission{
		{
			Scope:       "platform",
			Name:        "User Management",
			Code:        "user.manage",
			Module:      "user",
			Action:      "manage",
			Resource:    "user",
			Description: "Manage users",
			IsSystem:    true,
			IsMenu:      true,
			SortOrder:   1,
			Status:      1,
		},
		{
			Scope:       "platform",
			Name:        "Role Management",
			Code:        "role.manage",
			Module:      "role",
			Action:      "manage",
			Resource:    "role",
			Description: "Manage roles",
			IsSystem:    true,
			IsMenu:      true,
			SortOrder:   RoleManagementSortOrder,
			Status:      1,
		},
		{
			Scope:       "platform",
			Name:        "Permission Management",
			Code:        "permission.manage",
			Module:      "permission",
			Action:      "manage",
			Resource:    "permission",
			Description: "Manage permissions",
			IsSystem:    true,
			IsMenu:      true,
			SortOrder:   PermissionManagementSortOrder,
			Status:      1,
		},
	}

	// Create system permissions if they don't exist
	for _, perm := range systemPermissions {
		existing, err := s.permissionRepo.GetByCode(ctx, perm.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.ErrorContext(ctx, "Failed to check existing permission", logger.ErrorField(err))
			return errors.Wrap(err, "failed to check existing permission")
		}

		if existing == nil {
			if err := s.permissionRepo.Create(ctx, perm); err != nil {
				s.logger.ErrorContext(ctx, "Failed to create system permission", logger.ErrorField(err))
				return errors.Wrap(err, "failed to create system permission")
			}
			s.logger.InfoContext(ctx, "System permission created",
				logger.String("code", perm.Code))
		}
	}

	s.logger.InfoContext(ctx, "System permissions initialized successfully")
	return nil
}
