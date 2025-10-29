package service

import (
	"context"
	"io"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	pkgErrors "api-service/pkg/errors"
	"api-service/pkg/logger"
)

// MockSecretRepository mocks SecretRepository interface
type MockSecretRepository struct {
	mock.Mock
}

func (m *MockSecretRepository) Create(ctx context.Context, secret *model.Secret) error {
	args := m.Called(ctx, secret)
	if args.Error(0) == nil && secret.ID == 0 {
		secret.ID = 1
	}
	return args.Error(0)
}

func (m *MockSecretRepository) GetByID(ctx context.Context, id uint) (*model.Secret, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Secret), args.Error(1)
}

func (m *MockSecretRepository) GetByCode(ctx context.Context, code string) (*model.Secret, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Secret), args.Error(1)
}

func (m *MockSecretRepository) List(ctx context.Context, req *request.ListSecretsRequest, userID uint) ([]*model.Secret, int64, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.Secret), args.Get(1).(int64), args.Error(2)
}

func (m *MockSecretRepository) Update(ctx context.Context, secret *model.Secret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockSecretRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSecretRepository) ExistsByName(ctx context.Context, resourceGroupID uint, name string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, resourceGroupID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecretRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecretRepository) GetReferenceCount(ctx context.Context, secretID uint) (int, error) {
	args := m.Called(ctx, secretID)
	return args.Int(0), args.Error(1)
}

// MockSecretReferenceRepository mocks SecretReferenceRepository interface
type MockSecretReferenceRepository struct {
	mock.Mock
}

func (m *MockSecretReferenceRepository) Create(ctx context.Context, reference *model.SecretReference) error {
	args := m.Called(ctx, reference)
	if args.Error(0) == nil && reference.ID == 0 {
		reference.ID = 1
	}
	return args.Error(0)
}

func (m *MockSecretReferenceRepository) ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretReference, error) {
	args := m.Called(ctx, secretID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SecretReference), args.Error(1)
}

func (m *MockSecretReferenceRepository) ListByResourceCode(ctx context.Context, resourceCode string) ([]*model.SecretReference, error) {
	args := m.Called(ctx, resourceCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SecretReference), args.Error(1)
}

func (m *MockSecretReferenceRepository) Exists(ctx context.Context, secretID uint, resourceCode string) (bool, error) {
	args := m.Called(ctx, secretID, resourceCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecretReferenceRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSecretReferenceRepository) DeleteBySecretID(ctx context.Context, secretID uint) error {
	args := m.Called(ctx, secretID)
	return args.Error(0)
}

func (m *MockSecretReferenceRepository) HasActiveReferences(ctx context.Context, secretID uint) (bool, error) {
	args := m.Called(ctx, secretID)
	return args.Bool(0), args.Error(1)
}

// MockSecretAuthorizeRepository mocks SecretAuthorizeRepository interface
type MockSecretAuthorizeRepository struct {
	mock.Mock
}

func (m *MockSecretAuthorizeRepository) Create(ctx context.Context, authorize *model.SecretAuthorize) error {
	args := m.Called(ctx, authorize)
	if args.Error(0) == nil && authorize.ID == 0 {
		authorize.ID = 1
	}
	return args.Error(0)
}

func (m *MockSecretAuthorizeRepository) BatchCreate(ctx context.Context, authorizes []*model.SecretAuthorize) error {
	args := m.Called(ctx, authorizes)
	return args.Error(0)
}

func (m *MockSecretAuthorizeRepository) ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretAuthorize, error) {
	args := m.Called(ctx, secretID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SecretAuthorize), args.Error(1)
}

func (m *MockSecretAuthorizeRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSecretAuthorizeRepository) DeleteBySecretID(ctx context.Context, secretID uint) error {
	args := m.Called(ctx, secretID)
	return args.Error(0)
}

func (m *MockSecretAuthorizeRepository) IsAuthorized(ctx context.Context, secretID, userID uint) (bool, error) {
	args := m.Called(ctx, secretID, userID)
	return args.Bool(0), args.Error(1)
}

