package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Username    string         `json:"username" gorm:"uniqueIndex;not null;size:50" binding:"required"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null;size:100" binding:"required,email"`
	Password    string         `json:"-" gorm:"not null;size:255"`
	FirstName   string         `json:"first_name" gorm:"size:50"`
	LastName    string         `json:"last_name" gorm:"size:50"`
	Avatar      string         `json:"avatar" gorm:"size:255"`
	Status      string         `json:"status" gorm:"not null;default:active;size:20"` // active, inactive, banned
	Role        string         `json:"role" gorm:"not null;default:user;size:20"`     // admin, user, guest
	LastLoginAt *time.Time     `json:"last_login_at"`
	LoginCount  int            `json:"login_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 检查用户是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// GetFullName 获取用户全名
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// CanPerformAction 检查用户是否可以执行特定操作
func (u *User) CanPerformAction(action string) bool {
	if !u.IsActive() {
		return false
	}

	// 管理员可以执行所有操作
	if u.IsAdmin() {
		return true
	}

	// 根据角色检查权限
	switch u.Role {
	case "user":
		// 普通用户权限
		allowedActions := []string{
			"user:read", "user:update_self",
			"app:create", "app:read", "app:update_self", "app:delete_self",
		}
		for _, allowed := range allowedActions {
			if action == allowed {
				return true
			}
		}
	case "guest":
		// 访客权限
		allowedActions := []string{"user:read", "app:read"}
		for _, allowed := range allowedActions {
			if action == allowed {
				return true
			}
		}
	}

	return false
}
