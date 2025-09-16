package repository

import (
	"api-service/internal/model"
	"context"
)

// UserRepository defines user data access methods.
type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByIDWithRelations(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error

	// CreateUserRole creates a user-role association record.
	CreateUserRole(ctx context.Context, userRole *model.UserRole) error
	// GetRoleIDsByUserID returns all role IDs associated with the given user ID.
	GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
	// DeleteUserRole deletes the user-role association for the given user and role.
	DeleteUserRole(ctx context.Context, userID, roleID uint) error

	// List and search
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error)
	ListWithRelations(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error)

	// Business queries
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error)
	ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error)
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	// Statistics
	CountByStatus(ctx context.Context, status int) (int64, error)
	GetUserStats(ctx context.Context, userID uint) (*UserStats, error)
}

// UserStats user statistics information
type UserStats struct {
	LoginCount       int `json:"login_count"`
	ApplicationCount int `json:"application_count"`
	WorkflowCount    int
}
