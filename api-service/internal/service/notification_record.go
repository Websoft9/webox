package service

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

type notificationRecordService struct {
	notificationRepo repository.NotificationRecordRepository
	logger           logger.Logger
}

// NewNotificationRecordService creates a new notification record service instance
func NewNotificationRecordService(
	notificationRepo repository.NotificationRecordRepository,
	logger logger.Logger,
) service.NotificationRecordService {
	return &notificationRecordService{
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// logAndWrapError logs an error and wraps it with additional context
func (s *notificationRecordService) logAndWrapError(ctx context.Context, err error, message, wrapMessage string) error {
	s.logger.ErrorContext(ctx, message, logger.ErrorField(err))
	return errors.WrapError(err, errors.CodeRecordQueryFailed, wrapMessage)
}

// GetNotificationRecordList retrieves notification records list with pagination and filtering
func (s *notificationRecordService) GetNotificationRecordList(
	ctx context.Context,
	req *request.GetNotificationRecordListRequest,
) (*response.NotificationRecordListResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification record list",
		logger.String("service", "notification_record"),
		logger.String("operation", "GetNotificationRecordList"),
		logger.Int("page", req.PaginationRequest.Page),
		logger.Int("page_size", req.PaginationRequest.PageSize))

	// Get records from repository
	records, total, err := s.notificationRepo.GetList(ctx, req)
	if err != nil {
		return nil, s.logAndWrapError(ctx, err, "Failed to get notification records from repository", "failed to get notification records")
	}

	// Convert to response format
	responses := make([]response.NotificationRecordResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, *s.convertToResponse(record))
	}

	// Build paginated response
	pageSize := req.GetPageSize()
	page := req.Page
	if page <= 0 {
		page = 1
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	listResp := &response.NotificationRecordListResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      responses,
	}

	s.logger.InfoContext(ctx, "Notification record list retrieved successfully",
		logger.Int("total_count", int(total)),
		logger.Int("returned_count", len(responses)))

	return listResp, nil
}

// GetNotificationRecordByID retrieves a specific notification record by ID
func (s *notificationRecordService) GetNotificationRecordByID(ctx context.Context, id uint) (*response.NotificationRecordDetailResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification record by ID",
		logger.String("service", "notification_record"),
		logger.String("operation", "GetNotificationRecordByID"),
		logger.Uint("record_id", id))

	// Get record
	record, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "Notification record not found",
				logger.Uint("record_id", id))
			return nil, errors.NewAppError(errors.CodeRecordNotFound, "Notification record not found")
		}
		return nil, s.logAndWrapError(ctx, err, "Failed to get notification record from repository", "failed to get notification record")
	}

	// Convert to detail response
	detailResp := &response.NotificationRecordDetailResponse{
		NotificationRecordResponse: *s.convertToResponse(record),
	}

	s.logger.InfoContext(ctx, "Notification record retrieved successfully",
		logger.Uint("record_id", id))

	return detailResp, nil
}

// convertToResponse converts model.NotificationRecord to response.NotificationRecordResponse
func (s *notificationRecordService) convertToResponse(record *model.NotificationRecord) *response.NotificationRecordResponse {
	resp := &response.NotificationRecordResponse{
		ID:            record.ID,
		TemplateID:    record.TemplateID,
		ChannelType:   record.ChannelType,
		UserID:        record.UserID,
		Recipient:     record.Recipient,
		Subject:       record.Subject,
		Content:       record.Content,
		Status:        record.Status,
		RetryCount:    record.RetryCount,
		ReferenceID:   record.ReferenceID,
		ReferenceType: record.ReferenceType,
		IsRead:        record.IsRead,
		CreatedAt:     record.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     record.UpdatedAt.Format("2006-01-02 15:04:05"),
		ErrorMsg:      record.ErrorMsg,
	}

	// Format optional timestamp fields
	if record.SentAt != nil {
		sentAt := record.SentAt.Format("2006-01-02 15:04:05")
		resp.SentAt = &sentAt
	}

	if record.ReadAt != nil {
		readAt := record.ReadAt.Format("2006-01-02 15:04:05")
		resp.ReadAt = &readAt
	}

	return resp
}
