package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// NotificationTemplateService defines the interface for notification template business logic
type NotificationTemplateService interface {
	// CreateTemplate creates a new notification template
	CreateTemplate(ctx context.Context, req *request.CreateNotificationTemplateRequest) (*response.NotificationTemplateDetailResponse, error)

	// GetTemplate retrieves a notification template by ID
	GetTemplate(ctx context.Context, id uint) (*response.NotificationTemplateDetailResponse, error)

	// GetTemplateList retrieves a paginated list of notification templates
	GetTemplateList(ctx context.Context, req *request.GetNotificationTemplateListRequest) (*common.PaginationResponse, error)

	// UpdateTemplate updates an existing notification template
	UpdateTemplate(ctx context.Context, id uint, req *request.UpdateNotificationTemplateRequest) (*response.NotificationTemplateDetailResponse, error)

	// DeleteTemplate deletes a notification template
	DeleteTemplate(ctx context.Context, id uint) error

	// TestTemplate tests a template by sending it
	TestTemplate(ctx context.Context, id uint, req *request.TestNotificationTemplateRequest) error
}
