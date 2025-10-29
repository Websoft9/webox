package response

import "time"

// SecretResponse represents the basic secret response
type SecretResponse struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Description     *string   `json:"description"`
	ResourceGroupID uint      `json:"resource_group_id"`
	OwnerID         uint      `json:"owner_id"`
	OwnerName       string    `json:"owner_name"`
	ReferenceCount  int       `json:"reference_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SecretDetailResponse represents the detailed secret response
type SecretDetailResponse struct {
	ID                uint                      `json:"id"`
	Code              string                    `json:"code"`
	Name              string                    `json:"name"`
	Type              string                    `json:"type"`
	Description       *string                   `json:"description"`
	ResourceGroupID   uint                      `json:"resource_group_id"`
	ResourceGroupName string                    `json:"resource_group_name"`
	ExpiresAt         *time.Time                `json:"expires_at"`
	OwnerID           uint                      `json:"owner_id"`
	OwnerName         string                    `json:"owner_name"`
	References        []SecretReferenceResponse `json:"references"`
	AuthorizedUsers   []SecretAuthorizeResponse `json:"authorized_users"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

// SecretReferenceResponse represents a secret reference response
type SecretReferenceResponse struct {
	ID           uint      `json:"id"`
	ResourceCode string    `json:"resource_code"`
	CreatedAt    time.Time `json:"created_at"`
}

// SecretAuthorizeResponse represents a secret authorization response
type SecretAuthorizeResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

// SecretListResponse represents the paginated secret list response
type SecretListResponse struct {
	Total int64            `json:"total"`
	Items []SecretResponse `json:"items"`
}

// ReferenceResponse represents the response after creating a reference
type ReferenceResponse struct {
	ID           uint      `json:"id"`
	SecretID     uint      `json:"secret_id"`
	ResourceCode string    `json:"resource_code"`
	CreatedAt    time.Time `json:"created_at"`
}
