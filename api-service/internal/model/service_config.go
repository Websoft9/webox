package model

import (
	"time"
)

// ServiceConfig represents the service configuration model
type ServiceConfig struct {
	ID           uint       `json:"id" gorm:"primarykey;type:bigint unsigned"`
	Code         string     `json:"code" gorm:"uniqueIndex;not null;size:64;comment:Configuration code (unique identifier)" validate:"required,max=64"`
	ConfigKey    string     `json:"config_key" gorm:"not null;size:64;comment:Configuration key" validate:"required,max=64"`
	ConfigValue  string     `json:"config_value" gorm:"type:text;comment:Configuration value"`
	ConfigType   ConfigType `json:"config_type" gorm:"type:varchar(20);default:'STRING';index:idx_service_config_type;comment:Configuration type" validate:"required,oneof=STRING INTEGER BOOLEAN JSON FLOAT"`
	Category     string     `json:"category" gorm:"not null;size:32;index:idx_service_category;comment:Configuration category" validate:"required,max=32"`
	Description  string     `json:"description" gorm:"type:text;comment:Configuration description"`
	IsReadonly   bool       `json:"is_readonly" gorm:"type:tinyint(1);default:0;comment:Whether read-only"`
	IsEncrypted  bool       `json:"is_encrypted" gorm:"type:tinyint(1);default:0;comment:Whether encrypted"`
	DefaultValue string     `json:"default_value" gorm:"type:text;comment:Default value"`
	SortOrder    int        `json:"sort_order" gorm:"default:0;index:idx_service_sort_order;comment:Sort order"`
	OwnerID      uint       `json:"owner_id" gorm:"type:bigint unsigned;not null;comment:Owner ID" validate:"required"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
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
