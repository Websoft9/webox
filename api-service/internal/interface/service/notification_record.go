package service

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// NotificationRecordService defines the interface for notification record business logic
type NotificationRecordService interface {
	// GetNotificationRecordList retrieves notification records list with pagination and filtering
	GetNotificationRecordList(ctx context.Context, req *request.GetNotificationRecordListRequest) (*response.NotificationRecordListResponse, error)

	// GetNotificationRecordByID retrieves a specific notification record by ID
	GetNotificationRecordByID(ctx context.Context, id uint) (*response.NotificationRecordDetailResponse, error)
}
