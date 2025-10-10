package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// userProfileRepository is the implementation of the user profile repository
type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository instance
func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

// GetUserProfileByID retrieves user information by user ID
func (r *userProfileRepository) GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	// add status != -1 condition to only query non-deleted users
	result := r.db.Where("status != ?", -1).First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// LoadUserRoles loads the user's role information
func (r *userProfileRepository) LoadUserRoles(ctx context.Context, user *model.User) error {
	result := r.db.Model(user).Association("Roles").Find(&user.Roles)
	if result != nil {
		return result
	}

	return nil
}

// UpdateUserProfile updates user profile
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

// UpdateUserPassword updates user password
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

// GetLoginHistories retrieves user login history
func (r *userProfileRepository) GetLoginHistories(ctx context.Context, userID uint, req *request.LoginHistoryRequest) ([]*model.UserLoginHistory, int64, error) {
	var records []*model.UserLoginHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserLoginHistory{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order(req.GetSortOrder()).
		Offset(offset).Limit(limit).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetUserConfig gets a user configuration
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

// GetUserConfigsByCategory gets all configs under a category for a user
func (r *userProfileRepository) GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error) {
	var configs []*model.UserProfile
	result := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

// SaveUserConfig saves user config (create if not exists, update if exists)
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
