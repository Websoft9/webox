package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"

	"github.com/expr-lang/expr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AlertController handles alert related requests
type AlertController struct {
	alertService service.AlertService
	logger       logger.Logger
	i18n         *i18n.I18n
	validator    *validator.Validate
}

// NewAlertController creates a new alert controller
func NewAlertController(
	alertService service.AlertService,
	logger logger.Logger,
	i18n *i18n.I18n,
	validator *validator.Validate,
) *AlertController {
	return &AlertController{
		alertService: alertService,
		logger:       logger,
		i18n:         i18n,
		validator:    validator,
	}
}

// GetAlertRules handles getting alert rules list
// @Summary Get alert rules list
// @Description Get a list of alert rules with optional filtering
// @Tags Alert
// @Accept json
// @Produce json
// @Param page query int false "Page number, default 1"
// @Param page_size query int false "Page size, default 20"
// @Param rule_type query string false "Rule type filter"
// @Param target_type query string false "Target type filter"
// @Param is_enabled query bool false "Is enabled filter"
// @Param keyword query string false "Keyword for name search"
// @Success 200 {object} response.APIResponse{data=response.AlertRuleListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/rules [get]
func (c *AlertController) GetAlertRules(ctx *gin.Context) {
	var req request.AlertRuleQueryRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.alertService.ListAlertRules(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get alert rules", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert rule list successful")
	response.SuccessWithData(ctx, result)
}

// CreateAlertRule handles creating an alert rule
// @Summary Create alert rule
// @Description Create a new alert rule
// @Tags Alert
// @Accept json
// @Produce json
// @Param request body request.AlertRuleCreateRequest true "Alert rule create request"
// @Success 200 {object} response.APIResponse{data=response.AlertRuleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/rules [post]
func (c *AlertController) CreateAlertRule(ctx *gin.Context) {
	var req request.AlertRuleCreateRequest

	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Validate the condition expression format using regex
	if !isExpression(req.ConditionExpression) {
		c.logger.WarnContext(ctx, "Invalid condition expression format",
			logger.String("expression", req.ConditionExpression))
		response.BadRequest(ctx, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	result, err := c.alertService.CreateAlertRule(ctx, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create alert rule", logger.String("name", req.Name), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert rule create successful")
	response.SuccessWithData(ctx, result)
}

func isExpression(expression string) bool {
	_, err := expr.Compile(expression, expr.Env(nil))
	return err == nil
}

// GetAlertRule handles getting a single alert rule by ID
// @Summary Get alert rule by ID
// @Description Get details of a specific alert rule
// @Tags Alert
// @Accept json
// @Produce json
// @Param id path int true "Alert rule ID"
// @Success 200 {object} response.APIResponse{data=response.AlertRuleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/rules/{id} [get]
func (c *AlertController) GetAlertRule(ctx *gin.Context) {

	// Get alert rule ID from path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	result, err := c.alertService.GetAlertRuleByID(ctx, uint(id))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get alert rule",
			logger.Uint("id", uint(id)),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert rule retrieved successfully")
	response.SuccessWithData(ctx, result)
}

// UpdateAlertRule handles updating an alert rule
// @Summary Update alert rule
// @Description Update an existing alert rule
// @Tags Alert
// @Accept json
// @Produce json
// @Param id path int true "Alert rule ID"
// @Param request body request.AlertRuleUpdateRequest true "Alert rule update request"
// @Success 200 {object} response.APIResponse{data=response.AlertRuleResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/rules/{id} [put]
func (c *AlertController) UpdateAlertRule(ctx *gin.Context) {
	// Get role ID
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	var req request.AlertRuleUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Validate condition expression if it's provided in the update request
	if req.ConditionExpression != nil && !isExpression(*req.ConditionExpression) {
		c.logger.WarnContext(ctx, "Invalid condition expression format",
			logger.String("expression", *req.ConditionExpression))
		response.BadRequest(ctx, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	result, err := c.alertService.UpdateAlertRule(ctx, uint(id), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update alert rule",
			logger.Uint("id", uint(id)),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert rule updated successfully")
	response.SuccessWithData(ctx, result)
}

// DeleteAlertRule handles deleting an alert rule
// @Summary Delete alert rule
// @Description Delete an existing alert rule
// @Tags Alert
// @Accept json
// @Produce json
// @Param id path int true "Alert rule ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/rules/{id} [delete]
func (c *AlertController) DeleteAlertRule(ctx *gin.Context) {

	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	err := c.alertService.DeleteAlertRule(ctx, id)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete alert rule",
			logger.Uint("id", uint(id)),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert rule deleted successfully")
	response.Success(ctx)
}

// GetAlertRecords handles getting alert records list
// @Summary Get alert records list
// @Description Get a list of alert records with optional filtering
// @Tags Alert
// @Accept json
// @Produce json
// @Param page query int false "Page number, default 1"
// @Param page_size query int false "Page size, default 20"
// @Param status query string false "Status filter (FIRING, RESOLVED)"
// @Param severity query string false "Severity filter"
// @Param start_time query string false "Start time filter (ISO8601 format)"
// @Param end_time query string false "End time filter (ISO8601 format)"
// @Param alert_rule_id query int false "Alert rule ID filter"
// @Success 200 {object} response.APIResponse{data=response.AlertRecordListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/records [get]
func (c *AlertController) GetAlertRecords(ctx *gin.Context) {
	var req request.AlertRecordQueryRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.alertService.ListAlertRecords(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get alert records", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert records list successful")
	response.SuccessWithData(ctx, result)
}

// AcknowledgeAlertRecord handles acknowledging an alert record
// @Summary Acknowledge alert record
// @Description Acknowledge an existing alert record
// @Tags Alert
// @Accept json
// @Produce json
// @Param id path int true "Alert record ID"
// @Param request body request.AlertAcknowledgeRequest true "Alert acknowledge request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/records/{id}/acknowledge [put]
func (c *AlertController) AcknowledgeAlertRecord(ctx *gin.Context) {

	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	var req request.AlertAcknowledgeRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	err := c.alertService.AcknowledgeAlertRecord(ctx, uint(id), currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to acknowledge alert record",
			logger.Uint("id", uint(id)),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert record acknowledged successfully")
	response.Success(ctx)
}

// ResolveAlertRecord handles resolving an alert record
// @Summary Resolve alert record
// @Description Resolve an existing alert record
// @Tags Alert
// @Accept json
// @Produce json
// @Param id path int true "Alert record ID"
// @Param request body request.AlertResolveRequest true "Alert resolve request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/alert/records/{id}/resolve [put]
func (c *AlertController) ResolveAlertRecord(ctx *gin.Context) {
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	var req request.AlertResolveRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID from context
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	err := c.alertService.ResolveAlertRecord(ctx, uint(id), currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to resolve alert record",
			logger.Uint("id", uint(id)),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Alert record resolved successfully")
	response.Success(ctx)
}
