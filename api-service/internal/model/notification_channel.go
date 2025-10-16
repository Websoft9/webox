package model

import "time"

// NotificationChannelConfig represents notification channel configuration
type NotificationChannelConfig struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:Channel ID"`
	Code          string    `json:"code" gorm:"uniqueIndex;size:64;not null;comment:Channel unique code"`
	Name          string    `json:"name" gorm:"size:128;not null;comment:Channel display name"`
	Description   *string   `json:"description" gorm:"size:512;comment:Channel description"`
	ChannelType   string    `json:"channel_type" gorm:"size:16;not null;comment:Channel type (EMAIL, WEBHOOK)"`
	ChannelConfig JSON      `json:"channel_config" gorm:"type:json;comment:Channel configuration data"`
	OwnerID       uint      `json:"owner_id" gorm:"not null;comment:Channel owner user ID"`
	Status        int8      `json:"status" gorm:"default:1;comment:Channel status (0:disabled, 1:enabled)"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
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
