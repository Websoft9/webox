package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// DatabaseConnectionRepository defines the interface for database connection data access
type DatabaseConnectionRepository interface {
	// Create creates a new database connection
	Create(ctx context.Context, conn *model.DatabaseConnection) error

	// GetByID retrieves a database connection by ID
	GetByID(ctx context.Context, id uint) (*model.DatabaseConnection, error)

	// GetList retrieves a paginated list of database connections with filters
	GetList(ctx context.Context, req *request.GetDatabaseConnectionListRequest, ownerID uint) ([]*model.DatabaseConnection, int64, error)

	// Update updates an existing database connection
	Update(ctx context.Context, conn *model.DatabaseConnection) error

	// Delete deletes a database connection by ID
	Delete(ctx context.Context, id uint) error
}
