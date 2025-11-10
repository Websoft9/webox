package repository

import (
	"api-service/internal/model"
	"context"
)

// ConfigRepository defines the interface for configuration center operations
type ConfigRepository interface {
	// GetUserProfile retrieves user profile configuration by user ID
	GetUserProfile(ctx context.Context, userID uint) ([]*model.UserProfile, error)

	// GetUserProfileByKey retrieves a specific user profile configuration by user ID and key
	GetUserProfileByKey(ctx context.Context, userID uint, configKey string) (*model.UserProfile, error)

	// UpdateUserProfile updates user profile configuration
	UpdateUserProfile(ctx context.Context, profile *model.UserProfile) error

	// CreateUserProfile creates a new user profile configuration
	CreateUserProfile(ctx context.Context, profile *model.UserProfile) error

	// DeleteUserProfile deletes a user profile configuration
	DeleteUserProfile(ctx context.Context, userID uint, configKey string) error
}
