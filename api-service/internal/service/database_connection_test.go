package service

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/model"
	pkgErrors "api-service/pkg/errors"
	"api-service/pkg/logger"
)

// MockDatabaseConnectionRepository mocks DatabaseConnectionRepository interface
type MockDatabaseConnectionRepository struct {
	mock.Mock
}

func (m *MockDatabaseConnectionRepository) Create(ctx context.Context, conn *model.DatabaseConnection) error {
	args := m.Called(ctx, conn)
	return args.Error(0)
}

func (m *MockDatabaseConnectionRepository) GetByID(ctx context.Context, id uint) (*model.DatabaseConnection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DatabaseConnection), args.Error(1)
}

func (m *MockDatabaseConnectionRepository) GetList(ctx context.Context, req *request.GetDatabaseConnectionListRequest, ownerID uint) ([]*model.DatabaseConnection, int64, error) {
	args := m.Called(ctx, req, ownerID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.DatabaseConnection), args.Get(1).(int64), args.Error(2)
}

func (m *MockDatabaseConnectionRepository) Update(ctx context.Context, conn *model.DatabaseConnection) error {
	args := m.Called(ctx, conn)
	return args.Error(0)
}

func (m *MockDatabaseConnectionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupDatabaseConnectionTest sets up test dependencies
func setupDatabaseConnectionTest(t *testing.T) (*databaseConnectionService, *MockDatabaseConnectionRepository) {
	mockRepo := new(MockDatabaseConnectionRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, io.Discard)

	service := &databaseConnectionService{
		repo:   mockRepo,
		logger: mockLogger,
	}

	return service, mockRepo
}

func TestCreateConnection(t *testing.T) {
	service, mockRepo := setupDatabaseConnectionTest(t)
	ctx := context.Background()
	ownerID := uint(1)

	testCases := []struct {
		name          string
		req           *request.CreateDatabaseConnectionRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - MySQL connection",
			req: &request.CreateDatabaseConnectionRequest{
				Name:   "Test MySQL",
				DBType: "mysql",
				Host:   "localhost",
				Port:   3306,
			},
			setupMock: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.DatabaseConnection")).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Success - PostgreSQL with config",
			req: &request.CreateDatabaseConnectionRequest{
				Name:   "Test PostgreSQL",
				DBType: "postgresql",
				Host:   "localhost",
				Port:   5432,
				Config: map[string]interface{}{
					"charset": "utf8mb4",
					"timeout": 30,
				},
			},
			setupMock: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.DatabaseConnection")).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Repository error",
			req: &request.CreateDatabaseConnectionRequest{
				Name:   "Test Connection",
				DBType: "mysql",
				Host:   "localhost",
				Port:   3306,
			},
			setupMock: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.DatabaseConnection")).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordCreateFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordCreateFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.CreateConnection(ctx, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.req.Name, result.Name)
				assert.Equal(t, tc.req.DBType, result.DBType)
				assert.Equal(t, tc.req.Host, result.Host)
				assert.Equal(t, tc.req.Port, result.Port)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetConnection(t *testing.T) {
	service, mockRepo := setupDatabaseConnectionTest(t)
	ctx := context.Background()
	connID := uint(1)
	ownerID := uint(1)

	testCases := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Get connection",
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					Name:    "Test Connection",
					Code:    "db_conn_test123",
					DBType:  "mysql",
					Host:    "localhost",
					Port:    3306,
					OwnerID: ownerID,
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Connection not found",
			setupMock: func() {
				mockRepo.On("GetByID", ctx, connID).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Access denied (different owner)",
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					Name:    "Test Connection",
					OwnerID: uint(999), // Different owner
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetConnection(ctx, connID, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, connID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetConnectionList(t *testing.T) {
	service, mockRepo := setupDatabaseConnectionTest(t)
	ctx := context.Background()
	ownerID := uint(1)

	testCases := []struct {
		name          string
		req           *request.GetDatabaseConnectionListRequest
		setupMock     func()
		expectedError error
		expectedTotal int64
	}{
		{
			name: "Success - Get list",
			req: &request.GetDatabaseConnectionListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				connections := []*model.DatabaseConnection{
					{
						ID:      1,
						Name:    "Connection 1",
						Code:    "db_conn_test1",
						DBType:  "mysql",
						OwnerID: ownerID,
					},
					{
						ID:      2,
						Name:    "Connection 2",
						Code:    "db_conn_test2",
						DBType:  "postgresql",
						OwnerID: ownerID,
					},
				}
				mockRepo.On("GetList", ctx, mock.AnythingOfType("*request.GetDatabaseConnectionListRequest"), ownerID).
					Return(connections, int64(2), nil).Once()
			},
			expectedError: nil,
			expectedTotal: 2,
		},
		{
			name: "Success - Empty list",
			req: &request.GetDatabaseConnectionListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockRepo.On("GetList", ctx, mock.AnythingOfType("*request.GetDatabaseConnectionListRequest"), ownerID).
					Return([]*model.DatabaseConnection{}, int64(0), nil).Once()
			},
			expectedError: nil,
			expectedTotal: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetConnectionList(ctx, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedTotal, result.Total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateConnection(t *testing.T) {
	service, mockRepo := setupDatabaseConnectionTest(t)
	ctx := context.Background()
	connID := uint(1)
	ownerID := uint(1)
	newName := "Updated Connection"
	newPort := 3307

	testCases := []struct {
		name          string
		req           *request.UpdateDatabaseConnectionRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Update connection",
			req: &request.UpdateDatabaseConnectionRequest{
				Name: &newName,
				Port: &newPort,
			},
			setupMock: func() {
				existingConn := &model.DatabaseConnection{
					ID:      connID,
					Name:    "Old Connection",
					Code:    "db_conn_test",
					DBType:  "mysql",
					Host:    "localhost",
					Port:    3306,
					OwnerID: ownerID,
				}
				updatedConn := &model.DatabaseConnection{
					ID:      connID,
					Name:    newName,
					Code:    "db_conn_test",
					DBType:  "mysql",
					Host:    "localhost",
					Port:    newPort,
					OwnerID: ownerID,
				}
				mockRepo.On("GetByID", ctx, connID).Return(existingConn, nil).Once()
				mockRepo.On("Update", ctx, mock.AnythingOfType("*model.DatabaseConnection")).Return(nil).Once()
				mockRepo.On("GetByID", ctx, connID).Return(updatedConn, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Connection not found",
			req: &request.UpdateDatabaseConnectionRequest{
				Name: &newName,
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, connID).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Access denied",
			req: &request.UpdateDatabaseConnectionRequest{
				Name: &newName,
			},
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					OwnerID: uint(999), // Different owner
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.UpdateConnection(ctx, connID, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.req.Name != nil {
					assert.Equal(t, *tc.req.Name, result.Name)
				}
				if tc.req.Port != nil {
					assert.Equal(t, *tc.req.Port, result.Port)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteConnection(t *testing.T) {
	service, mockRepo := setupDatabaseConnectionTest(t)
	ctx := context.Background()
	connID := uint(1)
	ownerID := uint(1)

	testCases := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Delete connection",
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					Name:    "Test Connection",
					OwnerID: ownerID,
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
				mockRepo.On("Delete", ctx, connID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Connection not found",
			setupMock: func() {
				mockRepo.On("GetByID", ctx, connID).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Access denied",
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					OwnerID: uint(999), // Different owner
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
		{
			name: "Failure - Delete error",
			setupMock: func() {
				conn := &model.DatabaseConnection{
					ID:      connID,
					OwnerID: ownerID,
				}
				mockRepo.On("GetByID", ctx, connID).Return(conn, nil).Once()
				mockRepo.On("Delete", ctx, connID).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordDeleteFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordDeleteFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := service.DeleteConnection(ctx, connID, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name          string
		config        map[string]interface{}
		expectedError bool
	}{
		{
			name: "Valid config",
			config: map[string]interface{}{
				"charset": "utf8mb4",
				"timeout": 30,
			},
			expectedError: false,
		},
		{
			name:          "Empty config",
			config:        map[string]interface{}{},
			expectedError: false,
		},
		{
			name:          "Nil config",
			config:        nil,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConfig(tc.config)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
