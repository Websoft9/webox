package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

// NotificationTemplateController handles notification template HTTP requests
type NotificationTemplateController struct {
	templateService interfaceService.NotificationTemplateService
	logger          logger.Logger
	i18n            *i18n.I18n
}

// NewNotificationTemplateController creates a new notification template controller
func NewNotificationTemplateController(
	templateService interfaceService.NotificationTemplateService,
	logger logger.Logger,
	i18n *i18n.I18n,
) *NotificationTemplateController {
	return &NotificationTemplateController{
		templateService: templateService,
		logger:          logger,
		i18n:            i18n,
	}
}

// CreateTemplate create notification template
// @Summary Create notification template
// @Description Create a new notification template with title, content and variables
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateNotificationTemplateRequest true "Notification template configuration"
// @Success 201 {object} response.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates [post]
func (ctrl *NotificationTemplateController) CreateTemplate(c *gin.Context) {
	var req request.CreateNotificationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseBadRequest(c, err, "validation.invalid_request_body", ctrl.i18n)
		return
	}

	// Call service
	result, err := ctrl.templateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template created successfully",
		logger.String("name", req.Name),
		logger.Uint("template_id", result.ID))

	ResponseOKWithData(c, result, "notification.template.create_success", ctrl.i18n)
}

// GetTemplate get notification template by ID
// @Summary Get notification template by ID
// @Description Get detailed information about a specific notification template
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 200 {object} response.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates/{id} [get]
func (ctrl *NotificationTemplateController) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ResponseBadRequest(c, err, "validation.invalid_id_format", ctrl.i18n)
		return
	}

	result, err := ctrl.templateService.GetTemplate(c.Request.Context(), uint(id))
	if err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	ResponseOKWithData(c, result, "notification.template.get_success", ctrl.i18n)
}

// GetTemplateList get notification templates list with pagination
// @Summary Get notification templates list
// @Description Get paginated list of notification templates with filtering
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keyword query string false "Search keyword"
// @Param template_type query string false "Template type filter" Enums(EMAIL,WEBHOOK,INTERNAL)
// @Param status query int false "Template status filter" Enums(0,1)
// @Param is_system query int false "System template filter" Enums(0,1)
// @Success 200 {object} response.APIResponse{data=response.NotificationTemplateListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates [get]
func (ctrl *NotificationTemplateController) GetTemplateList(c *gin.Context) {
	var req request.GetNotificationTemplateListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ResponseBadRequest(c, err, "validation.invalid_request_parameters", ctrl.i18n)
		return
	}

	result, err := ctrl.templateService.GetTemplateList(c.Request.Context(), &req)
	if err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	ResponseOKWithData(c, result, "notification.template.list_success", ctrl.i18n)
}

// UpdateTemplate update notification template
// @Summary Update notification template
// @Description Update an existing notification template configuration
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Param request body request.UpdateNotificationTemplateRequest true "Updated notification template configuration"
// @Success 200 {object} response.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates/{id} [put]
func (ctrl *NotificationTemplateController) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ResponseBadRequest(c, err, "validation.invalid_id_format", ctrl.i18n)
		return
	}

	var req request.UpdateNotificationTemplateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		ResponseBadRequest(c, bindErr, "validation.invalid_request_body", ctrl.i18n)
		return
	}

	// Call service
	result, err := ctrl.templateService.UpdateTemplate(c.Request.Context(), uint(id), &req)
	if err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template updated successfully",
		logger.Uint("template_id", uint(id)))

	ResponseOKWithData(c, result, "notification.template.update_success", ctrl.i18n)
}

// DeleteTemplate delete notification template
// @Summary Delete notification template
// @Description Delete an existing notification template by ID
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates/{id} [delete]
func (ctrl *NotificationTemplateController) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ResponseBadRequest(c, err, "validation.invalid_id_format", ctrl.i18n)
		return
	}

	// Call service
	if err := ctrl.templateService.DeleteTemplate(c.Request.Context(), uint(id)); err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template deleted successfully",
		logger.Uint("template_id", uint(id)))

	ResponseOK(c, "notification.template.delete_success", ctrl.i18n)
}

// TestTemplate test notification template
// @Summary Test notification template
// @Description Test a notification template by sending a test notification
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Param request body request.TestNotificationTemplateRequest true "Test template configuration"
// @Success 200 {object} response.APIResponse{data=response.NotificationTemplateTestResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/notifications/templates/{id}/test [post]
func (ctrl *NotificationTemplateController) TestTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ResponseBadRequest(c, err, "validation.invalid_id_format", ctrl.i18n)
		return
	}

	var req request.TestNotificationTemplateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		ResponseBadRequest(c, bindErr, "validation.invalid_request_body", ctrl.i18n)
		return
	}

	result, err := ctrl.templateService.TestTemplate(c.Request.Context(), uint(id), &req)
	if err != nil {
		ResponseWithError(c, err, ctrl.logger, ctrl.i18n)
		return
	}

	if result.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": ctrl.i18n.T(c.Request.Context(), "notification.template.test_success"),
			"data":    result,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": result.Message,
			"data":    result,
		})
	}
}
