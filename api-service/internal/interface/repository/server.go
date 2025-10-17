package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// ServerRepository defines the interface for server data access operations
type ServerRepository interface {
	// Basic CRUD operations
	CreateServer(ctx context.Context, server *model.Server) error
	GetServerByID(ctx context.Context, id uint) (*model.Server, error)
	GetServerByName(ctx context.Context, name string) (*model.Server, error)
	UpdateServer(ctx context.Context, server *model.Server) error
	DeleteServer(ctx context.Context, id uint) error

	// Query operations
	ListServers(ctx context.Context, req *request.ListServersRequest) ([]*model.Server, int64, error)
	GetServersByStatus(ctx context.Context, status string) ([]*model.Server, error)
	GetServersByIDs(ctx context.Context, ids []uint) ([]*model.Server, error)

	// Batch operations
	BatchUpdateServerStatus(ctx context.Context, ids []uint, status string) error

	// Status operations
	UpdateServerStatus(ctx context.Context, id uint, status string) error
	UpdateServerLastSeen(ctx context.Context, id uint) error

	// Utility methods
	ExistsServerByName(ctx context.Context, name string, excludeID ...uint) (bool, error)
	CountServersByStatus(ctx context.Context) (map[string]int64, error)
	GetServerAgentsByServerID(ctx context.Context, serverID uint) ([]*model.ServerAgent, error)
}

// ServerAgentRepository defines the interface for server agent data access operations
type ServerAgentRepository interface {
	// Basic CRUD operations
	CreateAgent(ctx context.Context, agent *model.ServerAgent) error
	GetAgentByID(ctx context.Context, id uint) (*model.ServerAgent, error)
	GetAgentByServerID(ctx context.Context, serverID uint) (*model.ServerAgent, error)
	UpdateAgent(ctx context.Context, agent *model.ServerAgent) error
	DeleteAgent(ctx context.Context, id uint) error

	// Status operations
	UpdateAgentLastSeen(ctx context.Context, id uint) error
}
