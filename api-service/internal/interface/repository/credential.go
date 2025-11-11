package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// CredentialCategoryRepository defines the interface for credential category data access
type CredentialCategoryRepository interface {
	// GetByID retrieves a credential category by ID
	GetByID(ctx context.Context, id uint) (*model.CredentialCategory, error)

	// List retrieves all credential categories
	List(ctx context.Context) ([]*model.CredentialCategory, error)
}

// CredentialTemplateRepository defines the interface for credential template data access
type CredentialTemplateRepository interface {
	// GetByID retrieves a credential template by ID
	GetByID(ctx context.Context, id uint) (*model.CredentialTemplate, error)

	// List retrieves credential templates with optional category filter
	List(ctx context.Context, categoryID *uint) ([]*model.CredentialTemplate, error)
}

// CredentialRepository defines the interface for credential data access
type CredentialRepository interface {
	// Create creates a new credential
	Create(ctx context.Context, credential *model.Credential) error

	// Update updates an existing credential
	Update(ctx context.Context, credential *model.Credential) error

	// Delete deletes a credential by ID
	Delete(ctx context.Context, id uint) error

	// GetByID retrieves a credential by ID
	GetByID(ctx context.Context, id uint) (*model.Credential, error)

	// GetByName retrieves a credential by name
	GetByName(ctx context.Context, name string) (*model.Credential, error)

	// ExistsByName checks if a credential name exists (excluding specific ID)
	ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error)

	// List retrieves credentials with pagination and filters
	List(ctx context.Context, req *request.ListCredentialsRequest) ([]*model.Credential, int64, error)
}
