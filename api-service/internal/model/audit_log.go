package model

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	// Path parts threshold for user agent masking
	pathPartsThreshold = 2
)

// AuditLog audit log model
type AuditLog struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	UserID         *uint     `json:"user_id" gorm:"index"`
	Username       string    `json:"username" gorm:"size:64;index"`
	Action         string    `json:"action" gorm:"size:32;not null;index"`
	Module         string    `json:"module" gorm:"size:32;not null;index"`
	ResourceType   string    `json:"resource_type" gorm:"size:32"`
	ResourceID     *uint     `json:"resource_id" gorm:"index"`
	ResourceName   string    `json:"resource_name" gorm:"size:64"`
	Description    string    `json:"description"`
	IPAddress      string    `json:"ip_address" gorm:"size:45"`
	UserAgent      string    `json:"user_agent" gorm:"size:255"`
	RequestMethod  string    `json:"request_method" gorm:"size:10"`
	RequestURL     string    `json:"request_url" gorm:"size:255"`
	RequestParams  string    `json:"request_params"`
	ResponseStatus *int      `json:"response_status"`
	ResponseTime   *int      `json:"response_time"`
	Success        bool      `json:"success" gorm:"not null;index"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// GetFormattedRequestParams formats request parameters with sensitive data masking
func (a *AuditLog) GetFormattedRequestParams() interface{} {
	if a.RequestParams == "" {
		return nil
	}

	var params interface{}
	if err := json.Unmarshal([]byte(a.RequestParams), &params); err != nil {
		return a.RequestParams
	}

	// Mask sensitive data
	return a.sanitizeData(params)
}

// SanitizeForExport sanitizes sensitive data for export
func (a *AuditLog) SanitizeForExport() *AuditLog {
	sanitized := *a

	// Sanitize request parameters
	if sanitized.RequestParams != "" {
		var params interface{}
		if err := json.Unmarshal([]byte(sanitized.RequestParams), &params); err == nil {
			sanitizedParams := a.sanitizeData(params)
			if data, err := json.Marshal(sanitizedParams); err == nil {
				sanitized.RequestParams = string(data)
			}
		}
	}

	// Sanitize user agent information (keep browser info, hide detailed version)
	if sanitized.UserAgent != "" {
		parts := strings.Fields(sanitized.UserAgent)
		if len(parts) > pathPartsThreshold {
			sanitized.UserAgent = parts[0] + " " + parts[1] + " [details masked]"
		}
	}

	return &sanitized
}

// sanitizeData recursively sanitizes sensitive data
func (a *AuditLog) sanitizeData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for key, value := range v {
			lowerKey := strings.ToLower(key)
			// Mask sensitive fields
			if strings.Contains(lowerKey, "password") ||
				strings.Contains(lowerKey, "secret") ||
				strings.Contains(lowerKey, "token") ||
				strings.Contains(lowerKey, "key") ||
				strings.Contains(lowerKey, "auth") {
				sanitized[key] = "[masked]"
			} else {
				sanitized[key] = a.sanitizeData(value)
			}
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(v))
		for i, item := range v {
			sanitized[i] = a.sanitizeData(item)
		}
		return sanitized
	default:
		return v
	}
}