// MockResourceTypeRepository mocks ResourceTypeRepository interface
type MockResourceTypeRepository struct {
	mock.Mock
}

func (m *MockResourceTypeRepository) GetByCode(ctx context.Context, code string) (*model.ResourceType, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceType), args.Error(1)
}

func (m *MockResourceTypeRepository) List(ctx context.Context) ([]*model.ResourceType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ResourceType), args.Error(1)
}

// setupSecretServiceTest sets up test dependencies
func setupSecretServiceTest(t *testing.T) (*secretService, *MockSecretRepository, *MockSecretReferenceRepository, *MockSecretAuthorizeRepository, *MockResourceGroupRepository, *MockResourceTypeRepository, *MockUserRepository, *gorm.DB) {
	mockSecretRepo := new(MockSecretRepository)
	mockReferenceRepo := new(MockSecretReferenceRepository)
	mockAuthorizeRepo := new(MockSecretAuthorizeRepository)
	mockResourceGroupRepo := new(MockResourceGroupRepository)
	mockResourceTypeRepo := new(MockResourceTypeRepository)
	mockUserRepo := new(MockUserRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, io.Discard)

	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate test tables
	err = db.AutoMigrate(&model.ResourceType{})
	assert.NoError(t, err)

	// Create test crypto
	testCrypto, err := crypto.NewAESCrypto("test-encryption-key-32-bytes!!")
	assert.NoError(t, err)

	// Create temp directory for file storage
	tempDir := t.TempDir()

	service := &secretService{
		db:                db,
		secretRepo:        mockSecretRepo,
		referenceRepo:     mockReferenceRepo,
		authorizeRepo:     mockAuthorizeRepo,
		resourceGroupRepo: mockResourceGroupRepo,
		resourceTypeRepo:  mockResourceTypeRepo,
		userRepo:          mockUserRepo,
		crypto:            testCrypto,
		fileStoragePath:   tempDir,
		logger:            mockLogger,
	}

	return service, mockSecretRepo, mockReferenceRepo, mockAuthorizeRepo, mockResourceGroupRepo, mockResourceTypeRepo, mockUserRepo, db
}

// TestCreateTextSecret tests the CreateTextSecret functionality
func TestCreateTextSecret(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.CreateTextSecretRequest
		ownerID       uint
		mockSetup     func(*MockSecretRepository, *MockResourceGroupRepository, *MockUserRepository, *MockSecretAuthorizeRepository, *MockSecretReferenceRepository, *MockResourceTypeRepository, *gorm.DB)
		expectError   bool
		expectedError error
	}{
		{
			name: "successfully create text secret",
			request: &request.CreateTextSecretRequest{
				ResourceGroupID: 1,
				Name:            "Test Secret",
				Description:     stringPtrForSecrets("Test Description"),
				SecretText:      "test-secret-value",
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				// Mock resource group exists
				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{
					ID:   1,
					Name: "Test Group",
				}, nil)

				// Mock name uniqueness check
				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Test Secret", (*uint)(nil)).Return(false, nil)

				// Mock code uniqueness check
				secretRepo.On("ExistsByCode", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)

				// Mock create
				secretRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)

				// Mock reference count
				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)

				// Mock owner
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "testuser",
				}, nil)
			},
			expectError: false,
		},
		{
			name: "resource group not found",
			request: &request.CreateTextSecretRequest{
				ResourceGroupID: 999,
				Name:            "Test Secret",
				SecretText:      "test-secret-value",
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rgRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound))
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "duplicate secret name",
			request: &request.CreateTextSecretRequest{
				ResourceGroupID: 1,
				Name:            "Duplicate Secret",
				SecretText:      "test-secret-value",
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{ID: 1}, nil)
				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Duplicate Secret", (*uint)(nil)).Return(true, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeResourceAlreadyExists),
		},
		{
			name: "create secret with authorized users",
			request: &request.CreateTextSecretRequest{
				ResourceGroupID: 1,
				Name:            "Secret with Auth",
				SecretText:      "test-secret-value",
				AuthorizedUsers: []uint{2, 3},
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{ID: 1}, nil)
				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Secret with Auth", (*uint)(nil)).Return(false, nil)
				secretRepo.On("ExistsByCode", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
				secretRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)

				// Mock authorized users
				userRepo.On("GetByID", mock.Anything, uint(2)).Return(&model.User{
					ID:       2,
					Username: "user2",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(3)).Return(&model.User{
					ID:       3,
					Username: "user3",
					Status:   UserStatusActive,
				}, nil)
				authRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretAuthorize")).Return(nil).Times(2)

				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Username: "owner"}, nil)
			},
			expectError: false,
		},
		{
			name: "create secret with resource reference",
			request: &request.CreateTextSecretRequest{
				ResourceGroupID: 1,
				Name:            "Secret with Resource",
				SecretText:      "test-secret-value",
				ResourceCode:    stringPtrForSecrets("server_test123"),
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{ID: 1}, nil)
				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Secret with Resource", (*uint)(nil)).Return(false, nil)
				secretRepo.On("ExistsByCode", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
				secretRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)

				// Mock resource type validation
				rtRepo.On("GetByCode", mock.Anything, "server").Return(&model.ResourceType{
					ID:    1,
					Code:  "server",
					Table: "servers",
				}, nil)

				// Insert test server record
				db.Exec("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY, code TEXT)")
				db.Exec("INSERT INTO servers (code) VALUES (?)", "server_test123")

				refRepo.On("Exists", mock.Anything, uint(1), "server_test123").Return(false, nil)
				refRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretReference")).Return(nil)

				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(1, nil)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Username: "owner"}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, refRepo, authRepo, rgRepo, rtRepo, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, rgRepo, userRepo, authRepo, refRepo, rtRepo, db)
			}

			result, err := service.CreateTextSecret(context.Background(), tt.request, tt.ownerID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, string(model.SecretTypeText), result.Type)
			}

			secretRepo.AssertExpectations(t)
			rgRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

