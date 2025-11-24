package model

import (
	"time"
)

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Name         string    `gorm:"size:64;not null;comment:Template name" json:"name"`
	TemplateType string    `gorm:"size:20;not null;column:template_type;index:idx_template_type;comment:Notification type" json:"template_type"` // EMAIL, WEBHOOK, INTERNAL
	Subject      *string   `gorm:"size:255;comment:Notification subject" json:"subject"`
	Content      string    `gorm:"type:text;not null;comment:Notification content template" json:"content"`
	IsSystem     int       `gorm:"type:tinyint(1);default:0;column:is_system;comment:Whether system template" json:"is_system"`                          // 0-user template, 1-system template
	Status       int       `gorm:"type:tinyint(1);default:1;index:idx_notification_template_status;comment:Status: 0-disabled, 1-enabled" json:"status"` // 0-disabled, 1-enabled
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName returns the table name for NotificationTemplate
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
