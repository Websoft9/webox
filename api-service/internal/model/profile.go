package model

import (
	"time"

	"gorm.io/gorm"
)

// UserProfile 用户个人中心配置表
type UserProfile struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index;uniqueIndex:idx_user_config,priority:1"`
	ConfigKey   string         `json:"config_key" gorm:"size:64;not null;uniqueIndex:idx_user_config,priority:2"`
	ConfigValue string         `json:"config_value" gorm:"type:text"`
	Description string         `json:"description" gorm:"size:255"`
	Category    string         `json:"category" gorm:"size:32;default:'general'"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联字段
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定用户配置表名
func (UserProfile) TableName() string {
	return "user_profile"
}

// UserLoginHistory 用户登录历史表
type UserLoginHistory struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	IPAddress  string         `json:"ip_address" gorm:"size:45;not null"`
	UserAgent  string         `json:"user_agent" gorm:"type:text"`
	Location   string         `json:"location" gorm:"size:100"`
	Device     string         `json:"device" gorm:"size:100"`
	Browser    string         `json:"browser" gorm:"size:100"`
	OS         string         `json:"os" gorm:"size:100"`
	Status     string         `json:"status" gorm:"size:20;default:'SUCCESS'"` // SUCCESS, FAILED, LOGOUT
	LoginTime  time.Time      `json:"login_time"`
	LogoutTime *time.Time     `json:"logout_time"`
	SessionID  string         `json:"session_id" gorm:"size:128"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联字段
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定登录历史表名
func (UserLoginHistory) TableName() string {
	return "user_login_history"
}

// IsActive 检查会话是否还在活跃状态
func (ulh *UserLoginHistory) IsActive() bool {
	return ulh.Status == "SUCCESS" && ulh.LogoutTime == nil
}

// SessionDuration 获取会话持续时间
func (ulh *UserLoginHistory) SessionDuration() time.Duration {
	if ulh.LogoutTime != nil {
		return ulh.LogoutTime.Sub(ulh.LoginTime)
	}
	return time.Since(ulh.LoginTime)
}
