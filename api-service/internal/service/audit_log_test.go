package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

// MockAuditLogRepository mock implementation of AuditLogRepository
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *model.AuditLog) error {
	args := m.Called(ctx, auditLog)
	return args.Error(0)
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) List(ctx context.Context, filter *request.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*model.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogRepository) GetStatistics(ctx context.Context, filter *request.StatisticsFilter) (*response.AuditLogStatistics, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*response.AuditLogStatistics), args.Error(1)
}

func (m *MockAuditLogRepository) Export(ctx context.Context, filter *request.AuditLogFilter) ([]*model.AuditLog, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*model.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) CleanupOldLogs(ctx context.Context, beforeDate time.Time) (int64, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).(int64), args.Error(1)
}

// MockLogger mock implementation of Logger
type MockLogger struct {
	mock.Mock
}

// MockUserService mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

// Only implement the methods needed by audit log service
func (m *MockUserService) GetUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

// Placeholder implementations for other methods (not used by audit service)
func (m *MockUserService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.UserLoginResponse), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, req *request.UserListRequest) (*response.UserListResponse, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.UserListResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) CreateUser(ctx context.Context, req *request.UserCreateRequest) (*response.UserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) UpdateUserPassword(ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) ValidateUserAccess(ctx context.Context, userID uint, resource string) error {
	args := m.Called(ctx, userID, resource)
	return args.Error(0)
}

func (m *MockUserService) CheckUserQuota(ctx context.Context, userID uint, resourceType string) error {
	args := m.Called(ctx, userID, resourceType)
	return args.Error(0)
}

func (m *MockLogger) Debug(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) InfoContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) WarnContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) ErrorContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) IsDebugEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) IsInfoEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) IsWarnEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) IsErrorEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) WithFields(fields ...logger.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) SetLevel(level logger.Level) {
	m.Called(level)
}

func (m *MockLogger) SetOutput(w io.Writer) {
	m.Called(w)
}

