package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

// 用户状态常量
const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)

type userService struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

func NewUserService(userRepo repository.UserRepository, logger logger.Logger) service.UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Starting user registration", logger.String("username", req.Username))

	// 1. 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check username", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to check username")
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeUserAlreadyExists, "Username already exists")
	}

	// 2. 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check email", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to check email")
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeEmailAlreadyExists, "Email already exists")
	}

	// 3. 加密密码
	hashedPassword := utils.SHA256Hash(req.Password)

	// 4. Creating user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       UserStatusActive,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to create user")
	}

	s.logger.InfoContext(ctx, "User registration successful", logger.Uint("user_id", user.ID))

	return s.buildUserResponse(user), nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	s.logger.InfoContext(ctx, "Starting user login", logger.String("username", req.Username))

	// 1. 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "User not found", logger.String("username", req.Username))
			return nil, errors.ErrInvalidCredentials
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to find user")
	}

	// 2. 检查用户状态
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "User has been disabled", logger.String("username", req.Username))
		return nil, errors.NewAppError(errors.CodeForbidden, "User has been disabled")
	}

	// 3. 验证密码
	if user.PasswordHash != utils.SHA256Hash(req.Password) {
		s.logger.WarnContext(ctx, "Password verification failed", logger.String("username", req.Username))
		return nil, errors.ErrInvalidCredentials
	}

	// 4. 生成 JWT
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate token", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to generate token")
	}

	// 5. 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WarnContext(ctx, "Failed to update login time", logger.ErrorField(err))
		// 这里不返回错误，因为登录已经成功
	}

	s.logger.InfoContext(ctx, "User login successful", logger.Uint("user_id", user.ID))

	return &response.UserLoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(constants.TokenExpireHours * time.Hour), // 24小时过期
		User:      *s.buildUserResponse(user),
	}, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "Changing user password", logger.Uint("user_id", userID))

	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeInternalError, "Failed to get user information")
	}

	// 2. 验证旧密码
	if user.PasswordHash != utils.SHA256Hash(req.OldPassword) {
		s.logger.WarnContext(ctx, "Old password verification failed", logger.Uint("user_id", userID))
		return errors.ErrInvalidCredentials
	}

	// 3. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 4. 更新密码
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to update password")
	}

	s.logger.InfoContext(ctx, "Password changed successfully", logger.Uint("user_id", userID))

	return nil
}

// ListUsers Getting user list
func (s *userService) ListUsers(ctx context.Context,
	req *request.UserListRequest) (*response.UserListResponse, int64, error) {
	s.logger.InfoContext(ctx, "Getting user list",
		logger.Int("page", req.Page),
		logger.Int("pageSize", req.PageSize))

	// 构建过滤器
	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.Gender != nil {
		filters["gender"] = *req.Gender
	}
	if req.Language != nil {
		filters["language"] = *req.Language
	}
	if req.Keyword != nil && *req.Keyword != "" {
		filters["keyword"] = *req.Keyword
	}

	// Getting user list
	users, total, err := s.userRepo.ListWithRelations(ctx, req.GetOffset(), req.GetPageSize(), filters)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user list", logger.ErrorField(err))
		return nil, 0, errors.WrapError(err, errors.CodeInternalError, "Failed to get user list")
	}

	// 转换为响应格式
	userList := make([]response.UserResponse, 0, len(users))
	for _, user := range users {
		userList = append(userList, *s.buildUserResponse(user))
	}

	return &response.UserListResponse{
		Users: userList,
		Total: total,
	}, total, nil
}

// GetUser Getting user details
func (s *userService) GetUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Getting user details", logger.Uint("user_id", userID))

	user, err := s.userRepo.GetByIDWithRelations(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to get user")
	}

	return s.buildUserResponse(user), nil
}

// CreateUser Creating user
func (s *userService) CreateUser(ctx context.Context, req *request.UserCreateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Creating user", logger.String("username", req.Username))

	// 1. 验证用户创建
	if err := s.validateUserCreation(ctx, req); err != nil {
		return nil, err
	}

	// 2. Creating user对象
	user := s.createUserFromRequest(req)

	// 3. 保存用户
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to create user")
	}

	s.logger.InfoContext(ctx, "User created successfully", logger.Uint("user_id", user.ID))

	return s.buildUserResponse(user), nil
}

