package service

import (
	"context"
	"testing"

	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCredentialRepository is a mock implementation of CredentialRepository
type MockCredentialRepository struct {
	mock.Mock
}

func (m *MockCredentialRepository) Create(ctx context.Context, credential *model.Credential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockCredentialRepository) Update(ctx context.Context, credential *model.Credential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockCredentialRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCredentialRepository) GetByID(ctx context.Context, id uint) (*model.Credential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Credential), args.Error(1)
}

func (m *MockCredentialRepository) GetByName(ctx context.Context, name string) (*model.Credential, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Credential), args.Error(1)
}

func (m *MockCredentialRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCredentialRepository) List(ctx context.Context, req *request.ListCredentialsRequest) ([]*model.Credential, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Credential), args.Get(1).(int64), args.Error(2)
}

// MockCredentialCategoryRepository is a mock implementation
type MockCredentialCategoryRepository struct {
	mock.Mock
}

func (m *MockCredentialCategoryRepository) GetByID(ctx context.Context, id uint) (*model.CredentialCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CredentialCategory), args.Error(1)
}

func (m *MockCredentialCategoryRepository) List(ctx context.Context) ([]*model.CredentialCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CredentialCategory), args.Error(1)
}

// MockCredentialTemplateRepository is a mock implementation
type MockCredentialTemplateRepository struct {
	mock.Mock
}

func (m *MockCredentialTemplateRepository) GetByID(ctx context.Context, id uint) (*model.CredentialTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CredentialTemplate), args.Error(1)
}

func (m *MockCredentialTemplateRepository) List(ctx context.Context, categoryID *uint) ([]*model.CredentialTemplate, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CredentialTemplate), args.Error(1)
}

// TestValidateCredentialName tests credential name validation
func TestValidateCredentialName(t *testing.T) {
	// Setup
	mockRepo := new(MockCredentialRepository)
	mockCategoryRepo := new(MockCredentialCategoryRepository)
	mockTemplateRepo := new(MockCredentialTemplateRepository)

	cryptoInstance, err := crypto.NewAESCrypto("test-secret-key")
	assert.NoError(t, err)

	log := logger.NewZapLogger(logger.InfoLevel, nil)

	service, err := NewCredentialService(
		mockRepo,
		mockCategoryRepo,
		mockTemplateRepo,
		cryptoInstance,
		log,
	)
	assert.NoError(t, err)

	credService := service.(*credentialService)

	// Test cases
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "Valid name with alphanumeric and underscore",
			input:     "my_database_123",
			wantError: false,
		},
		{
			name:      "Valid name minimum length",
			input:     "abc",
			wantError: false,
		},
		{
			name:      "Invalid name too short",
			input:     "ab",
			wantError: true,
		},
		{
			name:      "Invalid name too long",
			input:     "this_is_a_very_long_credential_name_that_exceeds_fifty_characters",
			wantError: true,
		},
		{
			name:      "Invalid name with special characters",
			input:     "my-database",
			wantError: true,
		},
		{
			name:      "Invalid name with spaces",
			input:     "my database",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := credService.validateCredentialName(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestResolveCredentialReference tests credential reference resolution
func TestResolveCredentialReference(t *testing.T) {
	// Setup
	mockRepo := new(MockCredentialRepository)
	mockCategoryRepo := new(MockCredentialCategoryRepository)
	mockTemplateRepo := new(MockCredentialTemplateRepository)

	cryptoInstance, err := crypto.NewAESCrypto("test-secret-key")
	assert.NoError(t, err)

	log := logger.NewZapLogger(logger.InfoLevel, nil)

	service, err := NewCredentialService(
		mockRepo,
		mockCategoryRepo,
		mockTemplateRepo,
		cryptoInstance,
		log,
	)
	assert.NoError(t, err)

	ctx := context.Background()

	// Test valid reference
	t.Run("Valid reference with encrypted parameter", func(t *testing.T) {
		// Encrypt test value
		encryptedValue, err := cryptoInstance.Encrypt("test_password")
		assert.NoError(t, err)

		// Create parameters as JSON map (model.JSON is map[string]interface{})
		// Store with numeric keys to maintain array-like structure
		paramsMap := model.JSON{
			"0": map[string]interface{}{
				"input_name":   "password",
				"input_value":  encryptedValue,
				"is_encrypted": true,
			},
		}

		credential := &model.Credential{
			ID:         1,
			Name:       "my_db",
			Parameters: paramsMap,
		}

		mockRepo.On("GetByName", ctx, "my_db").Return(credential, nil)

		// Test resolution
		result, err := service.ResolveCredentialReference(ctx, "{{ credentials.my_db.password }}")
		assert.NoError(t, err)
		assert.Equal(t, "test_password", result)

		mockRepo.AssertExpectations(t)
	})

	// Test invalid reference format
	t.Run("Invalid reference format", func(t *testing.T) {
		_, err := service.ResolveCredentialReference(ctx, "invalid_format")
		assert.Error(t, err)
	})
}
