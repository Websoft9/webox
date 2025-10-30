package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

const (
	resourceGroupCodePrefix = "rg" // Prefix for resource group code
	resourceGroupCodeLength = 10   // Length of random string in resource group code
)

// resourceGroupRepository implements ResourceGroupRepository
type resourceGroupRepository struct {
	db *gorm.DB
}

// NewResourceGroupRepository creates a new resource group repository
func NewResourceGroupRepository(db *gorm.DB) repository.ResourceGroupRepository {
	return &resourceGroupRepository{
		db: db,
	}
}

// Create creates a new resource group and generates its code
func (r *resourceGroupRepository) Create(ctx context.Context, rg *model.ResourceGroup) error {
	// Generate code before creating
	rg.Code = model.GenerateCode(resourceGroupCodePrefix, resourceGroupCodeLength)

	// Create the resource group
	if err := r.db.WithContext(ctx).Create(rg).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByID retrieves a resource group by ID
func (r *resourceGroupRepository) GetByID(ctx context.Context, id uint) (*model.ResourceGroup, error) {
	var rg model.ResourceGroup
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &rg, nil
}

// GetByCode retrieves a resource group by code
func (r *resourceGroupRepository) GetByCode(ctx context.Context, code string) (*model.ResourceGroup, error) {
	var rg model.ResourceGroup
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &rg, nil
}

// GetList retrieves a paginated list of resource groups with filters
func (r *resourceGroupRepository) GetList(ctx context.Context, req *request.GetResourceGroupListRequest, ownerID uint) ([]*model.ResourceGroup, int64, error) {
	var resourceGroups []*model.ResourceGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ResourceGroup{})

	// Filter by owner
	query = query.Where("owner_id = ?", ownerID)

	// Apply filters
	if req.ProjectID != nil {
		query = query.Where("project_id = ?", *req.ProjectID)
	}

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination and sorting
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	orderBy := req.GetSortOrder()

	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&resourceGroups).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return resourceGroups, total, nil
}

