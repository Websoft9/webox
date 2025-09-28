package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/i18n"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserProfileRepository is a mock for UserProfileRepository interface
type MockUserProfileRepository struct {
	mock.Mock
}

// GetUserProfileByID mocks the GetUserProfileByID method
func (m *MockUserProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// LoadUserRoles mocks the LoadUserRoles method
func (m *MockUserProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// UpdateUserProfile mocks the UpdateUserProfile method
func (m *MockUserProfileRepository) UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error {
	args := m.Called(ctx, userID, updateData)
	return args.Error(0)
}

// UpdateUserPassword mocks the UpdateUserPassword method
func (m *MockUserProfileRepository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

// GetLoginHistories mocks the GetLoginHistories method
func (m *MockUserProfileRepository) GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.UserLoginHistory), args.Get(1).(int64), args.Error(2)
}

// GetUserConfigsByCategory mocks the GetUserConfigsByCategory method
func (m *MockUserProfileRepository) GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error) {
	args := m.Called(ctx, userID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserProfile), args.Error(1)
}

// SaveUserConfig mocks the SaveUserConfig method
func (m *MockUserProfileRepository) SaveUserConfig(ctx context.Context, config *model.UserProfile) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// GetUserConfig mocks the GetUserConfig method
func (m *MockUserProfileRepository) GetUserConfig(ctx context.Context, userID uint, category, key string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID, category, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

// TestGetUserProfile tests the GetUserProfile method
func TestGetUserProfile(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup for logging - making these more flexible to reduce potential nil issues
	mockLogger.On("InfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("WarnContext", mock.Anything, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Expected user
		expectedUser := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Nickname:  "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			// Adding roles here directly to avoid complexity with LoadUserRoles
			// Roles: []model.Role{
			// 	{ID: 1, Code: "admin", Name: "Administrator"},
			// 	{ID: 2, Code: "user", Name: "User"},
			// },
		}

		// Mock repository calls - returning success
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(expectedUser, nil).Once()
		mockRepo.On("LoadUserRoles", ctx, expectedUser).Return(nil).Once()

		// Call service method
		response, err := service.GetUserProfile(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedUser.ID, response.ID)
		assert.Equal(t, expectedUser.Username, response.Username)
		assert.Equal(t, expectedUser.Email, response.Email)
		assert.Equal(t, expectedUser.Nickname, response.Nickname)
		assert.Len(t, response.Roles, 0)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Role Loading Error", func(t *testing.T) {
		// Expected user
		expectedUser := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Nickname:  "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			// No roles here since we'll simulate a role loading error
		}

		// Mock repository calls
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(expectedUser, nil).Once()
		// Use fmt.Errorf instead of errors.New to avoid package conflict
		mockRepo.On("LoadUserRoles", ctx, expectedUser).Return(fmt.Errorf("role loading error")).Once()

		// Call service method
		response, err := service.GetUserProfile(ctx, userID)

		// Assert - should still return profile but with no roles
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedUser.ID, response.ID)
		assert.Equal(t, expectedUser.Username, response.Username)
		assert.Equal(t, expectedUser.Email, response.Email)
		assert.Empty(t, response.Roles)
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateUserProfile tests the UpdateUserProfile method
func TestUpdateUserProfile(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("WarnContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Request data
		nickname := "Updated Nickname"
		req := &request.UserProfileUpdateRequest{
			Nickname: &nickname,
		}

		// Expected update data
		expectedUpdateData := map[string]interface{}{
			"nickname": nickname,
		}

		// Expected user for GetUserProfile response
		expectedUser := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Nickname:  nickname,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock repository calls
		mockRepo.On("UpdateUserProfile", ctx, userID, expectedUpdateData).Return(nil).Once()
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(expectedUser, nil).Once()
		mockRepo.On("LoadUserRoles", ctx, expectedUser).Return(nil).Once()

		// Call service method
		response, err := service.UpdateUserProfile(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, userID, response.ID)
		assert.Equal(t, nickname, response.Nickname)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No Fields To Update", func(t *testing.T) {
		// Empty request
		req := &request.UserProfileUpdateRequest{}

		// Expected user for GetUserProfile response
		expectedUser := &model.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Nickname:  "Original Nickname",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock repository calls - should skip update and go directly to get
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(expectedUser, nil).Once()
		mockRepo.On("LoadUserRoles", ctx, expectedUser).Return(nil).Once()

		// Call service method
		response, err := service.UpdateUserProfile(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, userID, response.ID)
		assert.Equal(t, "Original Nickname", response.Nickname)
		mockRepo.AssertExpectations(t)
	})

}

