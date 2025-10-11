package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"api-service/pkg/utils"

	"context"

	"gorm.io/gorm"
)

type permissionService struct {
	permissionRepo repository.PermissionRepository
	db             *gorm.DB
	logger         logger.Logger
}

// NewPermissionService creates a new permission service instance
func NewPermissionService(
	permissionRepo repository.PermissionRepository,
	db *gorm.DB,
	logger logger.Logger,
) service.PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		db:             db,
		logger:         logger,
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
	if err != nil && !errors.Is(err, errors.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing permission", logger.ErrorField(err))
		return nil, err
	}
	if existing != nil {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Get parent code if ParentID is provided
	var parentCode string
	if req.ParentID != nil {
		parent, err := s.permissionRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		parentCode = parent.Code
	}

	// Create permission entity
	permission := &model.Permission{
		ParentCode:  parentCode,
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
		Status:      1,
		CreatedBy:   &createdBy,
		UpdatedBy:   &createdBy,
	}

	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create permission", logger.ErrorField(err))
		return nil, err
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
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return nil, err
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
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return nil, err
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return nil, errors.NewAppError(errors.CodeRecordUpdateFailed)
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
		return nil, err
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
		s.logger.ErrorContext(ctx, "Failed to get permission", logger.ErrorField(err))
		return err
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return errors.NewAppError(errors.CodeRecordDeleteFailed)
	}

	// Check if permission has associated roles
	roleCount, err := s.permissionRepo.CountRoles(ctx, id)
	if err != nil {
		return err
	}
	if roleCount > 0 {
		return errors.NewAppError(errors.CodeResourceInUse)
	}

	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete permission", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Permission deleted successfully",
		logger.Uint("permission_id", id))

	return nil
}

// ListPermissions lists permissions with pagination
func (s *permissionService) ListPermissions(ctx context.Context, req *request.ListPermissionsRequest) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Listing permissions",
		logger.String("service", "permission"),
		logger.String("operation", "ListPermissions"))

	// Extract language from context for translation support
	lang := s.getLanguageFromContext(ctx)

	permissions, total, err := s.permissionRepo.List(ctx, req, lang)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list permissions", logger.ErrorField(err))
		return nil, err
	}

	// Convert to response format
	items := make([]response.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		items[i] = *response.ConvertToPermissionResponse(permission)
	}

	return common.NewPaginationResponse(
		req.GetPage(),
		req.GetPageSize(),
		total,
		items,
	), nil
}

// GetPermissionTree retrieves permissions in tree structure
func (s *permissionService) GetPermissionTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*response.PermissionTreeResponse, error) {
	s.logger.InfoContext(ctx, "Getting permission tree",
		logger.String("service", "permission"),
		logger.String("operation", "GetPermissionTree"))

	permissions, err := s.permissionRepo.GetTree(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get permission tree", logger.ErrorField(err))
		return nil, err
	}

	tree := response.ConvertToPermissionTreeResponse(permissions)
	return tree, nil
}

// GetPermissionRoles retrieves roles associated with a permission
func (s *permissionService) GetPermissionRoles(ctx context.Context, id uint, req *common.PaginationRequest) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Getting permission roles",
		logger.String("service", "permission"),
		logger.String("operation", "GetPermissionRoles"),
		logger.Uint("permission_id", id))

	roles, total, err := s.permissionRepo.GetRoles(ctx, id, req.GetOffset(), req.GetPageSize())
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get permission roles", logger.ErrorField(err))
		return nil, err
	}

	// Convert to response format
	items := make([]response.RoleResponse, len(roles))
	for i := range roles {
		items[i] = *response.ConvertToRoleResponse(&roles[i])
	}

	return common.NewPaginationResponse(
		req.GetPage(),
		req.GetPageSize(),
		total,
		items,
	), nil
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
		return false, err
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
		return nil, err
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
		return err
	}

	s.logger.InfoContext(ctx, "Permission status updated successfully",
		logger.Int("count", len(ids)))

	return nil
}

// getLanguageFromContext extracts language from context
func (s *permissionService) getLanguageFromContext(ctx context.Context) string {
	userID, exists := utils.GetUserIDFromContext(ctx)
	if exists {
		redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID)
		language, err := redis.HGet(ctx, redisKey, constants.UserLanguage)

		if err == nil && language != "" {
			return language
		}
	}
	return constants.DefaultLanguage
}
