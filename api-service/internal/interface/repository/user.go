package repository

import (
	"api-service/internal/model"
	"context"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByIDWithRelations(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error

	// 列表和搜索
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error)
	ListWithRelations(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error)

	// 业务查询
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error)
	ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint) (bool, error)
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	// 统计信息
	CountByStatus(ctx context.Context, status int) (int64, error)
	GetUserStats(ctx context.Context, userID uint) (*UserStats, error)
}

// UserStats 用户统计信息
type UserStats struct {
	LoginCount       int `json:"login_count"`
	ApplicationCount int `json:"application_count"`
	WorkflowCount    int `json:"workflow_count"`
}
