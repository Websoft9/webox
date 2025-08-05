package repository

import (
	"api-service/internal/model"

	"gorm.io/gorm"
)

// UserFilter 用户查询过滤条件
type UserFilter struct {
	Keyword string
	Status  *int8
	RoleID  uint
	Sort    string
	Order   string
}

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByIDWithRelations(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.User, int64, error)
	ListWithFilter(offset, limit int, filter *UserFilter) ([]*model.User, int64, error)
	CreateWithRoles(user *model.User, roleIDs []uint) error
	UpdateWithRoles(user *model.User, roleIDs []uint) error
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

	err := r.db.Preload("Group").Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// GetByIDWithRelations 根据ID获取用户（包含关联数据）
func (r *userRepository) GetByIDWithRelations(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Group").Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListWithFilter 带条件的用户列表查询
func (r *userRepository) ListWithFilter(offset, limit int, filter *UserFilter) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	// 构建查询条件
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", keyword, keyword, keyword)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.RoleID != 0 {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ?", filter.RoleID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "created_at DESC"
	if filter.Sort != "" {
		orderBy = filter.Sort
		if filter.Order != "" {
			orderBy += " " + filter.Order
		}
	}
	query = query.Order(orderBy)

	// 分页查询
	err := query.Preload("Group").Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// CreateWithRoles 创建用户并分配角色
func (r *userRepository) CreateWithRoles(user *model.User, roleIDs []uint) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 创建用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 分配角色
	if len(roleIDs) > 0 {
		for _, roleID := range roleIDs {
			userRole := model.UserRole{
				UserID: user.ID,
				RoleID: roleID,
				Status: 1,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// UpdateWithRoles 更新用户并更新角色分配
func (r *userRepository) UpdateWithRoles(user *model.User, roleIDs []uint) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新用户基本信息
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果提供了角色ID列表，更新用户角色
	if roleIDs != nil {
		// 删除现有角色关联
		if err := tx.Where("user_id = ?", user.ID).Delete(&model.UserRole{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 添加新的角色关联
		for _, roleID := range roleIDs {
			userRole := model.UserRole{
				UserID: user.ID,
				RoleID: roleID,
				Status: 1,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}
