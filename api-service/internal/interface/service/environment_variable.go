package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// EnvironmentVariableService defines the interface for environment variable business logic
type EnvironmentVariableService interface {
	// CreatePlatformEnvVar creates a new platform-level environment variable
	CreatePlatformEnvVar(ctx context.Context, req *request.CreatePlatformEnvVarRequest, ownerID uint) (*response.EnvironmentVariableResponse, error)

	// CreateProjectEnvVar creates a new project-level environment variable
	CreateProjectEnvVar(ctx context.Context, req *request.CreateProjectEnvVarRequest, ownerID uint) (*response.EnvironmentVariableResponse, error)

	// GetEnvVar retrieves an environment variable by ID
	GetEnvVar(ctx context.Context, id uint) (*response.EnvironmentVariableDetailResponse, error)

	// GetPlatformEnvVarList retrieves a paginated list of platform-level environment variables
	GetPlatformEnvVarList(ctx context.Context, req *request.GetEnvVarListRequest) (*common.PaginationResponse, error)

	// GetProjectEnvVarList retrieves a paginated list of project-level environment variables
	GetProjectEnvVarList(ctx context.Context, req *request.GetProjectEnvVarListRequest) (*common.PaginationResponse, error)

	// UpdateEnvVar updates an existing environment variable
	UpdateEnvVar(ctx context.Context, id uint, req *request.UpdateEnvVarRequest) (*response.EnvironmentVariableResponse, error)

	// DeleteEnvVar deletes an environment variable
	DeleteEnvVar(ctx context.Context, id uint) error

	// ResolveEnvVar resolves environment variable interpolation in template string
	ResolveEnvVar(ctx context.Context, req *request.ResolveEnvVarRequest) (*response.ResolveEnvVarResponse, error)
}
