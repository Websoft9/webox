package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/validator"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// userService 用户业务逻辑实现
type userService struct {
	userRepo repository.UserRepository
	jwtAuth  *auth.JWTAuth
	logger   logger.Logger
}

// NewUserService 创建新的用户Service实例
func NewUserService(
	userRepo repository.UserRepository,
	jwtAuth *auth.JWTAuth,
	logger logger.Logger,
) service.UserService {
	return &userService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
		logger:   logger,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "用户注册开始", logger.String("username", req.Username))

	// 1. 业务验证
	if err := validator.ValidateUsername(req.Username); err != nil {
		s.logger.WarnContext(ctx, "用户名验证失败", logger.String("username", req.Username), logger.ErrorField(err))
		return nil, err
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		s.logger.WarnContext(ctx, "邮箱验证失败", logger.String("email", req.Email), logger.ErrorField(err))
		return nil, err
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		s.logger.WarnContext(ctx, "密码验证失败", logger.ErrorField(err))
		return nil, err
	}

	// 2. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查用户名存在性失败", logger.String("username", req.Username), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "检查用户名失败")
	}
	if exists {
		return nil, errors.ErrUserAlreadyExists
	}

	// 3. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "检查邮箱存在性失败", logger.String("email", req.Email), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "检查邮箱失败")
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// 4. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.ErrorContext(ctx, "密码加密失败", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "密码加密失败")
	}

	// 5. 创建用户模型
	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Status:    "active",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 6. 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "创建用户失败", logger.String("username", req.Username), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "创建用户失败")
	}

	s.logger.InfoContext(ctx, "用户注册成功", logger.String("username", req.Username), logger.Uint("user_id", user.ID))

	// 7. 返回响应
	return s.convertToUserResponse(user), nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	s.logger.InfoContext(ctx, "用户登录开始", logger.String("username", req.Username))

	// 1. 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "用户不存在", logger.String("username", req.Username))
			return nil, errors.ErrInvalidCredentials
		}
		s.logger.ErrorContext(ctx, "获取用户失败", logger.String("username", req.Username), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 检查用户状态
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "用户账号未激活", logger.String("username", req.Username), logger.String("status", user.Status))
		return nil, errors.ErrUserInactive
	}

	// 3. 验证密码
	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); bcryptErr != nil {
		s.logger.WarnContext(ctx, "密码验证失败", logger.String("username", req.Username))
		return nil, errors.ErrInvalidCredentials
	}

	// 4. 生成JWT Token
	token, expiresAt, err := s.jwtAuth.GenerateTokenWithUserInfo(user.ID, user.Username, user.Role)
	if err != nil {
		s.logger.ErrorContext(ctx, "生成Token失败", logger.Uint("user_id", user.ID), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "生成认证令牌失败")
	}

	// 5. 更新登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LoginCount++
	user.UpdatedAt = now

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.WarnContext(ctx, "更新登录信息失败", logger.Uint("user_id", user.ID), logger.ErrorField(err))
		// 这里不返回错误，因为登录已经成功，只是统计信息更新失败
	}

	s.logger.InfoContext(ctx, "用户登录成功", logger.String("username", req.Username), logger.Uint("user_id", user.ID))

	// 6. 返回响应
	return &response.UserLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *s.convertToUserResponse(user),
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	s.logger.InfoContext(ctx, "获取用户资料", logger.Uint("user_id", userID))

	// 1. 获取用户基本信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		s.logger.ErrorContext(ctx, "获取用户失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 获取用户统计信息
	stats, err := s.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		s.logger.WarnContext(ctx, "获取用户统计信息失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		// 统计信息获取失败时使用默认值
		stats = &repository.UserStats{
			LoginCount:       user.LoginCount,
			ApplicationCount: 0,
			WorkflowCount:    0,
		}
	}

	// 3. 构建响应
	profile := &response.UserProfileResponse{
		UserResponse:     *s.convertToUserResponse(user),
		LastLoginAt:      user.LastLoginAt,
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

	// 1. 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 验证邮箱唯一性（如果要更新邮箱）
	if req.Email != nil && *req.Email != user.Email {
		if err := validator.ValidateEmail(*req.Email); err != nil {
			return nil, err
		}

		exists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, errors.WrapError(err, errors.CodeInternalError, "检查邮箱失败")
		}
		if exists {
			return nil, errors.ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}

	// 3. 更新其他字段
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	user.UpdatedAt = time.Now()

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新用户资料失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, "更新用户资料失败")
	}

	s.logger.InfoContext(ctx, "用户资料更新成功", logger.Uint("user_id", userID))
	return s.convertToUserResponse(user), nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "修改用户密码", logger.Uint("user_id", userID))

	// 1. 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 验证旧密码
	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); bcryptErr != nil {
		s.logger.WarnContext(ctx, "旧密码验证失败", logger.Uint("user_id", userID))
		return errors.ErrInvalidPassword
	}

	// 3. 验证新密码强度
	if validErr := validator.ValidatePassword(req.NewPassword); validErr != nil {
		return validErr
	}

	// 4. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.WrapError(err, errors.CodeInternalError, "密码加密失败")
	}

	// 5. 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新密码失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "更新密码失败")
	}

	s.logger.InfoContext(ctx, "用户密码修改成功", logger.Uint("user_id", userID))
	return nil
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(
	ctx context.Context,
	req *request.UserListRequest,
) (*response.UserListResponse, int64, error) {
	s.logger.InfoContext(ctx, "获取用户列表")

	// 1. 构建过滤条件
	filters := make(map[string]interface{})
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Role != "" {
		filters["role"] = req.Role
	}

	// 2. 获取数据
	var users []*model.User
	var total int64
	var err error

	if req.Keyword != "" {
		// 搜索模式
		users, total, err = s.userRepo.Search(ctx, req.Keyword, req.GetOffset(), req.GetPageSize())
	} else {
		// 列表模式
		users, total, err = s.userRepo.List(ctx, req.GetOffset(), req.GetPageSize(), filters)
	}

	if err != nil {
		s.logger.ErrorContext(ctx, "获取用户列表失败", logger.ErrorField(err))
		return nil, 0, errors.WrapError(err, errors.CodeInternalError, "获取用户列表失败")
	}

	// 3. 转换响应
	userResponses := make([]response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.convertToUserResponse(user)
	}

	return &response.UserListResponse{Users: userResponses}, total, nil
}

