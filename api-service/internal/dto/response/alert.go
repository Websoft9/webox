package response

import (
	"api-service/internal/model"
	"time"
)

// AlertRuleResponse defines the response data for a single alert rule.
type AlertRuleResponse struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	RuleType             model.AlertRuleType `json:"rule_type"`
	TargetType           model.TargetType    `json:"target_type"`
	TargetID             *uint               `json:"target_id,omitempty"`
	MetricName           string              `json:"metric_name,omitempty"`
	ConditionExpression  string              `json:"condition_expression"`
	NotificationChannels string              `json:"notification_channels,omitempty"`
	IsEnabled            bool                `json:"is_enabled"`
	OwnerID              uint                `json:"owner_id"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

// AlertRuleListResponse defines the response data for a list of alert rules.
type AlertRuleListResponse struct {
	Total int64               `json:"total"`
	Items []AlertRuleResponse `json:"items"`
}

// AlertRecordResponse represents the response data for a single alert record
type AlertRecordResponse struct {
	ID               uint       `json:"id"`
	AlertRuleID      uint       `json:"alert_rule_id"`
	AlertID          string     `json:"alert_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Severity         string     `json:"severity"`
	FiredAt          time.Time  `json:"fired_at"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	AcknowledgedAt   *time.Time `json:"acknowledged_at,omitempty"`
	AcknowledgedBy   *uint      `json:"acknowledged_by,omitempty"`
	AcknowledgeNote  string     `json:"acknowledge_note,omitempty"`
	ResolutionNote   string     `json:"resolution_note,omitempty"`
	NotificationSent bool       `json:"notification_sent"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// AlertRecordListResponse represents a paginated list of alert records
type AlertRecordListResponse struct {
	Records  []AlertRecordResponse `json:"records"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}
