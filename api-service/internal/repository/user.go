package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"fmt"

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
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "status":
				query = query.Where("status = ?", value)
			case "role":
				query = query.Where("role = ?", value)
			}
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// Search 搜索用户
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 构建搜索条件
	searchCondition := fmt.Sprintf("%%%s%%", keyword)
	query := r.db.WithContext(ctx).Model(&model.User{}).Where(
		"username LIKE ? OR email LIKE ? OR nickname LIKE ? OR phone LIKE ?",
		searchCondition, searchCondition, searchCondition, searchCondition,
	)

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
	err := query.Offset(offset).Limit(limit).Order("last_login_at desc").Find(&users).Error
	return users, total, err
}

// GetUsersByRole 根据角色获取用户列表 (移除此方法，因为角色管理不在当前范围)

// CountByStatus 根据状态统计用户数量
func (r *userRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetUserStats 获取用户统计信息
func (r *userRepository) GetUserStats(ctx context.Context, userID uint) (*repository.UserStats, error) {
	stats := &repository.UserStats{}

	// 获取用户基本信息
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}

	stats.LoginCount = user.LoginCount

	// 这里可以添加更多统计查询，比如：
	// - 应用数量统计
	// - 工作流数量统计
	// 由于这些表可能还没有创建，暂时设置为0
	stats.ApplicationCount = 0
	stats.WorkflowCount = 0

	return stats, nil
}
