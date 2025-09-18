package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
)

// NotificationRecordController handles notification record related HTTP requests
type NotificationRecordController struct {
	notificationRecordService service.NotificationRecordService
	i18n                      *i18n.I18n
}

// NewNotificationRecordController creates a new notification record controller instance
func NewNotificationRecordController(
	notificationRecordService service.NotificationRecordService,
	i18n *i18n.I18n,
) *NotificationRecordController {
	return &NotificationRecordController{
		notificationRecordService: notificationRecordService,
		i18n:                      i18n,
	}
}

// GetNotificationRecords retrieves notification records list with pagination and filtering
// @Summary Get notification records list
// @Description Retrieve notification records list with pagination and filtering options
// @Tags Notification Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Param channel_type query string false "Channel type filter" Enums(EMAIL,WEBHOOK,INTERNAL)
// @Param status query string false "Status filter" Enums(PENDING,SENT,FAILED,RETRY)
// @Param recipient query string false "Recipient filter"
// @Param sent_start query string false "Sent time start filter" format(2006-01-02 15:04:05)
// @Param sent_end query string false "Sent time end filter" format(2006-01-02 15:04:05)
// @Success 200 {object} response.Response{data=response.NotificationRecordListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/notifications/records [get]
// @Security BearerAuth
func (c *NotificationRecordController) GetNotificationRecords(ctx *gin.Context) {
	var req request.GetNotificationRecordListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	records, err := c.notificationRecordService.GetNotificationRecordList(ctx.Request.Context(), &req)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ResponseOKWithData(ctx, records, "notification.get_records_success", c.i18n)
}

// GetNotificationRecord retrieves a specific notification record by ID
// @Summary Get notification record details
// @Description Retrieve detailed information of a specific notification record
// @Tags Notification Management
// @Accept json
// @Produce json
// @Param id path int true "Notification record ID"
// @Success 200 {object} response.Response{data=response.NotificationRecordDetailResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/notifications/records/{id} [get]
// @Security BearerAuth
func (c *NotificationRecordController) GetNotificationRecord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	record, err := c.notificationRecordService.GetNotificationRecordByID(ctx.Request.Context(), uint(id))
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ResponseOKWithData(ctx, record, "notification.get_record_success", c.i18n)
}
