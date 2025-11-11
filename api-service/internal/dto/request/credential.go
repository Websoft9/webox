package request

import "api-service/internal/dto/common"

// CreateCredentialRequest represents the request to create a credential
type CreateCredentialRequest struct {
	Name        string                    `json:"name" validate:"required,min=3,max=50,alphanum_underscore"` // Credential name (alphanumeric and underscore only)
	Description *string                   `json:"description" validate:"omitempty,max=200"`                  // Credential description
	TemplateID  uint                      `json:"template_id" validate:"required"`                           // Template ID
	Parameters  []CredentialParameterItem `json:"parameters" validate:"required,dive"`                       // Credential parameters
}

// CredentialParameterItem represents a credential parameter item
type CredentialParameterItem struct {
	InputName   string `json:"input_name" validate:"required"`  // Parameter name
	InputValue  string `json:"input_value" validate:"required"` // Parameter value
	IsEncrypted bool   `json:"is_encrypted"`                    // Is encrypted
}

// UpdateCredentialRequest represents the request to update a credential
type UpdateCredentialRequest struct {
	Description *string                   `json:"description" validate:"omitempty,max=200"` // Credential description
	Parameters  []CredentialParameterItem `json:"parameters" validate:"omitempty,dive"`     // Credential parameters
}

// ListCredentialsRequest represents the request to list credentials
type ListCredentialsRequest struct {
	CategoryID *uint `form:"category_id" binding:"omitempty"` // Filter by category ID
	TemplateID *uint `form:"template_id" binding:"omitempty"` // Filter by template ID
	common.BaseListRequest
}

// ListCredentialTemplatesRequest represents the request to list credential templates
type ListCredentialTemplatesRequest struct {
	CategoryID *uint `form:"category_id" binding:"omitempty"` // Filter by category ID
}
