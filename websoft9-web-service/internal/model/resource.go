package model

import (
	"time"

	"gorm.io/gorm"
)

// ========================================
// 3.3 资源管理 (Resource Management)
// ========================================

// ResourceGroup 资源组表
type ResourceGroup struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null" binding:"required"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	Description string         `json:"description" gorm:"type:text"`
	OwnerID     uint           `json:"owner_id" gorm:"not null"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	Status      int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Servers             []Server             `json:"servers" gorm:"foreignKey:ResourceGroupID"`
	AppInstances        []AppInstance        `json:"app_instances" gorm:"foreignKey:ResourceGroupID"`
	DatabaseConnections []DatabaseConnection `json:"database_connections" gorm:"foreignKey:ResourceGroupID"`
	AppGateways         []AppGateway         `json:"app_gateways" gorm:"foreignKey:ResourceGroupID"`
}

// DatabaseConnection 数据库连接表
type DatabaseConnection struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	Name              string         `json:"name" gorm:"not null" binding:"required"`
	DBType            string         `json:"db_type" gorm:"not null"` // mysql, postgresql, redis, mongodb
	Host              string         `json:"host" gorm:"not null" binding:"required"`
	Port              int            `json:"port" gorm:"not null"`
	Database          string         `json:"database"`
	Username          string         `json:"username"`
	Password          string         `json:"password"` // 加密存储
	SSLEnabled        int8           `json:"ssl_enabled" gorm:"default:0"`
	ConnectionTimeout int            `json:"connection_timeout" gorm:"default:30"` // 秒
	MaxConnections    int            `json:"max_connections" gorm:"default:10"`
	Status            string         `json:"status" gorm:"default:CONNECTED"` // CONNECTED, DISCONNECTED, ERROR
	Version           string         `json:"version"`
	Charset           string         `json:"charset"`
	Description       string         `json:"description" gorm:"type:text"`
	LastConnectedAt   *time.Time     `json:"last_connected_at"`
	OwnerID           uint           `json:"owner_id" gorm:"not null"`
	Owner             User           `json:"owner" gorm:"foreignKey:OwnerID"`
	ResourceGroupID   *uint          `json:"resource_group_id"`
	ResourceGroup     *ResourceGroup `json:"resource_group" gorm:"foreignKey:ResourceGroupID"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// Server 服务器表
type Server struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Name            string         `json:"name" gorm:"not null" binding:"required"`
	Hostname        string         `json:"hostname" gorm:"not null"`
	IPAddress       string         `json:"ip_address" gorm:"not null"`
	InternalIP      string         `json:"internal_ip"`
	SSHPort         int            `json:"ssh_port" gorm:"default:22"`
	OSType          string         `json:"os_type" gorm:"not null"`
	OSVersion       string         `json:"os_version"`
	KernelVersion   string         `json:"kernel_version"`
	CPUCores        int            `json:"cpu_cores" gorm:"default:0"`
	MemoryTotal     int64          `json:"memory_total" gorm:"default:0"` // MB
	DiskTotal       int64          `json:"disk_total" gorm:"default:0"`   // MB
	Architecture    string         `json:"architecture"`
	Status          string         `json:"status" gorm:"default:UNKNOWN"` // UNKNOWN, RUNNING, STOPPED, ERROR
	LastHeartbeatAt *time.Time     `json:"last_heartbeat_at"`
	ResourceGroupID *uint          `json:"resource_group_id"`
	ResourceGroup   *ResourceGroup `json:"resource_group" gorm:"foreignKey:ResourceGroupID"`
	OwnerID         uint           `json:"owner_id" gorm:"not null"`
	Owner           User           `json:"owner" gorm:"foreignKey:OwnerID"`
	Description     string         `json:"description" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Agents       []ServerAgent `json:"agents" gorm:"foreignKey:ServerID"`
	AppInstances []AppInstance `json:"app_instances" gorm:"foreignKey:ServerID"`
	AppGateways  []AppGateway  `json:"app_gateways" gorm:"foreignKey:ServerID"`
}

// ServerAgent 客户端表
type ServerAgent struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	ServerID        uint           `json:"server_id" gorm:"not null"`
	Server          Server         `json:"server" gorm:"foreignKey:ServerID"`
	ContainerID     string         `json:"container_id"`
	AgentIP         string         `json:"agent_ip"`
	AgentPort       int            `json:"agent_port" gorm:"default:22"`
	Version         string         `json:"version"`
	Status          string         `json:"status" gorm:"default:UNKNOWN"` // UNKNOWN, ONLINE, OFFLINE
	LastHeartbeatAt *time.Time     `json:"last_heartbeat_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// AppInstance 应用实例表
type AppInstance struct {
	ID              uint             `json:"id" gorm:"primarykey"`
	Name            string           `json:"name" gorm:"not null" binding:"required"`
	TemplateID      uint             `json:"template_id" gorm:"not null"`
	Template        AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	ServerID        uint             `json:"server_id" gorm:"not null"`
	Server          Server           `json:"server" gorm:"foreignKey:ServerID"`
	ContainerID     string           `json:"container_id"`
	ContainerName   string           `json:"container_name"`
	ImageName       string           `json:"image_name"`
	ImageTag        string           `json:"image_tag"`
	Status          string           `json:"status" gorm:"default:DEFAULT"` // DEFAULT, DEPLOYMENT, RUNNING, PAUSED, STOPPED, UPDATE
	StartedAt       *time.Time       `json:"started_at"`
	StoppedAt       *time.Time       `json:"stopped_at"`
	ResourceGroupID *uint            `json:"resource_group_id"`
	ResourceGroup   *ResourceGroup   `json:"resource_group" gorm:"foreignKey:ResourceGroupID"`
	OwnerID         uint             `json:"owner_id" gorm:"not null"`
	Owner           User             `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"-" gorm:"index"`

	// 关联关系
	GatewayPublishes []AppGatewayPublish `json:"gateway_publishes" gorm:"foreignKey:AppInstanceID"`
	Shortcuts        []AppShortcut       `json:"shortcuts" gorm:"foreignKey:AppInstanceID"`
	Deployments      []AppDeployment     `json:"deployments" gorm:"foreignKey:AppInstanceID"`
}

// 为了避免循环引用，这里使用接口或者在需要的地方重新定义
// 这些类型在其他文件中已经定义，这里只是为了完整性而声明

// Application 应用表 (兼容现有代码)
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
