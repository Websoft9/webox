package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;size:64;not null" validate:"required,min=2,max=64"`
	Code        string `json:"code" gorm:"uniqueIndex;size:32;not null" validate:"required,min=2,max=32"`
	Description string `json:"description" gorm:"type:text" validate:"max=500"`
	IsSystem    bool   `json:"is_system" gorm:"default:false"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"` // -1:删除, 0:禁用, 1:启用
	CreatedBy   *uint  `json:"created_by" gorm:"index"`
	UpdatedBy   *uint  `json:"updated_by" gorm:"index"`

	// 关联关系
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles"`

	// 统计字段 (不映射到数据库)
	PermissionCount int64 `json:"permission_count,omitempty" gorm:"-"`
	UserCount       int64 `json:"user_count,omitempty" gorm:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// IsActive 检查角色是否启用
func (r *Role) IsActive() bool {
	return r.Status == 1
}

// Permission 权限模型
type Permission struct {
	BaseModel
	ParentID    *uint  `json:"parent_id" gorm:"index"`
	Scope       string `json:"scope" gorm:"size:64;not null" validate:"required,oneof=platform project"`
	Name        string `json:"name" gorm:"size:64;not null" validate:"required,min=2,max=64"`
	Code        string `json:"code" gorm:"uniqueIndex;size:64;not null" validate:"required,min=2,max=64"`
	Module      string `json:"module" gorm:"size:32;not null" validate:"required,min=2,max=32"`
	Action      string `json:"action" gorm:"size:32;not null" validate:"required,min=2,max=32"`
	Resource    string `json:"resource" gorm:"size:256" validate:"max=64"`
	Description string `json:"description" gorm:"type:text" validate:"max=500"`
	IsSystem    bool   `json:"is_system" gorm:"default:false"`
	IsMenu      bool   `json:"is_menu" gorm:"default:false"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"` // -1:删除, 0:禁用, 1:启用
	CreatedBy   *uint  `json:"created_by" gorm:"index"`
	UpdatedBy   *uint  `json:"updated_by" gorm:"index"`

	// 关联关系
	Parent   *Permission  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Permission `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Roles    []Role       `json:"roles,omitempty" gorm:"many2many:role_permissions"`

	// 统计字段 (不映射到数据库)
	RoleCount int64 `json:"role_count,omitempty" gorm:"-"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// IsActive 检查权限是否启用
func (p *Permission) IsActive() bool {
	return p.Status == 1
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	RoleID    uint       `json:"role_id" gorm:"not null;index"`
	GrantedBy *uint      `json:"granted_by" gorm:"index"`
	GrantedAt time.Time  `json:"granted_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt *time.Time `json:"expires_at"`
	Status    int        `json:"status" gorm:"default:1"` //  -1:删除, 0:禁用, 1:启用
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// 关联关系
	User Role `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	RoleID       uint      `json:"role_id" gorm:"not null;index"`
	PermissionID uint      `json:"permission_id" gorm:"not null;index"`
	GrantedBy    *uint     `json:"granted_by" gorm:"index"`
	GrantedAt    time.Time `json:"granted_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status       int       `json:"status" gorm:"default:1"` //  -1:删除, 0:禁用, 1:启用
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联关系
	Role       Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// APIToken API访问令牌表
type APIToken struct {
	BaseModel
	Name        string     `json:"name" gorm:"size:64;not null" validate:"required,min=2,max=64"`
	Token       string     `json:"token" gorm:"uniqueIndex;size:255;not null"`
	TokenHash   string     `json:"-" gorm:"uniqueIndex;size:64;not null"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Scopes      JSON       `json:"scopes" gorm:"type:json"`
	Description string     `json:"description" gorm:"type:text" validate:"max=500"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	LastUsedIP  string     `json:"last_used_ip" gorm:"size:45"`
	ExpiresAt   *time.Time `json:"expires_at"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// 显示字段 (不映射到数据库)
	Username string `json:"username,omitempty" gorm:"-"`
}

// TableName 指定表名
func (APIToken) TableName() string {
	return "api_tokens"
}

// IsExpired 检查令牌是否过期
func (t *APIToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// UserTwoFactor 用户双因子认证表
type UserTwoFactor struct {
	BaseModel
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Method      string     `json:"method" gorm:"size:32;not null"` // TOTP, EMAIL
	Secret      string     `json:"-" gorm:"size:255"`              // 加密存储
	BackupCodes JSON       `json:"-" gorm:"type:json"`             // 备用码
	Email       string     `json:"email" gorm:"size:255"`
	Enabled     bool       `json:"enabled" gorm:"default:false"`
	VerifiedAt  *time.Time `json:"verified_at"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (UserTwoFactor) TableName() string {
	return "user_two_factor"
}

// JSON 自定义JSON类型
type JSON map[string]interface{}

// Scan implements sql.Scanner interface for GORM
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot convert %T to JSON", value)
	}

	if len(bytes) == 0 {
		*j = make(map[string]interface{})
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface for GORM
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// AuthConfig 认证配置
type AuthConfig struct {
	APIAuth struct {
		TokenAuthEnabled bool `json:"token_auth_enabled"`
		OAuth2Enabled    bool `json:"oauth2_enabled"`
		JWTConfig        struct {
			Algorithm        string `json:"algorithm"`
			ExpiresIn        int    `json:"expires_in"`
			RefreshExpiresIn int    `json:"refresh_expires_in"`
			AutoRefresh      bool   `json:"auto_refresh"`
		} `json:"jwt_config"`
	} `json:"api_auth"`

	UserAuth struct {
		OAuth2Enabled          bool             `json:"oauth2_enabled"`
		OAuth2Providers        []OAuth2Provider `json:"oauth2_providers"`
		TwoFactorEnabled       bool             `json:"two_factor_enabled"`
		TwoFactorMethods       []string         `json:"two_factor_methods"`
		TwoFactorRequiredRoles []string         `json:"two_factor_required_roles"`
		PasswordPolicy         PasswordPolicy   `json:"password_policy"`
		LoginSecurity          LoginSecurity    `json:"login_security"`
	} `json:"user_auth"`

	SessionConfig struct {
		Timeout               int  `json:"timeout"`
		MaxConcurrentSessions int  `json:"max_concurrent_sessions"`
		RememberMeEnabled     bool `json:"remember_me_enabled"`
		RememberMeDuration    int  `json:"remember_me_duration"`
	} `json:"session_config"`
}

// OAuth2Provider OAuth2提供商配置
type OAuth2Provider struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	Enabled      bool              `json:"enabled"`
	AutoRegister bool              `json:"auto_register"`
	UserMapping  map[string]string `json:"user_mapping"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength           int  `json:"min_length"`
	MaxLength           int  `json:"max_length"`
	RequireUppercase    bool `json:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers"`
	RequireSymbols      bool `json:"require_symbols"`
	PasswordHistory     int  `json:"password_history"`
	PasswordExpiresDays int  `json:"password_expires_days"`
}

// LoginSecurity 登录安全配置
type LoginSecurity struct {
	MaxLoginAttempts     int      `json:"max_login_attempts"`
	LockoutDuration      int      `json:"lockout_duration"`
	IPWhitelistEnabled   bool     `json:"ip_whitelist_enabled"`
	IPWhitelist          []string `json:"ip_whitelist"`
	LoginTimeRestriction bool     `json:"login_time_restriction"`
	AllowedLoginHours    string   `json:"allowed_login_hours"`
}
