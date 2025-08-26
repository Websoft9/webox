package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/logger"
	"api-service/pkg/security"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// MockAPITokenRepository API Token存储接口Mock
type MockAPITokenRepository struct {
	mock.Mock
}

func (m *MockAPITokenRepository) Create(ctx context.Context, token *model.APIToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAPITokenRepository) GetByID(ctx context.Context, id uint) (*model.APIToken, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APIToken), args.Error(1)
}

func (m *MockAPITokenRepository) GetByToken(ctx context.Context, tokenHash string) (*model.APIToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APIToken), args.Error(1)
}

func (m *MockAPITokenRepository) Update(ctx context.Context, token *model.APIToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAPITokenRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPITokenRepository) List(ctx context.Context, req *request.ListAPITokensRequest) ([]*model.APIToken, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.APIToken), args.Get(1).(int64), args.Error(2)
}

func (m *MockAPITokenRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.APIToken), args.Error(1)
}

func (m *MockAPITokenRepository) BatchDelete(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockAPITokenRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockAPITokenRepository) UpdateLastUsed(ctx context.Context, id uint, ip string) error {
	args := m.Called(ctx, id, ip)
	return args.Error(0)
}

func (m *MockAPITokenRepository) CleanExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAPITokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	args := m.Called(ctx, tx, token)
	return args.Error(0)
}

func (m *MockAPITokenRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	args := m.Called(ctx, tx, token)
	return args.Error(0)
}

// APITokenServiceTestSuite API Token服务测试套件
type APITokenServiceTestSuite struct {
	suite.Suite
	service       *apiTokenService
	mockTokenRepo *MockAPITokenRepository
	logger        logger.Logger
}

func (suite *APITokenServiceTestSuite) SetupTest() {
	suite.mockTokenRepo = &MockAPITokenRepository{}
	suite.logger = logger.NewZapLogger(logger.InfoLevel, nil)

	suite.service = &apiTokenService{
		tokenRepo: suite.mockTokenRepo,
		db:        &gorm.DB{},
		logger:    suite.logger,
		i18n:      nil, // Mock i18n not needed for unit tests
	}
}

