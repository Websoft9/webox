package service

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/security"
)

type profileService struct {
	profileRepo repository.ProfileRepository
	logger      logger.Logger
}

// NewProfileService creates a new profile service instance
func NewProfileService(profileRepo repository.ProfileRepository, logger logger.Logger) service.ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		logger:      logger,
	}
}

// GetProfile get user profile information
func (s *profileService) GetProfile(ctx context.Context, userID uint) (*response.ProfileResponse, error) {
	user, err := s.profileRepo.GetUserWithProfile(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get profile")
	}

	// Get profile configurations
	profiles, _ := s.profileRepo.GetUserProfiles(ctx, userID, "")

	// Convert to map for easier access
	profileMap := make(map[string]string)
	for i := range profiles {
		profileMap[profiles[i].ConfigKey] = profiles[i].ConfigValue
	}

	// Get two-factor status
	twoFactors, _ := s.profileRepo.GetUserTwoFactors(ctx, userID)

	// Parse gender
	gender := 0
	if genderStr, exists := profileMap["gender"]; exists {
		if g, err := strconv.Atoi(genderStr); err == nil {
			gender = g
		}
	}

	return &response.ProfileResponse{
		ID:               user.ID,
		Username:         user.Username,
		Email:            user.Email,
		Nickname:         profileMap["nickname"],
		Avatar:           profileMap["avatar"],
		Phone:            profileMap["phone"],
		Gender:           gender,
		Signature:        profileMap["signature"],
		Timezone:         profileMap["timezone"],
		Language:         profileMap["language"],
		TwoFactorEnabled: len(twoFactors) > 0,
		LastLoginAt:      user.LastLoginAt,
		LastLoginIP:      profileMap["last_login_ip"],
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}, nil
}

// UpdateProfile update user profile information
func (s *profileService) UpdateProfile(ctx context.Context, userID uint, req *request.ProfileUpdateRequest) error {
	// Check if user exists
	_, err := s.profileRepo.GetUserWithProfile(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return errors.NewAppError(errors.CodeInternalError, "Failed to get profile")
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Signature != nil {
		updates["signature"] = *req.Signature
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}

	// If no fields to update, return early
	if len(updates) == 0 {
		return nil
	}

	// Update all fields in a single operation
	err = s.profileRepo.SetUserProfile(ctx, userID, updates)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "Failed to update profile")
	}

	return nil
}

// ChangePassword change user password
func (s *profileService) ChangePassword(ctx context.Context, userID uint, req *request.PasswordChangeRequest) error {
	user, err := s.profileRepo.GetUserWithProfile(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return errors.NewAppError(errors.CodeInternalError, "Failed to get profile")
	}

	// Verify current password
	if !security.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.NewAppError(errors.CodeInvalidRequest, "Current password is incorrect")
	}

	// Verify new password confirmation
	if req.NewPassword != req.ConfirmPassword {
		return errors.NewAppError(errors.CodeInvalidRequest, "Password confirmation mismatch")
	}

	// Hash new password
	newPasswordHash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "Failed to hash password")
	}

	// Update password
	err = s.profileRepo.UpdateUserPassword(ctx, userID, newPasswordHash)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "Failed to update password")
	}

	return nil
}

// GetTwoFactorStatus get user two-factor authentication status
func (s *profileService) GetTwoFactorStatus(ctx context.Context, userID uint) (*response.TwoFactorStatusResponse, error) {
	twoFactors, err := s.profileRepo.GetUserTwoFactors(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get two-factor status")
	}

	return &response.TwoFactorStatusResponse{
		Enabled: len(twoFactors) > 0,
	}, nil
}

// EnableTwoFactor enable two-factor authentication
func (s *profileService) EnableTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorEnableRequest) (*response.TwoFactorEnableResponse, error) {
	// Check if two-factor is already enabled
	twoFactors, err := s.profileRepo.GetUserTwoFactors(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get two-factor status")
	}

	if len(twoFactors) > 0 {
		return nil, errors.NewAppError(errors.CodeInvalidRequest, "Two-factor authentication already enabled")
	}

	// TODO: Implement two-factor enable logic when security package is complete
	return &response.TwoFactorEnableResponse{}, nil
}

// VerifyTwoFactor verify and activate two-factor authentication
func (s *profileService) VerifyTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorVerifyRequest) error {
	// TODO: Implement two-factor verify logic when security package is complete
	return nil
}

// DisableTwoFactor disable two-factor authentication
func (s *profileService) DisableTwoFactor(ctx context.Context, userID uint, req *request.TwoFactorDisableRequest) error {
	// TODO: Implement two-factor disable logic when security package is complete
	return nil
}