// TestCreateAccountSecret tests the CreateAccountSecret functionality
func TestCreateAccountSecret(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.CreateAccountSecretRequest
		ownerID       uint
		mockSetup     func(*MockSecretRepository, *MockResourceGroupRepository, *MockUserRepository)
		expectError   bool
		expectedError error
	}{
		{
			name: "successfully create account secret",
			request: &request.CreateAccountSecretRequest{
				ResourceGroupID: 1,
				Name:            "Database Account",
				Description:     stringPtrForSecrets("MySQL Admin Account"),
				SecretUsername:  "admin",
				SecretPassword:  "P@ssw0rd123!",
			},
			ownerID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository) {
				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{ID: 1}, nil)
				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Database Account", (*uint)(nil)).Return(false, nil)
				secretRepo.On("ExistsByCode", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
				secretRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)
				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, Username: "owner"}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, _, _, rgRepo, _, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, rgRepo, userRepo)
			}

			result, err := service.CreateAccountSecret(context.Background(), tt.request, tt.ownerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, string(model.SecretTypeAccount), result.Type)
			}

			secretRepo.AssertExpectations(t)
		})
	}
}

// TestValidateFileUpload tests file upload validation
func TestValidateFileUpload(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		fileSize    int64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid PEM file",
			filename:    "certificate.pem",
			fileSize:    1024,
			expectError: false,
		},
		{
			name:        "valid KEY file",
			filename:    "private.key",
			fileSize:    2048,
			expectError: false,
		},
		{
			name:        "file size exceeds limit",
			filename:    "large.pem",
			fileSize:    6 * 1024 * 1024, // 6MB
			expectError: true,
		},
		{
			name:        "unsupported file type",
			filename:    "document.pdf",
			fileSize:    1024,
			expectError: true,
		},
		{
			name:        "file without extension",
			filename:    "noextension",
			fileSize:    1024,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, _, _, _, _, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			req := &request.CreateFileSecretRequest{
				SecretFile: &multipart.FileHeader{
					Filename: tt.filename,
					Size:     tt.fileSize,
				},
			}

			ext, err := service.validateFileUpload(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, ext)
			}
		})
	}
}

