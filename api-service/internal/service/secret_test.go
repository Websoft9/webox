package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-service/internal/config"
	"api-service/internal/dto/request"
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
		ID:              1,
		Name:            "Test API Key",
		KeyType:         model.SecretKeyTypeAPIKey,
		EncryptedValue:  "encrypted-value",
		Description:     &description,
		ResourceGroupID: &resourceGroupID,
		ExpiresAt:       &expiresAt,
		OwnerID:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// createTestSecretCreateRequest creates a test request for creating a secret key
func createTestSecretCreateRequest() *request.SecretKeyCreateRequest {
	description := "Test description"
	resourceGroupID := uint(1)
	expiresAt := time.Now().Add(24 * time.Hour)

	return &request.SecretKeyCreateRequest{
		Name:            "Test API Key",
		KeyType:         model.SecretKeyTypeAPIKey,
		EncryptedValue:  "test-value", // This is the plaintext value
		Description:     &description,
		ResourceGroupID: &resourceGroupID,
		ExpiresAt:       &expiresAt,
		AuthorizedUsers: []uint{1, 2, 3},
	}
}

// createTestUpdateRequest creates a test request for updating a secret key
func createTestUpdateRequest() *request.SecretKeyUpdateRequest {
	return &request.SecretKeyUpdateRequest{
		KeyType:        model.SecretKeyTypeAPIKey,
		EncryptedValue: "new-test-value", // This is the plaintext value
	}
}

