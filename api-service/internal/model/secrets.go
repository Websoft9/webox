package model

import (
	"time"
)

// SecretType represents the type of secret
type SecretType string

const (
	SecretTypeText    SecretType = "text"    // Text type secret
	SecretTypeAccount SecretType = "account" // Account type secret (username + password)
	SecretTypeFile    SecretType = "file"    // File type secret (certificate, key file, etc.)
)

// Secret represents a secret for storing sensitive credentials
type Secret struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// Secret code, globally unique, format: secrets_{random_string}
	Code            string     `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Name            string     `gorm:"size:64;not null" json:"name"`                                                           // Secret name
	Type            SecretType `gorm:"type:enum('text','account','file');not null" json:"type"`                                // Secret type
	Description     *string    `gorm:"type:text" json:"description"`                                                           // Secret description
	SecretFields    JSON       `gorm:"type:json;not null" json:"secret_fields"`                                                // Encrypted secret data (JSON format)
	ExpiresAt       *time.Time `gorm:"type:datetime;serializer:datetime" json:"expires_at"`                                    // Expiration time
	ResourceGroupID uint       `gorm:"not null;column:resource_group_id;index:idx_resource_group_id" json:"resource_group_id"` // Resource group ID
	OwnerID         uint       `gorm:"not null;column:owner_id;index:idx_owner_id" json:"owner_id"`                            // Owner ID
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                              // Creation time
	UpdatedAt       time.Time  `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                              // Last update time
}

// TableName returns the table name for Secret
func (Secret) TableName() string {
	return "secrets"
}

// SecretReference represents the reference relationship between secrets and resources
type SecretReference struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SecretID     uint      `gorm:"not null;column:secret_id;index:idx_secret_id" json:"secret_id"` // Secret ID
	ResourceCode string    `gorm:"size:128;not null;index:idx_resource_code" json:"resource_code"` // Resource code
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`      // Creation time
}

// TableName returns the table name for SecretReference
func (SecretReference) TableName() string {
	return "secret_references"
}

// SecretAuthorize represents the authorization relationship between secrets and users
type SecretAuthorize struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SecretID         uint      `gorm:"not null;column:secret_id;index:idx_secret_id" json:"secret_id"`                            // Secret ID
	AuthorizedUserID uint      `gorm:"not null;column:authorized_user_id;index:idx_authorized_user_id" json:"authorized_user_id"` // Authorized user ID
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                                 // Creation time
}

// TableName returns the table name for SecretAuthorize
func (SecretAuthorize) TableName() string {
	return "secret_authorizes"
}
