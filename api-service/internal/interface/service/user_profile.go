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
}
