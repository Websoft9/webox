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
	ID uint `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:Secret ID" json:"id"`
	// Secret code, globally unique, format: secrets_{random_string}
	Code            string     `gorm:"size:128;uniqueIndex:uk_secret_code;not null;comment:Secret code, globally unique, format: secrets_{random_string}" json:"code"`
	Name            string     `gorm:"size:64;not null;uniqueIndex:uk_resource_group_name;comment:Secret name" json:"name"`                                                                                             // Secret name
	Type            SecretType `gorm:"type:varchar(20);not null;index:idx_secret_type;comment:Secret type" json:"type"`                                                                                                 // Secret type
	Description     *string    `gorm:"type:text;comment:Secret description" json:"description"`                                                                                                                         // Secret description
	SecretFields    JSON       `gorm:"type:json;not null;comment:Encrypted secret data (JSON format)" json:"secret_fields"`                                                                                             // Encrypted secret data (JSON format)
	ExpiresAt       *time.Time `gorm:"type:datetime;serializer:datetime;comment:Expiration time" json:"expires_at"`                                                                                                     // Expiration time
	ResourceGroupID uint       `gorm:"type:bigint unsigned;not null;column:resource_group_id;uniqueIndex:uk_resource_group_name;index:idx_secret_resource_group_id;comment:Resource group ID" json:"resource_group_id"` // Resource group ID
	OwnerID         uint       `gorm:"type:bigint unsigned;not null;column:owner_id;index:idx_secret_owner_id;comment:Owner ID" json:"owner_id"`                                                                        // Owner ID
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_secret_created_at;comment:Creation time"`                                                            // Creation time
	UpdatedAt       time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`                                                                                          // Last update time
}

// TableName returns the table name for Secret
func (Secret) TableName() string {
	return "secrets"
}

// SecretReference represents the reference relationship between secrets and resources
type SecretReference struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:Reference ID" json:"id"`
	SecretID     uint      `gorm:"type:bigint unsigned;not null;column:secret_id;uniqueIndex:uk_secret_resource;comment:Secret ID" json:"secret_id"`              // Secret ID
	ResourceCode string    `gorm:"size:128;not null;uniqueIndex:uk_secret_resource;index:idx_reference_resource_code;comment:Resource code" json:"resource_code"` // Resource code
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`                                      // Creation time
}

// TableName returns the table name for SecretReference
func (SecretReference) TableName() string {
	return "secret_references"
}

// SecretAuthorize represents the authorization relationship between secrets and users
type SecretAuthorize struct {
	ID               uint      `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:Authorization ID" json:"id"`
	SecretID         uint      `gorm:"type:bigint unsigned;not null;column:secret_id;uniqueIndex:uk_secret_user;comment:Secret ID" json:"secret_id"`                                                         // Secret ID
	AuthorizedUserID uint      `gorm:"type:bigint unsigned;not null;column:authorized_user_id;uniqueIndex:uk_secret_user;index:idx_authorized_user_id;comment:Authorized user ID" json:"authorized_user_id"` // Authorized user ID
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`                                                                             // Creation time
}

// TableName returns the table name for SecretAuthorize
func (SecretAuthorize) TableName() string {
	return "secret_authorizes"
}
