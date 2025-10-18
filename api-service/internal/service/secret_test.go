package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-service/internal/config"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/i18n"
)

// MockSecretKeyRepository mocks SecretKeyRepository interface
type MockSecretKeyRepository struct {
	mock.Mock
}

func (m *MockSecretKeyRepository) Create(ctx context.Context, secretKey *model.SecretKey) error {
	args := m.Called(ctx, secretKey)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) GetByID(ctx context.Context, id uint) (*model.SecretKey, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SecretKey), args.Error(1)
}

func (m *MockSecretKeyRepository) Update(ctx context.Context, secretKey *model.SecretKey) error {
	args := m.Called(ctx, secretKey)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) List(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) ([]*model.SecretKey, int64, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).([]*model.SecretKey), args.Get(1).(int64), args.Error(2)
}

func (m *MockSecretKeyRepository) ListAll(ctx context.Context, userID uint) ([]*model.SecretKey, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SecretKey), args.Error(1)
}

func (m *MockSecretKeyRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]*model.SecretKey, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]*model.SecretKey), args.Error(1)
}

func (m *MockSecretKeyRepository) ExistsByName(ctx context.Context, name string, ownerID uint, excludeID ...uint) (bool, error) {
	args := m.Called(ctx, name, ownerID, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecretKeyRepository) CountByType(ctx context.Context, keyType model.SecretKeyType, ownerID uint) (int64, error) {
	args := m.Called(ctx, keyType, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSecretKeyRepository) CreateUserSecret(ctx context.Context, userSecret *model.UserSecret) error {
	args := m.Called(ctx, userSecret)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) DeleteUserSecretsBySecretKeyID(ctx context.Context, secretKeyID uint) error {
	args := m.Called(ctx, secretKeyID)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) CheckUserSecretAccess(ctx context.Context, secretKeyID uint, userID uint) (bool, error) {
	args := m.Called(ctx, secretKeyID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecretKeyRepository) CreateSecretReference(ctx context.Context, reference *model.SecretReference) error {
	args := m.Called(ctx, reference)
	return args.Error(0)
}

func (m *MockSecretKeyRepository) DeleteSecretReferencesBySecretID(ctx context.Context, secretID uint) error {
	args := m.Called(ctx, secretID)
	return args.Error(0)
}

// createTestConfig creates a test configuration with RSA keys
func createTestConfig() *config.Config {
	// Generate test RSA key pair
	privateKeyPEM, publicKeyPEM, err := crypto.GenerateKeyPair(2048)
	if err != nil {
		panic("Failed to generate test RSA keys: " + err.Error())
	}

	return &config.Config{
		Security: config.SecurityConfig{
			AesKey:        "test-32-char-aes-key-for-testing!",
			RSAPrivateKey: privateKeyPEM,
			RSAPublicKey:  publicKeyPEM,
		},
		Upload: config.UploadConfig{
			SecretStorage: "/tmp/test-secrets",
		},
	}
}

// setupSecretKeyService creates service with mocked dependencies
func setupSecretKeyService() (service.SecretKeyService, *MockSecretKeyRepository, *MockLogger) {
	mockRepo := &MockSecretKeyRepository{}
	mockLogger := &MockLogger{}

	// Initialize i18n for testing
	_ = i18n.Init() // Initialize with default config
	mockI18n := i18n.NewI18n()

	// Setup logger mock expectations
	mockLogger.On("InfoContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("ErrorContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("DebugContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("WarnContext", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe().Return()
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Maybe().Return()

	// Create a test config with RSA keys for testing
	testConfig := createTestConfig()

	secretKeyService := NewSecretKeyService(mockRepo, mockLogger, mockI18n, testConfig)
	return secretKeyService, mockRepo, mockLogger
}

// createTestSecretKey creates a test secret key model
func createTestSecretKey() *model.SecretKey {
	description := "Test description"
	resourceGroupID := uint(1)
	expiresAt := time.Now().Add(24 * time.Hour)

	return &model.SecretKey{
		ID:          1,
		Name:        "Test API Key",
		KeyType:     model.SecretKeyTypeText,
		Description: &description,
		SecretFields: model.SecretFields{
			"secret_key": "encrypted-secret-key-value",
		},
		ResourceGroupID: &resourceGroupID,
		ExpiresAt:       &expiresAt,
		OwnerID:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// createTestSecretCreateTextRequest creates a test request for creating a text secret key
func createTestSecretCreateTextRequest() *request.SecretKeyCreateTextRequest {
	description := "Test description"
	resourceGroupID := uint(1)
	expiresAt := time.Now().Add(24 * time.Hour)

	return &request.SecretKeyCreateTextRequest{
		Name:        "Test API Key",
		Description: &description,
		SecretFields: map[string]interface{}{
			"secret_key": "sk-1234567890abcdef",
		},
		ResourceGroupID: &resourceGroupID,
		ExpiresAt:       &expiresAt,
		AuthorizedUsers: []uint{1, 2, 3},
	}
}

// createTestUpdateRequest creates a test request for updating a secret key
func createTestUpdateRequest() *request.SecretKeyUpdateRequest {
	name := "Updated Secret Key"
	description := "Updated description"
	expiresAt := time.Now().Add(90 * 24 * time.Hour)
	return &request.SecretKeyUpdateRequest{
		Name:        &name,
		Description: &description,
		ExpiresAt:   &expiresAt,
	}
}

func TestSecretKeyService_CreateSecretKeyText_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateTextRequest()
	userID := uint(1)

	// Mock ExistsByName call
	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(false, nil)

	// Mock Create call
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil).Run(func(args mock.Arguments) {
		secretKey := args.Get(1).(*model.SecretKey)
		secretKey.ID = 1 // Set ID to simulate database auto-increment
	})

	// Mock CreateUserSecret calls for each authorized user
	for range req.AuthorizedUsers {
		mockRepo.On("CreateUserSecret", ctx, mock.AnythingOfType("*model.UserSecret")).Return(nil)
	}

	// Mock CreateSecretReference call if resource_code is provided
	if req.ResourceCode != nil {
		mockRepo.On("CreateSecretReference", ctx, mock.AnythingOfType("*model.SecretReference")).Return(nil)
	}

	// Execute
	result, err := service.CreateSecretKeyText(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, uint(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_CreateSecretKeyText_NameAlreadyExists(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateTextRequest()
	userID := uint(1)

	// Mock ExistsByName call returning already exists
	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(true, nil)

	// Execute
	result, err := service.CreateSecretKeyText(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_CreateSecretKeyText_RepositoryError(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateTextRequest()
	userID := uint(1)

	// Mock ExistsByName call
	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(false, nil)

	// Mock Create call returning an error
	expectedError := errors.New("database error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(expectedError)

	// Execute
	result, err := service.CreateSecretKeyText(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

// Tests for GetSecretKey
func TestSecretKeyService_GetSecretKey_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Execute
	result, err := service.GetSecretKey(ctx, keyID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testKey.ID, result.ID)
	assert.Equal(t, testKey.Name, result.Name)
	assert.Equal(t, testKey.KeyType, result.KeyType)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_GetSecretKey_NotFound(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	userID := uint(1)
	keyID := uint(999)

	// Mock GetByID call returning not found
	mockRepo.On("GetByID", ctx, keyID).Return((*model.SecretKey)(nil), errors.New("record not found"))

	// Execute
	result, err := service.GetSecretKey(ctx, keyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_GetSecretKey_AccessDenied(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	testKey.OwnerID = uint(1)
	wrongUserID := uint(2)
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Mock CheckUserSecretAccess call returning false
	mockRepo.On("CheckUserSecretAccess", ctx, wrongUserID, keyID).Return(false, nil)

	// Execute
	result, err := service.GetSecretKey(ctx, keyID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests for GetSecretKeyValue
func TestSecretKeyService_GetSecretKeyValue_Expired(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)
	keyID := uint(1)

	// Set expiration time to the past
	expiredTime := time.Now().Add(-24 * time.Hour)
	testKey.ExpiresAt = &expiredTime

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Execute
	result, err := service.GetSecretKeyValue(ctx, keyID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests for UpdateSecretKey
func TestSecretKeyService_UpdateSecretKey_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)
	keyID := uint(1)
	updateReq := createTestUpdateRequest()

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Mock Update call
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil)

	// Execute
	result, err := service.UpdateSecretKey(ctx, keyID, userID, updateReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, *updateReq.Name, result.Name)
	assert.Equal(t, *updateReq.Description, *result.Description)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_DeleteSecretKey_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Mock DeleteUserSecretsBySecretKeyID call
	mockRepo.On("DeleteUserSecretsBySecretKeyID", ctx, keyID).Return(nil)

	// Mock DeleteSecretReferencesBySecretID call
	mockRepo.On("DeleteSecretReferencesBySecretID", ctx, keyID).Return(nil)

	// Mock Delete call
	mockRepo.On("Delete", ctx, keyID).Return(nil)

	// Execute
	err := service.DeleteSecretKey(ctx, keyID, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Tests for ListSecretKeys
func TestSecretKeyService_ListSecretKeys_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKeys := []*model.SecretKey{createTestSecretKey()}
	userID := uint(1)
	req := &request.SecretKeyQueryRequest{
		PaginationRequest: common.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	// Mock List call
	mockRepo.On("List", ctx, req, userID).Return(testKeys, int64(1), nil)

	// Execute
	result, err := service.ListSecretKeys(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)

	items, ok := result.Items.([]response.SecretKeyResponse)
	assert.True(t, ok, "Items should be of type []response.SecretKeyResponse")
	assert.Equal(t, 1, len(items))
	assert.Equal(t, testKeys[0].Name, items[0].Name)

	mockRepo.AssertExpectations(t)
}

// Tests for ExportSecretKeys
func TestSecretKeyService_ExportSecretKeys_CSV_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKeys := []*model.SecretKey{createTestSecretKey()}
	userID := uint(1)
	req := &request.SecretKeyExportRequest{
		Format: "csv",
	}

	// Mock ListAll call
	mockRepo.On("ListAll", ctx, userID).Return(testKeys, nil)

	// Execute
	data, filename, err := service.ExportSecretKeys(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, filename)
	assert.Contains(t, filename, ".csv")
	assert.Contains(t, string(data), "ID,Name,Type,Description,Created At,Updated At,Expires At")
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_ExportSecretKeys_JSON_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKeys := []*model.SecretKey{createTestSecretKey()}
	userID := uint(1)
	req := &request.SecretKeyExportRequest{
		Format: "json",
	}

	// Mock ListAll call
	mockRepo.On("ListAll", ctx, userID).Return(testKeys, nil)

	// Execute
	data, filename, err := service.ExportSecretKeys(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, filename)
	assert.Contains(t, filename, ".json")
	assert.Contains(t, string(data), "\"id\":")
	assert.Contains(t, string(data), "\"name\":")
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_ExportSecretKeys_UnsupportedFormat(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	userID := uint(1)
	req := &request.SecretKeyExportRequest{
		Format: "xml",
	}

	testKeys := []*model.SecretKey{createTestSecretKey()}
	mockRepo.On("ListAll", ctx, userID).Return(testKeys, nil)

	// Execute
	data, filename, err := service.ExportSecretKeys(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Empty(t, filename)
	mockRepo.AssertExpectations(t)
}

// Tests for ValidateSecretKeyOwnership
func TestSecretKeyService_ValidateSecretKeyOwnership_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Execute
	err := service.ValidateSecretKeyOwnership(ctx, keyID, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_ValidateSecretKeyOwnership_AccessDenied(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	wrongUserID := uint(2)
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Execute
	err := service.ValidateSecretKeyOwnership(ctx, keyID, wrongUserID)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// Integration tests
func TestSecretKeyService_FullWorkflow(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	userID := uint(1)

	// Step 1: Create secret key
	createReq := createTestSecretCreateTextRequest()
	mockRepo.On("ExistsByName", ctx, createReq.Name, userID, []uint(nil)).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil).Run(func(args mock.Arguments) {
		secretKey := args.Get(1).(*model.SecretKey)
		secretKey.ID = 1
	})

	for range createReq.AuthorizedUsers {
		mockRepo.On("CreateUserSecret", ctx, mock.AnythingOfType("*model.UserSecret")).Return(nil)
	}

	createResult, err := service.CreateSecretKeyText(ctx, createReq, userID)
	assert.NoError(t, err)
	assert.NotNil(t, createResult)

	// Step 2: Get secret key
	testKey := createTestSecretKey()
	mockRepo.On("GetByID", ctx, uint(1)).Return(testKey, nil)

	getResult, err := service.GetSecretKey(ctx, 1, userID)
	assert.NoError(t, err)
	assert.Equal(t, testKey.ID, getResult.ID)

	// Step 3: Update secret key
	updateReq := createTestUpdateRequest()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil)

	updateResult, err := service.UpdateSecretKey(ctx, 1, userID, updateReq)
	assert.NoError(t, err)
	assert.Equal(t, testKey.ID, updateResult.ID)

	// Step 4: List secret keys
	listReq := &request.SecretKeyQueryRequest{
		PaginationRequest: common.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}
	mockRepo.On("List", ctx, listReq, userID).Return([]*model.SecretKey{testKey}, int64(1), nil)

	listResult, err := service.ListSecretKeys(ctx, listReq, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), listResult.Total)
	items, ok := listResult.Items.([]response.SecretKeyResponse)
	assert.True(t, ok)
	assert.Equal(t, 1, len(items))

	// Step 5: Delete secret key
	mockRepo.On("DeleteUserSecretsBySecretKeyID", ctx, uint(1)).Return(nil)
	mockRepo.On("DeleteSecretReferencesBySecretID", ctx, uint(1)).Return(nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err = service.DeleteSecretKey(ctx, 1, userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkSecretKeyService_CreateSecretKeyText(b *testing.B) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateTextRequest()
	userID := uint(1)

	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateSecretKeyText(ctx, req, userID)
	}
}

func BenchmarkSecretKeyService_GetSecretKey(b *testing.B) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	testKey := createTestSecretKey()
	userID := uint(1)

	mockRepo.On("GetByID", ctx, uint(1)).Return(testKey, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetSecretKey(ctx, 1, userID)
	}
}