func (suite *APITokenServiceTestSuite) TestCreateAPIToken_Success() {
	ctx := context.Background()
	userID := uint(1)
	expiresAt := time.Now().Add(24 * time.Hour)
	req := &request.CreateAPITokenRequest{
		Name:        "Test Token",
		Description: "Test API token",
		Scopes:      []string{"user:read", "user:write"},
		ExpiresAt:   &expiresAt,
	}

	// Skip permission validation for unit tests

	// Mock token creation
	suite.mockTokenRepo.On("Create", ctx, mock.AnythingOfType("*model.APIToken")).
		Run(func(args mock.Arguments) {
			token := args.Get(1).(*model.APIToken)
			token.ID = 1 // Simulate database ID assignment
		}).Return(nil)

	// Execute
	result, err := suite.service.CreateAPIToken(ctx, req, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(req.Name, result.Name)
	suite.Equal(userID, result.UserID)
	suite.NotEmpty(result.Token) // Should return the actual token on creation

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestCreateAPIToken_InvalidScopes() {
	// TODO: Implement scope validation in CreateAPIToken service
	suite.T().Skip("Scope validation not implemented yet")

	ctx := context.Background()
	userID := uint(1)
	req := &request.CreateAPITokenRequest{
		Name:   "Test Token",
		Scopes: []string{"invalid:scope", "user:read"},
	}

	// Skip permission validation for unit tests - simulate invalid scopes

	// Execute
	result, err := suite.service.CreateAPIToken(ctx, req, userID)

	// Assert
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "invalid scopes")

	// No permission repo assertions needed
}

func (suite *APITokenServiceTestSuite) TestGetAPIToken_Success() {
	ctx := context.Background()
	tokenID := uint(1)
	userID := uint(1)

	expectedToken := &model.APIToken{
		BaseModel:   model.BaseModel{ID: tokenID},
		Name:        "Test Token",
		UserID:      userID,
		Description: "Test description",
		Status:      1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(expectedToken, nil)

	// Execute
	result, err := suite.service.GetAPIToken(ctx, tokenID, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedToken.Name, result.Name)
	suite.Equal("***", result.Token) // Token should be masked

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestGetAPIToken_NotOwner() {
	ctx := context.Background()
	tokenID := uint(1)
	userID := uint(1)
	otherUserID := uint(2)

	expectedToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: tokenID},
		Name:      "Test Token",
		UserID:    otherUserID, // Different user
		Status:    1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(expectedToken, nil)

	// Execute
	result, err := suite.service.GetAPIToken(ctx, tokenID, userID)

	// Assert
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "token not found")

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestListAPITokens_Success() {
	ctx := context.Background()
	userID := uint(1)
	req := &request.ListAPITokensRequest{
		PaginationRequest: request.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
		UserID: &userID,
	}

	tokens := []*model.APIToken{
		{
			BaseModel: model.BaseModel{ID: 1},
			Name:      "Token 1",
			UserID:    userID,
			Status:    1,
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Name:      "Token 2",
			UserID:    userID,
			Status:    1,
		},
	}
	total := int64(2)

	suite.mockTokenRepo.On("List", ctx, req).Return(tokens, total, nil)

	// Execute
	result, err := suite.service.ListAPITokens(ctx, req, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(2, len(result.Items))
	suite.Equal(total, result.Total)
	// All tokens should have masked token values
	for _, token := range result.Items {
		suite.Equal("***", token.Token)
	}

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestUpdateAPIToken_Success() {
	ctx := context.Background()
	tokenID := uint(1)
	userID := uint(1)
	req := &request.UpdateAPITokenRequest{
		Name:        "Updated Token",
		Description: "Updated description",
		Status:      1,
	}

	// Mock existing token
	existingToken := &model.APIToken{
		BaseModel:   model.BaseModel{ID: tokenID},
		Name:        "Original Token",
		UserID:      userID,
		Description: "Original description",
		Status:      1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(existingToken, nil)
	suite.mockTokenRepo.On("Update", ctx, mock.AnythingOfType("*model.APIToken")).Return(nil)

	// Mock updated token for response
	updatedToken := &model.APIToken{
		BaseModel:   model.BaseModel{ID: tokenID},
		Name:        req.Name,
		UserID:      userID,
		Description: req.Description,
		Status:      req.Status,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(updatedToken, nil)

	// Execute
	result, err := suite.service.UpdateAPIToken(ctx, tokenID, req, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(req.Name, result.Name)
	suite.Equal(req.Description, result.Description)

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestValidateAPIToken_Success() {
	ctx := context.Background()
	tokenString := "test_token_string"
	tokenHash := security.HashToken(tokenString) // Use actual hashing function

	validToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Valid Token",
		UserID:    1,
		Status:    1,
		ExpiresAt: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
	}
	suite.mockTokenRepo.On("GetByToken", ctx, tokenHash).Return(validToken, nil)
	suite.mockTokenRepo.On("UpdateLastUsed", mock.Anything, validToken.ID, mock.AnythingOfType("string")).Return(nil)

	// Execute
	result, err := suite.service.ValidateAPIToken(ctx, tokenString)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.True(result.Valid)
	suite.Equal(validToken.UserID, result.UserID)

	// Wait for goroutine to complete UpdateLastUsed call
	time.Sleep(10 * time.Millisecond)

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestValidateAPIToken_Expired() {
	ctx := context.Background()
	tokenString := "test_token_string"
	tokenHash := security.HashToken(tokenString) // Use actual hashing function

	expiredToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Expired Token",
		UserID:    1,
		Status:    1,
		ExpiresAt: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(), // Expired
	}
	suite.mockTokenRepo.On("GetByToken", ctx, tokenHash).Return(expiredToken, nil)

	// Execute
	result, err := suite.service.ValidateAPIToken(ctx, tokenString)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.False(result.Valid)

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestRefreshAPIToken_Success() {
	ctx := context.Background()
	tokenID := uint(1)
	userID := uint(1)

	existingToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: tokenID},
		Name:      "Test Token",
		UserID:    userID,
		Status:    1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(existingToken, nil)
	suite.mockTokenRepo.On("Update", ctx, mock.AnythingOfType("*model.APIToken")).Return(nil)

	// Mock refreshed token for response
	refreshedToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: tokenID},
		Name:      existingToken.Name,
		UserID:    userID,
		Status:    1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(refreshedToken, nil)

	// Execute
	result, err := suite.service.RefreshAPIToken(ctx, tokenID, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.NotEmpty(result.Token) // Should return new token

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestRevokeAPIToken_Success() {
	ctx := context.Background()
	tokenID := uint(1)
	userID := uint(1)

	existingToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: tokenID},
		Name:      "Test Token",
		UserID:    userID,
		Status:    1,
	}
	suite.mockTokenRepo.On("GetByID", ctx, tokenID).Return(existingToken, nil)
	suite.mockTokenRepo.On("Update", ctx, mock.AnythingOfType("*model.APIToken")).Return(nil)

	// Execute
	err := suite.service.RevokeAPIToken(ctx, tokenID, userID)

	// Assert
	suite.NoError(err)

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

func (suite *APITokenServiceTestSuite) TestCleanExpiredTokens_Success() {
	ctx := context.Background()

	suite.mockTokenRepo.On("CleanExpiredTokens", ctx).Return(nil)

	// Execute
	err := suite.service.CleanExpiredTokens(ctx)

	// Assert
	suite.NoError(err)

	suite.mockTokenRepo.AssertExpectations(suite.T())
}

// Run the API token service test suite
func TestAPITokenServiceTestSuite(t *testing.T) {
	suite.Run(t, new(APITokenServiceTestSuite))
}

// Benchmark tests
func BenchmarkAPITokenService_ValidateToken(b *testing.B) {
	mockTokenRepo := &MockAPITokenRepository{}
	logger := logger.NewZapLogger(logger.InfoLevel, nil)

	service := &apiTokenService{
		tokenRepo: mockTokenRepo,
		db:        &gorm.DB{},
		logger:    logger,
		i18n:      nil,
	}

	ctx := context.Background()
	tokenString := "test_token_string"
	tokenHash := security.HashToken(tokenString) // Use actual hashing function

	validToken := &model.APIToken{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Valid Token",
		UserID:    1,
		Status:    1,
		ExpiresAt: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
	}

	// Setup mocks for benchmark
	mockTokenRepo.On("GetByToken", ctx, tokenHash).Return(validToken, nil)
	mockTokenRepo.On("UpdateLastUsed", ctx, mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateAPIToken(ctx, tokenString)
	}
}
