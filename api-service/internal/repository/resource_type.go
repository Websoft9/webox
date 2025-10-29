package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// resourceTypeRepository implements ResourceTypeRepository
type resourceTypeRepository struct {
	db *gorm.DB
}

// NewResourceTypeRepository creates a new resource type repository
func NewResourceTypeRepository(db *gorm.DB) repository.ResourceTypeRepository {
	return &resourceTypeRepository{
		db: db,
	}
}

// GetByCode retrieves a resource type by code
func (r *resourceTypeRepository) GetByCode(ctx context.Context, code string) (*model.ResourceType, error) {
	var resourceType model.ResourceType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&resourceType).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &resourceType, nil
}

// List retrieves all resource types
func (r *resourceTypeRepository) List(ctx context.Context) ([]*model.ResourceType, error) {
	var resourceTypes []*model.ResourceType
	if err := r.db.WithContext(ctx).Find(&resourceTypes).Error; err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return resourceTypes, nil
}
