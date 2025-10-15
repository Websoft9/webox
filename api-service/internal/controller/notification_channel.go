package controller

import (
	"net/http"

	response "api-service/internal/dto/common"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// NotificationChannelController handles notification channel related HTTP requests
type NotificationChannelController struct {
	channelService service.NotificationChannelService
	logger         logger.Logger
	validator      *validator.Validate
}

// NewNotificationChannelController creates a new notification channel controller
func NewNotificationChannelController(
	channelService service.NotificationChannelService,
	logger logger.Logger,
	validator *validator.Validate,
) *NotificationChannelController {
	return &NotificationChannelController{
		channelService: channelService,
		logger:         logger,
		validator:      validator,
	}
}

// GetChannelList get notification channels list with pagination
// @Summary Get notification channels list
// @Description Get paginated list of notification channels with filtering
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param channel_type query string false "Channel type filter" Enums(EMAIL,WEBHOOK)
// @Success		200		{object}	common.APIResponse	"Test Success"
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels [get]
func (ctrl *NotificationChannelController) GetChannelList(c *gin.Context) {
	var req request.GetNotificationChannelListRequest

	// Bind and validate query parameters
	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.channelService.GetChannelList(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// GetChannelByCode get notification channel details by code
// @Summary Get notification channel by code
// @Description Get detailed information about a specific notification channel
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Success 200 {object} common.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/{code} [get]
func (ctrl *NotificationChannelController) GetChannelByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		response.WithError(c, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	// Call service
	result, err := ctrl.channelService.GetChannelByCode(c.Request.Context(), code)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification channel retrieved successfully",
		logger.String("code", code))

	response.SuccessWithData(c, result)
}

// CreateEmailChannel create email notification channel
// @Summary Create email notification channel
// @Description Create a new email notification channel with SMTP configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateEmailChannelRequest true "Email channel configuration"
// @Success 201 {object} common.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/email [post]
func (ctrl *NotificationChannelController) CreateEmailChannel(c *gin.Context) {
	var req request.CreateEmailChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.channelService.CreateEmailChannel(c.Request.Context(), &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel created successfully",
		logger.String("code", req.Code))

	response.SuccessWithData(c, result)
}

// CreateWebhookChannel create webhook notification channel
// @Summary Create webhook notification channel
// @Description Create a new webhook notification channel with URL and headers configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateWebhookChannelRequest true "Webhook channel configuration"
// @Success 201 {object} common.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/webhook [post]
func (ctrl *NotificationChannelController) CreateWebhookChannel(c *gin.Context) {
	var req request.CreateWebhookChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.channelService.CreateWebhookChannel(c.Request.Context(), &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel created successfully",
		logger.String("code", req.Code))

	response.SuccessWithData(c, result)
}

// UpdateEmailChannel update email notification channel
// @Summary Update email notification channel
// @Description Update an existing email notification channel configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Param request body request.UpdateEmailChannelRequest true "Updated email channel configuration"
// @Success 200 {object} common.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/{code}/email [put]
func (ctrl *NotificationChannelController) UpdateEmailChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		response.WithError(c, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	var req request.UpdateEmailChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.channelService.UpdateEmailChannel(c.Request.Context(), code, &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel updated successfully",
		logger.String("code", code))

	response.SuccessWithData(c, result)
}

// UpdateWebhookChannel update webhook notification channel
// @Summary Update webhook notification channel
// @Description Update an existing webhook notification channel configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Param request body request.UpdateWebhookChannelRequest true "Updated webhook channel configuration"
// @Success 200 {object} common.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/{code}/webhook [put]
func (ctrl *NotificationChannelController) UpdateWebhookChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		response.WithError(c, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	var req request.UpdateWebhookChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.channelService.UpdateWebhookChannel(c.Request.Context(), code, &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel updated successfully",
		logger.String("code", code))

	response.SuccessWithData(c, result)
}

// DeleteChannel delete notification channel
// @Summary Delete notification channel
// @Description Delete an existing notification channel by code
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Success 204 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/{code} [delete]
func (ctrl *NotificationChannelController) DeleteChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		response.WithError(c, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	err := ctrl.channelService.DeleteChannel(c.Request.Context(), code, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification channel deleted successfully",
		logger.String("code", code))

	c.Status(http.StatusNoContent)
}

// TestEmailChannel test email notification channel
// @Summary Test email notification channel
// @Description Test an email notification channel by sending a test message
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TestEmailChannelRequest true "Test email configuration"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/test/email [post]
func (ctrl *NotificationChannelController) TestEmailChannel(c *gin.Context) {
	var req request.TestEmailChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	err := ctrl.channelService.TestEmailChannel(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel tested",
		logger.String("code", req.Code),
		logger.Bool("success", true))

	response.Success(c)
}

// TestWebhookChannel test webhook notification channel
// @Summary Test webhook notification channel
// @Description Test a webhook notification channel by sending a test message
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TestWebhookChannelRequest true "Test webhook configuration"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/channels/test/webhook [post]
func (ctrl *NotificationChannelController) TestWebhookChannel(c *gin.Context) {
	var req request.TestWebhookChannelRequest

	// Bind and validate request
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	err := ctrl.channelService.TestWebhookChannel(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel tested",
		logger.String("code", req.Code),
		logger.Bool("success", true))

	response.Success(c)
}
