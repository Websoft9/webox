package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/logger"

	"github.com/go-playground/validator/v10"
)

// NotificationTemplateController handles notification template HTTP requests
type NotificationTemplateController struct {
	templateService interfaceService.NotificationTemplateService
	validator       *validator.Validate
	logger          logger.Logger
}

// NewNotificationTemplateController creates a new notification template controller
func NewNotificationTemplateController(
	templateService interfaceService.NotificationTemplateService,
	validator *validator.Validate,
	logger logger.Logger,
) *NotificationTemplateController {
	return &NotificationTemplateController{
		templateService: templateService,
		validator:       validator,
		logger:          logger,
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
// @Success 201 {object} common.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates [post]
func (ctrl *NotificationTemplateController) CreateTemplate(c *gin.Context) {
	var req request.CreateNotificationTemplateRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.templateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template created successfully",
		logger.String("name", req.Name),
		logger.Uint("template_id", result.ID))

	response.SuccessWithData(c, result)
}

// GetTemplate get notification template by ID
// @Summary Get notification template by ID
// @Description Get detailed information about a specific notification template
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 200 {object} common.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates/{id} [get]
func (ctrl *NotificationTemplateController) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	result, err := ctrl.templateService.GetTemplate(c.Request.Context(), uint(id))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
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
// @Success 200 {object} common.APIResponse{data=[]response.NotificationTemplateResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates [get]
func (ctrl *NotificationTemplateController) GetTemplateList(c *gin.Context) {
	var req request.GetNotificationTemplateListRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	result, err := ctrl.templateService.GetTemplateList(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
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
// @Success 200 {object} common.APIResponse{data=response.NotificationTemplateResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates/{id} [put]
func (ctrl *NotificationTemplateController) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	var req request.UpdateNotificationTemplateRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.templateService.UpdateTemplate(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template updated successfully",
		logger.Uint("template_id", uint(id)))

	response.SuccessWithData(c, result)
}

// DeleteTemplate delete notification template
// @Summary Delete notification template
// @Description Delete an existing notification template by ID
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Template ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates/{id} [delete]
func (ctrl *NotificationTemplateController) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	// Call service
	if err := ctrl.templateService.DeleteTemplate(c.Request.Context(), uint(id)); err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Notification template deleted successfully",
		logger.Uint("template_id", uint(id)))

	response.Success(c)
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
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/notifications/templates/{id}/test [post]
func (ctrl *NotificationTemplateController) TestTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	var req request.TestNotificationTemplateRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	err = ctrl.templateService.TestTemplate(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}
	response.Success(c)
}