func TestSecretKeyService_CreateSecretKey_Success(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateRequest()
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

	// Execute
	result, err := service.CreateSecretKey(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.KeyType, result.KeyType)
	assert.Equal(t, uint(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_CreateSecretKey_NameAlreadyExists(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateRequest()
	userID := uint(1)

	// Mock ExistsByName call returning already exists
	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(true, nil)

	// Execute
	result, err := service.CreateSecretKey(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "e")
	mockRepo.AssertExpectations(t)
}

func TestSecretKeyService_CreateSecretKey_RepositoryError(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateRequest()
	userID := uint(1)

	// Mock ExistsByName call
	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(false, nil)

	// Mock Create call returning an error
	expectedError := errors.New("database error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(expectedError)

	// Execute
	result, err := service.CreateSecretKey(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed")
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
	wrongUserID := uint(2) // Not the owner of the key
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Mock CheckUserSecretAccess call returning false (no shared access)
	mockRepo.On("CheckUserSecretAccess", ctx, keyID, wrongUserID).Return(false, nil)

	// Execute
	result, err := service.GetSecretKey(ctx, keyID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "denied")
	mockRepo.AssertExpectations(t)
}

// Tests for GetSecretKeyValue
// func TestSecretKeyService_GetSecretKeyValue_Success(t *testing.T) {
// 	service, mockRepo, _ := setupSecretKeyService()
// 	ctx := context.Background()
// 	testKey := createTestSecretKey()
// 	userID := uint(1)
// 	keyID := uint(1)

// 	// Create a test config and RSA crypto with the same keys used in service
// 	testConfig := createTestConfig()
// 	rsaCrypto, err := crypto.NewRSACryptoFromKeys(testConfig.Security.RSAPrivateKey, testConfig.Security.RSAPublicKey)
// 	if err != nil {
// 		t.Fatalf("Failed to create RSA crypto: %v", err)
// 	}

// 	// Create a real encrypted value for testing using the same keys
// 	plainValue := "test-secret-value"
// 	encryptedValue, err := rsaCrypto.EncryptString(plainValue)
// 	if err != nil {
// 		t.Fatalf("Failed to encrypt test value: %v", err)
// 	}
// 	testKey.EncryptedValue = encryptedValue

// 	// Mock GetByID call
// 	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

// 	// Execute
// 	result, err := service.GetSecretKeyValue(ctx, keyID, userID)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	// Now we can correctly compare the decrypted result
// 	assert.Equal(t, plainValue, result.Value)
// 	mockRepo.AssertExpectations(t)
// }

// func TestSecretKeyService_GetSecretKeyValue_Expired(t *testing.T) {
// 	service, mockRepo, _ := setupSecretKeyService()
// 	ctx := context.Background()
// 	testKey := createTestSecretKey()
// 	userID := uint(1)
// 	keyID := uint(1)

// 	// Set expiration time to the past
// 	expiredTime := time.Now().Add(-24 * time.Hour)
// 	testKey.ExpiresAt = &expiredTime

// 	// Mock GetByID call
// 	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

// 	// Execute
// 	result, err := service.GetSecretKeyValue(ctx, keyID, userID)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.Contains(t, err.Error(), "expired")
// 	mockRepo.AssertExpectations(t)
// }

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
	assert.Equal(t, keyID, result.ID)
	assert.Equal(t, updateReq.KeyType, result.KeyType)
	mockRepo.AssertExpectations(t)
}

// Tests for DeleteSecretKey
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
		Page:     1,
		PageSize: 10,
	}

	// Mock List call
	mockRepo.On("List", ctx, req, userID).Return(testKeys, int64(1), nil)

	// Execute
	result, err := service.ListSecretKeys(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, testKeys[0].Name, result.Items[0].Name)
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

	// Mock List call
	mockRepo.On("List", ctx, mock.AnythingOfType("*request.SecretKeyQueryRequest"), userID).Return(testKeys, int64(1), nil)

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

	// Mock List call
	mockRepo.On("List", ctx, mock.AnythingOfType("*request.SecretKeyQueryRequest"), userID).Return(testKeys, int64(1), nil)

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

// func TestSecretKeyService_ExportSecretKeys_UnsupportedFormat(t *testing.T) {
// 	service, mockRepo, _ := setupSecretKeyService()
// 	ctx := context.Background()
// 	userID := uint(1)
// 	req := &request.SecretKeyExportRequest{
// 		Format: "xml", // Unsupported format
// 	}

// 	// Execute
// 	data, filename, err := service.ExportSecretKeys(ctx, req, userID)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Nil(t, data)
// 	assert.Empty(t, filename)
// 	assert.Contains(t, err.Error(), "invalid_export_format")
// 	// List should not be called
// 	mockRepo.AssertNotCalled(t, "List")
// }

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
	wrongUserID := uint(2) // Not the owner of the key
	keyID := uint(1)

	// Mock GetByID call
	mockRepo.On("GetByID", ctx, keyID).Return(testKey, nil)

	// Mock CheckUserSecretAccess call returning false (no shared access)
	mockRepo.On("CheckUserSecretAccess", ctx, keyID, wrongUserID).Return(false, nil)

	// Execute
	err := service.ValidateSecretKeyOwnership(ctx, keyID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "e")
	mockRepo.AssertExpectations(t)
}

// Integration tests
func TestSecretKeyService_FullWorkflow(t *testing.T) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	userID := uint(1)

	// Step 1: Create secret key
	createReq := createTestSecretCreateRequest()
	// Mock ExistsByName call
	mockRepo.On("ExistsByName", ctx, createReq.Name, userID, []uint(nil)).Return(false, nil)

	// Mock Create call
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil).Run(func(args mock.Arguments) {
		secretKey := args.Get(1).(*model.SecretKey)
		secretKey.ID = 1 // Set ID to simulate database auto-increment
	})

	createResult, err := service.CreateSecretKey(ctx, createReq, userID)
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
	listReq := &request.SecretKeyQueryRequest{Page: 1, PageSize: 10}
	mockRepo.On("List", ctx, listReq, userID).Return([]*model.SecretKey{testKey}, int64(1), nil)

	listResult, err := service.ListSecretKeys(ctx, listReq, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), listResult.Total)

	// Step 5: Delete secret key
	mockRepo.On("DeleteUserSecretsBySecretKeyID", ctx, uint(1)).Return(nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err = service.DeleteSecretKey(ctx, 1, userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkSecretKeyService_CreateSecretKey(b *testing.B) {
	service, mockRepo, _ := setupSecretKeyService()
	ctx := context.Background()
	req := createTestSecretCreateRequest()
	userID := uint(1)

	mockRepo.On("ExistsByName", ctx, req.Name, userID, []uint(nil)).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SecretKey")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateSecretKey(ctx, req, userID)
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
