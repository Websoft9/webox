package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/response"
)

// NotificationChannelController handles notification channel related HTTP requests
type NotificationChannelController struct {
	channelService service.NotificationChannelService
	logger         logger.Logger
	i18n           *i18n.I18n
}

// NewNotificationChannelController creates a new notification channel controller
func NewNotificationChannelController(
	channelService service.NotificationChannelService,
	logger logger.Logger,
	i18n *i18n.I18n,
) *NotificationChannelController {
	return &NotificationChannelController{
		channelService: channelService,
		logger:         logger,
		i18n:           i18n,
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
// @Param channel_type query string false "Channel type filter" Enums(EMAIL,WEBHOOK,INTERNAL)
// @Success 200 {object} response.APIResponse{data=response.NotificationChannelListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels [get]
func (ctrl *NotificationChannelController) GetChannelList(c *gin.Context) {
	var req request.GetNotificationChannelListRequest

	// Bind and validate query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, "Invalid query parameters"))
		return
	}

	// Set default pagination if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Call service
	result, err := ctrl.channelService.GetChannelList(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to get notification channel list", logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification channel list retrieved successfully",
		logger.Int("count", len(result.Items)))

	ResponseOKWithData(c, result, "notification.channel.list_success", ctrl.i18n)
}

// GetChannelByCode get notification channel details by code
// @Summary Get notification channel by code
// @Description Get detailed information about a specific notification channel
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Success 200 {object} response.APIResponse{data=response.NotificationChannelDetailResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/{code} [get]
func (ctrl *NotificationChannelController) GetChannelByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, "Channel code is required"))
		return
	}

	// Call service
	result, err := ctrl.channelService.GetChannelByCode(c.Request.Context(), code)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to get notification channel",
			logger.String("code", code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification channel retrieved successfully",
		logger.String("code", code))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.get_success"), result)
}

// CreateEmailChannel create email notification channel
// @Summary Create email notification channel
// @Description Create a new email notification channel with SMTP configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateEmailChannelRequest true "Email channel configuration"
// @Success 201 {object} response.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/email [post]
func (ctrl *NotificationChannelController) CreateEmailChannel(c *gin.Context) {
	var req request.CreateEmailChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Get user ID from context (from JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.logger.ErrorContext(c.Request.Context(), "User ID not found in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInvalidToken, "User authentication required"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user ID type in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, "Invalid user authentication"))
		return
	}

	// Call service
	result, err := ctrl.channelService.CreateEmailChannel(c.Request.Context(), &req, userIDUint)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to create email notification channel",
			logger.String("code", req.Code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel created successfully",
		logger.String("code", req.Code))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.create_success"), result)
}

// CreateWebhookChannel create webhook notification channel
// @Summary Create webhook notification channel
// @Description Create a new webhook notification channel with URL and headers configuration
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateWebhookChannelRequest true "Webhook channel configuration"
// @Success 201 {object} response.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/webhook [post]
func (ctrl *NotificationChannelController) CreateWebhookChannel(c *gin.Context) {
	var req request.CreateWebhookChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Get user ID from context (from JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.logger.ErrorContext(c.Request.Context(), "User ID not found in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInvalidToken, "User authentication required"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user ID type in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, "Invalid user authentication"))
		return
	}

	// Call service
	result, err := ctrl.channelService.CreateWebhookChannel(c.Request.Context(), &req, userIDUint)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to create webhook notification channel",
			logger.String("code", req.Code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel created successfully",
		logger.String("code", req.Code))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.create_success"), result)
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
// @Success 200 {object} response.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/{code}/email [put]
func (ctrl *NotificationChannelController) UpdateEmailChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, "Channel code is required"))
		return
	}

	var req request.UpdateEmailChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Get user ID from context (from JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.logger.ErrorContext(c.Request.Context(), "User ID not found in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInvalidToken, "User authentication required"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user ID type in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, "Invalid user authentication"))
		return
	}

	// Call service
	result, err := ctrl.channelService.UpdateEmailChannel(c.Request.Context(), code, &req, userIDUint)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to update email notification channel",
			logger.String("code", code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel updated successfully",
		logger.String("code", code))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.update_success"), result)
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
// @Success 200 {object} response.APIResponse{data=response.NotificationChannelResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/{code}/webhook [put]
func (ctrl *NotificationChannelController) UpdateWebhookChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, "Channel code is required"))
		return
	}

	var req request.UpdateWebhookChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Get user ID from context (from JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.logger.ErrorContext(c.Request.Context(), "User ID not found in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInvalidToken, "User authentication required"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user ID type in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, "Invalid user authentication"))
		return
	}

	// Call service
	result, err := ctrl.channelService.UpdateWebhookChannel(c.Request.Context(), code, &req, userIDUint)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to update webhook notification channel",
			logger.String("code", code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel updated successfully",
		logger.String("code", code))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.update_success"), result)
}

// DeleteChannel delete notification channel
// @Summary Delete notification channel
// @Description Delete an existing notification channel by code
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Channel code"
// @Success 204 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/{code} [delete]
func (ctrl *NotificationChannelController) DeleteChannel(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		ctrl.logger.WarnContext(c.Request.Context(), "Channel code is required")
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, "Channel code is required"))
		return
	}

	// Get user ID from context (from JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.logger.ErrorContext(c.Request.Context(), "User ID not found in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInvalidToken, "User authentication required"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid user ID type in context")
		errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, "Invalid user authentication"))
		return
	}

	// Call service
	err := ctrl.channelService.DeleteChannel(c.Request.Context(), code, userIDUint)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to delete notification channel",
			logger.String("code", code), logger.ErrorField(err))
		errors.HandleError(c, err)
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
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/test/email [post]
func (ctrl *NotificationChannelController) TestEmailChannel(c *gin.Context) {
	var req request.TestEmailChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Call service
	result, err := ctrl.channelService.TestEmailChannel(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to test email notification channel",
			logger.String("code", req.Code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Email notification channel tested",
		logger.String("code", req.Code),
		logger.Bool("success", result.Success))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.test_complete"), result)
}

// TestWebhookChannel test webhook notification channel
// @Summary Test webhook notification channel
// @Description Test a webhook notification channel by sending a test message
// @Tags Notification Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TestWebhookChannelRequest true "Test webhook configuration"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/channels/test/webhook [post]
func (ctrl *NotificationChannelController) TestWebhookChannel(c *gin.Context) {
	var req request.TestWebhookChannelRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WarnContext(c.Request.Context(), "Invalid request body", logger.ErrorField(err))
		errors.HandleError(c, errors.NewAppError(errors.CodeValidationFailed, err.Error()))
		return
	}

	// Call service
	result, err := ctrl.channelService.TestWebhookChannel(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Failed to test webhook notification channel",
			logger.String("code", req.Code), logger.ErrorField(err))
		errors.HandleError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Webhook notification channel tested",
		logger.String("code", req.Code),
		logger.Bool("success", result.Success))

	response.Success(c, ctrl.i18n.T(c.Request.Context(), "notification.channel.test_complete"), result)
}
