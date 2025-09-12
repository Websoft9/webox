package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	serviceiface "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/internal/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

// MockUserProfileRepository 是一个模拟的用户个人资料仓储
type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserProfileRepository) UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error {
	args := m.Called(ctx, userID, updateData)
	return args.Error(0)
}

func (m *MockUserProfileRepository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

func (m *MockUserProfileRepository) GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).([]model.UserLoginHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserProfileRepository) GetUserConfig(ctx context.Context, userID uint, category, key string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID, category, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error) {
	args := m.Called(ctx, userID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) SaveUserConfig(ctx context.Context, config *model.UserProfile) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// 创建测试用户数据
func createTestUser(id uint) *model.User {
	now := time.Now()
	return &model.User{
		ID:          id,
		Username:    "testuser",
		Email:       "test@example.com",
		Nickname:    "Test User",
		Avatar:      "https://example.com/avatar.jpg",
		Phone:       "13800138000",
		Gender:      1,
		Signature:   "This is a test signature",
		Status:      1,
		LastLoginAt: &now,
		LastLoginIP: "192.168.1.1",
		Timezone:    "Asia/Shanghai",
		Language:    "zh-CN",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// 创建一个简单的日志实现
type mockLogger struct{}

func (l *mockLogger) InfoContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (l *mockLogger) WarnContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (l *mockLogger) ErrorContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (l *mockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (l *mockLogger) FatalContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (l *mockLogger) Info(msg string, fields ...logger.Field)                              {}
func (l *mockLogger) Warn(msg string, fields ...logger.Field)                              {}
func (l *mockLogger) Error(msg string, fields ...logger.Field)                             {}
func (l *mockLogger) Debug(msg string, fields ...logger.Field)                             {}
func (l *mockLogger) Fatal(msg string, fields ...logger.Field)                             {}

// 设置测试辅助函数
func setupUserProfileTest() (serviceiface.UserProfileService, *MockUserProfileRepository, *i18n.I18n) {
	mockRepo := new(MockUserProfileRepository)
	mockLogger := &struct {
		logger.Logger
	}{}
	i18nInstance := &i18n.I18n{}

	// 根据 NewUserProfileService 的实际签名调整
	// 假设它返回了 serviceiface.UserProfileService
	svc := service.NewUserProfileService(mockRepo, mockLogger, i18nInstance)

	return svc, mockRepo, i18nInstance
}

// 测试 GetUserProfile 方法
func TestGetUserProfile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟数据
		user := createTestUser(userID)

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("LoadUserRoles", ctx, user).Return(nil)

		// 执行测试
		result, err := svc.GetUserProfile(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Nickname, result.Nickname)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		// 执行测试
		result, err := svc.GetUserProfile(ctx, userID)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeRecordNotFound, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("LoadUserRolesFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟数据
		user := createTestUser(userID)

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("LoadUserRoles", ctx, user).Return(errors.NewAppError(errors.CodeInternalError, "Failed to load roles"))

		// 执行测试
		result, err := svc.GetUserProfile(ctx, userID)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 UpdateUserProfile 方法
func TestUpdateUserProfile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		nickname := "New Nickname"
		req := &request.UserProfileUpdateRequest{
			Nickname: &nickname,
		}

		// 模拟数据
		user := createTestUser(userID)
		user.Nickname = nickname

		// 设置期望 - 更新及获取更新后的配置文件
		mockRepo.On("UpdateUserProfile", ctx, userID, mock.MatchedBy(func(data map[string]interface{}) bool {
			return data["nickname"] == nickname
		})).Return(nil)
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("LoadUserRoles", ctx, user).Return(nil)

		// 执行测试
		result, err := svc.UpdateUserProfile(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, nickname, result.Nickname)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NoFieldsToUpdate", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据 - 没有要更新的字段
		req := &request.UserProfileUpdateRequest{}

		// 模拟数据
		user := createTestUser(userID)

		// 设置期望 - 直接获取当前配置文件
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("LoadUserRoles", ctx, user).Return(nil)

		// 执行测试
		result, err := svc.UpdateUserProfile(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Username, result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		nickname := "New Nickname"
		req := &request.UserProfileUpdateRequest{
			Nickname: &nickname,
		}

		// 设置期望 - 更新失败
		mockRepo.On("UpdateUserProfile", ctx, userID, mock.Anything).Return(gorm.ErrRecordNotFound)

		// 执行测试
		result, err := svc.UpdateUserProfile(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeRecordNotFound, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetProfileAfterUpdateFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		nickname := "New Nickname"
		req := &request.UserProfileUpdateRequest{
			Nickname: &nickname,
		}

		// 设置期望
		mockRepo.On("UpdateUserProfile", ctx, userID, mock.Anything).Return(nil)
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		// 执行测试
		result, err := svc.UpdateUserProfile(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeRecordNotFound, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 ChangeProfilePassword 方法
func TestChangeProfilePassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     "old_password",
			NewPassword:     "new_password",
			ConfirmPassword: "new_password",
		}

		// 模拟数据
		now := time.Now()
		user := &model.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a", // SHA-256 of "old_password"
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("UpdateUserPassword", ctx, userID, mock.MatchedBy(func(hash string) bool {
			// 这里假设密码哈希算法是SHA-256
			return hash != "" && hash != user.PasswordHash
		})).Return(nil)

		// 执行测试
		err := svc.ChangeProfilePassword(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     "old_password",
			NewPassword:     "new_password",
			ConfirmPassword: "new_password",
		}

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		// 执行测试
		err := svc.ChangeProfilePassword(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeRecordNotFound, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WrongOldPassword", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     "wrong_password",
			NewPassword:     "new_password",
			ConfirmPassword: "new_password",
		}

		// 模拟数据
		now := time.Now()
		user := &model.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a", // SHA-256 of "old_password"
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)

		// 执行测试
		err := svc.ChangeProfilePassword(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInvalidCredentials, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("PasswordMismatch", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     "old_password",
			NewPassword:     "new_password",
			ConfirmPassword: "different_password",
		}

		// 模拟数据
		now := time.Now()
		user := &model.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a", // SHA-256 of "old_password"
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)

		// 执行测试
		err := svc.ChangeProfilePassword(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeValidationFailed, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePasswordFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     "old_password",
			NewPassword:     "new_password",
			ConfirmPassword: "new_password",
		}

		// 模拟数据
		now := time.Now()
		user := &model.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "e7cf3ef4f17c3999a94f2c6f612e8a888e5b1026878e4e19398b23bd38ec221a", // SHA-256 of "old_password"
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 设置期望
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(user, nil)
		mockRepo.On("UpdateUserPassword", ctx, userID, mock.Anything).Return(gorm.ErrInvalidDB)

		// 执行测试
		err := svc.ChangeProfilePassword(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 GetLoginHistories 方法
func TestGetLoginHistories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.LoginHistoryRequest{
			Page:     1,
			PageSize: 10,
		}

		// 模拟数据
		now := time.Now()
		records := []model.UserLoginHistory{
			{
				ID:        1,
				UserID:    userID,
				IPAddress: "192.168.1.1",
				UserAgent: "Mozilla/5.0",
				Device:    "PC",
				Browser:   "Chrome",
				Location:  "Beijing",
				LoginTime: now,
			},
		}
		total := int64(1)

		// 设置期望
		mockRepo.On("GetLoginHistories", ctx, userID, 1, 10).Return(records, total, nil)

		// 执行测试
		result, err := svc.GetLoginHistories(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, uint(1), result.Items[0].ID)
		assert.Equal(t, "192.168.1.1", result.Items[0].IPAddress)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DefaultPageValues", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据 - 无页码参数，应使用默认值
		req := &request.LoginHistoryRequest{}

		// 模拟数据
		records := []model.UserLoginHistory{}
		total := int64(0)

		// 设置期望 - 应使用默认值page=1, pageSize=20
		mockRepo.On("GetLoginHistories", ctx, userID, 1, 20).Return(records, total, nil)

		// 执行测试
		result, err := svc.GetLoginHistories(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Items)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.LoginHistoryRequest{
			Page:     1,
			PageSize: 10,
		}

		// 设置期望 - 仓库返回错误
		mockRepo.On("GetLoginHistories", ctx, userID, 1, 10).Return(
			[]model.UserLoginHistory{},
			int64(0),
			errors.NewAppError(errors.CodeInternalError, "Database error"),
		)

		// 执行测试
		result, err := svc.GetLoginHistories(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 GetNotificationSettings 方法
func TestGetNotificationSettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟数据
		configs := []*model.UserProfile{
			{
				ID:          1,
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "email_notifications",
				ConfigValue: constants.StringTrue,
			},
			{
				ID:          2,
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "sms_notifications",
				ConfigValue: constants.StringTrue,
			},
		}

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "notification").Return(configs, nil)

		// 执行测试
		result, err := svc.GetNotificationSettings(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.EmailNotifications)
		assert.True(t, result.SmsNotifications)
		assert.True(t, result.PushNotifications) // 默认值
		assert.False(t, result.MarketingEmails)  // 默认值
		mockRepo.AssertExpectations(t)
	})

	t.Run("NoConfigsFound", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "notification").Return([]*model.UserProfile{}, nil)

		// 执行测试
		result, err := svc.GetNotificationSettings(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.EmailNotifications) // 默认值
		assert.False(t, result.SmsNotifications)  // 默认值
		assert.True(t, result.PushNotifications)  // 默认值
		assert.False(t, result.MarketingEmails)   // 默认值
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "notification").Return(
			nil,
			errors.NewAppError(errors.CodeInternalError, "Database error"),
		)

		// 执行测试
		result, err := svc.GetNotificationSettings(ctx, userID)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 UpdateNotificationSettings 方法
func TestUpdateNotificationSettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.NotificationSettingsRequest{
			EmailNotifications: true,
			SmsNotifications:   true,
			PushNotifications:  false,
			MarketingEmails:    true,
		}

		// 设置期望
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "email_notifications" && config.ConfigValue == constants.StringTrue
		})).Return(nil)

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "sms_notifications" && config.ConfigValue == constants.StringTrue
		})).Return(nil)

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "push_notifications" && config.ConfigValue == constants.StringFalse
		})).Return(nil)

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "marketing_emails" && config.ConfigValue == constants.StringTrue
		})).Return(nil)

		// 执行测试
		err := svc.UpdateNotificationSettings(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveConfigFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.NotificationSettingsRequest{
			EmailNotifications: true,
		}

		// 设置期望 - 保存配置失败
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification"
		})).Return(gorm.ErrInvalidDB)

		// 执行测试
		err := svc.UpdateNotificationSettings(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}

