package dto

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
		p.PageSize = 10 // 默认每页10条
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // 最大每页100条
	}
	return p.PageSize
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页数量
	Total      int64       `json:"total"`       // 总记录数
	TotalPages int         `json:"total_pages"` // 总页数
	Data       interface{} `json:"data"`        // 数据列表
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64, data interface{}) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
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

	order := "asc"
	if s.Order == "desc" {
		order = "desc"
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
