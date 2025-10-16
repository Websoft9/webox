package service

import (
	"context"
	"encoding/json"
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
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
)

const (
	// Content types
	contentTypeCSV         = "text/csv"
	contentTypeOctetStream = "application/octet-stream"

	// Response limits
	maxErrorMessageLength = 200

	// Password masking constants
	maskedPassword = "******"
)

// Sensitive field names that should be masked
var sensitiveFieldNames = []string{
	// Password fields
	"password", "confirm_password", "old_password", "new_password",
	// Token and key fields
	"token", "access_token", "refresh_token", "api_token", "bearer_token",
	"api_key", "apikey", "secret", "secret_key", "api_secret",
	// Authorization fields
	"authorization", "auth", "credential", "credentials",
	// Certificate fields
	"private_key", "public_key", "certificate", "cert",
}

// auditLogService audit log service implementation
type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
	moduleRepo   repository.ModuleRepository
	userService  service.UserService
	db           *gorm.DB
	logger       logger.Logger
	config       *config.Config
	workerPool   *AuditLogWorkerPool
}

// NewAuditLogService creates audit log service instance
func NewAuditLogService(
	auditLogRepo repository.AuditLogRepository,
	moduleRepo repository.ModuleRepository,
	userService service.UserService,
	db *gorm.DB,
	logger logger.Logger,
	config *config.Config,
) service.AuditLogService {
	svc := &auditLogService{
		auditLogRepo: auditLogRepo,
		moduleRepo:   moduleRepo,
		userService:  userService,
		db:           db,
		logger:       logger,
		config:       config,
	}

	// Initialize worker pool
	svc.workerPool = NewAuditLogWorkerPool(svc, logger)

	return svc
}

// RecordLog records audit log (internal use)
func (s *auditLogService) RecordLog(ctx context.Context, req *request.CreateAuditLogRequest) error {
	if req == nil {
		s.logger.ErrorContext(ctx, "Audit log request is required")
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	auditLog := &model.AuditLog{
		UserID:         req.UserID,
		Username:       req.Username,
		Action:         req.Action,
		Module:         req.Module,
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
	startTime, endTime, err := timeRange.GetParsedTimeRange()
	if err != nil {
		s.logger.ErrorContext(ctx, "Invalid time range parameters", logger.ErrorField(err))
		return nil, "", errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	// Handle user ID based on permissions and request
	userID := req.UserID

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
		contentType = contentTypeOctetStream
	case constants.FormatExcel:
		exportFormat = utils.FormatExcel
		contentType = contentTypeOctetStream
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

	// Record audit log using existing method with background context
	return s.RecordLog(backgroundCtx, auditReq)
}

// SubmitAuditJob submits audit job to worker pool (non-blocking)
func (s *auditLogService) SubmitAuditJob(ginCtx *gin.Context, responseBody []byte, responseTime int) bool {
	if s.workerPool == nil {
		s.logger.Error("Worker pool not initialized")
		return false
	}

	job := &AuditJob{
		ginCtx:       ginCtx,
		responseBody: responseBody,
		responseTime: responseTime,
	}

	return s.workerPool.Submit(job)
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

	// Determine module
	module := s.getModuleName(ctx, ctx.Request.RequestURI)

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
func (s *auditLogService) getActionFromMethod(method, _ string) string {
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

// getModuleName extracts module code from URL path and returns translated name
func (s *auditLogService) getModuleName(ctx context.Context, path string) string {
	// Remove query parameters first
	pathWithoutQuery := s.extractPathWithoutQuery(path)
	pathSegments := strings.Split(strings.Trim(pathWithoutQuery, "/"), "/")

	// Expected format: /api/v1/{module}/...
	if len(pathSegments) < constants.MinPathSegments {
		return i18n.T("common.unknown", constants.DefaultLanguage)
	}

	moduleCode := pathSegments[2]

	// Query module name from database
	if moduleRecord, err := s.moduleRepo.GetByCode(ctx, moduleCode); err == nil && moduleRecord != nil {
		return i18n.T(moduleRecord.Name, constants.DefaultLanguage)
	} else if err != nil {
		s.logger.Warn("Failed to query module by code",
			logger.Field{Key: "code", Value: moduleCode},
			logger.Field{Key: "error", Value: err})
	}

	// Fallback to unknown
	return i18n.T("common.unknown", constants.DefaultLanguage)
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

// addFormData adds form data to params with sensitive data masking
func (s *auditLogService) addFormData(params map[string]interface{}, ctx *gin.Context) {
	for key, values := range ctx.Request.PostForm {
		if s.isSensitiveField(key) {
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
		// Recursively mask sensitive fields in nested JSON
		s.maskSensitiveData(jsonParams, params)
	}
}

// maskSensitiveData recursively masks sensitive fields in nested structures
func (s *auditLogService) maskSensitiveData(source, dest map[string]interface{}) {
	for key, value := range source {
		if s.isSensitiveField(key) {
			dest[key] = maskedPassword
			continue
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively handle nested objects
			maskedNested := make(map[string]interface{})
			s.maskSensitiveData(v, maskedNested)
			dest[key] = maskedNested
		case []interface{}:
			// Handle arrays
			dest[key] = s.maskSensitiveArray(v)
		default:
			dest[key] = value
		}
	}
}

// maskSensitiveArray masks sensitive data in arrays
func (s *auditLogService) maskSensitiveArray(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, item := range arr {
		if itemMap, ok := item.(map[string]interface{}); ok {
			maskedItem := make(map[string]interface{})
			s.maskSensitiveData(itemMap, maskedItem)
			result[i] = maskedItem
		} else {
			result[i] = item
		}
	}
	return result
}

// isSensitiveField checks if a field name is sensitive and should be masked
func (s *auditLogService) isSensitiveField(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, field := range sensitiveFieldNames {
		if lowerKey == field || strings.Contains(lowerKey, field) {
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
