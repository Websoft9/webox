package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// SecretKeyType represents the type of secret key
type SecretKeyType string

const (
	SecretKeyTypeText SecretKeyType = "TEXT"
	SecretKeyTypeAccount   SecretKeyType = "ACCOUNT"
	SecretKeyTypeFile      SecretKeyType = "FILE"
)

// ValidSecretKeyTypes returns all valid secret key types
func ValidSecretKeyTypes() []SecretKeyType {
	return []SecretKeyType{
		SecretKeyTypeText,
		SecretKeyTypeAccount,
		SecretKeyTypeFile,
	}
}

// IsValidSecretKeyType checks if a given type is valid
func IsValidSecretKeyType(keyType SecretKeyType) bool {
	validTypes := ValidSecretKeyTypes()
	for _, t := range validTypes {
		if t == keyType {
			return true
		}
	}
	return false
}

// SecretFields represents JSON custom fields
type SecretFields map[string]interface{}

// Value implements driver.Valuer interface for database storage
func (cf SecretFields) Value() (driver.Value, error) {
	if cf == nil {
		return nil, nil
	}
	return json.Marshal(cf)
}

// Scan implements sql.Scanner interface for database retrieval
func (cf *SecretFields) Scan(value interface{}) error {
	if value == nil {
		*cf = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, cf)
}

// SecretKey represents a secret key in the system
type SecretKey struct {
	ID              uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string        `json:"name" gorm:"size:64;not null"`
	KeyType         SecretKeyType `json:"key_type" gorm:"size:20;not null"`
	Description     *string       `json:"description" gorm:"type:text"`
	SecretFields    SecretFields  `json:"secret_fields" gorm:"type:text"`
	ExpiresAt       *time.Time    `json:"expires_at" gorm:"type:datetime;serializer:datetime"`
	ResourceGroupID *uint         `json:"resource_group_id"`
	OwnerID         uint          `json:"owner_id" gorm:"not null"`
	CreatedAt       time.Time     `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`

	// Relationships
	Owner *User `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
}

// TableName returns the table name for GORM
func (SecretKey) TableName() string {
	return "secret_keys"
}

// SecretReference represents secret key reference record
type SecretReference struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	SecretID     uint      `gorm:"not null;index:idx_secret_references_secret_id;uniqueIndex:idx_secret_references_secret_resource" json:"secret_id"`
	ResourceCode string    `gorm:"type:varchar(64);not null;index:idx_secret_references_resource_code;uniqueIndex:idx_secret_references_secret_resource" json:"resource_code"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`

	// Associations
	SecretKey SecretKey `gorm:"foreignKey:SecretID" json:"-"`
}

// UserSecret represents the relationship between users and secret keys
type UserSecret struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	SecretKeyID uint       `gorm:"not null;index" json:"secret_key_id"`
	GrantedBy   *uint      `gorm:"index" json:"granted_by,omitempty"`
	GrantedAt   time.Time  `json:"granted_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;serializer:datetime"`
	ExpiresAt   *time.Time `json:"expires_at" gorm:"type:datetime;serializer:datetime"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`

	User          User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SecretKey     SecretKey `gorm:"foreignKey:SecretKeyID;references:ID" json:"secret_key,omitempty"`
	GrantedByUser *User     `gorm:"foreignKey:GrantedBy;references:ID" json:"granted_by_user,omitempty"`
}

// TableName returns the table name for GORM
func (UserSecret) TableName() string {
	return "user_secret"
}

// BeforeCreate hook for GORM
func (s *SecretKey) BeforeCreate(tx *gorm.DB) error {
	if s.Name == "" {
		return errors.New("secret key name is required")
	}
	if s.KeyType == "" {
		return errors.New("secret key type is required")
	}
	if !IsValidSecretKeyType(s.KeyType) {
		return errors.New("invalid secret key type")
	}
	if s.OwnerID == 0 {
		return errors.New("secret key owner is required")
	}
	return nil
}

// BeforeUpdate hook for GORM
func (s *SecretKey) BeforeUpdate(tx *gorm.DB) error {
	if s.KeyType != "" && !IsValidSecretKeyType(s.KeyType) {
		return errors.New("invalid secret key type")
	}
	return nil
}

// IsExpired checks if the secret key is expired
func (s *SecretKey) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

// GetCustomFieldValue gets a value from custom fields
func (s *SecretKey) GetCustomFieldValue(key string) (interface{}, bool) {
	if s.SecretFields == nil {
		return nil, false
	}
	value, exists := s.SecretFields[key]
	return value, exists
}

// SetCustomFieldValue sets a value in custom fields
func (s *SecretKey) SetCustomFieldValue(key string, value interface{}) {
	if s.SecretFields == nil {
		s.SecretFields = make(SecretFields)
	}
	s.SecretFields[key] = value
}
