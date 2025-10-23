package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// DatabaseConnectionService defines the interface for database connection business logic
type DatabaseConnectionService interface {
	// CreateConnection creates a new database connection (without credentials)
	CreateConnection(ctx context.Context, req *request.CreateDatabaseConnectionRequest, ownerID uint) (*response.DatabaseConnectionResponse, error)

	// GetConnection retrieves a database connection by ID (without credentials)
	GetConnection(ctx context.Context, id uint, ownerID uint) (*response.DatabaseConnectionDetailResponse, error)

	// GetConnectionList retrieves a paginated list of database connections
	GetConnectionList(ctx context.Context, req *request.GetDatabaseConnectionListRequest, ownerID uint) (*common.PaginationResponse, error)

	// UpdateConnection updates an existing database connection (without credentials)
	UpdateConnection(ctx context.Context, id uint, req *request.UpdateDatabaseConnectionRequest, ownerID uint) (*response.DatabaseConnectionResponse, error)

	// DeleteConnection deletes a database connection
	DeleteConnection(ctx context.Context, id uint, ownerID uint) error
}
