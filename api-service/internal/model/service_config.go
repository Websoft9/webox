package model

import (
	"time"
)

// ServiceConfig represents the service configuration model
type ServiceConfig struct {
	ID           uint       `json:"id" gorm:"primarykey"`
	Code         string     `json:"code" gorm:"uniqueIndex;not null;size:64" validate:"required,max=64"`
	ConfigKey    string     `json:"config_key" gorm:"not null;size:64" validate:"required,max=64"`
	ConfigValue  string     `json:"config_value" gorm:"type:text"`
	ConfigType   ConfigType `json:"config_type" gorm:"type:enum('STRING','BOOLEAN','NUMBER','JSON');default:'STRING'" validate:"required,oneof=STRING BOOLEAN NUMBER JSON"`
	Category     string     `json:"category" gorm:"not null;size:32" validate:"required,max=32"`
	Description  string     `json:"description" gorm:"type:text"`
	IsReadonly   bool       `json:"is_readonly" gorm:"default:false"`
	IsEncrypted  bool       `json:"is_encrypted" gorm:"default:false"`
	DefaultValue string     `json:"default_value" gorm:"type:text"`
	SortOrder    int        `json:"sort_order" gorm:"default:0"`
	OwnerID      uint       `json:"owner_id" gorm:"not null;index" validate:"required"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for the ServiceConfig model
func (ServiceConfig) TableName() string {
	return "service_configs"
}

// GetEffectiveValue returns the effective configuration value (returns default value if current value is empty)
func (sc *ServiceConfig) GetEffectiveValue() string {
	if sc.ConfigValue != "" {
		return sc.ConfigValue
	}
	return sc.DefaultValue
}
