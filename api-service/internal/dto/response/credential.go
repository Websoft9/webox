package response

import "time"

// CredentialCategoryResponse represents a credential category response
type CredentialCategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CredentialTemplateResponse represents a credential template response
type CredentialTemplateResponse struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	Description  *string             `json:"description"`
	CategoryID   uint                `json:"category_id"`
	CategoryName string              `json:"category_name"`
	FormSchema   []FormFieldResponse `json:"form_schema"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// FormFieldResponse represents a form field in template schema
type FormFieldResponse struct {
	InputLabel        string `json:"input_label"`
	InputName         string `json:"input_name"`
	InputType         string `json:"input_type"`
	InputDefaultValue string `json:"input_default_value"`
	InputMinLength    int    `json:"input_min_length"`
	InputMaxLength    int    `json:"input_max_length"`
	InputPlaceholder  string `json:"input_placeholder"`
	InputRequired     bool   `json:"input_required"`
	InputPattern      string `json:"input_pattern"`
	IsEncrypted       bool   `json:"is_encrypted"`
}

// CredentialResponse represents a credential response
type CredentialResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	TemplateID   uint      `json:"template_id"`
	TemplateName string    `json:"template_name"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	OwnerID      uint      `json:"owner_id"`
	OwnerName    string    `json:"owner_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CredentialDetailResponse represents a detailed credential response
type CredentialDetailResponse struct {
	ID           uint                          `json:"id"`
	Name         string                        `json:"name"`
	Description  *string                       `json:"description"`
	TemplateID   uint                          `json:"template_id"`
	TemplateName string                        `json:"template_name"`
	CategoryID   uint                          `json:"category_id"`
	CategoryName string                        `json:"category_name"`
	Parameters   []CredentialParameterResponse `json:"parameters"`
	OwnerID      uint                          `json:"owner_id"`
	OwnerName    string                        `json:"owner_name"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// CredentialParameterResponse represents a credential parameter response
type CredentialParameterResponse struct {
	InputName   string `json:"input_name"`
	InputValue  string `json:"input_value"` // Encrypted fields will show "******"
	IsEncrypted bool   `json:"is_encrypted"`
}
