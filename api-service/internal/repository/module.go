package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

type moduleRepository struct {
	db *gorm.DB
}

// NewModuleRepository creates a module repository instance
func NewModuleRepository(db *gorm.DB) repository.ModuleRepository {
	return &moduleRepository{db: db}
}

// GetByCode retrieves a module by code
func (r *moduleRepository) GetByCode(ctx context.Context, code string) (*model.Module, error) {
	var module model.Module
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&module).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &module, nil
}
