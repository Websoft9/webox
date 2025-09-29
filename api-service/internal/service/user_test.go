package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/logger"
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDWithRelations(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	args := m.Called(ctx, usernameOrEmail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error) {
	args := m.Called(ctx, offset, limit, filters)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ListWithRelations(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error) {
	args := m.Called(ctx, offset, limit, filters)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	args := m.Called(ctx, keyword, offset, limit)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	args := m.Called(ctx, username, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error) {
	args := m.Called(ctx, email, excludeID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetUserStats(ctx context.Context, userID uint) (*repository.UserStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.UserStats), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

// Role related mock methods
func (m *MockUserRepository) CreateUserRole(ctx context.Context, userRole *model.UserRole) error {
	args := m.Called(ctx, userRole)
	return args.Error(0)
}

func (m *MockUserRepository) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockUserRepository) DeleteUserRole(ctx context.Context, userID uint, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// Test setup
func setupUserServiceTestFixed() (*userService, *MockUserRepository) {
	mockUserRepo := &MockUserRepository{}
	mockLogger := logger.NewZapLogger(logger.InfoLevel, nil)

	service := &userService{
		userRepo: mockUserRepo,
		logger:   mockLogger,
	}

	return service, mockUserRepo
}

// Helper function to create test user
func createTestUserFixed() *model.User {
	now := time.Now()
	passwordHash := auth.HashToken("password123")
	return &model.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Nickname:     "Test User",
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastLoginAt:  &now,
	}
}

// ===== ListUsers Tests =====
func TestUserService_ListUsers_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	req := &request.UserListRequest{}
	req.Page = 1
	req.PageSize = 10

	users := []*model.User{
		{ID: 1, Username: "user1", Email: "user1@example.com", Status: UserStatusActive},
		{ID: 2, Username: "user2", Email: "user2@example.com", Status: UserStatusActive},
	}
	total := int64(2)

	mockRepo.On("ListWithRelations", ctx, 0, 10, mock.AnythingOfType("map[string]interface {}")).Return(users, total, nil)

	result, err := service.ListUsers(ctx, req) // 注意：这里移除了 totalCount 返回值

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)

	// 类型断言：将 interface{} 转换为具体的切片类型
	items, ok := result.Items.([]response.UserResponse)
	assert.True(t, ok, "Items should be of type []response.UserResponse")
	assert.Equal(t, len(users), len(items))

	mockRepo.AssertExpectations(t)
}

// ===== GetUser Tests =====

func TestUserService_GetUser_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)

	mockRepo.On("GetByIDWithRelations", ctx, userID).Return(testUser, nil)

	result, err := service.GetUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.ID)
	assert.Equal(t, testUser.Username, result.Username)
	assert.Equal(t, testUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	userID := uint(999)

	mockRepo.On("GetByIDWithRelations", ctx, userID).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetUser(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== CreateUser Tests =====

func TestUserService_CreateUser_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	nickname := "New User"
	status := UserStatusActive
	req := &request.UserCreateRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Nickname: &nickname,
		Status:   &status,
		RoleIDs:  []uint{1, 2},
	}

	mockRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
	mockRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*model.User)
			user.ID = 1
		}).Return(nil)
	mockRepo.On("CreateUserRole", ctx, mock.AnythingOfType("*model.UserRole")).Return(nil).Twice()

	currentUserID := uint(100)
	result, err := service.CreateUser(ctx, currentUserID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Username, result.Username)
	assert.Equal(t, req.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

// ===== UpdateUser Tests =====

func TestUserService_UpdateUser_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)
	currentUserID := uint(100)
	newEmail := "updated@example.com"
	newNickname := "Updated User"
	newRoleIDs := []uint{1, 3}

	req := &request.UserUpdateRequest{
		Email:    &newEmail,
		Nickname: &newNickname,
		RoleIDs:  newRoleIDs,
	}

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)
	mockRepo.On("ExistsByEmailExcludeID", ctx, newEmail, userID).Return(false, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)
	mockRepo.On("GetRoleIDsByUserID", ctx, userID).Return([]uint{1, 2}, nil)
	mockRepo.On("CreateUserRole", ctx, mock.AnythingOfType("*model.UserRole")).Return(nil)
	mockRepo.On("DeleteUserRole", ctx, userID, uint(2)).Return(nil)

	result, err := service.UpdateUser(ctx, currentUserID, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

// ===== UpdateUserStatus Tests =====

func TestUserService_UpdateUserStatus_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)

	req := &request.UserUpdateStatusRequest{
		Status: UserStatusInactive,
	}

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	err := service.UpdateUserStatus(ctx, userID, req)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// ===== DeleteUser Tests =====

func TestUserService_DeleteUser_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := service.DeleteUser(ctx, userID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)
	expectedError := errors.New("database error")

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)
	mockRepo.On("Delete", ctx, userID).Return(expectedError)

	err := service.DeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

// ===== ChangePassword Tests =====

func TestUserService_ChangePassword_Success(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)

	req := &request.UserChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	}

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	err := service.ChangePassword(ctx, userID, req)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_ChangePassword_WrongOldPassword(t *testing.T) {
	service, mockRepo := setupUserServiceTestFixed()
	ctx := context.Background()

	testUser := createTestUserFixed()
	userID := uint(1)

	req := &request.UserChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	mockRepo.On("GetByID", ctx, userID).Return(testUser, nil)

	err := service.ChangePassword(ctx, userID, req)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
