package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserProfileRepository is a mock implementation of the UserProfileRepository interface
type MockUserProfileRepository struct {
	mock.Mock
}

// GetUserProfileByID mocks the repository method to get a user profile by ID
func (m *MockUserProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// LoadUserRoles mocks the repository method to load user roles
func (m *MockUserProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// UpdateUserProfile mocks the repository method to update a user profile
func (m *MockUserProfileRepository) UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error {
	args := m.Called(ctx, userID, updateData)
	return args.Error(0)
}

// UpdateUserPassword mocks the repository method to update a user password
func (m *MockUserProfileRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error {
	args := m.Called(ctx, userID, hashedPassword)
	return args.Error(0)
}

// GetLoginHistories mocks the repository method to get login histories
func (m *MockUserProfileRepository) GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.UserLoginHistory), args.Get(1).(int64), args.Error(2)
}

// GetUserConfigsByCategory mocks the repository method to get user configs by category
func (m *MockUserProfileRepository) GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error) {
	args := m.Called(ctx, userID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserProfile), args.Error(1)
}

// SaveUserConfig mocks the repository method to save user config
func (m *MockUserProfileRepository) SaveUserConfig(ctx context.Context, config *model.UserProfile) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// GetUserConfig mocks the repository method to get a single user config by key
func (m *MockUserProfileRepository) GetUserConfig(ctx context.Context, userID uint, category, configKey string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID, category, configKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

// MockLogger is a mock implementation of the logger.Logger interface
type UserProfileMockLogger struct {
	mock.Mock
}

// InfoContext mocks the InfoContext method
func (m *MockLogger) UserProfileInfoContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

// WarnContext mocks the WarnContext method
func (m *MockLogger) UserProfileWarnContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

// ErrorContext mocks the ErrorContext method
func (m *MockLogger) UserProfileErrorContext(ctx context.Context, msg string, fields ...logger.Field) {
	m.Called(ctx, msg, fields)
}

// Debug mocks the Debug method
func (m *MockLogger) UserProfileDebug(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

// Info mocks the Info method
func (m *MockLogger) UserProfileInfo(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

// Warn mocks the Warn method
func (m *MockLogger) UserProfileWarn(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

// Error mocks the Error method
func (m *MockLogger) UserProfileError(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

// Fatal mocks the Fatal method
func (m *MockLogger) UserProfileFatal(msg string, fields ...logger.Field) {
	m.Called(msg, fields)
}

// TestGetUserProfile tests the GetUserProfile method
func TestGetUserProfile(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		userID        uint
		mockUser      *model.User
		mockError     error
		roleError     error
		expectedError error
		expectedRoles int
	}{
		{
			name:   "success_with_roles",
			userID: 1,
			mockUser: &model.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				Nickname:  "Test User",
				Avatar:    "avatar.png",
				Phone:     "1234567890",
				Gender:    1,
				Signature: "Test Signature",
				Status:    1,
				LastLoginAt: func() *time.Time {
					t, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
					return &t
				}(),
				LastLoginIP: "127.0.0.1",
				Timezone:    "UTC",
				Language:    "en",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Roles: []model.Role{
					{Name: "admin", Description: "Administrator"},
				},
			},
			mockError:     nil,
			roleError:     nil,
			expectedError: nil,
			expectedRoles: 1,
		},
		{
			name:          "user_not_found",
			userID:        2,
			mockUser:      nil,
			mockError:     gorm.ErrRecordNotFound,
			roleError:     nil,
			expectedError: errors.NewAppError(errors.CodeRecordNotFound, "user_profile.not_found"),
			expectedRoles: 0,
		},
		{
			name:   "role_load_error",
			userID: 3,
			mockUser: &model.User{
				ID:       3,
				Username: "testuser3",
				Email:    "test3@example.com",
			},
			mockError:     nil,
			roleError:     fmt.Errorf("role load failed"),
			expectedError: nil, // Role loading error should not prevent profile retrieval
			expectedRoles: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			mockRepo.On("GetUserProfileByID", mock.Anything, tc.userID).Return(tc.mockUser, tc.mockError)
			if tc.mockUser != nil {
				mockRepo.On("LoadUserRoles", mock.Anything, tc.mockUser).Return(tc.roleError)
			}

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("InfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("WarnContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("ErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			profile, err := service.GetUserProfile(context.Background(), tc.userID)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, tc.mockUser.ID, profile.ID)
				assert.Equal(t, tc.mockUser.Username, profile.Username)
				assert.Equal(t, tc.mockUser.Email, profile.Email)
				assert.Len(t, profile.Roles, tc.expectedRoles)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestUpdateUserProfile tests the UpdateUserProfile method
func TestUpdateUserProfile(t *testing.T) {
	// Helper function to create string pointer
	str := func(s string) *string {
		return &s
	}
	genderValue := 1
	// Test cases
	testCases := []struct {
		name           string
		userID         uint
		updateRequest  *request.UserProfileUpdateRequest
		updateError    error
		getUserError   error
		expectedError  error
		expectedFields map[string]interface{}
	}{
		{
			name:   "success_update_all_fields",
			userID: 1,
			updateRequest: &request.UserProfileUpdateRequest{
				Nickname:  str("New Nickname"),
				Avatar:    str("new_avatar.png"),
				Phone:     str("9876543210"),
				Gender:    &genderValue,
				Signature: str("New Signature"),
				Timezone:  str("GMT+8"),
				Language:  str("zh-CN"),
			},
			updateError:   nil,
			getUserError:  nil,
			expectedError: nil,
			expectedFields: map[string]interface{}{
				"nickname":  "New Nickname",
				"avatar":    "new_avatar.png",
				"phone":     "9876543210",
				"gender":    "female",
				"signature": "New Signature",
				"timezone":  "GMT+8",
				"language":  "zh-CN",
			},
		},
		{
			name:   "success_update_partial_fields",
			userID: 2,
			updateRequest: &request.UserProfileUpdateRequest{
				Nickname: str("New Nickname"),
				Language: str("zh-CN"),
			},
			updateError:   nil,
			getUserError:  nil,
			expectedError: nil,
			expectedFields: map[string]interface{}{
				"nickname": "New Nickname",
				"language": "zh-CN",
			},
		},
		{
			name:           "no_fields_to_update",
			userID:         3,
			updateRequest:  &request.UserProfileUpdateRequest{},
			updateError:    nil,
			getUserError:   nil,
			expectedError:  nil,
			expectedFields: map[string]interface{}{
				// Empty map as there are no fields to update
			},
		},
		{
			name:   "user_not_found",
			userID: 4,
			updateRequest: &request.UserProfileUpdateRequest{
				Nickname: str("New Nickname"),
			},
			updateError:   gorm.ErrRecordNotFound,
			getUserError:  nil,
			expectedError: errors.NewAppError(errors.CodeRecordNotFound, "user_profile.not_found"),
			expectedFields: map[string]interface{}{
				"nickname": "New Nickname",
			},
		},
		{
			name:   "other_update_error",
			userID: 5,
			updateRequest: &request.UserProfileUpdateRequest{
				Nickname: str("New Nickname"),
			},
			updateError:   fmt.Errorf("database error"),
			getUserError:  nil,
			expectedError: errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.update_failed"),
			expectedFields: map[string]interface{}{
				"nickname": "New Nickname",
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			if len(tc.expectedFields) > 0 {
				mockRepo.On("UpdateUserProfile", mock.Anything, tc.userID, mock.MatchedBy(func(data map[string]interface{}) bool {
					// Check if all expected fields are in the update data
					for k, v := range tc.expectedFields {
						if data[k] != v {
							return false
						}
					}
					return true
				})).Return(tc.updateError)
			}

			// If no update error and no fields to update, GetUserProfile should be called directly
			if tc.updateError == nil && len(tc.expectedFields) == 0 {
				mockRepo.On("GetUserProfileByID", mock.Anything, tc.userID).Return(&model.User{
					ID:       tc.userID,
					Username: "testuser",
				}, tc.getUserError)
				mockRepo.On("LoadUserRoles", mock.Anything, mock.Anything).Return(nil)
			} else if tc.updateError == nil {
				// If update is successful, GetUserProfile will be called to return updated profile
				mockRepo.On("GetUserProfileByID", mock.Anything, tc.userID).Return(&model.User{
					ID:       tc.userID,
					Username: "testuser",
				}, tc.getUserError)
				mockRepo.On("LoadUserRoles", mock.Anything, mock.Anything).Return(nil)
			}

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileWarnContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			profile, err := service.UpdateUserProfile(context.Background(), tc.userID, tc.updateRequest)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				// Compare error messages since error types might be different due to wrapping
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				// If no fields to update or successful update, profile should be returned
				if len(tc.expectedFields) == 0 || tc.updateError == nil {
					assert.NotNil(t, profile)
					assert.Equal(t, tc.userID, profile.ID)
				}
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestChangeProfilePassword tests the ChangeProfilePassword method
func TestChangeProfilePassword(t *testing.T) {
	// Test cases
	testCases := []struct {
		name             string
		userID           uint
		passwordRequest  *request.ProfileChangePasswordRequest
		mockUser         *model.User
		getUserError     error
		updateError      error
		expectedError    error
		expectedPassword string
	}{
		{
			name:   "success_change_password",
			userID: 1,
			passwordRequest: &request.ProfileChangePasswordRequest{
				OldPassword:     "oldpassword123",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			mockUser: &model.User{
				ID:           1,
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", // SHA256 of "123"
			},
			getUserError:     nil,
			updateError:      nil,
			expectedError:    nil,
			expectedPassword: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", // SHA256 of "123"
		},
		{
			name:   "user_not_found",
			userID: 2,
			passwordRequest: &request.ProfileChangePasswordRequest{
				OldPassword:     "oldpassword123",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			mockUser:         nil,
			getUserError:     gorm.ErrRecordNotFound,
			updateError:      nil,
			expectedError:    errors.NewAppError(errors.CodeRecordNotFound, "user_profile.not_found"),
			expectedPassword: "",
		},
		{
			name:   "invalid_old_password",
			userID: 3,
			passwordRequest: &request.ProfileChangePasswordRequest{
				OldPassword:     "wrongpassword",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			mockUser: &model.User{
				ID:           3,
				Username:     "testuser3",
				Email:        "test3@example.com",
				PasswordHash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", // SHA256 of "123"
			},
			getUserError:     nil,
			updateError:      nil,
			expectedError:    errors.NewAppError(errors.CodeInvalidCredentials, "user_profile.password_verification_failed"),
			expectedPassword: "",
		},
		{
			name:   "password_mismatch",
			userID: 4,
			passwordRequest: &request.ProfileChangePasswordRequest{
				OldPassword:     "oldpassword123",
				NewPassword:     "newpassword123",
				ConfirmPassword: "differentpassword",
			},
			mockUser: &model.User{
				ID:           4,
				Username:     "testuser4",
				Email:        "test4@example.com",
				PasswordHash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", // SHA256 of "123"
			},
			getUserError:     nil,
			updateError:      nil,
			expectedError:    errors.NewAppError(errors.CodeValidationFailed, "user_profile.password_mismatch"),
			expectedPassword: "",
		},
		{
			name:   "update_error",
			userID: 5,
			passwordRequest: &request.ProfileChangePasswordRequest{
				OldPassword:     "oldpassword123",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			mockUser: &model.User{
				ID:           5,
				Username:     "testuser5",
				Email:        "test5@example.com",
				PasswordHash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", // SHA256 of "123"
			},
			getUserError:     nil,
			updateError:      fmt.Errorf("database error"),
			expectedError:    errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.password_update_failed"),
			expectedPassword: "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", // SHA256 of "newpassword123"
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			mockRepo.On("GetUserProfileByID", mock.Anything, tc.userID).Return(tc.mockUser, tc.getUserError)

			// Only expect UpdateUserPassword call if validation passes
			if tc.mockUser != nil &&
				tc.passwordRequest.OldPassword == "oldpassword123" && // This is just a test value that we've chosen
				tc.passwordRequest.NewPassword == tc.passwordRequest.ConfirmPassword {

				mockRepo.On("UpdateUserPassword", mock.Anything, tc.userID, mock.Anything).Return(tc.updateError)
			}

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileWarnContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			err := service.ChangeProfilePassword(context.Background(), tc.userID, tc.passwordRequest)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestGetLoginHistories tests the GetLoginHistories method
func TestGetLoginHistories(t *testing.T) {
	// Sample login history data
	sampleHistories := []model.UserLoginHistory{
		{
			ID:         1,
			UserID:     1,
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			Device:     "Desktop",
			Browser:    "Chrome",
			Location:   "Shanghai, China",
			LoginTime:  time.Now().Add(-1 * time.Hour),
			LogoutTime: func() *time.Time { t := time.Now().Add(-30 * time.Minute); return &t }(),
		},
		{
			ID:         2,
			UserID:     1,
			IPAddress:  "192.168.1.2",
			UserAgent:  "Mozilla/5.0",
			Device:     "Mobile",
			Browser:    "Safari",
			Location:   "Beijing, China",
			LoginTime:  time.Now().Add(-2 * time.Hour),
			LogoutTime: func() *time.Time { t := time.Now().Add(-1 * time.Hour); return &t }(),
		},
	}

	// Test cases
	testCases := []struct {
		name           string
		userID         uint
		request        *request.LoginHistoryRequest
		histories      []model.UserLoginHistory
		totalCount     int64
		repoError      error
		expectedError  error
		expectedLength int
	}{
		{
			name:   "success_with_data",
			userID: 1,
			request: &request.LoginHistoryRequest{
				Page:     1,
				PageSize: 10,
			},
			histories:      sampleHistories,
			totalCount:     2,
			repoError:      nil,
			expectedError:  nil,
			expectedLength: 2,
		},
		{
			name:   "success_empty_result",
			userID: 2,
			request: &request.LoginHistoryRequest{
				Page:     1,
				PageSize: 10,
			},
			histories:      []model.UserLoginHistory{},
			totalCount:     0,
			repoError:      nil,
			expectedError:  nil,
			expectedLength: 0,
		},
		{
			name:   "repository_error",
			userID: 3,
			request: &request.LoginHistoryRequest{
				Page:     1,
				PageSize: 10,
			},
			histories:      nil,
			totalCount:     0,
			repoError:      fmt.Errorf("database error"),
			expectedError:  errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.login_history_get_failed"),
			expectedLength: 0,
		},
		{
			name:   "default_pagination",
			userID: 1,
			request: &request.LoginHistoryRequest{
				Page:     0,
				PageSize: 0,
			},
			histories:      sampleHistories,
			totalCount:     2,
			repoError:      nil,
			expectedError:  nil,
			expectedLength: 2,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Calculate expected page and page size
			expectedPage := tc.request.Page
			if expectedPage <= 0 {
				expectedPage = 1
			}
			expectedPageSize := tc.request.PageSize
			if expectedPageSize <= 0 {
				expectedPageSize = 20
			}

			// Set up mock expectations
			mockRepo.On("GetLoginHistories", mock.Anything, tc.userID, expectedPage, expectedPageSize).
				Return(tc.histories, tc.totalCount, tc.repoError)

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileWarnContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			result, err := service.GetLoginHistories(context.Background(), tc.userID, tc.request)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Items, tc.expectedLength)

				// Check item mapping if there are items
				if tc.expectedLength > 0 {
					for i, hist := range tc.histories {
						assert.Equal(t, hist.ID, result.Items[i].ID)
						assert.Equal(t, hist.IPAddress, result.Items[i].IPAddress)
						assert.Equal(t, hist.UserAgent, result.Items[i].UserAgent)
						assert.Equal(t, hist.Device, result.Items[i].Device)
						assert.Equal(t, hist.Browser, result.Items[i].Browser)
						assert.Equal(t, hist.Location, result.Items[i].Location)
						assert.Equal(t, hist.LoginTime, result.Items[i].LoginTime)
						assert.Equal(t, hist.LogoutTime, result.Items[i].LogoutTime)
					}
				}
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestGetNotificationSettings tests the GetNotificationSettings method
func TestGetNotificationSettings(t *testing.T) {
	// Test cases
	testCases := []struct {
		name         string
		userID       uint
		configs      []model.UserProfile
		repoError    error
		expectedResp *response.NotificationSettingsResponse
	}{
		{
			name:   "success_with_configs",
			userID: 1,
			configs: []model.UserProfile{
				{
					UserID:      1,
					Category:    "notification",
					ConfigKey:   "email_notifications",
					ConfigValue: constants.StringTrue,
				},
				{
					UserID:      1,
					Category:    "notification",
					ConfigKey:   "sms_notifications",
					ConfigValue: constants.StringTrue,
				},
				{
					UserID:      1,
					Category:    "notification",
					ConfigKey:   "push_notifications",
					ConfigValue: constants.StringFalse,
				},
				{
					UserID:      1,
					Category:    "notification",
					ConfigKey:   "marketing_emails",
					ConfigValue: constants.StringTrue,
				},
			},
			repoError: nil,
			expectedResp: &response.NotificationSettingsResponse{
				EmailNotifications: true,
				SmsNotifications:   true,
				PushNotifications:  false,
				MarketingEmails:    true,
			},
		},
		{
			name:      "success_no_configs",
			userID:    2,
			configs:   []model.UserProfile{},
			repoError: nil,
			expectedResp: &response.NotificationSettingsResponse{
				EmailNotifications: true,  // default
				SmsNotifications:   false, // default
				PushNotifications:  true,  // default
				MarketingEmails:    false, // default
			},
		},
		{
			name:      "record_not_found",
			userID:    3,
			configs:   nil,
			repoError: gorm.ErrRecordNotFound,
			expectedResp: &response.NotificationSettingsResponse{
				EmailNotifications: true,  // default
				SmsNotifications:   false, // default
				PushNotifications:  true,  // default
				MarketingEmails:    false, // default
			},
		},
		{
			name:         "repository_error",
			userID:       4,
			configs:      nil,
			repoError:    fmt.Errorf("database error"),
			expectedResp: nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			mockRepo.On("GetUserConfigsByCategory", mock.Anything, tc.userID, "notification").Return(tc.configs, tc.repoError)

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileWarnContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			result, err := service.GetNotificationSettings(context.Background(), tc.userID)

			// Assertions
			if tc.expectedResp == nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResp.EmailNotifications, result.EmailNotifications)
				assert.Equal(t, tc.expectedResp.SmsNotifications, result.SmsNotifications)
				assert.Equal(t, tc.expectedResp.PushNotifications, result.PushNotifications)
				assert.Equal(t, tc.expectedResp.MarketingEmails, result.MarketingEmails)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestUpdateNotificationSettings tests the UpdateNotificationSettings method
func TestUpdateNotificationSettings(t *testing.T) {
	// Test cases
	testCases := []struct {
		name             string
		userID           uint
		request          *request.NotificationSettingsRequest
		savingErrors     map[string]error
		expectSaveConfig bool
		expectedError    error
	}{
		{
			name:   "success",
			userID: 1,
			request: &request.NotificationSettingsRequest{
				EmailNotifications: true,
				SmsNotifications:   false,
				PushNotifications:  true,
				MarketingEmails:    false,
			},
			savingErrors:     map[string]error{},
			expectSaveConfig: true,
			expectedError:    nil,
		},
		{
			name:   "error_saving_email_notifications",
			userID: 2,
			request: &request.NotificationSettingsRequest{
				EmailNotifications: true,
				SmsNotifications:   false,
				PushNotifications:  true,
				MarketingEmails:    false,
			},
			savingErrors: map[string]error{
				"email_notifications": fmt.Errorf("database error"),
			},
			expectSaveConfig: true,
			expectedError:    errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.notification_settings_update_failed"),
		},
		{
			name:   "error_saving_sms_notifications",
			userID: 3,
			request: &request.NotificationSettingsRequest{
				EmailNotifications: true,
				SmsNotifications:   true,
				PushNotifications:  true,
				MarketingEmails:    false,
			},
			savingErrors: map[string]error{
				"sms_notifications": fmt.Errorf("database error"),
			},
			expectSaveConfig: true,
			expectedError:    errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.notification_settings_update_failed"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations for each setting
			if tc.expectSaveConfig {
				// Email notifications
				mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
					return config.UserID == tc.userID &&
						config.Category == "notification" &&
						config.ConfigKey == "email_notifications" &&
						config.ConfigValue == boolToString(tc.request.EmailNotifications)
				})).Return(tc.savingErrors["email_notifications"])

				// If the first save succeeded, continue with the next ones
				if tc.savingErrors["email_notifications"] == nil {
					// SMS notifications
					mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
						return config.UserID == tc.userID &&
							config.Category == "notification" &&
							config.ConfigKey == "sms_notifications" &&
							config.ConfigValue == boolToString(tc.request.SmsNotifications)
					})).Return(tc.savingErrors["sms_notifications"])

					if tc.savingErrors["sms_notifications"] == nil {
						// Push notifications
						mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
							return config.UserID == tc.userID &&
								config.Category == "notification" &&
								config.ConfigKey == "push_notifications" &&
								config.ConfigValue == boolToString(tc.request.PushNotifications)
						})).Return(tc.savingErrors["push_notifications"])

						if tc.savingErrors["push_notifications"] == nil {
							// Marketing emails
							mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
								return config.UserID == tc.userID &&
									config.Category == "notification" &&
									config.ConfigKey == "marketing_emails" &&
									config.ConfigValue == boolToString(tc.request.MarketingEmails)
							})).Return(tc.savingErrors["marketing_emails"])
						}
					}
				}
			}

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			err := service.UpdateNotificationSettings(context.Background(), tc.userID, tc.request)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestGetSecuritySettings tests the GetSecuritySettings method
func TestGetSecuritySettings(t *testing.T) {
	// Test cases
	testCases := []struct {
		name         string
		userID       uint
		configs      []model.UserProfile
		repoError    error
		expectedResp *response.SecuritySettingsResponse
	}{
		{
			name:   "success_with_configs",
			userID: 1,
			configs: []model.UserProfile{
				{
					UserID:      1,
					Category:    "security",
					ConfigKey:   "login_alerts",
					ConfigValue: "true",
				},
				{
					UserID:      1,
					Category:    "security",
					ConfigKey:   "session_timeout",
					ConfigValue: "3600", // 1 hour
				},
			},
			repoError: nil,
			expectedResp: &response.SecuritySettingsResponse{
				LoginAlerts:    true,
				SessionTimeout: 3600,
			},
		},
		{
			name:      "success_no_configs",
			userID:    2,
			configs:   []model.UserProfile{},
			repoError: nil,
			expectedResp: &response.SecuritySettingsResponse{
				LoginAlerts:    true,                                   // default
				SessionTimeout: constants.DefaultSessionTimeoutSeconds, // default
			},
		},
		{
			name:      "record_not_found",
			userID:    3,
			configs:   nil,
			repoError: gorm.ErrRecordNotFound,
			expectedResp: &response.SecuritySettingsResponse{
				LoginAlerts:    true,                                   // default
				SessionTimeout: constants.DefaultSessionTimeoutSeconds, // default
			},
		},
		{
			name:         "repository_error",
			userID:       4,
			configs:      nil,
			repoError:    fmt.Errorf("database error"),
			expectedResp: nil,
		},
		{
			name:   "invalid_session_timeout",
			userID: 5,
			configs: []model.UserProfile{
				{
					UserID:      5,
					Category:    "security",
					ConfigKey:   "session_timeout",
					ConfigValue: "invalid", // not a valid number
				},
			},
			repoError: nil,
			expectedResp: &response.SecuritySettingsResponse{
				LoginAlerts:    true,                                   // default
				SessionTimeout: constants.DefaultSessionTimeoutSeconds, // default since invalid value provided
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			mockRepo.On("GetUserConfigsByCategory", mock.Anything, tc.userID, "security").Return(tc.configs, tc.repoError)

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("UserProfileErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			result, err := service.GetSecuritySettings(context.Background(), tc.userID)

			// Assertions
			if tc.expectedResp == nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResp.LoginAlerts, result.LoginAlerts)
				assert.Equal(t, tc.expectedResp.SessionTimeout, result.SessionTimeout)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

// TestUpdateSecuritySettings tests the UpdateSecuritySettings method
func TestUpdateSecuritySettings(t *testing.T) {
	// Test cases
	testCases := []struct {
		name             string
		userID           uint
		request          *request.SecuritySettingsRequest
		loginAlertsError error
		timeoutError     error
		expectedError    error
	}{
		{
			name:   "success",
			userID: 1,
			request: &request.SecuritySettingsRequest{
				LoginAlerts:    true,
				SessionTimeout: 3600,
			},
			loginAlertsError: nil,
			timeoutError:     nil,
			expectedError:    nil,
		},
		{
			name:   "error_saving_login_alerts",
			userID: 2,
			request: &request.SecuritySettingsRequest{
				LoginAlerts:    false,
				SessionTimeout: 1800,
			},
			loginAlertsError: fmt.Errorf("database error"),
			timeoutError:     nil,
			expectedError:    errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.security_settings_update_failed"),
		},
		{
			name:   "error_saving_session_timeout",
			userID: 3,
			request: &request.SecuritySettingsRequest{
				LoginAlerts:    true,
				SessionTimeout: 7200,
			},
			loginAlertsError: nil,
			timeoutError:     fmt.Errorf("database error"),
			expectedError:    errors.WrapError(fmt.Errorf("database error"), errors.CodeInternalError, "user_profile.security_settings_update_failed"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(MockUserProfileRepository)
			mockLogger := new(MockLogger)
			mockI18n := i18n.NewI18n()

			// Set up mock expectations
			// Login alerts setting
			mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
				return config.UserID == tc.userID &&
					config.Category == "security" &&
					config.ConfigKey == "login_alerts" &&
					config.ConfigValue == boolToString(tc.request.LoginAlerts)
			})).Return(tc.loginAlertsError)

			// Only expect session timeout saving if login alerts was successful
			if tc.loginAlertsError == nil {
				// Session timeout setting
				timeoutStr, _ := json.Marshal(tc.request.SessionTimeout)
				mockRepo.On("SaveUserConfig", mock.Anything, mock.MatchedBy(func(config *model.UserProfile) bool {
					return config.UserID == tc.userID &&
						config.Category == "security" &&
						config.ConfigKey == "session_timeout" &&
						config.ConfigValue == string(timeoutStr)
				})).Return(tc.timeoutError)
			}

			// Mock logger calls (ignore details for simplicity)
			mockLogger.On("UserProfileInfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("ErrorContext", mock.Anything, mock.Anything, mock.Anything).Return()

			// Create service with mocks
			service := NewUserProfileService(mockRepo, mockLogger, mockI18n)

			// Call the method being tested
			err := service.UpdateSecuritySettings(context.Background(), tc.userID, tc.request)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockRepo.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
