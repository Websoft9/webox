package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// NotificationRecordRepository defines the interface for notification record data access
type NotificationRecordRepository interface {
	// GetByID retrieves a notification record by ID
	GetByID(ctx context.Context, id uint) (*model.NotificationRecord, error)

	// GetList retrieves notification records with pagination and filtering
	GetList(ctx context.Context, req *request.GetNotificationRecordListRequest) ([]*model.NotificationRecord, int64, error)
}
