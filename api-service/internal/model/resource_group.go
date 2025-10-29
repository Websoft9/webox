package model

import (
	"time"
)

// ResourceGroup represents a resource group for organizing resources
type ResourceGroup struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   uint      `gorm:"not null;column:project_id;index:idx_project_id" json:"project_id"`        // Project ID
	Name        string    `gorm:"size:64;not null" json:"name"`                                             // Resource group name, unique within project
	Code        string    `gorm:"size:32;uniqueIndex;not null" json:"code"`                                 // Resource group code, globally unique, format: rg_{id}
	Description *string   `gorm:"type:text" json:"description"`                                             // Resource group description
	OwnerID     uint      `gorm:"not null;column:owner_id;index:idx_owner_id" json:"owner_id"`              // Resource group owner ID
	IsDefault   bool      `gorm:"default:false" json:"is_default"`                                          // Whether this is the default resource group (one per project)
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                                              // Sort order, smaller values appear first
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                // Creation time
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"` // Last update time
}

// TableName returns the table name for ResourceGroup
func (ResourceGroup) TableName() string {
	return "resource_groups"
}
