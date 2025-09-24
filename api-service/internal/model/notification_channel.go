package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// NotificationChannelConfig represents notification channel configuration
type NotificationChannelConfig struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	Code          string            `json:"code" gorm:"uniqueIndex;size:64;not null;comment:Channel unique code"`
	Name          string            `json:"name" gorm:"size:128;not null;comment:Channel display name"`
	Description   *string           `json:"description" gorm:"size:512;comment:Channel description"`
	ChannelType   string            `json:"channel_type" gorm:"size:16;not null;comment:Channel type (EMAIL, WEBHOOK)"`
	ChannelConfig JSONChannelConfig `json:"channel_config" gorm:"type:json;comment:Channel configuration data"`
	OwnerID       uint              `json:"owner_id" gorm:"not null;comment:Channel owner user ID"`
	Status        int8              `json:"status" gorm:"default:1;comment:Channel status (-1:deleted, 0:disabled, 1:enabled)"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
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

// JSONChannelConfig is a wrapper for storing channel config in database
type JSONChannelConfig map[string]interface{}

// Value implements driver.Valuer interface for GORM
func (j JSONChannelConfig) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	data, err := json.Marshal(j)
	return string(data), err
}

// Scan implements sql.Scanner interface for GORM
func (j *JSONChannelConfig) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("cannot scan %T into JSONChannelConfig", value)
	}

	return json.Unmarshal(data, j)
}

// TableName returns the table name for NotificationChannelConfig
func (NotificationChannelConfig) TableName() string {
	return "notification_channels"
}