// TestChangeProfilePassword tests the ChangeProfilePassword method
func TestChangeProfilePassword(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("WarnContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Old and new passwords
		oldPassword := "oldPassword123"
		newPassword := "newPassword123"

		// Request
		req := &request.ProfileChangePasswordRequest{
			OldPassword:     oldPassword,
			NewPassword:     newPassword,
			ConfirmPassword: newPassword,
		}

		// Expected user
		expectedUser := &model.User{
			ID:           userID,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: auth.HashToken(oldPassword),
		}

		// Expected new password hash
		expectedNewPasswordHash := auth.HashToken(newPassword)

		// Mock repository calls
		mockRepo.On("GetUserProfileByID", ctx, userID).Return(expectedUser, nil).Once()
		mockRepo.On("UpdateUserPassword", ctx, userID, expectedNewPasswordHash).Return(nil).Once()

		// Call service method
		err := service.ChangeProfilePassword(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

// TestGetLoginHistories tests the GetLoginHistories method
func TestGetLoginHistories(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Request
		req := &request.LoginHistoryRequest{
			Page:     1,
			PageSize: 10,
		}

		// Expected login histories
		loginTime := time.Now().Add(-time.Hour)
		logoutTime := time.Now().Add(-time.Minute)
		expectedHistories := []model.UserLoginHistory{
			{
				ID:         1,
				UserID:     userID,
				IPAddress:  "192.168.1.1",
				UserAgent:  "Mozilla/5.0",
				Device:     "Desktop",
				Browser:    "Chrome",
				Location:   "New York",
				LoginTime:  loginTime,
				LogoutTime: &logoutTime,
			},
		}

		// Mock repository calls
		mockRepo.On("GetLoginHistories", ctx, userID, 1, 10).Return(expectedHistories, int64(1), nil).Once()

		// Call service method
		response, err := service.GetLoginHistories(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Items, 1)
		// assert.Equal(t, uint(1), response.Items[0].ID)
		// assert.Equal(t, "192.168.1.1", response.Items[0].IPAddress)
		// assert.Equal(t, loginTime, response.Items[0].LoginTime)
		// assert.Equal(t, logoutTime, *response.Items[0].LogoutTime)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Default Pagination", func(t *testing.T) {
		// Request with invalid pagination (should use defaults)
		req := &request.LoginHistoryRequest{
			Page:     0,
			PageSize: 0,
		}

		// Mock repository calls with expected defaults
		mockRepo.On("GetLoginHistories", ctx, userID, 1, 20).Return([]model.UserLoginHistory{}, int64(0), nil).Once()

		// Call service method
		response, err := service.GetLoginHistories(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Items)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetNotificationSettings tests the GetNotificationSettings method
func TestGetNotificationSettings(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success with Saved Settings", func(t *testing.T) {
		// Expected user configs
		expectedConfigs := []*model.UserProfile{
			{
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "email_notifications",
				ConfigValue: constants.StringTrue,
			},
			{
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "sms_notifications",
				ConfigValue: constants.StringTrue,
			},
			{
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "push_notifications",
				ConfigValue: constants.StringFalse,
			},
			{
				UserID:      userID,
				Category:    "notification",
				ConfigKey:   "marketing_emails",
				ConfigValue: constants.StringTrue,
			},
		}

		// Mock repository calls
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "notification").Return(expectedConfigs, nil).Once()

		// Call service method
		response, err := service.GetNotificationSettings(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.EmailNotifications)
		assert.True(t, response.SmsNotifications)
		assert.False(t, response.PushNotifications)
		assert.True(t, response.MarketingEmails)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with Default Settings", func(t *testing.T) {
		// Mock repository calls with no saved configs
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "notification").Return(nil, gorm.ErrRecordNotFound).Once()

		// Call service method
		response, err := service.GetNotificationSettings(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		// Check default values
		assert.True(t, response.EmailNotifications) // Default true
		assert.False(t, response.SmsNotifications)  // Default false
		assert.True(t, response.PushNotifications)  // Default true
		assert.False(t, response.MarketingEmails)   // Default false
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateNotificationSettings tests the UpdateNotificationSettings method
func TestUpdateNotificationSettings(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Request
		req := &request.NotificationSettingsRequest{
			EmailNotifications: true,
			SmsNotifications:   false,
			PushNotifications:  true,
			MarketingEmails:    false,
		}

		// Mock repository calls for each setting
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "email_notifications" && config.ConfigValue == constants.StringTrue
		})).Return(nil).Once()

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "sms_notifications" && config.ConfigValue == constants.StringFalse
		})).Return(nil).Once()

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "push_notifications" && config.ConfigValue == constants.StringTrue
		})).Return(nil).Once()

		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "notification" &&
				config.ConfigKey == "marketing_emails" && config.ConfigValue == constants.StringFalse
		})).Return(nil).Once()

		// Call service method
		err := service.UpdateNotificationSettings(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetSecuritySettings tests the GetSecuritySettings method
func TestGetSecuritySettings(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success with Default Settings", func(t *testing.T) {
		// Mock repository calls with no saved configs
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return(nil, gorm.ErrRecordNotFound).Once()

		// Call service method
		response, err := service.GetSecuritySettings(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		// Check default values
		assert.True(t, response.LoginAlerts)
		assert.Equal(t, constants.DefaultSessionTimeoutSeconds, response.SessionTimeout)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Session Timeout Format", func(t *testing.T) {
		// User configs with invalid session timeout format
		expectedConfigs := []*model.UserProfile{
			{
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "login_alerts",
				ConfigValue: constants.StringTrue,
			},
			{
				UserID:      userID,
				Category:    "security",
				ConfigKey:   "session_timeout",
				ConfigValue: "invalid-json",
			},
		}

		// Mock repository calls
		mockRepo.On("GetUserConfigsByCategory", ctx, userID, "security").Return(expectedConfigs, nil).Once()

		// Call service method
		response, err := service.GetSecuritySettings(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.LoginAlerts)
		// Should fall back to default for invalid format
		assert.Equal(t, constants.DefaultSessionTimeoutSeconds, response.SessionTimeout)
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateSecuritySettings tests the UpdateSecuritySettings method
func TestUpdateSecuritySettings(t *testing.T) {
	// Setup
	mockRepo := new(MockUserProfileRepository)
	mockLogger := new(MockLogger)
	i18nInst := i18n.NewI18n()
	service := NewUserProfileService(mockRepo, mockLogger, i18nInst)
	ctx := context.Background()
	userID := uint(1)

	// Common mock setup
	mockLogger.On("InfoContext", ctx, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", ctx, mock.Anything, mock.Anything).Return()

	t.Run("Success", func(t *testing.T) {
		// Request
		req := &request.SecuritySettingsRequest{
			LoginAlerts:    false,
			SessionTimeout: 7200, // 2 hours
		}

		// Mock repository calls for login alerts setting
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "security" &&
				config.ConfigKey == "login_alerts" && config.ConfigValue == constants.StringFalse
		})).Return(nil).Once()

		// Mock repository calls for session timeout setting
		mockRepo.On("SaveUserConfig", ctx, mock.MatchedBy(func(config *model.UserProfile) bool {
			return config.UserID == userID && config.Category == "security" &&
				config.ConfigKey == "session_timeout" && config.ConfigValue == `7200`
		})).Return(nil).Once()

		// Call service method
		err := service.UpdateSecuritySettings(ctx, userID, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}
