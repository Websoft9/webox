package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SystemConfigController struct {
	SystemConfigService service.SystemConfigService
	validator           *validator.Validate
	logger              logger.Logger
}

func NewSystemConfigController(
	systemConfigService service.SystemConfigService,
	validator *validator.Validate,
	logger logger.Logger,
) *SystemConfigController {
	return &SystemConfigController{
		SystemConfigService: systemConfigService,
		validator:           validator,
		logger:              logger,
	}
}

// ListSystemConfigs lists system configurations with pagination and filtering
// @Summary List system configurations
// @Description List system configurations with pagination and filtering
// @Tags System Config
// @Security BearerAuth
// @Param category query string false "Filter by category"
// @Param keyword query string false "Filter by keyword"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs [get]
func (c *SystemConfigController) ListSystemConfigs(ctx *gin.Context) {
	var req request.ListSystemConfigsRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	result, err := c.SystemConfigService.ListSystemConfigs(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, result)
}

// TestSMTP tests the SMTP configuration by sending a test email
// @Summary Test SMTP configuration
// @Description Test SMTP configuration by sending a test email
// @Tags System Config
// @Accept json
// @Param request body request.TestSMTPRequest true "Test email request"
// @Security BearerAuth
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs/smtp/test [post]
func (c *SystemConfigController) TestSMTP(ctx *gin.Context) {
	var req request.TestSMTPRequest
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.SystemConfigService.TestSMTP(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.Success(ctx)
}

// ListBasicConfigs lists basic category system configurations
// @Summary List basic system configurations
// @Description List basic category system configurations
// @Tags System Config
// @Security BearerAuth
// @Param keyword query string false "Filter by keyword"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs/basic [get]
func (c *SystemConfigController) ListBasicConfigs(ctx *gin.Context) {
	var req request.ListSystemConfigsRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}
	req.Category = "basic"

	result, err := c.SystemConfigService.ListSystemConfigs(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, result)
}

// ListSecurityConfigs lists security category system configurations
// @Summary List security system configurations
// @Description List security category system configurations
// @Tags System Config
// @Security BearerAuth
// @Param keyword query string false "Filter by keyword"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs/security [get]
func (c *SystemConfigController) ListSecurityConfigs(ctx *gin.Context) {
	var req request.ListSystemConfigsRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	req.Category = "security"

	result, err := c.SystemConfigService.ListSystemConfigs(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, result)
}

// ListEmailConfigs lists email category system configurations
// @Summary List email system configurations
// @Description List email category system configurations
// @Tags System Config
// @Security BearerAuth
// @Param keyword query string false "Filter by keyword"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs/email [get]
func (c *SystemConfigController) ListEmailConfigs(ctx *gin.Context) {
	var req request.ListSystemConfigsRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	req.Category = "email"

	result, err := c.SystemConfigService.ListSystemConfigs(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, result)
}

// BatchUpdateSystemConfigs batch updates system configurations
// @Summary Batch update system configurations
// @Description Batch update system configurations
// @Tags System Config
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param configs body request.BatchUpdateSystemConfigsRequest true "Batch System Configuration Data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/system-configs [put]
func (c *SystemConfigController) BatchUpdateSystemConfigs(ctx *gin.Context) {
	var req request.BatchUpdateSystemConfigsRequest

	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	err := c.SystemConfigService.BatchUpdateSystemConfigs(ctx.Request.Context(), &req)
	if err != nil {
		response.WithError(ctx, err)
		return
	}

	response.Success(ctx)
}
