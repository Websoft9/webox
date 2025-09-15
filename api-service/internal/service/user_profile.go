package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

// userProfileService is the implementation of the user profile service
type userProfileService struct {
	profileRepo repository.UserProfileRepository
	logger      logger.Logger
	i18n        *i18n.I18n
}

// NewUserProfileService creates an instance of the user profile service
func NewUserProfileService(
	profileRepo repository.UserProfileRepository,
	logger logger.Logger,
	i18n *i18n.I18n,
) *userProfileService {
	return &userProfileService{
		profileRepo: profileRepo,
		logger:      logger,
		i18n:        i18n,
	}
}

// GetUserProfile retrieves the profile of the specified user
func (s *userProfileService) GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "Getting user profile data",
		logger.Uint("userID", userID))

	// Fetch user information from the repository layer
	user, err := s.profileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user profile from repository",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
	}

	// Build response DTO
	profileResp := &response.UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Phone:       user.Phone,
		Gender:      user.Gender,
		Signature:   user.Signature,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastLoginIP,
		Timezone:    user.Timezone,
		Language:    user.Language,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	// Load user role information
	err = s.profileRepo.LoadUserRoles(ctx, user)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to load user roles",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		// Role loading failure does not affect returning main profile data
	} else {
		// Add role information
		for i := range user.Roles {
			profileResp.Roles = append(profileResp.Roles, response.RoleResponse{
				ID:   user.Roles[i].ID,
				Code: user.Roles[i].Code,
				Name: user.Roles[i].Name,
			})
		}
	}

	s.logger.InfoContext(ctx, "User profile retrieved successfully", logger.Uint("userID", userID))
	return profileResp, nil
}

// UpdateUserProfile updates the user's profile
func (s *userProfileService) UpdateUserProfile(ctx context.Context, userID uint, req *request.UserProfileUpdateRequest) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "Updating user profile", logger.Uint("userID", userID))

	// Build update data
	updateData := make(map[string]interface{})

	if req.Nickname != nil {
		updateData["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		updateData["avatar"] = *req.Avatar
	}
	if req.Phone != nil {
		updateData["phone"] = *req.Phone
	}
	if req.Gender != nil {
		updateData["gender"] = *req.Gender
	}
	if req.Signature != nil {
		updateData["signature"] = *req.Signature
	}
	if req.Timezone != nil {
		updateData["timezone"] = *req.Timezone
	}
	if req.Language != nil {
		updateData["language"] = *req.Language
	}

	// If there are no fields to update, return the current profile directly
	if len(updateData) == 0 {
		s.logger.WarnContext(ctx, "No fields to update for user profile", logger.Uint("userID", userID))
		return s.GetUserProfile(ctx, userID)
	}

	// Update profile
	err := s.profileRepo.UpdateUserProfile(ctx, userID, updateData)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user profile",
			logger.Uint("userID", userID),
			logger.ErrorField(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}

		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "user_profile.update_failed"))
	}

	// Return the updated profile
	return s.GetUserProfile(ctx, userID)
}

// ChangeProfilePassword changes the user's personal password
func (s *userProfileService) ChangeProfilePassword(ctx context.Context, userID uint, req *request.ProfileChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "Changing user profile password", logger.Uint("userID", userID))

	// 1. Retrieve user information
	user, err := s.profileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "user_profile.get_failed"))
	}

	// 2. Verify old password
	if user.PasswordHash != auth.HashToken(req.OldPassword) {
		s.logger.WarnContext(ctx, "Old password verification failed", logger.Uint("userID", userID))
		return errors.NewAppError(errors.CodeInvalidCredentials, s.i18n.T(ctx, "user_profile.password_verification_failed"))
	}

	// 3. Confirm new password and confirmation match
	if req.NewPassword != req.ConfirmPassword {
		s.logger.WarnContext(ctx, "Password confirmation mismatch", logger.Uint("userID", userID))
		return errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "user_profile.password_mismatch"))
	}

	// 4. Hash the new password
	hashedPassword := auth.HashToken(req.NewPassword)

	// 5. Update the password - use a dedicated password update method
	err = s.profileRepo.UpdateUserPassword(ctx, userID, hashedPassword)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.password_update_failed"))
	}

	s.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", userID))
	return nil
}

// GetLoginHistories retrieves the user's login history
func (s *userProfileService) GetLoginHistories(ctx context.Context, userID uint, req *request.LoginHistoryRequest) (*response.LoginHistoryResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1 // default to page 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20 // default 20 items per page
	}

	// Fetch data
	records, _, err := s.profileRepo.GetLoginHistories(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "user_profile.login_history_get_failed"))
	}

	// Build response
	result := &response.LoginHistoryResponse{
		Items: make([]response.LoginHistoryItem, 0, len(records)),
	}

	// Convert record format
	for i := range records {
		item := response.LoginHistoryItem{
			ID:         records[i].ID,
			IPAddress:  records[i].IPAddress,
			UserAgent:  records[i].UserAgent,
			Device:     records[i].Device,
			Browser:    records[i].Browser,
			Location:   records[i].Location,
			LoginTime:  records[i].LoginTime,
			LogoutTime: records[i].LogoutTime,
		}

		result.Items = append(result.Items, item)
	}

	return result, nil
}

