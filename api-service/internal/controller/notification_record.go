package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"

	response "api-service/internal/dto/common"

	"github.com/go-playground/validator/v10"
)

// NotificationRecordController handles notification record related HTTP requests
type NotificationRecordController struct {
	notificationRecordService service.NotificationRecordService
	validator                 *validator.Validate
	logger                    logger.Logger
}

// NewNotificationRecordController creates a new notification record controller instance
func NewNotificationRecordController(
	notificationRecordService service.NotificationRecordService,
	validator *validator.Validate,
	logger logger.Logger,
) *NotificationRecordController {
	return &NotificationRecordController{
		notificationRecordService: notificationRecordService,
		validator:                 validator,
		logger:                    logger,
	}
}

// GetNotificationRecords retrieves notification records list with pagination and filtering
// @Summary Get notification records list
// @Description Retrieve notification records list with pagination and filtering options
// @Tags Notification Records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Param channel_type query string false "Channel type filter" Enums(EMAIL,WEBHOOK,INTERNAL)
// @Param status query string false "Status filter" Enums(PENDING,SENT,FAILED,RETRY)
// @Param recipient query string false "Recipient filter"
// @Param sent_start query string false "Sent time start filter" format(2006-01-02 15:04:05)
// @Param sent_end query string false "Sent time end filter" format(2006-01-02 15:04:05)
// @Success 200 {object} common.APIResponse{data=[]response.NotificationRecordResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/records [get]
// @Security BearerAuth
func (c *NotificationRecordController) GetNotificationRecords(ctx *gin.Context) {
	var req request.GetNotificationRecordListRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	records, err := c.notificationRecordService.GetNotificationRecordList(ctx.Request.Context(), &req)

	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, records)
}

// GetNotificationRecord retrieves a specific notification record by ID
// @Summary Get notification record details
// @Description Retrieve detailed information of a specific notification record
// @Tags Notification Records
// @Accept json
// @Produce json
// @Param id path int true "Notification record ID"
// @Success 200 {object} common.APIResponse{data=response.NotificationRecordResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/records/{id} [get]
// @Security BearerAuth
func (c *NotificationRecordController) GetNotificationRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Channel code is required")
		response.WithError(ctx, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	record, err := c.notificationRecordService.GetNotificationRecordByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, record)
}
