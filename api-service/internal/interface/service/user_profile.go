package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// UserProfileService 定义用户个人资料相关的业务逻辑接口
type UserProfileService interface {
	// GetUserProfile 获取用户个人资料
	GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error)

	// UpdateUserProfile 更新用户个人资料
	UpdateUserProfile(ctx context.Context, userID uint, req *request.UserProfileUpdateRequest) (*response.UserProfileResponse, error)

	// ChangeProfilePassword 修改用户个人密码
	ChangeProfilePassword(ctx context.Context, userID uint, req *request.ProfileChangePasswordRequest) error

	// GetLoginHistories 获取用户登录历史
	GetLoginHistories(ctx context.Context, userID uint, req *request.LoginHistoryRequest) (*response.LoginHistoryResponse, error)

	// 获取通知设置
	GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error)

	// 更新通知设置
	UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error

	// 获取安全设置
	GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error)

	// 更新安全设置
	UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error
}
