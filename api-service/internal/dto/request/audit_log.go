package request

import (
	"time"
)

// CreateAuditLogRequest request for creating audit log (internal use)
type CreateAuditLogRequest struct {
	UserID         *uint  `json:"user_id"`
	Username       string `json:"username"`
	Action         string `json:"action" binding:"required"`
	Module         string `json:"module" binding:"required,max=32"`
	ResourceType   string `json:"resource_type" binding:"max=32"`
	ResourceID     *uint  `json:"resource_id"`
	ResourceName   string `json:"resource_name" binding:"max=64"`
	Description    string `json:"description"`
	IPAddress      string `json:"ip_address" binding:"max=45"`
	UserAgent      string `json:"user_agent" binding:"max=255"`
	RequestMethod  string `json:"request_method" binding:"max=10"`
	RequestURL     string `json:"request_url" binding:"max=255"`
	RequestParams  string `json:"request_params"`
	ResponseStatus *int   `json:"response_status"`
	ResponseTime   *int   `json:"response_time"`
	Success        bool   `json:"success"`
	ErrorMessage   string `json:"error_message"`
}

// ListAuditLogRequest request for querying audit log list
type ListAuditLogRequest struct {
	Page         int        `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize     int        `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	UserID       *uint      `form:"user_id" json:"user_id"`
	Action       string     `form:"action" json:"action" binding:"max=32"`
	Module       string     `form:"module" json:"module" binding:"max=32"`
	ResourceType string     `form:"resource_type" json:"resource_type" binding:"max=32"`
	ResourceID   *uint      `form:"resource_id" json:"resource_id"`
	StartTime    *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime      *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	IPAddress    string     `form:"ip_address" json:"ip_address" binding:"max=45"`
	Success      *bool      `form:"success" json:"success"`
}

// AuditLogStatisticsRequest request for audit log statistics
type AuditLogStatisticsRequest struct {
	StartTime *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	GroupBy   string     `form:"group_by" json:"group_by" binding:"omitempty,oneof=hour day week month"`
}

// ExportAuditLogRequest request for exporting audit logs
type ExportAuditLogRequest struct {
	Format       string     `form:"format" json:"format" binding:"omitempty,oneof=csv excel json"`
	StartTime    *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime      *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02T15:04:05Z07:00"`
	UserID       *uint      `form:"user_id" json:"user_id"`
	Action       string     `form:"action" json:"action" binding:"max=32"`
	Module       string     `form:"module" json:"module" binding:"max=32"`
	ResourceType string     `form:"resource_type" json:"resource_type" binding:"max=32"`
}

// SetDefaults sets default values for list request
func (r *ListAuditLogRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
}

// SetDefaults sets default values for statistics request
func (r *AuditLogStatisticsRequest) SetDefaults() {
	if r.GroupBy == "" {
		r.GroupBy = "day"
	}
}

// SetDefaults sets default values for export request
func (r *ExportAuditLogRequest) SetDefaults() {
	if r.Format == "" {
		r.Format = "csv"
	}
}