// TestValidateResourceCode tests resource code validation
func TestValidateResourceCode(t *testing.T) {
	tests := []struct {
		name          string
		resourceCode  string
		mockSetup     func(*MockResourceTypeRepository, *gorm.DB)
		expectError   bool
		expectedError error
	}{
		{
			name:         "valid server resource code",
			resourceCode: "server_abc123",
			mockSetup: func(rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rtRepo.On("GetByCode", mock.Anything, "server").Return(&model.ResourceType{
					ID:    1,
					Code:  "server",
					Table: "servers",
				}, nil)

				db.Exec("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY, code TEXT)")
				db.Exec("INSERT INTO servers (code) VALUES (?)", "server_abc123")
			},
			expectError: false,
		},
		{
			name:         "invalid resource code format",
			resourceCode: "invalidcode",
			mockSetup:    func(rtRepo *MockResourceTypeRepository, db *gorm.DB) {},
			expectError:  true,
		},
		{
			name:         "empty resource code",
			resourceCode: "",
			mockSetup:    func(rtRepo *MockResourceTypeRepository, db *gorm.DB) {},
			expectError:  true,
		},
		{
			name:         "resource type not found",
			resourceCode: "unknown_abc123",
			mockSetup: func(rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rtRepo.On("GetByCode", mock.Anything, "unknown").Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound))
			},
			expectError: true,
		},
		{
			name:         "resource not found",
			resourceCode: "server_notexist",
			mockSetup: func(rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				rtRepo.On("GetByCode", mock.Anything, "server").Return(&model.ResourceType{
					ID:    1,
					Code:  "server",
					Table: "servers",
				}, nil)

				db.Exec("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY, code TEXT)")
				// 不插入记录，模拟资源不存在
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, _, _, rtRepo, _, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(rtRepo, db)
			}

			err := service.validateResourceCode(context.Background(), tt.resourceCode)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			rtRepo.AssertExpectations(t)
		})
	}
}

// TestGetSecret tests the GetSecret functionality
func TestGetSecret(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		secretID      uint
		userID        uint
		mockSetup     func(*MockSecretRepository, *MockResourceGroupRepository, *MockUserRepository, *MockSecretReferenceRepository, *MockSecretAuthorizeRepository)
		expectError   bool
		expectedError error
	}{
		{
			name:     "owner views secret details",
			secretID: 1,
			userID:   1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					Description:     stringPtrForSecrets("Test Description"),
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{
					ID:   1,
					Name: "Test Group",
				}, nil)

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil)

				refRepo.On("ListBySecretID", mock.Anything, uint(1)).Return([]*model.SecretReference{}, nil)
				authRepo.On("ListBySecretID", mock.Anything, uint(1)).Return([]*model.SecretAuthorize{}, nil)
			},
			expectError: false,
		},
		{
			name:     "authorized user views secret details",
			secretID: 1,
			userID:   2,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				authRepo.On("IsAuthorized", mock.Anything, uint(1), uint(2)).Return(true, nil)

				rgRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.ResourceGroup{
					ID:   1,
					Name: "Test Group",
				}, nil)

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil)

				refRepo.On("ListBySecretID", mock.Anything, uint(1)).Return([]*model.SecretReference{}, nil)
				authRepo.On("ListBySecretID", mock.Anything, uint(1)).Return([]*model.SecretAuthorize{}, nil)
			},
			expectError: false,
		},
		{
			name:     "unauthorized user access denied",
			secretID: 1,
			userID:   3,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				authRepo.On("IsAuthorized", mock.Anything, uint(1), uint(3)).Return(false, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
		{
			name:     "secret not found",
			secretID: 999,
			userID:   1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, userRepo *MockUserRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound))
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, refRepo, authRepo, rgRepo, _, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, rgRepo, userRepo, refRepo, authRepo)
			}

			result, err := service.GetSecret(context.Background(), tt.secretID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.secretID, result.ID)
			}

			secretRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateSecret tests the UpdateSecret functionality
