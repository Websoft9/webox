package repository

import (
	"api-service/internal/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.User, int64, error)
	
	// 新增方法
	GetByIDWithRoles(id uint) (*model.User, error)
	ListWithConditions(offset, limit int, conditions map[string]interface{}, sort, order string) ([]*model.User, int64, error)
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

func (r *userRepository) GetByIDWithRoles(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Preload("Group").First(&user, id).Error
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

	err := r.db.Preload("Roles").Preload("Group").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// ListWithConditions 根据条件查询用户列表
func (r *userRepository) ListWithConditions(offset, limit int, conditions map[string]interface{}, sort, order string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	// 构建查询条件
	if keyword, ok := conditions["keyword"].(string); ok && keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ? OR phone LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status, ok := conditions["status"].(int8); ok {
		query = query.Where("status = ?", status)
	}

	if roleID, ok := conditions["role_id"].(uint); ok {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ?", roleID)
	}

	// 排序
	if sort != "" {
		if order == "" {
			order = "desc"
		}
		orderClause := fmt.Sprintf("%s %s", sort, strings.ToUpper(order))
		query = query.Order(orderClause)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	err := query.Preload("Roles").Preload("Group").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// CreateWithRoles 创建用户并分配角色
func (r *userRepository) CreateWithRoles(user *model.User, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 分配角色
		if len(roleIDs) > 0 {
			var roles []model.Role
			if err := tx.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
				return err
			}
			if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateWithRoles 更新用户并更新角色
func (r *userRepository) UpdateWithRoles(user *model.User, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新用户基本信息
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		// 更新角色关联（如果提供了角色ID列表）
		if roleIDs != nil {
			// 清除现有角色
			if err := tx.Model(user).Association("Roles").Clear(); err != nil {
				return err
			}

			// 添加新角色
			if len(roleIDs) > 0 {
				var roles []model.Role
				if err := tx.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
					return err
				}
				if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
