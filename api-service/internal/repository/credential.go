package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// credentialRepository implements CredentialRepository
type credentialRepository struct {
	db *gorm.DB
}

// NewCredentialRepository creates a new credential repository
func NewCredentialRepository(db *gorm.DB) repository.CredentialRepository {
	return &credentialRepository{
		db: db,
	}
}

// Create creates a new credential
func (r *credentialRepository) Create(ctx context.Context, credential *model.Credential) error {
	err := r.db.WithContext(ctx).Create(credential).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// Update updates an existing credential
func (r *credentialRepository) Update(ctx context.Context, credential *model.Credential) error {
	err := r.db.WithContext(ctx).Save(credential).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}
	return nil
}

// Delete deletes a credential by ID
func (r *credentialRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&model.Credential{}, id).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// GetByID retrieves a credential by ID
func (r *credentialRepository) GetByID(ctx context.Context, id uint) (*model.Credential, error) {
	var credential model.Credential
	err := r.db.WithContext(ctx).
		Preload("Template").
		Preload("Template.Category").
		Preload("Owner").
		First(&credential, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &credential, nil
}

// GetByName retrieves a credential by name
func (r *credentialRepository) GetByName(ctx context.Context, name string) (*model.Credential, error) {
	var credential model.Credential
	err := r.db.WithContext(ctx).
		Preload("Template").
		Preload("Template.Category").
		Preload("Owner").
		Where("name = ?", name).
		First(&credential).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &credential, nil
}

// ExistsByName checks if a credential name exists (excluding specific ID)
func (r *credentialRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Credential{}).Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}

// List retrieves credentials with pagination and filters
func (r *credentialRepository) List(ctx context.Context, req *request.ListCredentialsRequest) ([]*model.Credential, int64, error) {
	var credentials []*model.Credential
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Credential{})

	// Apply filters
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	if req.CategoryID != nil {
		query = query.Joins("JOIN credential_templates ON credentials.template_id = credential_templates.id").
			Where("credential_templates.category_id = ?", *req.CategoryID)
	}

	if req.TemplateID != nil {
		query = query.Where("template_id = ?", *req.TemplateID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	query = query.Offset(offset).Order(req.GetSortOrder()).Limit(pageSize)

	// Preload associations
	query = query.Preload("Template").
		Preload("Template.Category").
		Preload("Owner")

	// Execute query
	if err := query.Find(&credentials).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return credentials, total, nil
}
