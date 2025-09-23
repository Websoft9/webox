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
	SecretKeyTypeAPIKey      SecretKeyType = "API_KEY"
	SecretKeyTypeDatabase    SecretKeyType = "DATABASE"
	SecretKeyTypeSSH         SecretKeyType = "SSH"
	SecretKeyTypeCertificate SecretKeyType = "CERTIFICATE"
	SecretKeyTypeCustom      SecretKeyType = "CUSTOM"
)

// CustomFields represents JSON custom fields
type CustomFields map[string]interface{}

// Value implements driver.Valuer interface for database storage
func (cf CustomFields) Value() (driver.Value, error) {
	if cf == nil {
		return nil, nil
	}
	return json.Marshal(cf)
}

// Scan implements sql.Scanner interface for database retrieval
func (cf *CustomFields) Scan(value interface{}) error {
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
	EncryptedValue  string        `json:"encrypted_value" gorm:"type:text;not null"`
	Description     *string       `json:"description" gorm:"type:text"`
	CustomFields    CustomFields  `json:"custom_fields" gorm:"type:text"`
	ExpiresAt       *time.Time    `json:"expires_at"`
	ResourceGroupID *uint         `json:"resource_group_id"`
	OwnerID         uint          `json:"owner_id" gorm:"not null"`
	CreatedAt       time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Owner *User `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
}

// TableName returns the table name for GORM
func (SecretKey) TableName() string {
	return "secret_keys"
}

// BeforeCreate hook for GORM
func (s *SecretKey) BeforeCreate(tx *gorm.DB) error {
	if s.Name == "" {
		return errors.New("secret key name is required")
	}
	if s.KeyType == "" {
		return errors.New("secret key type is required")
	}
	if s.EncryptedValue == "" {
		return errors.New("secret key value is required")
	}
	if s.OwnerID == 0 {
		return errors.New("secret key owner is required")
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
