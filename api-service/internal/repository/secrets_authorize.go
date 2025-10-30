package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// secretAuthorizeRepository implements SecretAuthorizeRepository
type secretAuthorizeRepository struct {
	db *gorm.DB
}

// NewSecretAuthorizeRepository creates a new secret authorize repository
func NewSecretAuthorizeRepository(db *gorm.DB) repository.SecretAuthorizeRepository {
	return &secretAuthorizeRepository{
		db: db,
	}
}

// Create creates a new secret authorization
func (r *secretAuthorizeRepository) Create(ctx context.Context, authorize *model.SecretAuthorize) error {
	if err := r.db.WithContext(ctx).Create(authorize).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// BatchCreate creates multiple secret authorizations
func (r *secretAuthorizeRepository) BatchCreate(ctx context.Context, authorizes []*model.SecretAuthorize) error {
	if len(authorizes) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Create(&authorizes).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// ListBySecretID retrieves all authorizations for a secret
func (r *secretAuthorizeRepository) ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretAuthorize, error) {
	var authorizes []*model.SecretAuthorize
	if err := r.db.WithContext(ctx).Where("secret_id = ?", secretID).Find(&authorizes).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return authorizes, nil
}

// Delete deletes a specific authorization
func (r *secretAuthorizeRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.SecretAuthorize{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// DeleteBySecretID deletes all authorizations for a secret
func (r *secretAuthorizeRepository) DeleteBySecretID(ctx context.Context, secretID uint) error {
	if err := r.db.WithContext(ctx).Where("secret_id = ?", secretID).Delete(&model.SecretAuthorize{}).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// IsAuthorized checks if a user is authorized to access a secret
func (r *secretAuthorizeRepository) IsAuthorized(ctx context.Context, secretID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.SecretAuthorize{}).
		Where("secret_id = ? AND authorized_user_id = ?", secretID, userID).
		Count(&count).Error; err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}
