package model

import "time"

// NotificationRecord represents a notification record in the database
type NotificationRecord struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement;comment:Record ID"`
	TemplateID    *uint      `json:"template_id" gorm:"index;comment:Template ID"`
	ChannelType   string     `json:"channel_type" gorm:"type:varchar(20);not null;comment:Notification channel"`
	Recipient     string     `json:"recipient" gorm:"type:varchar(255);not null;comment:Recipient"`
	Subject       *string    `json:"subject" gorm:"type:varchar(255);comment:Notification subject"`
	Content       string     `json:"content" gorm:"type:text;not null;comment:Notification content"`
	Status        string     `json:"status" gorm:"type:varchar(20);default:'PENDING';comment:Sending status"`
	SentAt        *time.Time `json:"sent_at" gorm:"comment:Sent time"`
	ErrorMsg      *string    `json:"error_msg" gorm:"type:text;comment:Error message"`
	RetryCount    int        `json:"retry_count" gorm:"default:0;comment:Retry count"`
	ReferenceID   *string    `json:"reference_id" gorm:"type:varchar(64);comment:Reference record ID"`
	ReferenceType *string    `json:"reference_type" gorm:"type:varchar(32);comment:Reference type (table_name)"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime;comment:Created at"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:Updated at"`
}

// TableName returns the table name for NotificationRecord
func (NotificationRecord) TableName() string {
	return "notification_records"
}
