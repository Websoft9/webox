package model

import (
	"time"
)

// CredentialCategory represents a credential category
type CredentialCategory struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `gorm:"size:100;not null;uniqueIndex:uk_lcategory_name;comment:Category name" json:"name"` // Category name
	Description *string   `gorm:"size:500;comment:Category description" json:"description"`                          // Category description
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName returns the table name for CredentialCategory
func (CredentialCategory) TableName() string {
	return "credential_categories"
}

// CredentialTemplate represents a credential template
type CredentialTemplate struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `gorm:"size:100;not null;comment:Template name" json:"name"`                                                   // Template name
	Description *string   `gorm:"size:500;comment:Template description" json:"description"`                                              // Template description
	CategoryID  uint      `gorm:"type:integer;not null;column:category_id;index:idx_category_id;comment:Category ID" json:"category_id"` // Category ID
	FormSchema  JSON      `gorm:"type:json;not null;comment:Form schema (JSON format)" json:"form_schema"`                               // Form schema (JSON format)
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`              // Creation time
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`                // Update time

	// Association fields
	Category *CredentialCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName returns the table name for CredentialTemplate
func (CredentialTemplate) TableName() string {
	return "credential_templates"
}

// Credential represents a credential record
type Credential struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `gorm:"size:100;not null;uniqueIndex:uk_credential_name;comment:Credential name (alphanumeric and underscore only)" json:"name"` // Credential name (alphanumeric and underscore only)
	Description *string   `gorm:"size:500;comment:Credential description" json:"description"`                                                              // Credential description
	TemplateID  uint      `gorm:"type:integer;not null;column:template_id;index:idx_credential_template_id;comment:Template ID" json:"template_id"`        // Template ID
	Parameters  JSON      `gorm:"type:json;not null;comment:Credential parameters (JSON format)" json:"parameters"`                                        // Credential parameters (JSON format)
	OwnerID     uint      `gorm:"type:integer;not null;column:owner_id;index:idx_cred_owner_id;comment:Owner user ID" json:"owner_id"`                     // Owner user ID
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`                                // Creation time
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`                                  // Update time

	// Association fields
	Template *CredentialTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	Owner    *User               `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
}

// TableName returns the table name for Credential
func (Credential) TableName() string {
	return "credentials"
}

// FormField represents a form field definition in template schema
type FormField struct {
	InputLabel        string `json:"input_label"`         // Field label
	InputName         string `json:"input_name"`          // Field name
	InputType         string `json:"input_type"`          // Field type (text, password, number, etc.)
	InputDefaultValue string `json:"input_default_value"` // Default value
	InputMinLength    int    `json:"input_min_length"`    // Minimum length
	InputMaxLength    int    `json:"input_max_length"`    // Maximum length
	InputPlaceholder  string `json:"input_placeholder"`   // Placeholder text
	InputRequired     bool   `json:"input_required"`      // Is required
	InputPattern      string `json:"input_pattern"`       // Validation pattern (regex)
	IsEncrypted       bool   `json:"is_encrypted"`        // Should be encrypted
}

// CredentialParameter represents a credential parameter
type CredentialParameter struct {
	InputName   string `json:"input_name"`   // Parameter name
	InputValue  string `json:"input_value"`  // Parameter value
	IsEncrypted bool   `json:"is_encrypted"` // Is encrypted
}
