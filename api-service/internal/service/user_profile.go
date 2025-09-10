package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"

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
func (s *userProfileService) GetUserProfile(ctx context.Context, req *request.UserProfileRequest) (*response.UserProfileResponse, error) {
	// 确定要查询的用户ID
	userID := req.UserID

	s.logger.InfoContext(ctx, "Getting user profile data",
		logger.Uint("userID", userID))

	// 从仓储层获取用户信息
	user, err := s.profileRepo.GetUserProfileByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user profile from repository",
			logger.Uint("userID", userID),
			logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeNotFound, s.i18n.T(ctx, "user_profile.not_found"))
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
	for i := range user.Roles {
		profileResp.Roles = append(profileResp.Roles, response.RoleResponse{
			ID:   user.Roles[i].ID,
			Code: user.Roles[i].Code,
			Name: user.Roles[i].Name,
		})
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
		return s.GetUserProfile(ctx, &request.UserProfileRequest{UserID: userID})
	}

	// 更新资料
	err := s.profileRepo.UpdateUserProfile(ctx, userID, updateData)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user profile",
			logger.Uint("userID", userID),
			logger.ErrorField(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}

		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.update_failed"))
	}

	// 返回更新后的资料
	return s.GetUserProfile(ctx, &request.UserProfileRequest{UserID: userID})
}

// ChangeProfilePassword 修改用户个人密码
func (s *userProfileService) ChangeProfilePassword(ctx context.Context, userID uint, req *request.ProfileChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "Changing user profile password", logger.Uint("userID", userID))

	// 使用请求中的UserID而不是传入的userID
	targetUserID := req.UserID

	// 1. 获取用户信息
	user, err := s.profileRepo.GetUserProfileByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeNotFound, s.i18n.T(ctx, "user_profile.not_found"))
		}
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.get_failed"))
	}

	// 2. 验证旧密码
	if user.PasswordHash != utils.SHA256Hash(req.OldPassword) {
		s.logger.WarnContext(ctx, "Old password verification failed", logger.Uint("userID", targetUserID))
		return errors.NewAppError(errors.CodeInvalidCredentials, s.i18n.T(ctx, "user_profile.password_verification_failed"))
	}

	// 3. 确认新密码与确认密码一致
	if req.NewPassword != req.ConfirmPassword {
		s.logger.WarnContext(ctx, "Password confirmation mismatch", logger.Uint("userID", targetUserID))
		return errors.NewAppError(errors.CodeValidationError, s.i18n.T(ctx, "user_profile.password_mismatch"))
	}

	// 4. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 5. 更新密码 - 使用专门的密码更新方法
	err = s.profileRepo.UpdateUserPassword(ctx, targetUserID, hashedPassword)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password",
			logger.Uint("userID", targetUserID),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "user_profile.password_update_failed"))
	}

	s.logger.InfoContext(ctx, "User password changed successfully", logger.Uint("userID", targetUserID))
	return nil
}
