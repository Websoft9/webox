package response

import "time"

// DatabaseConnectionResponse represents a database connection in response
type DatabaseConnectionResponse struct {
	ID              uint                   `json:"id"`
	Name            string                 `json:"name"`
	Code            string                 `json:"code"` // Unique connection code
	DBType          string                 `json:"db_type"`
	Host            string                 `json:"host"`
	Port            int                    `json:"port"`
	Database        *string                `json:"database"`         // Optional
	Description     *string                `json:"description"`      // Connection description
	Config          map[string]interface{} `json:"config,omitempty"` // Additional configuration
	OwnerID         uint                   `json:"owner_id"`
	ResourceGroupID *uint                  `json:"resource_group_id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// DatabaseConnectionDetailResponse represents detailed database connection information
// Credentials are managed separately via secret management system
// Frontend should use secret_references API to query associated secrets
type DatabaseConnectionDetailResponse struct {
	DatabaseConnectionResponse
}
