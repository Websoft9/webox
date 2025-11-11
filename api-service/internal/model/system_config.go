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
	ID              uint       `json:"id" gorm:"primarykey;type:bigint unsigned"`
	ConfigKey       string     `json:"config_key" gorm:"uniqueIndex;not null;size:64;comment:Configuration key" validate:"required,max=64"`
	ConfigValue     string     `json:"config_value" gorm:"type:text;comment:Configuration value"`
	ConfigType      ConfigType `json:"config_type" gorm:"type:varchar(20);default:'STRING';index:idx_system_config_type;comment:Configuration type" validate:"required,oneof=STRING INTEGER BOOLEAN JSON FLOAT"`
	Category        string     `json:"category" gorm:"not null;size:32;index:idx_system_category;comment:Configuration category" validate:"required,max=32"`
	Description     string     `json:"description" gorm:"type:text;comment:Configuration description"`
	IsReadonly      bool       `json:"is_readonly" gorm:"type:tinyint(1);default:0;comment:Whether read-only"`
	IsEncrypted     bool       `json:"is_encrypted" gorm:"type:tinyint(1);default:0;comment:Whether encrypted"`
	DefaultValue    string     `json:"default_value" gorm:"type:text;comment:Default value"`
	ValidationRules string     `json:"validation_rules" gorm:"type:json;comment:Validation rules"`
	SortOrder       int        `json:"sort_order" gorm:"default:0;index:idx_system_sort_order;comment:Sort order"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
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
