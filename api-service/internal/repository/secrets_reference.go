package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// secretReferenceRepository implements SecretReferenceRepository
type secretReferenceRepository struct {
	db *gorm.DB
}

// NewSecretReferenceRepository creates a new secret reference repository
func NewSecretReferenceRepository(db *gorm.DB) repository.SecretReferenceRepository {
	return &secretReferenceRepository{
		db: db,
	}
}

// Create creates a new secret reference
func (r *secretReferenceRepository) Create(ctx context.Context, reference *model.SecretReference) error {
	if err := r.db.WithContext(ctx).Create(reference).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// ListBySecretID retrieves all references for a secret
func (r *secretReferenceRepository) ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretReference, error) {
	var references []*model.SecretReference
	if err := r.db.WithContext(ctx).Where("secret_id = ?", secretID).Find(&references).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return references, nil
}

// ListByResourceCode retrieves all references for a resource
func (r *secretReferenceRepository) ListByResourceCode(ctx context.Context, resourceCode string) ([]*model.SecretReference, error) {
	var references []*model.SecretReference
	if err := r.db.WithContext(ctx).Where("resource_code = ?", resourceCode).Find(&references).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return references, nil
}

// Exists checks if a reference exists
func (r *secretReferenceRepository) Exists(ctx context.Context, secretID uint, resourceCode string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.SecretReference{}).
		Where("secret_id = ? AND resource_code = ?", secretID, resourceCode).
		Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// Delete deletes a specific reference
func (r *secretReferenceRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.SecretReference{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// DeleteBySecretID deletes all references for a secret
func (r *secretReferenceRepository) DeleteBySecretID(ctx context.Context, secretID uint) error {
	if err := r.db.WithContext(ctx).Where("secret_id = ?", secretID).Delete(&model.SecretReference{}).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// HasActiveReferences checks if a secret has any active references
func (r *secretReferenceRepository) HasActiveReferences(ctx context.Context, secretID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.SecretReference{}).
		Where("secret_id = ?", secretID).
		Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}
