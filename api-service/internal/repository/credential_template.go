package repository

import (
	"context"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// credentialTemplateRepository implements CredentialTemplateRepository
type credentialTemplateRepository struct {
	db *gorm.DB
}

// NewCredentialTemplateRepository creates a new credential template repository
func NewCredentialTemplateRepository(db *gorm.DB) repository.CredentialTemplateRepository {
	return &credentialTemplateRepository{
		db: db,
	}
}

// GetByID retrieves a credential template by ID
func (r *credentialTemplateRepository) GetByID(ctx context.Context, id uint) (*model.CredentialTemplate, error) {
	var template model.CredentialTemplate
	err := r.db.WithContext(ctx).Preload("Category").First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &template, nil
}

// List retrieves credential templates with optional category filter
func (r *credentialTemplateRepository) List(ctx context.Context, categoryID *uint) ([]*model.CredentialTemplate, error) {
	var templates []*model.CredentialTemplate
	query := r.db.WithContext(ctx).Preload("Category")

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	err := query.Order("category_id ASC, id ASC").Find(&templates).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return templates, nil
}
