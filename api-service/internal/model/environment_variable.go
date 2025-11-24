package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EnvVarScope represents the scope of an environment variable
type EnvVarScope string

const (
	EnvVarScopePlatform EnvVarScope = "PLATFORM" // Platform-level scope, visible to all projects
	EnvVarScopeProject  EnvVarScope = "PROJECT"  // Project-level scope, visible only to specific project
)

// ValidEnvVarScopes returns all valid environment variable scopes
func ValidEnvVarScopes() []EnvVarScope {
	return []EnvVarScope{
		EnvVarScopePlatform,
		EnvVarScopeProject,
	}
}

// IsValidEnvVarScope checks if a given scope is valid
func IsValidEnvVarScope(scope EnvVarScope) bool {
	validScopes := ValidEnvVarScopes()
	for _, s := range validScopes {
		if s == scope {
			return true
		}
	}
	return false
}

// EnvironmentVariable represents an environment variable in the system
type EnvironmentVariable struct {
	ID          uint        `json:"id" gorm:"primarykey"`
	Name        string      `gorm:"size:64;not null;uniqueIndex:idx_name_scope_project;comment:Variable name" json:"name"`                            // Variable name
	Value       string      `gorm:"type:text;not null;comment:Variable value (encrypted if sensitive)" json:"value"`                                  // Variable value, encrypted if is_sensitive=true
	Scope       EnvVarScope `gorm:"size:20;not null;uniqueIndex:idx_name_scope_project;index:idx_scope;comment:Scope: PLATFORM/PROJECT" json:"scope"` // Scope: PLATFORM or PROJECT
	ProjectID   *uint       `gorm:"type:integer;uniqueIndex:idx_name_scope_project;index:idx_env_project_id;comment:Project ID" json:"project_id"`    // Project ID, required when scope=PROJECT
	Description *string     `gorm:"type:text;comment:Variable description" json:"description"`                                                        // Variable description
	IsSensitive bool        `gorm:"default:false;not null;comment:Whether sensitive variable" json:"is_sensitive"`                                    // Whether the variable is sensitive (encrypted storage)
	OwnerID     uint        `gorm:"type:integer;not null;index:idx_env_owner_id;comment:Owner user ID" json:"owner_id"`                               // Owner user ID
	CreatedAt   time.Time   `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`                         // Creation time
	UpdatedAt   time.Time   `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`                           // Last update time
}

// TableName returns the table name for EnvironmentVariable
func (EnvironmentVariable) TableName() string {
	return "environment_variables"
}

// BeforeCreate GORM hook before creating
// Only validates data integrity constraints, not business rules
func (ev *EnvironmentVariable) BeforeCreate(tx *gorm.DB) error {
	// Validate scope (data integrity)
	if !IsValidEnvVarScope(ev.Scope) {
		return fmt.Errorf("invalid scope: %s", ev.Scope)
	}

	// Validate project_id consistency with scope (data integrity)
	if ev.Scope == EnvVarScopeProject && ev.ProjectID == nil {
		return fmt.Errorf("project_id is required when scope is PROJECT")
	}

	if ev.Scope == EnvVarScopePlatform && ev.ProjectID != nil {
		return fmt.Errorf("project_id must be null when scope is PLATFORM")
	}

	return nil
}

// BeforeUpdate GORM hook before updating
// Prevents modification of immutable fields
func (ev *EnvironmentVariable) BeforeUpdate(tx *gorm.DB) error {
	// Use SELECT to only fetch the fields we need to compare (performance optimization)
	var oldRecord EnvironmentVariable
	if err := tx.Select("scope", "name", "project_id").First(&oldRecord, ev.ID).Error; err == nil {
		if oldRecord.Scope != ev.Scope {
			return fmt.Errorf("scope cannot be modified after creation")
		}
		if oldRecord.Name != ev.Name {
			return fmt.Errorf("name cannot be modified after creation")
		}
		if oldRecord.ProjectID != ev.ProjectID {
			// Handle nil pointer comparison
			if (oldRecord.ProjectID == nil && ev.ProjectID != nil) ||
				(oldRecord.ProjectID != nil && ev.ProjectID == nil) ||
				(oldRecord.ProjectID != nil && ev.ProjectID != nil && *oldRecord.ProjectID != *ev.ProjectID) {
				return fmt.Errorf("project_id cannot be modified after creation")
			}
		}
	}

	return nil
}
