package response

import (
	"encoding/json"
	"time"

	"api-service/internal/model"
)

// AuditLogResponse audit log response structure
type AuditLogResponse struct {
	ID             uint                  `json:"id" example:"1"`
	User           *AuditLogUserResponse `json:"user,omitempty"`
	Action         string                `json:"action" example:"CREATE"`
	Module         string                `json:"module" example:"APP"`
	Description    string                `json:"description" example:"Created a new application instance"`
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
}

// FromAuditLog converts audit log model to response structure
func (r *AuditLogResponse) FromAuditLog(audit *model.AuditLog) {
	r.ID = audit.ID
	r.Action = audit.Action
	r.Module = audit.Module
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
