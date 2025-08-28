package service

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// ProfileService defines the interface for profile-related operations
type ProfileService interface {
	// Profile management
	GetProfile(ctx context.Context, userID uint) (*response.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *request.ProfileUpdateRequest) error

	// Password management
	ChangePassword(ctx context.Context, userID uint, req *request.PasswordChangeRequest) error

	// Two-factor authentication
	GetTwoFactorStatus(ctx context.Context, userID uint) (*response.TwoFactorStatusResponse, error)
	EnableTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorEnableRequest) (*response.TwoFactorEnableResponse, error)
	VerifyTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorVerifyRequest) error
	DisableTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorDisableRequest) error

	// Login history
	GetLoginHistory(ctx context.Context, userID uint, req *request.LoginHistoryListRequest) (*response.LoginHistoryListResponse, error)

	// Settings
	GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error)
	UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error
	GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error)
	UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error
}
