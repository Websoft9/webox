package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

const (
	// Export formats
	exportFormatCSV = "csv"

	// System constants
	systemModuleName = "System"

	// Path parsing constants
	minPathSegments = 3

	// Response limits
	maxErrorMessageLength = 200

	// Audit log display constants
	statusSuccess = "true"
	statusFailure = "false"

	// Password masking constants
	maskedPassword       = "******"
	passwordField        = "password"
	newPasswordField     = "new_password"
	oldPasswordField     = "old_password"
	confirmPasswordField = "confirm_password"

	// Excel styling constants
	excelHeaderFontSize   = 12
	excelRowStartIndex    = 2
	successRateMultiplier = 100

	// Column widths for Excel export
	colWidthID           = 8
	colWidthUserID       = 10
	colWidthUsername     = 15
	colWidthAction       = 12
	colWidthModule       = 12
	colWidthResourceType = 15
	colWidthResourceID   = 10
	colWidthResourceName = 20
	colWidthDescription  = 30
	colWidthIP           = 15
	colWidthUserAgent    = 25
	colWidthMethod       = 10
	colWidthURL          = 30
	colWidthParams       = 25
	colWidthStatus       = 10
	colWidthTime         = 12
	colWidthSuccess      = 8
	colWidthError        = 20
	colWidthCreatedAt    = 20
)

// auditLogService audit log service implementation
type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
	userService  service.UserService
	db           *gorm.DB
	logger       logger.Logger
	i18n         *i18n.I18n
}

// NewAuditLogService creates audit log service instance
func NewAuditLogService(
	auditLogRepo repository.AuditLogRepository,
	userService service.UserService,
	db *gorm.DB,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.AuditLogService {
	return &auditLogService{
		auditLogRepo: auditLogRepo,
		userService:  userService,
		db:           db,
		logger:       logger,
		i18n:         i18n,
	}
}

// RecordLog records audit log (internal use)
func (s *auditLogService) RecordLog(ctx context.Context, req *request.CreateAuditLogRequest) error {
	if req == nil {
		s.logger.ErrorContext(ctx, "Audit log request is required")
		return errors.New("audit log request is required")
	}

	s.logger.InfoContext(ctx, "Recording audit log",
		logger.String("service", "audit_log"),
		logger.String("operation", "RecordLog"),
		logger.String("action", req.Action),
		logger.String("module", req.Module))

	auditLog := &model.AuditLog{
		UserID:         req.UserID,
		Username:       req.Username,
		Action:         req.Action,
		Module:         req.Module,
		ResourceType:   req.ResourceType,
		ResourceID:     req.ResourceID,
		ResourceName:   req.ResourceName,
		Description:    req.Description,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		RequestMethod:  req.RequestMethod,
		RequestURL:     req.RequestURL,
		RequestParams:  req.RequestParams,
		ResponseStatus: req.ResponseStatus,
		ResponseTime:   req.ResponseTime,
		Success:        req.Success,
		ErrorMessage:   req.ErrorMessage,
	}

	if err := s.auditLogRepo.Create(ctx, auditLog); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create audit log", logger.ErrorField(err))
		return errors.Wrap(err, "failed to record audit log")
	}

	s.logger.InfoContext(ctx, "Audit log recorded successfully",
		logger.Uint("audit_log_id", auditLog.ID))
	return nil
}

// GetAuditLog gets audit log details by ID
func (s *auditLogService) GetAuditLog(ctx context.Context, id uint) (*response.AuditLogResponse, error) {
	s.logger.InfoContext(ctx, "Getting audit log",
		logger.String("service", "audit_log"),
		logger.String("operation", "GetAuditLog"),
		logger.Uint("audit_log_id", id))

	auditLog, err := s.auditLogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get audit log", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get audit log")
	}

	s.logger.InfoContext(ctx, "Audit log retrieved successfully",
		logger.Uint("audit_log_id", id))

	var resp response.AuditLogResponse
	resp.FromAuditLog(auditLog)

	// Fill complete user information
	if resp.User != nil && resp.User.ID > 0 {
		if userInfo, err := s.userService.GetUser(ctx, resp.User.ID); err == nil {
			resp.User.ID = userInfo.ID
			resp.User.Username = userInfo.Username
			resp.User.Nickname = userInfo.Nickname
		}
	}

	return &resp, nil
}

