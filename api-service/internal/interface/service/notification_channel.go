package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// NotificationChannelService defines the interface for notification channel service
type NotificationChannelService interface {
	// GetChannelList retrieves notification channels list with pagination and filtering
	GetChannelList(ctx context.Context, req *request.GetNotificationChannelListRequest) (*common.PaginationResponse, error)

	// GetChannelByCode retrieves a specific notification channel by code
	GetChannelByCode(ctx context.Context, code string) (*response.NotificationChannelResponse, error)

	// CreateEmailChannel creates a new email notification channel
	CreateEmailChannel(ctx context.Context, req *request.CreateEmailChannelRequest, userID uint) (*response.NotificationChannelResponse, error)

	// CreateWebhookChannel creates a new webhook notification channel
	CreateWebhookChannel(ctx context.Context, req *request.CreateWebhookChannelRequest, userID uint) (*response.NotificationChannelResponse, error)

	// UpdateEmailChannel updates an existing email notification channel
	UpdateEmailChannel(ctx context.Context, code string, req *request.UpdateEmailChannelRequest, userID uint) (*response.NotificationChannelResponse, error)

	// UpdateWebhookChannel updates an existing webhook notification channel
	UpdateWebhookChannel(ctx context.Context, code string, req *request.UpdateWebhookChannelRequest, userID uint) (*response.NotificationChannelResponse, error)

	// DeleteChannel soft deletes a notification channel
	DeleteChannel(ctx context.Context, code string, userID uint) error

	// TestEmailChannel tests email channel configuration
	TestEmailChannel(ctx context.Context, req *request.TestEmailChannelRequest) error

	// TestWebhookChannel tests webhook channel configuration
	TestWebhookChannel(ctx context.Context, req *request.TestWebhookChannelRequest) error
}
