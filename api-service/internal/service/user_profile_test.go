package service

import (
        "api-service/internal/dto/request"
        "api-service/internal/model"
        "api-service/pkg/i18n"
        "api-service/pkg/logger"
        "api-service/pkg/utils"
        "context"
        "fmt"
        "testing"
        "time"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/mock"
        "gorm.io/gorm"
)

// 创建用于测试的 mock 仓库
type MockUserProfileRepository struct {
        mock.Mock
}

func (m *MockUserProfileRepository) GetUserProfileByID(ctx context.Context, id uint) (*model.User, error) {
        args := m.Called(ctx, id)
        if args.Get(0) == nil {
                return nil, args.Error(1)
        }
        return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserProfileRepository) UpdateUserProfile(ctx context.Context, id uint, data map[string]interface{}) error {
        args := m.Called(ctx, id, data)
        return args.Error(0)
}

func (m *MockUserProfileRepository) UpdateUserPassword(ctx context.Context, id uint, passwordHash string) error {
        args := m.Called(ctx, id, passwordHash)
        return args.Error(0)
}

func (m *MockUserProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
        args := m.Called(ctx, user)
        return args.Error(0)
}

// Mock Logger
type MockLogger struct {
        mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...logger.Field) {}
func (m *MockLogger) Info(msg string, fields ...logger.Field)  {}
func (m *MockLogger) Warn(msg string, fields ...logger.Field)  {}
func (m *MockLogger) Error(msg string, fields ...logger.Field) {}
func (m *MockLogger) Fatal(msg string, fields ...logger.Field) {}

func (m *MockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (m *MockLogger) InfoContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (m *MockLogger) WarnContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (m *MockLogger) ErrorContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (m *MockLogger) FatalContext(ctx context.Context, msg string, fields ...logger.Field) {}

func (m *MockLogger) IsDebugEnabled() bool { return true }
func (m *MockLogger) IsInfoEnabled() bool  { return true }
func (m *MockLogger) IsWarnEnabled() bool  { return true }
func (m *MockLogger) IsErrorEnabled() bool { return true }

func (m *MockLogger) SetLevel(level logger.Level) {}

func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
        return m
}

