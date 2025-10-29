package service

import (
	"context"
	"io"
	"testing"

	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEnvironmentVariableRepository is a mock implementation of EnvironmentVariableRepository
type MockEnvironmentVariableRepository struct {
	mock.Mock
}

func (m *MockEnvironmentVariableRepository) Create(ctx context.Context, envVar *model.EnvironmentVariable) error {
	args := m.Called(ctx, envVar)
	return args.Error(0)
}

func (m *MockEnvironmentVariableRepository) GetByID(ctx context.Context, id uint) (*model.EnvironmentVariable, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EnvironmentVariable), args.Error(1)
}

func (m *MockEnvironmentVariableRepository) GetByNameAndScope(ctx context.Context, name string, scope model.EnvVarScope, projectID *uint) (*model.EnvironmentVariable, error) {
	args := m.Called(ctx, name, scope, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EnvironmentVariable), args.Error(1)
}

func (m *MockEnvironmentVariableRepository) GetPlatformList(ctx context.Context, req *request.GetEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.EnvironmentVariable), args.Get(1).(int64), args.Error(2)
}

func (m *MockEnvironmentVariableRepository) GetProjectList(ctx context.Context, req *request.GetProjectEnvVarListRequest) ([]*model.EnvironmentVariable, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.EnvironmentVariable), args.Get(1).(int64), args.Error(2)
}

func (m *MockEnvironmentVariableRepository) GetAllByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) ([]*model.EnvironmentVariable, error) {
	args := m.Called(ctx, scope, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.EnvironmentVariable), args.Error(1)
}

func (m *MockEnvironmentVariableRepository) Update(ctx context.Context, envVar *model.EnvironmentVariable) error {
	args := m.Called(ctx, envVar)
	return args.Error(0)
}

