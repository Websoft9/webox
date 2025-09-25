package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"context"
	"strings"
	"time"

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
		return errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to create permission")
	}
	return nil
}

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("status != -1").First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission by ID")
	}
	return &permission, nil
}

// GetByCode retrieves a permission by code
func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ? AND status != -1", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission by code")
	}
	return &permission, nil
}

// Update updates a permission
func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	// Check if permission exists and is not deleted or disabled
	var existingPermission model.Permission
	if err := r.db.WithContext(ctx).Where("status != -1").First(&existingPermission, permission.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission for update")
	}

	// Check if permission is disabled
	if existingPermission.Status == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordIsDisabled, "cannot update disabled permission")
	}

	result := r.db.WithContext(ctx).Model(permission).
		Omit("created_at", "code", "scope", "module", "action", "resource").
		Updates(permission)

	if result.Error != nil {
		return errors.WrapError(result.Error, errors.CodeRecordUpdateFailed, "failed to update permission")
	}

	if result.RowsAffected == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNoAffected, "no rows affected")
	}

	return nil
}

// Delete deletes a permission (soft delete)
func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	// Check if permission exists and is not already deleted
	var permission model.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission")
	}

	// Check if already deleted
	if permission.Status == -1 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordDeleteFailed, "permission already deleted")
	}

	// Check if there are associated roles
	var roleCount int64
	if err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("permission_code = ?", id).Count(&roleCount).Error; err != nil {
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to check permission roles")
	}

	if roleCount > 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordDeleteFailed, "cannot delete permission with associated roles")
	}

	if permission.IsSystem {
		return errors.NewAppErrorWithMessage(errors.CodeRecordDeleteDenied, "cannot delete system permission")
	}

	// Check if there are child permissions that are not deleted
	var childCount int64
	if err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("parent_code = ? AND status != -1", id).Count(&childCount).Error; err != nil {
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to check child permissions")
	}

	if childCount > 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordDeleteFailed, "cannot delete permission with active child permissions")
	}

	// Perform soft delete by setting status to -1
	result := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("id = ?", id).
		Update("status", -1)

	if result.Error != nil {
		return errors.WrapError(result.Error, errors.CodeRecordDeleteFailed, "failed to delete permission")
	}

	if result.RowsAffected == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNoAffected, "no rows affected")
	}

	return nil
}

// List retrieves a list of permissions
func (r *permissionRepository) List(ctx context.Context, req *request.ListPermissionsRequest, lang string) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Permission{}).Where("status != -1")

	// Build query conditions with translation support
	if req.Search != "" {
		// Find permission codes that match the translated names
		matchedCodes, err := r.findMatchingPermissionCodes(ctx, req.Search, lang)
		if err != nil {
			// If translation matching fails, fallback to original search
			query = query.Where("code LIKE ? OR description LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
		} else if len(matchedCodes) > 0 {
			// Use matched codes from translation, also include code and description search
			query = query.Where("code IN (?) OR code LIKE ? OR description LIKE ?", matchedCodes, "%"+req.Search+"%", "%"+req.Search+"%")
		} else {
			// No translation matches found, fallback to original search
			query = query.Where("code LIKE ? OR description LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
		}
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
		return nil, 0, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to count permissions")
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order("module ASC, sort_order ASC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&permissions).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to list permissions")
	}

	// Get statistical information
	for _, permission := range permissions {
		// Get role count
		var roleCount int64
		r.db.WithContext(ctx).Model(&model.RolePermission{}).
			Where("permission_code = ? AND status = 1", permission.ID).Count(&roleCount)
		permission.RoleCount = roleCount
	}

	return permissions, total, nil
}

// GetTree retrieves the permission tree
func (r *permissionRepository) GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error) {
	query := r.buildBaseQuery(ctx, req)

	var permissions []*model.Permission
	var err error

	if req.Scope != "" {
		permissions, err = r.getPermissionsWithScope(ctx, query, req.Scope)
	} else {
		permissions, err = r.getAllPermissions(query)
	}

	if err != nil {
		return nil, err
	}

	tree := r.buildPermissionTree(permissions)

	if req.Scope != "" {
		return r.filterTreeByScope(tree, req.Scope), nil
	}

	return tree, nil
}

// buildBaseQuery creates the base query with common conditions
func (r *permissionRepository) buildBaseQuery(ctx context.Context, req *request.PermissionTreeRequest) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&model.Permission{}).Where("status != -1")

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		query = query.Where("status = 1") // Show only enabled permissions by default
	}

	return query
}

