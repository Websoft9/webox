package service

import (
	"api-service/internal/model"
	"api-service/internal/repository"
	"api-service/pkg/auth"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// UserQueryParams 用户查询参数
type UserQueryParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1"`
	RoleID   uint   `form:"role_id"`
	Sort     string `form:"sort"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	GroupID   uint     `json:"group_id" binding:"required"`
	Username  string   `json:"username" binding:"required,min=3,max=32"`
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=8,max=32"`
	Nickname  string   `json:"nickname"`
	Phone     string   `json:"phone"`
	Gender    int8     `json:"gender"`
	Signature string   `json:"signature"`
	Timezone  string   `json:"timezone"`
	Language  string   `json:"language"`
	RoleIDs   []uint   `json:"role_ids"`
	Status    int8     `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	GroupID   uint     `json:"group_id"`
	Nickname  string   `json:"nickname"`
	Email     string   `json:"email" binding:"omitempty,email"`
	Phone     string   `json:"phone"`
	Gender    int8     `json:"gender"`
	Signature string   `json:"signature"`
	Timezone  string   `json:"timezone"`
	Language  string   `json:"language"`
	RoleIDs   []uint   `json:"role_ids"`
	Status    int8     `json:"status"`
}

// UserListResult 用户列表结果
type UserListResult struct {
	Items      []*model.User `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type UserService interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, error)
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
	
	// 新增的用户管理接口
	ListUsersWithFilter(params *UserQueryParams) (*UserListResult, error)
	GetUserByID(userID uint) (*model.User, error)
	CreateUser(req *CreateUserRequest) (*model.User, error)
	UpdateUser(userID uint, req *UpdateUserRequest) (*model.User, error)
	DeleteUser(userID uint) error
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

func (s *userService) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// ListUsersWithFilter 带条件分页查询用户
func (s *userService) ListUsersWithFilter(params *UserQueryParams) (*UserListResult, error) {
	filter := repository.UserFilter{
		Keyword: params.Keyword,
		Status:  params.Status,
		RoleID:  params.RoleID,
		Sort:    params.Sort,
		Order:   params.Order,
	}

	offset := (params.Page - 1) * params.PageSize
	users, total, err := s.userRepo.ListWithFilter(offset, params.PageSize, &filter)
	if err != nil {
		return nil, err
	}

	// 计算分页信息
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	pagination := PaginationInfo{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}

	return &UserListResult{
		Items:      users,
		Pagination: pagination,
	}, nil
}

// GetUserByID 根据ID获取用户详情
func (s *userService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.GetByIDWithRelations(userID)
}

// CreateUser 创建用户
func (s *userService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// 验证用户名格式
	if !isValidUsername(req.Username) {
		return nil, errors.New("username must be 3-32 characters long and contain only letters, numbers, and underscores")
	}

	// 验证密码强度
	if !isValidPassword(req.Password) {
		return nil, errors.New("password must be 8-32 characters long and contain at least one uppercase letter, one lowercase letter, and one number")
	}

	// 验证手机号格式（如果提供）
	if req.Phone != "" && !isValidPhone(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}
	if req.Language == "" {
		req.Language = "zh-CN"
	}
	if req.Status == 0 {
		req.Status = 1 // 默认启用
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
		Status:       req.Status,
		Timezone:     req.Timezone,
		Language:     req.Language,
	}

	if err := s.userRepo.CreateWithRoles(user, req.RoleIDs); err != nil {
		return nil, err
	}

	// 返回创建的用户（包含关联数据）
	return s.userRepo.GetByIDWithRelations(user.ID)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(userID uint, req *UpdateUserRequest) (*model.User, error) {
	// 获取现有用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 如果更新邮箱，检查是否已存在
	if req.Email != "" && req.Email != user.Email {
		if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	// 验证手机号格式（如果提供）
	if req.Phone != "" && !isValidPhone(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// 更新字段
	if req.GroupID != 0 {
		user.GroupID = req.GroupID
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.Gender = req.Gender
	if req.Signature != "" {
		user.Signature = req.Signature
	}
	if req.Timezone != "" {
		user.Timezone = req.Timezone
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Status != 0 {
		user.Status = req.Status
	}

	if err := s.userRepo.UpdateWithRoles(user, req.RoleIDs); err != nil {
		return nil, err
	}

	// 返回更新后的用户（包含关联数据）
	return s.userRepo.GetByIDWithRelations(user.ID)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// TODO: 检查用户是否有关联的资源或正在执行的任务

	return s.userRepo.Delete(userID)
}

// 验证函数
func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}
	// 至少包含一个大写字母、一个小写字母和一个数字
	hasUpper, _ := regexp.MatchString("[A-Z]", password)
	hasLower, _ := regexp.MatchString("[a-z]", password)
	hasNumber, _ := regexp.MatchString("[0-9]", password)
	return hasUpper && hasLower && hasNumber
}

func isValidPhone(phone string) bool {
	// 简单的手机号验证，支持国内手机号
	matched, _ := regexp.MatchString("^1[3-9]\\d{9}$", phone)
	return matched
}
