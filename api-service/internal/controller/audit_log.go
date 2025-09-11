package controller

import (
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
)

// AuditLogController audit log controller
type AuditLogController struct {
	auditLogService service.AuditLogService
	validator       *validator.Validate
	logger          logger.Logger
	i18n            *i18n.I18n
}

// NewAuditLogController creates audit log controller instance
func NewAuditLogController(
	auditLogService service.AuditLogService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *AuditLogController {
	return &AuditLogController{
		auditLogService: auditLogService,
		validator:       validator,
		logger:          logger,
		i18n:            i18n,
	}
}

// GetAuditLog gets single audit log details
// @Summary Get audit log details
// @Description Get audit log details by ID
// @Tags Audit Log
// @Security BearerAuth
// @Param id path int true "Audit Log ID"
// @Success 200 {object} response.AuditLogResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/audit-logs/{id} [get]
func (c *AuditLogController) GetAuditLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ResponseBadRequest(ctx, err, "validation.invalid_query_parameters", c.i18n)
		return
	}

	auditLog, err := c.auditLogService.GetAuditLog(ctx.Request.Context(), uint(id))
	if err != nil {
		ResponseNotFound(ctx, "common.not_found", c.i18n)
		return
	}

	ResponseOKWithData(ctx, auditLog, "common.success", c.i18n)
}

// ListAuditLogs get audit log list
// @Summary Get audit log list
// @Description Get paginated audit log list with multiple filter conditions
// @Tags Audit Log
// @Security BearerAuth
// @Param page query int false "Page Number" default(1)
// @Param page_size query int false "Page Size" default(20)
// @Param user_id query int false "User ID"
// @Param action query string false "Action Type"
// @Param resource_type query string false "Resource Type"
// @Param resource_id query int false "Resource ID"
// @Param start_time query string false "Start Time" format(date-time)
// @Param end_time query string false "End Time" format(date-time)
// @Param ip_address query string false "IP Address"
// @Success 200 {object} response.AuditLogListResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/audit-logs [get]
func (c *AuditLogController) ListAuditLogs(ctx *gin.Context) {
	var req request.ListAuditLogRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	result, err := c.auditLogService.ListAuditLogs(ctx.Request.Context(), &req)
	if err != nil {
		ResponseInternalError(ctx, err, "common.failed", c.logger, c.i18n)
		return
	}

	ResponseOKWithData(ctx, result, "common.success", c.i18n)
}

// GetAuditLogStatistics get audit log statistics
// @Summary Get audit log statistics
// @Description Get audit log statistical analysis data
// @Tags Audit Log
// @Security BearerAuth
// @Param start_time query string false "Start Time" format(date-time)
// @Param end_time query string false "End Time" format(date-time)
// @Param group_by query string false "Group By" Enums(hour,day,week,month) default(day)
// @Success 200 {object} response.AuditLogStatisticsResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/audit-logs/statistics [get]
func (c *AuditLogController) GetAuditLogStatistics(ctx *gin.Context) {
	var req request.AuditLogStatisticsRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	statistics, err := c.auditLogService.GetStatistics(ctx.Request.Context(), &req)
	if err != nil {
		ResponseInternalError(ctx, err, "common.failed", c.logger, c.i18n)
		return
	}

	ResponseOKWithData(ctx, statistics, "audit_log.success", c.i18n)
}

// ExportAuditLogs export audit logs
// @Summary Export audit logs
// @Description Export audit logs by conditions, supports multiple formats
// @Tags Audit Log
// @Security BearerAuth
// @Param format query string false "Export Format" Enums(csv,excel,json) default(csv)
// @Param start_time query string false "Start Time" format(date-time)
// @Param end_time query string false "End Time" format(date-time)
// @Param user_id query int false "User ID"
// @Success 200 {file} file "Export File"
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/audit-logs/export [get]
func (c *AuditLogController) ExportAuditLogs(ctx *gin.Context) {
	var req request.ExportAuditLogRequest

	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger, c.i18n) {
		return
	}

	data, contentType, err := c.auditLogService.ExportAuditLogs(ctx.Request.Context(), ctx, &req)
	if err != nil {
		ResponseInternalError(ctx, err, "common.failed", c.logger, c.i18n)
		return
	}

	// Set response headers
	filename := generateExportFilename(req.Format)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Length", strconv.Itoa(len(data)))

	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.Write(data)
}

// generateExportFilename generates filename with current date in format audit-logs-YYYYMMDD.ext
func generateExportFilename(format string) string {
	dateStr := time.Now().Format("20060102") // YYYYMMDD format
	extension := getFileExtensionByFormat(format)
	return fmt.Sprintf("audit-logs-%s%s", dateStr, extension)
}

// getFileExtensionByFormat returns the correct file extension for the given format
func getFileExtensionByFormat(format string) string {
	switch format {
	case "excel":
		return ".xlsx"
	case "json":
		return ".json"
	default:
		return ".csv"
	}
}
