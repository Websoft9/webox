package request

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
	ResourceID uint64   `json:"resourceId" binding:"required" example:"123"`
	TagIDs     []uint64 `json:"tagIds"`
	TagNames   []string `json:"tagNames"`
	ReplaceAll bool     `json:"replaceAll" example:"false"`
}

// TagUnassignRequest represents a request to unassign tags from a resource
type TagUnassignRequest struct {
	ResourceID uint64   `json:"resourceId" binding:"required" example:"123"`
	TagIDs     []uint64 `json:"tagIds" binding:"required"`
}

// TagSearchRequest represents a request to search resources by tags
type TagSearchRequest struct {
	TagIDs    []uint64 `form:"tagIds" binding:"omitempty"`
	TagNames  []string `form:"tagNames" binding:"omitempty"`
	Operation string   `form:"operation" binding:"omitempty,oneof=AND OR" example:"AND"`
	Page      int      `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize  int      `form:"pageSize" binding:"omitempty,min=1,max=100" example:"20"`
}

// TaggingListRequest represents a request to get resource tags
type TaggingListRequest struct {
	ResourceID uint64 `form:"resourceId" binding:"required" example:"123"`
}

// TagNameSearchRequest represents a request to search tags by name
type TagNameSearchRequest struct {
	Q string `form:"q" binding:"required" example:"prod"`
}
