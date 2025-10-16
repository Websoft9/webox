package model

import (
	"api-service/internal/constants"
	"time"
)

// ConfigType represents the type of configuration value
type ConfigType string

const (
	ConfigTypeString  ConfigType = constants.ConfigTypeString
	ConfigTypeBoolean ConfigType = constants.ConfigTypeBoolean
	ConfigTypeNumber  ConfigType = constants.ConfigTypeNumber
	ConfigTypeJSON    ConfigType = constants.ConfigTypeJSON
)

// SystemConfig represents the system configuration model
type SystemConfig struct {
	ID           uint       `json:"id" gorm:"primarykey"`
	ConfigKey    string     `json:"config_key" gorm:"uniqueIndex;not null;size:64" validate:"required,max=64"`
	ConfigValue  string     `json:"config_value" gorm:"type:text"`
	ConfigType   ConfigType `json:"config_type" gorm:"type:enum('STRING','BOOLEAN','NUMBER','JSON');default:'STRING'" validate:"required,oneof=STRING BOOLEAN NUMBER JSON"`
	Category     string     `json:"category" gorm:"not null;size:32" validate:"required,max=32"`
	Description  string     `json:"description" gorm:"type:text"`
	IsReadonly   bool       `json:"is_readonly" gorm:"default:false"`
	IsEncrypted  bool       `json:"is_encrypted" gorm:"default:false"`
	DefaultValue string     `json:"default_value" gorm:"type:text"`
	SortOrder    int        `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for the SystemConfig model
func (SystemConfig) TableName() string {
	return "system_configs"
}

// GetEffectiveValue returns the effective configuration value (returns default value if current value is empty)
func (sc *SystemConfig) GetEffectiveValue() string {
	if sc.ConfigValue != "" {
		return sc.ConfigValue
	}
	return sc.DefaultValue
}