// getAllPermissions retrieves all permissions without scope filtering
func (r *permissionRepository) getAllPermissions(query *gorm.DB) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := query.Order("module ASC, sort_order ASC, created_at ASC").Find(&permissions).Error
	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission tree")
	}
	return permissions, nil
}

// getPermissionsWithScope retrieves permissions with scope filtering and their parents
func (r *permissionRepository) getPermissionsWithScope(ctx context.Context, baseQuery *gorm.DB, scope string) ([]*model.Permission, error) {
	scopePermissions, err := r.getScopedPermissions(baseQuery, scope)
	if err != nil {
		return nil, err
	}

	if len(scopePermissions) == 0 {
		return []*model.Permission{}, nil
	}

	parentCodes := r.collectParentCodes(scopePermissions)
	if len(parentCodes) == 0 {
		return scopePermissions, nil
	}

	parentPermissions, err := r.getParentPermissions(ctx, baseQuery, parentCodes)
	if err != nil {
		return nil, err
	}

	return r.combinePermissions(scopePermissions, parentPermissions), nil
}

// getScopedPermissions retrieves permissions for a specific scope
func (r *permissionRepository) getScopedPermissions(query *gorm.DB, scope string) ([]*model.Permission, error) {
	var scopePermissions []*model.Permission
	scopeQuery := query.Where("scope = ?", scope)
	err := scopeQuery.Order("module ASC, sort_order ASC, created_at ASC").Find(&scopePermissions).Error
	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get scoped permissions")
	}
	return scopePermissions, nil
}

// collectParentCodes extracts parent codes from permissions
func (r *permissionRepository) collectParentCodes(permissions []*model.Permission) []string {
	parentCodes := make(map[string]bool)
	for _, perm := range permissions {
		if perm.ParentCode != "" {
			parentCodes[perm.ParentCode] = true
		}
	}

	if len(parentCodes) == 0 {
		return nil
	}

	parentCodeSlice := make([]string, 0, len(parentCodes))
	for code := range parentCodes {
		parentCodeSlice = append(parentCodeSlice, code)
	}

	return parentCodeSlice
}

// getParentPermissions retrieves parent permissions by codes
func (r *permissionRepository) getParentPermissions(ctx context.Context, baseQuery *gorm.DB, parentCodes []string) ([]*model.Permission, error) {
	var parentPermissions []*model.Permission
	parentQuery := r.db.WithContext(ctx).Model(&model.Permission{}).
		Where("code IN ? AND status != -1", parentCodes)

	// Apply the same status filter as base query
	if baseQuery != nil {
		// Extract status condition from base query if needed
		parentQuery = parentQuery.Where("status = 1") // Default to enabled
	}

	err := parentQuery.Order("module ASC, sort_order ASC, created_at ASC").Find(&parentPermissions).Error
	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get parent permissions")
	}

	return parentPermissions, nil
}

// combinePermissions merges scoped permissions with their parents
func (r *permissionRepository) combinePermissions(scopePermissions, parentPermissions []*model.Permission) []*model.Permission {
	permissionMap := make(map[string]*model.Permission)

	// Add scoped permissions first
	for _, perm := range scopePermissions {
		permissionMap[perm.Code] = perm
	}

	// Add parent permissions if not already present
	for _, perm := range parentPermissions {
		if _, exists := permissionMap[perm.Code]; !exists {
			permissionMap[perm.Code] = perm
		}
	}

	// Convert map back to slice
	permissions := make([]*model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions
}

// buildPermissionTree builds the permission tree
func (r *permissionRepository) buildPermissionTree(permissions []*model.Permission) []*model.Permission {
	// Initialize Children slice for all permissions
	for _, perm := range permissions {
		perm.Children = []*model.Permission{}
	}

	// Create code to permission mapping
	permissionMap := make(map[string]*model.Permission)
	for _, perm := range permissions {
		permissionMap[perm.Code] = perm
	}

	// Build tree structure
	var roots []*model.Permission
	for _, perm := range permissions {
		if perm.ParentCode == "" {
			// Root node
			roots = append(roots, perm)
		} else {
			// Find parent and add this as child
			if parent, exists := permissionMap[perm.ParentCode]; exists {
				parent.Children = append(parent.Children, perm)
			}
		}
	}

	return roots
}

// filterTreeByScope filters the tree to show only trees that contain permissions of the specified scope
func (r *permissionRepository) filterTreeByScope(roots []*model.Permission, scope string) []*model.Permission {
	var filteredRoots []*model.Permission

	for _, root := range roots {
		if filteredRoot := r.filterNodeByScope(root, scope); filteredRoot != nil {
			filteredRoots = append(filteredRoots, filteredRoot)
		}
	}

	return filteredRoots
}

// filterNodeByScope recursively filters a node and its children by scope
func (r *permissionRepository) filterNodeByScope(node *model.Permission, scope string) *model.Permission {
	// Create a copy of the node
	filteredNode := *node
	filteredNode.Children = []*model.Permission{}

	// Check if this node has the target scope
	nodeHasTargetScope := node.Scope == scope

	// Check if any descendants have the target scope
	hasTargetScopeInChildren := false

	// Recursively filter children
	for _, child := range node.Children {
		if filteredChild := r.filterNodeByScope(child, scope); filteredChild != nil {
			filteredNode.Children = append(filteredNode.Children, filteredChild)
			hasTargetScopeInChildren = true
		}
	}

	// Return the node if:
	// 1. The node itself has the target scope, OR
	// 2. Any of its descendants have the target scope
	if nodeHasTargetScope || hasTargetScopeInChildren {
		return &filteredNode
	}

	return nil
}

// GetWithRoles retrieves a permission with its roles
func (r *permissionRepository) GetWithRoles(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).
		Preload("Roles", "status != -1").
		First(&permission, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission with roles")
	}

	return &permission, nil
}

