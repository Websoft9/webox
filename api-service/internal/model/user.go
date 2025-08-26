package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	GroupID           uint           `json:"group_id" gorm:"not null"`
	Username          string         `json:"username" gorm:"uniqueIndex;not null;size:64"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash      string         `json:"-" gorm:"column:password_hash;not null;size:255"`
	Password          string         `json:"-" gorm:"-"` // 临时字段，不映射到数据库
	Nickname          string         `json:"nickname" gorm:"size:64"`
	Avatar            string         `json:"avatar" gorm:"size:255"`
	Phone             string         `json:"phone" gorm:"size:20"`
	Gender            int            `json:"gender" gorm:"default:0"` // 0:未知, 1:男, 2:女
	Signature         string         `json:"signature" gorm:"size:255"`
	Status            int            `json:"status" gorm:"default:1"` // 1:active, 0:inactive
	LastLoginAt       *time.Time     `json:"last_login_at"`
	LastLoginIP       string         `json:"last_login_ip" gorm:"size:45"`
	Timezone          string         `json:"timezone" gorm:"size:64;default:UTC"`
	Language          string         `json:"language" gorm:"size:10;default:zh-CN"`
	TwoFactorEnabled  bool           `json:"two_factor_enabled" gorm:"default:false"`
	PasswordChangedAt *time.Time     `json:"password_changed_at"`
	LoginAttempts     int            `json:"login_attempts" gorm:"default:0"`
	LockedUntil       *time.Time     `json:"locked_until"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联字段 (不会直接映射到数据库，需要时通过 Preload 加载)
	Group      *UserGroup      `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	Roles      []Role          `json:"roles,omitempty" gorm:"many2many:user_roles"`
	APITokens  []APIToken      `json:"api_tokens,omitempty" gorm:"foreignKey:UserID"`
	TwoFactors []UserTwoFactor `json:"two_factors,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 检查用户是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == 1
}

// GetDisplayName 获取用户显示名称
func (u *User) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

// UserGroup 用户组模型
type UserGroup struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null;size:64"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null;size:32"`
	Description string    `json:"description" gorm:"type:text"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定用户组表名
func (UserGroup) TableName() string {
	return "user_groups"
}

// 移除重复的 Role 定义，使用 security.go 中的定义
