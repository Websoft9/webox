package request

import (
	"api-service/internal/dto/common"
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
	common.PaginationRequest
	UserID       *uint  `form:"user_id" json:"user_id"`
	Action       string `form:"action" json:"action" validate:"max=32"`
	ResourceType string `form:"resource_type" json:"resource_type" validate:"max=32"`
	ResourceID   *uint  `form:"resource_id" json:"resource_id"`
	common.TimeRangeRequest
	IPAddress string `form:"ip_address" json:"ip_address" validate:"max=45"`
	Success   *bool  `form:"success" json:"success"`
}

// ExportAuditLogRequest request for exporting audit logs
type ExportAuditLogRequest struct {
	Format string `form:"format" json:"format" validate:"omitempty,oneof=csv excel json" default:"csv"`
	common.TimeRangeRequest
	UserID *uint `form:"user_id" json:"user_id"`
}
