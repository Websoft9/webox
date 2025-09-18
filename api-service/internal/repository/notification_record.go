package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
)

type notificationRecordRepository struct {
	db *gorm.DB
}

// NewNotificationRecordRepository creates a new notification record repository instance
func NewNotificationRecordRepository(db *gorm.DB) repository.NotificationRecordRepository {
	return &notificationRecordRepository{db: db}
}

// GetByID retrieves a notification record by ID
func (r *notificationRecordRepository) GetByID(ctx context.Context, id uint) (*model.NotificationRecord, error) {
	var record model.NotificationRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetList retrieves notification records with pagination and filtering
func (r *notificationRecordRepository) GetList(ctx context.Context, req *request.GetNotificationRecordListRequest) ([]*model.NotificationRecord, int64, error) {
	var records []*model.NotificationRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.NotificationRecord{})

	// Apply filters
	if req.ChannelType != "" {
		query = query.Where("channel_type = ?", req.ChannelType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Recipient != "" {
		query = query.Where("recipient LIKE ?", "%"+req.Recipient+"%")
	}
	if req.SentStart != "" {
		sentStart, err := time.Parse("2006-01-02 15:04:05", req.SentStart)
		if err == nil {
			query = query.Where("sent_at >= ?", sentStart)
		}
	}
	if req.SentEnd != "" {
		sentEnd, err := time.Parse("2006-01-02 15:04:05", req.SentEnd)
		if err == nil {
			query = query.Where("sent_at <= ?", sentEnd)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and order
	offset := req.GetOffset()
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.GetPageSize()).
		Find(&records).Error

	return records, total, err
}
