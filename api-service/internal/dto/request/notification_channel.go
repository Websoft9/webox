package request

import "api-service/internal/dto"

// GetNotificationChannelListRequest represents the request for getting notification channel list
type GetNotificationChannelListRequest struct {
	dto.BaseListRequest
	ChannelType string `form:"channel_type" binding:"omitempty,oneof=EMAIL WEBHOOK" json:"channel_type"`
	Status      *int8  `form:"status" binding:"omitempty,oneof=-1 0 1" json:"status"`
	OwnerID     *uint  `form:"owner_id" json:"owner_id"`
	Search      string `form:"search" json:"search"`
}

// CreateEmailChannelRequest represents the request for creating email channel
type CreateEmailChannelRequest struct {
	Code        string `json:"code" binding:"required,min=3,max=64,alphanum"`
	Name        string `json:"name" binding:"required,min=2,max=128"`
	Description string `json:"description"`

	// Email specific configuration
	SMTPHost     string `json:"smtp_host" binding:"required"`
	SMTPPort     int    `json:"smtp_port" binding:"required,min=1,max=65535"`
	SMTPSecurity string `json:"smtp_security" binding:"omitempty,oneof=none ssl tls"`
	SMTPTimeout  int    `json:"smtp_timeout" binding:"omitempty,min=1,max=300"`
	SMTPUsername string `json:"smtp_username" binding:"required"`
	SMTPPassword string `json:"smtp_password" binding:"required"`
	SenderEmail  string `json:"sender_email" binding:"required,email"`
	SenderName   string `json:"sender_name"`
	Encoding     string `json:"encoding" binding:"omitempty,oneof=utf-8 gbk gb2312"`
	RateLimit    int    `json:"rate_limit" binding:"omitempty,min=1"`
	RetryCount   int    `json:"retry_count" binding:"omitempty,min=0,max=10"`
	QuietHours   string `json:"quiet_hours"`
}

// CreateWebhookChannelRequest represents the request for creating webhook channel
type CreateWebhookChannelRequest struct {
	Code        string `json:"code" binding:"required,min=3,max=64,alphanum"`
	Name        string `json:"name" binding:"required,min=2,max=128"`
	Description string `json:"description"`

	// Webhook specific configuration
	URL        string            `json:"url" binding:"required,url"`
	Method     string            `json:"method" binding:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Headers    map[string]string `json:"headers"`
	Secret     string            `json:"secret"`
	Timeout    int               `json:"timeout" binding:"omitempty,min=1,max=300"`
	Encoding   string            `json:"encoding" binding:"omitempty,oneof=utf-8 gbk gb2312"`
	RateLimit  int               `json:"rate_limit" binding:"omitempty,min=1"`
	RetryCount int               `json:"retry_count" binding:"omitempty,min=0,max=10"`
	QuietHours string            `json:"quiet_hours"`
}

// UpdateEmailChannelRequest represents the request for updating email channel
type UpdateEmailChannelRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=128"`
	Description *string `json:"description"`

	// Email specific configuration - all optional for partial updates
	SMTPHost     *string `json:"smtp_host"`
	SMTPPort     *int    `json:"smtp_port" binding:"omitempty,min=1,max=65535"`
	SMTPSecurity *string `json:"smtp_security" binding:"omitempty,oneof=none ssl tls"`
	SMTPTimeout  *int    `json:"smtp_timeout" binding:"omitempty,min=1,max=300"`
	SMTPUsername *string `json:"smtp_username"`
	SMTPPassword *string `json:"smtp_password"`
	SenderEmail  *string `json:"sender_email" binding:"omitempty,email"`
	SenderName   *string `json:"sender_name"`
	Encoding     *string `json:"encoding" binding:"omitempty,oneof=utf-8 gbk gb2312"`
	RateLimit    *int    `json:"rate_limit" binding:"omitempty,min=1"`
	RetryCount   *int    `json:"retry_count" binding:"omitempty,min=0,max=10"`
	QuietHours   *string `json:"quiet_hours"`
}

// UpdateWebhookChannelRequest represents the request for updating webhook channel
type UpdateWebhookChannelRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=128"`
	Description *string `json:"description"`

	// Webhook specific configuration - all optional for partial updates
	URL        *string            `json:"url" binding:"omitempty,url"`
	Method     *string            `json:"method" binding:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Headers    *map[string]string `json:"headers"`
	Secret     *string            `json:"secret"`
	Timeout    *int               `json:"timeout" binding:"omitempty,min=1,max=300"`
	Encoding   *string            `json:"encoding" binding:"omitempty,oneof=utf-8 gbk gb2312"`
	RateLimit  *int               `json:"rate_limit" binding:"omitempty,min=1"`
	RetryCount *int               `json:"retry_count" binding:"omitempty,min=0,max=10"`
	QuietHours *string            `json:"quiet_hours"`
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
