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
	s.logger.InfoContext(ctx, "开始用户注册", logger.String("username", req.Username))

	// 1. 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查用户名失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "检查用户名失败")
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeUserAlreadyExists, "用户名已存在")
	}

	// 2. 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查邮箱失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "检查邮箱失败")
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeEmailAlreadyExists, "邮箱已存在")
	}

	// 3. 加密密码
	hashedPassword := utils.SHA256Hash(req.Password)

	// 4. 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       UserStatusActive,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "创建用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "创建用户失败")
	}

	s.logger.InfoContext(ctx, "用户注册成功", logger.Uint("user_id", user.ID))

	return s.buildUserResponse(user), nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	s.logger.InfoContext(ctx, "开始用户登录", logger.String("username", req.Username))

	// 1. 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "用户不存在", logger.String("username", req.Username))
			return nil, errors.ErrInvalidCredentials
		}
		s.logger.ErrorContext(ctx, "查找用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "查找用户失败")
	}

	// 2. 检查用户状态
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "用户已被禁用", logger.String("username", req.Username))
		return nil, errors.NewAppError(errors.CodeForbidden, "用户已被禁用")
	}

	// 3. 验证密码
	if user.PasswordHash != utils.SHA256Hash(req.Password) {
		s.logger.WarnContext(ctx, "密码验证失败", logger.String("username", req.Username))
		return nil, errors.ErrInvalidCredentials
	}

	// 4. 生成 JWT
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "生成token失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "生成token失败")
	}

	// 5. 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WarnContext(ctx, "更新登录时间失败", logger.ErrorField(err))
		// 这里不返回错误，因为登录已经成功
	}

	s.logger.InfoContext(ctx, "用户登录成功", logger.Uint("user_id", user.ID))

	return &response.UserLoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(constants.TokenExpireHours * time.Hour), // 24小时过期
		User:      *s.buildUserResponse(user),
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "获取用户资料", logger.Uint("user_id", userID))

	user, err := s.userRepo.GetByIDWithRelations(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		s.logger.ErrorContext(ctx, "获取用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户失败")
	}

	// 获取统计信息
	stats, err := s.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.WarnContext(ctx, "获取用户统计失败", logger.ErrorField(err))
		// 使用默认值
		stats = &repository.UserStats{}
	}

	profile := &response.UserProfileResponse{
		UserResponse:     *s.buildUserResponse(user),
		LoginCount:       stats.LoginCount,
		ApplicationCount: stats.ApplicationCount,
		WorkflowCount:    stats.WorkflowCount,
	}

	return profile, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(
	ctx context.Context,
	userID uint,
	req *request.UserUpdateProfileRequest,
) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "更新用户资料", logger.Uint("user_id", userID))

	// 1. 获取当前用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 检查邮箱唯一性（如果要更新邮箱）
	if req.Email != nil && *req.Email != user.Email {
		if err := s.validateEmailUniqueness(ctx, *req.Email, userID); err != nil {
			return nil, err
		}
	}

	// 3. 更新用户字段
	s.updateUserFields(user, req)

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "更新用户失败")
	}

	s.logger.InfoContext(ctx, "用户资料更新成功", logger.Uint("user_id", userID))

	return s.buildUserResponse(user), nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "修改用户密码", logger.Uint("user_id", userID))

	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 验证旧密码
	if user.PasswordHash != utils.SHA256Hash(req.OldPassword) {
		s.logger.WarnContext(ctx, "旧密码验证失败", logger.Uint("user_id", userID))
		return errors.ErrInvalidCredentials
	}

	// 3. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 4. 更新密码
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新密码失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "更新密码失败")
	}

	s.logger.InfoContext(ctx, "密码修改成功", logger.Uint("user_id", userID))

	return nil
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context,
	req *request.UserListRequest) (*response.UserListResponse, int64, error) {
	s.logger.InfoContext(ctx, "获取用户列表",
		logger.Int("page", req.Page),
		logger.Int("pageSize", req.PageSize))

	// 构建过滤器
	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.GroupID != nil {
		filters["group_id"] = *req.GroupID
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

	// 获取用户列表
	users, total, err := s.userRepo.ListWithRelations(ctx, req.GetOffset(), req.GetPageSize(), filters)
	if err != nil {
		s.logger.ErrorContext(ctx, "获取用户列表失败", logger.ErrorField(err))
		return nil, 0, errors.WrapError(err, errors.CodeInternalError, "获取用户列表失败")
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

// GetUser 获取用户详情
func (s *userService) GetUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "获取用户详情", logger.Uint("user_id", userID))

	user, err := s.userRepo.GetByIDWithRelations(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		s.logger.ErrorContext(ctx, "获取用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户失败")
	}

	return s.buildUserResponse(user), nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *request.UserCreateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "创建用户", logger.String("username", req.Username))

	// 1. 验证用户创建
	if err := s.validateUserCreation(ctx, req); err != nil {
		return nil, err
	}

	// 2. 创建用户对象
	user := s.createUserFromRequest(req)

	// 3. 保存用户
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "创建用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "创建用户失败")
	}

	s.logger.InfoContext(ctx, "用户创建成功", logger.Uint("user_id", user.ID))

	return s.buildUserResponse(user), nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context,
	userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "更新用户", logger.Uint("user_id", userID))

	// 1. 获取当前用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		s.logger.ErrorContext(ctx, "获取用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户失败")
	}

	// 2. 如果要更新邮箱，检查邮箱是否已存在
	if req.Email != nil && *req.Email != user.Email {
		if err := s.validateEmailUniqueness(ctx, *req.Email, userID); err != nil {
			return nil, err
		}
	}

	// 3. 更新用户信息
	s.updateUserFromRequest(user, req)

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新用户失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "更新用户失败")
	}

	s.logger.InfoContext(ctx, "用户更新成功", logger.Uint("user_id", userID))

	return s.buildUserResponse(user), nil
}

