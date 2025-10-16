package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
)

const (
	// Content types
	contentTypeCSV   = "text/csv"
	contentTypeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	// Response limits
	maxErrorMessageLength = 200

	// Password masking constants
	maskedPassword       = "******"
	passwordField        = "password"
	newPasswordField     = "new_password"
	oldPasswordField     = "old_password"
	confirmPasswordField = "confirm_password"
)

// passwordFieldNames contains field names that should be masked
var passwordFieldNames = []string{passwordField, confirmPasswordField, oldPasswordField, newPasswordField}

// auditLogService audit log service implementation
type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
	userService  service.UserService
	db           *gorm.DB
	logger       logger.Logger
	config       *config.Config
}

// NewAuditLogService creates audit log service instance
func NewAuditLogService(
	auditLogRepo repository.AuditLogRepository,
	userService service.UserService,
	db *gorm.DB,
	logger logger.Logger,
	config *config.Config,
) service.AuditLogService {
	return &auditLogService{
		auditLogRepo: auditLogRepo,
		userService:  userService,
		db:           db,
		logger:       logger,
		config:       config,
	}
}

// RecordLog records audit log (internal use)
func (s *auditLogService) RecordLog(ctx context.Context, req *request.CreateAuditLogRequest) error {
	if req == nil {
		s.logger.ErrorContext(ctx, "Audit log request is required")
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
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
		return err
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
		return nil, err
	}

	s.logger.InfoContext(ctx, "Audit log retrieved successfully", logger.Uint("audit_log_id", id))

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
func (s *auditLogService) ListAuditLogs(ctx context.Context, req *request.ListAuditLogRequest) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Listing audit logs",
		logger.String("service", "audit_log"),
		logger.String("operation", "ListAuditLogs"),
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize))

	auditLogs, total, err := s.auditLogRepo.List(ctx, req)
	if err != nil {
		return nil, err
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
	return common.NewPaginationResponse(
		req.GetPage(),
		req.GetPageSize(),
		total,
		responses,
	), nil
}

// ExportAuditLogs exports audit logs in specified format
func (s *auditLogService) ExportAuditLogs(ctx context.Context, ginCtx *gin.Context, req *request.ExportAuditLogRequest) (data []byte, contentType string, err error) {
	s.logger.InfoContext(ctx, "Exporting audit logs",
		logger.String("service", "audit_log"),
		logger.String("operation", "ExportAuditLogs"),
		logger.String("format", req.Format))

	// Parse and validate time range using TimeRangeRequest
	timeRange := &common.TimeRangeRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	startTime, endTime, err := timeRange.GetTimeRange(false)
	if err != nil {
		s.logger.ErrorContext(ctx, "Invalid time range parameters", logger.ErrorField(err))
		return nil, "", errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Handle user ID based on permissions and request
	userID := s.resolveExportUserID(ginCtx, req.UserID)

	// Build query conditions using QueryBuilder pattern
	queryBuilder := s.buildExportQueryBuilder(userID, startTime, endTime)

	// Determine export format
	format := strings.ToLower(req.Format)
	if format == "" {
		format = constants.FormatCSV // default format
	}

	var exportFormat utils.ExportFormat
	switch format {
	case constants.FormatJSON:
		exportFormat = utils.FormatJSON
		contentType = "application/json"
	case constants.FormatExcel:
		exportFormat = utils.FormatExcel
		contentType = contentTypeExcel
	case constants.FormatCSV:
		exportFormat = utils.FormatCSV
		contentType = contentTypeCSV
	default:
		s.logger.WarnContext(ctx, "Unsupported export format, using CSV", logger.String("format", format))
		exportFormat = utils.FormatCSV
		contentType = contentTypeCSV
	}

	// Create export configuration
	exporter := utils.NewDBExporter(s.db)
	config := utils.NewExportConfigBuilder().
		TableName("audit_logs").
		Fields("*").
		QueryBuilder(queryBuilder).
		Format(exportFormat).
		OrderBy("created_at DESC").
		Build()

	// Export data using DBExporter
	exportData, err := exporter.Export(config)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to export audit logs", logger.ErrorField(err))
		return nil, "", errors.NewAppErrorWrapError(err, errors.CodeRecordExportFailed)
	}

	s.logger.InfoContext(ctx, "Audit logs exported successfully",
		logger.String("format", format),
		logger.Int("data_size", len(exportData)))

	return exportData, contentType, nil
}

