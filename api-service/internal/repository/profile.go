package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/logger"
)

type profileRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewProfileRepository creates a new profile repository instance
func NewProfileRepository(db *gorm.DB, logger logger.Logger) repository.ProfileRepository {
	return &profileRepository{
		db:     db,
		logger: logger,
	}
}

// GetUserProfile get user profile configuration by key
func (r *profileRepository) GetUserProfile(ctx context.Context, userID uint, configKey string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ? AND config_key = ?", userID, configKey).First(&profile).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get user profile",
			logger.Uint("user_id", userID),
			logger.String("config_key", configKey),
			logger.ErrorField(err))
		return nil, err
	}
	return &profile, nil
}

// SetProfile sets or updates a user profile config (user_profile表)
func (r *profileRepository) SetProfile(ctx context.Context, userID uint, configKey, configValue string) error {
	var profile model.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND config_key = ?", userID, configKey).
		First(&profile).Error

	now := time.Now()

	if err == gorm.ErrRecordNotFound {
		// 插入新配置
		profile = model.UserProfile{
			UserID:      userID,
			ConfigKey:   configKey,
			ConfigValue: configValue,
			UpdatedAt:   now,
		}
		err = r.db.WithContext(ctx).Create(&profile).Error
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to create user profile config",
				logger.Uint("user_id", userID),
				logger.String("config_key", configKey),
				logger.ErrorField(err))
			return err
		}
		return nil
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query user profile config",
			logger.Uint("user_id", userID),
			logger.String("config_key", configKey),
			logger.ErrorField(err))
		return err
	}

	// 更新已有配置
	err = r.db.WithContext(ctx).
		Model(&model.UserProfile{}).
		Where("user_id = ? AND config_key = ?", userID, configKey).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"updated_at":   now,
		}).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update user profile config",
			logger.Uint("user_id", userID),
			logger.String("config_key", configKey),
			logger.ErrorField(err))
	}
	return err
}

// SetUserProfile set or update user profile fields
func (r *profileRepository) SetUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update user profile",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
	}
	return err
}

// GetUserProfiles get user profile configurations by category
func (r *profileRepository) GetUserProfiles(ctx context.Context, userID uint, category string) ([]model.UserProfile, error) {
	var profiles []model.UserProfile
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&profiles).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get user profiles",
			logger.Uint("user_id", userID),
			logger.String("category", category),
			logger.ErrorField(err))
		return nil, err
	}
	return profiles, nil
}

// GetUserWithProfile get user with basic profile information
func (r *profileRepository) GetUserWithProfile(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get user",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, err
	}
	return &user, nil
}

// UpdateUserPassword update user password hash
func (r *profileRepository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":       passwordHash,
			"password_changed_at": &now,
			"updated_at":          now,
		}).Error

	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to update user password",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
	}
	return err
}

// GetUserTwoFactors get user's two-factor authentication configurations
func (r *profileRepository) GetUserTwoFactors(ctx context.Context, userID uint) ([]model.UserTwoFactor, error) {
	var twoFactors []model.UserTwoFactor
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_enabled = ?", userID, true).Find(&twoFactors).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get user two factors",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, err
	}
	return twoFactors, nil
}

// GetUserTwoFactorByType get user's two-factor authentication by type
func (r *profileRepository) GetUserTwoFactorByType(ctx context.Context, userID uint, tfaType string) (*model.UserTwoFactor, error) {
	var twoFactor model.UserTwoFactor
	err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, tfaType).First(&twoFactor).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		r.logger.ErrorContext(ctx, "Failed to get user two factor by type",
			logger.Uint("user_id", userID),
			logger.String("type", tfaType),
			logger.ErrorField(err))
		return nil, err
	}
	return &twoFactor, nil
}

// CreateUserTwoFactor create a new two-factor authentication configuration
func (r *profileRepository) CreateUserTwoFactor(ctx context.Context, twoFactor *model.UserTwoFactor) error {
	err := r.db.WithContext(ctx).Create(twoFactor).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create user two factor",
			logger.Uint("user_id", twoFactor.UserID),
			logger.ErrorField(err))
	}
	return err
}

// DeleteUserTwoFactor delete user's two-factor authentication by type
func (r *profileRepository) DeleteUserTwoFactor(ctx context.Context, userID uint, tfaType string) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, tfaType).Delete(&model.UserTwoFactor{}).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to delete user two factor",
			logger.Uint("user_id", userID),
			logger.String("type", tfaType),
			logger.ErrorField(err))
	}
	return err
}

// GetUserLoginHistory get user's login history with pagination
func (r *profileRepository) GetLoginHistory(ctx context.Context, userID uint, req *request.LoginHistoryListRequest) ([]model.UserLoginHistory, int64, error) {
	var history []model.UserLoginHistory
	var total int64

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	query := r.db.WithContext(ctx).Model(&model.UserLoginHistory{}).Where("user_id = ?", userID)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to count login history",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Order("login_time DESC").Offset(offset).Limit(pageSize).Find(&history).Error
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get login history",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, 0, err
	}

	return history, total, nil
}

// GetActiveSessionCount get count of active sessions for user
func (r *profileRepository) GetActiveSessionCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	// This would typically query a sessions table
	// For now, return a mock value
	return count, nil
}