// ListAuditLogs lists audit logs with pagination and filters
func (s *auditLogService) ListAuditLogs(ctx context.Context, req *request.ListAuditLogRequest) (*response.AuditLogListResponse, error) {
	s.logger.InfoContext(ctx, "Listing audit logs",
		logger.String("service", "audit_log"),
		logger.String("operation", "ListAuditLogs"),
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize))

	// Set default values
	req.SetDefaults()

	// Build filter conditions
	filter := &repository.AuditLogFilter{
		Page:         req.Page,
		PageSize:     req.PageSize,
		UserID:       req.UserID,
		Action:       req.Action,
		Module:       req.Module,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		IPAddress:    req.IPAddress,
		Success:      req.Success,
	}

	auditLogs, total, err := s.auditLogRepo.List(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list audit logs", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to list audit logs")
	}

	// Convert to response structure and collect user IDs
	responses := make([]response.AuditLogResponse, 0, len(auditLogs))
	userIDs := make(map[uint]bool)

	for _, auditLog := range auditLogs {
		var resp response.AuditLogResponse
		resp.FromAuditLog(auditLog)
		responses = append(responses, resp)

		// Collect user IDs for batch fetching
		if auditLog.UserID != nil {
			userIDs[*auditLog.UserID] = true
		}
	}

	// Batch fetch user information
	userInfoMap := make(map[uint]*response.UserResponse)
	for userID := range userIDs {
		if userInfo, err := s.userService.GetUser(ctx, userID); err == nil {
			userInfoMap[userID] = userInfo
		}
	}

	// Fill user information
	for i := range responses {
		if responses[i].User != nil && responses[i].User.ID > 0 {
			if userInfo, exists := userInfoMap[responses[i].User.ID]; exists {
				responses[i].User.ID = userInfo.ID
				responses[i].User.Nickname = userInfo.Nickname
				responses[i].User.Username = userInfo.Username
			}
		}
	}

	// Build paginated response
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	listResp := &response.AuditLogListResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      responses,
	}

	s.logger.InfoContext(ctx, "Audit logs listed successfully",
		logger.Int("total_count", int(total)),
		logger.Int("returned_count", len(responses)))

	return listResp, nil
}

// GetStatistics gets audit log statistics
func (s *auditLogService) GetStatistics(ctx context.Context, req *request.AuditLogStatisticsRequest) (*response.AuditLogStatisticsResponse, error) {
	s.logger.InfoContext(ctx, "Getting audit log statistics",
		logger.String("service", "audit_log"),
		logger.String("operation", "GetStatistics"),
		logger.String("group_by", req.GroupBy))

	// Set default values
	req.SetDefaults()

	// Build statistics filter conditions
	filter := &repository.StatisticsFilter{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		GroupBy:   req.GroupBy,
	}

	stats, err := s.auditLogRepo.GetStatistics(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get audit log statistics", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get audit log statistics")
	}

	// Convert to response structure
	resp := &response.AuditLogStatisticsResponse{
		TotalOperations:   stats.TotalOperations,
		SuccessOperations: stats.SuccessOperations,
		FailedOperations:  stats.FailedOperations,
	}

	// Calculate success rate
	if stats.TotalOperations > 0 {
		resp.SuccessRate = float64(stats.SuccessOperations) / float64(stats.TotalOperations) * successRateMultiplier
	}

	// Convert user statistics
	for _, userStat := range stats.TopUsers {
		resp.TopUsers = append(resp.TopUsers, response.UserStatItem{
			UserID:         userStat.UserID,
			Username:       userStat.Username,
			OperationCount: userStat.OperationCount,
		})
	}

	// Convert action statistics
	for _, actionStat := range stats.TopActions {
		resp.TopActions = append(resp.TopActions, response.ActionStatItem{
			Action: actionStat.Action,
			Count:  actionStat.Count,
		})
	}

	// Convert timeline statistics
	for _, timelineStat := range stats.Timeline {
		resp.Timeline = append(resp.Timeline, response.TimelineStatItem{
			Date:  timelineStat.Date,
			Count: timelineStat.Count,
		})
	}

	s.logger.InfoContext(ctx, "Audit log statistics retrieved successfully",
		logger.Int("timeline_points", len(stats.Timeline)))

	return resp, nil
}

// ExportAuditLogs exports audit logs in specified format
func (s *auditLogService) ExportAuditLogs(ctx context.Context, ginCtx *gin.Context, req *request.ExportAuditLogRequest) (data []byte, contentType string, err error) {
	s.logger.InfoContext(ctx, "Exporting audit logs",
		logger.String("service", "audit_log"),
		logger.String("operation", "ExportAuditLogs"),
		logger.String("format", req.Format))

	// Set default values
	req.SetDefaults()

	// Build filter conditions
	filter := &repository.AuditLogFilter{
		UserID:       req.UserID,
		Action:       req.Action,
		Module:       req.Module,
		ResourceType: req.ResourceType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	}

	auditLogs, err := s.auditLogRepo.Export(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to export audit logs", logger.ErrorField(err))
		return nil, "", errors.Wrap(err, "failed to export audit logs")
	}

	// Sanitize logs for export
	sanitizedLogs := make([]*model.AuditLog, 0, len(auditLogs))
	for _, auditLog := range auditLogs {
		sanitizedLogs = append(sanitizedLogs, auditLog.SanitizeForExport())
	}

	// Export by format
	switch req.Format {
	case "json":
		return s.exportToJSON(sanitizedLogs)
	case "excel":
		return s.exportToExcel(ginCtx, sanitizedLogs)
	case "csv":
		return s.exportToCSV(ginCtx, sanitizedLogs)
	default:
		s.logger.ErrorContext(ctx, "Unsupported export format", logger.String("format", req.Format))
		return nil, "", errors.New("unsupported export format: " + req.Format)
	}
}

