package service

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
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

// GetNotificationRecordList retrieves notification records list with pagination and filtering
func (s *notificationRecordService) GetNotificationRecordList(
	ctx context.Context,
	req *request.GetNotificationRecordListRequest,
) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification record list",
		logger.String("service", "notification_record"),
		logger.String("operation", "GetNotificationRecordList"),
		logger.Int("page", req.PaginationRequest.Page),
		logger.Int("page_size", req.PaginationRequest.PageSize))

	// Get records from repository
	records, total, err := s.notificationRepo.GetList(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification records from repository", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Convert to response format
	responses := make([]response.NotificationRecordResponse, 0, len(records))
	for _, record := range records {
		resp := response.NotificationRecordResponse(*record)
		responses = append(responses, resp)
	}

	s.logger.InfoContext(ctx, "Notification record list retrieved successfully",
		logger.Int("total_count", int(total)),
		logger.Int("returned_count", len(responses)))

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	result := &common.PaginationResponse{
		Page:       req.GetPage(),
		PageSize:   req.GetPageSize(),
		Total:      total,
		TotalPages: totalPages,
		Items:      responses,
	}
	return result, nil
}

// GetNotificationRecordByID retrieves a specific notification record by ID
func (s *notificationRecordService) GetNotificationRecordByID(ctx context.Context, id uint) (*response.NotificationRecordResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification record by ID",
		logger.String("service", "notification_record"),
		logger.String("operation", "GetNotificationRecordByID"),
		logger.Uint("record_id", id))

	// Get record
	record, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "Notification record not found", logger.Uint("record_id", id))
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		s.logger.ErrorContext(ctx, "Failed to get notification record from repository", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	s.logger.InfoContext(ctx, "Notification record retrieved successfully",
		logger.Uint("record_id", id))

	resp := response.NotificationRecordResponse(*record)
	return &resp, nil
}