// 测试 GetSecuritySettings 方法
func TestGetSecuritySettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟数据
		timeoutValue, _ := json.Marshal(3600)
		configs := []*model.UserProfile{
			{
				ID:          1,
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "login_alerts",
				ConfigValue: constants.StringTrue,
			},
			{
				ID:          2,
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "session_timeout",
				ConfigValue: string(timeoutValue),
			},
		}

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return(configs, nil)

		// 执行测试
		result, err := svc.GetSecuritySettings(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.LoginAlerts)
		assert.Equal(t, 3600, result.SessionTimeout)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NoConfigsFound", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return([]*model.UserProfile{}, nil)

		// 执行测试
		result, err := svc.GetSecuritySettings(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.LoginAlerts)                                             // 默认值
		assert.Equal(t, constants.DefaultSessionTimeoutSeconds, result.SessionTimeout) // 默认值
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return(
			nil,
			errors.NewAppError(errors.CodeInternalError, "Database error"),
		)

		// 执行测试
		result, err := svc.GetSecuritySettings(ctx, userID)

		// 验证结果
		require.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidSessionTimeoutValue", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟数据 - 会话超时值无效
		configs := []*model.UserProfile{
			{
				ID:          1,
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "login_alerts",
				ConfigValue: constants.StringTrue,
			},
			{
				ID:          2,
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "session_timeout",
				ConfigValue: "invalid-json",
			},
		}

		// 设置期望
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return(configs, nil)

		// 执行测试
		result, err := svc.GetSecuritySettings(ctx, userID)

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.LoginAlerts)
		assert.Equal(t, constants.DefaultSessionTimeoutSeconds, result.SessionTimeout) // 使用默认值
		mockRepo.AssertExpectations(t)
	})
}

