package model

import (
	"time"
)

// UserLoginHistory User login history model
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

	// Association
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserProfile User profile settings table
type UserProfile struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID       uint      `gorm:"column:user_id;not null" json:"user_id"`
	Category     string    `gorm:"column:category;default:general" json:"category"`
	ConfigKey    string    `gorm:"column:config_key;not null;uniqueIndex" json:"config_key"`
	ConfigValue  string    `gorm:"column:config_value;type:text" json:"config_value"`
	Description  string    `gorm:"column:description;type:text" json:"description"`
	IsReadonly   bool      `gorm:"column:is_readonly;default:0" json:"is_readonly"`
	IsEncrypted  bool      `gorm:"column:is_encrypted;default:0" json:"is_encrypted"`
	DefaultValue string    `gorm:"column:default_value;type:text" json:"default_value"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name
func (UserProfile) TableName() string {
	return "user_profile"
}

// TableName specifies the table name
func (UserLoginHistory) TableName() string {
	return "user_login_history"
}
