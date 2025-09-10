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
}
