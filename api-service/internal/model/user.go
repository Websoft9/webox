package model

import (
	"time"
)

// User user model
type User struct {
	ID           uint       `json:"id" gorm:"primarykey"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null;size:64"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;not null;size:255"`
	Nickname     string     `json:"nickname" gorm:"size:64"`
	Avatar       string     `json:"avatar" gorm:"size:255"`
	Phone        string     `json:"phone" gorm:"size:20"`
	Gender       int        `json:"gender" gorm:"default:0"` // 0:unknown, 1:male, 2:female
	Signature    string     `json:"signature" gorm:"size:255"`
	Status       int        `json:"status" gorm:"default:1"` // 1:active, 0:inactive
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `json:"last_login_ip" gorm:"size:45"`
	Timezone     string     `json:"timezone" gorm:"size:64;default:UTC"`
	Language     string     `json:"language" gorm:"size:10;default:zh-CN"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Association fields (not directly mapped to the database, loaded via Preload when needed)
	Roles      []Role          `json:"roles,omitempty" gorm:"many2many:user_roles"`
	APITokens  []APIToken      `json:"api_tokens,omitempty" gorm:"foreignKey:UserID"`
	TwoFactors []UserTwoFactor `json:"two_factors,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == 1
}

// GetDisplayName gets the user's display name
func (u *User) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}
