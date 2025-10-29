package service

import (
	"api-service/internal/config"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecretService implements service.SecretService for testing
type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) CreateTextSecret(ctx context.Context, req *request.CreateTextSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	args := m.Called(ctx, req, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SecretResponse), args.Error(1)
}

func (m *MockSecretService) CreateAccountSecret(ctx context.Context, req *request.CreateAccountSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	args := m.Called(ctx, req, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SecretResponse), args.Error(1)
}

func (m *MockSecretService) CreateFileSecret(ctx context.Context, req *request.CreateFileSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	args := m.Called(ctx, req, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SecretResponse), args.Error(1)
}

func (m *MockSecretService) GetSecret(ctx context.Context, id uint, userID uint) (*response.SecretDetailResponse, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SecretDetailResponse), args.Error(1)
}

func (m *MockSecretService) ListSecrets(ctx context.Context, req *request.ListSecretsRequest, userID uint) (*common.PaginationResponse, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.PaginationResponse), args.Error(1)
}

func (m *MockSecretService) UpdateSecret(ctx context.Context, id uint, req *request.UpdateSecretRequest, userID uint) (*response.SecretResponse, error) {
	args := m.Called(ctx, id, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SecretResponse), args.Error(1)
}

func (m *MockSecretService) DeleteSecret(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockSecretService) GetSecretReferencesByResourceCode(ctx context.Context, resourceCode string) ([]*response.SecretReferenceResponse, error) {
	args := m.Called(ctx, resourceCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*response.SecretReferenceResponse), args.Error(1)
}

func (m *MockSecretService) DeleteSecretReferencesByResourceCode(ctx context.Context, resourceCode string) error {
	args := m.Called(ctx, resourceCode)
	return args.Error(0)
}

func (m *MockSecretService) DownloadSecretFile(ctx context.Context, id uint, userID uint) ([]byte, string, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]byte), args.Get(1).(string), args.Error(2)
}

func (m *MockSecretService) CreateReference(ctx context.Context, req *request.CreateReferenceRequest) (*response.ReferenceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ReferenceResponse), args.Error(1)
}

// MockServerRepository implements repository.ServerRepository for testing
type MockServerRepository struct {
	mock.Mock
}

func (m *MockServerRepository) CreateServer(ctx context.Context, server *model.Server) error {
	args := m.Called(ctx, server)
	return args.Error(0)
}

func (m *MockServerRepository) GetServerByID(ctx context.Context, id uint) (*model.Server, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Server), args.Error(1)
}

func (m *MockServerRepository) GetServerByName(ctx context.Context, name string) (*model.Server, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Server), args.Error(1)
}

func (m *MockServerRepository) UpdateServer(ctx context.Context, server *model.Server) error {
	args := m.Called(ctx, server)
	return args.Error(0)
}

func (m *MockServerRepository) DeleteServer(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockServerRepository) ListServers(ctx context.Context, req *request.ListServersRequest) ([]*model.Server, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Server), args.Get(1).(int64), args.Error(2)
}

func (m *MockServerRepository) ExistsServerByName(ctx context.Context, name string, excludeID ...uint) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockServerRepository) ExistsServerByHost(ctx context.Context, host string) (bool, error) {
	args := m.Called(ctx, host)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockServerRepository) BatchUpdateServerStatus(ctx context.Context, ids []uint, status string) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockServerRepository) CountServersByStatus(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockServerRepository) GetServerAgentsByServerID(ctx context.Context, serverID uint) ([]*model.ServerAgent, error) {
	args := m.Called(ctx, serverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ServerAgent), args.Error(1)
}

func (m *MockServerRepository) GetServersByIDs(ctx context.Context, ids []uint) ([]*model.Server, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Server), args.Error(1)
}

func (m *MockServerRepository) GetServersByStatus(ctx context.Context, status string) ([]*model.Server, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Server), args.Error(1)
}