func TestUpdateSecret(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		secretID      uint
		request       *request.UpdateSecretRequest
		userID        uint
		mockSetup     func(*MockSecretRepository, *MockResourceGroupRepository, *MockSecretAuthorizeRepository, *MockUserRepository)
		expectError   bool
		expectedError error
	}{
		{
			name:     "successfully update secret name",
			secretID: 1,
			request: &request.UpdateSecretRequest{
				Name: stringPtrForSecrets("Updated Secret Name"),
			},
			userID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, authRepo *MockSecretAuthorizeRepository, userRepo *MockUserRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Original Name",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				secretRepo.On("ExistsByName", mock.Anything, uint(1), "Updated Secret Name", uintPtrForSecrets(1)).Return(false, nil)
				secretRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)
				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)

				// Mock user repo for owner lookup in toResponse
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil)
			},
			expectError: false,
		},
		{
			name:     "non-owner cannot update",
			secretID: 1,
			request: &request.UpdateSecretRequest{
				Name: stringPtrForSecrets("Updated Name"),
			},
			userID: 2,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, authRepo *MockSecretAuthorizeRepository, userRepo *MockUserRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Original Name",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
		{
			name:     "update authorized users list",
			secretID: 1,
			request: &request.UpdateSecretRequest{
				AuthorizedUsers: &[]uint{2, 3, 4},
			},
			userID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, rgRepo *MockResourceGroupRepository, authRepo *MockSecretAuthorizeRepository, userRepo *MockUserRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				authRepo.On("DeleteBySecretID", mock.Anything, uint(1)).Return(nil)

				// Mock authorized users
				userRepo.On("GetByID", mock.Anything, uint(2)).Return(&model.User{
					ID:       2,
					Username: "user2",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(3)).Return(&model.User{
					ID:       3,
					Username: "user3",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(4)).Return(&model.User{
					ID:       4,
					Username: "user4",
					Status:   UserStatusActive,
				}, nil)
				authRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretAuthorize")).Return(nil).Times(3)

				secretRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Secret")).Return(nil)
				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)

				// Mock owner for toResponse
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, _, authRepo, rgRepo, _, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, rgRepo, authRepo, userRepo)
			}

			result, err := service.UpdateSecret(context.Background(), tt.secretID, tt.request, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			secretRepo.AssertExpectations(t)
		})
	}
}

// TestDeleteSecret tests the DeleteSecret functionality with transaction
func TestDeleteSecret(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		secretID      uint
		userID        uint
		mockSetup     func(*MockSecretRepository, *MockSecretReferenceRepository, *MockSecretAuthorizeRepository)
		expectError   bool
		expectedError error
	}{
		{
			name:     "successfully delete secret without references",
			secretID: 1,
			userID:   1,
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					SecretFields:    model.JSON{},
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				refRepo.On("HasActiveReferences", mock.Anything, uint(1)).Return(false, nil)
				authRepo.On("DeleteBySecretID", mock.Anything, uint(1)).Return(nil)
				secretRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "fail to delete secret with references",
			secretID: 1,
			userID:   1,
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					SecretFields:    model.JSON{},
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)

				refRepo.On("HasActiveReferences", mock.Anything, uint(1)).Return(true, nil)
				refRepo.On("ListBySecretID", mock.Anything, uint(1)).Return([]*model.SecretReference{
					{ID: 1, SecretID: 1, ResourceCode: "server_abc123"},
				}, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppErrorWithI18n(pkgErrors.CodeResourceInUse, "secret.has_active_references"),
		},
		{
			name:     "non-owner cannot delete",
			secretID: 1,
			userID:   2,
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test Secret",
					Type:            model.SecretTypeText,
					ResourceGroupID: 1,
					OwnerID:         1,
					CreatedAt:       now,
					UpdatedAt:       now,
				}, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeAccessDenied),
		},
		{
			name:     "delete file type secret",
			secretID: 1,
			userID:   1,
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, authRepo *MockSecretAuthorizeRepository) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:              1,
					Code:            "secrets_test123",
					Name:            "Test File Secret",
					Type:            model.SecretTypeFile,
					ResourceGroupID: 1,
					OwnerID:         1,
					SecretFields: model.JSON{
						"secret_filename": "test_file.pem",
					},
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				refRepo.On("HasActiveReferences", mock.Anything, uint(1)).Return(false, nil)
				authRepo.On("DeleteBySecretID", mock.Anything, uint(1)).Return(nil)
				secretRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, refRepo, authRepo, _, _, _, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, refRepo, authRepo)
			}

			err := service.DeleteSecret(context.Background(), tt.secretID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}

			secretRepo.AssertExpectations(t)
			refRepo.AssertExpectations(t)
		})
	}
}

