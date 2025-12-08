package model

import "time"

// Module represents a system module
type Module struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"size:64;not null;comment:module name"`
	Code        string    `json:"code" gorm:"size:64;not null;uniqueIndex;index:idx_module_code;comment:module code"`
	Description string    `json:"description" gorm:"type:text;comment:module description"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName specifies the table name for Module model
func (Module) TableName() string {
	return "modules"
}
