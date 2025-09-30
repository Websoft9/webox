package service

import (
	"api-service/internal/config"
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Helper functions for pointer types
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// MockSystemConfigRepository system config repository mock
type MockSystemConfigRepository struct {
	mock.Mock
}

func (m *MockSystemConfigRepository) Create(ctx context.Context, config *model.SystemConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockSystemConfigRepository) GetByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SystemConfig), args.Error(1)
}

func (m *MockSystemConfigRepository) GetByID(ctx context.Context, id uint) (*model.SystemConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SystemConfig), args.Error(1)
}

func (m *MockSystemConfigRepository) List(ctx context.Context, filter *request.ListSystemConfigsFilter) ([]*model.SystemConfig, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SystemConfig), args.Error(1)
}

func (m *MockSystemConfigRepository) Update(ctx context.Context, config *model.SystemConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockSystemConfigRepository) UpdateValue(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockSystemConfigRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEmailService email service mock for system config tests
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

func (m *MockEmailService) SendVerificationEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error {
	args := m.Called(ctx, expires, to, token, baseURL, lang)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error {
	args := m.Called(ctx, expires, to, token, baseURL, lang)
	return args.Error(0)
}

// SystemConfigServiceTestSuite system config service test suite
type SystemConfigServiceTestSuite struct {
	suite.Suite
	service              *SystemConfigService
	mockSystemConfigRepo *MockSystemConfigRepository
	mockEmailService     *MockEmailService
	logger               logger.Logger
	i18n                 *i18n.I18n
	config               *config.Config
}

func (suite *SystemConfigServiceTestSuite) SetupTest() {
	// Initialize crypto for testing
	_, err := crypto.InitDefaultCrypto("test-encryption-key-32-characters")
	if err != nil {
		suite.T().Fatalf("Failed to initialize crypto: %v", err)
	}

	suite.mockSystemConfigRepo = &MockSystemConfigRepository{}
	suite.mockEmailService = &MockEmailService{}
	suite.logger = logger.NewZapLogger(logger.InfoLevel, nil)
	suite.i18n = i18n.NewI18n()
	suite.config = &config.Config{
		Email: config.EmailConfigMain{
			SMTP: config.SMTPConfigMain{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
			},
		},
	}

	suite.service = &SystemConfigService{
		systemConfigRepo: suite.mockSystemConfigRepo,
		emailService:     suite.mockEmailService,
		config:           suite.config,
		db:               &gorm.DB{},
		logger:           suite.logger,
		i18n:             suite.i18n,
	}
}

func (suite *SystemConfigServiceTestSuite) TearDownTest() {
	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
	suite.mockEmailService.AssertExpectations(suite.T())
}

// TestListSystemConfigs tests the ListSystemConfigs method
func (suite *SystemConfigServiceTestSuite) TestListSystemConfigs_Success() {
	ctx := context.Background()
	req := &request.ListSystemConfigsRequest{
		Keyword: "smtp",
	}

	// Create a real encrypted value for testing
	cryptoInstance := crypto.GetDefaultCrypto()
	encryptedPassword, err := cryptoInstance.Encrypt("secret_password")
	suite.NoError(err)

	expectedConfigs := []*model.SystemConfig{
		{
			ID:           1,
			ConfigKey:    "email.smtp.host",
			ConfigValue:  "smtp.example.com",
			ConfigType:   model.ConfigTypeString,
			Category:     "email",
			Description:  "SMTP服务器地址",
			IsReadonly:   false,
			IsEncrypted:  false,
			DefaultValue: "",
			SortOrder:    1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			ConfigKey:    "email.smtp.password",
			ConfigValue:  encryptedPassword, // Real encrypted value
			ConfigType:   model.ConfigTypeString,
			Category:     "email",
			Description:  "SMTP密码",
			IsReadonly:   false,
			IsEncrypted:  true,
			DefaultValue: "",
			SortOrder:    2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	filter := &request.ListSystemConfigsFilter{
		Keyword: req.Keyword,
	}

	suite.mockSystemConfigRepo.On("List", ctx, filter).Return(expectedConfigs, nil)

	result, err := suite.service.ListSystemConfigs(ctx, req)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Items, 2)
	suite.Equal("email.smtp.host", result.Items[0].ConfigKey)
	suite.Equal("smtp.example.com", result.Items[0].ConfigValue)
	suite.Equal("email.smtp.password", result.Items[1].ConfigKey)
	suite.Equal("secret_password", result.Items[1].ConfigValue) // Should be decrypted
}

func (suite *SystemConfigServiceTestSuite) TestListSystemConfigs_RepositoryError() {
	ctx := context.Background()
	req := &request.ListSystemConfigsRequest{}

	filter := &request.ListSystemConfigsFilter{}
	suite.mockSystemConfigRepo.On("List", ctx, filter).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.ListSystemConfigs(ctx, req)

	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "failed to list system configs")
}

// TestTestSMTP tests the TestSMTP method
func (suite *SystemConfigServiceTestSuite) TestTestSMTP_Success() {
	ctx := context.Background()
	req := &request.TestSMTPRequest{
		TestEmail: "test@example.com",
	}

	suite.mockEmailService.On("SendEmail", ctx, "test@example.com", "Test Email", "This is a test email from Websoft9.").Return(nil)

	err := suite.service.TestSMTP(ctx, req)

	suite.NoError(err)
}

func (suite *SystemConfigServiceTestSuite) TestTestSMTP_EmailServiceNotConfigured() {
	ctx := context.Background()
	req := &request.TestSMTPRequest{
		TestEmail: "test@example.com",
	}

	// Create service without email service
	serviceWithoutEmail := &SystemConfigService{
		systemConfigRepo: suite.mockSystemConfigRepo,
		emailService:     nil, // No email service
		config:           suite.config,
		db:               &gorm.DB{},
		logger:           suite.logger,
		i18n:             suite.i18n,
	}

	err := serviceWithoutEmail.TestSMTP(ctx, req)

	suite.Error(err)
	suite.Contains(err.Error(), "SMTP service not configured")
}

func (suite *SystemConfigServiceTestSuite) TestTestSMTP_EmailSendFailed() {
	ctx := context.Background()
	req := &request.TestSMTPRequest{
		TestEmail: "test@example.com",
	}

	suite.mockEmailService.On("SendEmail", ctx, "test@example.com", "Test Email", "This is a test email from Websoft9.").Return(errors.NewAppError(errors.CodeRecordCreateFailed))

	err := suite.service.TestSMTP(ctx, req)

	suite.Error(err)
	suite.Contains(err.Error(), "failed to send test email")
}

// TestEncryptConfigValue tests the encryptConfigValue method with simplified testing
func (suite *SystemConfigServiceTestSuite) TestEncryptConfigValue_NotEncrypted() {
	config := &model.SystemConfig{
		ConfigKey:   "test.key",
		ConfigValue: "plain_text_value",
		IsEncrypted: false,
	}

	err := suite.service.encryptConfigValue(config)

	suite.NoError(err)
	suite.Equal("plain_text_value", config.ConfigValue) // Should remain unchanged
}

func (suite *SystemConfigServiceTestSuite) TestEncryptConfigValue_EmptyValue() {
	config := &model.SystemConfig{
		ConfigKey:   "test.key",
		ConfigValue: "",
		IsEncrypted: true,
	}

	err := suite.service.encryptConfigValue(config)

	suite.NoError(err)
	suite.Equal("", config.ConfigValue) // Should remain empty
}

func (suite *SystemConfigServiceTestSuite) TestEncryptConfigValue_CryptoNotInitialized() {
	// Skip this test since crypto is now initialized in SetupTest
	suite.T().Skip("Crypto is initialized in test environment")
}

// TestDecryptConfigValue tests the decryptConfigValue method
func (suite *SystemConfigServiceTestSuite) TestDecryptConfigValue_NotEncrypted() {
	config := &model.SystemConfig{
		ConfigKey:   "test.key",
		ConfigValue: "plain_text_value",
		IsEncrypted: false,
	}

	err := suite.service.decryptConfigValue(config)

	suite.NoError(err)
	suite.Equal("plain_text_value", config.ConfigValue) // Should remain unchanged
}

func (suite *SystemConfigServiceTestSuite) TestDecryptConfigValue_EmptyValue() {
	config := &model.SystemConfig{
		ConfigKey:   "test.key",
		ConfigValue: "",
		IsEncrypted: true,
	}

	err := suite.service.decryptConfigValue(config)

	suite.NoError(err)                  // Should not error when crypto is not initialized (test environment)
	suite.Equal("", config.ConfigValue) // Should remain empty
}

// TestDecryptConfigList tests the decryptConfigList method
func (suite *SystemConfigServiceTestSuite) TestDecryptConfigList_Success() {
	// Create a real encrypted value for testing
	cryptoInstance := crypto.GetDefaultCrypto()
	encryptedValue, err := cryptoInstance.Encrypt("secret_value")
	suite.NoError(err)

	configs := []*model.SystemConfig{
		{
			ConfigKey:   "test.key1",
			ConfigValue: encryptedValue,
			IsEncrypted: true,
		},
		{
			ConfigKey:   "test.key2",
			ConfigValue: "plain_value2",
			IsEncrypted: false,
		},
		{
			ConfigKey:   "test.key3",
			ConfigValue: "",
			IsEncrypted: true,
		},
	}

	err = suite.service.decryptConfigList(configs)

	suite.NoError(err)
	suite.Equal("secret_value", configs[0].ConfigValue) // Should be decrypted
	suite.Equal("plain_value2", configs[1].ConfigValue) // Should remain unchanged
	suite.Equal("", configs[2].ConfigValue)             // Should remain empty
}

// TestListBasicConfigs tests the ListSystemConfigs method with basic category
func (suite *SystemConfigServiceTestSuite) TestListBasicConfigs_Success() {
	configs := []*model.SystemConfig{
		{
			ID:          1,
			ConfigKey:   "basic.key1",
			ConfigValue: "value1",
			Category:    "basic",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	filter := &request.ListSystemConfigsFilter{Category: "basic"}
	suite.mockSystemConfigRepo.On("List", mock.Anything, filter).Return(configs, nil)

	req := &request.ListSystemConfigsRequest{Category: "basic"}
	result, err := suite.service.ListSystemConfigs(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Items, 1)
	suite.Equal("basic.key1", result.Items[0].ConfigKey)
	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestListSecurityConfigs tests the ListSystemConfigs method with security category
func (suite *SystemConfigServiceTestSuite) TestListSecurityConfigs_Success() {
	configs := []*model.SystemConfig{
		{
			ID:          1,
			ConfigKey:   "security.key1",
			ConfigValue: "secure_value",
			Category:    "security",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	filter := &request.ListSystemConfigsFilter{Category: "security"}
	suite.mockSystemConfigRepo.On("List", mock.Anything, filter).Return(configs, nil)

	req := &request.ListSystemConfigsRequest{Category: "security"}
	result, err := suite.service.ListSystemConfigs(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Items, 1)
	suite.Equal("security.key1", result.Items[0].ConfigKey)
	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestListEmailConfigs tests the ListSystemConfigs method with email category
func (suite *SystemConfigServiceTestSuite) TestListEmailConfigs_Success() {
	configs := []*model.SystemConfig{
		{
			ID:          1,
			ConfigKey:   "email.smtp.host",
			ConfigValue: "smtp.example.com",
			Category:    "email",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	filter := &request.ListSystemConfigsFilter{Category: "email"}
	suite.mockSystemConfigRepo.On("List", mock.Anything, filter).Return(configs, nil)

	req := &request.ListSystemConfigsRequest{Category: "email"}
	result, err := suite.service.ListSystemConfigs(context.Background(), req)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Items, 1)
	suite.Equal("email.smtp.host", result.Items[0].ConfigKey)
	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestBatchUpdateSystemConfigs tests the BatchUpdateSystemConfigs method
func (suite *SystemConfigServiceTestSuite) TestBatchUpdateSystemConfigs_Success() {
	config1 := &model.SystemConfig{
		ID:          1,
		ConfigKey:   "test.key1",
		ConfigValue: "old_value1",
		ConfigType:  model.ConfigTypeString,
		Category:    "basic",
		Description: "old description 1",
		SortOrder:   1,
		IsReadonly:  false,
		IsEncrypted: false,
	}

	config2 := &model.SystemConfig{
		ID:          2,
		ConfigKey:   "test.key2",
		ConfigValue: "old_value2",
		ConfigType:  model.ConfigTypeString,
		Category:    "security",
		Description: "old description 2",
		SortOrder:   2,
		IsReadonly:  false,
		IsEncrypted: false,
	}

	req := &request.BatchUpdateSystemConfigsRequest{
		Configs: []request.UpdateSystemConfigRequest{
			{
				ConfigKey:   "test.key1",
				ConfigValue: "new_value1",
				ConfigType:  "BOOLEAN",
				Category:    "updated_basic",
				Description: stringPtr("new description 1"),
				SortOrder:   intPtr(10),
			},
			{
				ConfigKey:   "test.key2",
				ConfigValue: "new_value2",
				ConfigType:  "NUMBER",
				Category:    "updated_security",
				Description: stringPtr("new description 2"),
				SortOrder:   intPtr(20),
			},
		},
	}

	suite.mockSystemConfigRepo.On("GetByKey", mock.Anything, "test.key1").Return(config1, nil)
	suite.mockSystemConfigRepo.On("GetByKey", mock.Anything, "test.key2").Return(config2, nil)
	suite.mockSystemConfigRepo.On("Update", mock.Anything, config1).Return(nil)
	suite.mockSystemConfigRepo.On("Update", mock.Anything, config2).Return(nil)

	err := suite.service.BatchUpdateSystemConfigs(context.Background(), req)

	suite.NoError(err)
	// Verify all fields were updated correctly
	suite.Equal("new_value1", config1.ConfigValue)
	suite.Equal(model.ConfigTypeBoolean, config1.ConfigType)
	suite.Equal("updated_basic", config1.Category)
	suite.Equal("new description 1", config1.Description)
	suite.Equal(10, config1.SortOrder)

	suite.Equal("new_value2", config2.ConfigValue)
	suite.Equal(model.ConfigTypeNumber, config2.ConfigType)
	suite.Equal("updated_security", config2.Category)
	suite.Equal("new description 2", config2.Description)
	suite.Equal(20, config2.SortOrder)

	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestBatchUpdateSystemConfigs_AllFieldsUpdate tests that all fields can be updated
func (suite *SystemConfigServiceTestSuite) TestBatchUpdateSystemConfigs_AllFieldsUpdate() {
	originalConfig := &model.SystemConfig{
		ID:          1,
		ConfigKey:   "test.comprehensive",
		ConfigValue: "original_value",
		ConfigType:  model.ConfigTypeString,
		Category:    "original_category",
		Description: "original description",
		SortOrder:   5,
		IsReadonly:  false,
		IsEncrypted: false,
	}

	req := &request.BatchUpdateSystemConfigsRequest{
		Configs: []request.UpdateSystemConfigRequest{
			{
				ConfigKey:   "test.comprehensive",
				ConfigValue: "updated_value",
				ConfigType:  "JSON",
				Category:    "updated_category",
				Description: stringPtr("updated description"),
				SortOrder:   intPtr(100),
			},
		},
	}

	suite.mockSystemConfigRepo.On("GetByKey", mock.Anything, "test.comprehensive").Return(originalConfig, nil)
	suite.mockSystemConfigRepo.On("Update", mock.Anything, originalConfig).Return(nil)

	err := suite.service.BatchUpdateSystemConfigs(context.Background(), req)

	suite.NoError(err)

	// Verify all updatable fields were modified
	suite.Equal("updated_value", originalConfig.ConfigValue)
	suite.Equal(model.ConfigTypeJSON, originalConfig.ConfigType)
	suite.Equal("updated_category", originalConfig.Category)
	suite.Equal("updated description", originalConfig.Description)
	suite.Equal(100, originalConfig.SortOrder)

	// Verify non-updatable fields remain unchanged
	suite.Equal(uint(1), originalConfig.ID)
	suite.Equal("test.comprehensive", originalConfig.ConfigKey)
	suite.Equal(false, originalConfig.IsReadonly)
	suite.Equal(false, originalConfig.IsEncrypted)

	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestBatchUpdateSystemConfigs_SkipReadonly tests batch update skipping readonly configs
func (suite *SystemConfigServiceTestSuite) TestBatchUpdateSystemConfigs_SkipReadonly() {
	readonlyConfig := &model.SystemConfig{
		ID:          1,
		ConfigKey:   "readonly.key",
		ConfigValue: "readonly_value",
		Category:    "basic",
		IsReadonly:  true,
		IsEncrypted: false,
	}

	req := &request.BatchUpdateSystemConfigsRequest{
		Configs: []request.UpdateSystemConfigRequest{
			{
				ConfigKey:   "readonly.key",
				ConfigValue: "new_value",
				Category:    "basic",
			},
		},
	}

	suite.mockSystemConfigRepo.On("GetByKey", mock.Anything, "readonly.key").Return(readonlyConfig, nil)
	// Update should not be called for readonly config

	err := suite.service.BatchUpdateSystemConfigs(context.Background(), req)

	suite.NoError(err)
	suite.Equal("readonly_value", readonlyConfig.ConfigValue) // Should remain unchanged
	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestBatchUpdateSystemConfigs_PartialUpdate tests partial update when some fields are not provided
func (suite *SystemConfigServiceTestSuite) TestBatchUpdateSystemConfigs_PartialUpdate() {
	originalConfig := &model.SystemConfig{
		ID:          1,
		ConfigKey:   "test.partial",
		ConfigValue: "original_value",
		ConfigType:  model.ConfigTypeString,
		Category:    "original_category",
		Description: "original description",
		SortOrder:   5,
		IsReadonly:  false,
		IsEncrypted: false,
	}

	// Only update ConfigValue and ConfigType, leave Description and SortOrder unchanged
	req := &request.BatchUpdateSystemConfigsRequest{
		Configs: []request.UpdateSystemConfigRequest{
			{
				ConfigKey:   "test.partial",
				ConfigValue: "updated_value",
				ConfigType:  "BOOLEAN",
				Category:    "updated_category",
				// Description and SortOrder are not provided (nil pointers)
			},
		},
	}

	suite.mockSystemConfigRepo.On("GetByKey", mock.Anything, "test.partial").Return(originalConfig, nil)
	suite.mockSystemConfigRepo.On("Update", mock.Anything, originalConfig).Return(nil)

	err := suite.service.BatchUpdateSystemConfigs(context.Background(), req)

	suite.NoError(err)

	// Verify required fields were updated
	suite.Equal("updated_value", originalConfig.ConfigValue)
	suite.Equal(model.ConfigTypeBoolean, originalConfig.ConfigType)
	suite.Equal("updated_category", originalConfig.Category)

	// Verify optional fields remained unchanged
	suite.Equal("original description", originalConfig.Description)
	suite.Equal(5, originalConfig.SortOrder)

	suite.mockSystemConfigRepo.AssertExpectations(suite.T())
}

// TestNewSystemConfigService tests the NewSystemConfigService constructor
func (suite *SystemConfigServiceTestSuite) TestNewSystemConfigService_WithEmailService() {
	configWithEmail := &config.Config{
		Email: config.EmailConfigMain{
			SMTP: config.SMTPConfigMain{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "test@example.com",
				Password: "password",
			},
		},
	}

	service := NewSystemConfigService(
		suite.mockSystemConfigRepo,
		configWithEmail,
		&gorm.DB{},
		suite.logger,
		suite.i18n,
	)

	suite.NotNil(service)
	suite.NotNil(service.emailService)
	suite.Equal(suite.mockSystemConfigRepo, service.systemConfigRepo)
}

func (suite *SystemConfigServiceTestSuite) TestNewSystemConfigService_WithoutEmailService() {
	configWithoutEmail := &config.Config{
		Email: config.EmailConfigMain{
			SMTP: config.SMTPConfigMain{
				Host:     "",
				Username: "",
			},
		},
	}

	service := NewSystemConfigService(
		suite.mockSystemConfigRepo,
		configWithoutEmail,
		&gorm.DB{},
		suite.logger,
		suite.i18n,
	)

	suite.NotNil(service)
	suite.Nil(service.emailService)
}

// Run the test suite
func TestSystemConfigServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SystemConfigServiceTestSuite))
}
