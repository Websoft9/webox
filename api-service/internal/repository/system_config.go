package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"
	"strings"

	"gorm.io/gorm"
)

// systemConfigRepository implements the SystemConfigRepository interface
type systemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository creates a new SystemConfig repository instance
func NewSystemConfigRepository(db *gorm.DB) repository.SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

// Create creates a new system configuration
func (r *systemConfigRepository) Create(ctx context.Context, config *model.SystemConfig) error {
	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByKey retrieves a system configuration by key
func (r *systemConfigRepository) GetByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig

	err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &config, nil
}

// GetByID retrieves a system configuration by ID
func (r *systemConfigRepository) GetByID(ctx context.Context, id uint) (*model.SystemConfig, error) {
	var config model.SystemConfig

	err := r.db.WithContext(ctx).First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &config, nil
}

// List retrieves system configurations with filtering
func (r *systemConfigRepository) List(ctx context.Context, filter *request.ListSystemConfigsFilter) ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig

	query := r.db.WithContext(ctx).Model(&model.SystemConfig{})

	// Apply category filter
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	// Apply keyword filter (search in config_key and description)
	if filter.Keyword != "" {
		keyword := "%" + strings.ToLower(filter.Keyword) + "%"
		query = query.Where("LOWER(config_key) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	// Order by sort_order, then by config_key
	query = query.Order("sort_order ASC, config_key ASC")

	err := query.Find(&configs).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return configs, nil
}

// Update updates an existing system configuration
func (r *systemConfigRepository) Update(ctx context.Context, config *model.SystemConfig) error {
	result := r.db.WithContext(ctx).Updates(config)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}

// UpdateValue updates the value of a system configuration by key
func (r *systemConfigRepository) UpdateValue(ctx context.Context, key, value string) error {
	result := r.db.WithContext(ctx).
		Model(&model.SystemConfig{}).
		Where("config_key = ?", key).
		Update("config_value", value)

	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}

// Delete deletes a system configuration by ID
func (r *systemConfigRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.SystemConfig{}, id)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}

	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}

	return nil
}
