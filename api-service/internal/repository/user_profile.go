package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/logger"
	"context"

	"gorm.io/gorm"
)

// userProfileRepository 是用户资料仓储的实现
type userProfileRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserProfileRepository 创建用户资料仓储实例
func NewUserProfileRepository(db *gorm.DB, logger logger.Logger) repository.UserProfileRepository {
	return &userProfileRepository{
		db:     db,
		logger: logger,
	}
}

// GetUserProfileByID 根据用户ID获取用户资料
func (r *userProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	r.logger.InfoContext(ctx, "Getting user profile by ID", logger.Uint("userID", userID))

	var user model.User
	// 添加 status != -1 条件，只查询未删除的用户
	result := r.db.Where("status != ?", -1).First(&user, userID)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to get user profile", logger.ErrorField(result.Error))
		return nil, result.Error
	}

	return &user, nil
}

// LoadUserRoles 加载用户的角色信息
func (r *userProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	r.logger.InfoContext(ctx, "Loading user roles", logger.Uint("userID", user.ID))

	result := r.db.Model(user).Association("Roles").Find(&user.Roles)
	if result != nil {
		r.logger.ErrorContext(ctx, "Failed to load user roles", logger.ErrorField(result))
		return result
	}

	return nil
}

// UpdateUserProfile 更新用户资料
func (r *userProfileRepository) UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error {
	r.logger.InfoContext(ctx, "Updating user profile", logger.Uint("userID", userID))

	result := r.db.Model(&model.User{}).Where("id = ? AND status != ?", userID, -1).Updates(updateData)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to update user profile",
			logger.Uint("userID", userID),
			logger.ErrorField(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.WarnContext(ctx, "No user found to update", logger.Uint("userID", userID))
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateUserPassword 更新用户密码
func (r *userProfileRepository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	r.logger.InfoContext(ctx, "Updating user password", logger.Uint("userID", userID))

	result := r.db.Model(&model.User{}).Where("id = ? AND status != ?", userID, -1).Update("password_hash", passwordHash)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to update user password",
			logger.Uint("userID", userID),
			logger.ErrorField(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.WarnContext(ctx, "No user found to update password", logger.Uint("userID", userID))
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetLoginHistories 获取用户登录历史记录
func (r *userProfileRepository) GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error) {
	var records []model.UserLoginHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserLoginHistory{}).Where("user_id = ?", userID)

	// 计算总记录数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("login_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
