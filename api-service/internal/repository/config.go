package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"

	"gorm.io/gorm"
)

// configRepository implements the ConfigRepository interface
type configRepository struct {
	db *gorm.DB
}

// NewConfigRepository creates a new Config repository instance
func NewConfigRepository(db *gorm.DB) repository.ConfigRepository {
	return &configRepository{db: db}
}

// GetUserProfile retrieves user profile configuration by user ID
func (r *configRepository) GetUserProfile(ctx context.Context, userID uint) ([]*model.UserProfile, error) {
	var profiles []*model.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("sort_order ASC, config_key ASC").
		Find(&profiles).Error
	if err != nil {
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return profiles, nil
}

// GetUserProfileByKey retrieves a specific user profile configuration by user ID and key
func (r *configRepository) GetUserProfileByKey(ctx context.Context, userID uint, configKey string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND config_key = ?", userID, configKey).
		First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &profile, nil
}

// UpdateUserProfile updates user profile configuration
func (r *configRepository) UpdateUserProfile(ctx context.Context, profile *model.UserProfile) error {
	result := r.db.WithContext(ctx).Updates(profile)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}

// CreateUserProfile creates a new user profile configuration
func (r *configRepository) CreateUserProfile(ctx context.Context, profile *model.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// DeleteUserProfile deletes a user profile configuration
func (r *configRepository) DeleteUserProfile(ctx context.Context, userID uint, configKey string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND config_key = ?", userID, configKey).
		Delete(&model.UserProfile{})
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}
