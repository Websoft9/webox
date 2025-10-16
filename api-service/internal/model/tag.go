package model

import (
	"time"
)

// Tag represents a tag entity in the system
type Tag struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:128" example:"production"`
	Color       string    `json:"color" gorm:"size:16" example:"#ff0000"`
	Description string    `json:"description" gorm:"type:text" example:"Production environment tag"`
	CreatedBy   uint      `json:"created_by" gorm:"default:0" example:"1"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`

	// Relations
	Taggings []Tagging `json:"taggings,omitempty" gorm:"foreignKey:TagID"`
}

// TableName specifies the table name for Tag model
func (Tag) TableName() string {
	return "tags"
}

// Tagging represents the association between tags and resources
type Tagging struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	TagID      uint      `json:"tag_id" gorm:"not null;index" example:"1"`
	ResourceID uint      `json:"resource_id" gorm:"not null;index" example:"123"`
	CreatedBy  uint      `json:"created_by" gorm:"default:0" example:"1"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`

	// Relations
	Tag *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID"`
}

// TableName specifies the table name for Tagging model
func (Tagging) TableName() string {
	return "taggings"
}
