package response

import "time"

// ResourceGroupResponse represents the basic resource group response
type ResourceGroupResponse struct {
	ID          uint      `json:"id"`
	ProjectID   uint      `json:"project_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description *string   `json:"description"`
	OwnerID     uint      `json:"owner_id"`
	IsDefault   bool      `json:"is_default"` // Whether this is the default resource group
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ResourceGroupDetailResponse represents the detailed resource group response
type ResourceGroupDetailResponse struct {
	ResourceGroupResponse
}

// ResourceItemResponse represents a single resource in the resource group
type ResourceItemResponse struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	ResourceType string    `json:"resource_type"` // Resource type code (server, database, secret, etc.)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ResourceGroupResourcesResponse represents the resources list in a resource group
type ResourceGroupResourcesResponse struct {
	Total     int64                  `json:"total"`
	Resources []ResourceItemResponse `json:"resources"`
}

// ResourceStatisticsResponse represents the resource statistics response
type ResourceStatisticsResponse struct {
	Total         int64          `json:"total"`          // Total number of resources
	ResourceTypes map[string]int `json:"resource_types"` // Statistics by resource type, e.g., {"server": 10, "database": 5}
}
