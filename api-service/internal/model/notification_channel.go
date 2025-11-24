package model

import "time"

// NotificationChannelConfig represents notification channel configuration
type NotificationChannelConfig struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	Code          string    `json:"code" gorm:"uniqueIndex:uk_nfc_code;size:64;not null;comment:Channel code (unique identifier)"`
	Name          string    `json:"name" gorm:"size:128;not null;comment:Channel name"`
	Description   *string   `json:"description" gorm:"type:text;comment:Channel description"`
	ChannelType   string    `json:"channel_type" gorm:"type:varchar(20);not null;index:idx_channel_type;comment:Channel type"`
	ChannelConfig JSON      `json:"channel_config" gorm:"type:json;not null;comment:Channel configuration (JSON format)"`
	OwnerID       uint      `json:"owner_id" gorm:"type:integer;not null;index:idx_nfc_owner_id;comment:Owner user ID"`
	Status        int8      `json:"status" gorm:"type:tinyint(1);default:1;index:idx_nfc_status;comment:Status (0-disabled, 1-enabled)"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// EmailConfig represents email channel configuration
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPSecurity string `json:"smtp_security"`
	SMTPTimeout  int    `json:"smtp_timeout"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SenderEmail  string `json:"sender_email"`
	SenderName   string `json:"sender_name"`
	Encoding     string `json:"encoding"`
	RateLimit    int    `json:"rate_limit"`
	RetryCount   int    `json:"retry_count"`
	QuietHours   string `json:"quiet_hours"`
}

// WebhookConfig represents webhook channel configuration
type WebhookConfig struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Secret     string            `json:"secret"`
	Timeout    int               `json:"timeout"`
	Encoding   string            `json:"encoding"`
	RateLimit  int               `json:"rate_limit"`
	RetryCount int               `json:"retry_count"`
	QuietHours string            `json:"quiet_hours"`
}

// TableName returns the table name for NotificationChannelConfig
func (NotificationChannelConfig) TableName() string {
	return "notification_channels"
}
