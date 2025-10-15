package request

import "api-service/internal/dto/common"

// GetNotificationChannelListRequest represents the request for getting notification channel list
type GetNotificationChannelListRequest struct {
	common.BaseListRequest
	ChannelType string `form:"channel_type" binding:"omitempty,oneof=EMAIL WEBHOOK" json:"channel_type"`
	Status      *int8  `form:"status" binding:"omitempty,oneof=0 1" json:"status"`
	OwnerID     *uint  `form:"owner_id" json:"owner_id"`
	Search      string `form:"search" json:"search"`
}

// CreateEmailChannelRequest represents the request to create an email notification channel
type CreateEmailChannelRequest struct {
	Code         string  `json:"code" binding:"required" validate:"required,min=3,max=50" example:"email-channel-001"`
	Name         string  `json:"name" binding:"required" validate:"required,min=1,max=100" example:"My Email Channel"`
	Description  *string `json:"description" validate:"max=500" example:"Email notification channel for alerts"`
	SMTPHost     string  `json:"smtp_host" binding:"required" validate:"required" example:"smtp.gmail.com"`
	SMTPPort     int     `json:"smtp_port" binding:"required" validate:"required,min=1,max=65535" example:"587"`
	SMTPSecurity string  `json:"smtp_security" binding:"required" validate:"required,oneof=none tls ssl" example:"tls"`
	SMTPTimeout  int     `json:"smtp_timeout" validate:"min=0,max=300" example:"30"`
	SMTPUsername string  `json:"smtp_username" binding:"required" validate:"required" example:"user@example.com"`
	SMTPPassword string  `json:"smtp_password" binding:"required" validate:"required" example:"your-password"`
	SenderEmail  string  `json:"sender_email" binding:"required" validate:"required,email" example:"noreply@example.com"`
	SenderName   string  `json:"sender_name" binding:"required" validate:"required" example:"System Notifications"`
	Encoding     string  `json:"encoding" validate:"oneof=utf-8 gbk gb2312" example:"utf-8"`
	RateLimit    int     `json:"rate_limit" validate:"min=0,max=1000" example:"3"`
	RetryCount   int     `json:"retry_count" validate:"min=0,max=10" example:"3"`
	QuietHours   string  `json:"quiet_hours" validate:"max=100" example:"22:00-08:00"`
}

// CreateWebhookChannelRequest represents the request to create a webhook notification channel
type CreateWebhookChannelRequest struct {
	Code        string            `json:"code" binding:"required" validate:"required,min=3,max=50" example:"webhook-channel-001"`
	Name        string            `json:"name" binding:"required" validate:"required,min=1,max=100" example:"My Webhook Channel"`
	Description *string           `json:"description" validate:"max=500" example:"Webhook notification channel for alerts"`
	URL         string            `json:"url" binding:"required" validate:"required,url" example:"https://api.example.com/webhook"`
	Method      string            `json:"method" binding:"required" validate:"required,oneof=GET POST PUT PATCH DELETE" example:"POST"`
	Headers     map[string]string `json:"headers" validate:"dive,max=200"`
	Secret      string            `json:"secret" validate:"max=100" example:"your-webhook-secret"`
	Timeout     int               `json:"timeout" validate:"min=0,max=300" example:"30"`
	Encoding    string            `json:"encoding" validate:"oneof=utf-8 gbk gb2312" example:"utf-8"`
	RateLimit   int               `json:"rate_limit" validate:"min=0,max=1000" example:"3"`
	RetryCount  int               `json:"retry_count" validate:"min=0,max=10" example:"3"`
	QuietHours  string            `json:"quiet_hours" validate:"max=100" example:"22:00-08:00"`
}

// UpdateEmailChannelRequest represents the request to update an email notification channel
type UpdateEmailChannelRequest struct {
	Name         *string `json:"name" validate:"min=1,max=100" example:"Updated Email Channel"`
	Description  *string `json:"description" validate:"max=500" example:"Updated email notification channel"`
	SMTPHost     *string `json:"smtp_host" example:"smtp.gmail.com"`
	SMTPPort     *int    `json:"smtp_port" validate:"min=1,max=65535" example:"587"`
	SMTPSecurity *string `json:"smtp_security" validate:"oneof=none tls ssl" example:"tls"`
	SMTPTimeout  *int    `json:"smtp_timeout" validate:"min=0,max=300" example:"30"`
	SMTPUsername *string `json:"smtp_username" example:"user@example.com"`
	SMTPPassword *string `json:"smtp_password" example:"your-password"`
	SenderEmail  *string `json:"sender_email" validate:"email" example:"noreply@example.com"`
	SenderName   *string `json:"sender_name" example:"System Notifications"`
	Encoding     *string `json:"encoding" validate:"oneof=utf-8 gbk gb2312" example:"utf-8"`
	RateLimit    *int    `json:"rate_limit" validate:"min=0,max=1000" example:"3"`
	RetryCount   *int    `json:"retry_count" validate:"min=0,max=10" example:"3"`
	QuietHours   *string `json:"quiet_hours" validate:"max=100" example:"22:00-08:00"`
}

// UpdateWebhookChannelRequest represents the request to update a webhook notification channel
type UpdateWebhookChannelRequest struct {
	Name        *string            `json:"name" validate:"min=1,max=100" example:"Updated Webhook Channel"`
	Description *string            `json:"description" validate:"max=500" example:"Updated webhook notification channel"`
	URL         *string            `json:"url" validate:"url" example:"https://api.example.com/webhook"`
	Method      *string            `json:"method" validate:"oneof=GET POST PUT PATCH DELETE" example:"POST"`
	Headers     *map[string]string `json:"headers" validate:"dive,max=200"`
	Secret      *string            `json:"secret" validate:"max=100" example:"your-webhook-secret"`
	Timeout     *int               `json:"timeout" validate:"min=0,max=300" example:"30"`
	Encoding    *string            `json:"encoding" validate:"oneof=utf-8 gbk gb2312" example:"utf-8"`
	RateLimit   *int               `json:"rate_limit" validate:"min=0,max=1000" example:"3"`
	RetryCount  *int               `json:"retry_count" validate:"min=0,max=10" example:"3"`
	QuietHours  *string            `json:"quiet_hours" validate:"max=100" example:"22:00-08:00"`
}

// TestEmailChannelRequest represents the request for testing email channel
type TestEmailChannelRequest struct {
	Code      string `json:"code" binding:"required"`
	Recipient string `json:"recipient" binding:"required,email"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

// TestWebhookChannelRequest represents the request for testing webhook channel
type TestWebhookChannelRequest struct {
	Code    string      `json:"code" binding:"required"`
	Payload interface{} `json:"payload"`
}
