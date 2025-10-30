package request

import "api-service/internal/dto/common"

// CreateEnvVarRequest represents the request to create an environment variable
type CreateEnvVarRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=64"`    // Variable name (letters, digits, and underscores only, must start with letter or underscore)
	Value       string  `json:"value" validate:"max=4096"`                // Variable value (can be empty for placeholder variables with default values)
	Description *string `json:"description" validate:"omitempty,max=500"` // Variable description
	IsSensitive bool    `json:"is_sensitive"`                             // Whether sensitive variable
}

// CreatePlatformEnvVarRequest represents the request to create a platform-level environment variable
type CreatePlatformEnvVarRequest struct {
	CreateEnvVarRequest
}

// CreateProjectEnvVarRequest represents the request to create a project-level environment variable
type CreateProjectEnvVarRequest struct {
	CreateEnvVarRequest
	ProjectID uint `json:"project_id" validate:"required"` // Project ID (required for project-level variables)
}

// UpdateEnvVarRequest represents the request to update an environment variable
type UpdateEnvVarRequest struct {
	Value       *string `json:"value" validate:"omitempty,max=4096"`      // Variable value
	Description *string `json:"description" validate:"omitempty,max=500"` // Variable description
	IsSensitive *bool   `json:"is_sensitive"`                             // Whether sensitive variable
}

// GetEnvVarListRequest represents the request to get environment variable list
type GetEnvVarListRequest struct {
	common.BaseListRequest
	Keywords    *string `form:"keywords" binding:"omitempty" json:"keywords"`         // Search by name or description
	IsSensitive *bool   `form:"is_sensitive" binding:"omitempty" json:"is_sensitive"` // Filter by sensitive flag
}

// GetProjectEnvVarListRequest represents the request to get project-level environment variable list
type GetProjectEnvVarListRequest struct {
	GetEnvVarListRequest
	ProjectID uint `form:"project_id" binding:"required" json:"project_id"` // Project ID (required)
}

// ResolveEnvVarRequest represents the request to resolve environment variable interpolation
type ResolveEnvVarRequest struct {
	Template  string `json:"template" validate:"required"`                     // Template string containing ${VAR_NAME} syntax
	Scope     string `json:"scope" validate:"required,oneof=platform project"` // Scope: platform or project
	ProjectID *uint  `json:"project_id" validate:"required_if=Scope project"`  // Project ID (required when scope=project)
}
