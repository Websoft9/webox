package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// userRepository user data access implementation
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("status != ?", -1).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDWithRelations gets a user by ID (with related data)
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

// GetByUsername gets a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ? AND status != ?", username, -1).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ? AND status != ?", email, -1).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserRole creates a user-role association record.
func (r *userRepository) CreateUserRole(ctx context.Context, userRole *model.UserRole) error {
	if err := r.db.WithContext(ctx).Create(userRole).Error; err != nil {
		return err
	}
	return nil
}

// GetByUsernameOrEmail gets a user by username or email
func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").
		Where("(username = ? OR email = ?) AND status != ?", usernameOrEmail, usernameOrEmail, -1).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// soft delete: set status = -1 and update updated_at
	updates := map[string]interface{}{
		"status":     -1,
		"updated_at": time.Now(),
	}
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// List gets a list of users
func (r *userRepository) List(
	ctx context.Context,
	offset, limit int,
	filters map[string]interface{},
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)

	// Apply filters
	query = r.applyFilters(query, filters)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// ListWithRelations retrieves a list of users with their relationships loaded
func (r *userRepository) ListWithRelations(
	ctx context.Context,
	offset, limit int,
	filters map[string]interface{},
) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)

	// Handle role_id filter
	if roleID, ok := filters["role_id"]; ok && roleID != nil {
		// Join with user_roles table to filter users by role
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ? AND user_roles.status = ?", roleID, 1)
		// Remove role_id from filters to avoid applying it twice in applyFilters
		delete(filters, "role_id")
	}

	// Apply other filters
	query = r.applyFilters(query, filters)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with relationships
	err := query.Preload("Roles").Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// Search searches users
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status != ?", -1)
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// ExistsByUsername checks if the username exists
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ? AND status != ?", username, -1).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if the email exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ? AND status != ?", email, -1).Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcludeID checks if the username exists (excluding the specified ID)
func (r *userRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ? AND id != ? AND status != ?", username, excludeID, -1).
		Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcludeID checks if the email exists (excluding the specified ID)
func (r *userRepository) ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ? AND id != ? AND status != ?", email, excludeID, -1).
		Count(&count).Error
	return count > 0, err
}

// GetActiveUsers gets a list of active users
func (r *userRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", 1)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

// CountByStatus counts users by status
func (r *userRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetUserStats gets user statistics information
func (r *userRepository) GetUserStats(ctx context.Context, userID uint) (*repository.UserStats, error) {
	// Implement statistics logic according to actual business requirements
	// Currently returns default values, can add actual statistics queries later
	stats := &repository.UserStats{
		LoginCount:       0,
		ApplicationCount: 0,
		WorkflowCount:    0,
	}

	// Get login count (if there is a related table)
	// You can query the audit_logs table or other related tables to get actual statistics

	return stats, nil
}

// applyFilters applies query filters
func (r *userRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value != nil {
			query = query.Where(key+" = ?", value)
		}
	}
	return query
}

// ExistsByID checks if the ID exists
func (r *userRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ? AND status != ?", id, -1).Count(&count).Error
	return count > 0, err
}

// GetRoleIDsByUserID returns all role IDs associated with the given user ID.
func (r *userRepository) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("user_id = ? AND status = ?", userID, 1).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

// DeleteUserRole deletes the user-role association for the given user and role.
func (r *userRepository) DeleteUserRole(ctx context.Context, userID, roleID uint) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
	if err != nil {
		return err
	}
	return nil
}
