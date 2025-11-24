package model

import (
	"time"
)

// ResourceGroup represents a resource group for organizing resources
type ResourceGroup struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	ProjectID   uint      `gorm:"type:integer;not null;column:project_id;index:idx_group_project_id;comment:Project ID" json:"project_id"`         // Project ID
	Name        string    `gorm:"size:64;not null;comment:Resource group name" json:"name"`                                                        // Resource group name, unique within project
	Code        string    `gorm:"size:32;uniqueIndex:uk_group_code;not null;comment:Resource group code" json:"code"`                              // Resource group code, globally unique, format: rg_{id}
	Description *string   `gorm:"type:text;comment:Resource group description" json:"description"`                                                 // Resource group description
	OwnerID     uint      `gorm:"type:integer;not null;column:owner_id;index:idx_rs_owner_id;comment:Owner ID" json:"owner_id"`                    // Resource group owner ID
	IsDefault   bool      `gorm:"type:tinyint(1);default:0;index:idx_is_default;comment:Is default resource group: 0-no, 1-yes" json:"is_default"` // Whether this is the default resource group (one per project)
	SortOrder   int       `gorm:"default:0;comment:Sort order" json:"sort_order"`                                                                  // Sort order, smaller values appear first
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`                        // Creation time
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`                          // Last update time
}

// TableName returns the table name for ResourceGroup
func (ResourceGroup) TableName() string {
	return "resource_groups"
}
