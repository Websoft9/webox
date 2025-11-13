package model

import (
	"time"
)

// DatabaseConnection represents a database connection configuration
type DatabaseConnection struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;type:bigint unsigned" json:"id"`
	Name            string    `gorm:"size:64;not null;uniqueIndex:uk_owner_name;comment:Connection name" json:"name"`
	Code            string    `gorm:"size:64;uniqueIndex:uk_db_code;not null;comment:Connection code (format: db_conn_{id}, auto-generated)" json:"code"`                             // Unique connection code, format: db_conn_{id}
	DBType          string    `gorm:"size:32;not null;column:db_type;index:idx_db_type;comment:Database type (mysql, postgresql, mariadb, sqlserver, oracle, sqlite)" json:"db_type"` // mysql, postgresql, mariadb, sqlserver, oracle, sqlite
	Host            string    `gorm:"size:255;not null;comment:Host address" json:"host"`
	Port            int       `gorm:"type:int;not null;comment:Port number" json:"port"`
	Database        *string   `gorm:"size:64;comment:Database name (optional for some scenarios)" json:"database"`                        // Optional, for scenarios like PostgreSQL server connection
	Description     *string   `gorm:"size:255;comment:Connection description" json:"description"`                                         // Connection description
	Config          *JSON     `gorm:"type:text;comment:Extra configuration (JSON format for database-specific parameters)" json:"config"` // Additional configuration in JSON format (uses common.JSON type)
	OwnerID         uint      `gorm:"type:bigint unsigned;not null;column:owner_id;index:idx_db_owner_id;comment:Owner ID" json:"owner_id"`
	ResourceGroupID *uint     `gorm:"type:bigint unsigned;column:resource_group_id;index:idx_db_resource_group_id;comment:Resource group ID" json:"resource_group_id"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
}

// TableName returns the table name for DatabaseConnection
func (DatabaseConnection) TableName() string {
	return "database_connections"
}