// Test setup helpers
// setupAuditLogService creates service with mocked dependencies
func setupAuditLogService() (service.AuditLogService, *MockAuditLogRepository, *MockUserService, *MockLogger) {
	mockRepo := &MockAuditLogRepository{}
	mockUserService := &MockUserService{}
	mockLogger := &MockLogger{}

	// Create test database for export functionality
	mockDB := setupAuditTestDB()

	// Initialize i18n for testing
	_ = i18n.Init() // Initialize with default config
	mockI18n := i18n.NewI18n()

	// Create test configuration
	testConfig := &config.Config{
		AuditLog: config.AuditLogConfig{
			SkipPaths: []string{"/health", "/api/v1/health"},
			SensitiveGetPaths: []string{
				"/api/v1/users/profile",
				"/api/v1/api-tokens",
				"/api/v1/roles",
				"/api/v1/permissions",
				"/api/v1/users/:id/two-factor",
			},
			AuditMethods: []string{"POST", "PUT", "DELETE", "PATCH"},
			SkipMethods:  []string{"OPTIONS", "HEAD"},
		},
	}

	// Setup logger mock expectations
	mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("DebugContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("InfoContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("WarnContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("ErrorContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()

	auditLogService := NewAuditLogService(mockRepo, mockUserService, mockDB, mockLogger, mockI18n, testConfig)
	return auditLogService, mockRepo, mockUserService, mockLogger
}

// createTestGinContext creates a mock gin.Context for testing
func createTestGinContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set("language", "zh-CN") // Set default language for testing
	return c
}

func createTestAuditLog() *model.AuditLog {
	now := time.Now()
	userID := uint(1)
	responseStatus := 200
	responseTime := 150

	return &model.AuditLog{
		ID:             1, // Set for testing purposes - in real app this would be auto-generated
		UserID:         &userID,
		Username:       "testuser",
		Action:         constants.ActionCreate,
		Module:         "User Management",
		ResourceType:   "USER",
		ResourceID:     &userID,
		ResourceName:   "Test User",
		Description:    "Create user test",
		IPAddress:      "192.168.1.1",
		UserAgent:      "Mozilla/5.0 Test",
		RequestMethod:  "POST",
		RequestURL:     "/api/v1/users",
		RequestParams:  `{"name": "test"}`,
		ResponseStatus: &responseStatus,
		ResponseTime:   &responseTime,
		Success:        true,
		ErrorMessage:   "",
		CreatedAt:      now,
	}
}

func createTestCreateRequest() *request.CreateAuditLogRequest {
	userID := uint(1)
	responseStatus := 200
	responseTime := 150

	return &request.CreateAuditLogRequest{
		UserID:         &userID,
		Username:       "testuser",
		Action:         constants.ActionCreate,
		Module:         "User Management",
		ResourceType:   "USER",
		ResourceID:     &userID,
		ResourceName:   "Test User",
		Description:    "Create user test",
		IPAddress:      "192.168.1.1",
		UserAgent:      "Mozilla/5.0 Test",
		RequestMethod:  "POST",
		RequestURL:     "/api/v1/users",
		RequestParams:  `{"name": "test"}`,
		ResponseStatus: &responseStatus,
		ResponseTime:   &responseTime,
		Success:        true,
		ErrorMessage:   "",
	}
}

// Tests for RecordLog
func TestAuditLogService_RecordLog_Success(t *testing.T) {
	service, mockRepo, _, mockLogger := setupAuditLogService()
	ctx := context.Background()
	req := createTestCreateRequest()

	// Mock repository call
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	// Mock logger calls - RecordLog calls InfoContext twice
	mockLogger.On("InfoContext", ctx, "Recording audit log", mock.Anything).Return()
	mockLogger.On("InfoContext", ctx, "Audit log recorded successfully", mock.Anything).Return()

	// Execute
	err := service.RecordLog(ctx, req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuditLogService_RecordLog_RepositoryError(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	req := createTestCreateRequest()

	// Mock repository error
	expectedError := errors.New("database connection failed")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(expectedError)

	// Execute
	err := service.RecordLog(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to record audit log")
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_RecordLog_NilRequest(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()

	// Execute
	err := service.RecordLog(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "audit log request is required")
}

// Tests for GetAuditLog
func TestAuditLogService_GetAuditLog_Success(t *testing.T) {
	service, mockRepo, mockUserService, _ := setupAuditLogService()
	ctx := context.Background()
	testLog := createTestAuditLog()

	// Mock repository call
	mockRepo.On("GetByID", ctx, uint(1)).Return(testLog, nil)

	// Mock user service call for user info
	userResp := &response.UserResponse{
		ID:       1,
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
	}
	mockUserService.On("GetUser", ctx, uint(1)).Return(userResp, nil)

	// Execute
	result, err := service.GetAuditLog(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testLog.ID, result.ID)
	assert.Equal(t, testLog.Action, result.Action)
	assert.Equal(t, testLog.Description, result.Description)
	assert.Equal(t, "Test User", result.User.Nickname)
	// Note: AuditLogUserResponse no longer contains Email field
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_GetAuditLog_NotFound(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()

	// Mock repository call returning nil
	mockRepo.On("GetByID", ctx, uint(999)).Return((*model.AuditLog)(nil), errors.New("record not found"))

	// Execute
	result, err := service.GetAuditLog(ctx, 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get audit log")
	mockRepo.AssertExpectations(t)
}

// Tests for ListAuditLogs
func TestAuditLogService_ListAuditLogs_Success(t *testing.T) {
	service, mockRepo, mockUserService, _ := setupAuditLogService()
	ctx := context.Background()
	testLogs := []*model.AuditLog{createTestAuditLog()}
	req := &request.ListAuditLogRequest{
		Page:     1,
		PageSize: 20,
		Action:   "CREATE",
	}

	// Mock repository call
	mockRepo.On("List", ctx, mock.AnythingOfType("*request.AuditLogFilter")).Return(testLogs, int64(1), nil)

	// Mock user service call for user info
	userResp := &response.UserResponse{
		ID:       1,
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
	}
	mockUserService.On("GetUser", ctx, uint(1)).Return(userResp, nil)

	// Execute
	result, err := service.ListAuditLogs(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, "Test User", result.Items[0].User.Nickname)
	// Note: AuditLogUserResponse no longer contains Email field
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_ListAuditLogs_WithDefaults(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	testLogs := []*model.AuditLog{}
	req := &request.ListAuditLogRequest{} // Empty request, should use defaults

	// Mock repository call
	mockRepo.On("List", ctx, mock.AnythingOfType("*request.AuditLogFilter")).Return(testLogs, int64(0), nil)

	// Execute
	result, err := service.ListAuditLogs(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)      // Default page
	assert.Equal(t, 20, result.PageSize) // Default page size
	mockRepo.AssertExpectations(t)
}

// Tests for GetStatistics
func TestAuditLogService_GetStatistics_Success(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	req := &request.AuditLogStatisticsRequest{
		GroupBy: "day",
	}

	mockStats := &response.AuditLogStatistics{
		TotalOperations:   100,
		SuccessOperations: 95,
		FailedOperations:  5,
		SuccessRate:       95.0,
		TopUsers: []response.UserOperationCount{
			{UserID: 1, Username: "testuser", OperationCount: 50},
		},
		TopActions: []response.ActionCount{
			{Action: "CREATE", Count: 30},
		},
		Timeline: []response.TimelineCount{
			{Date: "2025-08-28", Count: 25},
		},
	}

	// Mock repository call
	mockRepo.On("GetStatistics", ctx, mock.AnythingOfType("*request.StatisticsFilter")).Return(mockStats, nil)

	// Execute
	result, err := service.GetStatistics(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.TotalOperations)
	assert.Equal(t, int64(95), result.SuccessOperations)
	assert.Equal(t, int64(5), result.FailedOperations)
	assert.Equal(t, 1, len(result.TopUsers))
	assert.Equal(t, 1, len(result.TopActions))
	mockRepo.AssertExpectations(t)
}

// Tests for ExportAuditLogs
func TestAuditLogService_ExportAuditLogs_CSV_Success(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()
	ginCtx := createTestGinContext()
	req := &request.ExportAuditLogRequest{
		Format: "csv",
	}

	// Execute
	data, contentType, err := service.ExportAuditLogs(ctx, ginCtx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "application/octet-stream", contentType)
	assert.Contains(t, string(data), "action,created_at,description") // CSV header with database fields
	// Mock repository is no longer called since we use DBExporter now
}

func TestAuditLogService_ExportAuditLogs_JSON_Success(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()
	ginCtx := createTestGinContext()
	req := &request.ExportAuditLogRequest{
		Format: "json",
	}

	// Execute
	data, contentType, err := service.ExportAuditLogs(ctx, ginCtx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "application/octet-stream", contentType)
	assert.Contains(t, string(data), `"id": 2`) // JSON content with actual test data (note the space)
}

func TestAuditLogService_ExportAuditLogs_UnsupportedFormat(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()
	ginCtx := createTestGinContext()
	req := &request.ExportAuditLogRequest{
		Format: "pdf", // Unsupported format, should default to CSV
	}

	// Execute
	data, contentType, err := service.ExportAuditLogs(ctx, ginCtx, req)

	// Assert - unsupported format defaults to CSV
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "application/octet-stream", contentType)
	assert.Contains(t, string(data), "action,created_at,description") // CSV header with database fields
}

func TestAuditLogService_ExportAuditLogs_Excel_Success(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()
	ginCtx := createTestGinContext()
	req := &request.ExportAuditLogRequest{
		Format: "excel",
	}

	// Execute
	data, contentType, err := service.ExportAuditLogs(ctx, ginCtx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Equal(t, "application/octet-stream", contentType)

	// Verify that data is not empty and contains Excel magic bytes
	assert.True(t, len(data) > 100) // Excel files are typically larger than 100 bytes

	// Check for Excel file signature (first few bytes should indicate it's a zip-based format)
	assert.Equal(t, "PK", string(data[0:2])) // Excel files start with "PK" (ZIP signature)
}

// TestAuditLogService_ExportAuditLogs_DefaultValues tests export with default user ID and start time
func TestAuditLogService_ExportAuditLogs_DefaultValues(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	ctx := context.Background()
	ginCtx := createTestGinContext()

	// Set user ID directly in context to simulate authenticated user
	testUserID := uint(2) // Use existing user ID from test data
	ginCtx.Set("user_id", testUserID)
	ginCtx.Set("username", "testuser")

	req := &request.ExportAuditLogRequest{
		Format: "csv",
		// UserID and StartTime intentionally left nil to test defaults
	}

	// Execute
	data, contentType, err := service.ExportAuditLogs(ctx, ginCtx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "application/octet-stream", contentType)
	// Should return CSV header even if no data matches the user filter
	assert.Contains(t, string(data), "action,created_at,description")
}

// Tests for CleanupExpiredLogs
func TestAuditLogService_CleanupExpiredLogs_Success(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	retentionDays := 90

	// Mock repository call
	mockRepo.On("CleanupOldLogs", ctx, mock.AnythingOfType("time.Time")).Return(int64(25), nil)

	// Execute
	deletedCount, err := service.CleanupExpiredLogs(ctx, retentionDays)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(25), deletedCount)
	mockRepo.AssertExpectations(t)
}

func TestAuditLogService_CleanupExpiredLogs_InvalidRetentionDays(t *testing.T) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()

	// Mock repository call (service will use default 90 days for negative input)
	mockRepo.On("CleanupOldLogs", ctx, mock.AnythingOfType("time.Time")).Return(int64(0), nil)

	// Test with negative retention days (should be converted to default 90 days)
	deletedCount, err := service.CleanupExpiredLogs(ctx, -1)

	// Assert (no error expected as service handles negative values by using default)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), deletedCount)
	mockRepo.AssertExpectations(t)
}

// Integration-style tests
func TestAuditLogService_FullWorkflow(t *testing.T) {
	service, mockRepo, mockUserService, _ := setupAuditLogService()
	ctx := context.Background()

	// Mock user service for all user info calls
	userResp := &response.UserResponse{
		ID:       1,
		Username: "testuser",
		Nickname: "Test User",
		Email:    "test@example.com",
	}
	mockUserService.On("GetUser", ctx, uint(1)).Return(userResp, nil)

	// Step 1: Record a log
	createReq := createTestCreateRequest()
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	err := service.RecordLog(ctx, createReq)
	assert.NoError(t, err)

	// Step 2: Retrieve the log
	testLog := createTestAuditLog()
	mockRepo.On("GetByID", ctx, uint(1)).Return(testLog, nil)

	result, err := service.GetAuditLog(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, testLog.ID, result.ID)
	assert.Equal(t, "Test User", result.User.Nickname)

	// Step 3: List logs
	mockRepo.On("List", ctx, mock.AnythingOfType("*request.AuditLogFilter")).Return([]*model.AuditLog{testLog}, int64(1), nil)

	listReq := &request.ListAuditLogRequest{Page: 1, PageSize: 20}
	listResult, err := service.ListAuditLogs(ctx, listReq)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), listResult.Total)

	mockRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkAuditLogService_RecordLog(b *testing.B) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	req := createTestCreateRequest()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.AuditLog")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.RecordLog(ctx, req)
	}
}

func BenchmarkAuditLogService_GetAuditLog(b *testing.B) {
	service, mockRepo, _, _ := setupAuditLogService()
	ctx := context.Background()
	testLog := createTestAuditLog()

	mockRepo.On("GetByID", ctx, uint(1)).Return(testLog, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetAuditLog(ctx, 1)
	}
}

// TestAuditLogService_ConfigurablePaths tests that the audit service uses configurable paths
func TestAuditLogService_ConfigurablePaths(t *testing.T) {
	service, _, _, _ := setupAuditLogService()
	auditService := service.(*auditLogService)

	// Test skip paths from configuration
	assert.True(t, auditService.ShouldSkipAudit("GET", "/health"))
	assert.True(t, auditService.ShouldSkipAudit("GET", "/api/v1/health"))
	assert.False(t, auditService.ShouldSkipAudit("POST", "/api/v1/users"))

	// Test sensitive GET paths from configuration
	assert.False(t, auditService.ShouldSkipAudit("GET", "/api/v1/users/profile"))
	assert.False(t, auditService.ShouldSkipAudit("GET", "/api/v1/api-tokens"))
	assert.False(t, auditService.ShouldSkipAudit("GET", "/api/v1/roles"))
	assert.False(t, auditService.ShouldSkipAudit("GET", "/api/v1/permissions"))

	// Test audit methods from configuration
	assert.False(t, auditService.ShouldSkipAudit("POST", "/api/v1/users"))
	assert.False(t, auditService.ShouldSkipAudit("PUT", "/api/v1/users/1"))
	assert.False(t, auditService.ShouldSkipAudit("DELETE", "/api/v1/users/1"))
	assert.False(t, auditService.ShouldSkipAudit("PATCH", "/api/v1/users/1"))

	// Test skip methods from configuration
	assert.True(t, auditService.ShouldSkipAudit("OPTIONS", "/api/v1/users"))
	assert.True(t, auditService.ShouldSkipAudit("HEAD", "/api/v1/users"))

	// Test regular GET requests should be skipped
	assert.True(t, auditService.ShouldSkipAudit("GET", "/api/v1/users"))
	assert.True(t, auditService.ShouldSkipAudit("GET", "/api/v1/some/random/path"))
}

// setupAuditTestDB creates an in-memory SQLite database for testing
func setupAuditTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// Auto-migrate the AuditLog model for testing
	err = db.AutoMigrate(&model.AuditLog{})
	if err != nil {
		panic("failed to auto-migrate test database: " + err.Error())
	}

	// Create some test data for export functionality
	userID1 := uint(1)
	userID2 := uint(2)
	statusCode1 := 201
	statusCode2 := 200
	responseTime1 := 100
	responseTime2 := 150

	testAuditLogs := []*model.AuditLog{
		{
			ID:             1,
			UserID:         &userID1,
			Username:       "testuser1",
			Action:         "CREATE",
			Module:         "USER",
			ResourceType:   "user",
			RequestMethod:  "POST",
			RequestURL:     "/api/v1/users",
			ResponseStatus: &statusCode1,
			ResponseTime:   &responseTime1,
			IPAddress:      "192.168.1.1",
			UserAgent:      "TestAgent/1.0",
			Description:    "Created new user",
			Success:        true,
			CreatedAt:      time.Now().Add(-24 * time.Hour),
		},
		{
			ID:             2,
			UserID:         &userID2,
			Username:       "testuser2",
			Action:         "UPDATE",
			Module:         "USER",
			ResourceType:   "user",
			RequestMethod:  "PUT",
			RequestURL:     "/api/v1/users/2",
			ResponseStatus: &statusCode2,
			ResponseTime:   &responseTime2,
			IPAddress:      "192.168.1.2",
			UserAgent:      "TestAgent/1.0",
			Description:    "Updated user profile",
			Success:        true,
			CreatedAt:      time.Now().Add(-12 * time.Hour),
		},
	}

	for _, log := range testAuditLogs {
		db.Create(log)
	}

	return db
}
