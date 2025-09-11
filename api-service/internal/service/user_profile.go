package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

// userProfileService 是用户个人资料服务的实现
type userProfileService struct {
	profileRepo repository.UserProfileRepository
	logger      logger.Logger
	i18n        *i18n.I18n
}

// NewUserProfileService 创建用户个人资料服务的实例
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

// GetUserProfile 获取指定用户的个人资料
func (s *userProfileService) GetUserProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "Getting user profile data",
		logger.Uint("userID", userID))

	// 从仓储层获取用户信息
	user, err := s.profileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user profile from repository",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
	}

	// 构建响应DTO
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

	// 加载用户角色信息
	err = s.profileRepo.LoadUserRoles(ctx, user)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to load user roles",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		// 角色加载失败不影响主体数据返回
	} else {
		// 添加角色信息
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

// UpdateUserProfile 更新用户个人资料
func (s *userProfileService) UpdateUserProfile(ctx context.Context, userID uint, req *request.UserProfileUpdateRequest) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "Updating user profile", logger.Uint("userID", userID))

	// 构建更新数据
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

	// 如果没有需要更新的字段，直接返回当前资料
	if len(updateData) == 0 {
		s.logger.WarnContext(ctx, "No fields to update for user profile", logger.Uint("userID", userID))
		return s.GetUserProfile(ctx, userID)
	}

	// 更新资料
	err := s.profileRepo.UpdateUserProfile(ctx, userID, updateData)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user profile",
			logger.Uint("userID", userID),
			logger.ErrorField(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}

		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.update_failed"))
	}

	// 返回更新后的资料
	return s.GetUserProfile(ctx, userID)
}

// ChangeProfilePassword 修改用户个人密码
func (s *userProfileService) ChangeProfilePassword(ctx context.Context, userID uint, req *request.ProfileChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "Changing user profile password", logger.Uint("userID", userID))

	// 1. 获取用户信息
	user, err := s.profileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.get_failed"))
	}

	// 2. 验证旧密码
	if user.PasswordHash != utils.SHA256Hash(req.OldPassword) {
		s.logger.WarnContext(ctx, "Old password verification failed", logger.Uint("userID", userID))
		return errors.NewAppError(errors.CodeInvalidCredentials, s.i18n.T(ctx, "user_profile.password_verification_failed"))
	}

	// 3. 确认新密码与确认密码一致
	if req.NewPassword != req.ConfirmPassword {
		s.logger.WarnContext(ctx, "Password confirmation mismatch", logger.Uint("userID", userID))
		return errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "user_profile.password_mismatch"))
	}

	// 4. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 5. 更新密码 - 使用专门的密码更新方法
	err = s.profileRepo.UpdateUserPassword(ctx, userID, hashedPassword)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.password_update_failed"))
	}

	s.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", userID))
	return nil
}

// GetLoginHistories 获取用户登录历史
func (s *userProfileService) GetLoginHistories(ctx context.Context, userID uint, req *request.LoginHistoryRequest) (*response.LoginHistoryResponse, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1 // 默认第1页
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20 // 默认每页20条
	}

	// 获取数据
	records, _, err := s.profileRepo.GetLoginHistories(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.login_history_get_failed"))
	}

	// 构建响应
	result := &response.LoginHistoryResponse{
		Items: make([]response.LoginHistoryItem, 0, len(records)),
	}

	// 转换记录格式
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

// GetNotificationSettings 获取通知设置
func (s *userProfileService) GetNotificationSettings(ctx context.Context, userID uint) (*response.NotificationSettingsResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification settings", logger.Uint("userID", userID))

	configs, err := s.profileRepo.GetUserConfigsByCategory(ctx, userID, "notification")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to get notification settings",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.notification_settings_get_failed"))
	}

	// 设置默认值
	settings := &response.NotificationSettingsResponse{
		EmailNotifications: true,  // 默认开启邮件通知
		SmsNotifications:   false, // 默认关闭短信通知
		PushNotifications:  true,  // 默认开启推送通知
		MarketingEmails:    false, // 默认关闭营销邮件
	}

	// 如果找到配置，则使用配置值
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

// UpdateNotificationSettings 更新通知设置
func (s *userProfileService) UpdateNotificationSettings(ctx context.Context, userID uint, req *request.NotificationSettingsRequest) error {
	s.logger.InfoContext(ctx, "Updating notification settings", logger.Uint("userID", userID))

	// 更新邮件通知设置
	if err := s.saveUserConfig(ctx, userID, "notification", "email_notifications", boolToString(req.EmailNotifications), "Email notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save email notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.EmailNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// 更新短信通知设置
	if err := s.saveUserConfig(ctx, userID, "notification", "sms_notifications", boolToString(req.SmsNotifications), "SMS notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save SMS notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.SmsNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// 更新推送通知设置
	if err := s.saveUserConfig(ctx, userID, "notification", "push_notifications", boolToString(req.PushNotifications), "Push notifications setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save push notification setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.PushNotifications),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	// 更新营销邮件设置
	if err := s.saveUserConfig(ctx, userID, "notification", "marketing_emails", boolToString(req.MarketingEmails), "Marketing emails setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save marketing emails setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.MarketingEmails),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.notification_settings_update_failed"))
	}

	s.logger.InfoContext(ctx, "Notification settings updated successfully", logger.Uint("userID", userID))
	return nil
}

// GetSecuritySettings 获取安全设置
func (s *userProfileService) GetSecuritySettings(ctx context.Context, userID uint) (*response.SecuritySettingsResponse, error) {
	s.logger.InfoContext(ctx, "Getting security settings", logger.Uint("userID", userID))

	configs, err := s.profileRepo.GetUserConfigsByCategory(ctx, userID, "security")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to get security settings",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.security_settings_get_failed"))
	}

	// 设置默认值
	settings := &response.SecuritySettingsResponse{
		LoginAlerts:    true,                                   // 默认开启登录提醒
		SessionTimeout: constants.DefaultSessionTimeoutSeconds, // 默认会话超时时间30分钟
	}

	// 如果找到配置，则使用配置值
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

// UpdateSecuritySettings 更新安全设置
func (s *userProfileService) UpdateSecuritySettings(ctx context.Context, userID uint, req *request.SecuritySettingsRequest) error {
	s.logger.InfoContext(ctx, "Updating security settings",
		logger.Uint("userID", userID),
		logger.Bool("loginAlerts", req.LoginAlerts),
		logger.Int("sessionTimeout", req.SessionTimeout))

	// 更新登录提醒设置
	if err := s.saveUserConfig(ctx, userID, "security", "login_alerts", boolToString(req.LoginAlerts), "Login alerts setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save login alerts setting",
			logger.Uint("userID", userID),
			logger.Bool("value", req.LoginAlerts),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.security_settings_update_failed"))
	}

	// 更新会话超时设置
	timeoutStr, _ := json.Marshal(req.SessionTimeout)
	if err := s.saveUserConfig(ctx, userID, "security", "session_timeout", string(timeoutStr), "Session timeout setting"); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save session timeout setting",
			logger.Uint("userID", userID),
			logger.Int("value", req.SessionTimeout),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.security_settings_update_failed"))
	}

	s.logger.InfoContext(ctx, "Security settings updated successfully", logger.Uint("userID", userID))
	return nil
}

// 保存用户配置的辅助方法
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

// 将bool转换为字符串的辅助方法
func boolToString(b bool) string {
	if b {
		return constants.StringTrue
	}
	return constants.StringFalse
}