// Update updates an existing resource group
func (r *resourceGroupRepository) Update(ctx context.Context, rg *model.ResourceGroup) error {
	if err := r.db.WithContext(ctx).Model(rg).Updates(rg).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// Delete deletes a resource group by ID
func (r *resourceGroupRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.ResourceGroup{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	return nil
}

// CheckNameExists checks if a resource group name exists in a project
func (r *resourceGroupRepository) CheckNameExists(ctx context.Context, projectID uint, name string, excludeID *uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&model.ResourceGroup{}).
		Where("project_id = ? AND name = ?", projectID, name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// GetDefaultResourceGroup retrieves the default resource group for a project
func (r *resourceGroupRepository) GetDefaultResourceGroup(ctx context.Context, projectID uint) (*model.ResourceGroup, error) {
	var rg model.ResourceGroup
	if err := r.db.WithContext(ctx).Where("project_id = ? AND is_default = ?", projectID, true).First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &rg, nil
}

// GetProjectIDFromResourceCode retrieves the project_id from a resource code
// It parses the code to identify resource type, queries the resource to get resource_group_id,
// then queries the resource group to get project_id
func (r *resourceGroupRepository) GetProjectIDFromResourceCode(ctx context.Context, resourceCode string) (uint, error) {
	// Get resource type mapping from resource_types table
	var resourceTypes []struct {
		Code      string
		TableName string
	}
	if err := r.db.WithContext(ctx).Table("resource_types").Select("code, table_name").Find(&resourceTypes).Error; err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Build a map from resource type prefix to table name
	typeToTableMap := make(map[string]string)
	for _, rt := range resourceTypes {
		typeToTableMap[rt.Code] = rt.TableName
	}

	// Extract resource type prefix from code (e.g., "database_xxx" -> "database")
	// Convention: code format is "{resource_type}_{identifier}"
	var tableName string
	for prefix := range typeToTableMap {
		if len(resourceCode) > len(prefix) && resourceCode[:len(prefix)] == prefix && resourceCode[len(prefix)] == '_' {
			tableName = typeToTableMap[prefix]
			break
		}
	}

	if tableName == "" {
		return 0, errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Query the resource table to get resource_group_id
	var result struct {
		ResourceGroupID *uint
	}

	query := `SELECT resource_group_id FROM ` + tableName + ` WHERE code = ? LIMIT 1`
	if err := r.db.WithContext(ctx).Raw(query, resourceCode).Scan(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// If resource has no resource_group_id, it should belong to a default group
	// but we can't determine the project. This is an invalid state.
	if result.ResourceGroupID == nil {
		return 0, errors.NewAppError(errors.CodeRecordNotFound)
	}

	// Query resource_groups table to get project_id
	var rg model.ResourceGroup
	if err := r.db.WithContext(ctx).First(&rg, *result.ResourceGroupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return rg.ProjectID, nil
}

// GetResourcesByGroupID retrieves resources in a resource group
func (r *resourceGroupRepository) GetResourcesByGroupID(
	ctx context.Context,
	groupID uint,
	req *request.GetResourceGroupResourcesRequest,
) ([]response.ResourceItemResponse, int64, error) {
	var resources []response.ResourceItemResponse
	var total int64

	// Get resource type mapping from resource_types table
	var resourceTypes []struct {
		Code      string
		TableName string
	}

	query := r.db.WithContext(ctx).Table("resource_types").Select("code, table_name")

	// If resource type is specified, only query that type
	if req.ResourceType != nil && *req.ResourceType != "" {
		query = query.Where("code = ?", *req.ResourceType)
	}

	if err := query.Find(&resourceTypes).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Query resources for each resource type
	for _, rt := range resourceTypes {
		typeResources, _, err := r.queryResourcesByType(ctx, groupID, rt.Code, rt.TableName)
		if err != nil {
			return nil, 0, err
		}
		resources = append(resources, typeResources...)
	}

	total = int64(len(resources))

	// Apply pagination
	offset := req.GetOffset()
	pageSize := req.GetPageSize()

	if offset >= len(resources) {
		return []response.ResourceItemResponse{}, total, nil
	}

	end := offset + pageSize
	if end > len(resources) {
		end = len(resources)
	}

	return resources[offset:end], total, nil
}

// queryResourcesByType queries resources of a specific type
func (r *resourceGroupRepository) queryResourcesByType(
	ctx context.Context,
	groupID uint,
	resourceType, tableName string,
) ([]response.ResourceItemResponse, int64, error) {
	var resources []response.ResourceItemResponse

	query := `
		SELECT 
			id,
			code,
			name,
			? as resource_type,
			created_at,
			updated_at
		FROM ` + tableName + ` 
		WHERE resource_group_id = ?
		ORDER BY created_at DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, resourceType, groupID).Scan(&resources).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return resources, int64(len(resources)), nil
}

// MoveResourcesToGroup moves multiple resources to a target resource group
// It identifies the resource type by code prefix and table_name from resource_types table
// If targetGroupID is nil (not provided), resources are moved to the project's default resource group
func (r *resourceGroupRepository) MoveResourcesToGroup(ctx context.Context, resourceCodes []string, targetGroupID *uint, projectID uint) error {
	// If targetGroupID is nil (not provided), get the default resource group
	var finalGroupID uint
	if targetGroupID == nil {
		defaultGroup, err := r.GetDefaultResourceGroup(ctx, projectID)
		if err != nil {
			return err
		}
		finalGroupID = defaultGroup.ID
	} else {
		finalGroupID = *targetGroupID
	}

	// Get resource type mapping from resource_types table
	var resourceTypes []struct {
		Code      string
		TableName string
	}
	if err := r.db.WithContext(ctx).Table("resource_types").Select("code, table_name").Find(&resourceTypes).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Build a map from resource type prefix to table name
	typeToTableMap := make(map[string]string)
	for _, rt := range resourceTypes {
		typeToTableMap[rt.Code] = rt.TableName
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Group codes by resource type
		codesByType := make(map[string][]string)
		for _, code := range resourceCodes {
			// Extract resource type prefix from code (e.g., "server_123" -> "server")
			// Convention: code format is "{resource_type}_{identifier}"
			var resourceType string
			for prefix := range typeToTableMap {
				if len(code) > len(prefix) && code[:len(prefix)] == prefix && code[len(prefix)] == '_' {
					resourceType = prefix
					break
				}
			}

			if resourceType != "" {
				codesByType[resourceType] = append(codesByType[resourceType], code)
			}
		}

		// Execute update for each resource type
		for resourceType, codes := range codesByType {
			tableName, ok := typeToTableMap[resourceType]
			if !ok {
				continue
			}

			updateQuery := `
				UPDATE ` + tableName + `
				SET resource_group_id = ?, updated_at = CURRENT_TIMESTAMP
				WHERE code IN ?`

			if err := tx.Exec(updateQuery, finalGroupID, codes).Error; err != nil {
				return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
			}
		}
		return nil
	})
}

// GetResourceStatistics gets resource statistics by project or resource group
func (r *resourceGroupRepository) GetResourceStatistics(
	ctx context.Context,
	projectID, resourceGroupID *uint,
) (map[string]int, error) {
	statistics := make(map[string]int)

	// If querying by project_id, first get all resource group IDs in that project
	var resourceGroupIDs []uint
	if projectID != nil {
		if err := r.db.WithContext(ctx).
			Model(&model.ResourceGroup{}).
			Where("project_id = ?", *projectID).
			Pluck("id", &resourceGroupIDs).Error; err != nil {
			return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
		}
		// If no resource groups found in the project, return empty statistics
		if len(resourceGroupIDs) == 0 {
			return statistics, nil
		}
	}

	// Get resource type mapping from resource_types table
	var resourceTypes []struct {
		Code      string
		TableName string
	}
	if err := r.db.WithContext(ctx).Table("resource_types").Select("code, table_name").Find(&resourceTypes).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	for _, rt := range resourceTypes {
		query := r.db.WithContext(ctx).Table(rt.TableName)

		// Apply filters based on query conditions
		if projectID != nil {
			// Filter by resource groups that belong to the project
			query = query.Where("resource_group_id IN ?", resourceGroupIDs)
		} else if resourceGroupID != nil {
			// Filter by specific resource group
			query = query.Where("resource_group_id = ?", *resourceGroupID)
		}
		// If neither projectID nor resourceGroupID is provided, count all resources

		var count int64
		if err := query.Count(&count).Error; err != nil {
			return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
		}

		statistics[rt.Code] = int(count)
	}

	return statistics, nil
}