// CleanupExpiredLogs cleanup expired audit logs (system scheduled task)
func (s *auditLogService) CleanupExpiredLogs(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		retentionDays = 90 // Default retention 90 days
	}

	beforeDate := time.Now().AddDate(0, 0, -retentionDays)
	deletedCount, err := s.auditLogRepo.CleanupOldLogs(ctx, beforeDate)
	if err != nil {
		s.logger.Error("Failed to cleanup old audit logs",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "retention_days", Value: retentionDays})
		return 0, err
	}

	s.logger.Info("Old audit logs cleaned up successfully",
		logger.Field{Key: "deleted_count", Value: deletedCount},
		logger.Field{Key: "retention_days", Value: retentionDays},
		logger.Field{Key: "before_date", Value: beforeDate.Format("2006-01-02")})

	return deletedCount, nil
}

// LogUserAction records user operation (convenience method)
func (s *auditLogService) LogUserAction(
	ctx context.Context,
	userID *uint,
	username, action, module, description, ipAddress, userAgent string,
	success bool,
	errorMsg string,
) error {
	req := &request.CreateAuditLogRequest{
		UserID:       userID,
		Username:     username,
		Action:       action,
		Module:       module,
		Description:  description,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Success:      success,
		ErrorMessage: errorMsg,
	}

	return s.RecordLog(ctx, req)
}

// getLocalizedHeaders returns internationalized headers based on the language from gin context
func (s *auditLogService) getLocalizedHeaders(ginCtx *gin.Context) []string {
	// Return localized headers using i18n service
	return []string{
		s.i18n.T(ginCtx, "audit_log.header_id"),
		s.i18n.T(ginCtx, "audit_log.header_user_id"),
		s.i18n.T(ginCtx, "audit_log.header_username"),
		s.i18n.T(ginCtx, "audit_log.header_action"),
		s.i18n.T(ginCtx, "audit_log.header_module"),
		s.i18n.T(ginCtx, "audit_log.header_resource_type"),
		s.i18n.T(ginCtx, "audit_log.header_resource_id"),
		s.i18n.T(ginCtx, "audit_log.header_resource_name"),
		s.i18n.T(ginCtx, "audit_log.header_description"),
		s.i18n.T(ginCtx, "audit_log.header_ip_address"),
		s.i18n.T(ginCtx, "audit_log.header_user_agent"),
		s.i18n.T(ginCtx, "audit_log.header_request_method"),
		s.i18n.T(ginCtx, "audit_log.header_request_url"),
		s.i18n.T(ginCtx, "audit_log.header_request_params"),
		s.i18n.T(ginCtx, "audit_log.header_response_status"),
		s.i18n.T(ginCtx, "audit_log.header_response_time"),
		s.i18n.T(ginCtx, "audit_log.header_success"),
		s.i18n.T(ginCtx, "audit_log.header_error_message"),
		s.i18n.T(ginCtx, "audit_log.header_created_at"),
	}
}

// exportToJSON export to JSON format
func (s *auditLogService) exportToJSON(auditLogs []*model.AuditLog) (data []byte, contentType string, err error) {
	// Convert to response format and marshal to JSON
	responses := make([]response.AuditLogResponse, 0, len(auditLogs))
	for _, auditLog := range auditLogs {
		var resp response.AuditLogResponse
		resp.FromAuditLog(auditLog)
		responses = append(responses, resp)
	}

	data, err = json.Marshal(responses)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return data, "application/json", nil
}