func (m *MockEnvironmentVariableRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEnvironmentVariableRepository) CountByScope(ctx context.Context, scope model.EnvVarScope, projectID *uint) (int64, error) {
	args := m.Called(ctx, scope, projectID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEnvironmentVariableRepository) ExistsByName(ctx context.Context, name string, scope model.EnvVarScope, projectID, excludeID *uint) (bool, error) {
	args := m.Called(ctx, name, scope, projectID, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestValidateEnvVarName(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)
	envService := service.(*environmentVariableService)

	tests := []struct {
		name      string
		varName   string
		wantError bool
	}{
		{
			name:      "valid uppercase name",
			varName:   "DB_HOST",
			wantError: false,
		},
		{
			name:      "valid lowercase name",
			varName:   "db_host",
			wantError: false,
		},
		{
			name:      "valid mixed case name",
			varName:   "DbHost",
			wantError: false,
		},
		{
			name:      "valid name with numbers",
			varName:   "DB_HOST_123",
			wantError: false,
		},
		{
			name:      "valid name starting with underscore",
			varName:   "_private_var",
			wantError: false,
		},
		{
			name:      "empty name",
			varName:   "",
			wantError: true,
		},
		{
			name:      "name starting with number",
			varName:   "123_VAR",
			wantError: true,
		},
		{
			name:      "name with special characters",
			varName:   "DB-HOST",
			wantError: true,
		},
		{
			name:      "name with spaces",
			varName:   "DB HOST",
			wantError: true,
		},
		{
			name:      "name too long",
			varName:   "THIS_IS_A_VERY_LONG_VARIABLE_NAME_THAT_EXCEEDS_SIXTY_FOUR_CHARACTERS_LIMIT",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := envService.validateEnvVarName(tt.varName)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreatePlatformEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()
	ownerID := uint(1)

	t.Run("successful creation", func(t *testing.T) {
		req := &request.CreatePlatformEnvVarRequest{
			CreateEnvVarRequest: request.CreateEnvVarRequest{
				Name:        "TEST_VAR",
				Value:       "test_value",
				IsSensitive: false,
			},
		}

		mockRepo.On("ExistsByName", ctx, "TEST_VAR", model.EnvVarScopePlatform, (*uint)(nil), (*uint)(nil)).Return(false, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.EnvironmentVariable")).Return(nil).Once()

		result, err := service.CreatePlatformEnvVar(ctx, req, ownerID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "TEST_VAR", result.Name)
		assert.Equal(t, "test_value", result.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid name format", func(t *testing.T) {
		req := &request.CreatePlatformEnvVarRequest{
			CreateEnvVarRequest: request.CreateEnvVarRequest{
				Name:  "123-INVALID",
				Value: "test_value",
			},
		}

		result, err := service.CreatePlatformEnvVar(ctx, req, ownerID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("duplicate name", func(t *testing.T) {
		req := &request.CreatePlatformEnvVarRequest{
			CreateEnvVarRequest: request.CreateEnvVarRequest{
				Name:  "DUPLICATE_VAR",
				Value: "test_value",
			},
		}

		mockRepo.On("ExistsByName", ctx, "DUPLICATE_VAR", model.EnvVarScopePlatform, (*uint)(nil), (*uint)(nil)).Return(true, nil).Once()

		result, err := service.CreatePlatformEnvVar(ctx, req, ownerID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateProjectEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()
	ownerID := uint(1)
	projectID := uint(100)

	t.Run("successful creation", func(t *testing.T) {
		req := &request.CreateProjectEnvVarRequest{
			CreateEnvVarRequest: request.CreateEnvVarRequest{
				Name:        "PROJECT_VAR",
				Value:       "project_value",
				IsSensitive: false,
			},
			ProjectID: projectID,
		}

		mockRepo.On("ExistsByName", ctx, "PROJECT_VAR", model.EnvVarScopeProject, &projectID, (*uint)(nil)).Return(false, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.EnvironmentVariable")).Return(nil).Once()

		result, err := service.CreateProjectEnvVar(ctx, req, ownerID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "PROJECT_VAR", result.Name)
		assert.Equal(t, "project_value", result.Value)
		assert.Equal(t, &projectID, result.ProjectID)
		mockRepo.AssertExpectations(t)
	})
}

func TestResolveEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()
	projectID := uint(100)

	t.Run("resolve simple variable", func(t *testing.T) {
		req := &request.ResolveEnvVarRequest{
			Template:  "Database host is ${DB_HOST}",
			Scope:     "platform",
			ProjectID: nil,
		}

		platformVar := &model.EnvironmentVariable{
			ID:          1,
			Name:        "DB_HOST",
			Value:       "localhost",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
		}

		mockRepo.On("GetByNameAndScope", ctx, "DB_HOST", model.EnvVarScopePlatform, (*uint)(nil)).Return(platformVar, nil).Once()

		result, err := service.ResolveEnvVar(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Database host is localhost", result.Result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("resolve with default value", func(t *testing.T) {
		req := &request.ResolveEnvVarRequest{
			Template:  "Port is ${PORT:3306}",
			Scope:     "platform",
			ProjectID: nil,
		}

		platformVar := &model.EnvironmentVariable{
			ID:          2,
			Name:        "PORT",
			Value:       "", // Empty value
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
		}

		mockRepo.On("GetByNameAndScope", ctx, "PORT", model.EnvVarScopePlatform, (*uint)(nil)).Return(platformVar, nil).Once()

		result, err := service.ResolveEnvVar(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Port is 3306", result.Result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("resolve with project variable priority", func(t *testing.T) {
		req := &request.ResolveEnvVarRequest{
			Template:  "DB is ${DB_NAME}",
			Scope:     "project",
			ProjectID: &projectID,
		}

		projectVar := &model.EnvironmentVariable{
			ID:          3,
			Name:        "DB_NAME",
			Value:       "project_db",
			Scope:       model.EnvVarScopeProject,
			ProjectID:   &projectID,
			IsSensitive: false,
		}

		mockRepo.On("GetByNameAndScope", ctx, "DB_NAME", model.EnvVarScopeProject, &projectID).Return(projectVar, nil).Once()

		result, err := service.ResolveEnvVar(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "DB is project_db", result.Result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("variable not found keeps placeholder", func(t *testing.T) {
		req := &request.ResolveEnvVarRequest{
			Template:  "Unknown var ${UNKNOWN_VAR}",
			Scope:     "platform",
			ProjectID: nil,
		}

		mockRepo.On("GetByNameAndScope", ctx, "UNKNOWN_VAR", model.EnvVarScopePlatform, (*uint)(nil)).Return(nil, assert.AnError).Once()

		result, err := service.ResolveEnvVar(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Unknown var ${UNKNOWN_VAR}", result.Result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("resolve multiple variables", func(t *testing.T) {
		req := &request.ResolveEnvVarRequest{
			Template:  "Connect to ${DB_HOST}:${DB_PORT}",
			Scope:     "platform",
			ProjectID: nil,
		}

		hostVar := &model.EnvironmentVariable{
			ID:          4,
			Name:        "DB_HOST",
			Value:       "localhost",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
		}

		portVar := &model.EnvironmentVariable{
			ID:          5,
			Name:        "DB_PORT",
			Value:       "3306",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
		}

		mockRepo.On("GetByNameAndScope", ctx, "DB_HOST", model.EnvVarScopePlatform, (*uint)(nil)).Return(hostVar, nil).Once()
		mockRepo.On("GetByNameAndScope", ctx, "DB_PORT", model.EnvVarScopePlatform, (*uint)(nil)).Return(portVar, nil).Once()

		result, err := service.ResolveEnvVar(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Connect to localhost:3306", result.Result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()

	t.Run("get existing variable", func(t *testing.T) {
		envVar := &model.EnvironmentVariable{
			ID:          1,
			Name:        "TEST_VAR",
			Value:       "test_value",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
			OwnerID:     1,
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(envVar, nil).Once()

		result, err := service.GetEnvVar(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "TEST_VAR", result.Name)
		assert.Equal(t, "test_value", result.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get sensitive variable masks value", func(t *testing.T) {
		envVar := &model.EnvironmentVariable{
			ID:          2,
			Name:        "SECRET_KEY",
			Value:       "encrypted_secret",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: true,
			OwnerID:     1,
		}

		mockRepo.On("GetByID", ctx, uint(2)).Return(envVar, nil).Once()

		result, err := service.GetEnvVar(ctx, 2)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "SECRET_KEY", result.Name)
		assert.Equal(t, "******", result.Value) // Value should be masked
		assert.True(t, result.IsSensitive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get non-existent variable", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, assert.AnError).Once()

		result, err := service.GetEnvVar(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		existingVar := &model.EnvironmentVariable{
			ID:          1,
			Name:        "TEST_VAR",
			Value:       "old_value",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
			OwnerID:     1,
		}

		newValue := "new_value"
		req := &request.UpdateEnvVarRequest{
			Value: &newValue,
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingVar, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.EnvironmentVariable")).Return(nil).Once()

		result, err := service.UpdateEnvVar(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new_value", result.Value)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteEnvVar(t *testing.T) {
	mockRepo := new(MockEnvironmentVariableRepository)
	log := logger.NewZapLogger(logger.InfoLevel, io.Discard)
	service := NewEnvironmentVariableService(mockRepo, log)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		existingVar := &model.EnvironmentVariable{
			ID:          1,
			Name:        "TEST_VAR",
			Value:       "test_value",
			Scope:       model.EnvVarScopePlatform,
			IsSensitive: false,
			OwnerID:     1,
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingVar, nil).Once()
		mockRepo.On("Delete", ctx, uint(1)).Return(nil).Once()

		err := service.DeleteEnvVar(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete non-existent variable", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, assert.AnError).Once()

		err := service.DeleteEnvVar(ctx, 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
