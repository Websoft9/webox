package request

import "api-service/internal/dto/common"

// TagCreateRequest represents a request to create a tag
type TagCreateRequest struct {
	Name        string `json:"name" binding:"required,max=128" example:"production"`
	Color       string `json:"color" binding:"max=16" example:"#ff0000"`
	Description string `json:"description" binding:"max=500" example:"Production environment tag"`
}

// TagUpdateRequest represents a request to update a tag
type TagUpdateRequest struct {
	Name        string `json:"name" binding:"required,max=128" example:"production"`
	Color       string `json:"color" binding:"max=16" example:"#ff0000"`
	Description string `json:"description" binding:"max=500" example:"Production environment tag"`
}

// TagListRequest represents a request to list tags
type TagListRequest struct {
	Search     string `form:"search" json:"search" binding:"omitempty" example:"prod"`
	ExcludeIDs string `form:"excludeIds" json:"excludeIds" binding:"omitempty" example:"1,2,3"`
}

// TagAssignRequest represents a request to assign tags to a resource
type TagAssignRequest struct {
	ResourceID uint     `json:"resourceId" binding:"required" example:"123"`
	TagIDs     []uint   `json:"tagIds"`
	TagNames   []string `json:"tagNames"`
	ReplaceAll bool     `json:"replaceAll" example:"false"`
}

// TagUnassignRequest represents a request to unassign tags from a resource
type TagUnassignRequest struct {
	ResourceID uint   `json:"resourceId" binding:"required" example:"123"`
	TagIDs     []uint `json:"tagIds" binding:"required"`
}

// TagSearchRequest represents a request to search resources by tags
type TagSearchRequest struct {
	TagIDs    []uint   `form:"tagIds" binding:"omitempty"`
	TagNames  []string `form:"tagNames" binding:"omitempty"`
	Operation string   `form:"operation" binding:"omitempty,oneof=AND OR" example:"AND"`
	common.PaginationRequest
	common.SortRequest
}

// TaggingListRequest represents a request to get resource tags
type TaggingListRequest struct {
	ResourceID uint `form:"resourceId" binding:"required" example:"123"`
}

// TagNameSearchRequest represents a request to search tags by name
type TagNameSearchRequest struct {
	Q string `form:"q" binding:"required" example:"prod"`
}
