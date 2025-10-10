package request

import "api-service/internal/dto/common"

// CreateServerRequest represents the request to create a server
type CreateServerRequest struct {
	Name            string  `json:"name" binding:"required,max=64" validate:"required,max=64"`
	Hostname        string  `json:"hostname" binding:"required,max=255" validate:"required,max=255"`
	Host            string  `json:"host" binding:"required,max=255" validate:"required,max=255"`
	InternalIP      *string `json:"internal_ip" binding:"omitempty,ip" validate:"omitempty,ip"`
	IPv6Address     *string `json:"ipv6_address" binding:"omitempty,ipv6" validate:"omitempty,ipv6"`
	SSHPort         int     `json:"ssh_port" binding:"omitempty,min=1,max=65535" validate:"omitempty,min=1,max=65535"`
	SSHCredentialID *string `json:"ssh_credential_id" binding:"omitempty,max=64" validate:"omitempty,max=64"`
	// Direct credential fields for creation (will be stored via credential management)
	SSHUsername     string `json:"ssh_username" binding:"omitempty,max=64" validate:"omitempty,max=64"`
	SSHPassword     string `json:"ssh_password" binding:"omitempty" validate:"omitempty"`
	SSHKey          string `json:"ssh_key" binding:"omitempty" validate:"omitempty"`
	ResourceGroupID *uint  `json:"resource_group_id" binding:"omitempty" validate:"omitempty"`
	OwnerID         uint   `json:"owner_id" binding:"required" validate:"required"`
	Description     string `json:"description" binding:"omitempty" validate:"omitempty"`
}

// UpdateServerRequest represents the request to update a server
type UpdateServerRequest struct {
	Name            *string `json:"name" binding:"omitempty,max=64" validate:"omitempty,max=64"`
	Hostname        *string `json:"hostname" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	Host            *string `json:"host" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	InternalIP      *string `json:"internal_ip" binding:"omitempty,ip" validate:"omitempty,ip"`
	IPv6Address     *string `json:"ipv6_address" binding:"omitempty,ipv6" validate:"omitempty,ipv6"`
	SSHPort         *int    `json:"ssh_port" binding:"omitempty,min=1,max=65535" validate:"omitempty,min=1,max=65535"`
	SSHCredentialID *string `json:"ssh_credential_id" binding:"omitempty,max=64" validate:"omitempty,max=64"`
	// Direct credential fields for updates (will be stored via credential management)
	SSHUsername     *string `json:"ssh_username" binding:"omitempty,max=64" validate:"omitempty,max=64"`
	SSHPassword     *string `json:"ssh_password" binding:"omitempty" validate:"omitempty"`
	SSHKey          *string `json:"ssh_key" binding:"omitempty" validate:"omitempty"`
	ResourceGroupID *uint   `json:"resource_group_id" binding:"omitempty" validate:"omitempty"`
	Description     *string `json:"description" binding:"omitempty" validate:"omitempty"`
}

// ListServersRequest represents the request to list servers
type ListServersRequest struct {
	common.PaginationRequest
	Keyword         string `form:"keyword" json:"keyword" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	ResourceGroupID *uint  `form:"resource_group_id" json:"resource_group_id" binding:"omitempty" validate:"omitempty"`
	IncludeDeleted  bool   `form:"include_deleted" json:"include_deleted" binding:"omitempty" validate:"omitempty"`
	Status          string `form:"status" json:"status" binding:"omitempty" validate:"omitempty"`
	OS              string `form:"os" json:"os" binding:"omitempty,max=100" validate:"omitempty,max=100"`
	Search          string `form:"search" json:"search" binding:"omitempty,max=255" validate:"omitempty,max=255"`
	SortBy          string `form:"sort_by" json:"sort_by" binding:"omitempty" validate:"omitempty"`
	SortOrder       string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc" validate:"omitempty,oneof=asc desc"`
}

// DeleteServerRequest represents the request to delete a server
type DeleteServerRequest struct {
	Force          bool `form:"force" json:"force" binding:"omitempty" validate:"omitempty"`
	UninstallAgent bool `form:"uninstall_agent" json:"uninstall_agent" binding:"omitempty" validate:"omitempty"`
}

