package model

import (
	"time"
	"gorm.io/gorm"
)

type Application struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Version     string         `json:"version"`
	Status      string         `json:"status" gorm:"default:stopped"`
	ServerID    uint           `json:"server_id"`
	Server      Server         `json:"server" gorm:"foreignKey:ServerID"`
	Config      string         `json:"config" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Server struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"not null"`
	Host         string         `json:"host" gorm:"not null"`
	Port         int            `json:"port" gorm:"default:22"`
	Status       string         `json:"status" gorm:"default:offline"`
	OS           string         `json:"os"`
	Architecture string         `json:"architecture"`
	AgentVersion string         `json:"agent_version"`
	LastSeen     *time.Time     `json:"last_seen"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Gateway struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"not null"`
	Domain       string         `json:"domain" gorm:"uniqueIndex;not null"`
	ApplicationID uint          `json:"application_id"`
	Application  Application    `json:"application" gorm:"foreignKey:ApplicationID"`
	SSLEnabled   bool           `json:"ssl_enabled" gorm:"default:false"`
	SSLCert      string         `json:"ssl_cert"`
	AccessRules  string         `json:"access_rules" gorm:"type:text"`
	Status       string         `json:"status" gorm:"default:active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}