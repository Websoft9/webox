package model

import "time"

// NotificationRecord represents a notification record in the database
type NotificationRecord struct {
	ID            uint       `json:"id" gorm:"primarykey"`
	TemplateID    *uint      `json:"template_id" gorm:"type:integer;index:idx_notification_template_id;comment:Template ID"`
	ChannelType   string     `json:"channel_type" gorm:"type:varchar(20);not null;index:idx_notification_channel_type;comment:Notification type"`
	Recipient     string     `json:"recipient" gorm:"type:varchar(255);not null;comment:Recipient"`
	Subject       *string    `json:"subject" gorm:"type:varchar(255);comment:Notification subject"`
	Content       string     `json:"content" gorm:"type:text;not null;comment:Notification content"`
	Status        string     `json:"status" gorm:"type:varchar(20);default:'PENDING';index:idx_notification_status;comment:Send status"`
	SentAt        *time.Time `json:"sent_at" gorm:"type:datetime;serializer:datetime;comment:Send time"`
	ErrorMsg      *string    `json:"error_msg" gorm:"type:text;comment:Error message"`
	RetryCount    int        `json:"retry_count" gorm:"type:int;default:0;comment:Retry count"`
	ReferenceID   *string    `json:"reference_id" gorm:"type:varchar(64);index:idx_reference_id;comment:Reference ID"`
	ReferenceType *string    `json:"reference_type" gorm:"type:varchar(32);index:idx_reference_type;comment:Reference type"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName returns the table name for NotificationRecord
func (NotificationRecord) TableName() string {
	return "notification_records"
}
