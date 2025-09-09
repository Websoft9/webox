package request

import (
	"time"
)

// CreateAuditLogRequest request for creating audit log (internal use)
type CreateAuditLogRequest struct {
	UserID         *uint  `json:"user_id"`
	Username       string `json:"username"`
	Action         string `json:"action" validate:"required"`
	Module         string `json:"module" validate:"required,max=32"`
	ResourceType   string `json:"resource_type" validate:"max=32"`
	ResourceID     *uint  `json:"resource_id"`
	ResourceName   string `json:"resource_name" validate:"max=64"`
	Description    string `json:"description"`
	IPAddress      string `json:"ip_address" validate:"max=45"`
	UserAgent      string `json:"user_agent" validate:"max=255"`
	RequestMethod  string `json:"request_method" validate:"max=10"`
	RequestURL     string `json:"request_url" validate:"max=255"`
	RequestParams  string `json:"request_params"`
	ResponseStatus *int   `json:"response_status"`
	ResponseTime   *int   `json:"response_time"`
	Success        bool   `json:"success"`
	ErrorMessage   string `json:"error_message"`
}

// ListAuditLogRequest request for querying audit log list
type ListAuditLogRequest struct {
	Page         int        `form:"page" json:"page" validate:"omitempty,min=1" default:"1"`
	PageSize     int        `form:"page_size" json:"page_size" validate:"omitempty,min=1,max=100" default:"20"`
	UserID       *uint      `form:"user_id" json:"user_id"`
	Action       string     `form:"action" json:"action" validate:"max=32"`
	ResourceType string     `form:"resource_type" json:"resource_type" validate:"max=32"`
	ResourceID   *uint      `form:"resource_id" json:"resource_id"`
	StartTime    *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02 15:04:00"`
	EndTime      *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02 15:04:00"`
	IPAddress    string     `form:"ip_address" json:"ip_address" validate:"max=45"`
}

// AuditLogStatisticsRequest request for audit log statistics
type AuditLogStatisticsRequest struct {
	StartTime *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02 15:04:00"`
	EndTime   *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02 15:04:00"`
	GroupBy   string     `form:"group_by" json:"group_by" validate:"omitempty,oneof=hour day week month" default:"day"`
}

// ExportAuditLogRequest request for exporting audit logs
type ExportAuditLogRequest struct {
	Format    string     `form:"format" json:"format" validate:"omitempty,oneof=csv excel json" default:"csv"`
	StartTime *time.Time `form:"start_time" json:"start_time" time_format:"2006-01-02 15:04:00"`
	EndTime   *time.Time `form:"end_time" json:"end_time" time_format:"2006-01-02 15:04:00"`
	UserID    *uint      `form:"user_id" json:"user_id"`
}

// StatisticsFilter statistics query filter conditions
type StatisticsFilter struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	GroupBy   string     `json:"group_by"`
}

// AuditLogFilter audit log query filter conditions
type AuditLogFilter struct {
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	UserID       *uint      `json:"user_id"`
	Username     string     `json:"username"`
	Action       string     `json:"action"`
	Module       string     `json:"module"`
	ResourceType string     `json:"resource_type"`
	ResourceID   *uint      `json:"resource_id"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	IPAddress    string     `json:"ip_address"`
	Success      *bool      `json:"success"`
}