// ServerStatusCheckRequest represents the request to check server status
type ServerStatusCheckRequest struct {
	ServerIDs  []uint   `json:"server_ids" binding:"required,min=1" validate:"required,min=1"`
	CheckTypes []string `json:"check_types" binding:"omitempty" validate:"omitempty,dive,oneof=ssh agent docker"`
	Timeout    int      `json:"timeout" binding:"omitempty,min=5,max=300" validate:"omitempty,min=5,max=300"`
}

// ServerActionRequest represents the request to execute actions on servers (设计文档序号7)
type ServerActionRequest struct {
	ServerIDs []uint                 `json:"server_ids" binding:"required,min=1" validate:"required,min=1"`
	Action    string                 `json:"action" binding:"required" validate:"required,oneof=restart shutdown update_agent run_command"`
	Params    map[string]interface{} `json:"params" binding:"omitempty" validate:"omitempty"`
}

// ServerFileUploadRequest represents the request to upload a file to server
type ServerFileUploadRequest struct {
	FilePath string `json:"file_path" binding:"required" validate:"required"`
	Mode     string `json:"mode" binding:"omitempty" validate:"omitempty"`
	Owner    string `json:"owner" binding:"omitempty" validate:"omitempty"`
}

// ServerFileDownloadRequest represents the request to download a file from server
type ServerFileDownloadRequest struct {
	FilePath string `json:"file_path" binding:"required" validate:"required"`
}

// ServerFileDeleteRequest represents the request to delete a file from server
type ServerFileDeleteRequest struct {
	FilePath string `json:"file_path" binding:"required" validate:"required"`
}

// BatchUpdateServersRequest represents a request to update multiple servers
type BatchUpdateServersRequest struct {
	ServerIDs   []uint  `json:"server_ids" validate:"required,min=1" example:"[1,2,3]"`               // List of server IDs to update
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Updated"` // Optional description update
}

// DeployAgentRequest represents a request to deploy an agent
type DeployAgentRequest struct {
	ServerID       uint   `json:"server_id" validate:"required" example:"1"`                                 // Server ID where to deploy the agent
	DeploymentType string `json:"deployment_type" validate:"required,oneof=docker systemd" example:"docker"` // Deployment type: docker or systemd
}

// UpdateAgentRequest represents a request to update an agent
type UpdateAgentRequest struct {
	Version     *string `json:"version,omitempty" validate:"omitempty,max=50" example:"1.0.1"` // Agent version
	InstallPath *string `json:"install_path,omitempty" validate:"omitempty,max=255"`           // Installation path
	ConfigPath  *string `json:"config_path,omitempty" validate:"omitempty,max=255"`            // Configuration path
	LogPath     *string `json:"log_path,omitempty" validate:"omitempty,max=255"`               // Log path
}

// GetAgentLogsRequest represents a request to get agent logs
type GetAgentLogsRequest struct {
	AgentID uint `json:"agent_id" validate:"required" example:"1"`            // Agent ID
	Lines   int  `json:"lines,omitempty" validate:"omitempty,min=1,max=1000"` // Number of log lines to retrieve
	Follow  bool `json:"follow,omitempty"`                                    // Whether to follow logs (streaming)
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	Name        string `json:"name" validate:"required,max=100" example:"Server Deployment"`               // Task name
	Type        string `json:"type" validate:"required,max=50" example:"deployment"`                       // Task type
	Description string `json:"description,omitempty" validate:"omitempty,max=500" example:"Deploy server"` // Task description
}

// ListTasksRequest represents a request to list tasks
type ListTasksRequest struct {
	common.PaginationRequest
	Status string `form:"status" validate:"omitempty,oneof=pending running completed failed" example:"running"` // Filter by status
	Type   string `form:"type" validate:"omitempty,max=50" example:"deployment"`                                // Filter by type
}

// CreateSecretRequest represents a request to create a secret
type CreateSecretRequest struct {
	Name        string `json:"name" validate:"required,max=100" example:"database_password"`                   // Secret name
	Value       string `json:"value" validate:"required" example:"secret_value"`                               // Secret value
	Type        string `json:"type" validate:"required,max=50" example:"password"`                             // Secret type
	Description string `json:"description,omitempty" validate:"omitempty,max=500" example:"Database password"` // Secret description
}

// UpdateSecretRequest represents a request to update a secret
type UpdateSecretRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=100"`        // Secret name
	Value       *string `json:"value,omitempty"`                                    // Secret value
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"` // Secret description
}
