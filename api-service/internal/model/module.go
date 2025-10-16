package model

import "time"

// Module represents a system module
type Module struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:64;not null"`
	Code        string    `json:"code" gorm:"size:64;not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for Module model
func (Module) TableName() string {
	return "modules"
}
