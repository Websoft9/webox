package request

import (
	"api-service/internal/dto/common"
	"api-service/internal/model"
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
	common.PaginationRequest
	RuleType   model.AlertRuleType `form:"rule_type" binding:"omitempty,oneof=METRIC LOG EVENT"`
	TargetType model.TargetType    `form:"target_type" binding:"omitempty,oneof=SERVER APP_INSTANCE WORKFLOW"`
	IsEnabled  *bool               `form:"is_enabled"`
	Keyword    string              `form:"keyword"`
	common.SortRequest
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
	common.PaginationRequest
	Status      string `form:"status" binding:"omitempty"`
	Severity    string `form:"severity" binding:"omitempty"`
	AlertRuleID *uint  `form:"alert_rule_id" binding:"omitempty"`
	common.TimeRangeRequest
	common.SortRequest
}

// AlertAcknowledgeRequest represents a request to acknowledge an alert record
type AlertAcknowledgeRequest struct {
	Note string `json:"note" binding:"omitempty"`
}

// AlertResolveRequest represents a request to resolve an alert record
type AlertResolveRequest struct {
	ResolutionNote string `json:"resolution_note" binding:"omitempty"`
}
