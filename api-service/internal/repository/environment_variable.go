package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// environmentVariableRepository implements EnvironmentVariableRepository
type environmentVariableRepository struct {
	db *gorm.DB
}

// NewEnvironmentVariableRepository creates a new environment variable repository
func NewEnvironmentVariableRepository(db *gorm.DB) repository.EnvironmentVariableRepository {
	return &environmentVariableRepository{
		db: db,
	}
}

// Create creates a new environment variable
func (r *environmentVariableRepository) Create(ctx context.Context, envVar *model.EnvironmentVariable) error {
	if err := r.db.WithContext(ctx).Create(envVar).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// GetByID retrieves an environment variable by ID
func (r *environmentVariableRepository) GetByID(ctx context.Context, id uint) (*model.EnvironmentVariable, error) {
	var envVar model.EnvironmentVariable
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&envVar).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &envVar, nil
}

// GetByNameAndScope retrieves an environment variable by name, scope, and project_id
func (r *environmentVariableRepository) GetByNameAndScope(ctx context.Context, name string, scope model.EnvVarScope, projectID *uint) (*model.EnvironmentVariable, error) {
	var envVar model.EnvironmentVariable
	query := r.db.WithContext(ctx).Where("name = ? AND scope = ?", name, scope)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	} else {
		query = query.Where("project_id IS NULL")
	}

	if err := query.First(&envVar).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &envVar, nil
}

// GetPlatformList retrieves a paginated list of platform-level environment variables
func (r *environmentVariableRepository) GetPlatformList(ctx context.Context, req *request.GetEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error) {
	var envVars []*model.EnvironmentVariable
	var total int64

	query := r.db.WithContext(ctx).Model(&model.EnvironmentVariable{})

	// Filter by platform scope
	query = query.Where("scope = ?", model.EnvVarScopePlatform)

	// Filter by keywords (search in name or description)
	if req.Keywords != nil && *req.Keywords != "" {
		searchPattern := "%" + *req.Keywords + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// Filter by is_sensitive
	if req.IsSensitive != nil {
		query = query.Where("is_sensitive = ?", *req.IsSensitive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	// Execute query
	if err := query.Find(&envVars).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return envVars, total, nil
}

// GetProjectList retrieves a paginated list of project-level environment variables
func (r *environmentVariableRepository) GetProjectList(ctx context.Context, req *request.GetProjectEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error) {
	var envVars []*model.EnvironmentVariable
	var total int64

	query := r.db.WithContext(ctx).Model(&model.EnvironmentVariable{})

	// Filter by project scope and project_id
	query = query.Where("scope = ? AND project_id = ?", model.EnvVarScopeProject, req.ProjectID)

	// Filter by keywords (search in name or description)
	if req.Keywords != nil && *req.Keywords != "" {
		searchPattern := "%" + *req.Keywords + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// Filter by is_sensitive
	if req.IsSensitive != nil {
		query = query.Where("is_sensitive = ?", *req.IsSensitive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	// Execute query
	if err := query.Find(&envVars).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return envVars, total, nil
}

// GetAllByScope retrieves all environment variables by scope and optional project ID
func (r *environmentVariableRepository) GetAllByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) ([]*model.EnvironmentVariable, error) {
	var envVars []*model.EnvironmentVariable

	query := r.db.WithContext(ctx).Where("scope = ?", scope)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	} else {
		query = query.Where("project_id IS NULL")
	}

	if err := query.Find(&envVars).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return envVars, nil
}

// Update updates an existing environment variable
func (r *environmentVariableRepository) Update(ctx context.Context, envVar *model.EnvironmentVariable) error {
	if err := r.db.WithContext(ctx).Save(envVar).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}
	return nil
}

// Delete deletes an environment variable by ID
func (r *environmentVariableRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.EnvironmentVariable{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// CountByScope counts environment variables by scope and optional project ID
func (r *environmentVariableRepository) CountByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.EnvironmentVariable{}).Where("scope = ?", scope)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	} else {
		query = query.Where("project_id IS NULL")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count, nil
}

// ExistsByName checks if an environment variable exists with the given name, scope, and project_id
func (r *environmentVariableRepository) ExistsByName(
	ctx context.Context,
	name string,
	scope model.EnvVarScope,
	projectID, excludeID *uint,
) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.EnvironmentVariable{}).
		Where("name = ? AND scope = ?", name, scope)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	} else {
		query = query.Where("project_id IS NULL")
	}

	// Exclude specific ID (for update operations)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}