// GetChildren retrieves child permissions
func (r *permissionRepository) GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Where("parent_code = ? AND status != -1", parentID).
		Order("sort_order ASC, created_at ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get child permissions")
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
		Where("permission_code = ? AND status = 1", permissionID).Count(&count).Error

	if err != nil {
		return 0, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to count permission roles")
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
		return errors.WrapError(result.Error, errors.CodeRecordUpdateFailed, "failed to batch update permission status")
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
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permissions by IDs")
	}

	return permissions, nil
}

// GetUserPermissions retrieves permissions for a user
func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON permissions.code = role_permissions.permission_code").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.status != -1 AND role_permissions.status != -1 AND permissions.status != -1", userID).
		Find(&permissions).Error

	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get user permissions")
	}

	return permissions, nil
}

// CheckUserPermission checks if a user has a specific permission
func (r *permissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.code = role_permissions.permission_code").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Where("user_roles.status = 1 AND role_permissions.status != -1 AND permissions.status != -1").
		Count(&count).Error

	if err != nil {
		return false, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to check user permission")
	}

	return count > 0, nil
}

// CreateWithTx creates a permission within a transaction
func (r *permissionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	if err := tx.WithContext(ctx).Create(permission).Error; err != nil {
		return errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to create permission in transaction")
	}
	return nil
}

// UpdateWithTx updates a permission within a transaction
func (r *permissionRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	// Check if permission exists and is not deleted or disabled
	var existingPermission model.Permission
	if err := tx.WithContext(ctx).Where("status != -1").First(&existingPermission, permission.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission for update in transaction")
	}

	// Check if permission is disabled
	if existingPermission.Status == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordIsDisabled, "cannot update disabled permission")
	}

	result := tx.WithContext(ctx).Model(permission).
		Omit("created_at", "code", "scope", "module", "action", "resource").
		Updates(permission)

	if result.Error != nil {
		return errors.WrapError(result.Error, errors.CodeRecordUpdateFailed, "failed to update permission in transaction")
	}

	if result.RowsAffected == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNoAffected, "no rows affected in transaction")
	}

	return nil
}

// findMatchingPermissionCodes finds permission codes that match the search term after translation
func (r *permissionRepository) findMatchingPermissionCodes(ctx context.Context, searchTerm, lang string) ([]string, error) {
	// Create a struct to hold only the fields we need
	type PermissionForTranslation struct {
		Name string `gorm:"column:name"`
		Code string `gorm:"column:code"`
	}

	var permissions []PermissionForTranslation
	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Select("name, code").
		Where("status != -1").
		Order("id").
		Find(&permissions).Error

	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permissions for translation matching")
	}

	var matchedCodes []string
	searchTermLower := strings.ToLower(searchTerm)

	for _, perm := range permissions {
		// Translate the permission name using i18n
		translatedName := i18n.T(perm.Name, lang)
		translatedNameLower := strings.ToLower(translatedName)

		// Check if translated name contains the search term
		if strings.Contains(translatedNameLower, searchTermLower) {
			matchedCodes = append(matchedCodes, perm.Code)
		}
	}

	return matchedCodes, nil
}
