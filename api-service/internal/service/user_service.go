package service

import (
	"api-service/internal/model"
	"api-service/internal/repository"
	"api-service/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, error)
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
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