func (m *MockServerRepository) UpdateServerStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockServerRepository) UpdateServerLastSeen(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupServerService creates a server service with mock dependencies
func setupTestServerService() (*serverService, *MockServerRepository, *MockSystemConfigRepository) {
	mockRepo := new(MockServerRepository)
	mockSystemConfigRepo := new(MockSystemConfigRepository)
	mockSecretService := new(MockSecretService)
	mockSecretRefRepo := new(MockSecretReferenceRepository)
	mockSecretRepo := new(MockSecretRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, io.Discard)

	// Setup default mock for GetSecretReferencesByResourceCode to avoid unexpected call errors
	mockSecretService.On("GetSecretReferencesByResourceCode", mock.Anything, mock.Anything).
		Return(nil, errors.NewAppError(errors.CodeRecordNotFound)).Maybe()

	// Setup default mock for ListByResourceCode to avoid unexpected call errors
	mockSecretRefRepo.On("ListByResourceCode", mock.Anything, mock.Anything).
		Return([]*model.SecretReference{}, nil).Maybe()

	service := NewServerService(&ServerServiceConfig{
		Logger:           mockLogger,
		ServerRepo:       mockRepo,
		SecretService:    mockSecretService,
		SecretRefRepo:    mockSecretRefRepo,
		SecretRepo:       mockSecretRepo,
		SystemConfigRepo: mockSystemConfigRepo,
		Config:           &config.Config{},
	}).(*serverService)

	return service, mockRepo, mockSystemConfigRepo
}

// TestCreateServer tests the server creation functionality
func TestCreateServer(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.CreateServerRequest
		setupMocks    func(*MockServerRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name: "successful server creation",
			request: &request.CreateServerRequest{
				Name:    "test-server",
				Host:    "192.168.1.100",
				SSHPort: 22,
				// Note: Hostname removed - dynamically collected by Agent
				// Note: OwnerID removed - auto-populated from JWT token
			},
			setupMocks: func(repo *MockServerRepository) {
				repo.On("ExistsServerByName", mock.Anything, "test-server", mock.Anything).Return(false, nil)
				repo.On("CreateServer", mock.Anything, mock.AnythingOfType("*model.Server")).Return(nil).Run(func(args mock.Arguments) {
					// Set the ID for the created server
					server := args.Get(1).(*model.Server)
					server.ID = 1
				})
				// Mock for background connectivity test (these calls happen asynchronously)
				repo.On("GetServerByID", mock.Anything, uint(1)).Return(&model.Server{ID: 1, Name: "test-server"}, nil).Maybe()
				repo.On("UpdateServerStatus", mock.Anything, uint(1), mock.Anything).Return(nil).Maybe()
			},
			expectedError: false,
		},
		{
			name: "server name already exists",
			request: &request.CreateServerRequest{
				Name:    "existing-server",
				Host:    "192.168.1.100",
				SSHPort: 22,
				// Note: Hostname removed - dynamically collected by Agent
				// Note: OwnerID removed - auto-populated from JWT token
			},
			setupMocks: func(repo *MockServerRepository) {
				repo.On("ExistsServerByName", mock.Anything, "existing-server", mock.Anything).Return(true, nil)
			},
			expectedError: true,
			errorCode:     errors.CodeResourceAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _ := setupTestServerService()
			tt.setupMocks(mockRepo)

			// CreateServer now requires currentUserID parameter
			result, err := service.CreateServer(context.Background(), tt.request, uint(1))

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorCode != 0 {
					appErr, ok := err.(*errors.AppError)
					assert.True(t, ok, "Expected AppError")
					if ok {
						assert.Equal(t, tt.errorCode, appErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetServer tests retrieving a server by ID
func TestGetServer(t *testing.T) {
	tests := []struct {
		name          string
		serverID      uint
		setupMocks    func(*MockServerRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name:     "successful server retrieval",
			serverID: 1,
			setupMocks: func(repo *MockServerRepository) {
				server := &model.Server{
					ID:       1,
					Name:     "test-server",
					Hostname: "test.example.com",
					Host:     "192.168.1.100",
					SSHPort:  22,
				}
				repo.On("GetServerByID", mock.Anything, uint(1)).Return(server, nil)
			},
			expectedError: false,
		},
		{
			name:     "server not found",
			serverID: 999,
			setupMocks: func(repo *MockServerRepository) {
				repo.On("GetServerByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _ := setupTestServerService()
			tt.setupMocks(mockRepo)

			result, err := service.GetServer(context.Background(), tt.serverID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.serverID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetServerStatus tests checking the status of a single server
func TestGetServerStatus(t *testing.T) {
	tests := []struct {
		name          string
		serverID      uint
		setupMocks    func(*MockServerRepository, *MockSystemConfigRepository)
		expectedError bool
	}{
		{
			name:     "successful status check",
			serverID: 1,
			setupMocks: func(repo *MockServerRepository, sysConfigRepo *MockSystemConfigRepository) {
				server := &model.Server{
					ID:      1,
					Name:    "test-server",
					Host:    "127.0.0.1",
					SSHPort: 22,
				}
				repo.On("GetServerByID", mock.Anything, uint(1)).Return(server, nil)
				// Mock system config for SSH timeout
				sysConfigRepo.On("GetByKey", mock.Anything, "server.ssh_timeout").Return(nil, errors.NewAppError(errors.CodeRecordNotFound)).Maybe()
			},
			expectedError: false,
		},
		{
			name:     "server not found",
			serverID: 999,
			setupMocks: func(repo *MockServerRepository, sysConfigRepo *MockSystemConfigRepository) {
				repo.On("GetServerByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
				sysConfigRepo.On("GetByKey", mock.Anything, mock.Anything).Return(nil, errors.NewAppError(errors.CodeRecordNotFound)).Maybe()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockSysConfigRepo := setupTestServerService()
			tt.setupMocks(mockRepo, mockSysConfigRepo)

			result, err := service.GetServerStatus(context.Background(), tt.serverID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.serverID, result.ServerID)
				assert.NotNil(t, result.Checks)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateServer tests server update functionality
func TestUpdateServer(t *testing.T) {
	service, mockRepo, _ := setupTestServerService()

	existingServer := &model.Server{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "old-server",
		Hostname:  "old.example.com",
		Host:      "192.168.1.100",
		SSHPort:   22,
		OwnerID:   1,
	}

	updateReq := &request.UpdateServerRequest{
		Name: stringPtr("updated-server"),
		// Note: Hostname removed - dynamically collected by Agent, not updated via API
	}

	mockRepo.On("GetServerByID", mock.Anything, uint(1)).Return(existingServer, nil)
	mockRepo.On("ExistsServerByName", mock.Anything, "updated-server", []uint{uint(1)}).Return(false, nil)
	mockRepo.On("UpdateServer", mock.Anything, mock.AnythingOfType("*model.Server")).Return(nil)

	result, err := service.UpdateServer(context.Background(), 1, updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-server", result.Name)

	mockRepo.AssertExpectations(t)
}

// TestDeleteServer tests server deletion
func TestDeleteServer(t *testing.T) {
	service, mockRepo, _ := setupTestServerService()
	mockSecretService := new(MockSecretService)
	mockSecretRefRepo := new(MockSecretReferenceRepository)
	service.secretService = mockSecretService
	service.secretRefRepo = mockSecretRefRepo

	existingServer := &model.Server{
		ID:   1,
		Code: "srv123",
		Name: "test-server",
	}

	mockRepo.On("GetServerByID", mock.Anything, uint(1)).Return(existingServer, nil)
	mockSecretRefRepo.On("ListByResourceCode", mock.Anything, "srv123").Return([]*model.SecretReference{}, nil)
	mockRepo.On("DeleteServer", mock.Anything, uint(1)).Return(nil)

	err := service.DeleteServer(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestListServers tests server listing functionality
func TestListServers(t *testing.T) {
	service, mockRepo, _ := setupTestServerService()

	servers := []*model.Server{
		{
			ID:   1,
			Name: "server-1",
		},
		{
			ID:   2,
			Name: "server-2",
		},
	}

	mockRepo.On("ListServers", mock.Anything, mock.AnythingOfType("*request.ListServersRequest")).Return(servers, int64(2), nil)

	req := &request.ListServersRequest{
		PaginationRequest: common.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	result, err := service.ListServers(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Items))
	assert.Equal(t, int64(2), result.Total)

	mockRepo.AssertExpectations(t)
}
