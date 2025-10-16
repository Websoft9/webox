package response

import (
	"time"
)

// ServerResponse represents the server information in API responses
type ServerResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	Hostname        string  `json:"hostname"`
	Host            string  `json:"host"`
	InternalIP      *string `json:"internal_ip"`
	IPv6Address     *string `json:"ipv6_address"`
	SSHPort         int     `json:"ssh_port"`
	SSHCredentialID *string `json:"ssh_credential_id"`
	OSDistro        *string `json:"os_distro"`
	OSVersion       *string `json:"os_version"`
	KernelVersion   *string `json:"kernel_version"`
	CPUCores        int     `json:"cpu_cores"`
	MemoryTotal     int64   `json:"memory_total"`
	DiskTotal       int64   `json:"disk_total"`
	Architecture    *string `json:"architecture"`
	ResourceGroupID *uint   `json:"resource_group_id"`
	OwnerID         uint    `json:"owner_id"`
	Description     *string `json:"description"`
	// Status information (from Redis cache)
	SSHStatus     *string    `json:"ssh_status,omitempty"`
	AgentStatus   *string    `json:"agent_status,omitempty"`
	DockerStatus  *string    `json:"docker_status,omitempty"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	// Related data
	Owner     *UserResponse        `json:"owner,omitempty"`
	Agent     *ServerAgentResponse `json:"agent,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// ServerListResponse represents the paginated server list response
type ServerListResponse struct {
	Items      []ServerResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ServerStatusCheckResult represents the result of server status check (设计文档格式)
type ServerStatusCheckResult struct {
	ServerID     uint                         `json:"server_id"`
	ServerName   string                       `json:"server_name"`
	Status       string                       `json:"status,omitempty"` // success/failed
	ErrorCode    int                          `json:"error_code,omitempty"`
	ErrorMessage string                       `json:"error_message,omitempty"`
	Checks       map[string]ServiceStatusInfo `json:"checks"` // ssh/agent/docker status details
}

// ServiceStatusInfo represents status information for SSH/Agent/Docker services
type ServiceStatusInfo struct {
	Status          string     `json:"status"`                  // connected/online/running etc
	ResponseTime    int        `json:"response_time,omitempty"` // for SSH
	CheckedAt       *time.Time `json:"checked_at"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`   // for Agent
	Version         string     `json:"version,omitempty"`          // for Agent/Docker
	Uptime          int        `json:"uptime,omitempty"`           // for Agent
	ContainersCount int        `json:"containers_count,omitempty"` // for Docker
	ImagesCount     int        `json:"images_count,omitempty"`     // for Docker
	ErrorCode       int        `json:"error_code,omitempty"`
	ErrorMessage    string     `json:"error_message,omitempty"`
}

// ServerServiceStatus represents the status of a specific service
type ServerServiceStatus struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Version     *string                `json:"version,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Message     string                 `json:"message,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// ServerStatusCheckResponse represents the response for server status check
type ServerStatusCheckResponse struct {
	Results      []ServerStatusCheckResult `json:"results"`
	SuccessCount int                       `json:"success_count"`
	FailureCount int                       `json:"failure_count"`
}

// BatchServerStatusResponse represents the response for batch server status check (设计文档序号6)
type BatchServerStatusResponse struct {
	SuccessCount int                       `json:"success_count,omitempty"`
	FailedCount  int                       `json:"failed_count,omitempty"`
	Results      []ServerStatusCheckResult `json:"results"`
}

// ServerActionResponse represents the response for server actions
type ServerActionResponse struct {
	TaskID      string `json:"task_id"`
	Action      string `json:"action"`
	ServerCount int    `json:"server_count"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

// ServerAgentResponse represents agent information in API responses
type ServerAgentResponse struct {
	ID              uint       `json:"id"`
	ServerID        uint       `json:"server_id"`
	AgentID         string     `json:"agent_id"`
	DeploymentType  string     `json:"deployment_type"`
	ContainerID     *string    `json:"container_id"`
	ContainerName   *string    `json:"container_name"`
	ServiceName     *string    `json:"service_name"`
	BinaryPath      *string    `json:"binary_path"`
	ConfigPath      *string    `json:"config_path"`
	AgentIP         *string    `json:"agent_ip"`
	AgentPort       int        `json:"agent_port"`
	Version         *string    `json:"version"`
	PullMode        bool       `json:"pull_mode"`
	PullInterval    int        `json:"pull_interval"`
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at"`
	IsOnline        bool       `json:"is_online"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ServerStatusResponse represents the response for server status check
type ServerStatusResponse struct {
	ServerID  uint      `json:"server_id"`
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
	Message   string    `json:"message"`
}

// BatchOperationResponse represents the response for batch operations
type BatchOperationResponse struct {
	Results      []*OperationResult `json:"results"`
	SuccessCount int                `json:"success_count"`
	FailureCount int                `json:"failure_count"`
}

// OperationResult represents the result of a single operation
type OperationResult struct {
	ServerID uint   `json:"server_id"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// ServerStatisticsResponse represents server statistics
type ServerStatisticsResponse struct {
	Total        int64            `json:"total"`
	StatusCounts map[string]int64 `json:"status_counts"`
}

// AgentStatusResponse represents agent status check response
type AgentStatusResponse struct {
	AgentID   uint      `json:"agent_id"`
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
	Message   string    `json:"message"`
}

// AgentLogsResponse represents agent logs response
type AgentLogsResponse struct {
	AgentID   uint      `json:"agent_id"`
	Logs      []string  `json:"logs"`
	Timestamp time.Time `json:"timestamp"`
}

// TaskResponse represents a task in responses
type TaskResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskListResponse represents a list of tasks
type TaskListResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// SecretResponse represents a secret in responses (without value)
type SecretResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SecretValueResponse represents a secret with its value
type SecretValueResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// ServerFileUploadResponse represents the response for file upload (设计文档序号8)
type ServerFileUploadResponse struct {
	ServerID   uint      `json:"server_id"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	Message    string    `json:"message"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ServerFileDeleteResponse represents the response for file deletion (设计文档序号10)
type ServerFileDeleteResponse struct {
	ServerID  uint      `json:"server_id"`
	FilePath  string    `json:"file_path"`
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