// exportToExcel export to Excel format
func (s *auditLogService) exportToExcel(ginCtx *gin.Context, auditLogs []*model.AuditLog) (data []byte, contentType string, err error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			s.logger.Error("Failed to close Excel file", logger.Field{Key: "error", Value: closeErr})
		}
	}()

	sheetName := "审计日志"
	if setupErr := s.setupExcelSheet(f, sheetName); setupErr != nil {
		return nil, "", setupErr
	}

	if headerErr := s.writeExcelHeaders(f, sheetName, ginCtx); headerErr != nil {
		return nil, "", headerErr
	}

	if dataErr := s.writeExcelData(f, sheetName, auditLogs); dataErr != nil {
		return nil, "", dataErr
	}

	s.formatExcelSheet(f, sheetName, auditLogs)

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("failed to write Excel file to buffer: %w", err)
	}

	return buffer.Bytes(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

// setupExcelSheet creates and sets up the Excel sheet
func (s *auditLogService) setupExcelSheet(f *excelize.File, sheetName string) error {
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create worksheet: %w", err)
	}

	// Set the created sheet as the active sheet
	f.SetActiveSheet(index)
	return nil
}

// writeExcelHeaders writes headers to the Excel sheet
func (s *auditLogService) writeExcelHeaders(f *excelize.File, sheetName string, ginCtx *gin.Context) error {
	// Get internationalized headers
	headers := s.getLocalizedHeaders(ginCtx)

	// Set headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if setErr := f.SetCellValue(sheetName, cell, header); setErr != nil {
			return fmt.Errorf("failed to set header %s: %w", header, setErr)
		}
	}

	return s.applyHeaderStyle(f, sheetName)
}

// applyHeaderStyle applies styling to Excel headers
func (s *auditLogService) applyHeaderStyle(f *excelize.File, sheetName string) error {
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: excelHeaderFontSize,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E6E6FA"}, // Light purple background
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}

	// Apply header style
	if styleErr := f.SetRowStyle(sheetName, 1, 1, headerStyle); styleErr != nil {
		return fmt.Errorf("failed to apply header style: %w", styleErr)
	}

	return nil
}

// writeExcelData writes audit log data to Excel sheet
func (s *auditLogService) writeExcelData(f *excelize.File, sheetName string, auditLogs []*model.AuditLog) error {
	for i, log := range auditLogs {
		row := i + excelRowStartIndex // Start from row 2 (after header)
		values := s.prepareExcelRowData(log)

		for j, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+j, row)
			if cellErr := f.SetCellValue(sheetName, cell, value); cellErr != nil {
				return fmt.Errorf("failed to set cell %s: %w", cell, cellErr)
			}
		}
	}
	return nil
}

// prepareExcelRowData prepares a single row of data for Excel export
func (s *auditLogService) prepareExcelRowData(log *model.AuditLog) []interface{} {
	// Convert pointers to values for display
	userID := ""
	if log.UserID != nil {
		userID = fmt.Sprintf("%d", *log.UserID)
	}

	resourceID := ""
	if log.ResourceID != nil {
		resourceID = fmt.Sprintf("%d", *log.ResourceID)
	}

	responseStatus := ""
	if log.ResponseStatus != nil {
		responseStatus = fmt.Sprintf("%d", *log.ResponseStatus)
	}

	responseTime := ""
	if log.ResponseTime != nil {
		responseTime = fmt.Sprintf("%d", *log.ResponseTime)
	}

	success := statusFailure
	if log.Success {
		success = statusSuccess
	}

	return []interface{}{
		log.ID,            // A: ID
		userID,            // B: 用户ID
		log.Username,      // C: 用户名
		log.Action,        // D: 操作
		log.Module,        // E: 模块
		log.ResourceType,  // F: 资源类型
		resourceID,        // G: 资源ID
		log.ResourceName,  // H: 资源名称
		log.Description,   // I: 描述
		log.IPAddress,     // J: IP地址
		log.UserAgent,     // K: 用户代理
		log.RequestMethod, // L: 请求方法
		log.RequestURL,    // M: 请求URL
		log.RequestParams, // N: 请求参数
		responseStatus,    // O: 响应状态
		responseTime,      // P: 响应时间
		success,           // Q: 成功
		log.ErrorMessage,  // R: 错误信息
		log.CreatedAt.Format("2006-01-02 15:04:05"), // S: 创建时间
	}
}

// formatExcelSheet applies formatting to the Excel sheet
func (s *auditLogService) formatExcelSheet(f *excelize.File, sheetName string, auditLogs []*model.AuditLog) {
	// Set column widths for better readability
	columnWidths := s.getExcelColumnWidths()

	for col, width := range columnWidths {
		if widthErr := f.SetColWidth(sheetName, col, col, width); widthErr != nil {
			s.logger.Warn("Failed to set column width", logger.Field{Key: "column", Value: col}, logger.Field{Key: "error", Value: widthErr})
		}
	}

	// Auto-filter for the data
	if len(auditLogs) > 0 {
		dataRange := fmt.Sprintf("A1:S%d", len(auditLogs)+1)
		if filterErr := f.AutoFilter(sheetName, dataRange, nil); filterErr != nil {
			s.logger.Warn("Failed to set auto filter", logger.Field{Key: "error", Value: filterErr})
		}
	}

	// Delete default sheet if it exists and is different from our sheet
	if f.GetSheetName(0) != sheetName {
		if deleteErr := f.DeleteSheet("Sheet1"); deleteErr != nil {
			s.logger.Warn("Failed to delete default sheet", logger.Field{Key: "error", Value: deleteErr})
		}
	}
}

