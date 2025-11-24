package model

import "time"

// ResourceType represents a resource type definition
type ResourceType struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `gorm:"size:64;not null" json:"name"`                                              // Resource type name (i18n key)
	Code        string    `gorm:"size:64;not null;uniqueIndex" json:"code"`                                  // Resource type code (e.g., server, database)
	Table       string    `gorm:"size:64;not null;column:table_name;index:idx_table_name" json:"table_name"` // Database table name for this resource type
	Description string    `gorm:"type:text" json:"description"`                                              // Resource type description
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                 // Creation time
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`                 // Last update time
}

// TableName returns the table name for ResourceType
func (ResourceType) TableName() string {
	return "resource_types"
}
