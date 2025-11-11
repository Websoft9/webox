package repository

import (
	"context"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// credentialCategoryRepository implements CredentialCategoryRepository
type credentialCategoryRepository struct {
	db *gorm.DB
}

// NewCredentialCategoryRepository creates a new credential category repository
func NewCredentialCategoryRepository(db *gorm.DB) repository.CredentialCategoryRepository {
	return &credentialCategoryRepository{
		db: db,
	}
}

// GetByID retrieves a credential category by ID
func (r *credentialCategoryRepository) GetByID(ctx context.Context, id uint) (*model.CredentialCategory, error) {
	var category model.CredentialCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &category, nil
}

// List retrieves all credential categories
func (r *credentialCategoryRepository) List(ctx context.Context) ([]*model.CredentialCategory, error) {
	var categories []*model.CredentialCategory
	err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return categories, nil
}
