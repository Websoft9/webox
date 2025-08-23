package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// userRepository 用户数据访问实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建新的用户Repository实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDWithRelations 根据ID获取用户（包含关联数据）
func (r *userRepository) GetByIDWithRelations(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Group").
		Preload("Roles").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(
	ctx context.Context,
	offset, limit int,
	filters map[string]interface{},
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 应用过滤器
	query = r.applyFilters(query, filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// ListWithRelations 获取用户列表（包含关联数据）
func (r *userRepository) ListWithRelations(
	ctx context.Context,
	offset, limit int,
	filters map[string]interface{},
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 应用过滤器
	query = r.applyFilters(query, filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据（包含关联）
	err := query.
		Preload("Group").
		Preload("Roles").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&users).Error
	return users, total, err
}

// Search 搜索用户
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *userRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ? AND id != ?", username, excludeID).
		Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcludeID 检查邮箱是否存在（排除指定ID）
func (r *userRepository) ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ? AND id != ?", email, excludeID).
		Count(&count).Error
	return count > 0, err
}

// GetActiveUsers 获取活跃用户列表
func (r *userRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", 1)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// GetUsersByGroupID 根据用户组ID获取用户列表
func (r *userRepository) GetUsersByGroupID(
	ctx context.Context, groupID uint, offset, limit int,
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("group_id = ?", groupID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// CountByStatus 根据状态统计用户数
func (r *userRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetUserStats 获取用户统计信息
func (r *userRepository) GetUserStats(ctx context.Context, userID uint) (*repository.UserStats, error) {
	// 这里需要根据实际的业务逻辑来实现统计
	// 目前返回默认值，后续可以添加实际的统计查询
	stats := &repository.UserStats{
		LoginCount:       0,
		ApplicationCount: 0,
		WorkflowCount:    0,
	}

	// 获取登录次数（如果有相关表的话）
	// 可以通过查询 audit_logs 表或其他相关表来获取实际的统计数据

	return stats, nil
}

// applyFilters 应用查询过滤器
func (r *userRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "status":
				query = query.Where("status = ?", value)
			case "group_id":
				query = query.Where("group_id = ?", value)
			case "gender":
				query = query.Where("gender = ?", value)
			case "language":
				query = query.Where("language = ?", value)
			case "keyword":
				if keyword, ok := value.(string); ok && keyword != "" {
					searchPattern := "%" + keyword + "%"
					query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
						searchPattern, searchPattern, searchPattern)
				}
			}
		}
	}
	return query
}

// ExistsByID 检查ID是否存在
func (r *userRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
