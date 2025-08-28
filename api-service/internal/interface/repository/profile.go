package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"context"
)

// ProfileRepository 个人中心数据访问接口
type ProfileRepository interface {
	// 用户配置管理
	GetUserProfile(ctx context.Context, userID uint, configKey string) (*model.UserProfile, error)
	SetProfile(ctx context.Context, userID uint, configKey, configValue string) error
	SetUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	GetUserProfiles(ctx context.Context, userID uint, category string) ([]model.UserProfile, error)

	// 双因子认证管理
	GetUserTwoFactors(ctx context.Context, userID uint) ([]model.UserTwoFactor, error)
	CreateUserTwoFactor(ctx context.Context, twoFactor *model.UserTwoFactor) error
	DeleteUserTwoFactor(ctx context.Context, userID uint, tfaType string) error

	// 登录历史管理
	GetLoginHistory(ctx context.Context, userID uint, req *request.LoginHistoryListRequest) ([]model.UserLoginHistory, int64, error)

	// 用户信息扩展
	GetUserWithProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateUserPassword(ctx context.Context, userID uint, newPasswordHash string) error
}