// UpdateUserStatus 更新用户状态
func (s *userService) UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error {
	s.logger.InfoContext(ctx, "更新用户状态", logger.Uint("user_id", userID), logger.Int("status", req.Status))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	user.Status = req.Status
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新用户状态失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "更新用户状态失败")
	}

	s.logger.InfoContext(ctx, "用户状态更新成功", logger.Uint("user_id", userID))
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "删除用户", logger.Uint("user_id", userID))

	// 1. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		s.logger.ErrorContext(ctx, "检查用户失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "检查用户失败")
	}

	// 2. 删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "删除用户失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "删除用户失败")
	}

	s.logger.InfoContext(ctx, "用户删除成功", logger.Uint("user_id", userID))
	return nil
}

// UpdateUserPassword 更新用户密码（管理员操作）
func (s *userService) UpdateUserPassword(
	ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest,
) error {
	s.logger.InfoContext(ctx, "更新用户密码", logger.Uint("user_id", userID))

	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "用户不存在")
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 加密新密码
	hashedPassword := utils.SHA256Hash(req.NewPassword)

	// 3. 更新密码
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新密码失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "更新密码失败")
	}

	s.logger.InfoContext(ctx, "用户密码更新成功", logger.Uint("user_id", userID))
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
		return errors.WrapError(err, errors.CodeInternalError, "检查邮箱唯一性失败")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "邮箱已存在")
	}
	return nil
}

// updateUserFields 更新用户字段
func (s *userService) updateUserFields(user *model.User, req *request.UserUpdateProfileRequest) {
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

// validateUserCreation 验证用户创建
func (s *userService) validateUserCreation(ctx context.Context, req *request.UserCreateRequest) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查用户名失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "检查用户名失败")
	}
	if exists {
		return errors.NewAppError(errors.CodeUserAlreadyExists, "用户名已存在")
	}

	// 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查邮箱失败", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "检查邮箱失败")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "邮箱已存在")
	}
	return nil
}

// createUserFromRequest 从请求创建用户对象
func (s *userService) createUserFromRequest(req *request.UserCreateRequest) *model.User {
	user := &model.User{
		GroupID:      req.GroupID,
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

// updateUserFromRequest 从请求更新用户对象
func (s *userService) updateUserFromRequest(user *model.User, req *request.UserUpdateRequest) {
	if req.GroupID != nil {
		user.GroupID = *req.GroupID
	}
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
		GroupID:     user.GroupID,
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

	if user.Group != nil {
		resp.Group = &response.UserGroupResponse{
			ID:          user.Group.ID,
			Name:        user.Group.Name,
			Code:        user.Group.Code,
			Description: user.Group.Description,
		}
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
