package service

import (
	"api-service/internal/dto"
	"api-service/internal/model"
	"api-service/internal/repository"
	"api-service/pkg/auth"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, error)
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ListUsers(query *dto.UserListQuery) ([]*model.User, *dto.PaginationInfo, error)
	GetUserByID(userID uint) (*model.User, error)
	CreateUser(req *dto.CreateUserRequest) (*model.User, error)
	UpdateUser(userID uint, req *dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(userID uint) error
	ValidatePassword(password string) bool
}

type userService struct {
	userRepo repository.UserRepository
	jwtAuth  *auth.JWTAuth
}

func NewUserService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth) UserService {
	return &userService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (s *userService) Register(username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		GroupID:      1, // 默认用户组ID，需要确保存在
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Status:       1, // 1-启用
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.Status != 1 {
		return "", errors.New("account is disabled")
	}

	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); bcryptErr != nil {
		return "", errors.New("invalid credentials")
	}

	// 获取用户角色，这里简化处理，实际应该查询用户角色关联表
	userRole := "user" // 默认角色
	if len(user.Roles) > 0 {
		userRole = user.Roles[0].Code
	}

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Username, userRole)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) UpdateProfile(userID uint, updates map[string]interface{}) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}

	return s.userRepo.Update(user)
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); bcryptErr != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *userService) ListUsers(query *dto.UserListQuery) ([]*model.User, *dto.PaginationInfo, error) {
	offset := (query.Page - 1) * query.PageSize

	// 构建查询条件
	conditions := map[string]interface{}{}
	if query.Keyword != "" {
		conditions["keyword"] = query.Keyword
	}
	if query.Status != nil {
		conditions["status"] = *query.Status
	}
	if query.RoleID != nil {
		conditions["role_id"] = *query.RoleID
	}

	users, total, err := s.userRepo.ListWithConditions(offset, query.PageSize, conditions, query.Sort, query.Order)
	if err != nil {
		return nil, nil, err
	}

	// 计算分页信息
	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	pagination := &dto.PaginationInfo{
		Page:       query.Page,
		PageSize:   query.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    query.Page < totalPages,
		HasPrev:    query.Page > 1,
	}

	return users, pagination, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.GetByIDWithRoles(userID)
}

// CreateUser 创建用户
func (s *userService) CreateUser(req *dto.CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 验证手机号格式
	if req.Phone != "" && !s.validatePhone(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	status := req.Status
	if status == 0 {
		status = 1 // 默认启用
	}

	timezone := req.Timezone
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	language := req.Language
	if language == "" {
		language = "zh-CN"
	}

	user := &model.User{
		GroupID:      req.GroupID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Gender:       req.Gender,
		Signature:    req.Signature,
		Status:       status,
		Timezone:     timezone,
		Language:     language,
	}

	if err := s.userRepo.CreateWithRoles(user, req.RoleIDs); err != nil {
		return nil, err
	}

	// 重新获取用户信息（包含角色）
	return s.userRepo.GetByIDWithRoles(user.ID)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(userID uint, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 更新字段
	if req.GroupID != nil {
		user.GroupID = *req.GroupID
	}
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Email != nil {
		// 检查邮箱是否已存在（排除当前用户）
		if existingUser, err := s.userRepo.GetByEmail(*req.Email); err == nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}
	if req.Phone != nil {
		if *req.Phone != "" && !s.validatePhone(*req.Phone) {
			return nil, errors.New("invalid phone number format")
		}
		user.Phone = *req.Phone
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
	if req.Status != nil {
		user.Status = *req.Status
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateWithRoles(user, req.RoleIDs); err != nil {
		return nil, err
	}

	// 重新获取用户信息（包含角色）
	return s.userRepo.GetByIDWithRoles(userID)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 检查是否有Owner权限需要转移
	// TODO: 实现权限转移检查逻辑
	// 这里应该检查用户是否拥有某些项目的Owner权限
	// 如果有，需要先转移Owner权限给其他用户

	// 设置用户状态为已删除
	user.Status = -1 // -1表示已删除
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// ValidatePassword 验证密码复杂度
func (s *userService) ValidatePassword(password string) bool {
	// 密码长度检查
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	// 包含大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 包含小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 包含数字
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// validatePhone 验证手机号格式
func (s *userService) validatePhone(phone string) bool {
	// 简单的手机号验证（支持中国手机号）
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}