// getExcelColumnWidths returns column width configuration for Excel export
func (s *auditLogService) getExcelColumnWidths() map[string]float64 {
	return map[string]float64{
		"A": colWidthID,           // ID
		"B": colWidthUserID,       // 用户ID
		"C": colWidthUsername,     // 用户名
		"D": colWidthAction,       // 操作
		"E": colWidthModule,       // 模块
		"F": colWidthResourceType, // 资源类型
		"G": colWidthResourceID,   // 资源ID
		"H": colWidthResourceName, // 资源名称
		"I": colWidthDescription,  // 描述
		"J": colWidthIP,           // IP地址
		"K": colWidthUserAgent,    // 用户代理
		"L": colWidthMethod,       // 请求方法
		"M": colWidthURL,          // 请求URL
		"N": colWidthParams,       // 请求参数
		"O": colWidthStatus,       // 响应状态
		"P": colWidthTime,         // 响应时间
		"Q": colWidthSuccess,      // 成功
		"R": colWidthError,        // 错误信息
		"S": colWidthCreatedAt,    // 创建时间
	}
}

// exportToCSV export to CSV format
func (s *auditLogService) exportToCSV(ginCtx *gin.Context, auditLogs []*model.AuditLog) (data []byte, contentType string, err error) {
	// Create CSV with proper escaping for special characters
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// Get internationalized headers
	headers := s.getLocalizedHeaders(ginCtx)

	// Write header
	if err := writer.Write(headers); err != nil {
		return nil, "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, log := range auditLogs {
		// Convert pointers to values for display
		userID := ""
		if log.UserID != nil {
			userID = fmt.Sprintf("%d", *log.UserID)
		}

		resourceID := ""
		if log.ResourceID != nil {
			resourceID = fmt.Sprintf("%d", *log.ResourceID)
		}

		responseStatus := ""
		if log.ResponseStatus != nil {
			responseStatus = fmt.Sprintf("%d", *log.ResponseStatus)
		}

		responseTime := ""
		if log.ResponseTime != nil {
			responseTime = fmt.Sprintf("%d", *log.ResponseTime)
		}

		success := statusFailure
		if log.Success {
			success = statusSuccess
		}

		// Create row data - same order as Excel export
		row := []string{
			fmt.Sprintf("%d", log.ID), // ID
			userID,                    // 用户ID
			log.Username,              // 用户名
			log.Action,                // 操作
			log.Module,                // 模块
			log.ResourceType,          // 资源类型
			resourceID,                // 资源ID
			log.ResourceName,          // 资源名称
			log.Description,           // 描述
			log.IPAddress,             // IP地址
			log.UserAgent,             // 用户代理
			log.RequestMethod,         // 请求方法
			log.RequestURL,            // 请求URL
			log.RequestParams,         // 请求参数
			responseStatus,            // 响应状态
			responseTime,              // 响应时间
			success,                   // 成功
			log.ErrorMessage,          // 错误信息
			log.CreatedAt.Format("2006-01-02 15:04:05"), // 创建时间
		}

		// Write row
		if err := writer.Write(row); err != nil {
			return nil, "", fmt.Errorf("failed to write CSV row for log ID %d: %w", log.ID, err)
		}
	}

	// Flush writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buffer.Bytes(), "text/csv; charset=utf-8", nil
}

// ShouldSkipAudit determines whether to skip audit recording based on method and path
func (s *auditLogService) ShouldSkipAudit(method, path string) bool {
	return s.shouldSkip(path) || s.shouldSkipByMethod(method, path)
}

// RecordAuditFromRequest records audit log from HTTP request context
func (s *auditLogService) RecordAuditFromRequest(backgroundCtx context.Context, ginCtx *gin.Context, responseBody []byte, responseTime int) error {
	// Extract audit information from request context
	auditReq := s.extractAuditInfoFromRequest(ginCtx, responseBody, responseTime)

	// For login requests, try to get username from request body if not available from user context
	if auditReq.Action == constants.ActionLogin && auditReq.Username == "" {
		if username := s.extractUsernameFromRequestBody(ginCtx); username != "" {
			auditReq.Username = username
		}
	}

	// Record audit log using existing method with background context
	return s.RecordLog(backgroundCtx, auditReq)
}

