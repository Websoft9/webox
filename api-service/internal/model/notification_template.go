package model

import (
	"time"
)

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	TemplateType string    `gorm:"size:20;not null;column:template_type" json:"template_type"` // EMAIL, WEBHOOK, INTERNAL
	Subject      *string   `gorm:"size:255" json:"subject"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	IsSystem     int       `gorm:"default:0;column:is_system" json:"is_system"` // 0-user template, 1-system template
	Status       int       `gorm:"default:1" json:"status"`                     // 0-disabled, 1-enabled
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName returns the table name for NotificationTemplate
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
