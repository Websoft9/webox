package service

import (
	"api-service/internal/model"
	"api-service/pkg/logger"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// MockAPITokenRepository API token repository mock
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

func (m *MockAPITokenRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.APIToken), args.Error(1)
}

func (m *MockAPITokenRepository) GetActiveTokenByUserID(ctx context.Context, userID uint) (*model.APIToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APIToken), args.Error(1)
}

func (m *MockAPITokenRepository) BatchDelete(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
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
