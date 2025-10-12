package common

// 分页相关常量
const (
	DEFAULT_PAGE_SIZE int = 10  // 默认每页数量
	MAX_PAGE_SIZE     int = 100 // 最大每页数量
)

// 排序相关常量
const (
	SORT_ORDER_ASC  = "ASC"
	SORT_ORDER_DESC = "DESC"
)

// PaginationRequest 分页请求结构
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`                   // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量，最大100
}

// GetOffset 获取页码
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

// GetOffset 计算偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// GetPageSize 获取每页数量
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = DEFAULT_PAGE_SIZE
	}
	if p.PageSize > MAX_PAGE_SIZE {
		p.PageSize = MAX_PAGE_SIZE
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

// SortRequest 排序请求结构
type SortRequest struct {
	Field string `form:"sort_field" json:"sort_field" validate:"omitempty"`                // 排序字段
	Order string `form:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"` // 排序方向: asc, desc
}

// GetSortOrder 获取排序SQL片段
func (s *SortRequest) GetSortOrder() string {
	if s.Field == "" {
		return "created_at desc" // 默认按创建时间降序
	}

	order := SORT_ORDER_ASC
	if s.Order == SORT_ORDER_DESC {
		order = SORT_ORDER_DESC
	}

	return s.Field + " " + order
}

// SearchRequest 搜索请求结构
type SearchRequest struct {
	Keyword string `form:"keyword" json:"keyword" validate:"omitempty"` // 搜索关键词
}

// TimeRangeRequest 时间范围请求结构
type TimeRangeRequest struct {
	StartTime string `form:"start_time" validate:"omitempty,rfc3339"`
	EndTime   string `form:"end_time" validate:"omitempty,rfc3339"`
}

// BaseListRequest 基础列表请求（包含分页、排序、搜索）
type BaseListRequest struct {
	PaginationRequest
	SortRequest
	SearchRequest
}