// GetLoginHistory get user login history
func (s *profileService) GetLoginHistory(ctx context.Context, userID uint, req *request.LoginHistoryListRequest) (*response.LoginHistoryListResponse, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// Get login history
	loginHistories, total, err := s.profileRepo.GetLoginHistory(ctx, userID, req)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get login history")
	}

	items := make([]response.LoginHistoryResponse, len(loginHistories))
	for i := range loginHistories {
		items[i] = response.LoginHistoryResponse{
			ID:         loginHistories[i].ID,
			IPAddress:  loginHistories[i].IPAddress,
			UserAgent:  loginHistories[i].UserAgent,
			Location:   loginHistories[i].Location,
			Device:     loginHistories[i].Device,
			Browser:    loginHistories[i].Browser,
			OS:         loginHistories[i].OS,
			LoginTime:  loginHistories[i].LoginTime,
			LogoutTime: loginHistories[i].LogoutTime,
			Status:     loginHistories[i].Status,
		}
	}

	return &response.LoginHistoryListResponse{
		Items: items,
		Total: total,
	}, nil
}

// GetNotificationSettings get user notification settings
func (s *profileService) GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error) {
	// Get notification settings from profile configurations
	profiles, err := s.profileRepo.GetUserProfiles(ctx, userID, "notification_")
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get notification settings")
	}

	// Convert to map for easier access
	settingsMap := make(map[string]string)
	for i := range profiles {
		settingsMap[profiles[i].ConfigKey] = profiles[i].ConfigValue
	}

	// Parse settings with defaults
	emailEnabled := true
	smsEnabled := false
	pushEnabled := true

	if val, exists := settingsMap["notification_email"]; exists {
		emailEnabled = val == constants.BoolTrue
	}
	if val, exists := settingsMap["notification_sms"]; exists {
		smsEnabled = val == constants.BoolTrue
	}
	if val, exists := settingsMap["notification_push"]; exists {
		pushEnabled = val == constants.BoolTrue
	}

	return &response.NotificationSettingsResponse{
		EmailEnabled: emailEnabled,
		SMSEnabled:   smsEnabled,
		PushEnabled:  pushEnabled,
	}, nil
}

// UpdateNotificationSettings update user notification settings
func (s *profileService) UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error {
	// Update email notification setting
	if req.EmailNotifications != nil {
		value := constants.BoolFalse
		if *req.EmailNotifications {
			value = constants.BoolTrue
		}
		err := s.profileRepo.SetProfile(ctx, userID, "notification_email", value)
		if err != nil {
			return errors.NewAppError(errors.CodeInternalError, "Failed to update email notification settings")
		}
	}

	// Update SMS notification setting
	if req.SMSNotifications != nil {
		value := constants.BoolFalse
		if *req.SMSNotifications {
			value = constants.BoolTrue
		}
		err := s.profileRepo.SetProfile(ctx, userID, "notification_sms", value)
		if err != nil {
			return errors.NewAppError(errors.CodeInternalError, "Failed to update SMS notification settings")
		}
	}

	// Update push notification setting
	if req.PushNotifications != nil {
		value := constants.BoolFalse
		if *req.PushNotifications {
			value = constants.BoolTrue
		}
		err := s.profileRepo.SetProfile(ctx, userID, "notification_push", value)
		if err != nil {
			return errors.NewAppError(errors.CodeInternalError, "Failed to update push notification settings")
		}
	}

	return nil
}

// GetSecuritySettings get user security settings
func (s *profileService) GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error) {
	// Get security settings from profile configurations
	profiles, err := s.profileRepo.GetUserProfiles(ctx, userID, "security_")
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "Failed to get security settings")
	}

	// Convert to map for easier access
	settingsMap := make(map[string]string)
	for i := range profiles {
		settingsMap[profiles[i].ConfigKey] = profiles[i].ConfigValue
	}

	// Get two-factor status
	twoFactors, _ := s.profileRepo.GetUserTwoFactors(ctx, userID)

	// Parse settings with defaults
	loginNotification := true
	sessionTimeout := 30 // 30 minutes default

	if val, exists := settingsMap["security_login_notification"]; exists {
		loginNotification = val == constants.BoolTrue
	}
	if val, exists := settingsMap["security_session_timeout"]; exists {
		if timeout, err := strconv.Atoi(val); err == nil {
			sessionTimeout = timeout
		}
	}

	return &response.SecuritySettingsResponse{
		TwoFactorEnabled:  len(twoFactors) > 0,
		LoginNotification: loginNotification,
		SessionTimeout:    sessionTimeout,
	}, nil
}

// UpdateSecuritySettings update user security settings
func (s *profileService) UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error {
	// Update login notification setting
	if req.LoginNotifications != nil {
		value := constants.BoolFalse
		if *req.LoginNotifications {
			value = constants.BoolTrue
		}
		err := s.profileRepo.SetProfile(ctx, userID, "security_login_notification", value)
		if err != nil {
			return errors.NewAppError(errors.CodeInternalError, "Failed to update login notification settings")
		}
	}

	// Update session timeout setting
	if req.SessionTimeout != nil {
		value := fmt.Sprintf("%d", *req.SessionTimeout)
		err := s.profileRepo.SetProfile(ctx, userID, "security_session_timeout", value)
		if err != nil {
			return errors.NewAppError(errors.CodeInternalError, "Failed to update session timeout settings")
		}
	}

	return nil
}
