package model

import (
	"time"
)

// Tag represents a tag entity in the system
type Tag struct {
	ID          uint      `json:"id" gorm:"primarykey;type:bigint"`
	Name        string    `json:"name" gorm:"uniqueIndex:ux_name;not null;size:128;index:idx_tag_name;comment:Tag name" example:"production"`
	Color       string    `json:"color" gorm:"size:16;comment:Tag color" example:"#ff0000"`
	Description string    `json:"description" gorm:"type:text;comment:Tag description" example:"Production environment tag"`
	CreatedBy   uint      `json:"created_by" gorm:"type:bigint;default:0;index:idx_created_by;comment:Creator user ID" example:"1"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Taggings []Tagging `json:"taggings,omitempty" gorm:"foreignKey:TagID"`
}

// TableName specifies the table name for Tag model
func (Tag) TableName() string {
	return "tags"
}

// Tagging represents the association between tags and resources
type Tagging struct {
	ID           uint      `json:"id" gorm:"primarykey;type:bigint"`
	TagID        uint      `json:"tag_id" gorm:"type:bigint;not null;index:idx_tag_id;comment:Tag ID" example:"1"`
	ResourceCode string    `json:"resource_code" gorm:"not null;index:idx_tag_resource_code;size:64;comment:Resource code" example:"SERVER_001"`
	CreatedBy    uint      `json:"created_by" gorm:"type:bigint;default:0;comment:Association creator" example:"1"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Tag *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Tagging model
func (Tagging) TableName() string {
	return "taggings"
}
