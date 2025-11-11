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
	ID                   uint          `gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	Name                 string        `gorm:"size:64;not null;comment:Rule name"`
	RuleType             AlertRuleType `gorm:"type:varchar(20);not null;index:idx_rule_type;comment:Rule type"`
	TargetType           TargetType    `gorm:"type:varchar(20);not null;index:idx_target_type;comment:Target type"`
	TargetID             *uint         `gorm:"type:bigint unsigned;index:idx_target_id;comment:Target ID"`
	MetricName           string        `gorm:"size:64;comment:Metric name"`
	ConditionExpression  string        `gorm:"type:text;not null;comment:Condition expression"`
	NotificationChannels string        `gorm:"type:json;comment:Notification channels (JSON format)"`
	IsEnabled            bool          `gorm:"type:tinyint(1);default:1;index:idx_is_enabled;comment:Whether enabled"`
	OwnerID              uint          `gorm:"type:bigint unsigned;not null;index:idx_alert_owner_id;comment:Owner ID"`
	CreatedAt            time.Time     `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt            time.Time     `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName specifies the table name.
func (AlertRule) TableName() string {
	return "alert_rules"
}

// AlertRecord represents an alert record in the system
type AlertRecord struct {
	ID                   uint       `json:"id" gorm:"primarykey;type:bigint unsigned"`
	AlertRuleID          uint       `json:"alert_rule_id" gorm:"type:bigint unsigned;not null;index:idx_alert_records_rule_status;comment:Alert rule ID"`
	AlertID              string     `json:"alert_id" gorm:"not null;uniqueIndex;size:64;comment:Alert ID"`
	Title                string     `json:"title" gorm:"not null;size:255;comment:Alert title"`
	Description          string     `json:"description" gorm:"type:text;comment:Alert description"`
	Status               string     `json:"status" gorm:"type:varchar(20);default:'FIRING';index:idx_alert_records_rule_status;comment:Status"`
	FiredAt              time.Time  `json:"fired_at" gorm:"type:datetime;serializer:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_fired_at;comment:Fired time"`
	ResolvedAt           *time.Time `json:"resolved_at" gorm:"type:datetime;serializer:datetime;comment:Resolved time"`
	AcknowledgedAt       *time.Time `json:"acknowledged_at" gorm:"type:datetime;serializer:datetime;comment:Acknowledged time"`
	AcknowledgedBy       *uint      `json:"acknowledged_by" gorm:"type:bigint unsigned;index:idx_acknowledged_by;comment:Acknowledged by ID"`
	AcknowledgeNote      string     `json:"acknowledge_note" gorm:"type:text"`
	ResolutionNote       string     `json:"resolution_note" gorm:"type:text;comment:Resolution note"`
	NotificationSent     bool       `json:"notification_sent" gorm:"type:tinyint(1);default:0;comment:Notification sent"`
	NotificationChannels string     `json:"notification_channels" gorm:"type:json;comment:Notification channels (JSON format)"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName specifies the table name for AlertRecord
func (AlertRecord) TableName() string {
	return "alert_records"
}
