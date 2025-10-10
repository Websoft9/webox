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
	ID              uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string          `json:"name" gorm:"type:varchar(64);not null;index"`
	Hostname        string          `json:"hostname" gorm:"type:varchar(255);not null"`
	Host            string          `json:"host" gorm:"type:varchar(255);not null;index;comment:Host address (IP or domain) for SSH/Agent priority connection"`
	InternalIP      *string         `json:"internal_ip" gorm:"type:varchar(45);comment:Internal IP address"`
	IPv6Address     *string         `json:"ipv6_address" gorm:"type:varchar(45);comment:IPv6 address for future expansion"`
	SSHPort         int             `json:"ssh_port" gorm:"default:22;comment:SSH port"`
	SSHCredentialID *string         `json:"ssh_credential_id" gorm:"type:varchar(64);comment:SSH credential ID from key management"`
	OSDistro        *string         `json:"os_distro" gorm:"type:varchar(32);comment:Operating system distribution"`
	OSVersion       *string         `json:"os_version" gorm:"type:varchar(64);comment:Operating system version"`
	KernelVersion   *string         `json:"kernel_version" gorm:"type:varchar(64);comment:Kernel version"`
	CPUCores        int             `json:"cpu_cores" gorm:"default:0;comment:CPU cores count"`
	MemoryTotal     int64           `json:"memory_total" gorm:"default:0;comment:Total memory (MB)"`
	DiskTotal       int64           `json:"disk_total" gorm:"default:0;comment:Total disk space (MB)"`
	Architecture    *string         `json:"architecture" gorm:"type:varchar(16);comment:System architecture"`
	ResourceGroupID *uint           `json:"resource_group_id" gorm:"index;comment:Resource group ID"`
	OwnerID         uint            `json:"owner_id" gorm:"not null;index;comment:Owner user ID"`
	Description     *string         `json:"description" gorm:"type:text;comment:Server description"`
	CreatedAt       time.Time       `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:Soft delete time"`
}

// TableName returns the table name for Server model
func (Server) TableName() string {
	return "servers"
}

// ServerAgent represents agent deployment information for a server
type ServerAgent struct {
	ID              uint                `json:"id" gorm:"primaryKey;autoIncrement"`
	ServerID        uint                `json:"server_id" gorm:"uniqueIndex;not null;comment:Server ID"`
	AgentID         string              `json:"agent_id" gorm:"type:varchar(64);uniqueIndex;not null;comment:Agent unique identifier"`
	DeploymentType  AgentDeploymentType `json:"deployment_type" gorm:"type:enum('docker','systemd');default:'docker';comment:Agent deployment type"`
	ContainerID     *string             `json:"container_id" gorm:"type:varchar(64);comment:Container ID (for Docker deployment)"`
	ContainerName   *string             `json:"container_name" gorm:"type:varchar(128);comment:Container name (for Docker deployment)"`
	ServiceName     *string             `json:"service_name" gorm:"type:varchar(64);comment:Service name (for systemd deployment)"`
	BinaryPath      *string             `json:"binary_path" gorm:"type:varchar(255);comment:Binary file path (for systemd deployment)"`
	ConfigPath      *string             `json:"config_path" gorm:"type:varchar(255);comment:Configuration file path"`
	AgentIP         *string             `json:"agent_ip" gorm:"type:varchar(45);comment:Agent IP address"`
	AgentPort       int                 `json:"agent_port" gorm:"default:8080;comment:Agent communication port"`
	Version         *string             `json:"version" gorm:"type:varchar(32);comment:Agent version"`
	PullMode        bool                `json:"pull_mode" gorm:"default:true;comment:Whether Pull mode is enabled"`
	PullInterval    int                 `json:"pull_interval" gorm:"default:30;comment:Pull interval in seconds"`
	LastHeartbeatAt *time.Time          `json:"last_heartbeat_at" gorm:"comment:Last heartbeat time"`
	CreatedAt       time.Time           `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

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
