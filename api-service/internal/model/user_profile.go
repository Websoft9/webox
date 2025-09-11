// /opt/webox/api-service/internal/model/user_login_history.go

package model

import (
	"time"
)

// UserLoginHistory 用户登录历史记录模型
type UserLoginHistory struct {
	ID         uint       `json:"id" gorm:"primarykey"`
	UserID     uint       `json:"user_id" gorm:"column:user_id;not null;index"`
	IPAddress  string     `json:"ip_address" gorm:"column:ip_address;size:45"`
	UserAgent  string     `json:"user_agent" gorm:"column:user_agent;size:255"`
	Location   string     `json:"location" gorm:"column:location;size:100"`
	Device     string     `json:"device" gorm:"column:device;size:100"`
	Browser    string     `json:"browser" gorm:"column:browser;size:100"`
	LoginTime  time.Time  `json:"login_time" gorm:"column:login_time;not null;default:CURRENT_TIMESTAMP"`
	LogoutTime *time.Time `json:"logout_time" gorm:"column:logout_time"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (UserLoginHistory) TableName() string {
	return "user_login_history"
}
