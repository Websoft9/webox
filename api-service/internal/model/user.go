package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint       `json:"id" gorm:"primarykey"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null;size:64"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;not null;size:255"`
	Nickname     string     `json:"nickname" gorm:"size:64"`
	Avatar       string     `json:"avatar" gorm:"size:255"`
	Phone        string     `json:"phone" gorm:"size:20"`
	Gender       int        `json:"gender" gorm:"default:0"` // 0:未知, 1:男, 2:女
	Signature    string     `json:"signature" gorm:"size:255"`
	Status       int        `json:"status" gorm:"default:1"` // 1:active, 0:inactive
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `json:"last_login_ip" gorm:"size:45"`
	Timezone     string     `json:"timezone" gorm:"size:64;default:UTC"`
	Language     string     `json:"language" gorm:"size:10;default:zh-CN"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联字段 (不会直接映射到数据库，需要时通过 Preload 加载)
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

// 移除重复的 Role 定义，使用 security.go 中的定义