// shouldSkip determines whether to skip audit recording for specific paths
func (s *auditLogService) shouldSkip(path string) bool {
	skipPaths := []string{
		"/health",
		"/api/v1/health",
		"/metrics",
		"/api/v1/audit-logs", // Avoid recording audit log query operations themselves
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// shouldSkipByMethod determines whether to skip audit recording based on HTTP method
// Only record modification operations (POST, PUT, DELETE, PATCH) and specific sensitive operations
func (s *auditLogService) shouldSkipByMethod(method, path string) bool {
	// Always record modification operations regardless of path
	modificationMethods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, m := range modificationMethods {
		if method == m {
			return false // Don't skip - record modification operations
		}
	}

	// For GET requests, only record specific sensitive operations
	if method == "GET" {
		// Record export operations
		if strings.Contains(path, "/export") {
			return false // Don't skip - record exports
		}

		// Record sensitive authentication/authorization queries
		sensitiveGetPaths := []string{
			"/api/v1/users/profile",        // Profile access
			"/api/v1/api-tokens",           // API token listing
			"/api/v1/roles",                // Role listing
			"/api/v1/permissions",          // Permission listing
			"/api/v1/users/:id/two-factor", // 2FA status queries
		}

		for _, sensitivePath := range sensitiveGetPaths {
			if strings.Contains(path, strings.ReplaceAll(sensitivePath, "/:id/", "/")) {
				return false // Don't skip - record sensitive queries
			}
		}

		// Skip all other GET requests (regular queries)
		return true
	}

	// Skip OPTIONS, HEAD requests
	if method == "OPTIONS" || method == "HEAD" {
		return true
	}

	return false // Record by default for other methods
}

// extractAuditInfoFromRequest extracts audit information from HTTP request context
func (s *auditLogService) extractAuditInfoFromRequest(ctx *gin.Context, responseBody []byte, responseTime int) *request.CreateAuditLogRequest {
	// Get user information
	var userID *uint
	var username string
	if userClaims, exists := ctx.Get("user"); exists {
		if claims, ok := userClaims.(*auth.Claims); ok {
			userID = &claims.UserID
			username = claims.Username
		}
	}

	// For login operations, try to extract username from request body
	action := s.getActionFromMethod(ctx.Request.Method, ctx.Request.RequestURI)
	if action == constants.ActionLogin && username == "" {
		username = s.extractUsernameFromLoginRequest(ctx)
	}

	// Determine module
	module := s.getModuleFromPath(ctx.Request.RequestURI)

	// Build description
	description := s.buildDescription(ctx.Request.Method, ctx.Request.RequestURI, ctx.Writer.Status())

	// Get request parameters
	requestParams := s.getRequestParams(ctx)

	// Determine if operation was successful
	// Only 2xx status codes (200-299) are considered successful
	statusCode := ctx.Writer.Status()
	success := statusCode >= constants.StatusOK && statusCode < constants.StatusBadRequest

	// Debug: Print status and success for verification
	fmt.Printf("DEBUG: Status Code: %d, Success: %t\n", statusCode, success)

	// Get error message
	var errorMessage string
	if !success && len(responseBody) > 0 {
		errorMessage = s.extractErrorMessage(responseBody)
	}

	// Create audit log request with proper values
	return &request.CreateAuditLogRequest{
		UserID:         userID,
		Username:       username,
		Action:         action,
		Module:         module,
		ResourceType:   s.getResourceType(module),
		ResourceID:     s.extractResourceID(ctx),
		ResourceName:   s.getResourceTableName(module),
		Description:    description,
		IPAddress:      s.getClientIP(ctx),
		UserAgent:      ctx.Request.UserAgent(),
		RequestMethod:  ctx.Request.Method,
		RequestURL:     ctx.Request.RequestURI,
		RequestParams:  requestParams,
		ResponseStatus: &statusCode,
		ResponseTime:   &responseTime,
		Success:        success,
		ErrorMessage:   errorMessage,
	}
}

// getActionFromMethod determines action type based on HTTP method and path
func (s *auditLogService) getActionFromMethod(method, path string) string {
	// Check for export operations first
	if method == constants.HTTPMethodGET && strings.Contains(path, "/export") {
		return constants.ActionQuery
	}

	// Check for specific auth operations
	if method == constants.HTTPMethodPOST {
		switch {
		case strings.Contains(path, "/auth/login"):
			return constants.ActionLogin
		case strings.Contains(path, "/auth/logout"):
			return constants.ActionLogout
		case strings.Contains(path, "/auth/register"):
			return constants.ActionRegister
		case strings.Contains(path, "/auth/refresh"):
			return constants.ActionRefresh
		case strings.Contains(path, "/password"):
			return constants.ActionChangePassword
		}
	}

	// Default action based on HTTP method
	switch method {
	case constants.HTTPMethodPOST:
		return constants.ActionCreate
	case constants.HTTPMethodGET:
		return constants.ActionQuery
	case constants.HTTPMethodPUT, constants.HTTPMethodPATCH:
		return constants.ActionUpdate
	case constants.HTTPMethodDELETE:
		return constants.ActionDelete
	default:
		return constants.ActionQuery
	}
}

// getModuleFromPath determines module based on request path
func (s *auditLogService) getModuleFromPath(path string) string {
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")

	if len(pathSegments) < minPathSegments {
		return systemModuleName
	}

	// Skip "api" and "v1"
	module := pathSegments[2]

	moduleMap := map[string]string{
		"users":         constants.ModuleUser,
		"roles":         constants.ModuleRole,
		"permissions":   constants.ModulePermission,
		"auth":          constants.ModuleAuth,
		"api-tokens":    constants.ModuleAuth, // API tokens are part of auth module
		"two-factor":    constants.ModuleAuth, // 2FA is part of auth module
		"audit-logs":    constants.ModuleAuditLog,
		"auth-config":   constants.ModuleAuth,
		"i18n":          systemModuleName, // Keep system for i18n
		"applications":  constants.ModuleApplication,
		"projects":      constants.ModuleProject,
		"servers":       constants.ModuleServer,
		"databases":     constants.ModuleDatabase,
		"secrets":       constants.ModuleSecret,
		"certificates":  constants.ModuleCertificate,
		"workflows":     constants.ModuleWorkflow,
		"jobs":          constants.ModuleJob,
		"notifications": constants.ModuleNotification,
		"profile":       constants.ModuleProfile,
	}

	if moduleName, exists := moduleMap[module]; exists {
		return moduleName
	}

	return systemModuleName
}

// buildDescription builds operation description
func (s *auditLogService) buildDescription(method, path string, statusCode int) string {
	// Determine action
	action := s.getActionFromMethod(method, path)

	// Special handling for export operations
	if action == constants.ActionQuery {
		if statusCode >= 200 && statusCode < 300 {
			return "Export data successfully"
		} else {
			return "Export data failed"
		}
	}

	if statusCode >= 200 && statusCode < 300 {
		return action + " success"
	} else {
		return action + " failed"
	}
}

// getRequestParams gets request parameters with special handling for exports
func (s *auditLogService) getRequestParams(ctx *gin.Context) string {
	params := make(map[string]interface{})

	// Get query parameters
	s.addQueryParams(params, ctx)

	// Get path parameters
	s.addPathParams(params, ctx)

	// Handle POST request body data
	if ctx.Request.Method == constants.HTTPMethodPOST {
		s.addPostParams(params, ctx)
	}

	// Handle export operations
	s.addExportParams(params, ctx)

	// Convert to JSON string
	return s.paramsToJSON(params)
}

// addQueryParams adds query parameters to the params map
func (s *auditLogService) addQueryParams(params map[string]interface{}, ctx *gin.Context) {
	for key, values := range ctx.Request.URL.Query() {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}
}

// addPathParams adds path parameters to the params map
func (s *auditLogService) addPathParams(params map[string]interface{}, ctx *gin.Context) {
	for _, param := range ctx.Params {
		params[param.Key] = param.Value
	}
}

// addPostParams handles POST request form and JSON data
func (s *auditLogService) addPostParams(params map[string]interface{}, ctx *gin.Context) {
	// Try to parse form data
	if err := ctx.Request.ParseForm(); err == nil {
		s.addFormData(params, ctx)
	}

	// Also try to extract from JSON body stored in context
	s.addJSONBodyData(params, ctx)
}

// addFormData adds form data to params with password masking
func (s *auditLogService) addFormData(params map[string]interface{}, ctx *gin.Context) {
	for key, values := range ctx.Request.PostForm {
		if s.isPasswordField(key) {
			params[key] = maskedPassword
		} else if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}
}

// addJSONBodyData adds JSON body data to params with password masking
func (s *auditLogService) addJSONBodyData(params map[string]interface{}, ctx *gin.Context) {
	requestBodyInterface, exists := ctx.Get("audit_request_body")
	if !exists {
		return
	}

	requestBody, ok := requestBodyInterface.([]byte)
	if !ok || len(requestBody) == 0 {
		return
	}

	var jsonParams map[string]interface{}
	if err := json.Unmarshal(requestBody, &jsonParams); err == nil {
		for key, value := range jsonParams {
			if s.isPasswordField(key) {
				params[key] = maskedPassword
			} else {
				params[key] = value
			}
		}
	}
}

// isPasswordField checks if a field name is a password field
func (s *auditLogService) isPasswordField(key string) bool {
	passwordFields := []string{passwordField, confirmPasswordField, oldPasswordField, newPasswordField}
	for _, field := range passwordFields {
		if key == field {
			return true
		}
	}
	return false
}

// addExportParams adds export-specific metadata
func (s *auditLogService) addExportParams(params map[string]interface{}, ctx *gin.Context) {
	if !strings.Contains(ctx.Request.RequestURI, "/export") {
		return
	}

	// Add export-specific metadata
	if format, exists := params["format"]; exists {
		params["export_format"] = format
	} else {
		params["export_format"] = exportFormatCSV // default
	}

	// Record export time for audit purposes
	params["export_timestamp"] = time.Now().Format(time.RFC3339)
}

// paramsToJSON converts params map to JSON string
func (s *auditLogService) paramsToJSON(params map[string]interface{}) string {
	// Count of filter parameters (excluding export-specific ones)
	if strings.Contains(fmt.Sprintf("%v", params), "export") {
		filterCount := 0
		exportParams := map[string]bool{
			"format":           true,
			"page":             true,
			"page_size":        true,
			"export_format":    true,
			"export_timestamp": true,
		}

		for key := range params {
			if !exportParams[key] {
				filterCount++
			}
		}
		params["filter_count"] = filterCount
	}

	if len(params) == 0 {
		return ""
	}

	if jsonBytes, err := json.Marshal(params); err == nil {
		return string(jsonBytes)
	}

	return ""
}

// getClientIP gets client IP address
func (s *auditLogService) getClientIP(ctx *gin.Context) string {
	// Try to get real IP from proxy headers
	realIP := ctx.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Try to get IP from forwarded headers
	forwardedFor := ctx.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For may contain multiple IPs, take the first one
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Use RemoteAddr
	return ctx.ClientIP()
}

// extractErrorMessage extracts error message from response body
func (s *auditLogService) extractErrorMessage(responseBody []byte) string {
	var resp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(responseBody, &resp); err == nil {
		if resp.Error != "" {
			return resp.Error
		}
		if resp.Message != "" {
			return resp.Message
		}
	}

	// If unable to parse, return truncated response body
	if len(responseBody) > maxErrorMessageLength {
		return string(responseBody[:maxErrorMessageLength]) + "..."
	}
	return string(responseBody)
}

// extractResourceID extracts resource ID from URL parameters
func (s *auditLogService) extractResourceID(ctx *gin.Context) *uint {
	// Try to get resource ID from URL parameters (like /api/v1/users/123)
	if id := ctx.Param("id"); id != "" {
		if parsedID, err := strconv.ParseUint(id, 10, 32); err == nil {
			resourceID := uint(parsedID)
			return &resourceID
		}
	}
	return nil
}

// getResourceType returns the resource type (Chinese name) for a module
func (s *auditLogService) getResourceType(module string) string {
	return module
}

// getResourceTableName returns the database table name for a module
func (s *auditLogService) getResourceTableName(module string) string {
	return constants.GetModuleTableName(module)
}

// extractUsernameFromLoginRequest extracts username from login request body
func (s *auditLogService) extractUsernameFromLoginRequest(ctx *gin.Context) string {
	// Try to get username from request body for login operations
	if ctx.Request.Body == nil {
		return ""
	}

	// First try from form data (multipart form or URL encoded)
	if username := ctx.PostForm("username"); username != "" {
		return username
	}

	// Try from query parameters (less common for login)
	if username := ctx.Query("username"); username != "" {
		return username
	}

	// For JSON requests, we can't re-read the body here since it's already consumed
	// In a real implementation, you'd want to store this in middleware
	// For now, try to get from any available source

	return ""
}

// extractUsernameFromRequestBody extracts username from various request sources
func (s *auditLogService) extractUsernameFromRequestBody(ctx *gin.Context) string {
	// Try different sources for username

	// 1. Try from form data first
	if username := ctx.PostForm("username"); username != "" {
		return username
	}

	// 2. Try from query parameters
	if username := ctx.Query("username"); username != "" {
		return username
	}

	// 3. Try from JSON body stored in context by middleware
	if requestBodyInterface, exists := ctx.Get("audit_request_body"); exists {
		if requestBody, ok := requestBodyInterface.([]byte); ok && len(requestBody) > 0 {
			// Parse JSON to extract username
			var loginReq struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			}
			if err := json.Unmarshal(requestBody, &loginReq); err == nil {
				if loginReq.Username != "" {
					return loginReq.Username
				}
				// Some systems might use email as username
				if loginReq.Email != "" {
					return loginReq.Email
				}
			}
		}
	}

	return ""
}
