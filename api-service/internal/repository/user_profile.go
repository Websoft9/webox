package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// userProfileRepository 是用户资料仓储的实现
type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository 创建用户资料仓储实例
func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

// GetUserProfileByID 根据用户ID获取用户资料
func (r *userProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	// 添加 status != -1 条件，只查询未删除的用户
	result := r.db.Where("status != ?", -1).First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// LoadUserRoles 加载用户的角色信息
func (r *userProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	result := r.db.Model(user).Association("Roles").Find(&user.Roles)
	if result != nil {
		return result
	}

	return nil
}

// UpdateUserProfile 更新用户资料
func (r *userProfileRepository) UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error {
	result := r.db.Model(&model.User{}).Where("id = ? AND status != ?", userID, -1).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateUserPassword 更新用户密码
func (r *userProfileRepository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	result := r.db.Model(&model.User{}).Where("id = ? AND status != ?", userID, -1).Update("password_hash", passwordHash)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
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

// 获取用户配置
func (r *userProfileRepository) GetUserConfig(ctx context.Context, userID uint, category, configKey string) (*model.UserProfile, error) {
	var config model.UserProfile
	result := r.db.WithContext(ctx).Where("user_id = ? AND category = ? AND config_key = ?", userID, category, configKey).First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}

		return nil, result.Error
	}

	return &config, nil
}

// 获取用户分类下的所有配置
func (r *userProfileRepository) GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error) {
	var configs []*model.UserProfile
	result := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

// 保存用户配置（不存在则创建，存在则更新）
func (r *userProfileRepository) SaveUserConfig(ctx context.Context, userProfile *model.UserProfile) error {
	var existing model.UserProfile
	result := r.db.WithContext(ctx).Where("user_id = ? AND category = ? AND config_key = ?",
		userProfile.UserID, userProfile.Category, userProfile.ConfigKey).First(&existing)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(userProfile).Error; err != nil {
				return err
			}
			return nil
		}

		return result.Error
	}

	existing.ConfigValue = userProfile.ConfigValue
	existing.Description = userProfile.Description
	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return err
	}
	return nil
}