// GetNotificationSettings retrieves notification settings
func (s *userProfileService) GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification settings", logger.Uint("userID", userID))

	configs, err := s.profileRepo.GetUserConfigsByCategory(ctx, userID, "notification")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to get notification settings",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "user_profile.notification_settings_get_failed"))
	}

	// Set default values
	settings := &response.NotificationSettingsResponse{
		EmailNotifications: true,  // default enable email notifications
		SmsNotifications:   false, // default disable SMS notifications
		PushNotifications:  true,  // default enable push notifications
		MarketingEmails:    false, // default disable marketing emails
	}

	// If configurations are found, use their values
	for _, config := range configs {
		switch config.ConfigKey {
		case "email_notifications":
			settings.EmailNotifications = config.ConfigValue == constants.StringTrue
		case "sms_notifications":
			settings.SmsNotifications = config.ConfigValue == constants.StringTrue
		case "push_notifications":
			settings.PushNotifications = config.ConfigValue == constants.StringTrue
		case "marketing_emails":
			settings.MarketingEmails = config.ConfigValue == constants.StringTrue
		}
	}

	s.logger.InfoContext(ctx, "Notification settings retrieved successfully", logger.Uint("userID", userID))
	return settings, nil
}

// UpdateNotificationSettings updates notification settings
func (s *userProfileService) UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error {
	s.logger.InfoContext(ctx, "Updating notification settings", logger.Uint("userID", userID))

	// Update email notification setting
	if err := s.saveUserConfig(ctx, userID, "notification", "email_notifications", boolToString(req.EmailNotifications), "Email notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save email notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.EmailNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// Update SMS notification setting
	if err := s.saveUserConfig(ctx, userID, "notification", "sms_notifications", boolToString(req.SmsNotifications), "SMS notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save SMS notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.SmsNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// Update push notification setting
	if err := s.saveUserConfig(ctx, userID, "notification", "push_notifications", boolToString(req.PushNotifications), "Push notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save push notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.PushNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// Update marketing emails setting
	if err := s.saveUserConfig(ctx, userID, "notification", "marketing_emails", boolToString(req.MarketingEmails), "Marketing emails setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save marketing emails setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.MarketingEmails),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	s.logger.InfoContext(ctx, "Notification settings updated successfully", logger.Uint("userID", userID))
	return nil
}

// GetSecuritySettings retrieves security settings
func (s *userProfileService) GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error) {
	s.logger.InfoContext(ctx, "Getting security settings", logger.Uint("userID", userID))

	configs, err := s.profileRepo.GetUserConfigsByCategory(ctx, userID, "security")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to get security settings",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "user_profile.security_settings_get_failed"))
	}

	// Set default values
	settings := &response.SecuritySettingsResponse{
		LoginAlerts:    true,                                   // default enable login alerts
		SessionTimeout: constants.DefaultSessionTimeoutSeconds, // default session timeout 30 minutes
	}

	// If configurations are found, use their values
	for _, config := range configs {
		switch config.ConfigKey {
		case "login_alerts":
			settings.LoginAlerts = config.ConfigValue == "true"
		case "session_timeout":
			var timeout int
			if err := json.Unmarshal([]byte(config.ConfigValue), &timeout); err == nil && timeout > 0 {
				settings.SessionTimeout = timeout
			}
		}
	}

	s.logger.InfoContext(ctx, "Security settings retrieved successfully", logger.Uint("userID", userID))
	return settings, nil
}

// UpdateSecuritySettings updates security settings
func (s *userProfileService) UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error {
	s.logger.InfoContext(ctx, "Updating security settings",
		logger.Uint("userID", userID),
		logger.Bool("loginAlerts", req.LoginAlerts),
		logger.Int("sessionTimeout", req.SessionTimeout))

	// Update login alerts setting
	if err := s.saveUserConfig(ctx, userID, "security", "login_alerts", boolToString(req.LoginAlerts), "Login alerts setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save login alerts setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.LoginAlerts),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.security_settings_update_failed"))
	}

	// Update session timeout setting
	timeoutStr, _ := json.Marshal(req.SessionTimeout)
	if err := s.saveUserConfig(ctx, userID, "security", "session_timeout", string(timeoutStr), "Session timeout setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save session timeout setting",
			logger.Uint("userID", userID),
			logger.Int("value", req.SessionTimeout),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "user_profile.security_settings_update_failed"))
	}

	s.logger.InfoContext(ctx, "Security settings updated successfully", logger.Uint("userID", userID))
	return nil
}

// saveUserConfig is a helper method to persist user configuration
func (s *userProfileService) saveUserConfig(ctx context.Context, userID uint, category, configKey, configValue, description string) error {
	config := &model.UserProfile{
		UserID:      userID,
		Category:    category,
		ConfigKey:   configKey,
		ConfigValue: configValue,
		Description: description,
	}

	if err := s.profileRepo.SaveUserConfig(ctx, config); err != nil {
		return err
	}

	return nil
}

// boolToString converts a bool to its string representation
func boolToString(b bool) string {
	if b {
		return constants.StringTrue
	}
	return constants.StringFalse
}
