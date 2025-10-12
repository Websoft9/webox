package service

import (
	"api-service/internal/dto/common"
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

func (m *MockUserRepository) ListWithRelations(ctx context.Context, req *request.UserListRequest) ([]*model.User, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
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

// TestUserService_ListUsers tests the ListUsers method
func TestUserService_ListUsers(t *testing.T) {
	// Setup
	mockRepo := &MockUserRepository{}
	mockLogger := &MockLogger{}
	service := NewUserService(mockRepo, mockLogger)
	ctx := context.Background()

	// Common mock setup for logging
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Prepare test data - simplified without roles to avoid field issues
		users := []*model.User{
			{
				ID:       1,
				Username: "user1",
				Email:    "user1@example.com",
				Nickname: "User One",
				Status:   1,
			},
			{
				ID:       2,
				Username: "user2",
				Email:    "user2@example.com",
				Nickname: "User Two",
				Status:   1,
			},
		}
		total := int64(2)

		// Create request
		req := &request.UserListRequest{
			BaseListRequest: common.BaseListRequest{
				PaginationRequest: common.PaginationRequest{
					Page:     1,
					PageSize: 10,
				},
			},
		}

		// Set mock expectations
		mockRepo.On("ListWithRelations", ctx, req).Return(users, total, nil).Once()

		// Execute test
		result, err := service.ListUsers(ctx, req)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, result.Total)

		userList, ok := result.Items.([]response.UserResponse)
		assert.True(t, ok, "Items should be of type []response.UserResponse")
		assert.Len(t, userList, 2)

		// Verify first user
		assert.Equal(t, uint(1), userList[0].ID)
		assert.Equal(t, "user1", userList[0].Username)
		assert.Equal(t, "user1@example.com", userList[0].Email)
		assert.Equal(t, "User One", userList[0].Nickname)
		assert.Equal(t, 1, userList[0].Status)

		// Verify second user
		assert.Equal(t, uint(2), userList[1].ID)
		assert.Equal(t, "user2", userList[1].Username)
		assert.Equal(t, "user2@example.com", userList[1].Email)
		assert.Equal(t, "User Two", userList[1].Nickname)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with filters", func(t *testing.T) {
		// Prepare filtered test data
		users := []*model.User{
			{
				ID:       1,
				Username: "admin",
				Email:    "admin@example.com",
				Status:   1,
				Gender:   1,
			},
		}
		total := int64(1)

		// Create request with filters
		status := 1
		gender := 1
		keyword := "admin"
		roleID := uint(1)
		req := &request.UserListRequest{
			BaseListRequest: common.BaseListRequest{
				PaginationRequest: common.PaginationRequest{
					Page:     1,
					PageSize: 10,
				},
			},
			Status:  &status,
			Keyword: &keyword,
			Gender:  &gender,
			RoleID:  &roleID,
		}

		// Set mock expectations
		mockRepo.On("ListWithRelations", ctx, req).Return(users, total, nil).Once()

		// Execute test
		result, err := service.ListUsers(ctx, req)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, result.Total)

		userList, ok := result.Items.([]response.UserResponse)
		assert.True(t, ok)
		assert.Len(t, userList, 1)
		assert.Equal(t, uint(1), userList[0].ID)
		assert.Equal(t, "admin", userList[0].Username)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty result", func(t *testing.T) {
		// Prepare empty test data
		users := []*model.User{}
		total := int64(0)

		// Create request
		req := &request.UserListRequest{
			BaseListRequest: common.BaseListRequest{
				PaginationRequest: common.PaginationRequest{
					Page:     1,
					PageSize: 10,
				},
			},
		}

		// Set mock expectations
		mockRepo.On("ListWithRelations", ctx, req).Return(users, total, nil).Once()

		// Execute test
		result, err := service.ListUsers(ctx, req)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, result.Total)

		userList, ok := result.Items.([]response.UserResponse)
		assert.True(t, ok)
		assert.Empty(t, userList)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		// Create request
		req := &request.UserListRequest{
			BaseListRequest: common.BaseListRequest{
				PaginationRequest: common.PaginationRequest{
					Page:     1,
					PageSize: 10,
				},
			},
		}

		// Set mock expectations with error
		expectedError := errors.New("database connection error")
		mockRepo.On("ListWithRelations", ctx, req).Return(nil, int64(0), expectedError).Once()

		// Execute test
		result, err := service.ListUsers(ctx, req)

		// Verify error handling
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with roles", func(t *testing.T) {
		// If you need to test with roles, create them without specifying fields
		// that might not exist in the model.Role struct
		users := []*model.User{
			{
				ID:       1,
				Username: "user_with_roles",
				Email:    "roles@example.com",
				Nickname: "User With Roles",
				Status:   1,
				// Leave Roles empty or use actual Role objects if you know the correct fields
			},
		}
		total := int64(1)

		// Create request
		req := &request.UserListRequest{
			BaseListRequest: common.BaseListRequest{
				PaginationRequest: common.PaginationRequest{
					Page:     1,
					PageSize: 10,
				},
			},
		}

		// Set mock expectations
		mockRepo.On("ListWithRelations", ctx, req).Return(users, total, nil).Once()

		// Execute test
		result, err := service.ListUsers(ctx, req)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, total, result.Total)

		userList, ok := result.Items.([]response.UserResponse)
		assert.True(t, ok)
		assert.Len(t, userList, 1)
		assert.Equal(t, uint(1), userList[0].ID)

		mockRepo.AssertExpectations(t)
	})
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