// 测试 UpdateSecuritySettings 方法
func TestUpdateSecuritySettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.SecuritySettingsRequest{
			LoginAlerts:    false,
			SessionTimeout: 7200,
		}

		// 设置期望
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "security" &&
				config.ConfigKey == "login_alerts" && config.ConfigValue == constants.StringFalse
		})).Return(nil)

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "security" &&
				config.ConfigKey == "session_timeout" && config.ConfigValue != ""
		})).Return(nil)

		// 执行测试
		err := svc.UpdateSecuritySettings(ctx, userID, req)

		// 验证结果
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveConfigFailed", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据
		req := &request.SecuritySettingsRequest{
			LoginAlerts:    false,
			SessionTimeout: 7200,
		}

		// 设置期望 - 保存配置失败
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "security"
		})).Return(gorm.ErrInvalidDB)

		// 执行测试
		err := svc.UpdateSecuritySettings(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeInternalError, appErr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SessionTimeoutOutOfRange", func(t *testing.T) {
		// 设置测试环境
		svc, mockRepo, _ := setupUserProfileTest()
		ctx := context.Background()
		userID := uint(1)

		// 模拟请求数据 - 会话超时值超出范围
		req := &request.SecuritySettingsRequest{
			LoginAlerts:    true,
			SessionTimeout: 100000, // 假设这个值超出了允许范围
		}

		// 执行测试
		err := svc.UpdateSecuritySettings(ctx, userID, req)

		// 验证结果
		require.Error(t, err)
		appErr, ok := err.(*errors.AppError)
		assert.True(t, ok)
		assert.Equal(t, errors.CodeValidationFailed, appErr.Code)
		mockRepo.AssertExpectations(t)
	})
}
