package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// EnvironmentVariableRepository defines the interface for environment variable data access
type EnvironmentVariableRepository interface {
	// Create creates a new environment variable
	Create(ctx context.Context, envVar *model.EnvironmentVariable) error

	// GetByID retrieves an environment variable by ID
	GetByID(ctx context.Context, id uint) (*model.EnvironmentVariable, error)

	// GetByNameAndScope retrieves an environment variable by name, scope, and project_id
	GetByNameAndScope(ctx context.Context, name string, scope model.EnvVarScope, projectID *uint) (*model.EnvironmentVariable, error)

	// GetPlatformList retrieves a paginated list of platform-level environment variables
	GetPlatformList(ctx context.Context, req *request.GetEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error)

	// GetProjectList retrieves a paginated list of project-level environment variables
	GetProjectList(ctx context.Context, req *request.GetProjectEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error)

	// GetAllByScope retrieves all environment variables by scope and optional project ID
	GetAllByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) ([]*model.EnvironmentVariable, error)

	// Update updates an existing environment variable
	Update(ctx context.Context, envVar *model.EnvironmentVariable) error

	// Delete deletes an environment variable by ID
	Delete(ctx context.Context, id uint) error

	// CountByScope counts environment variables by scope and optional project ID
	CountByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) (int64, error)

	// ExistsByName checks if an environment variable exists with the given name, scope, and project_id
	ExistsByName(ctx context.Context, name string, scope model.EnvVarScope, projectID, excludeID *uint) (bool, error)
}