// GetUser 获取单个用户信息（管理员功能）
func (s *userService) GetUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	return s.convertToUserResponse(user), nil
}

// UpdateUserStatus 更新用户状态（管理员功能）
func (s *userService) UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error {
	s.logger.InfoContext(ctx, "更新用户状态", logger.Uint("user_id", userID), logger.String("new_status", req.Status))

	// 1. 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 验证状态转换
	if err := validator.ValidateUserStatus(user.Status, req.Status); err != nil {
		return err
	}

	// 3. 更新状态
	user.Status = req.Status
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "更新用户状态失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "更新用户状态失败")
	}

	s.logger.InfoContext(ctx, "用户状态更新成功", logger.Uint("user_id", userID), logger.String("new_status", req.Status))
	return nil
}

// DeleteUser 删除用户（管理员功能）
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "删除用户", logger.Uint("user_id", userID))

	// 1. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 2. 执行删除
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "删除用户失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "删除用户失败")
	}

	s.logger.InfoContext(ctx, "用户删除成功", logger.Uint("user_id", userID))
	return nil
}

// ValidateUserAccess 验证用户访问权限
func (s *userService) ValidateUserAccess(ctx context.Context, userID uint, resource string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return errors.WrapError(err, errors.CodeInternalError, "获取用户信息失败")
	}

	// 检查用户状态
	if user.Status != "active" {
		return errors.ErrUserInactive
	}

	// 检查权限
	return validator.ValidateUserPermission(user.Role, resource)
}

// CheckUserQuota 检查用户配额
func (s *userService) CheckUserQuota(ctx context.Context, userID uint, resourceType string) error {
	stats, err := s.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		return errors.WrapError(err, errors.CodeInternalError, "获取用户统计信息失败")
	}

	var currentCount int
	switch resourceType {
	case "applications":
		currentCount = stats.ApplicationCount
	case "workflows":
		currentCount = stats.WorkflowCount
	default:
		currentCount = 0
	}

	return validator.ValidateResourceQuota(userID, resourceType, currentCount)
}

// convertToUserResponse 转换用户模型为响应结构
func (s *userService) convertToUserResponse(user *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
