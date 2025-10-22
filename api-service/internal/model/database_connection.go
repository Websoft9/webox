package model

import (
	"time"
)

// DatabaseConnection represents a database connection configuration
type DatabaseConnection struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"size:64;not null" json:"name"`
	Code            string    `gorm:"size:64;uniqueIndex;not null" json:"code"`                         // Unique connection code, format: db_conn_{id}
	DBType          string    `gorm:"size:32;not null;column:db_type;index:idx_db_type" json:"db_type"` // mysql, postgresql, mariadb, sqlserver, oracle, sqlite
	Host            string    `gorm:"size:255;not null" json:"host"`
	Port            int       `gorm:"type:int;not null" json:"port"`
	Database        *string   `gorm:"size:64" json:"database"`     // Optional, for scenarios like PostgreSQL server connection
	Description     *string   `gorm:"size:255" json:"description"` // Connection description
	Config          *JSON     `gorm:"type:text" json:"config"`     // Additional configuration in JSON format (uses common.JSON type)
	OwnerID         uint      `gorm:"not null;column:owner_id;index:idx_owner_id" json:"owner_id"`
	ResourceGroupID *uint     `gorm:"column:resource_group_id;index:idx_resource_group_id" json:"resource_group_id"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for DatabaseConnection
func (DatabaseConnection) TableName() string {
	return "database_connections"
}
