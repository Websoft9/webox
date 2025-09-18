package response

import (
	"time"
)

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID          uint64    `json:"id" example:"1"`
	Name        string    `json:"name" example:"production"`
	Color       string    `json:"color" example:"#ff0000"`
	Description string    `json:"description" example:"Production environment tag"`
	CreatedBy   uint64    `json:"createdBy" example:"1"`
	CreatedAt   time.Time `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
	UsageCount  int64     `json:"usageCount,omitempty" example:"25"`
}

// TagSimpleResponse represents a simplified tag in API responses
type TagSimpleResponse struct {
	ID    uint64 `json:"id" example:"1"`
	Name  string `json:"name" example:"production"`
	Color string `json:"color" example:"#ff0000"`
}

// TagAssignResult represents the result of a tag assignment operation
type TagAssignResult struct {
	Name    string `json:"name" example:"new-feature"`
	TagID   uint64 `json:"tagId" example:"3"`
	Status  string `json:"status" example:"created"`
	Message string `json:"message" example:"Tag created and associated"`
}

// TagAssignResponse represents the response for tag assignment
type TagAssignResponse struct {
	Results      []TagAssignResult   `json:"results"`
	ResourceTags []TagSimpleResponse `json:"resourceTags"`
}

// TagUnassignResponse represents the response for tag unassignment
type TagUnassignResponse struct {
	RemovedCount int    `json:"removedCount" example:"2"`
	Message      string `json:"message" example:"Removed 2 tag associations"`
}

// TaggedResource represents a resource with its associated tags
type TaggedResource struct {
	ResourceID   uint64              `json:"resourceId" example:"1"`
	ResourceName string              `json:"resourceName" example:"WordPress Blog"`
	Tags         []TagSimpleResponse `json:"tags"`
	MatchedTags  []uint64            `json:"matchedTags" example:"[1,3]"`
	CreatedAt    time.Time           `json:"createdAt" example:"2024-01-15T10:30:00Z"`
}

// TagSearchResponse represents the response for tag search
type TagSearchResponse struct {
	Total     int64            `json:"total" example:"25"`
	Page      int              `json:"page" example:"1"`
	PageSize  int              `json:"pageSize" example:"20"`
	Resources []TaggedResource `json:"resources"`
}
