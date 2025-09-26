package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"

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
		return err
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
		return nil, err
	}

	return &secretKey, nil
}

// Update updates an existing secret key
func (r *secretKeyRepository) Update(ctx context.Context, secretKey *model.SecretKey) error {
	if err := r.db.WithContext(ctx).Save(secretKey).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a secret key by ID
func (r *secretKeyRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.SecretKey{}, id).Error; err != nil {
		return err
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

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Preload("Owner").
		Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&secretKeys).Error

	if err != nil {
		return nil, 0, err
	}

	return secretKeys, total, nil
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
		return nil, err
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
		return false, err
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
		return 0, err
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
