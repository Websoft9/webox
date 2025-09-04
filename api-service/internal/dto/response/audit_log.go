package response

import (
	"encoding/json"
	"time"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// AuditLogResponse audit log response structure
type AuditLogResponse struct {
	ID             uint                  `json:"id" example:"1"`
	User           *AuditLogUserResponse `json:"user,omitempty"`
	Action         string                `json:"action" example:"CREATE"`
	Module         string                `json:"module" example:"APP"`
	ResourceType   string                `json:"resource_type" example:"APP_INSTANCE"`
	ResourceID     *uint                 `json:"resource_id,omitempty" example:"123"`
	ResourceName   string                `json:"resource_name" example:"我的博客"`
	Description    string                `json:"description" example:"创建应用实例：我的博客"`
	IPAddress      string                `json:"ip_address" example:"192.168.1.100"`
	UserAgent      string                `json:"user_agent" example:"Mozilla/5.0"`
	RequestMethod  string                `json:"request_method" example:"POST"`
	RequestURL     string                `json:"request_url" example:"/api/v1/applications"`
	RequestData    interface{}           `json:"request_data,omitempty"`
	ResponseStatus *int                  `json:"response_status,omitempty" example:"200"`
	ResponseTime   *int                  `json:"response_time,omitempty" example:"1500"`
	Success        bool                  `json:"success" example:"true"`
	ErrorMessage   string                `json:"error_message,omitempty"`
	CreatedAt      time.Time             `json:"created_at" example:"2025-07-15T10:30:00Z"`
}

type AuditLogUserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
	Nickname string `json:"nickname" example:"管理员"`
}

// UserBasicResponse basic user information response (used in audit logs)
type UserBasicResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
	Nickname string `json:"nickname" example:"管理员"`
	Email    string `json:"email" example:"admin@example.com"`
}

// AuditLogListResponse audit log list response
type AuditLogListResponse struct {
	Page       int                `json:"page"`        // 当前页码
	PageSize   int                `json:"page_size"`   // 每页数量
	Total      int64              `json:"total"`       // 总记录数
	TotalPages int                `json:"total_pages"` // 总页数
	Items      []AuditLogResponse `json:"items"`       // 审计日志列表
}

// AuditLogStatisticsResponse audit log statistics response
type AuditLogStatisticsResponse struct {
	TotalOperations   int64              `json:"total_operations" example:"1250"`
	SuccessOperations int64              `json:"success_operations" example:"1200"`
	FailedOperations  int64              `json:"failed_operations" example:"50"`
	SuccessRate       float64            `json:"success_rate" example:"96.0"`
	TopUsers          []UserStatItem     `json:"top_users"`
	TopActions        []ActionStatItem   `json:"top_actions"`
	Timeline          []TimelineStatItem `json:"timeline"`
}

// OperationStatItem operation statistics item
type OperationStatItem struct {
	OperationType string `json:"operation_type" example:"查询"`
	Count         int64  `json:"count" example:"500"`
}

// UserStatItem user statistics item
type UserStatItem struct {
	UserID         uint   `json:"user_id" example:"1"`
	Username       string `json:"username" example:"admin"`
	OperationCount int64  `json:"operation_count" example:"500"`
}

// ActionStatItem action statistics item
type ActionStatItem struct {
	Action string `json:"action" example:"READ"`
	Count  int64  `json:"count" example:"800"`
}

// TimelineStatItem timeline statistics item
type TimelineStatItem struct {
	Date  string `json:"date" example:"2025-01-22"`
	Count int64  `json:"count" example:"150"`
}

// HourlyStatItem hourly statistics item
type HourlyStatItem struct {
	Hour  int   `json:"hour" example:"14"`
	Count int64 `json:"count" example:"50"`
}

// StatusStatItem status code statistics item
type StatusStatItem struct {
	StatusCode int   `json:"status_code" example:"200"`
	Count      int64 `json:"count" example:"800"`
}

// ModuleStatItem module statistics item
type ModuleStatItem struct {
	Module string `json:"module" example:"用户管理"`
	Count  int64  `json:"count" example:"300"`
}

// FromAuditLog converts audit log model to response structure
func (r *AuditLogResponse) FromAuditLog(audit *model.AuditLog) {
	r.ID = audit.ID
	r.Action = audit.Action
	r.Module = audit.Module
	r.ResourceType = audit.ResourceType
	r.ResourceID = audit.ResourceID
	r.ResourceName = audit.ResourceName
	r.Description = audit.Description
	r.IPAddress = audit.IPAddress
	r.UserAgent = audit.UserAgent
	r.RequestMethod = audit.RequestMethod
	r.RequestURL = audit.RequestURL
	r.ResponseStatus = audit.ResponseStatus
	r.ResponseTime = audit.ResponseTime
	r.Success = audit.Success
	r.ErrorMessage = audit.ErrorMessage
	r.CreatedAt = audit.CreatedAt

	if audit.RequestParams != "" {
		var requestData interface{}
		if err := json.Unmarshal([]byte(audit.RequestParams), &requestData); err == nil {
			r.RequestData = requestData
		} else {
			r.RequestData = audit.RequestParams
		}
	}

	// Set user information from audit log fields directly
	if audit.UserID != nil {
		r.User = &AuditLogUserResponse{
			ID:       *audit.UserID,
			Username: audit.Username, // May be empty, will be filled by service layer
		}
	}
}

// FromCreateRequest converts create request to response structure
func (r *AuditLogResponse) FromCreateRequest(req *request.CreateAuditLogRequest) {
	if req.UserID != nil {
		r.User = &AuditLogUserResponse{
			ID:       *req.UserID,
			Username: req.Username,
		}
	}

	r.Action = req.Action
	r.Module = req.Module
	r.ResourceType = req.ResourceType
	r.ResourceID = req.ResourceID
	r.ResourceName = req.ResourceName
	r.Description = req.Description
	r.IPAddress = req.IPAddress
	r.UserAgent = req.UserAgent
	r.RequestMethod = req.RequestMethod
	r.RequestURL = req.RequestURL
	r.ResponseStatus = req.ResponseStatus
	r.ResponseTime = req.ResponseTime
	r.Success = req.Success
	r.ErrorMessage = req.ErrorMessage

	// Handle request parameters
	if req.RequestParams != "" {
		var requestData interface{}
		if err := json.Unmarshal([]byte(req.RequestParams), &requestData); err == nil {
			r.RequestData = requestData
		} else {
			r.RequestData = req.RequestParams
		}
	}
}
