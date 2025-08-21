package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	GroupID     uint           `json:"group_id" gorm:"default:1"` // 用户组ID
	Username    string         `json:"username" gorm:"uniqueIndex;not null;size:50" binding:"required"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null;size:100" binding:"required,email"`
	Password    string         `json:"-" gorm:"not null;size:255"`
	Nickname    string         `json:"nickname" gorm:"size:100"` // 用户昵称
	Avatar      string         `json:"avatar" gorm:"size:255"`
	Phone       string         `json:"phone" gorm:"size:20"`                          // 手机号码
	Gender      int            `json:"gender" gorm:"default:0"`                       // 性别：1-男，2-女，0-未知
	Signature   string         `json:"signature" gorm:"size:500"`                     // 个性签名
	Status      int            `json:"status" gorm:"not null;default:1"`              // 状态：0-禁用，1-启用
	Timezone    string         `json:"timezone" gorm:"size:50;default:Asia/Shanghai"` // 时区
	Language    string         `json:"language" gorm:"size:10;default:zh-CN"`         // 语言
	LastLoginAt *time.Time     `json:"last_login_at"`                                 // 最后登录时间
	LastLoginIP string         `json:"last_login_ip" gorm:"size:45"`                  // 最后登录IP
	LoginCount  int            `json:"login_count" gorm:"default:0"`                  // 登录次数
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
	return u.Status == 1
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	// 注意：这里需要根据用户角色表来判断，暂时先返回false
	// 因为角色管理不在当前开发范围内
	return false
}

// GetDisplayName 获取用户显示名称
func (u *User) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
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

	// 暂时简化权限检查，后续根据角色表实现
	allowedActions := []string{
		"user:read", "user:update_self",
		"app:create", "app:read", "app:update_self", "app:delete_self",
	}
	for _, allowed := range allowedActions {
		if action == allowed {
			return true
		}
	}

	return false
}
