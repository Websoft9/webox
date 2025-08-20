package repository

import (
	"api-service/internal/model"

	"gorm.io/gorm"
)

// UserFilters 用户筛选条件
type UserFilters struct {
	Keyword string
	Status  *int8
	RoleID  *uint
	GroupID *uint
	Sort    string // created_at, updated_at, username
	Order   string // asc, desc
}

type UserRepository interface {
	// 现有方法保持不变
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.User, int64, error)

	// API设计中定义的查询方法
	GetUsersWithPagination(page, pageSize int, filters UserFilters) ([]*model.User, int64, error)
	GetUserByIDWithRoles(id uint) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(id uint, updates map[string]interface{}) error
	DeleteUser(id uint) error
	UpdatePassword(id uint, hashedPassword string) error
	SearchUsers(keyword string, page, pageSize int) ([]*model.User, int64, error)
	GetUserWithRoles(id uint) (*model.User, error)
	CheckUsernameExists(username string, excludeID *uint) (bool, error)
	CheckEmailExists(email string, excludeID *uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) List(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// GetUsersWithPagination 带筛选条件的分页查询
func (r *userRepository) GetUsersWithPagination(page, pageSize int, filters UserFilters) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{}).Preload("Group").Preload("Roles")

	// 应用筛选条件
	if filters.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+filters.Keyword+"%", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%")
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.GroupID != nil {
		query = query.Where("group_id = ?", *filters.GroupID)
	}

	if filters.RoleID != nil {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ?", *filters.RoleID)
	}

	// 应用排序
	orderBy := "created_at"
	if filters.Sort != "" {
		orderBy = filters.Sort
	}

	orderDir := "desc"
	if filters.Order != "" {
		orderDir = filters.Order
	}

	query = query.Order(orderBy + " " + orderDir)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// GetUserByIDWithRoles 根据ID获取用户，包含角色信息
func (r *userRepository) GetUserByIDWithRoles(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Group").Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUser 软删除用户
func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password_hash", hashedPassword).Error
}

// SearchUsers 搜索用户
func (r *userRepository) SearchUsers(keyword string, page, pageSize int) ([]*model.User, int64, error) {
	filters := UserFilters{
		Keyword: keyword,
		Sort:    "created_at",
		Order:   "desc",
	}
	return r.GetUsersWithPagination(page, pageSize, filters)
}

// GetUserWithRoles 获取用户及其角色信息
func (r *userRepository) GetUserWithRoles(id uint) (*model.User, error) {
	return r.GetUserByIDWithRoles(id)
}

// CheckUsernameExists 检查用户名是否存在
func (r *userRepository) CheckUsernameExists(username string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.User{}).Where("username = ?", username)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否存在
func (r *userRepository) CheckEmailExists(email string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.User{}).Where("email = ?", email)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
