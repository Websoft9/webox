package common

import (
	"time"

	"api-service/internal/constants"
	"api-service/internal/validator"
)

// pagination related constants
const (
	DEFAULT_PAGE_SIZE int = 20  // default page size
	MAX_PAGE_SIZE     int = 100 // maximum page size
)

// sorting related constants
const (
	SORT_ORDER_ASC  = "ASC"
	SORT_ORDER_DESC = "DESC"
)

// PaginationRequest pagination request structure
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`                   // page number, starting from 1
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"` // page size, maximum 100
}

// GetPage retrieves the current page number, defaulting to 1 if not set or invalid
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

// GetOffset retrieves the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// GetPageSize retrieves the page size
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = DEFAULT_PAGE_SIZE
	}
	if p.PageSize > MAX_PAGE_SIZE {
		p.PageSize = MAX_PAGE_SIZE
	}
	return p.PageSize
}

// PaginationResponse pagination response structure
type PaginationResponse struct {
	Page       int         `json:"page"`        // current page number
	PageSize   int         `json:"page_size"`   // page size
	Total      int64       `json:"total"`       // total record count
	TotalPages int         `json:"total_pages"` // total pages
	Items      interface{} `json:"items"`       // data list
}

// NewPaginationResponse creates a new pagination response
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

// SortRequest sorting request structure
type SortRequest struct {
	Field string `form:"sort_field" json:"sort_field" validate:"omitempty"`                // sorting field
	Order string `form:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"` // sorting direction: asc, desc
}

// GetSortOrder retrieves the sorting SQL fragment
func (s *SortRequest) GetSortOrder() string {
	if s.Field == "" {
		return "created_at desc" // default to descending by creation time
	}

	order := SORT_ORDER_ASC
	if s.Order == SORT_ORDER_DESC {
		order = SORT_ORDER_DESC
	}

	return s.Field + " " + order
}

// SearchRequest search request structure
type SearchRequest struct {
	Keyword string `form:"keyword" json:"keyword" validate:"omitempty"` // search keyword
}

// TimeRangeRequest represents a time range request with validation
type TimeRangeRequest struct {
	StartTime string `form:"start_time" json:"start_time" validate:"omitempty,time_range" example:"2023-01-01T00:00:00Z"`
	EndTime   string `form:"end_time" json:"end_time" validate:"omitempty,time_range" example:"2023-01-31T23:59:59Z"`
}

// GetParsedTimeRange parses and returns the time range with validation and adjustment
func (t *TimeRangeRequest) GetParsedTimeRange() (startTime, endTime time.Time, err error) {
	// Parse start time or use default (24 hours ago)
	if t.StartTime != "" {
		startTime, err = time.Parse(constants.DefaultTimeFormat, t.StartTime)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		startTime = time.Now().AddDate(0, 0, -1) // Default: 24 hours ago
	}

	// Parse end time or use current time
	if t.EndTime != "" {
		endTime, err = time.Parse(constants.DefaultTimeFormat, t.EndTime)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		endTime = time.Now()
	}

	// Validate and adjust time range (limit to configured maximum duration)
	maxDuration := validator.GetMaxTimeRangeDuration()
	if endTime.Sub(startTime) > maxDuration {
		endTime = startTime.Add(maxDuration)
	}

	return startTime, endTime, nil
}

// BaseListRequest base list request structure (includes pagination, sorting, and search)
type BaseListRequest struct {
	PaginationRequest
	SortRequest
	SearchRequest
}
