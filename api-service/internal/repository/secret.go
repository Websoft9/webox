package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// secretKeyRepository implements SecretKeyRepository interface
type secretKeyRepository struct {
	db *gorm.DB
}

// NewSecretKeyRepository creates a new secret key repository
func NewSecretKeyRepository(db *gorm.DB) repository.SecretKeyRepository {
	return &secretKeyRepository{
		db: db,
	}
}

// Create creates a new secret key
func (r *secretKeyRepository) Create(ctx context.Context, secretKey *model.SecretKey) error {
	if err := r.db.WithContext(ctx).Create(secretKey).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// GetByID retrieves a secret key by ID
func (r *secretKeyRepository) GetByID(ctx context.Context, id uint) (*model.SecretKey, error) {
	var secretKey model.SecretKey
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("id = ?", id).
		First(&secretKey).Error

	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordNotFound)
	}

	return &secretKey, nil
}

// Update updates an existing secret key
func (r *secretKeyRepository) Update(ctx context.Context, secretKey *model.SecretKey) error {
	if err := r.db.WithContext(ctx).Updates(secretKey).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}
	return nil
}

// Delete soft deletes a secret key by ID
func (r *secretKeyRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.SecretKey{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// List retrieves secret keys with pagination and filtering
func (r *secretKeyRepository) List(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) ([]*model.SecretKey, int64, error) {
	var secretKeys []*model.SecretKey
	var total int64

	query := r.db.WithContext(ctx).Model(&model.SecretKey{})

	// Apply filters
	query = query.Where("owner_id = ?", userID)

	if req.KeyType != nil {
		query = query.Where("key_type = ?", *req.KeyType)
	}

	// Keyword filter: fuzzy search on name and description
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// ResourceCode filter: join with secret_references table
	if req.ResourceCode != nil && *req.ResourceCode != "" {
		query = query.Joins("INNER JOIN secret_references ON secret_references.secret_id = secret_keys.id").
			Where("secret_references.resource_code = ?", *req.ResourceCode)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()
	err := query.
		Preload("Owner").
		Order("secret_keys.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&secretKeys).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return secretKeys, total, nil
}

// ListAll retrieves all secret keys for a user without pagination (for export)
func (r *secretKeyRepository) ListAll(ctx context.Context, userID uint) ([]*model.SecretKey, error) {
	var secretKeys []*model.SecretKey

	err := r.db.WithContext(ctx).
		Model(&model.SecretKey{}).
		Where("owner_id = ?", userID).
		Preload("Owner").
		Order("created_at DESC").
		Find(&secretKeys).Error

	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return secretKeys, nil
}

// GetByOwnerID retrieves secret keys by owner ID
func (r *secretKeyRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]*model.SecretKey, error) {
	var secretKeys []*model.SecretKey
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&secretKeys).Error

	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return secretKeys, nil
}

// ExistsByName checks if a secret key with the given name exists for a user
func (r *secretKeyRepository) ExistsByName(ctx context.Context, name string, ownerID uint, excludeID ...uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&model.SecretKey{}).
		Where("name = ? AND owner_id = ?", name, ownerID)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// CountByType counts secret keys by type for a user
func (r *secretKeyRepository) CountByType(ctx context.Context, keyType model.SecretKeyType, ownerID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.SecretKey{}).
		Where("key_type = ? AND owner_id = ?", keyType, ownerID).
		Count(&count).Error

	if err != nil {
		return 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count, nil
}

// CreateUserSecret creates a user secret relationship
func (r *secretKeyRepository) CreateUserSecret(ctx context.Context, userSecret *model.UserSecret) error {
	return r.db.WithContext(ctx).Create(userSecret).Error
}

// DeleteUserSecretsBySecretKeyID deletes all user secret relationships for a secret key
func (r *secretKeyRepository) DeleteUserSecretsBySecretKeyID(ctx context.Context, secretKeyID uint) error {
	return r.db.WithContext(ctx).
		Where("secret_key_id = ?", secretKeyID).
		Delete(&model.UserSecret{}).Error
}

// CheckUserSecretAccess checks if a user has access to a secret key
func (r *secretKeyRepository) CheckUserSecretAccess(ctx context.Context, userID, secretKeyID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserSecret{}).
		Where("user_id = ? AND secret_key_id = ?", userID, secretKeyID).
		Count(&count).Error

	if err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// CreateSecretReference creates a secret reference record
func (r *secretKeyRepository) CreateSecretReference(ctx context.Context, reference *model.SecretReference) error {
	return r.db.WithContext(ctx).Create(reference).Error
}

// DeleteSecretReferencesBySecretID deletes all references for a secret key
func (r *secretKeyRepository) DeleteSecretReferencesBySecretID(ctx context.Context, secretID uint) error {
	return r.db.WithContext(ctx).Where("secret_id = ?", secretID).Delete(&model.SecretReference{}).Error
}

// DeleteSecretReferencesByResourceCode deletes secret references by resource code
func (r *secretKeyRepository) DeleteSecretReferencesByResourceCode(ctx context.Context, resourceCode string) error {
	return r.db.WithContext(ctx).Where("resource_code = ?", resourceCode).Delete(&model.SecretReference{}).Error
}

// GetSecretReferencesByResourceCode gets secret references by resource code
func (r *secretKeyRepository) GetSecretReferencesByResourceCode(ctx context.Context, resourceCode string) ([]*model.SecretReference, error) {
	var references []*model.SecretReference
	err := r.db.WithContext(ctx).
		Where("resource_code = ?", resourceCode).
		Find(&references).Error

	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return references, nil
}

// AssignSecretToResource assigns a secret to a resource
func (r *secretKeyRepository) AssignSecretToResource(ctx context.Context, secretID uint, resourceCode string) error {
	// Check if the assignment already exists
	var existingRef model.SecretReference
	err := r.db.WithContext(ctx).Where("secret_id = ? AND resource_code = ?", secretID, resourceCode).First(&existingRef).Error
	if err == nil {
		// Assignment already exists, return success
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create new assignment
	reference := &model.SecretReference{
		SecretID:     secretID,
		ResourceCode: resourceCode,
	}

	return r.db.WithContext(ctx).Create(reference).Error
}

// UnassignSecretFromResource unassigns a secret from a resource
func (r *secretKeyRepository) UnassignSecretFromResource(ctx context.Context, secretID uint, resourceCode string) error {
	result := r.db.WithContext(ctx).Where("secret_id = ? AND resource_code = ?", secretID, resourceCode).Delete(&model.SecretReference{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetSecretReferences gets all resource references for a secret
func (r *secretKeyRepository) GetSecretReferences(ctx context.Context, secretID uint) ([]*model.SecretReference, error) {
	var references []*model.SecretReference
	err := r.db.WithContext(ctx).Where("secret_id = ?", secretID).Find(&references).Error
	if err != nil {
		return nil, err
	}
	return references, nil
}

// GetSecretReferenceBySecretAndResource gets a specific secret reference
func (r *secretKeyRepository) GetSecretReferenceBySecretAndResource(ctx context.Context, secretID uint, resourceCode string) (*model.SecretReference, error) {
	var reference model.SecretReference
	err := r.db.WithContext(ctx).Where("secret_id = ? AND resource_code = ?", secretID, resourceCode).First(&reference).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &reference, nil
}
