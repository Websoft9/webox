package service

import (
	"api-service/internal/model"
	"api-service/internal/repository"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserListResponse 用户列表响应
type UserListResponse struct {
	Items      []UserResponse `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	UserResponse
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID          uint       `json:"id"`
	GroupID     uint       `json:"group_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone"`
	Gender      int8       `json:"gender"`
	Signature   string     `json:"signature"`
	Status      int8       `json:"status"`
	Timezone    string     `json:"timezone"`
	Language    string     `json:"language"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	Roles       []RoleInfo `json:"roles"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	GroupID   uint   `json:"group_id" binding:"required"`
	Username  string `json:"username" binding:"required,min=3,max=32"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=32"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Gender    int8   `json:"gender"`
	Signature string `json:"signature"`
	Timezone  string `json:"timezone"`
	Language  string `json:"language"`
	RoleIDs   []uint `json:"role_ids"`
	Status    int8   `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	GroupID   *uint   `json:"group_id"`
	Username  *string `json:"username" binding:"omitempty,min=3,max=32"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Nickname  *string `json:"nickname"`
	Phone     *string `json:"phone"`
	Gender    *int8   `json:"gender"`
	Signature *string `json:"signature"`
	Timezone  *string `json:"timezone"`
	Language  *string `json:"language"`
	RoleIDs   []uint  `json:"role_ids"`
	Status    *int8   `json:"status"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserService interface {
	// 用户资料管理（保持兼容性）
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)

	// 对应API接口的业务方法
	GetUsersList(page, pageSize int, filters repository.UserFilters) (*UserListResponse, error)
	GetUserDetail(id uint) (*UserDetailResponse, error)
	CreateUser(req CreateUserRequest) (*model.User, error)
	UpdateUser(id uint, req UpdateUserRequest) error
	DeleteUser(id uint) error
	ChangeUserPassword(id uint, req ChangePasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
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

func (s *userService) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// GetUsersList 获取用户列表（带筛选和分页）
func (s *userService) GetUsersList(page, pageSize int, filters repository.UserFilters) (*UserListResponse, error) {
	users, total, err := s.userRepo.GetUsersWithPagination(page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.convertToUserResponse(user)
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	pagination := Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &UserListResponse{
		Items:      userResponses,
		Pagination: pagination,
	}, nil
}

// GetUserDetail 获取用户详情
func (s *userService) GetUserDetail(id uint) (*UserDetailResponse, error) {
	user, err := s.userRepo.GetUserWithRoles(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	userResponse := s.convertToUserResponse(user)
	return &UserDetailResponse{UserResponse: userResponse}, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(req CreateUserRequest) (*model.User, error) {
	// 验证用户名是否已存在
	exists, err := s.userRepo.CheckUsernameExists(req.Username, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 验证邮箱是否已存在
	exists, err = s.userRepo.CheckEmailExists(req.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	language := req.Language
	if language == "" {
		language = "zh-CN"
	}

	status := req.Status
	if status == 0 {
		status = 1 // 默认启用
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

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// 如果指定了角色，需要关联角色（这里简化处理，实际需要操作user_roles表）
	// TODO: 实现角色关联逻辑

	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id uint, req UpdateUserRequest) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	updates := make(map[string]interface{})

	// 验证用户名唯一性
	if req.Username != nil {
		exists, err := s.userRepo.CheckUsernameExists(*req.Username, &id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("username already exists")
		}
		updates["username"] = *req.Username
	}

	// 验证邮箱唯一性
	if req.Email != nil {
		exists, err := s.userRepo.CheckEmailExists(*req.Email, &id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("email already exists")
		}
		updates["email"] = *req.Email
	}

	// 其他字段更新
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
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
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	return s.userRepo.UpdateUser(id, updates)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.DeleteUser(id)
}

// ChangeUserPassword 修改用户密码
func (s *userService) ChangeUserPassword(id uint, req ChangePasswordRequest) error {
	// 验证新密码确认
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("password confirmation does not match")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return s.userRepo.UpdatePassword(id, string(hashedPassword))
}

// convertToUserResponse 转换为用户响应格式
func (s *userService) convertToUserResponse(user *model.User) UserResponse {
	roles := make([]RoleInfo, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		}
	}

	return UserResponse{
		ID:          user.ID,
		GroupID:     user.GroupID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Phone:       user.Phone,
		Gender:      user.Gender,
		Signature:   user.Signature,
		Status:      user.Status,
		Timezone:    user.Timezone,
		Language:    user.Language,
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastLoginIP,
		Roles:       roles,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
