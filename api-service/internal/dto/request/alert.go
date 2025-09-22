package request

import (
	"api-service/internal/model"
	"time"
)

// AlertRuleCreateRequest defines the request data for creating an alert rule.
type AlertRuleCreateRequest struct {
	Name                 string              `json:"name" binding:"required,max=64"`
	RuleType             model.AlertRuleType `json:"rule_type" binding:"required,oneof=METRIC LOG EVENT"`
	TargetType           model.TargetType    `json:"target_type" binding:"required,oneof=SERVER APP_INSTANCE WORKFLOW"`
	TargetID             *uint               `json:"target_id"`
	MetricName           string              `json:"metric_name"`
	ConditionExpression  string              `json:"condition_expression" binding:"required"`
	NotificationChannels string              `json:"notification_channels"`
	IsEnabled            *bool               `json:"is_enabled"`
}

// AlertRuleQueryRequest defines the request data for querying alert rules.
type AlertRuleQueryRequest struct {
	Page       int                 `form:"page,default=1" binding:"min=1"`
	PageSize   int                 `form:"page_size,default=20" binding:"min=5,max=100"`
	RuleType   model.AlertRuleType `form:"rule_type" binding:"omitempty,oneof=METRIC LOG EVENT"`
	TargetType model.TargetType    `form:"target_type" binding:"omitempty,oneof=SERVER APP_INSTANCE WORKFLOW"`
	IsEnabled  *bool               `form:"is_enabled"`
	Keyword    string              `form:"keyword"`
}

// AlertRuleUpdateRequest defines the request data for updating an alert rule.
type AlertRuleUpdateRequest struct {
	Name                 *string `json:"name" binding:"omitempty,max=64"`
	ConditionExpression  *string `json:"condition_expression"`
	NotificationChannels *string `json:"notification_channels"`
	IsEnabled            *bool   `json:"is_enabled"`
}

// AlertRecordQueryRequest represents a request for querying alert records
type AlertRecordQueryRequest struct {
	Page        int       `form:"page" binding:"omitempty,min=1"`
	PageSize    int       `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status      string    `form:"status" binding:"omitempty"`
	Severity    string    `form:"severity" binding:"omitempty"`
	StartTime   time.Time `form:"start_time" binding:"omitempty"`
	EndTime     time.Time `form:"end_time" binding:"omitempty"`
	AlertRuleID *uint     `form:"alert_rule_id" binding:"omitempty"`
}

// AlertAcknowledgeRequest represents a request to acknowledge an alert record
type AlertAcknowledgeRequest struct {
	Note string `json:"note" binding:"omitempty"`
}

// AlertResolveRequest represents a request to resolve an alert record
type AlertResolveRequest struct {
	ResolutionNote string `json:"resolution_note" binding:"omitempty"`
}
