package model

import (
	"time"
)

// UserLoginHistory User login history model
type UserLoginHistory struct {
	ID         uint       `json:"id" gorm:"primarykey;type:bigint unsigned"`
	UserID     uint       `json:"user_id" gorm:"column:user_id;type:bigint unsigned;not null;index:idx_log_user_id;comment:User ID"`
	IPAddress  string     `json:"ip_address" gorm:"column:ip_address;size:45;comment:IP address"`
	UserAgent  string     `json:"user_agent" gorm:"column:user_agent;size:255;comment:User agent"`
	Location   string     `json:"location" gorm:"column:location;size:100;comment:Login location"`
	Device     string     `json:"device" gorm:"column:device;size:100;comment:Device info"`
	Browser    string     `json:"browser" gorm:"column:browser;size:100;comment:Browser info"`
	LoginTime  time.Time  `json:"login_time" gorm:"column:login_time;not null;default:CURRENT_TIMESTAMP;type:datetime;serializer:datetime;index:idx_login_time;comment:Login time"`
	LogoutTime *time.Time `json:"logout_time" gorm:"column:logout_time;type:datetime;serializer:datetime;comment:Logout time"`
	Status     string     `json:"status" gorm:"type:varchar(20);default:'ACTIVE';index:idx_log_status;comment:Session status"`
	SessionID  string     `json:"session_id" gorm:"size:128;comment:Session ID"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`

	// Association
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// UserProfile User profile settings table
type UserProfile struct {
	ID           uint      `gorm:"primaryKey;column:id;type:bigint unsigned" json:"id"`
	UserID       uint      `gorm:"column:user_id;type:bigint unsigned;not null;uniqueIndex:uk_user_config_key;comment:User ID" json:"user_id"`
	Category     string    `gorm:"column:category;size:64;default:general;comment:Profile category" json:"category"`
	ConfigKey    string    `gorm:"column:config_key;size:64;not null;uniqueIndex:uk_user_config_key;comment:Configuration key" json:"config_key"`
	ConfigValue  string    `gorm:"column:config_value;type:text;comment:Configuration value" json:"config_value"`
	Description  string    `gorm:"column:description;type:text;comment:Configuration description" json:"description"`
	IsReadonly   int8      `gorm:"column:is_readonly;type:tinyint(1);default:0;comment:Whether read-only" json:"is_readonly"`
	IsEncrypted  int8      `gorm:"column:is_encrypted;type:tinyint(1);default:0;comment:Whether encrypted" json:"is_encrypted"`
	DefaultValue string    `gorm:"column:default_value;type:text;comment:Default value" json:"default_value"`
	SortOrder    int       `gorm:"column:sort_order;default:0;index:idx_profile_sort_order;comment:Sort order" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName specifies the table name
func (UserProfile) TableName() string {
	return "user_profile"
}

// GetEffectiveValue returns the effective value for the user profile
// If ConfigValue is empty, returns DefaultValue
func (up *UserProfile) GetEffectiveValue() string {
	if up.ConfigValue != "" {
		return up.ConfigValue
	}
	return up.DefaultValue
}

// TableName specifies the table name
func (UserLoginHistory) TableName() string {
	return "user_login_history"
}
