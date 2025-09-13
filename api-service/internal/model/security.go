package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// BaseModel base model
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role role model
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;size:64;not null" validate:"required,min=2,max=64"`
	Code        string `json:"code" gorm:"uniqueIndex;size:32;not null" validate:"required,min=2,max=32"`
	Description string `json:"description" gorm:"type:text" validate:"max=500"`
	IsSystem    bool   `json:"is_system" gorm:"default:false"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"` // -1:deleted, 0:disabled, 1:enabled

	// Associations
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles"`

	// Statistical fields (not mapped to database)
	PermissionCount int64 `json:"permission_count,omitempty" gorm:"-"`
	UserCount       int64 `json:"user_count,omitempty" gorm:"-"`
}

// TableName specify table name
func (Role) TableName() string {
	return "roles"
}

// IsActive check if role is enabled
func (r *Role) IsActive() bool {
	return r.Status == 1
}

// Permission permission model
type Permission struct {
	BaseModel
	ParentCode  string `json:"parent_code" gorm:"index;size:64;"`
	Scope       string `json:"scope" gorm:"size:64;not null" validate:"required,oneof=platform project"`
	Name        string `json:"name" gorm:"size:64;not null" validate:"required,min=2,max=64"`
	Code        string `json:"code" gorm:"uniqueIndex;size:64;not null" validate:"required,min=2,max=64"`
	Module      string `json:"module" gorm:"size:32;not null" validate:"required,min=2,max=32"`
	Action      string `json:"action" gorm:"size:32;not null" validate:"required,min=2,max=32"`
	Resource    string `json:"resource" gorm:"size:64" validate:"max=64"`
	Element     string `json:"element" gorm:"size:64" validate:"max=64"`
	Description string `json:"description" gorm:"type:text" validate:"max=500"`
	IsSystem    bool   `json:"is_system" gorm:"default:false"`
	IsMenu      bool   `json:"is_menu" gorm:"default:false"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"` // -1:deleted, 0:disabled, 1:enabled
	CreatedBy   *uint  `json:"created_by" gorm:"index"`
	UpdatedBy   *uint  `json:"updated_by" gorm:"index"`

	// Associations
	Parent   *Permission   `json:"parent,omitempty" gorm:"foreignKey:ParentCode"`
	Children []*Permission `json:"children,omitempty" gorm:"foreignKey:ParentCode"`
	Roles    []Role        `json:"roles,omitempty" gorm:"many2many:role_permissions"`

	// Statistical fields (not mapped to database)
	RoleCount int64 `json:"role_count,omitempty" gorm:"-"`
}

// TableName specify table name
func (Permission) TableName() string {
	return "permissions"
}

// IsActive check if permission is enabled
func (p *Permission) IsActive() bool {
	return p.Status == 1
}

// UserRole user role association table
type UserRole struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	RoleID    uint       `json:"role_id" gorm:"not null;index"`
	GrantedBy *uint      `json:"granted_by" gorm:"index"`
	GrantedAt time.Time  `json:"granted_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt *time.Time `json:"expires_at"`
	Status    int        `json:"status" gorm:"default:1"` //  -1:deleted, 0:disabled, 1:enabled
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Associations
	User Role `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName specify table name
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission role permission association table
type RolePermission struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	RoleID         uint      `json:"role_id" gorm:"not null;index"`
	PermissionCode string    `json:"permission_code" gorm:"not null;size:64;index"`
	GrantedBy      *uint     `json:"granted_by" gorm:"index"`
	GrantedAt      time.Time `json:"granted_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status         int       `json:"status" gorm:"default:1"` //  -1:deleted, 0:disabled, 1:enabled
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Associations
	Role       Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionCode"`
}

// TableName specify table name
func (RolePermission) TableName() string {
	return "role_permissions"
}

// APIToken API access token table
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

	// Associations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Display fields (not mapped to database)
	Username string `json:"username,omitempty" gorm:"-"`
}

// TableName specify table name
func (APIToken) TableName() string {
	return "api_tokens"
}