// buildExportQueryBuilder builds query conditions for export
func (s *auditLogService) buildExportQueryBuilder(userID *uint, startTime, endTime time.Time) utils.QueryBuilder {
	builders := []utils.QueryBuilder{
		utils.WhereBetween("created_at", startTime, endTime),
	}

	// Add user ID filter if specified
	if userID != nil {
		builders = append(builders, utils.WhereEqual("user_id", *userID))
	}

	return utils.CombineQueryBuilders(builders...)
}

// resolveExportUserID resolves the user ID for export - if not specified, use current user
func (s *auditLogService) resolveExportUserID(ginCtx *gin.Context, requestUserID *uint) *uint {
	// If user ID is explicitly provided in request, use it
	if requestUserID != nil {
		return requestUserID
	}

	// Fallback: try to get user ID from context (legacy support)
	if userIDInterface, exists := ginCtx.Get("user_id"); exists {
		if currentUserID, ok := userIDInterface.(uint); ok {
			return &currentUserID
		}
	}

	// Return nil if no user context (export all)
	return nil
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

// ShouldSkipAudit determines whether to skip audit recording based on method and path
func (s *auditLogService) ShouldSkipAudit(method, path string) bool {
	return s.shouldSkip(path) || s.shouldSkipByMethod(method, path)
}

// shouldSkip determines whether to skip audit recording for specific paths
func (s *auditLogService) shouldSkip(path string) bool {
	// Extract path without query parameters for accurate matching
	pathWithoutQuery := s.extractPathWithoutQuery(path)

	skipPaths := s.config.AuditLog.SkipPaths

	for _, skipPath := range skipPaths {
		// Support both exact match and prefix match with trailing slash
		if pathWithoutQuery == skipPath || strings.HasPrefix(pathWithoutQuery, skipPath+"/") {
			return true
		}
	}

	return false
}

// shouldSkipByMethod determines whether to skip audit recording based on HTTP method
func (s *auditLogService) shouldSkipByMethod(method, path string) bool {
	// Extract path without query parameters for accurate matching
	pathWithoutQuery := s.extractPathWithoutQuery(path)

	// Always record modification operations regardless of path
	modificationMethods := s.config.AuditLog.AuditMethods
	for _, m := range modificationMethods {
		if method == m {
			return false // Don't skip - record modification operations
		}
	}

	// For GET requests, only record specific sensitive operations
	if method == "GET" {
		// Record sensitive authentication/authorization queries
		sensitiveGetPaths := s.config.AuditLog.SensitiveGetPaths

		for _, sensitivePath := range sensitiveGetPaths {
			// Handle path parameters like /:id/
			normalizedPath := strings.ReplaceAll(sensitivePath, "/:id/", "/")
			normalizedPath = strings.ReplaceAll(normalizedPath, "/:id", "")

			// Check for exact match or path with parameters
			if pathWithoutQuery == sensitivePath ||
				pathWithoutQuery == normalizedPath ||
				strings.HasPrefix(pathWithoutQuery, normalizedPath+"/") ||
				(strings.Contains(sensitivePath, "/:id") && strings.HasPrefix(pathWithoutQuery, normalizedPath)) {
				return false // Don't skip - record sensitive queries
			}
		}

		// Skip all other GET requests (regular queries)
		return true
	}

	// Skip configured skip methods
	skipMethods := s.config.AuditLog.SkipMethods
	for _, skipMethod := range skipMethods {
		if method == skipMethod {
			return true
		}
	}

	return false // Record by default for other methods
}

// extractPathWithoutQuery removes query parameters from URL path
func (s *auditLogService) extractPathWithoutQuery(path string) string {
	if queryIndex := strings.Index(path, "?"); queryIndex != -1 {
		return path[:queryIndex]
	}
	return path
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

// extractAuditInfoFromRequest extracts audit information from HTTP request context
func (s *auditLogService) extractAuditInfoFromRequest(ctx *gin.Context, responseBody []byte, responseTime int) *request.CreateAuditLogRequest {
	// Get user information
	var userID *uint
	var username string
	if userClaims, exists := ctx.Get("claims"); exists {
		if claims, ok := userClaims.(*auth.Claims); ok {
			userID = &claims.UserID
			username = claims.Username
		}
	}

	action := s.getActionFromMethod(ctx.Request.Method, ctx.Request.RequestURI)
	if action == constants.ActionLogin && username == "" {
		// For login operations, try to extract username from request body
		username = s.extractUsernameFromRequestBody(ctx)
	}

	// Determine module
	module := s.getModuleFromPath(ctx.Request.RequestURI)
	module_type := constants.GetModuleType(module)

	// Build description
	description := s.buildDescription(ctx.Request.Method, ctx.Request.RequestURI, ctx.Writer.Status())

	// Get request parameters
	requestParams := s.getRequestParams(ctx)

	// Determine if operation was successful
	statusCode := ctx.Writer.Status()
	success := statusCode >= constants.StatusOK && statusCode < constants.StatusBadRequest

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
		ResourceType:   module_type,
		ResourceID:     s.extractResourceID(ctx),
		ResourceName:   constants.GetModuleTableName(module_type),
		Description:    description,
		IPAddress:      utils.GetRealIP(ctx),
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
	// Check for specific auth operations
	switch {
	case strings.HasSuffix(path, "/auth/login") || path == "/auth/login":
		return constants.ActionLogin
	case strings.HasSuffix(path, "/auth/logout") || path == "/auth/logout":
		return constants.ActionLogout
	case strings.HasSuffix(path, "/auth/register") || path == "/auth/register":
		return constants.ActionRegister
	case strings.HasSuffix(path, "/auth/refresh") || path == "/auth/refresh":
		return constants.ActionRefresh
	}

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

// buildDescription builds operation description
func (s *auditLogService) buildDescription(method, path string, statusCode int) string {
	action := s.getActionFromMethod(method, path)

	if statusCode >= 200 && statusCode < 300 {
		return action + " success"
	}
	return action + " failed"
}

// getModuleFromPath extracts module from URL path in format /api/v1/{module}
func (s *auditLogService) getModuleFromPath(path string) string {
	// Remove query parameters first
	pathWithoutQuery := s.extractPathWithoutQuery(path)
	pathSegments := strings.Split(strings.Trim(pathWithoutQuery, "/"), "/")

	// Expected format: /api/v1/{module}/...
	if len(pathSegments) < constants.MinPathSegments {
		return "unknown"
	}
	return pathSegments[2]
}

// getRequestParams gets request parameters with special handling for exports
func (s *auditLogService) getRequestParams(ctx *gin.Context) string {
	params := make(map[string]interface{})
	s.addQueryParams(params, ctx)
	s.addPathParams(params, ctx)
	if ctx.Request.Method == constants.HTTPMethodPOST {
		s.addPostParams(params, ctx)
	}
	return s.paramsToJSON(params)
}

// addQueryParams adds query parameters to the params map
func (s *auditLogService) addQueryParams(params map[string]interface{}, ctx *gin.Context) {
	for key, values := range ctx.Request.URL.Query() {
		params[key] = s.getValueFromSlice(values)
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
		} else {
			params[key] = s.getValueFromSlice(values)
		}
	}
}

// getValueFromSlice returns single value or slice based on length
func (s *auditLogService) getValueFromSlice(values []string) interface{} {
	if len(values) == 1 {
		return values[0]
	}
	return values
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
	for _, field := range passwordFieldNames {
		if key == field {
			return true
		}
	}
	return false
}

// paramsToJSON converts params map to JSON string
func (s *auditLogService) paramsToJSON(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	if jsonBytes, err := json.Marshal(params); err == nil {
		return string(jsonBytes)
	}

	return ""
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
	if id := ctx.Param("id"); id != "" {
		if parsedID, err := strconv.ParseUint(id, 10, 32); err == nil {
			resourceID := uint(parsedID)
			return &resourceID
		}
	}
	return nil
}

// extractUsernameFromRequestBody extracts username from JSON request body only
func (s *auditLogService) extractUsernameFromRequestBody(ctx *gin.Context) string {
	// Extract from JSON body stored in context by middleware
	if requestBodyInterface, exists := ctx.Get("audit_request_body"); exists {
		if requestBody, ok := requestBodyInterface.([]byte); ok && len(requestBody) > 0 {
			var loginReq struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			}
			if err := json.Unmarshal(requestBody, &loginReq); err == nil {
				if loginReq.Username != "" {
					return loginReq.Username
				}
				if loginReq.Email != "" {
					return loginReq.Email
				}
			}
		}
	}
	return ""
}
