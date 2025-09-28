package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

type UserProfileService interface {
	GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error)

	UpdateUserProfile(ctx context.Context, userID uint, req *request.UserProfileUpdateRequest) (*response.UserProfileResponse, error)

	ChangeProfilePassword(ctx context.Context, userID uint, req *request.ProfileChangePasswordRequest) error

	GetLoginHistories(ctx context.Context, userID uint, req *request.LoginHistoryRequest) (*common.PaginationResponse, error)

	GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error)

	UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error

	GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error)

	UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error
}
