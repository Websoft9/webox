package repository

import (
	"api-service/internal/model"
	"context"
)

// UserProfileRepository 用户资料仓储接口
type UserProfileRepository interface {
	// GetUserProfileByID 根据用户ID获取用户资料
	GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error)

	// LoadUserRoles 加载用户的角色信息
	LoadUserRoles(ctx context.Context, user *model.User) error

	// UpdateUserProfile 更新用户资料
	UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error

	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error

	// GetLoginHistories 获取用户登录历史记录
	GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error)

	// 获取用户配置
	GetUserConfig(ctx context.Context, userID uint, category, configKey string) (*model.UserProfile, error)

	// 获取用户分类下的所有配置
	GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error)

	// 保存用户配置（不存在则创建，存在则更新）
	SaveUserConfig(ctx context.Context, userProfile *model.UserProfile) error
}
