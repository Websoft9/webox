package request

import "api-service/internal/dto/common"

// CreateDatabaseConnectionRequest represents the request to create a database connection
type CreateDatabaseConnectionRequest struct {
	Name            string                 `json:"name" validate:"required,min=2,max=64"`
	DBType          string                 `json:"db_type" validate:"required,oneof=mysql postgresql mariadb sqlserver oracle sqlite"`
	Host            string                 `json:"host" validate:"required,max=255"`
	Port            int                    `json:"port" validate:"required,min=1,max=65535"`
	Database        *string                `json:"database" validate:"omitempty,max=64"`     // Optional, for scenarios like PostgreSQL server connection
	Description     *string                `json:"description" validate:"omitempty,max=255"` // Connection description
	Config          map[string]interface{} `json:"config" validate:"omitempty"`              // Additional configuration (charset, timeout, ssl, schema, etc.)
	ResourceGroupID *uint                  `json:"resource_group_id"`
}

// UpdateDatabaseConnectionRequest represents the request to update a database connection
type UpdateDatabaseConnectionRequest struct {
	Name            *string                `json:"name" validate:"omitempty,min=2,max=64"`
	Host            *string                `json:"host" validate:"omitempty,max=255"`
	Port            *int                   `json:"port" validate:"omitempty,min=1,max=65535"`
	Database        *string                `json:"database" validate:"omitempty,max=64"`
	Description     *string                `json:"description" validate:"omitempty,max=255"`
	Config          map[string]interface{} `json:"config" validate:"omitempty"`
	ResourceGroupID *uint                  `json:"resource_group_id"`
}

// GetDatabaseConnectionListRequest represents the request to get database connection list
type GetDatabaseConnectionListRequest struct {
	common.BaseListRequest
	DBType *string `form:"db_type" binding:"omitempty,oneof=mysql postgresql mariadb sqlserver oracle sqlite" json:"db_type"`
	Code   *string `form:"code" binding:"omitempty,max=64" json:"code"` // Filter by connection code
}
