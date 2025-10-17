package service

import (
	"api-service/internal/dto/request"
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

// MockServerAgentRepository implements repository.ServerAgentRepository for testing
type MockServerAgentRepository struct {
	mock.Mock
}

func (m *MockServerAgentRepository) CreateAgent(ctx context.Context, agent *model.ServerAgent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockServerAgentRepository) GetAgentByID(ctx context.Context, id uint) (*model.ServerAgent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServerAgent), args.Error(1)
}

func (m *MockServerAgentRepository) GetAgentByServerID(ctx context.Context, serverID uint) (*model.ServerAgent, error) {
	args := m.Called(ctx, serverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServerAgent), args.Error(1)
}

func (m *MockServerAgentRepository) UpdateAgent(ctx context.Context, agent *model.ServerAgent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockServerAgentRepository) DeleteAgent(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockServerAgentRepository) UpdateAgentLastSeen(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupServerAgentService creates a server agent service with mock dependencies
func setupTestServerAgentService() (*serverAgentService, *MockServerAgentRepository, *MockServerRepository) {
	mockAgentRepo := new(MockServerAgentRepository)
	mockServerRepo := new(MockServerRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, io.Discard)

	service := NewServerAgentService(ServerAgentServiceConfig{
		AgentRepo:  mockAgentRepo,
		ServerRepo: mockServerRepo,
		Logger:     mockLogger,
	}).(*serverAgentService)

	return service, mockAgentRepo, mockServerRepo
}

// TestDeployAgent tests the agent deployment functionality
func TestDeployAgent(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.DeployAgentRequest
		setupMocks    func(*MockServerAgentRepository, *MockServerRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name: "successful docker agent deployment",
			request: &request.DeployAgentRequest{
				ServerID:       1,
				DeploymentType: "docker",
			},
			setupMocks: func(agentRepo *MockServerAgentRepository, serverRepo *MockServerRepository) {
				server := &model.Server{
					ID:      1,
					Name:    "test-server",
					Host:    "192.168.1.100",
					SSHPort: 22,
					OwnerID: 1,
				}
				serverRepo.On("GetServerByID", mock.Anything, uint(1)).Return(server, nil)
				agentRepo.On("GetAgentByServerID", mock.Anything, uint(1)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
				agentRepo.On("CreateAgent", mock.Anything, mock.AnythingOfType("*model.ServerAgent")).
					Run(func(args mock.Arguments) {
						agent := args.Get(1).(*model.ServerAgent)
						agent.ID = 1 // Set the ID as it would be set by the database
					}).
					Return(nil)
				// Mock UpdateAgentLastSeen which is called asynchronously in executeAgentDeployment
				agentRepo.On("UpdateAgentLastSeen", mock.Anything, uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "successful systemd agent deployment",
			request: &request.DeployAgentRequest{
				ServerID:       2,
				DeploymentType: "systemd",
			},
			setupMocks: func(agentRepo *MockServerAgentRepository, serverRepo *MockServerRepository) {
				server := &model.Server{
					ID:      2,
					Name:    "test-server-2",
					Host:    "192.168.1.101",
					SSHPort: 22,
					OwnerID: 1,
				}
				serverRepo.On("GetServerByID", mock.Anything, uint(2)).Return(server, nil)
				agentRepo.On("GetAgentByServerID", mock.Anything, uint(2)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
				agentRepo.On("CreateAgent", mock.Anything, mock.AnythingOfType("*model.ServerAgent")).
					Run(func(args mock.Arguments) {
						agent := args.Get(1).(*model.ServerAgent)
						agent.ID = 2 // Set the ID as it would be set by the database
					}).
					Return(nil)
				// Mock UpdateAgentLastSeen which is called asynchronously in executeAgentDeployment
				agentRepo.On("UpdateAgentLastSeen", mock.Anything, uint(2)).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "agent already exists",
			request: &request.DeployAgentRequest{
				ServerID:       1,
				DeploymentType: "docker",
			},
			setupMocks: func(agentRepo *MockServerAgentRepository, serverRepo *MockServerRepository) {
				server := &model.Server{
					ID:      1,
					Name:    "test-server",
					Host:    "192.168.1.100",
					SSHPort: 22,
					OwnerID: 1,
				}
				existingAgent := &model.ServerAgent{
					ID:             1,
					ServerID:       1,
					DeploymentType: model.AgentDeploymentDocker,
				}
				serverRepo.On("GetServerByID", mock.Anything, uint(1)).Return(server, nil)
				agentRepo.On("GetAgentByServerID", mock.Anything, uint(1)).Return(existingAgent, nil)
			},
			expectedError: true,
			errorCode:     errors.CodeResourceAlreadyExists,
		},
		{
			name: "server not found",
			request: &request.DeployAgentRequest{
				ServerID:       999,
				DeploymentType: "docker",
			},
			setupMocks: func(agentRepo *MockServerAgentRepository, serverRepo *MockServerRepository) {
				serverRepo.On("GetServerByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
		{
			name: "invalid deployment type",
			request: &request.DeployAgentRequest{
				ServerID:       1,
				DeploymentType: "invalid",
			},
			setupMocks: func(agentRepo *MockServerAgentRepository, serverRepo *MockServerRepository) {
				server := &model.Server{
					ID:      1,
					Name:    "test-server",
					Host:    "192.168.1.100",
					SSHPort: 22,
					OwnerID: 1,
				}
				serverRepo.On("GetServerByID", mock.Anything, uint(1)).Return(server, nil)
				agentRepo.On("GetAgentByServerID", mock.Anything, uint(1)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeInvalidParameterFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockAgentRepo, mockServerRepo := setupTestServerAgentService()
			tt.setupMocks(mockAgentRepo, mockServerRepo)

			result, err := service.DeployAgent(context.Background(), tt.request)

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
				assert.Equal(t, tt.request.ServerID, result.ServerID)
				assert.Equal(t, tt.request.DeploymentType, result.DeploymentType)
				// Wait for async deployment to complete (deployment simulation time + buffer)
				time.Sleep(3 * time.Second)
			}

			mockAgentRepo.AssertExpectations(t)
			mockServerRepo.AssertExpectations(t)
		})
	}
}

// TestGetAgent tests retrieving an agent by ID
func TestGetAgent(t *testing.T) {
	tests := []struct {
		name          string
		agentID       uint
		setupMocks    func(*MockServerAgentRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name:    "successful agent retrieval",
			agentID: 1,
			setupMocks: func(repo *MockServerAgentRepository) {
				agent := &model.ServerAgent{
					ID:             1,
					ServerID:       1,
					DeploymentType: model.AgentDeploymentDocker,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				repo.On("GetAgentByID", mock.Anything, uint(1)).Return(agent, nil)
			},
			expectedError: false,
		},
		{
			name:    "agent not found",
			agentID: 999,
			setupMocks: func(repo *MockServerAgentRepository) {
				repo.On("GetAgentByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockAgentRepo, _ := setupTestServerAgentService()
			tt.setupMocks(mockAgentRepo)

			result, err := service.GetAgent(context.Background(), tt.agentID)

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
				assert.Equal(t, tt.agentID, result.ID)
			}

			mockAgentRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateAgent tests agent update functionality
func TestUpdateAgent(t *testing.T) {
	service, mockAgentRepo, _ := setupTestServerAgentService()

	existingAgent := &model.ServerAgent{
		ID:             1,
		ServerID:       1,
		DeploymentType: model.AgentDeploymentDocker,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	updateReq := &request.UpdateAgentRequest{}

	mockAgentRepo.On("GetAgentByID", mock.Anything, uint(1)).Return(existingAgent, nil)
	mockAgentRepo.On("UpdateAgent", mock.Anything, mock.AnythingOfType("*model.ServerAgent")).Return(nil)

	result, err := service.UpdateAgent(context.Background(), 1, updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)

	mockAgentRepo.AssertExpectations(t)
}

// TestUninstallAgent tests agent uninstallation
func TestUninstallAgent(t *testing.T) {
	tests := []struct {
		name          string
		agentID       uint
		setupMocks    func(*MockServerAgentRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name:    "successful agent uninstallation",
			agentID: 1,
			setupMocks: func(repo *MockServerAgentRepository) {
				agent := &model.ServerAgent{
					ID:             1,
					ServerID:       1,
					DeploymentType: model.AgentDeploymentDocker,
				}
				repo.On("GetAgentByID", mock.Anything, uint(1)).Return(agent, nil)
				// Mock DeleteAgent which is called asynchronously in executeAgentUninstallation
				repo.On("DeleteAgent", mock.Anything, uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "agent not found",
			agentID: 999,
			setupMocks: func(repo *MockServerAgentRepository) {
				repo.On("GetAgentByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockAgentRepo, _ := setupTestServerAgentService()
			tt.setupMocks(mockAgentRepo)

			err := service.UninstallAgent(context.Background(), tt.agentID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorCode != 0 {
					appErr, ok := err.(*errors.AppError)
					assert.True(t, ok, "Expected AppError")
					if ok {
						assert.Equal(t, tt.errorCode, appErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				// Wait for async goroutine to complete (executeAgentUninstallation)
				time.Sleep(2 * time.Second)
			}

			mockAgentRepo.AssertExpectations(t)
		})
	}
}

// TestRestartAgent tests agent restart functionality
func TestRestartAgent(t *testing.T) {
	tests := []struct {
		name          string
		agentID       uint
		setupMocks    func(*MockServerAgentRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name:    "successful agent restart",
			agentID: 1,
			setupMocks: func(repo *MockServerAgentRepository) {
				agent := &model.ServerAgent{
					ID:             1,
					ServerID:       1,
					DeploymentType: model.AgentDeploymentDocker,
				}
				repo.On("GetAgentByID", mock.Anything, uint(1)).Return(agent, nil)
				// Mock UpdateAgentLastSeen which is called asynchronously in executeAgentRestart
				repo.On("UpdateAgentLastSeen", mock.Anything, uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "agent not found",
			agentID: 999,
			setupMocks: func(repo *MockServerAgentRepository) {
				repo.On("GetAgentByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockAgentRepo, _ := setupTestServerAgentService()
			tt.setupMocks(mockAgentRepo)

			err := service.RestartAgent(context.Background(), tt.agentID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorCode != 0 {
					appErr, ok := err.(*errors.AppError)
					assert.True(t, ok, "Expected AppError")
					if ok {
						assert.Equal(t, tt.errorCode, appErr.Code)
					}
				}
			} else {
				assert.NoError(t, err)
				// Wait for async restart to complete (restart simulation time + buffer)
				time.Sleep(2 * time.Second)
			}

			mockAgentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAgentLogs tests retrieving agent logs
func TestGetAgentLogs(t *testing.T) {
	service, mockAgentRepo, _ := setupTestServerAgentService()

	agent := &model.ServerAgent{
		ID:             1,
		ServerID:       1,
		DeploymentType: model.AgentDeploymentDocker,
	}

	mockAgentRepo.On("GetAgentByID", mock.Anything, uint(1)).Return(agent, nil)

	req := &request.GetAgentLogsRequest{
		AgentID: 1,
		Lines:   10,
	}

	result, err := service.GetAgentLogs(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.AgentID)
	assert.NotEmpty(t, result.Logs)

	mockAgentRepo.AssertExpectations(t)
}

// TestGetAgentStatus tests checking agent status
func TestGetAgentStatus(t *testing.T) {
	tests := []struct {
		name          string
		agentID       uint
		setupMocks    func(*MockServerAgentRepository)
		expectedError bool
		errorCode     errors.ErrorCode
	}{
		{
			name:    "successful status check",
			agentID: 1,
			setupMocks: func(repo *MockServerAgentRepository) {
				now := time.Now()
				agent := &model.ServerAgent{
					ID:              1,
					ServerID:        1,
					DeploymentType:  model.AgentDeploymentDocker,
					LastHeartbeatAt: &now,
				}
				repo.On("GetAgentByID", mock.Anything, uint(1)).Return(agent, nil)
			},
			expectedError: false,
		},
		{
			name:    "agent not found",
			agentID: 999,
			setupMocks: func(repo *MockServerAgentRepository) {
				repo.On("GetAgentByID", mock.Anything, uint(999)).Return(nil, errors.NewAppError(errors.CodeRecordNotFound))
			},
			expectedError: true,
			errorCode:     errors.CodeRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockAgentRepo, _ := setupTestServerAgentService()
			tt.setupMocks(mockAgentRepo)

			result, err := service.GetAgentStatus(context.Background(), tt.agentID)

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
				assert.Equal(t, tt.agentID, result.AgentID)
				assert.NotEmpty(t, result.Status)
			}

			mockAgentRepo.AssertExpectations(t)
		})
	}
}
