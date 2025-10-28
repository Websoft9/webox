package response

import "time"

// EnvironmentVariableResponse represents the basic environment variable response
type EnvironmentVariableResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Value       string    `json:"value"` // Masked value for sensitive variables
	Scope       string    `json:"scope"` // PLATFORM or PROJECT
	ProjectID   *uint     `json:"project_id"`
	Description *string   `json:"description"`
	IsSensitive bool      `json:"is_sensitive"`
	OwnerID     uint      `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EnvironmentVariableDetailResponse represents the detailed environment variable response
type EnvironmentVariableDetailResponse struct {
	EnvironmentVariableResponse
}

// ResolveEnvVarResponse represents the response for variable interpolation resolution
type ResolveEnvVarResponse struct {
	Result string `json:"result"` // Resolved template string
}