// 测试 GetUserProfile 函数
func TestGetUserProfile(t *testing.T) {
        tests := []struct {
                name          string
                userID        uint
                setupMock     func(*MockUserProfileRepository)
                expectSuccess bool
                expectError   string
        }{
                {
                        name:   "成功获取用户资料",
                        userID: 1,
                        setupMock: func(repo *MockUserProfileRepository) {
                                now := time.Now()
                                user := &model.User{
                                        ID:          1,
                                        Username:    "testuser",
                                        Email:       "test@example.com",
                                        Nickname:    "Test User",
                                        Avatar:      "avatar.jpg",
                                        Status:      1,
                                        CreatedAt:   now,
                                        UpdatedAt:   now,
                                        LastLoginAt: &now,
                                        LastLoginIP: "127.0.0.1",
                                        Roles: []model.Role{
                                                {
                                                        BaseModel: model.BaseModel{ID: 1},
                                                        Code:      "user",
                                                        Name:      "普通用户",
                                                },
                                        },
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                        },
                        expectSuccess: true,
                },
                {
                        name:   "用户不存在",
                        userID: 999,
                        setupMock: func(repo *MockUserProfileRepository) {
                                repo.On("GetUserProfileByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.not_found",
                },
                {
                        name:   "数据库错误",
                        userID: 2,
                        setupMock: func(repo *MockUserProfileRepository) {
                                repo.On("GetUserProfileByID", mock.Anything, uint(2)).Return(nil, fmt.Errorf("database error"))
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.not_found",
                },
        }

        for _, tc := range tests {
                t.Run(tc.name, func(t *testing.T) {
                        mockRepo := new(MockUserProfileRepository)
                        mockLogger := new(MockLogger)
                        mockI18n := i18n.NewI18n()

                        tc.setupMock(mockRepo)

                        service := &userProfileService{
                                profileRepo: mockRepo,
                                logger:      mockLogger,
                                i18n:        mockI18n,
                        }
                        
                        ctx := context.Background()
                        req := &request.UserProfileRequest{UserID: tc.userID}

                        result, err := service.GetUserProfile(ctx, req)

                        if tc.expectSuccess {
                                assert.NoError(t, err)
                                assert.NotNil(t, result)
                                assert.Equal(t, tc.userID, result.ID)
                        } else {
                                assert.Error(t, err)
                                assert.Contains(t, err.Error(), tc.expectError)
                        }

                        mockRepo.AssertExpectations(t)
                })
        }
}

// 测试 UpdateUserProfile 函数
func TestUpdateUserProfile(t *testing.T) {
        nickname := "New Nickname"
        avatar := "new_avatar.jpg"
        phone := "1234567890"
        var gender int = 1
        signature := "Hello world"
        timezone := "Asia/Shanghai"
        language := "zh-CN"

        tests := []struct {
                name          string
                userID        uint
                request       *request.UserProfileUpdateRequest
                setupMock     func(*MockUserProfileRepository)
                expectSuccess bool
                expectError   string
        }{
                {
                        name:   "成功更新所有字段",
                        userID: 1,
                        request: &request.UserProfileUpdateRequest{
                                UserID:    1,
                                Nickname:  &nickname,
                                Avatar:    &avatar,
                                Phone:     &phone,
                                Gender:    &gender,
                                Signature: &signature,
                                Timezone:  &timezone,
                                Language:  &language,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                expectedData := map[string]interface{}{
                                        "nickname":  nickname,
                                        "avatar":    avatar,
                                        "phone":     phone,
                                        "gender":    gender,
                                        "signature": signature,
                                        "timezone":  timezone,
                                        "language":  language,
                                }
                                repo.On("UpdateUserProfile", mock.Anything, uint(1), expectedData).Return(nil)

                                now := time.Now()
                                user := &model.User{
                                        ID:        1,
                                        Username:  "testuser",
                                        Email:     "test@example.com",
                                        Nickname:  nickname,
                                        Avatar:    avatar,
                                        Phone:     phone,
                                        Gender:    gender,
                                        Signature: signature,
                                        Timezone:  timezone,
                                        Language:  language,
                                        Status:    1,
                                        CreatedAt: now,
                                        UpdatedAt: now,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                        },
                        expectSuccess: true,
                },
                {
                        name:    "没有需要更新的字段",
                        userID:  1,
                        request: &request.UserProfileUpdateRequest{UserID: 1},
                        setupMock: func(repo *MockUserProfileRepository) {
                                now := time.Now()
                                user := &model.User{
                                        ID:        1,
                                        Username:  "testuser",
                                        Email:     "test@example.com",
                                        Nickname:  "Original Nickname",
                                        Status:    1,
                                        CreatedAt: now,
                                        UpdatedAt: now,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                        },
                        expectSuccess: true,
                },
                {
                        name:   "用户不存在",
                        userID: 999,
                        request: &request.UserProfileUpdateRequest{
                                UserID:   999,
                                Nickname: &nickname,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                expectedData := map[string]interface{}{
                                        "nickname": nickname,
                                }
                                repo.On("UpdateUserProfile", mock.Anything, uint(999), expectedData).Return(gorm.ErrRecordNotFound)
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.not_found",
                },
                {
                        name:   "数据库更新错误",
                        userID: 2,
                        request: &request.UserProfileUpdateRequest{
                                UserID:   2,
                                Nickname: &nickname,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                expectedData := map[string]interface{}{
                                        "nickname": nickname,
                                }
                                repo.On("UpdateUserProfile", mock.Anything, uint(2), expectedData).Return(fmt.Errorf("database error"))
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.update_failed",
                },
        }

        for _, tc := range tests {
                t.Run(tc.name, func(t *testing.T) {
                        mockRepo := new(MockUserProfileRepository)
                        mockLogger := new(MockLogger)
                        mockI18n := i18n.NewI18n()

                        tc.setupMock(mockRepo)

                        service := &userProfileService{
                                profileRepo: mockRepo,
                                logger:      mockLogger,
                                i18n:        mockI18n,
                        }
                        
                        ctx := context.Background()

                        result, err := service.UpdateUserProfile(ctx, tc.userID, tc.request)

                        if tc.expectSuccess {
                                assert.NoError(t, err)
                                if tc.request.Nickname != nil {
                                        assert.Equal(t, *tc.request.Nickname, result.Nickname)
                                }
                        } else {
                                assert.Error(t, err)
                                assert.Contains(t, err.Error(), tc.expectError)
                        }

                        mockRepo.AssertExpectations(t)
                })
        }
}

// 测试 ChangeProfilePassword 函数
func TestChangeProfilePassword(t *testing.T) {
        // 创建测试用的密码哈希
        oldPassword := "oldPassword123"
        oldPasswordHash := utils.SHA256Hash(oldPassword)

        newPassword := "newPassword456"
        newPasswordHash := utils.SHA256Hash(newPassword)

        tests := []struct {
                name          string
                userID        uint
                request       *request.ProfileChangePasswordRequest
                setupMock     func(*MockUserProfileRepository)
                expectSuccess bool
                expectError   string
        }{
                {
                        name:   "成功修改密码",
                        userID: 1,
                        request: &request.ProfileChangePasswordRequest{
                                UserID:          1,
                                OldPassword:     oldPassword,
                                NewPassword:     newPassword,
                                ConfirmPassword: newPassword,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                user := &model.User{
                                        ID:           1,
                                        Username:     "testuser",
                                        Email:        "test@example.com",
                                        PasswordHash: oldPasswordHash,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                                repo.On("UpdateUserPassword", mock.Anything, uint(1), newPasswordHash).Return(nil)
                        },
                        expectSuccess: true,
                },
                {
                        name:   "用户不存在",
                        userID: 999,
                        request: &request.ProfileChangePasswordRequest{
                                UserID:          999,
                                OldPassword:     oldPassword,
                                NewPassword:     newPassword,
                                ConfirmPassword: newPassword,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                repo.On("GetUserProfileByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.not_found",
                },
                {
                        name:   "旧密码验证失败",
                        userID: 1,
                        request: &request.ProfileChangePasswordRequest{
                                UserID:          1,
                                OldPassword:     "wrongPassword",
                                NewPassword:     newPassword,
                                ConfirmPassword: newPassword,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                user := &model.User{
                                        ID:           1,
                                        Username:     "testuser",
                                        Email:        "test@example.com",
                                        PasswordHash: oldPasswordHash,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.password_verification_failed",
                },
                {
                        name:   "新密码和确认密码不匹配",
                        userID: 1,
                        request: &request.ProfileChangePasswordRequest{
                                UserID:          1,
                                OldPassword:     oldPassword,
                                NewPassword:     newPassword,
                                ConfirmPassword: "differentPassword",
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                user := &model.User{
                                        ID:           1,
                                        Username:     "testuser",
                                        Email:        "test@example.com",
                                        PasswordHash: oldPasswordHash,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.password_mismatch",
                },
                {
                        name:   "更新密码时发生数据库错误",
                        userID: 1,
                        request: &request.ProfileChangePasswordRequest{
                                UserID:          1,
                                OldPassword:     oldPassword,
                                NewPassword:     newPassword,
                                ConfirmPassword: newPassword,
                        },
                        setupMock: func(repo *MockUserProfileRepository) {
                                user := &model.User{
                                        ID:           1,
                                        Username:     "testuser",
                                        Email:        "test@example.com",
                                        PasswordHash: oldPasswordHash,
                                }
                                repo.On("GetUserProfileByID", mock.Anything, uint(1)).Return(user, nil)
                                repo.On("UpdateUserPassword", mock.Anything, uint(1), newPasswordHash).Return(fmt.Errorf("database error"))
                        },
                        expectSuccess: false,
                        expectError:   "user_profile.password_update_failed",
                },
        }

        for _, tc := range tests {
                t.Run(tc.name, func(t *testing.T) {
                        mockRepo := new(MockUserProfileRepository)
                        mockLogger := new(MockLogger)
                        mockI18n := i18n.NewI18n()

                        tc.setupMock(mockRepo)

                        service := &userProfileService{
                                profileRepo: mockRepo,
                                logger:      mockLogger,
                                i18n:        mockI18n,
                        }
                        
                        ctx := context.Background()

                        err := service.ChangeProfilePassword(ctx, tc.userID, tc.request)

                        if tc.expectSuccess {
                                assert.NoError(t, err)
                        } else {
                                assert.Error(t, err)
                                assert.Contains(t, err.Error(), tc.expectError)
                        }

                        mockRepo.AssertExpectations(t)
                })
        }
}