// IsExpired check if token is expired
func (t *APIToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// UserTwoFactor user two-factor authentication table
type UserTwoFactor struct {
	BaseModel
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Method      string     `json:"method" gorm:"size:32;not null"` // TOTP, EMAIL
	Secret      string     `json:"-" gorm:"size:255"`              // Encrypted storage
	BackupCodes JSON       `json:"-" gorm:"type:json"`             // Backup codes
	Email       string     `json:"email" gorm:"size:255"`
	Enabled     bool       `json:"enabled" gorm:"default:false"`
	VerifiedAt  *time.Time `json:"verified_at"`

	// Associations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specify table name
func (UserTwoFactor) TableName() string {
	return "user_two_factor"
}

// JSON custom JSON type
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

// AuthConfig authentication configuration - matches configs/auth.yaml structure
type AuthConfig struct {
	Version string `json:"version"`

	// API authentication configuration
	APIAuth struct {
		// Token authentication configuration
		TokenAuth struct {
			Algorithm        string `json:"algorithm"`
			Secret           string `json:"secret"`
			ExpiresIn        int    `json:"expires_in"`
			RefreshExpiresIn int    `json:"refresh_expires_in"`
			AutoRefresh      bool   `json:"auto_refresh"`
		} `json:"token_auth"`

		// OAuth2 API authentication configuration
		OAuth2 struct {
			Enabled           bool     `json:"enabled"`
			DefaultScopes     []string `json:"default_scopes"`
			TokenEndpoint     string   `json:"token_endpoint"`
			AuthorizeEndpoint string   `json:"authorize_endpoint"`
		} `json:"oauth2"`
	} `json:"api_auth"`

	// User authentication configuration
	UserAuth struct {
		// Basic authentication configuration
		BasicAuth struct {
			LoginMethods []string `json:"login_methods"`
		} `json:"basic_auth"`

		// Email authentication configuration
		EmailAuth struct {
			Enabled   bool `json:"enabled"`
			ExpiresIn int  `json:"expires_in"`
		} `json:"email_auth"`

		// Password policy configuration
		PasswordPolicy PasswordPolicy `json:"password_policy"`

		// Login security configuration
		LoginSecurity LoginSecurity `json:"login_security"`

		// OAuth2 login configuration
		OAuth2 struct {
			Enabled      bool                            `json:"enabled"`
			AutoRegister bool                            `json:"auto_register"`
			DefaultRole  string                          `json:"default_role"`
			Providers    map[string]OAuth2ProviderConfig `json:"providers"`
		} `json:"oauth2"`

		// Two-factor authentication configuration
		TwoFactor struct {
			Enabled       bool     `json:"enabled"`
			RequiredRoles []string `json:"required_roles"`
			Methods       struct {
				// TOTP authentication configuration
				TOTP struct {
					Enabled          bool   `json:"enabled"`
					Issuer           string `json:"issuer"`
					Algorithm        string `json:"algorithm"`
					Digits           int    `json:"digits"`
					Period           int    `json:"period"`
					BackupCodesCount int    `json:"backup_codes_count"`
				} `json:"totp"`

				// Email verification code configuration
				Email struct {
					Enabled    bool   `json:"enabled"`
					CodeLength int    `json:"code_length"`
					ExpiresIn  int    `json:"expires_in"`
					RateLimit  int    `json:"rate_limit"`
					Template   string `json:"template"`
				} `json:"email"`
			} `json:"methods"`
		} `json:"two_factor"`
	} `json:"user_auth"`

	// Session configuration
	SessionConfig struct {
		Timeout               int  `json:"timeout"`
		MaxConcurrentSessions int  `json:"max_concurrent_sessions"`
		RememberMeEnabled     bool `json:"remember_me_enabled"`
		RememberMeDuration    int  `json:"remember_me_duration"`
	} `json:"session_config"`
}

// OAuth2ProviderConfig OAuth2 provider configuration - matches configs/auth.yaml structure
type OAuth2ProviderConfig struct {
	Name         string            `json:"name"`
	Enabled      bool              `json:"enabled"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	AuthorizeURL string            `json:"authorize_url"`
	TokenURL     string            `json:"token_url"`
	UserInfoURL  string            `json:"user_info_url"`
	AutoRegister bool              `json:"auto_register"`
	UserMapping  map[string]string `json:"user_mapping"`
	SortOrder    int               `json:"sort_order"`
}

// PasswordPolicy password policy - matches configs/auth.yaml structure
type PasswordPolicy struct {
	MinLength           int  `json:"min_length"`
	MaxLength           int  `json:"max_length"`
	RequireUppercase    bool `json:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers"`
	RequireSymbols      bool `json:"require_symbols"`
	PasswordExpiresDays int  `json:"password_expires_days"`
}

// LoginSecurity login security configuration - matches configs/auth.yaml structure
type LoginSecurity struct {
	MaxLoginAttempts     int      `json:"max_login_attempts"`
	LockoutDuration      int      `json:"lockout_duration"`
	IPWhitelistEnabled   bool     `json:"ip_whitelist_enabled"`
	IPWhitelist          []string `json:"ip_whitelist"`
	LoginTimeRestriction bool     `json:"login_time_restriction"`
	AllowedLoginHours    string   `json:"allowed_login_hours"`
}
