package service

import (
	"api-service/internal/model"
	"api-service/internal/repository"
	"api-service/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	GroupID  uint   `json:"group_id"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Gender   int8   `json:"gender"`
	Avatar   string `json:"avatar"`
	GroupID  uint   `json:"group_id"`
	Status   int8   `json:"status"`
}

type UserService interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, error)
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ListUsers(page, pageSize int, keyword, status, groupID string) ([]*model.User, int64, error)
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

func (s *userService) ListUsers(page, pageSize int, keyword, status, groupID string) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.ListWithFilter(offset, pageSize, keyword, status, groupID)
}

func (s *userService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		GroupID:      req.GroupID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Gender:       req.Gender,
		Status:       1, // 默认启用
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(userID uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户名重复（排除自己）
	if req.Username != "" && req.Username != user.Username {
		if existingUser, err := s.userRepo.GetByUsername(req.Username); err == nil && existingUser.ID != userID {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}

	// 检查邮箱重复（排除自己）
	if req.Email != "" && req.Email != user.Email {
		if existingUser, err := s.userRepo.GetByEmail(req.Email); err == nil && existingUser.ID != userID {
			return nil, errors.New("邮箱已存在")
		}
		user.Email = req.Email
	}

	// 更新其他字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.GroupID != 0 {
		user.GroupID = req.GroupID
	}
	if req.Status >= 0 {
		user.Status = req.Status
	}
	user.Gender = req.Gender

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	return s.userRepo.Delete(userID)
}
