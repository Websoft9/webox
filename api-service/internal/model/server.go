package model

import (
	"time"

	"gorm.io/gorm"
)

// AgentDeploymentType defines the agent deployment type
type AgentDeploymentType string

const (
	AgentDeploymentDocker  AgentDeploymentType = "docker"
	AgentDeploymentSystemd AgentDeploymentType = "systemd"
)

// Server represents a managed server in the platform
type Server struct {
	ID              uint            `json:"id" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	Name            string          `json:"name" gorm:"type:varchar(64);not null;index:idx_server_name;comment:Server name"`
	Code            string          `json:"code" gorm:"type:varchar(64);not null;index:idx_server_code;comment:Server code"`
	Hostname        string          `json:"hostname" gorm:"type:varchar(255);not null;comment:Hostname"`
	Host            string          `json:"host" gorm:"type:varchar(255);not null;index:idx_host;comment:Host address (IP or domain) for SSH/Agent priority connection"`
	InternalIP      *string         `json:"internal_ip" gorm:"type:varchar(45);comment:Internal IP address"`
	IPv6Address     *string         `json:"ipv6_address" gorm:"type:varchar(45);comment:IPv6 address for future expansion"`
	SSHPort         int             `json:"ssh_port" gorm:"type:int;default:22;comment:SSH port"`
	OSDistro        *string         `json:"os_distro" gorm:"type:varchar(32);comment:Operating system distribution"`
	OSVersion       *string         `json:"os_version" gorm:"type:varchar(64);comment:Operating system version"`
	KernelVersion   *string         `json:"kernel_version" gorm:"type:varchar(64);comment:Kernel version"`
	CPUCores        int             `json:"cpu_cores" gorm:"type:int;default:0;comment:CPU cores count"`
	MemoryTotal     int64           `json:"memory_total" gorm:"type:bigint;default:0;comment:Total memory (MB)"`
	DiskTotal       int64           `json:"disk_total" gorm:"type:bigint;default:0;comment:Total disk space (MB)"`
	Architecture    *string         `json:"architecture" gorm:"type:varchar(16);comment:System architecture"`
	ResourceGroupID *uint           `json:"resource_group_id" gorm:"type:bigint unsigned;index:idx_server_resource_group_id;comment:Resource group ID"`
	OwnerID         uint            `json:"owner_id" gorm:"type:bigint unsigned;not null;index:idx_server_owner_id;comment:Owner user ID"`
	Description     *string         `json:"description" gorm:"type:text;comment:Server description"`
	CreatedAt       time.Time       `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at" gorm:"index:idx_deleted_at;comment:Soft delete time"`
}

// TableName returns the table name for Server model
func (Server) TableName() string {
	return "servers"
}

// ServerAgent represents agent deployment information for a server
type ServerAgent struct {
	ID              uint                `json:"id" gorm:"primaryKey;autoIncrement;type:bigint unsigned"`
	ServerID        uint                `json:"server_id" gorm:"type:bigint unsigned;uniqueIndex:uk_server_id;not null;index:idx_server_id;comment:Server ID"`
	AgentID         string              `json:"agent_id" gorm:"type:varchar(64);uniqueIndex:uk_agent_id;not null;comment:Agent unique identifier"`
	DeploymentType  AgentDeploymentType `json:"deployment_type" gorm:"type:varchar(20);default:'docker';index:idx_deployment_type;comment:Agent deployment type: docker container or systemd service"`
	ContainerID     *string             `json:"container_id" gorm:"type:varchar(64);comment:Container ID (for Docker deployment)"`
	ContainerName   *string             `json:"container_name" gorm:"type:varchar(128);comment:Container name (for Docker deployment)"`
	ServiceName     *string             `json:"service_name" gorm:"type:varchar(64);comment:Service name (for systemd deployment)"`
	BinaryPath      *string             `json:"binary_path" gorm:"type:varchar(255);comment:Binary file path (for systemd deployment)"`
	ConfigPath      *string             `json:"config_path" gorm:"type:varchar(255);comment:Configuration file path"`
	AgentIP         *string             `json:"agent_ip" gorm:"type:varchar(45);comment:Agent IP address"`
	AgentPort       int                 `json:"agent_port" gorm:"type:int;default:8080;comment:Agent communication port"`
	Version         *string             `json:"version" gorm:"type:varchar(32);comment:Agent version"`
	PullMode        bool                `json:"pull_mode" gorm:"type:boolean;default:true;comment:Whether Pull mode is enabled"`
	PullInterval    int                 `json:"pull_interval" gorm:"type:int;default:30;comment:Pull interval in seconds"`
	LastHeartbeatAt *time.Time          `json:"last_heartbeat_at" gorm:"type:datetime;index:idx_last_heartbeat;comment:Last heartbeat time"`
	CreatedAt       time.Time           `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Creation time"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:Update time"`

	// Relations
	Server *Server `json:"server,omitempty" gorm:"foreignKey:ServerID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for ServerAgent model
func (ServerAgent) TableName() string {
	return "server_agents"
}

// IsOnline checks if the agent is online based on last heartbeat
func (sa *ServerAgent) IsOnline() bool {
	if sa.LastHeartbeatAt == nil {
		return false
	}
	// Consider agent offline if no heartbeat for more than 2 minutes
	return time.Since(*sa.LastHeartbeatAt) <= 2*time.Minute
}

// IsDockerDeployment checks if agent is deployed as Docker container
func (sa *ServerAgent) IsDockerDeployment() bool {
	return sa.DeploymentType == AgentDeploymentDocker
}

// IsSystemdDeployment checks if agent is deployed as systemd service
func (sa *ServerAgent) IsSystemdDeployment() bool {
	return sa.DeploymentType == AgentDeploymentSystemd
}
