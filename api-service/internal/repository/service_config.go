package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"

	"gorm.io/gorm"
)

// serviceConfigRepository implements the ServiceConfigRepository interface
type serviceConfigRepository struct {
	db *gorm.DB
}

// NewServiceConfigRepository creates a new ServiceConfig repository instance
func NewServiceConfigRepository(db *gorm.DB) repository.ServiceConfigRepository {
	return &serviceConfigRepository{db: db}
}

// Create creates a new service configuration
func (r *serviceConfigRepository) Create(ctx context.Context, config *model.ServiceConfig) error {
	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// GetByCode retrieves a service configuration by code
func (r *serviceConfigRepository) GetByCode(ctx context.Context, code string) (*model.ServiceConfig, error) {
	var config model.ServiceConfig
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &config, nil
}

// GetByID retrieves a service configuration by ID
func (r *serviceConfigRepository) GetByID(ctx context.Context, id uint) (*model.ServiceConfig, error) {
	var config model.ServiceConfig
	err := r.db.WithContext(ctx).First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &config, nil
}

// List retrieves all service configurations
func (r *serviceConfigRepository) List(ctx context.Context) ([]*model.ServiceConfig, error) {
	var configs []*model.ServiceConfig
	err := r.db.WithContext(ctx).Order("sort_order ASC, code ASC").Find(&configs).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return configs, nil
}

// ListByOwner retrieves service configurations by owner ID
func (r *serviceConfigRepository) ListByOwner(ctx context.Context, ownerID uint) ([]*model.ServiceConfig, error) {
	var configs []*model.ServiceConfig
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("sort_order ASC, code ASC").
		Find(&configs).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return configs, nil
}

// ListByCategory retrieves service configurations by category
func (r *serviceConfigRepository) ListByCategory(ctx context.Context, category string) ([]*model.ServiceConfig, error) {
	var configs []*model.ServiceConfig
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("sort_order ASC, code ASC").
		Find(&configs).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return configs, nil
}

// Update updates an existing service configuration
func (r *serviceConfigRepository) Update(ctx context.Context, config *model.ServiceConfig) error {
	result := r.db.WithContext(ctx).Updates(config)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}

// UpdateValue updates the value of a service configuration by code
func (r *serviceConfigRepository) UpdateValue(ctx context.Context, code, value string) error {
	result := r.db.WithContext(ctx).
		Model(&model.ServiceConfig{}).
		Where("code = ?", code).
		Update("config_value", value)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}

// Delete deletes a service configuration by ID
func (r *serviceConfigRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.ServiceConfig{}, id)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}
