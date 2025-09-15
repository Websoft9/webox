package repository

import (
	"api-service/internal/model"
	"context"
)

// UserProfileRepository is the interface for user profile repository
type UserProfileRepository interface {
	// GetUserProfileByID retrieves user profile by user ID
	GetUserProfileByID(ctx context.Context, userID uint) (*model.User, error)

	// LoadUserRoles loads the roles associated with the user
	LoadUserRoles(ctx context.Context, user *model.User) error

	// UpdateUserProfile update profile information
	UpdateUserProfile(ctx context.Context, userID uint, updateData map[string]interface{}) error

	// UpdateUserPassword updates the user's password
	UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error

	// GetLoginHistories gets the login history for a user with pagination
	GetLoginHistories(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginHistory, int64, error)

	// get user configuration by user ID, category, and config key
	GetUserConfig(ctx context.Context, userID uint, category, configKey string) (*model.UserProfile, error)

	// get all user configurations by user ID and category
	GetUserConfigsByCategory(ctx context.Context, userID uint, category string) ([]*model.UserProfile, error)

	// save or update user configuration
	SaveUserConfig(ctx context.Context, userProfile *model.UserProfile) error
}
