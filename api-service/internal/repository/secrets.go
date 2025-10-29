package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// secretRepository implements SecretRepository
type secretRepository struct {
	db *gorm.DB
}

// NewSecretRepository creates a new secret repository
func NewSecretRepository(db *gorm.DB) repository.SecretRepository {
	return &secretRepository{
		db: db,
	}
}

// Create creates a new secret
func (r *secretRepository) Create(ctx context.Context, secret *model.Secret) error {
	if err := r.db.WithContext(ctx).Create(secret).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// GetByID retrieves a secret by ID
func (r *secretRepository) GetByID(ctx context.Context, id uint) (*model.Secret, error) {
	var secret model.Secret
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&secret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &secret, nil
}

// GetByCode retrieves a secret by code
func (r *secretRepository) GetByCode(ctx context.Context, code string) (*model.Secret, error) {
	var secret model.Secret
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&secret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &secret, nil
}

// List retrieves a paginated list of secrets with filters
func (r *secretRepository) List(ctx context.Context, req *request.ListSecretsRequest, userID uint) ([]*model.Secret, int64, error) {
	var secrets []*model.Secret
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Secret{})

	// Filter by user permissions: owner or authorized user
	query = query.Where("owner_id = ? OR id IN (SELECT secret_id FROM secret_authorizes WHERE authorized_user_id = ?)", userID, userID)

	// Keyword search (name and description)
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// Filter by resource code (join with secret_references)
	if req.ResourceCode != nil && *req.ResourceCode != "" {
		query = query.Where("id IN (SELECT secret_id FROM secret_references WHERE resource_code = ?)", *req.ResourceCode)
	}

	// Time range filter
	startTime, endTime, parseErr := req.GetTimeRange(true)
	if parseErr != nil {
		return nil, 0, errors.NewAppErrorWrapError(parseErr, errors.CodeRecordQueryFailed)
	}
	if startTime != "" && endTime != "" {
		query = query.Where("updated_at BETWEEN ? AND ?", startTime, endTime)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination and sorting
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	orderBy := req.GetSortOrder()

	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&secrets).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return secrets, total, nil
}

// Update updates an existing secret
func (r *secretRepository) Update(ctx context.Context, secret *model.Secret) error {
	if err := r.db.WithContext(ctx).Model(secret).Updates(secret).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}
	return nil
}

// Delete deletes a secret by ID
func (r *secretRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Secret{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// ExistsByName checks if a secret name exists in a resource group
func (r *secretRepository) ExistsByName(ctx context.Context, resourceGroupID uint, name string, excludeID *uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&model.Secret{}).
		Where("resource_group_id = ? AND name = ?", resourceGroupID, name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// ExistsByCode checks if a secret code exists
func (r *secretRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Secret{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// GetReferenceCount gets the count of references for a secret
func (r *secretRepository) GetReferenceCount(ctx context.Context, secretID uint) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.SecretReference{}).Where("secret_id = ?", secretID).Count(&count).Error; err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return int(count), nil
}
