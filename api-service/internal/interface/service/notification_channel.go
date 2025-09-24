package service

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// NotificationChannelService defines the interface for notification channel service
type NotificationChannelService interface {
	// GetChannelList retrieves notification channels list with pagination and filtering
	GetChannelList(ctx context.Context, req *request.GetNotificationChannelListRequest) (*response.NotificationChannelListResponse, error)

	// GetChannelByCode retrieves a specific notification channel by code
	GetChannelByCode(ctx context.Context, code string) (*response.NotificationChannelDetailResponse, error)

	// CreateEmailChannel creates a new email notification channel
	CreateEmailChannel(ctx context.Context, req *request.CreateEmailChannelRequest, userID uint) (*response.NotificationChannelDetailResponse, error)

	// CreateWebhookChannel creates a new webhook notification channel
	CreateWebhookChannel(ctx context.Context, req *request.CreateWebhookChannelRequest, userID uint) (*response.NotificationChannelDetailResponse, error)

	// UpdateEmailChannel updates an existing email notification channel
	UpdateEmailChannel(ctx context.Context, code string, req *request.UpdateEmailChannelRequest, userID uint) (*response.NotificationChannelDetailResponse, error)

	// UpdateWebhookChannel updates an existing webhook notification channel
	UpdateWebhookChannel(ctx context.Context, code string, req *request.UpdateWebhookChannelRequest, userID uint) (*response.NotificationChannelDetailResponse, error)

	// DeleteChannel soft deletes a notification channel
	DeleteChannel(ctx context.Context, code string, userID uint) error

	// TestEmailChannel tests email channel configuration
	TestEmailChannel(ctx context.Context, req *request.TestEmailChannelRequest) (*response.TestChannelResponse, error)

	// TestWebhookChannel tests webhook channel configuration
	TestWebhookChannel(ctx context.Context, req *request.TestWebhookChannelRequest) (*response.TestChannelResponse, error)
}
