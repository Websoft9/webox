package dto

// 分页相关常量
const (
	DefaultPageSize int = 10  // 默认每页数量
	MaxPageSize     int = 100 // 最大每页数量
)

// 排序相关常量
const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// PaginationRequest 分页请求结构
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`                   // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量，最大100
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetPageSize()
}

// GetPageSize 获取每页数量
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p.PageSize
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页数量
	Total      int64       `json:"total"`       // 总记录数
	TotalPages int         `json:"total_pages"` // 总页数
	Items      interface{} `json:"items"`       // 数据列表
}

// ListResponse 通用列表响应结构（包含分页信息和数据列表）
type ListResponse struct {
	Items      interface{}         `json:"items"`      // 数据列表
	Pagination *PaginationResponse `json:"pagination"` // 分页信息
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64, items interface{}) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
	}
}

// NewListResponse 创建列表响应（包含分页信息）
func NewListResponse(page, pageSize int, total int64, items interface{}) *ListResponse {
	pagination := &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}
	if pagination.TotalPages == 0 {
		pagination.TotalPages = 1
	}

	return &ListResponse{
		Items:      items,
		Pagination: pagination,
	}
}

// SortRequest 排序请求结构
type SortRequest struct {
	Field string `form:"sort_field" json:"sort_field"` // 排序字段
	Order string `form:"sort_order" json:"sort_order"` // 排序方向: asc, desc
}

// GetSortOrder 获取排序SQL片段
func (s *SortRequest) GetSortOrder() string {
	if s.Field == "" {
		return "id desc" // 默认按ID降序
	}

	order := SortOrderAsc
	if s.Order == SortOrderDesc {
		order = SortOrderDesc
	}

	return s.Field + " " + order
}

// SearchRequest 搜索请求结构
type SearchRequest struct {
	Keyword string `form:"keyword" json:"keyword"` // 搜索关键词
}

// BaseListRequest 基础列表请求（包含分页、排序、搜索）
type BaseListRequest struct {
	PaginationRequest
	SortRequest
	SearchRequest
}
