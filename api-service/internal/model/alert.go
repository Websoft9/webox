package model

import (
	"encoding/json"
	"time"
)

// AlertRuleType defines the type of alert rule.
type AlertRuleType string

const (
	AlertRuleTypeMetric AlertRuleType = "METRIC"
	AlertRuleTypeLog    AlertRuleType = "LOG"
	AlertRuleTypeEvent  AlertRuleType = "EVENT"
)

// TargetType defines the type of alert target.
type TargetType string

const (
	TargetTypeServer      TargetType = "SERVER"
	TargetTypeAppInstance TargetType = "APP_INSTANCE"
	TargetTypeWorkflow    TargetType = "WORKFLOW"
)

// NotificationChannel defines the notification channel.
type NotificationChannel struct {
	Type    string          `json:"type"`    // Notification channel type: email, sms, webhook, etc.
	Config  json.RawMessage `json:"config"`  // Notification channel configuration
	Enabled bool            `json:"enabled"` // Whether the channel is enabled
}

// AlertRule is the model for alert rules.
type AlertRule struct {
	ID                   uint          `gorm:"primaryKey;autoIncrement"`
	Name                 string        `gorm:"size:64;not null"`
	RuleType             AlertRuleType `gorm:"type:enum('METRIC','LOG','EVENT');not null"`
	TargetType           TargetType    `gorm:"type:enum('SERVER','APP_INSTANCE','WORKFLOW');not null"`
	TargetID             *uint         `gorm:"index"`
	MetricName           string        `gorm:"size:64"`
	ConditionExpression  string        `gorm:"type:text;not null"`
	NotificationChannels string        `gorm:"type:text"`
	IsEnabled            bool          `gorm:"default:true"`
	OwnerID              uint          `gorm:"not null;index"`
	CreatedAt            time.Time     `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time     `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name.
func (AlertRule) TableName() string {
	return "alert_rules"
}

// AlertRecord represents an alert record in the system
type AlertRecord struct {
	ID                   uint       `json:"id" gorm:"primarykey"`
	AlertRuleID          uint       `json:"alert_rule_id" gorm:"not null"`
	AlertID              string     `json:"alert_id" gorm:"not null;unique;size:64"`
	Title                string     `json:"title" gorm:"not null;size:255"`
	Description          string     `json:"description" gorm:"type:text"`
	Status               string     `json:"status" gorm:"default:'FIRING'"`
	FiredAt              time.Time  `json:"fired_at" gorm:"type:datetime;serializer:datetime;not null"`
	ResolvedAt           *time.Time `json:"resolved_at" gorm:"type:datetime;serializer:datetime"`
	AcknowledgedAt       *time.Time `json:"acknowledged_at" gorm:"type:datetime;serializer:datetime"`
	AcknowledgedBy       *uint      `json:"acknowledged_by"`
	AcknowledgeNote      string     `json:"acknowledge_note" gorm:"type:text"`
	ResolutionNote       string     `json:"resolution_note" gorm:"type:text"`
	NotificationSent     bool       `json:"notification_sent" gorm:"default:false"`
	NotificationChannels string     `json:"notification_channels" gorm:"type:json"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for AlertRecord
func (AlertRecord) TableName() string {
	return "alert_records"
}
