package request

import "api-service/internal/dto/common"

// CreateNotificationTemplateRequest represents the request to create a notification template
type CreateNotificationTemplateRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=64"`
	TemplateType string  `json:"template_type" validate:"required,oneof=EMAIL WEBHOOK INTERNAL"`
	Subject      *string `json:"subject" validate:"required,min=2,max=255"`
	Content      string  `json:"content" validate:"required"`
}

// UpdateNotificationTemplateRequest represents the request to update a notification template
type UpdateNotificationTemplateRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=2,max=64"`
	Subject *string `json:"subject" validate:"required,min=2,max=255"`
	Content *string `json:"content" validate:"required"`
}

// GetNotificationTemplateListRequest represents the request to get notification template list
type GetNotificationTemplateListRequest struct {
	common.BaseListRequest
	TemplateType string `form:"template_type" binding:"omitempty,oneof=EMAIL WEBHOOK INTERNAL" json:"template_type"`
	Status       *int   `form:"status" binding:"omitempty,oneof=0 1" json:"status"`
	IsSystem     *int   `form:"is_system" binding:"omitempty,oneof=0 1" json:"is_system"`
}

// TestNotificationTemplateRequest represents the request to test a notification template
type TestNotificationTemplateRequest struct {
	Code         string                 `json:"code" binding:"required" validate:"required"`
	Recipient    string                 `json:"recipient" binding:"required,email" validate:"required,email"`
	VariableData map[string]interface{} `json:"variable_data"`
}
