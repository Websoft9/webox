package model

import (
	"time"
)

// User user model
type User struct {
	ID           uint       `json:"id" gorm:"primarykey"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null;size:64;comment:Username"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null;size:255;index:idx_email;comment:Email"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;not null;size:255;comment:Password hash"`
	Nickname     string     `json:"nickname" gorm:"size:64;comment:Nickname"`
	Avatar       string     `json:"avatar" gorm:"size:255;comment:Avatar URL"`
	Phone        string     `json:"phone" gorm:"size:20;comment:Phone number"`
	Gender       int        `json:"gender" gorm:"type:tinyint(1);default:0;comment:Gender: 0-unknown, 1-male, 2-female"` // 0:unknown, 1:male, 2:female
	Signature    string     `json:"signature" gorm:"size:255;comment:Personal signature"`
	Status       int        `json:"status" gorm:"type:tinyint(1);default:1;index:idx_user_status;comment:Status: -1-deleted, 0-disabled, 1-enabled"` // -1:deleted, 0:disabled, 1:enabled
	LastLoginAt  *time.Time `json:"last_login_at" gorm:"type:datetime;serializer:datetime;comment:Last login time"`
	LastLoginIP  string     `json:"last_login_ip" gorm:"size:45;comment:Last login IP"`
	Timezone     string     `json:"timezone" gorm:"size:64;default:UTC;comment:Timezone"`
	Language     string     `json:"language" gorm:"size:10;default:zh-CN;comment:Language"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`

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
