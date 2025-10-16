package repository

import (
	"context"

	"api-service/internal/model"
)

// ModuleRepository module repository interface
type ModuleRepository interface {
	// GetByCode retrieves a module by code
	GetByCode(ctx context.Context, code string) (*model.Module, error)
}