// TestCreateReference tests the CreateReference functionality
func TestCreateReference(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.CreateReferenceRequest
		mockSetup     func(*MockSecretRepository, *MockSecretReferenceRepository, *MockResourceTypeRepository, *gorm.DB)
		expectError   bool
		expectedError error
	}{
		{
			name: "successfully create secret reference",
			request: &request.CreateReferenceRequest{
				SecretID:     1,
				ResourceCode: "server_test123",
			},
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:   1,
					Name: "Test Secret",
				}, nil)

				// Mock resource type validation
				rtRepo.On("GetByCode", mock.Anything, "server").Return(&model.ResourceType{
					ID:    1,
					Code:  "server",
					Table: "servers",
				}, nil)

				db.Exec("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY, code TEXT)")
				db.Exec("INSERT INTO servers (code) VALUES (?)", "server_test123")

				refRepo.On("Exists", mock.Anything, uint(1), "server_test123").Return(false, nil)
				refRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretReference")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "secret not found",
			request: &request.CreateReferenceRequest{
				SecretID:     999,
				ResourceCode: "server_test123",
			},
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				secretRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound))
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "reference already exists",
			request: &request.CreateReferenceRequest{
				SecretID:     1,
				ResourceCode: "server_test123",
			},
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:   1,
					Name: "Test Secret",
				}, nil)

				rtRepo.On("GetByCode", mock.Anything, "server").Return(&model.ResourceType{
					ID:    1,
					Code:  "server",
					Table: "servers",
				}, nil)

				db.Exec("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY, code TEXT)")
				db.Exec("INSERT INTO servers (code) VALUES (?)", "server_test123")

				refRepo.On("Exists", mock.Anything, uint(1), "server_test123").Return(true, nil)
			},
			expectError:   true,
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeResourceAlreadyExists),
		},
		{
			name: "invalid resource code format",
			request: &request.CreateReferenceRequest{
				SecretID:     1,
				ResourceCode: "invalidcode",
			},
			mockSetup: func(secretRepo *MockSecretRepository, refRepo *MockSecretReferenceRepository, rtRepo *MockResourceTypeRepository, db *gorm.DB) {
				secretRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.Secret{
					ID:   1,
					Name: "Test Secret",
				}, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, refRepo, _, _, rtRepo, _, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, refRepo, rtRepo, db)
			}

			result, err := service.CreateReference(context.Background(), tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.SecretID, result.SecretID)
				assert.Equal(t, tt.request.ResourceCode, result.ResourceCode)
			}

			secretRepo.AssertExpectations(t)
			refRepo.AssertExpectations(t)
		})
	}
}

// TestCreateAuthorizations tests authorization creation with user status check
func TestCreateAuthorizations(t *testing.T) {
	tests := []struct {
		name      string
		secretID  uint
		userIDs   []uint
		mockSetup func(*MockUserRepository, *MockSecretAuthorizeRepository)
	}{
		{
			name:     "successfully create authorizations for active users",
			secretID: 1,
			userIDs:  []uint{2, 3},
			mockSetup: func(userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository) {
				userRepo.On("GetByID", mock.Anything, uint(2)).Return(&model.User{
					ID:       2,
					Username: "user2",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(3)).Return(&model.User{
					ID:       3,
					Username: "user3",
					Status:   UserStatusActive,
				}, nil)

				authRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretAuthorize")).Return(nil).Times(2)
			},
		},
		{
			name:     "skip inactive users",
			secretID: 1,
			userIDs:  []uint{2, 3},
			mockSetup: func(userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository) {
				userRepo.On("GetByID", mock.Anything, uint(2)).Return(&model.User{
					ID:       2,
					Username: "user2",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(3)).Return(&model.User{
					ID:       3,
					Username: "user3",
					Status:   0, // Inactive
				}, nil)

				// Only one authorization should be created
				authRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretAuthorize")).Return(nil).Once()
			},
		},
		{
			name:     "skip non-existent users",
			secretID: 1,
			userIDs:  []uint{2, 999},
			mockSetup: func(userRepo *MockUserRepository, authRepo *MockSecretAuthorizeRepository) {
				userRepo.On("GetByID", mock.Anything, uint(2)).Return(&model.User{
					ID:       2,
					Username: "user2",
					Status:   UserStatusActive,
				}, nil)
				userRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound))

				authRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.SecretAuthorize")).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, authRepo, _, _, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(userRepo, authRepo)
			}

			service.createAuthorizations(context.Background(), tt.secretID, tt.userIDs)

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
		})
	}
}