// UpdateUser Updating user
func (s *userService) UpdateUser(ctx context.Context,
	userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Updating user", logger.Uint("user_id", userID))

	// 1. 获取当前用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to get user")
	}

	// 2. 如果要更新邮箱，检查邮箱是否已存在
	if req.Email != nil && *req.Email != user.Email {
		if err := s.validateEmailUniqueness(ctx, *req.Email, userID); err != nil {
			return nil, err
		}
	}

	// 3. Updating user信息
	s.updateUserFromRequest(user, req)

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to update user")
	}

	s.logger.InfoContext(ctx, "User updated successfully", logger.Uint("user_id", userID))

	return s.buildUserResponse(user), nil
}

// UpdateUserStatus Updating user状态
func (s *userService) UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error {
	s.logger.InfoContext(ctx, "Updating user status", logger.Uint("user_id", userID), logger.Int("status", req.Status))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeInternalError, "Failed to get user information")
	}

	user.Status = req.Status
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user status", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to update user status")
	}

	s.logger.InfoContext(ctx, "User status updated successfully", logger.Uint("user_id", userID))
	return nil
}

// DeleteUser Deleting user
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "Deleting user", logger.Uint("user_id", userID))

	// 1. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to check user", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to check user")
	}

	// 2. Deleting user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete user", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to delete user")
	}

	s.logger.InfoContext(ctx, "User deleted successfully", logger.Uint("user_id", userID))
	return nil
}

// UpdateUserPassword Updating user密码（管理员操作）
func (s *userService) UpdateUserPassword(
	ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest,
) error {
	s.logger.InfoContext(ctx, "Updating user password", logger.Uint("user_id", userID))

	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeInternalError, "Failed to get user information")
	}

	// 2. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 3. 更新密码
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to update password")
	}

	s.logger.InfoContext(ctx, "User password updated successfully", logger.Uint("user_id", userID))
	return nil
}

// ValidateUserAccess 验证用户访问权限
func (s *userService) ValidateUserAccess(ctx context.Context, userID uint, resource string) error {
	// TODO: 实现具体的权限验证逻辑
	return nil
}

// CheckUserQuota 检查用户配额
func (s *userService) CheckUserQuota(ctx context.Context, userID uint, resourceType string) error {
	// TODO: 实现具体的配额检查逻辑
	return nil
}

// 辅助函数

// validateEmailUniqueness 验证邮箱唯一性
func (s *userService) validateEmailUniqueness(ctx context.Context, email string, userID uint) error {
	exists, err := s.userRepo.ExistsByEmailExcludeID(ctx, email, userID)
	if err != nil {
		return errors.WrapError(err, errors.CodeInternalError, "Failed to check email uniqueness")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "Email already exists")
	}
	return nil
}

// validateUserCreation 验证用户创建
func (s *userService) validateUserCreation(ctx context.Context, req *request.UserCreateRequest) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check username", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to check username")
	}
	if exists {
		return errors.NewAppError(errors.CodeUserAlreadyExists, "Username already exists")
	}

	// 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check email", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "Failed to check email")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "Email already exists")
	}
	return nil
}

// createUserFromRequest 从请求Creating user对象
func (s *userService) createUserFromRequest(req *request.UserCreateRequest) *model.User {
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: utils.SHA256Hash(req.Password),
		Status:       1, // 默认激活
	}

	// 设置可选字段
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Signature != nil {
		user.Signature = *req.Signature
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.Language != nil {
		user.Language = *req.Language
	}

	return user
}

func (s *userService) updateUserFromRequest(user *model.User, req *request.UserUpdateRequest) {
	// Updating user信息
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Signature != nil {
		user.Signature = *req.Signature
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.Language != nil {
		user.Language = *req.Language
	}
}

// buildUserResponse 构建用户响应
func (s *userService) buildUserResponse(user *model.User) *response.UserResponse {
	resp := &response.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
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

	if len(user.Roles) > 0 {
		resp.Roles = make([]response.RoleResponse, len(user.Roles))
		for i := range user.Roles {
			role := &user.Roles[i]
			resp.Roles[i] = response.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Code:        role.Code,
				Description: role.Description,
			}
		}
	}

	return resp
}
