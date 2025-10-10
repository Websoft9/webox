package service

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// ServerService defines the interface for server business operations
type ServerService interface {
	// Basic CRUD operations
	CreateServer(ctx context.Context, req *request.CreateServerRequest) (*response.ServerResponse, error)
	GetServer(ctx context.Context, id uint) (*response.ServerResponse, error)
	UpdateServer(ctx context.Context, id uint, req *request.UpdateServerRequest) (*response.ServerResponse, error)
	DeleteServer(ctx context.Context, id uint) error
	ListServers(ctx context.Context, req *request.ListServersRequest) (*response.ServerListResponse, error)

	// Status and action operations (按设计文档)
	CheckServersStatus(ctx context.Context, req *request.ServerStatusCheckRequest) (*response.BatchServerStatusResponse, error)
	ExecuteServerActions(ctx context.Context, req *request.ServerActionRequest) (*response.ServerActionResponse, error)

	// File management operations (设计文档序号8-10)
	UploadFile(ctx context.Context, serverID uint, filePath string, fileData []byte) (*response.ServerFileUploadResponse, error)
	DownloadFile(ctx context.Context, serverID uint, filePath string) ([]byte, string, error) // data, filename, error
	DeleteFile(ctx context.Context, serverID uint, filePath string) (*response.ServerFileDeleteResponse, error)
}

// ServerAgentService defines the interface for server agent business operations
type ServerAgentService interface {
	// Basic operations
	DeployAgent(ctx context.Context, req *request.DeployAgentRequest) (*response.ServerAgentResponse, error)
	GetAgent(ctx context.Context, id uint) (*response.ServerAgentResponse, error)
	UpdateAgent(ctx context.Context, id uint, req *request.UpdateAgentRequest) (*response.ServerAgentResponse, error)
	UninstallAgent(ctx context.Context, id uint) error
	RestartAgent(ctx context.Context, id uint) error

	// Status and monitoring operations
	GetAgentStatus(ctx context.Context, id uint) (*response.AgentStatusResponse, error)
	GetAgentLogs(ctx context.Context, req *request.GetAgentLogsRequest) (*response.AgentLogsResponse, error)
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID      string             `json:"task_id"`
	Status      string             `json:"status"`
	Progress    int                `json:"progress"`
	Results     []TaskServerResult `json:"results"`
	CreatedAt   string             `json:"created_at"`
	CompletedAt *string            `json:"completed_at"`
}

// TaskServerResult represents the result of a task execution on a specific server
type TaskServerResult struct {
	ServerID    uint    `json:"server_id"`
	ServerName  string  `json:"server_name"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	ErrorCode   *int    `json:"error_code,omitempty"`
	Output      string  `json:"output,omitempty"`
	CompletedAt *string `json:"completed_at"`
}
