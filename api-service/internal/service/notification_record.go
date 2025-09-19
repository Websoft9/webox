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
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

type notificationRecordService struct {
	notificationRepo repository.NotificationRecordRepository
	logger           logger.Logger
	i18n             *i18n.I18n
}

// NewNotificationRecordService creates a new notification record service instance
func NewNotificationRecordService(
	notificationRepo repository.NotificationRecordRepository,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.NotificationRecordService {
	return &notificationRecordService{
		notificationRepo: notificationRepo,
		logger:           logger,
		i18n:             i18n,
	}
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
		s.logger.ErrorContext(ctx, "Failed to get notification records from repository", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.get_records_failed"))
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
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.record_not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification record from repository", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.get_record_failed"))
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
		Recipient:     record.Recipient,
		Subject:       record.Subject,
		Content:       record.Content,
		Status:        record.Status,
		RetryCount:    record.RetryCount,
		ReferenceID:   record.ReferenceID,
		ReferenceType: record.ReferenceType,
		CreatedAt:     record.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     record.UpdatedAt.Format("2006-01-02 15:04:05"),
		ErrorMsg:      record.ErrorMsg,
	}

	// Format optional timestamp fields
	if record.SentAt != nil {
		sentAt := record.SentAt.Format("2006-01-02 15:04:05")
		resp.SentAt = &sentAt
	}

	return resp
}