// Helper functions for secrets tests
func stringPtrForSecrets(s string) *string {
	return &s
}

func uintPtrForSecrets(u uint) *uint {
	return &u
}

// TestListSecrets tests the ListSecrets functionality
func TestListSecrets(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		request     *request.ListSecretsRequest
		userID      uint
		mockSetup   func(*MockSecretRepository, *MockUserRepository)
		expectError bool
	}{
		{
			name: "successfully list secrets with pagination",
			request: &request.ListSecretsRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			userID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, userRepo *MockUserRepository) {
				secrets := []*model.Secret{
					{
						ID:              1,
						Code:            "secrets_test1",
						Name:            "Test Secret 1",
						Type:            model.SecretTypeText,
						ResourceGroupID: 1,
						OwnerID:         1,
						CreatedAt:       now,
						UpdatedAt:       now,
					},
					{
						ID:              2,
						Code:            "secrets_test2",
						Name:            "Test Secret 2",
						Type:            model.SecretTypeAccount,
						ResourceGroupID: 1,
						OwnerID:         1,
						CreatedAt:       now,
						UpdatedAt:       now,
					},
				}

				secretRepo.On("List", mock.Anything, mock.AnythingOfType("*request.ListSecretsRequest"), uint(1)).
					Return(secrets, int64(2), nil)

				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)
				secretRepo.On("GetReferenceCount", mock.Anything, uint(2)).Return(1, nil)

				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil).Times(2)
			},
			expectError: false,
		},
		{
			name: "list secrets with keyword filter",
			request: &request.ListSecretsRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
					SearchRequest: common.SearchRequest{
						Keyword: "OpenAI",
					},
				},
			},
			userID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, userRepo *MockUserRepository) {
				secrets := []*model.Secret{
					{
						ID:              1,
						Code:            "secrets_openai",
						Name:            "OpenAI API Key",
						Type:            model.SecretTypeText,
						ResourceGroupID: 1,
						OwnerID:         1,
						CreatedAt:       now,
						UpdatedAt:       now,
					},
				}

				secretRepo.On("List", mock.Anything, mock.AnythingOfType("*request.ListSecretsRequest"), uint(1)).
					Return(secrets, int64(1), nil)

				secretRepo.On("GetReferenceCount", mock.Anything, uint(1)).Return(0, nil)
				userRepo.On("GetByID", mock.Anything, uint(1)).Return(&model.User{
					ID:       1,
					Username: "owner",
				}, nil)
			},
			expectError: false,
		},
		{
			name: "list secrets returns empty list",
			request: &request.ListSecretsRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			userID: 1,
			mockSetup: func(secretRepo *MockSecretRepository, userRepo *MockUserRepository) {
				secretRepo.On("List", mock.Anything, mock.AnythingOfType("*request.ListSecretsRequest"), uint(1)).
					Return([]*model.Secret{}, int64(0), nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, secretRepo, _, _, _, _, userRepo, db := setupSecretServiceTest(t)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tt.mockSetup != nil {
				tt.mockSetup(secretRepo, userRepo)
			}

			result, err := service.ListSecrets(context.Background(), tt.request, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			secretRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
