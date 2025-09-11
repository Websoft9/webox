package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"time"

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
	err := r.db.WithContext(ctx).Where("status != ?", -1).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDWithRelations 根据ID获取用户（包含关联数据）
func (r *userRepository) GetByIDWithRelations(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("status != ?", -1).
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
	err := r.db.WithContext(ctx).Where("username = ? AND status != ?", username, -1).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ? AND status != ?", email, -1).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("(username = ? OR email = ?) AND status != ?", usernameOrEmail, usernameOrEmail, -1).
		First(&user).Error
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
	// soft delete: set status = -1 and update updated_at
	updates := map[string]interface{}{
		"status":     -1,
		"updated_at": time.Now(),
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// List 获取用户列表
func (r *userRepository) List(
	ctx context.Context,
	offset, limit int,
	filters map[string]interface{},
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)

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

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)

	// 应用过滤器
	query = r.applyFilters(query, filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Preload("Roles").Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// Search 搜索用户
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)
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
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ? AND status != ?", username, -1).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ? AND status != ?", email, -1).Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *userRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ? AND id != ? AND status != ?", username, excludeID, -1).
		Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcludeID 检查邮箱是否存在（排除指定ID）
func (r *userRepository) ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ? AND id != ? AND status != ?", email, excludeID, -1).
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
		if value != nil {
			query = query.Where(key+" = ?", value)
		}
	}
	return query
}

// ExistsByID 检查ID是否存在
func (r *userRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ? AND status != ?", id, -1).Count(&count).Error
	return count > 0, err
}
